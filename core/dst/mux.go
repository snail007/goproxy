// Copyright 2014 The DST Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package dst

import (
	"fmt"
	"runtime/debug"

	"net"
	"sync"
	"time"
)

const (
	maxIncomingRequests = 1024
	maxPacketSize       = 500
	handshakeTimeout    = 5 * time.Second
	handshakeInterval   = 1 * time.Second
)

// Mux is a UDP multiplexer of DST connections.
type Mux struct {
	conn       net.PacketConn
	packetSize int

	conns      map[connectionID]*Conn
	handshakes map[connectionID]chan packet
	connsMut   sync.Mutex

	incoming  chan *Conn
	closed    chan struct{}
	closeOnce sync.Once

	buffers *sync.Pool
}

// NewMux creates a new DST Mux on top of a packet connection.
func NewMux(conn net.PacketConn, packetSize int) *Mux {
	if packetSize <= 0 {
		packetSize = maxPacketSize
	}
	m := &Mux{
		conn:       conn,
		packetSize: packetSize,
		conns:      map[connectionID]*Conn{},
		handshakes: make(map[connectionID]chan packet),
		incoming:   make(chan *Conn, maxIncomingRequests),
		closed:     make(chan struct{}),
		buffers: &sync.Pool{
			New: func() interface{} {
				return make([]byte, packetSize)
			},
		},
	}

	// Attempt to maximize buffer space. Start at 16 MB and work downwards 0.5
	// MB at a time.

	if conn, ok := conn.(*net.UDPConn); ok {
		for buf := 16384 * 1024; buf >= 512*1024; buf -= 512 * 1024 {
			err := conn.SetReadBuffer(buf)
			if err == nil {
				if debugMux {
					log.Println(m, "read buffer is", buf)
				}
				break
			}
		}
		for buf := 16384 * 1024; buf >= 512*1024; buf -= 512 * 1024 {
			err := conn.SetWriteBuffer(buf)
			if err == nil {
				if debugMux {
					log.Println(m, "write buffer is", buf)
				}
				break
			}
		}
	}

	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("crashed, err: %s\nstack:%s", e, string(debug.Stack()))
			}
		}()
		m.readerLoop()
	}()
	return m
}

// Accept waits for and returns the next connection to the listener.
func (m *Mux) Accept() (net.Conn, error) {
	return m.AcceptDST()
}

// AcceptDST waits for and returns the next connection to the listener.
func (m *Mux) AcceptDST() (*Conn, error) {
	conn, ok := <-m.incoming
	if !ok {
		return nil, ErrClosedMux
	}
	return conn, nil
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (m *Mux) Close() error {
	var err error = ErrClosedMux
	m.closeOnce.Do(func() {
		err = m.conn.Close()
		close(m.incoming)
		close(m.closed)
	})
	return err
}

// Addr returns the listener's network address.
func (m *Mux) Addr() net.Addr {
	return m.conn.LocalAddr()
}

// Dial connects to the address on the named network.
//
// Network must be "dst".
//
// Addresses have the form host:port. If host is a literal IPv6 address or
// host name, it must be enclosed in square brackets as in "[::1]:80",
// "[ipv6-host]:http" or "[ipv6-host%zone]:80". The functions JoinHostPort and
// SplitHostPort manipulate addresses in this form.
//
// Examples:
//	Dial("dst", "12.34.56.78:80")
//	Dial("dst", "google.com:http")
//	Dial("dst", "[2001:db8::1]:http")
//	Dial("dst", "[fe80::1%lo0]:80")
func (m *Mux) Dial(network, addr string) (net.Conn, error) {
	return m.DialDST(network, addr)
}

// Dial connects to the address on the named network.
//
// Network must be "dst".
//
// Addresses have the form host:port. If host is a literal IPv6 address or
// host name, it must be enclosed in square brackets as in "[::1]:80",
// "[ipv6-host]:http" or "[ipv6-host%zone]:80". The functions JoinHostPort and
// SplitHostPort manipulate addresses in this form.
//
// Examples:
//	Dial("dst", "12.34.56.78:80")
//	Dial("dst", "google.com:http")
//	Dial("dst", "[2001:db8::1]:http")
//	Dial("dst", "[fe80::1%lo0]:80")
func (m *Mux) DialDST(network, addr string) (*Conn, error) {
	if network != "dst" {
		return nil, ErrNotDST
	}

	dst, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	resp := make(chan packet)

	m.connsMut.Lock()
	connID := m.newConnID()
	m.handshakes[connID] = resp
	m.connsMut.Unlock()

	conn, err := m.clientHandshake(dst, connID, resp)

	m.connsMut.Lock()
	defer m.connsMut.Unlock()
	delete(m.handshakes, connID)

	if err != nil {
		return nil, err
	}

	m.conns[connID] = conn
	return conn, nil
}

// handshake performs the client side handshake (i.e. Dial)
func (m *Mux) clientHandshake(dst net.Addr, connID connectionID, resp chan packet) (*Conn, error) {
	if debugMux {
		log.Printf("%v dial %v connID %v", m, dst, connID)
	}

	nextHandshake := time.NewTimer(0)
	defer nextHandshake.Stop()

	handshakeTimeout := time.NewTimer(handshakeTimeout)
	defer handshakeTimeout.Stop()

	var remoteCookie uint32
	seqNo := randomSeqNo()

	for {
		select {
		case <-m.closed:
			// Failure. The mux has been closed.
			return nil, ErrClosedConn

		case <-handshakeTimeout.C:
			// Handshake timeout. Close and abort.
			return nil, ErrHandshakeTimeout

		case <-nextHandshake.C:
			// Send a handshake request.

			m.write(packet{
				src: connID,
				dst: dst,
				hdr: header{
					packetType: typeHandshake,
					flags:      flagRequest,
					connID:     0,
					sequenceNo: seqNo,
					timestamp:  timestampMicros(),
				},
				data: handshakeData{uint32(m.packetSize), connID, remoteCookie}.marshal(),
			})
			nextHandshake.Reset(handshakeInterval)

		case pkt := <-resp:
			hd := unmarshalHandshakeData(pkt.data)

			if pkt.hdr.flags&flagCookie == flagCookie {
				// We should resend the handshake request with a different cookie value.
				remoteCookie = hd.cookie
				nextHandshake.Reset(0)
			} else if pkt.hdr.flags&flagResponse == flagResponse {
				// Successfull handshake response.
				conn := newConn(m, dst)

				conn.connID = connID
				conn.remoteConnID = hd.connID
				conn.nextRecvSeqNo = pkt.hdr.sequenceNo + 1
				conn.packetSize = int(hd.packetSize)
				if conn.packetSize > m.packetSize {
					conn.packetSize = m.packetSize
				}

				conn.nextSeqNo = seqNo + 1

				conn.start()

				return conn, nil
			}
		}
	}
}

func (m *Mux) readerLoop() {
	buf := make([]byte, m.packetSize)
	for {
		buf = buf[:cap(buf)]
		n, from, err := m.conn.ReadFrom(buf)
		if err != nil {
			m.Close()
			return
		}
		buf = buf[:n]

		hdr := unmarshalHeader(buf)

		var bufCopy []byte
		if len(buf) > dstHeaderLen {
			bufCopy = m.buffers.Get().([]byte)[:len(buf)-dstHeaderLen]
			copy(bufCopy, buf[dstHeaderLen:])
		}

		pkt := packet{hdr: hdr, data: bufCopy}
		if debugMux {
			log.Println(m, "read", pkt)
		}

		if hdr.packetType == typeHandshake {
			m.incomingHandshake(from, hdr, bufCopy)
		} else {
			m.connsMut.Lock()
			conn, ok := m.conns[hdr.connID]
			m.connsMut.Unlock()

			if ok {
				conn.in <- packet{
					dst:  nil,
					hdr:  hdr,
					data: bufCopy,
				}
			} else if debugMux && hdr.packetType != typeShutdown {
				log.Printf("packet %v for unknown conn %v", hdr, hdr.connID)
			}
		}
	}
}

func (m *Mux) incomingHandshake(from net.Addr, hdr header, data []byte) {
	if hdr.connID == 0 {
		// A new incoming handshake request.
		m.incomingHandshakeRequest(from, hdr, data)
	} else {
		// A response to an ongoing handshake.
		m.incomingHandshakeResponse(from, hdr, data)
	}
}

func (m *Mux) incomingHandshakeRequest(from net.Addr, hdr header, data []byte) {
	if hdr.flags&flagRequest != flagRequest {
		log.Printf("Handshake pattern with flags 0x%x to connID zero", hdr.flags)
		return
	}

	hd := unmarshalHandshakeData(data)

	correctCookie := cookie(from)
	if hd.cookie != correctCookie {
		// Incorrect or missing SYN cookie. Send back a handshake
		// with the expected one.
		m.write(packet{
			dst: from,
			hdr: header{
				packetType: typeHandshake,
				flags:      flagResponse | flagCookie,
				connID:     hd.connID,
				timestamp:  timestampMicros(),
			},
			data: handshakeData{
				packetSize: uint32(m.packetSize),
				cookie:     correctCookie,
			}.marshal(),
		})
		return
	}

	seqNo := randomSeqNo()

	m.connsMut.Lock()
	connID := m.newConnID()

	conn := newConn(m, from)
	conn.connID = connID
	conn.remoteConnID = hd.connID
	conn.nextSeqNo = seqNo + 1
	conn.nextRecvSeqNo = hdr.sequenceNo + 1
	conn.packetSize = int(hd.packetSize)
	if conn.packetSize > m.packetSize {
		conn.packetSize = m.packetSize
	}
	conn.start()

	m.conns[connID] = conn
	m.connsMut.Unlock()

	m.write(packet{
		dst: from,
		hdr: header{
			packetType: typeHandshake,
			flags:      flagResponse,
			connID:     hd.connID,
			sequenceNo: seqNo,
			timestamp:  timestampMicros(),
		},
		data: handshakeData{
			connID:     conn.connID,
			packetSize: uint32(conn.packetSize),
		}.marshal(),
	})

	m.incoming <- conn
}

func (m *Mux) incomingHandshakeResponse(from net.Addr, hdr header, data []byte) {
	m.connsMut.Lock()
	handShake, ok := m.handshakes[hdr.connID]
	m.connsMut.Unlock()

	if ok {
		// This is a response to a handshake in progress.
		handShake <- packet{
			dst:  nil,
			hdr:  hdr,
			data: data,
		}
	} else if debugMux && hdr.packetType != typeShutdown {
		log.Printf("Handshake packet %v for unknown conn %v", hdr, hdr.connID)
	}
}

func (m *Mux) write(pkt packet) (int, error) {
	buf := m.buffers.Get().([]byte)
	buf = buf[:dstHeaderLen+len(pkt.data)]
	pkt.hdr.marshal(buf)
	copy(buf[dstHeaderLen:], pkt.data)
	if debugMux {
		log.Println(m, "write", pkt)
	}
	n, err := m.conn.WriteTo(buf, pkt.dst)
	m.buffers.Put(buf)
	return n, err
}

func (m *Mux) String() string {
	return fmt.Sprintf("Mux-%v", m.Addr())
}

// Find a unique connection ID
func (m *Mux) newConnID() connectionID {
	for {
		connID := randomConnID()
		if _, ok := m.conns[connID]; ok {
			continue
		}
		if _, ok := m.handshakes[connID]; ok {
			continue
		}
		return connID
	}
}

func (m *Mux) removeConn(c *Conn) {
	m.connsMut.Lock()
	delete(m.conns, c.connID)
	m.connsMut.Unlock()
}

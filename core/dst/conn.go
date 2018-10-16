// Copyright 2014 The DST Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package dst

import (
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"runtime/debug"

	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defExpTime       = 100 * time.Millisecond // N * (4 * RTT + RTTVar + SYN)
	expCountClose    = 8                      // close connection after this many Exps
	minTimeClose     = 5 * time.Second        // if at least this long has passed
	maxInputBuffer   = 8 << 20                // bytes
	muxBufferPackets = 128                    // buffer size of channel between mux and reader routine
	rttMeasureWindow = 32                     // number of packets to track for RTT averaging
	rttMeasureSample = 128                    // Sample every ... packet for RTT

	// number of bytes to subtract from MTU when chunking data, to try to
	// avoid fragmentation
	sliceOverhead = 8 /*pppoe, similar*/ + 20 /*ipv4*/ + 8 /*udp*/ + 16 /*dst*/
)

func init() {
	// Properly seed the random number generator that we use for sequence
	// numbers and stuff.
	buf := make([]byte, 8)
	if n, err := crand.Read(buf); n != 8 || err != nil {
		panic("init random failure")
	}
	rand.Seed(int64(binary.BigEndian.Uint64(buf)))
}

// TODO: export this interface when it's usable from the outside
type congestionController interface {
	Ack()
	NegAck()
	Exp()
	SendWindow() int
	PacketRate() int // PPS
	UpdateRTT(time.Duration)
}

// Conn is an SDT connection carried over a Mux.
type Conn struct {
	// Set at creation, thereafter immutable:

	mux          *Mux
	dst          net.Addr
	connID       connectionID
	remoteConnID connectionID
	in           chan packet
	cc           congestionController
	packetSize   int
	closed       chan struct{}
	closeOnce    sync.Once

	// Touched by more than one goroutine, needs locking.

	nextSeqNoMut sync.Mutex
	nextSeqNo    sequenceNo

	inbufMut  sync.Mutex
	inbufCond *sync.Cond
	inbuf     bytes.Buffer

	expMut sync.Mutex
	exp    *time.Timer

	sendBuffer *sendBuffer // goroutine safe

	packetDelays     [rttMeasureWindow]time.Duration
	packetDelaysSlot int
	packetDelaysMut  sync.Mutex

	// Owned by the reader routine, needs no locking

	recvBuffer        packetList
	nextRecvSeqNo     sequenceNo
	lastAckedSeqNo    sequenceNo
	lastNegAckedSeqNo sequenceNo
	expCount          int
	expReset          time.Time

	// Only accessed atomically

	packetsIn         int64
	packetsOut        int64
	bytesIn           int64
	bytesOut          int64
	resentPackets     int64
	droppedPackets    int64
	outOfOrderPackets int64

	// Special

	debugResetRecvSeqNo chan sequenceNo
}

func newConn(m *Mux, dst net.Addr) *Conn {
	conn := &Conn{
		mux:                 m,
		dst:                 dst,
		nextSeqNo:           sequenceNo(rand.Uint32()),
		packetSize:          maxPacketSize,
		in:                  make(chan packet, muxBufferPackets),
		closed:              make(chan struct{}),
		sendBuffer:          newSendBuffer(m),
		exp:                 time.NewTimer(defExpTime),
		debugResetRecvSeqNo: make(chan sequenceNo),
		expReset:            time.Now(),
	}

	conn.lastAckedSeqNo = conn.nextSeqNo - 1
	conn.inbufCond = sync.NewCond(&conn.inbufMut)

	conn.cc = newWindowCC()
	conn.sendBuffer.SetWindowAndRate(conn.cc.SendWindow(), conn.cc.PacketRate())
	conn.recvBuffer.Resize(128)

	return conn
}

func (c *Conn) start() {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("crashed, err: %s\nstack:\n%s", e, string(debug.Stack()))
			}
		}()
		c.reader()
	}()
}

func (c *Conn) reader() {
	if debugConnection {
		log.Println(c, "reader() starting")
		defer log.Println(c, "reader() exiting")
	}

	for {
		select {
		case <-c.closed:
			// Ack any received but not yet acked messages.
			c.sendAck(0)

			// Send a shutdown message.
			c.nextSeqNoMut.Lock()
			c.mux.write(packet{
				src: c.connID,
				dst: c.dst,
				hdr: header{
					packetType: typeShutdown,
					connID:     c.remoteConnID,
					sequenceNo: c.nextSeqNo,
				},
			})
			c.nextSeqNo++
			c.nextSeqNoMut.Unlock()
			atomic.AddInt64(&c.packetsOut, 1)
			atomic.AddInt64(&c.bytesOut, dstHeaderLen)
			return

		case pkt := <-c.in:
			atomic.AddInt64(&c.packetsIn, 1)
			atomic.AddInt64(&c.bytesIn, dstHeaderLen+int64(len(pkt.data)))

			c.expCount = 1

			switch pkt.hdr.packetType {
			case typeData:
				c.rcvData(pkt)
			case typeAck:
				c.rcvAck(pkt)
			case typeNegAck:
				c.rcvNegAck(pkt)
			case typeShutdown:
				c.rcvShutdown(pkt)
			default:
				log.Println("Unhandled packet", pkt)
				continue
			}

		case <-c.exp.C:
			c.eventExp()
			c.resetExp()

		case n := <-c.debugResetRecvSeqNo:
			// Back door for testing
			c.lastAckedSeqNo = n - 1
			c.nextRecvSeqNo = n
		}
	}
}

func (c *Conn) eventExp() {
	c.expCount++

	if c.sendBuffer.lost.Len() > 0 || c.sendBuffer.send.Len() > 0 {
		c.cc.Exp()
		c.sendBuffer.SetWindowAndRate(c.cc.SendWindow(), c.cc.PacketRate())
		c.sendBuffer.ScheduleResend()

		if debugConnection {
			log.Println(c, "did resends due to Exp")
		}

		if c.expCount > expCountClose && time.Since(c.expReset) > minTimeClose {
			if debugConnection {
				log.Println(c, "close due to Exp")
			}

			// We're shutting down due to repeated exp:s. Don't wait for the
			// send buffer to drain, which it would otherwise do in
			// c.Close()..
			c.sendBuffer.CrashStop()

			c.Close()
		}
	}
}

func (c *Conn) rcvAck(pkt packet) {
	ack := pkt.hdr.sequenceNo

	if debugConnection {
		log.Printf("%v read Ack %v", c, ack)
	}

	c.cc.Ack()

	if ack%rttMeasureSample == 0 {
		if ts := timestamp(binary.BigEndian.Uint32(pkt.data)); ts > 0 {
			if delay := time.Duration(timestampMicros()-ts) * time.Microsecond; delay > 0 {
				c.packetDelaysMut.Lock()
				c.packetDelays[c.packetDelaysSlot] = delay
				c.packetDelaysSlot = (c.packetDelaysSlot + 1) % len(c.packetDelays)
				c.packetDelaysMut.Unlock()

				if rtt, n := c.averageDelay(); n > 8 {
					c.cc.UpdateRTT(rtt)
				}
			}
		}
	}

	c.sendBuffer.Acknowledge(ack)
	c.sendBuffer.SetWindowAndRate(c.cc.SendWindow(), c.cc.PacketRate())

	c.resetExp()
}

func (c *Conn) averageDelay() (time.Duration, int) {
	var total time.Duration
	var n int

	c.packetDelaysMut.Lock()
	for _, d := range c.packetDelays {
		if d != 0 {
			total += d
			n++
		}
	}
	c.packetDelaysMut.Unlock()

	if n == 0 {
		return 0, 0
	}
	return total / time.Duration(n), n
}

func (c *Conn) rcvNegAck(pkt packet) {
	nak := pkt.hdr.sequenceNo

	if debugConnection {
		log.Printf("%v read NegAck %v", c, nak)
	}

	c.sendBuffer.NegativeAck(nak)

	//c.cc.NegAck()
	c.resetExp()
}

func (c *Conn) rcvShutdown(pkt packet) {
	// XXX: We accept shutdown packets somewhat from the future since the
	// sender will number the shutdown after any packets that might still be
	// in the write buffer. This should be fixed to let the write buffer empty
	// on close and reduce the window here.
	if pkt.LessSeq(c.nextRecvSeqNo + 128) {
		if debugConnection {
			log.Println(c, "close due to shutdown")
		}
		c.Close()
	}
}

func (c *Conn) rcvData(pkt packet) {
	if debugConnection {
		log.Println(c, "recv data", pkt.hdr)
	}

	if pkt.LessSeq(c.nextRecvSeqNo) {
		if debugConnection {
			log.Printf("%v old packet received; seq %v, expected %v", c, pkt.hdr.sequenceNo, c.nextRecvSeqNo)
		}
		atomic.AddInt64(&c.droppedPackets, 1)
		return
	}

	if debugConnection {
		log.Println(c, "into recv buffer:", pkt)
	}
	c.recvBuffer.InsertSorted(pkt)
	if c.recvBuffer.LowestSeq() == c.nextRecvSeqNo {
		for _, pkt := range c.recvBuffer.PopSequence(^sequenceNo(0)) {
			if debugConnection {
				log.Println(c, "from recv buffer:", pkt)
			}

			// An in-sequence packet.

			c.nextRecvSeqNo = pkt.hdr.sequenceNo + 1

			c.sendAck(pkt.hdr.timestamp)

			c.inbufMut.Lock()
			for c.inbuf.Len() > len(pkt.data)+maxInputBuffer {
				c.inbufCond.Wait()
				select {
				case <-c.closed:
					return
				default:
				}
			}

			c.inbuf.Write(pkt.data)
			c.inbufCond.Broadcast()
			c.inbufMut.Unlock()
		}
	} else {
		if debugConnection {
			log.Printf("%v lost; seq %v, expected %v", c, pkt.hdr.sequenceNo, c.nextRecvSeqNo)
		}
		c.recvBuffer.InsertSorted(pkt)
		c.sendNegAck()
		atomic.AddInt64(&c.outOfOrderPackets, 1)
	}
}

func (c *Conn) sendAck(ts timestamp) {
	if c.lastAckedSeqNo == c.nextRecvSeqNo {
		return
	}

	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(ts))
	c.mux.write(packet{
		src: c.connID,
		dst: c.dst,
		hdr: header{
			packetType: typeAck,
			connID:     c.remoteConnID,
			sequenceNo: c.nextRecvSeqNo,
		},
		data: buf[:],
	})

	atomic.AddInt64(&c.packetsOut, 1)
	atomic.AddInt64(&c.bytesOut, dstHeaderLen)
	if debugConnection {
		log.Printf("%v send Ack %v", c, c.nextRecvSeqNo)
	}

	c.lastAckedSeqNo = c.nextRecvSeqNo
}

func (c *Conn) sendNegAck() {
	if c.lastNegAckedSeqNo == c.nextRecvSeqNo {
		return
	}

	c.mux.write(packet{
		src: c.connID,
		dst: c.dst,
		hdr: header{
			packetType: typeNegAck,
			connID:     c.remoteConnID,
			sequenceNo: c.nextRecvSeqNo,
		},
	})

	atomic.AddInt64(&c.packetsOut, 1)
	atomic.AddInt64(&c.bytesOut, dstHeaderLen)
	if debugConnection {
		log.Printf("%v send NegAck %v", c, c.nextRecvSeqNo)
	}

	c.lastNegAckedSeqNo = c.nextRecvSeqNo
}

func (c *Conn) resetExp() {
	d, _ := c.averageDelay()
	d = d*4 + 10*time.Millisecond

	if d < defExpTime {
		d = defExpTime
	}

	c.expMut.Lock()
	c.exp.Reset(d)
	c.expMut.Unlock()
}

// String returns a string representation of the connection.
func (c *Conn) String() string {
	return fmt.Sprintf("%v/%v/%v", c.connID, c.LocalAddr(), c.RemoteAddr())
}

// Read reads data from the connection.
// Read can be made to time out and return a Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.
func (c *Conn) Read(b []byte) (n int, err error) {
	defer func() {
		if e := recover(); e != nil {
			n = 0
			err = io.EOF
		}
	}()
	c.inbufMut.Lock()
	defer c.inbufMut.Unlock()
	for c.inbuf.Len() == 0 {
		select {
		case <-c.closed:
			return 0, io.EOF
		default:
		}
		c.inbufCond.Wait()
	}
	return c.inbuf.Read(b)
}

// Write writes data to the connection.
// Write can be made to time out and return a Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline.
func (c *Conn) Write(b []byte) (n int, err error) {
	select {
	case <-c.closed:
		return 0, ErrClosedConn
	default:
	}

	sent := 0
	sliceSize := c.packetSize - sliceOverhead
	for i := 0; i < len(b); i += sliceSize {
		nxt := i + sliceSize
		if nxt > len(b) {
			nxt = len(b)
		}
		slice := b[i:nxt]
		sliceCopy := c.mux.buffers.Get().([]byte)[:len(slice)]
		copy(sliceCopy, slice)

		c.nextSeqNoMut.Lock()
		pkt := packet{
			src: c.connID,
			dst: c.dst,
			hdr: header{
				packetType: typeData,
				sequenceNo: c.nextSeqNo,
				connID:     c.remoteConnID,
			},
			data: sliceCopy,
		}
		c.nextSeqNo++
		c.nextSeqNoMut.Unlock()

		if err := c.sendBuffer.Write(pkt); err != nil {
			return sent, err
		}

		atomic.AddInt64(&c.packetsOut, 1)
		atomic.AddInt64(&c.bytesOut, int64(len(slice)+dstHeaderLen))

		sent += len(slice)
		c.resetExp()
	}
	return sent, nil
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c *Conn) Close() error {
	defer func() {
		_ = recover()
	}()
	c.closeOnce.Do(func() {
		if debugConnection {
			log.Println(c, "explicit close start")
			defer log.Println(c, "explicit close done")
		}

		// XXX: Ugly hack to implement lingering sockets...
		time.Sleep(4 * defExpTime)

		c.sendBuffer.Stop()
		c.mux.removeConn(c)
		close(c.closed)

		c.inbufMut.Lock()
		c.inbufCond.Broadcast()
		c.inbufMut.Unlock()
	})
	return nil
}

// LocalAddr returns the local network address.
func (c *Conn) LocalAddr() net.Addr {
	return c.mux.Addr()
}

// RemoteAddr returns the remote network address.
func (c *Conn) RemoteAddr() net.Addr {
	return c.dst
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail with a timeout (see type Error) instead of
// blocking. The deadline applies to all future I/O, not just
// the immediately following call to Read or Write.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
//
// BUG(jb): SetDeadline is not implemented.
func (c *Conn) SetDeadline(t time.Time) error {
	return ErrNotImplemented
}

// SetReadDeadline sets the deadline for future Read calls.
// A zero value for t means Read will not time out.
//
// BUG(jb): SetReadDeadline is not implemented.
func (c *Conn) SetReadDeadline(t time.Time) error {
	return ErrNotImplemented
}

// SetWriteDeadline sets the deadline for future Write calls.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
//
// BUG(jb): SetWriteDeadline is not implemented.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	return ErrNotImplemented
}

type Statistics struct {
	DataPacketsIn     int64
	DataPacketsOut    int64
	DataBytesIn       int64
	DataBytesOut      int64
	ResentPackets     int64
	DroppedPackets    int64
	OutOfOrderPackets int64
}

// String returns a printable represetnation of the Statistics.
func (s Statistics) String() string {
	return fmt.Sprintf("PktsIn: %d, PktsOut: %d, BytesIn: %d, BytesOut: %d, PktsResent: %d, PktsDropped: %d, PktsOutOfOrder: %d",
		s.DataPacketsIn, s.DataPacketsOut, s.DataBytesIn, s.DataBytesOut, s.ResentPackets, s.DroppedPackets, s.OutOfOrderPackets)
}

// GetStatistics returns a snapsht of the current connection statistics.
func (c *Conn) GetStatistics() Statistics {
	return Statistics{
		DataPacketsIn:     atomic.LoadInt64(&c.packetsIn),
		DataPacketsOut:    atomic.LoadInt64(&c.packetsOut),
		DataBytesIn:       atomic.LoadInt64(&c.bytesIn),
		DataBytesOut:      atomic.LoadInt64(&c.bytesOut),
		ResentPackets:     atomic.LoadInt64(&c.resentPackets),
		DroppedPackets:    atomic.LoadInt64(&c.droppedPackets),
		OutOfOrderPackets: atomic.LoadInt64(&c.outOfOrderPackets),
	}
}

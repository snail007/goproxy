package socks5

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	socks5c "github.com/snail007/goproxy/core/lib/socks5"
)

type Dialer struct {
	timeout          time.Duration
	usernamePassword *socks5c.UsernamePassword
}

// NewDialer returns a new Dialer that dials through the provided
// proxy server's network and address.
func NewDialer(auth *socks5c.UsernamePassword, timeout time.Duration) *Dialer {
	if auth != nil && auth.Password == "" && auth.Username == "" {
		auth = nil
	}
	return &Dialer{
		usernamePassword: auth,
		timeout:          timeout,
	}
}

func (d *Dialer) DialConn(conn *net.Conn, network, addr string) (err error) {
	client := NewClientConn(conn, network, addr, d.timeout, d.usernamePassword, nil)
	err = client._Handshake()
	return
}

type ClientConn struct {
	user     string
	password string
	conn     *net.Conn
	header   []byte
	timeout  time.Duration
	addr     string
	network  string
	udpAddr  string
}

// SOCKS5 returns a Dialer that makes SOCKSv5 connections to the given address
// with an optional username and password. See RFC 1928 and RFC 1929.
// target must be a canonical address with a host and port.
// network : tcp udp
func NewClientConn(conn *net.Conn, network, target string, timeout time.Duration, auth *socks5c.UsernamePassword, header []byte) *ClientConn {
	s := &ClientConn{
		conn:    conn,
		network: network,
		timeout: timeout,
	}
	if auth != nil {
		s.user = auth.Username
		s.password = auth.Password
	}
	if header != nil && len(header) > 0 {
		s.header = header
	}
	if network == "udp" && target == "" {
		target = "0.0.0.0:1"
	}
	s.addr = target
	return s
}

// connect takes an existing connection to a socks5 proxy server,
// and commands the server to extend that connection to target,
// which must be a canonical address with a host and port.
func (s *ClientConn) _Handshake() error {
	host, portStr, err := net.SplitHostPort(s.addr)
	if err != nil {
		return err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return errors.New("proxy: failed to parse port number: " + portStr)
	}
	if port < 1 || port > 0xffff {
		return errors.New("proxy: port number out of range: " + portStr)
	}

	if err := s.auth(host); err != nil {
		return err
	}
	buf := []byte{}
	if s.network == "tcp" {
		buf = append(buf, socks5c.VERSION_V5, socks5c.CMD_CONNECT, 0 /* reserved */)

	} else {
		buf = append(buf, socks5c.VERSION_V5, socks5c.CMD_ASSOCIATE, 0 /* reserved */)
	}
	if ip := net.ParseIP(host); ip != nil {
		if ip4 := ip.To4(); ip4 != nil {
			buf = append(buf, socks5c.ATYP_IPV4)
			ip = ip4
		} else {
			buf = append(buf, socks5c.ATYP_IPV6)
		}
		buf = append(buf, ip...)
	} else {
		if len(host) > 255 {
			return errors.New("proxy: destination host name too long: " + host)
		}
		buf = append(buf, socks5c.ATYP_DOMAIN)
		buf = append(buf, byte(len(host)))
		buf = append(buf, host...)
	}
	buf = append(buf, byte(port>>8), byte(port))
	(*s.conn).SetDeadline(time.Now().Add(s.timeout))
	if _, err := (*s.conn).Write(buf); err != nil {
		return errors.New("proxy: failed to write connect request to SOCKS5 proxy at " + s.addr + ": " + err.Error())
	}
	(*s.conn).SetDeadline(time.Time{})
	(*s.conn).SetDeadline(time.Now().Add(s.timeout))
	if _, err := io.ReadFull((*s.conn), buf[:4]); err != nil {
		return errors.New("proxy: failed to read connect reply from SOCKS5 proxy at " + s.addr + ": " + err.Error())
	}
	(*s.conn).SetDeadline(time.Time{})
	failure := "unknown error"
	if int(buf[1]) < len(socks5c.Socks5Errors) {
		failure = socks5c.Socks5Errors[buf[1]]
	}

	if len(failure) > 0 {
		return errors.New("proxy: SOCKS5 proxy at " + s.addr + " failed to connect: " + failure)
	}

	bytesToDiscard := 0
	switch buf[3] {
	case socks5c.ATYP_IPV4:
		bytesToDiscard = net.IPv4len
	case socks5c.ATYP_IPV6:
		bytesToDiscard = net.IPv6len
	case socks5c.ATYP_DOMAIN:
		(*s.conn).SetDeadline(time.Now().Add(s.timeout))
		_, err := io.ReadFull((*s.conn), buf[:1])
		(*s.conn).SetDeadline(time.Time{})
		if err != nil {
			return errors.New("proxy: failed to read domain length from SOCKS5 proxy at " + s.addr + ": " + err.Error())
		}
		bytesToDiscard = int(buf[0])
	default:
		return errors.New("proxy: got unknown address type " + strconv.Itoa(int(buf[3])) + " from SOCKS5 proxy at " + s.addr)
	}

	if cap(buf) < bytesToDiscard {
		buf = make([]byte, bytesToDiscard)
	} else {
		buf = buf[:bytesToDiscard]
	}
	(*s.conn).SetDeadline(time.Now().Add(s.timeout))
	if _, err := io.ReadFull((*s.conn), buf); err != nil {
		return errors.New("proxy: failed to read address from SOCKS5 proxy at " + s.addr + ": " + err.Error())
	}
	(*s.conn).SetDeadline(time.Time{})
	var ip net.IP
	ip = buf
	ipStr := ""
	if bytesToDiscard == net.IPv4len || bytesToDiscard == net.IPv6len {
		if ipv4 := ip.To4(); ipv4 != nil {
			ipStr = ipv4.String()
		} else {
			ipStr = ip.To16().String()
		}
	}
	//log.Printf("%v", ipStr)
	// Also need to discard the port number
	(*s.conn).SetDeadline(time.Now().Add(s.timeout))
	if _, err := io.ReadFull((*s.conn), buf[:2]); err != nil {
		return errors.New("proxy: failed to read port from SOCKS5 proxy at " + s.addr + ": " + err.Error())
	}
	p := binary.BigEndian.Uint16([]byte{buf[0], buf[1]})
	//log.Printf("%v", p)
	s.udpAddr = net.JoinHostPort(ipStr, fmt.Sprintf("%d", p))
	//log.Printf("%v", s.udpAddr)
	(*s.conn).SetDeadline(time.Time{})
	return nil
}
func (s *ClientConn) SendUDP(data []byte, addr string) (respData []byte, err error) {

	c, err := net.DialTimeout("udp", s.udpAddr, s.timeout)
	if err != nil {
		return
	}
	conn := c.(*net.UDPConn)

	p := socks5c.NewPacketUDP()
	p.Build(addr, data)
	conn.SetDeadline(time.Now().Add(s.timeout))
	conn.Write(p.Bytes())
	conn.SetDeadline(time.Time{})

	buf := make([]byte, 1024)
	conn.SetDeadline(time.Now().Add(s.timeout))
	n, _, err := conn.ReadFrom(buf)
	conn.SetDeadline(time.Time{})
	if err != nil {
		return
	}
	respData = buf[:n]
	return
}
func (s *ClientConn) auth(host string) error {

	// the size here is just an estimate
	buf := make([]byte, 0, 6+len(host))

	buf = append(buf, socks5c.VERSION_V5)
	if len(s.user) > 0 && len(s.user) < 256 && len(s.password) < 256 {
		buf = append(buf, 2 /* num auth methods */, socks5c.Method_NO_AUTH, socks5c.Method_USER_PASS)
	} else {
		buf = append(buf, 1 /* num auth methods */, socks5c.Method_NO_AUTH)
	}
	(*s.conn).SetDeadline(time.Now().Add(s.timeout))
	if _, err := (*s.conn).Write(buf); err != nil {
		return errors.New("proxy: failed to write greeting to SOCKS5 proxy at " + s.addr + ": " + err.Error())
	}
	(*s.conn).SetDeadline(time.Time{})

	(*s.conn).SetDeadline(time.Now().Add(s.timeout))
	if _, err := io.ReadFull((*s.conn), buf[:2]); err != nil {
		return errors.New("proxy: failed to read greeting from SOCKS5 proxy at " + s.addr + ": " + err.Error())
	}
	(*s.conn).SetDeadline(time.Time{})

	if buf[0] != 5 {
		return errors.New("proxy: SOCKS5 proxy at " + s.addr + " has unexpected version " + strconv.Itoa(int(buf[0])))
	}
	if buf[1] == 0xff {
		return errors.New("proxy: SOCKS5 proxy at " + s.addr + " requires authentication")
	}

	// See RFC 1929
	if buf[1] == socks5c.Method_USER_PASS {
		buf = buf[:0]
		buf = append(buf, 1 /* password protocol version */)
		buf = append(buf, uint8(len(s.user)))
		buf = append(buf, s.user...)
		buf = append(buf, uint8(len(s.password)))
		buf = append(buf, s.password...)
		(*s.conn).SetDeadline(time.Now().Add(s.timeout))
		if _, err := (*s.conn).Write(buf); err != nil {
			return errors.New("proxy: failed to write authentication request to SOCKS5 proxy at " + s.addr + ": " + err.Error())
		}
		(*s.conn).SetDeadline(time.Time{})
		(*s.conn).SetDeadline(time.Now().Add(s.timeout))
		if _, err := io.ReadFull((*s.conn), buf[:2]); err != nil {
			return errors.New("proxy: failed to read authentication reply from SOCKS5 proxy at " + s.addr + ": " + err.Error())
		}
		(*s.conn).SetDeadline(time.Time{})
		if buf[1] != 0 {
			return errors.New("proxy: SOCKS5 proxy at " + s.addr + " rejected username/password")
		}
	}
	return nil
}

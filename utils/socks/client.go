package socks

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

var socks5Errors = []string{
	"",
	"general failure",
	"connection forbidden",
	"network unreachable",
	"host unreachable",
	"connection refused",
	"TTL expired",
	"command not supported",
	"address type not supported",
}

type Auth struct {
	User, Password string
}
type ClientConn struct {
	user     string
	password string
	conn     *net.Conn
	header   []byte
	timeout  time.Duration
	addr     string
	network  string
	UDPAddr  string
}

// SOCKS5 returns a Dialer that makes SOCKSv5 connections to the given address
// with an optional username and password. See RFC 1928 and RFC 1929.
// target must be a canonical address with a host and port.
// network : tcp udp
func NewClientConn(conn *net.Conn, network, target string, timeout time.Duration, auth *Auth, header []byte) *ClientConn {
	s := &ClientConn{
		conn:    conn,
		network: network,
		timeout: timeout,
	}
	if auth != nil {
		s.user = auth.User
		s.password = auth.Password
	}
	if header != nil && len(header) > 0 {
		s.header = header
	}
	if network == "udp" && target == "" {
		target = "0.0.0.0:0"
	}
	s.addr = target
	return s
}

// connect takes an existing connection to a socks5 proxy server,
// and commands the server to extend that connection to target,
// which must be a canonical address with a host and port.
func (s *ClientConn) Handshake() error {
	host, portStr, err := net.SplitHostPort(s.addr)
	if err != nil {
		return err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return errors.New("proxy: failed to parse port number: " + portStr)
	}
	if s.network == "tcp" && (port < 1 || port > 0xffff) {
		return errors.New("proxy: port number out of range: " + portStr)
	}

	if err := s.handshake(host); err != nil {
		return err
	}
	buf := []byte{}
	if s.network == "tcp" {
		buf = append(buf, VERSION_V5, CMD_CONNECT, 0 /* reserved */)

	} else {
		buf = append(buf, VERSION_V5, CMD_ASSOCIATE, 0 /* reserved */)
	}
	if ip := net.ParseIP(host); ip != nil {
		if ip4 := ip.To4(); ip4 != nil {
			buf = append(buf, ATYP_IPV4)
			ip = ip4
		} else {
			buf = append(buf, ATYP_IPV6)
		}
		buf = append(buf, ip...)
	} else {
		if len(host) > 255 {
			return errors.New("proxy: destination host name too long: " + host)
		}
		buf = append(buf, ATYP_DOMAIN)
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
	if int(buf[1]) < len(socks5Errors) {
		failure = socks5Errors[buf[1]]
	}

	if len(failure) > 0 {
		return errors.New("proxy: SOCKS5 proxy at " + s.addr + " failed to connect: " + failure)
	}

	bytesToDiscard := 0
	switch buf[3] {
	case ATYP_IPV4:
		bytesToDiscard = net.IPv4len
	case ATYP_IPV6:
		bytesToDiscard = net.IPv6len
	case ATYP_DOMAIN:
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
	s.UDPAddr = net.JoinHostPort(ipStr, fmt.Sprintf("%d", p))
	//log.Printf("%v", s.udpAddr)
	(*s.conn).SetDeadline(time.Time{})
	return nil
}
func (s *ClientConn) SendUDP(data []byte, addr string) (respData []byte, err error) {

	c, err := net.DialTimeout("udp", s.UDPAddr, s.timeout)
	if err != nil {
		return
	}
	conn := c.(*net.UDPConn)

	p := NewPacketUDP()
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
func (s *ClientConn) handshake(host string) error {

	// the size here is just an estimate
	buf := make([]byte, 0, 6+len(host))

	buf = append(buf, VERSION_V5)
	if len(s.user) > 0 && len(s.user) < 256 && len(s.password) < 256 {
		buf = append(buf, 2 /* num auth methods */, Method_NO_AUTH, Method_USER_PASS)
	} else {
		buf = append(buf, 1 /* num auth methods */, Method_NO_AUTH)
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
	if buf[1] == Method_USER_PASS {
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

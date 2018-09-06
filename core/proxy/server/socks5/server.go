package socks5

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	socks5c "github.com/snail007/goproxy/core/lib/socks5"
)

type BasicAuther interface {
	CheckUserPass(username, password, fromIP, ToTarget string) bool
}
type Request struct {
	ver         uint8
	cmd         uint8
	reserve     uint8
	addressType uint8
	dstAddr     string
	dstPort     string
	dstHost     string
	bytes       []byte
	rw          io.ReadWriter
}

func NewRequest(rw io.ReadWriter, header ...[]byte) (req Request, err interface{}) {
	var b = make([]byte, 1024)
	var n int
	req = Request{rw: rw}
	if header != nil && len(header) == 1 && len(header[0]) > 1 {
		b = header[0]
		n = len(header[0])
	} else {
		n, err = rw.Read(b[:])
		if err != nil {
			err = fmt.Errorf("read req data fail,ERR: %s", err)
			return
		}
	}
	req.ver = uint8(b[0])
	req.cmd = uint8(b[1])
	req.reserve = uint8(b[2])
	req.addressType = uint8(b[3])
	if b[0] != 0x5 {
		err = fmt.Errorf("sosck version supported")
		req.TCPReply(socks5c.REP_REQ_FAIL)
		return
	}
	switch b[3] {
	case 0x01: //IP V4
		req.dstHost = net.IPv4(b[4], b[5], b[6], b[7]).String()
	case 0x03: //域名
		req.dstHost = string(b[5 : n-2]) //b[4]表示域名的长度
	case 0x04: //IP V6
		req.dstHost = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
	}
	req.dstPort = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))
	req.dstAddr = net.JoinHostPort(req.dstHost, req.dstPort)
	req.bytes = b[:n]
	return
}
func (s *Request) Bytes() []byte {
	return s.bytes
}
func (s *Request) Addr() string {
	return s.dstAddr
}
func (s *Request) Host() string {
	return s.dstHost
}
func (s *Request) Port() string {
	return s.dstPort
}
func (s *Request) AType() uint8 {
	return s.addressType
}
func (s *Request) CMD() uint8 {
	return s.cmd
}

func (s *Request) TCPReply(rep uint8) (err error) {
	_, err = s.rw.Write(s.NewReply(rep, "0.0.0.0:0"))
	return
}
func (s *Request) UDPReply(rep uint8, addr string) (err error) {
	_, err = s.rw.Write(s.NewReply(rep, addr))
	return
}
func (s *Request) NewReply(rep uint8, addr string) []byte {
	var response bytes.Buffer
	host, port, _ := net.SplitHostPort(addr)
	ip := net.ParseIP(host)
	ipb := ip.To4()
	atyp := socks5c.ATYP_IPV4
	ipv6 := ip.To16()
	zeroiIPv6 := fmt.Sprintf("%d%d%d%d%d%d%d%d%d%d%d%d",
		ipv6[0], ipv6[1], ipv6[2], ipv6[3],
		ipv6[4], ipv6[5], ipv6[6], ipv6[7],
		ipv6[8], ipv6[9], ipv6[10], ipv6[11],
	)
	if ipb == nil && ipv6 != nil && "0000000000255255" != zeroiIPv6 {
		atyp = socks5c.ATYP_IPV6
		ipb = ip.To16()
	}
	porti, _ := strconv.Atoi(port)
	portb := make([]byte, 2)
	binary.BigEndian.PutUint16(portb, uint16(porti))
	// log.Printf("atyp : %v", atyp)
	// log.Printf("ip : %v", []byte(ip))
	response.WriteByte(socks5c.VERSION_V5)
	response.WriteByte(rep)
	response.WriteByte(socks5c.RSV)
	response.WriteByte(atyp)
	response.Write(ipb)
	response.Write(portb)
	return response.Bytes()
}

type MethodsRequest struct {
	ver          uint8
	methodsCount uint8
	methods      []uint8
	bytes        []byte
	rw           *io.ReadWriter
}

func NewMethodsRequest(r io.ReadWriter, header ...[]byte) (s MethodsRequest, err interface{}) {
	defer func() {
		if err == nil {
			err = recover()
		}
	}()
	s = MethodsRequest{}
	s.rw = &r
	var buf = make([]byte, 300)
	var n int
	if header != nil && len(header) == 1 && len(header[0]) > 1 {
		buf = header[0]
		n = len(header[0])
	} else {
		n, err = r.Read(buf)
		if err != nil {
			return
		}
	}
	if buf[0] != 0x05 {
		err = fmt.Errorf("socks version not supported")
		return
	}
	if n != int(buf[1])+int(2) {
		err = fmt.Errorf("socks methods data length error")
		return
	}
	s.ver = buf[0]
	s.methodsCount = buf[1]
	s.methods = buf[2:n]
	s.bytes = buf[:n]
	return
}
func (s *MethodsRequest) Version() uint8 {
	return s.ver
}
func (s *MethodsRequest) MethodsCount() uint8 {
	return s.methodsCount
}
func (s *MethodsRequest) Methods() []uint8 {
	return s.methods
}
func (s *MethodsRequest) Select(method uint8) bool {
	for _, m := range s.methods {
		if m == method {
			return true
		}
	}
	return false
}
func (s *MethodsRequest) Reply(method uint8) (err error) {
	_, err = (*s.rw).Write([]byte{byte(socks5c.VERSION_V5), byte(method)})
	return
}
func (s *MethodsRequest) Bytes() []byte {
	return s.bytes
}

type ServerConn struct {
	target   string
	user     string
	password string
	conn     *net.Conn
	timeout  time.Duration
	auth     *BasicAuther
	header   []byte
	ver      uint8
	//method
	methodsCount uint8
	methods      []uint8
	method       uint8
	//request
	cmd         uint8
	reserve     uint8
	addressType uint8
	dstAddr     string
	dstPort     string
	dstHost     string
	udpAddress  string
}

func NewServerConn(conn *net.Conn, timeout time.Duration, auth *BasicAuther, udpAddress string, header []byte) *ServerConn {
	if udpAddress == "" {
		udpAddress = "0.0.0.0:16666"
	}
	s := &ServerConn{
		conn:       conn,
		timeout:    timeout,
		auth:       auth,
		header:     header,
		ver:        socks5c.VERSION_V5,
		udpAddress: udpAddress,
	}
	return s

}
func (s *ServerConn) Close() {
	(*s.conn).Close()
}
func (s *ServerConn) AuthData() socks5c.UsernamePassword {
	return socks5c.UsernamePassword{s.user, s.password}
}
func (s *ServerConn) Method() uint8 {
	return s.method
}
func (s *ServerConn) Target() string {
	return s.target
}
func (s *ServerConn) Handshake() (err error) {
	remoteAddr := (*s.conn).RemoteAddr()
	//协商开始
	//method select request
	var methodReq MethodsRequest
	(*s.conn).SetReadDeadline(time.Now().Add(time.Second * s.timeout))

	methodReq, e := NewMethodsRequest((*s.conn), s.header)
	(*s.conn).SetReadDeadline(time.Time{})
	if e != nil {
		(*s.conn).SetReadDeadline(time.Now().Add(time.Second * s.timeout))
		methodReq.Reply(socks5c.Method_NONE_ACCEPTABLE)
		(*s.conn).SetReadDeadline(time.Time{})
		err = fmt.Errorf("new methods request fail,ERR: %s", e)
		return
	}
	//log.Printf("%v,s.auth == %v && methodReq.Select(Method_NO_AUTH) %v", methodReq.methods, s.auth, methodReq.Select(Method_NO_AUTH))
	if s.auth == nil && methodReq.Select(socks5c.Method_NO_AUTH) && !methodReq.Select(socks5c.Method_USER_PASS) {
		// if !methodReq.Select(Method_NO_AUTH) {
		// 	(*s.conn).SetReadDeadline(time.Now().Add(time.Second * s.timeout))
		// 	methodReq.Reply(Method_NONE_ACCEPTABLE)
		// 	(*s.conn).SetReadDeadline(time.Time{})
		// 	err = fmt.Errorf("none method found : Method_NO_AUTH")
		// 	return
		// }
		s.method = socks5c.Method_NO_AUTH
		//method select reply
		(*s.conn).SetReadDeadline(time.Now().Add(time.Second * s.timeout))
		err = methodReq.Reply(socks5c.Method_NO_AUTH)
		(*s.conn).SetReadDeadline(time.Time{})
		if err != nil {
			err = fmt.Errorf("reply answer data fail,ERR: %s", err)
			return
		}
		// err = fmt.Errorf("% x", methodReq.Bytes())
	} else {
		//auth
		if !methodReq.Select(socks5c.Method_USER_PASS) {
			(*s.conn).SetReadDeadline(time.Now().Add(time.Second * s.timeout))
			methodReq.Reply(socks5c.Method_NONE_ACCEPTABLE)
			(*s.conn).SetReadDeadline(time.Time{})
			err = fmt.Errorf("none method found : Method_USER_PASS")
			return
		}
		s.method = socks5c.Method_USER_PASS
		//method reply need auth
		(*s.conn).SetReadDeadline(time.Now().Add(time.Second * s.timeout))
		err = methodReq.Reply(socks5c.Method_USER_PASS)
		(*s.conn).SetReadDeadline(time.Time{})
		if err != nil {
			err = fmt.Errorf("reply answer data fail,ERR: %s", err)
			return
		}
		//read auth
		buf := make([]byte, 500)
		var n int
		(*s.conn).SetReadDeadline(time.Now().Add(time.Second * s.timeout))
		n, err = (*s.conn).Read(buf)
		(*s.conn).SetReadDeadline(time.Time{})
		if err != nil {
			err = fmt.Errorf("read auth info fail,ERR: %s", err)
			return
		}
		r := buf[:n]
		s.user = string(r[2 : r[1]+2])
		s.password = string(r[2+r[1]+1:])
		//err = fmt.Errorf("user:%s,pass:%s", user, pass)
		//auth
		_addr := strings.Split(remoteAddr.String(), ":")
		if s.auth == nil || (*s.auth).CheckUserPass(s.user, s.password, _addr[0], "") {
			(*s.conn).SetDeadline(time.Now().Add(time.Millisecond * time.Duration(s.timeout)))
			_, err = (*s.conn).Write([]byte{0x01, 0x00})
			(*s.conn).SetDeadline(time.Time{})
			if err != nil {
				err = fmt.Errorf("answer auth success to %s fail,ERR: %s", remoteAddr, err)
				return
			}
		} else {
			(*s.conn).SetDeadline(time.Now().Add(time.Millisecond * time.Duration(s.timeout)))
			_, err = (*s.conn).Write([]byte{0x01, 0x01})
			(*s.conn).SetDeadline(time.Time{})
			if err != nil {
				err = fmt.Errorf("answer auth fail to %s fail,ERR: %s", remoteAddr, err)
				return
			}
			err = fmt.Errorf("auth fail from %s", remoteAddr)
			return
		}
	}
	//request detail
	(*s.conn).SetReadDeadline(time.Now().Add(time.Second * s.timeout))
	request, e := NewRequest(*s.conn)
	(*s.conn).SetReadDeadline(time.Time{})
	if e != nil {
		err = fmt.Errorf("read request data fail,ERR: %s", e)
		return
	}
	//协商结束

	switch request.CMD() {
	case socks5c.CMD_BIND:
		err = request.TCPReply(socks5c.REP_UNKNOWN)
		if err != nil {
			err = fmt.Errorf("TCPReply REP_UNKNOWN to %s fail,ERR: %s", remoteAddr, err)
			return
		}
		err = fmt.Errorf("cmd bind not supported, form: %s", remoteAddr)
		return
	case socks5c.CMD_CONNECT:
		err = request.TCPReply(socks5c.REP_SUCCESS)
		if err != nil {
			err = fmt.Errorf("TCPReply REP_SUCCESS to %s fail,ERR: %s", remoteAddr, err)
			return
		}
	case socks5c.CMD_ASSOCIATE:
		err = request.UDPReply(socks5c.REP_SUCCESS, s.udpAddress)
		if err != nil {
			err = fmt.Errorf("UDPReply REP_SUCCESS to %s fail,ERR: %s", remoteAddr, err)
			return
		}
	}

	//fill socks info
	s.target = request.Addr()
	s.methodsCount = methodReq.MethodsCount()
	s.methods = methodReq.Methods()
	s.cmd = request.CMD()
	s.reserve = request.reserve
	s.addressType = request.addressType
	s.dstAddr = request.dstAddr
	s.dstHost = request.dstHost
	s.dstPort = request.dstPort
	return
}

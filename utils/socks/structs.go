package socks

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
)

const (
	Method_NO_AUTH         = uint8(0x00)
	Method_GSSAPI          = uint8(0x01)
	Method_USER_PASS       = uint8(0x02)
	Method_IANA            = uint8(0x7F)
	Method_RESVERVE        = uint8(0x80)
	Method_NONE_ACCEPTABLE = uint8(0xFF)
	VERSION_V5             = uint8(0x05)
	CMD_CONNECT            = uint8(0x01)
	CMD_BIND               = uint8(0x02)
	CMD_ASSOCIATE          = uint8(0x03)
	ATYP_IPV4              = uint8(0x01)
	ATYP_DOMAIN            = uint8(0x03)
	ATYP_IPV6              = uint8(0x04)
	REP_SUCCESS            = uint8(0x00)
	REP_REQ_FAIL           = uint8(0x01)
	REP_RULE_FORBIDDEN     = uint8(0x02)
	REP_NETWOR_UNREACHABLE = uint8(0x03)
	REP_HOST_UNREACHABLE   = uint8(0x04)
	REP_CONNECTION_REFUSED = uint8(0x05)
	REP_TTL_TIMEOUT        = uint8(0x06)
	REP_CMD_UNSUPPORTED    = uint8(0x07)
	REP_ATYP_UNSUPPORTED   = uint8(0x08)
	REP_UNKNOWN            = uint8(0x09)
	RSV                    = uint8(0x00)
)

var (
	ZERO_IP   = []byte{0x00, 0x00, 0x00, 0x00}
	ZERO_PORT = []byte{0x00, 0x00}
)

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

func NewRequest(rw io.ReadWriter) (req Request, err interface{}) {
	var b [1024]byte
	var n int
	req = Request{rw: rw}
	n, err = rw.Read(b[:])
	if err != nil {
		err = fmt.Errorf("read req data fail,ERR: %s", err)
		return
	}
	req.ver = uint8(b[0])
	req.cmd = uint8(b[1])
	req.reserve = uint8(b[2])
	req.addressType = uint8(b[3])

	if b[0] != 0x5 {
		err = fmt.Errorf("sosck version supported")
		req.TCPReply(REP_REQ_FAIL)
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
	atyp := ATYP_IPV4
	ipv6 := ip.To16()
	zeroiIPv6 := fmt.Sprintf("%d%d%d%d%d%d%d%d%d%d%d%d",
		ipv6[0], ipv6[1], ipv6[2], ipv6[3],
		ipv6[4], ipv6[5], ipv6[6], ipv6[7],
		ipv6[8], ipv6[9], ipv6[10], ipv6[11],
	)
	if ipv6 != nil && "0000000000255255" != zeroiIPv6 {
		atyp = ATYP_IPV6
		ipb = ip.To16()
	}
	porti, _ := strconv.Atoi(port)
	portb := make([]byte, 2)
	binary.BigEndian.PutUint16(portb, uint16(porti))
	// log.Printf("atyp : %v", atyp)
	// log.Printf("ip : %v", []byte(ip))
	response.WriteByte(VERSION_V5)
	response.WriteByte(rep)
	response.WriteByte(RSV)
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

func NewMethodsRequest(r io.ReadWriter) (s MethodsRequest, err interface{}) {
	defer func() {
		if err == nil {
			err = recover()
		}
	}()
	s = MethodsRequest{}
	s.rw = &r
	var buf = make([]byte, 300)
	var n int
	n, err = r.Read(buf)
	if err != nil {
		return
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
func (s *MethodsRequest) Select(method uint8) bool {
	for _, m := range s.methods {
		if m == method {
			return true
		}
	}
	return false
}
func (s *MethodsRequest) Reply(method uint8) (err error) {
	_, err = (*s.rw).Write([]byte{byte(VERSION_V5), byte(method)})
	return
}
func (s *MethodsRequest) Bytes() []byte {
	return s.bytes
}

type UDPPacket struct {
	rsv     uint16
	frag    uint8
	atype   uint8
	dstHost string
	dstPort string
	data    []byte
	header  []byte
	bytes   []byte
}

func ParseUDPPacket(b []byte) (p UDPPacket, err error) {
	p = UDPPacket{}
	p.frag = uint8(b[2])
	p.bytes = b
	if p.frag != 0 {
		err = fmt.Errorf("FRAG only support for 0 , %v ,%v", p.frag, b[:4])
		return
	}
	portIndex := 0
	p.atype = b[3]
	switch p.atype {
	case ATYP_IPV4: //IP V4
		p.dstHost = net.IPv4(b[4], b[5], b[6], b[7]).String()
		portIndex = 8
	case ATYP_DOMAIN: //域名
		domainLen := uint8(b[4])
		p.dstHost = string(b[5 : 5+domainLen]) //b[4]表示域名的长度
		portIndex = int(5 + domainLen)
	case ATYP_IPV6: //IP V6
		p.dstHost = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		portIndex = 20
	}
	p.dstPort = strconv.Itoa(int(b[portIndex])<<8 | int(b[portIndex+1]))
	p.data = b[portIndex+2:]
	p.header = b[:portIndex+2]
	return
}
func (s *UDPPacket) Header() []byte {
	return s.header
}
func (s *UDPPacket) Host() string {
	return s.dstHost
}

func (s *UDPPacket) Port() string {
	return s.dstPort
}
func (s *UDPPacket) Data() []byte {
	return s.data
}

package socks

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
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
	if ipb == nil && ipv6 != nil && "0000000000255255" != zeroiIPv6 {
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
	_, err = (*s.rw).Write([]byte{byte(VERSION_V5), byte(method)})
	return
}
func (s *MethodsRequest) Bytes() []byte {
	return s.bytes
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

func (s *UDPPacket) Header() []byte {
	return s.header
}
func (s *UDPPacket) NewReply(data []byte) []byte {
	var buf bytes.Buffer
	buf.Write(s.header)
	buf.Write(data)
	return buf.Bytes()
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

type PacketUDP struct {
	rsv     uint16
	frag    uint8
	atype   uint8
	dstHost string
	dstPort string
	data    []byte
}

func NewPacketUDP() (p PacketUDP) {
	return PacketUDP{}
}
func (p *PacketUDP) Build(destAddr string, data []byte) (err error) {
	host, port, err := net.SplitHostPort(destAddr)
	if err != nil {
		return
	}
	p.rsv = 0
	p.frag = 0
	p.dstHost = host
	p.dstPort = port
	p.atype = ATYP_IPV4
	if ip := net.ParseIP(host); ip != nil {
		if ip4 := ip.To4(); ip4 != nil {
			p.atype = ATYP_IPV4
			ip = ip4
		} else {
			p.atype = ATYP_IPV6
		}
	} else {
		if len(host) > 255 {
			err = errors.New("proxy: destination host name too long: " + host)
			return
		}
		p.atype = ATYP_DOMAIN
	}
	p.data = data

	return
}
func (p *PacketUDP) Parse(b []byte) (err error) {
	if len(b) < 9 {
		return fmt.Errorf("too short packet")
	}
	p.frag = uint8(b[2])
	if p.frag != 0 {
		err = fmt.Errorf("FRAG only support for 0 , %v ,%v", p.frag, b[:4])
		return
	}
	portIndex := 0
	p.atype = b[3]
	switch p.atype {
	case ATYP_IPV4: //IP V4
		if len(b) < 11 {
			return fmt.Errorf("too short packet")
		}
		p.dstHost = net.IPv4(b[4], b[5], b[6], b[7]).String()
		portIndex = 8
	case ATYP_DOMAIN: //域名
		domainLen := uint8(b[4])
		if len(b) < int(domainLen)+7 {
			return fmt.Errorf("too short packet")
		}
		p.dstHost = string(b[5 : 5+domainLen]) //b[4]表示域名的长度
		portIndex = int(5 + domainLen)
	case ATYP_IPV6: //IP V6
		if len(b) < 22 {
			return fmt.Errorf("too short packet")
		}
		p.dstHost = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		portIndex = 20
	}
	p.dstPort = strconv.Itoa(int(b[portIndex])<<8 | int(b[portIndex+1]))
	p.data = b[portIndex+2:]
	return
}
func (p *PacketUDP) Header() []byte {
	header := new(bytes.Buffer)
	header.Write([]byte{0x00, 0x00, p.frag, p.atype})
	if p.atype == ATYP_IPV4 {
		ip := net.ParseIP(p.dstHost)
		header.Write(ip.To4())
	} else if p.atype == ATYP_IPV6 {
		ip := net.ParseIP(p.dstHost)
		header.Write(ip.To16())
	} else if p.atype == ATYP_DOMAIN {
		hBytes := []byte(p.dstHost)
		header.WriteByte(byte(len(hBytes)))
		header.Write(hBytes)
	}
	port, _ := strconv.ParseUint(p.dstPort, 10, 64)
	portBytes := new(bytes.Buffer)
	binary.Write(portBytes, binary.BigEndian, port)
	header.Write(portBytes.Bytes()[portBytes.Len()-2:])
	return header.Bytes()
}
func (p *PacketUDP) Bytes() []byte {
	packBytes := new(bytes.Buffer)
	packBytes.Write(p.Header())
	packBytes.Write(p.data)
	return packBytes.Bytes()
}
func (p *PacketUDP) Host() string {
	return p.dstHost
}
func (p *PacketUDP) Addr() string {
	return net.JoinHostPort(p.dstHost, p.dstPort)
}
func (p *PacketUDP) Port() string {
	return p.dstPort
}
func (p *PacketUDP) Data() []byte {
	return p.data
}

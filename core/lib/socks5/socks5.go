package socks5

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
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
var Socks5Errors = []string{
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

// Auth contains authentication parameters that specific Dialers may require.
type UsernamePassword struct {
	Username, Password string
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
	p.frag = uint8(b[2])
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

func (p *PacketUDP) Port() string {
	return p.dstPort
}
func (p *PacketUDP) Data() []byte {
	return p.data
}

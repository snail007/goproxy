package socks

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/snail007/goproxy/utils"
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

type ServerConn struct {
	target   string
	user     string
	password string
	conn     *net.Conn
	timeout  time.Duration
	auth     *utils.BasicAuth
	header   []byte
	ver      uint8
	//method
	methodsCount uint8
	methods      []uint8
	method       uint8
	//request
	cmd             uint8
	reserve         uint8
	addressType     uint8
	dstAddr         string
	dstPort         string
	dstHost         string
	UDPConnListener *net.UDPConn
	enableUDP       bool
	udpIP           string
}

func NewServerConn(conn *net.Conn, timeout time.Duration, auth *utils.BasicAuth, enableUDP bool, udpHost string, header []byte) *ServerConn {

	s := &ServerConn{
		conn:      conn,
		timeout:   timeout,
		auth:      auth,
		header:    header,
		ver:       VERSION_V5,
		enableUDP: enableUDP,
		udpIP:     udpHost,
	}
	return s

}
func (s *ServerConn) Close() {
	utils.CloseConn(s.conn)
}
func (s *ServerConn) AuthData() Auth {
	return Auth{s.user, s.password}
}
func (s *ServerConn) IsUDP() bool {
	return s.cmd == CMD_ASSOCIATE
}
func (s *ServerConn) IsTCP() bool {
	return s.cmd == CMD_CONNECT
}
func (s *ServerConn) Method() uint8 {
	return s.method
}
func (s *ServerConn) Target() string {
	return s.target
}
func (s *ServerConn) Host() string {
	return s.dstHost
}
func (s *ServerConn) Port() string {
	return s.dstPort
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
		methodReq.Reply(Method_NONE_ACCEPTABLE)
		(*s.conn).SetReadDeadline(time.Time{})
		err = fmt.Errorf("new methods request fail,ERR: %s", e)
		return
	}
	//log.Printf("%v,s.auth == %v && methodReq.Select(Method_NO_AUTH) %v", methodReq.methods, s.auth, methodReq.Select(Method_NO_AUTH))
	if s.auth == nil && methodReq.Select(Method_NO_AUTH) && !methodReq.Select(Method_USER_PASS) {
		// if !methodReq.Select(Method_NO_AUTH) {
		// 	(*s.conn).SetReadDeadline(time.Now().Add(time.Second * s.timeout))
		// 	methodReq.Reply(Method_NONE_ACCEPTABLE)
		// 	(*s.conn).SetReadDeadline(time.Time{})
		// 	err = fmt.Errorf("none method found : Method_NO_AUTH")
		// 	return
		// }
		s.method = Method_NO_AUTH
		//method select reply
		(*s.conn).SetReadDeadline(time.Now().Add(time.Second * s.timeout))
		err = methodReq.Reply(Method_NO_AUTH)
		(*s.conn).SetReadDeadline(time.Time{})
		if err != nil {
			err = fmt.Errorf("reply answer data fail,ERR: %s", err)
			return
		}
		// err = fmt.Errorf("% x", methodReq.Bytes())
	} else {
		//auth
		if !methodReq.Select(Method_USER_PASS) {
			(*s.conn).SetReadDeadline(time.Now().Add(time.Second * s.timeout))
			methodReq.Reply(Method_NONE_ACCEPTABLE)
			(*s.conn).SetReadDeadline(time.Time{})
			err = fmt.Errorf("none method found : Method_USER_PASS")
			return
		}
		s.method = Method_USER_PASS
		//method reply need auth
		(*s.conn).SetReadDeadline(time.Now().Add(time.Second * s.timeout))
		err = methodReq.Reply(Method_USER_PASS)
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
		if s.auth == nil || s.auth.CheckUserPass(s.user, s.password, _addr[0], "") {
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
	case CMD_BIND:
		err = request.TCPReply(REP_UNKNOWN)
		if err != nil {
			err = fmt.Errorf("TCPReply REP_UNKNOWN to %s fail,ERR: %s", remoteAddr, err)
			return
		}
		err = fmt.Errorf("cmd bind not supported, form: %s", remoteAddr)
		return
	case CMD_CONNECT:
		err = request.TCPReply(REP_SUCCESS)
		if err != nil {
			err = fmt.Errorf("TCPReply REP_SUCCESS to %s fail,ERR: %s", remoteAddr, err)
			return
		}
	case CMD_ASSOCIATE:
		if !s.enableUDP {
			request.UDPReply(REP_UNKNOWN, "0.0.0.0:0")
			if err != nil {
				err = fmt.Errorf("UDPReply REP_UNKNOWN to %s fail,ERR: %s", remoteAddr, err)
				return
			}
			err = fmt.Errorf("cmd associate not supported, form: %s", remoteAddr)
			return
		}
		a, _ := net.ResolveUDPAddr("udp", ":0")
		s.UDPConnListener, err = net.ListenUDP("udp", a)
		if err != nil {
			request.UDPReply(REP_UNKNOWN, "0.0.0.0:0")
			err = fmt.Errorf("udp bind fail,ERR: %s , for %s", err, remoteAddr)
			return
		}
		_, port, _ := net.SplitHostPort(s.UDPConnListener.LocalAddr().String())
		err = request.UDPReply(REP_SUCCESS, net.JoinHostPort(s.udpIP, port))
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

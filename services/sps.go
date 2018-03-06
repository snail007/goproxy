package services

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"runtime/debug"
	"snail007/proxy/utils"
	"snail007/proxy/utils/socks"
	"strconv"
	"strings"
)

type SPS struct {
	outPool        utils.OutPool
	cfg            SPSArgs
	domainResolver utils.DomainResolver
}

func NewSPS() Service {
	return &SPS{
		outPool: utils.OutPool{},
		cfg:     SPSArgs{},
	}
}
func (s *SPS) CheckArgs() {
	if *s.cfg.Parent == "" {
		log.Fatalf("parent required for %s %s", s.cfg.Protocol(), *s.cfg.Local)
	}
	if *s.cfg.ParentType == "" {
		log.Fatalf("parent type unkown,use -T <tls|tcp|kcp>")
	}
	if *s.cfg.ParentType == TYPE_TLS || *s.cfg.LocalType == TYPE_TLS {
		s.cfg.CertBytes, s.cfg.KeyBytes = utils.TlsBytes(*s.cfg.CertFile, *s.cfg.KeyFile)
	}
}
func (s *SPS) InitService() {
	s.InitOutConnPool()
}
func (s *SPS) InitOutConnPool() {
	if *s.cfg.ParentType == TYPE_TLS || *s.cfg.ParentType == TYPE_TCP || *s.cfg.ParentType == TYPE_KCP {
		//dur int, isTLS bool, certBytes, keyBytes []byte,
		//parent string, timeout int, InitialCap int, MaxCap int
		s.outPool = utils.NewOutPool(
			0,
			*s.cfg.ParentType,
			*s.cfg.KCPMethod,
			*s.cfg.KCPKey,
			s.cfg.CertBytes, s.cfg.KeyBytes,
			*s.cfg.Parent,
			*s.cfg.Timeout,
			0,
			0,
		)
	}
}

func (s *SPS) StopService() {
	if s.outPool.Pool != nil {
		s.outPool.Pool.ReleaseAll()
	}
}
func (s *SPS) Start(args interface{}) (err error) {
	s.cfg = args.(SPSArgs)
	s.CheckArgs()
	log.Printf("use %s %s parent %s", *s.cfg.ParentType, *s.cfg.ParentServiceType, *s.cfg.Parent)
	s.InitService()

	for _, addr := range strings.Split(*s.cfg.Local, ",") {
		if addr != "" {
			host, port, _ := net.SplitHostPort(*s.cfg.Local)
			p, _ := strconv.Atoi(port)
			sc := utils.NewServerChannel(host, p)
			if *s.cfg.LocalType == TYPE_TCP {
				err = sc.ListenTCP(s.callback)
			} else if *s.cfg.LocalType == TYPE_TLS {
				err = sc.ListenTls(s.cfg.CertBytes, s.cfg.KeyBytes, s.callback)
			} else if *s.cfg.LocalType == TYPE_KCP {
				err = sc.ListenKCP(*s.cfg.KCPMethod, *s.cfg.KCPKey, s.callback)
			}
			if err != nil {
				return
			}
			log.Printf("%s http(s)+socks proxy on %s", s.cfg.Protocol(), (*sc.Listener).Addr())
		}
	}
	return
}

func (s *SPS) Clean() {
	s.StopService()
}
func (s *SPS) callback(inConn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%s conn handler crashed with err : %s \nstack: %s", s.cfg.Protocol(), err, string(debug.Stack()))
		}
	}()
	var err error
	switch *s.cfg.ParentType {
	case TYPE_KCP:
		fallthrough
	case TYPE_TCP:
		fallthrough
	case TYPE_TLS:
		err = s.OutToTCP(&inConn)
	default:
		err = fmt.Errorf("unkown parent type %s", *s.cfg.ParentType)
	}
	if err != nil {
		log.Printf("connect to %s parent %s fail, ERR:%s", *s.cfg.ParentType, *s.cfg.Parent, err)
		utils.CloseConn(&inConn)
	}
}
func (s *SPS) OutToTCP(inConn *net.Conn) (err error) {
	buf := make([]byte, 1024)
	n, err := (*inConn).Read(buf)
	header := buf[:n]
	if err != nil {
		log.Printf("ERR:%s", err)
		utils.CloseConn(inConn)
		return
	}
	address := ""
	var forwardBytes []byte
	//fmt.Printf("%v", header)
	if header[0] == socks.VERSION_V5 {
		//socks
		methodReq, e := socks.NewMethodsRequest(*inConn, header)
		if e != nil {
			log.Printf("new method request err:%s", e)
			utils.CloseConn(inConn)
			err = e.(error)
			return
		}
		if !methodReq.Select(socks.Method_NO_AUTH) {
			methodReq.Reply(socks.Method_NONE_ACCEPTABLE)
			utils.CloseConn(inConn)
			log.Printf("none method found : Method_NO_AUTH")
			return
		}
		//method select reply
		err = methodReq.Reply(socks.Method_NO_AUTH)
		if err != nil {
			log.Printf("reply answer data fail,ERR: %s", err)
			utils.CloseConn(inConn)
			return
		}
		//request detail
		request, e := socks.NewRequest(*inConn)
		if e != nil {
			log.Printf("read request data fail,ERR: %s", e)
			utils.CloseConn(inConn)
			err = e.(error)
			return
		}
		if request.CMD() != socks.CMD_CONNECT {
			//只支持tcp
			request.TCPReply(socks.REP_UNKNOWN)
			utils.CloseConn(inConn)
			err = errors.New("cmd not supported")
			return
		}
		address = request.Addr()
		request.TCPReply(socks.REP_SUCCESS)
	} else if bytes.IndexByte(header, '\n') != -1 {
		//http
		var request utils.HTTPRequest
		request, err = utils.NewHTTPRequest(inConn, 1024, false, nil, header)
		if err != nil {
			log.Printf("new http request fail,ERR: %s", err)
			utils.CloseConn(inConn)
			return
		}
		if len(header) >= 7 && strings.ToLower(string(header[:7])) == "connect" {
			//https
			request.HTTPSReply()
			//log.Printf("https reply: %s", request.Host)
		} else {
			forwardBytes = request.HeadBuf
		}
		address = request.Host
	} else {
		log.Printf("unknown request from: %s,%s", (*inConn).RemoteAddr(), string(header))
		utils.CloseConn(inConn)
		err = errors.New("unknown request")
		return
	}
	//connect to parent
	var outConn net.Conn
	var _outConn interface{}
	_outConn, err = s.outPool.Pool.Get()
	if err == nil {
		outConn = _outConn.(net.Conn)
	}
	if err != nil {
		log.Printf("connect to %s , err:%s", *s.cfg.Parent, err)
		utils.CloseConn(inConn)
		return
	}
	//ask parent for connect to target address
	if *s.cfg.ParentServiceType == "http" {
		//http parent
		fmt.Fprintf(outConn, "CONNECT %s HTTP/1.1\r\n", address)
		reply := make([]byte, 100)
		n, err = outConn.Read(reply)
		if err != nil {
			log.Printf("read reply from %s , err:%s", *s.cfg.Parent, err)
			utils.CloseConn(inConn)
			utils.CloseConn(&outConn)
			return
		}
		//log.Printf("reply: %s", string(reply[:n]))
	} else {
		log.Printf("connect %s", address)
		//socks parent
		//send auth type
		_, err = outConn.Write([]byte{0x05, 0x01, 0x00})
		if err != nil {
			log.Printf("write method to %s fail, err:%s", *s.cfg.Parent, err)
			utils.CloseConn(inConn)
			utils.CloseConn(&outConn)
			return
		}
		//read reply
		reply := make([]byte, 512)
		n, err = outConn.Read(reply)
		if err != nil {
			log.Printf("read reply from %s , err:%s", *s.cfg.Parent, err)
			utils.CloseConn(inConn)
			utils.CloseConn(&outConn)
			return
		}
		//log.Printf("method reply %v", reply[:n])

		//build request
		buf, err = s.buildRequest(address)
		if err != nil {
			log.Printf("build request to %s fail , err:%s", *s.cfg.Parent, err)
			utils.CloseConn(inConn)
			utils.CloseConn(&outConn)
			return
		}
		//send address request
		_, err = outConn.Write(buf)
		if err != nil {
			log.Printf("write request to %s fail, err:%s", *s.cfg.Parent, err)
			utils.CloseConn(inConn)
			utils.CloseConn(&outConn)
			return
		}
		//read reply
		reply = make([]byte, 512)
		n, err = outConn.Read(reply)
		if err != nil {
			log.Printf("read reply from %s , err:%s", *s.cfg.Parent, err)
			utils.CloseConn(inConn)
			utils.CloseConn(&outConn)
			return
		}

		//log.Printf("request reply %v", reply[:n])
	}
	//forward client data to target,if necessary.
	if len(forwardBytes) > 0 {
		outConn.Write(forwardBytes)
	}
	//bind
	inAddr := (*inConn).RemoteAddr().String()
	outAddr := outConn.RemoteAddr().String()
	utils.IoBind((*inConn), outConn, func(err interface{}) {
		log.Printf("conn %s - %s released", inAddr, outAddr)
	})
	log.Printf("conn %s - %s connected", inAddr, outAddr)
	return
}
func (s *SPS) buildRequest(address string) (buf []byte, err error) {
	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		err = errors.New("proxy: failed to parse port number: " + portStr)
		return
	}
	if port < 1 || port > 0xffff {
		err = errors.New("proxy: port number out of range: " + portStr)
		return
	}
	buf = buf[:0]
	buf = append(buf, 0x05, 0x01, 0 /* reserved */)

	if ip := net.ParseIP(host); ip != nil {
		if ip4 := ip.To4(); ip4 != nil {
			buf = append(buf, 0x01)
			ip = ip4
		} else {
			buf = append(buf, 0x04)
		}
		buf = append(buf, ip...)
	} else {
		if len(host) > 255 {
			err = errors.New("proxy: destination host name too long: " + host)
			return
		}
		buf = append(buf, 0x03)
		buf = append(buf, byte(len(host)))
		buf = append(buf, host...)
	}
	buf = append(buf, byte(port>>8), byte(port))
	return
}
func (s *SPS) Resolve(address string) string {
	if *s.cfg.DNSAddress == "" {
		return address
	}
	ip, err := s.domainResolver.Resolve(address)
	if err != nil {
		log.Printf("dns error %s , ERR:%s", address, err)
	}
	return ip
}

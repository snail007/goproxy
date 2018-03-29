package services

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"runtime/debug"
	"snail007/proxy/utils"
	"snail007/proxy/utils/socks"
	"strconv"
	"strings"
	"time"
)

type SPS struct {
	outPool        utils.OutPool
	cfg            SPSArgs
	domainResolver utils.DomainResolver
	basicAuth      utils.BasicAuth
}

func NewSPS() Service {
	return &SPS{
		outPool:   utils.OutPool{},
		cfg:       SPSArgs{},
		basicAuth: utils.BasicAuth{},
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
		if *s.cfg.CaCertFile != "" {
			var err error
			s.cfg.CaCertBytes, err = ioutil.ReadFile(*s.cfg.CaCertFile)
			if err != nil {
				log.Fatalf("read ca file error,ERR:%s", err)
			}
		}
	}
}
func (s *SPS) InitService() {
	s.InitOutConnPool()
	if *s.cfg.DNSAddress != "" {
		(*s).domainResolver = utils.NewDomainResolver(*s.cfg.DNSAddress, *s.cfg.DNSTTL)
	}
	s.InitBasicAuth()
}
func (s *SPS) InitOutConnPool() {
	if *s.cfg.ParentType == TYPE_TLS || *s.cfg.ParentType == TYPE_TCP || *s.cfg.ParentType == TYPE_KCP {
		//dur int, isTLS bool, certBytes, keyBytes []byte,
		//parent string, timeout int, InitialCap int, MaxCap int
		s.outPool = utils.NewOutPool(
			0,
			*s.cfg.ParentType,
			s.cfg.KCP,
			s.cfg.CertBytes, s.cfg.KeyBytes, nil,
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
				err = sc.ListenTls(s.cfg.CertBytes, s.cfg.KeyBytes, nil, s.callback)
			} else if *s.cfg.LocalType == TYPE_KCP {
				err = sc.ListenKCP(s.cfg.KCP, s.callback)
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
		log.Printf("connect to %s parent %s fail, ERR:%s from %s", *s.cfg.ParentType, *s.cfg.Parent, err, inConn.RemoteAddr())
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
	var auth socks.Auth
	var forwardBytes []byte
	//fmt.Printf("%v", header)
	if header[0] == socks.VERSION_V5 {
		//socks5 server
		var serverConn *socks.ServerConn
		if s.IsBasicAuth() {
			serverConn = socks.NewServerConn(inConn, time.Millisecond*time.Duration(*s.cfg.Timeout), &s.basicAuth, "", header)
		} else {
			serverConn = socks.NewServerConn(inConn, time.Millisecond*time.Duration(*s.cfg.Timeout), nil, "", header)
		}
		if err = serverConn.Handshake(); err != nil {
			return
		}
		address = serverConn.Target()
		auth = serverConn.AuthData()
	} else if bytes.IndexByte(header, '\n') != -1 {
		//http
		var request utils.HTTPRequest
		(*inConn).SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
		if s.IsBasicAuth() {
			request, err = utils.NewHTTPRequest(inConn, 1024, true, &s.basicAuth, header)
		} else {
			request, err = utils.NewHTTPRequest(inConn, 1024, false, nil, header)
		}
		(*inConn).SetDeadline(time.Time{})
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
		var userpass string
		if s.IsBasicAuth() {
			userpass, err = request.GetAuthDataStr()
			if err != nil {
				return
			}
			userpassA := strings.Split(userpass, ":")
			if len(userpassA) == 2 {
				auth = socks.Auth{User: userpassA[0], Password: userpassA[1]}
			}
		}
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
		pb := new(bytes.Buffer)
		pb.Write([]byte(fmt.Sprintf("CONNECT %s HTTP/1.1\r\nProxy-Connection: Keep-Alive\r\n", address)))
		//Proxy-Authorization:\r\n
		u := ""
		if *s.cfg.ParentAuth != "" {
			a := strings.Split(*s.cfg.ParentAuth, ":")
			if len(a) != 2 {
				err = fmt.Errorf("parent auth data format error")
				return
			}
			u = fmt.Sprintf("%s:%s", a[0], a[1])
		} else {
			if !s.IsBasicAuth() && auth.Password != "" && auth.User != "" {
				u = fmt.Sprintf("%s:%s", auth.User, auth.Password)
			}
		}
		if u != "" {
			pb.Write([]byte(fmt.Sprintf("Proxy-Authorization:Basic %s\r\n", base64.StdEncoding.EncodeToString([]byte(u)))))
		}
		pb.Write([]byte("\r\n"))
		outConn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
		_, err = outConn.Write(pb.Bytes())
		outConn.SetDeadline(time.Time{})
		if err != nil {
			log.Printf("write CONNECT to %s , err:%s", *s.cfg.Parent, err)
			utils.CloseConn(inConn)
			utils.CloseConn(&outConn)
			return
		}
		reply := make([]byte, 1024)
		outConn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
		_, err = outConn.Read(reply)
		outConn.SetDeadline(time.Time{})
		if err != nil {
			log.Printf("read reply from %s , err:%s", *s.cfg.Parent, err)
			utils.CloseConn(inConn)
			utils.CloseConn(&outConn)
			return
		}
		//log.Printf("reply: %s", string(reply[:n]))
	} else {
		log.Printf("connect %s", address)
		//socks client
		var clientConn *socks.ClientConn
		if *s.cfg.ParentAuth != "" {
			a := strings.Split(*s.cfg.ParentAuth, ":")
			if len(a) != 2 {
				err = fmt.Errorf("parent auth data format error")
				return
			}
			clientConn = socks.NewClientConn(&outConn, "tcp", address, time.Millisecond*time.Duration(*s.cfg.Timeout), &socks.Auth{User: a[0], Password: a[1]}, header)
		} else {
			if !s.IsBasicAuth() && auth.Password != "" && auth.User != "" {
				clientConn = socks.NewClientConn(&outConn, "tcp", address, time.Millisecond*time.Duration(*s.cfg.Timeout), &auth, header)
			} else {
				clientConn = socks.NewClientConn(&outConn, "tcp", address, time.Millisecond*time.Duration(*s.cfg.Timeout), nil, header)
			}
		}
		if err = clientConn.Handshake(); err != nil {
			return
		}
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
func (s *SPS) InitBasicAuth() (err error) {
	if *s.cfg.DNSAddress != "" {
		s.basicAuth = utils.NewBasicAuth(&(*s).domainResolver)
	} else {
		s.basicAuth = utils.NewBasicAuth(nil)
	}
	if *s.cfg.AuthURL != "" {
		s.basicAuth.SetAuthURL(*s.cfg.AuthURL, *s.cfg.AuthURLOkCode, *s.cfg.AuthURLTimeout, *s.cfg.AuthURLRetry)
		log.Printf("auth from %s", *s.cfg.AuthURL)
	}
	if *s.cfg.AuthFile != "" {
		var n = 0
		n, err = s.basicAuth.AddFromFile(*s.cfg.AuthFile)
		if err != nil {
			err = fmt.Errorf("auth-file ERR:%s", err)
			return
		}
		log.Printf("auth data added from file %d , total:%d", n, s.basicAuth.Total())
	}
	if len(*s.cfg.Auth) > 0 {
		n := s.basicAuth.Add(*s.cfg.Auth)
		log.Printf("auth data added %d, total:%d", n, s.basicAuth.Total())
	}
	return
}
func (s *SPS) IsBasicAuth() bool {
	return *s.cfg.AuthFile != "" || len(*s.cfg.Auth) > 0 || *s.cfg.AuthURL != ""
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

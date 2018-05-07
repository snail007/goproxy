package services

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	logger "log"
	"net"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/snail007/goproxy/utils"
	"github.com/snail007/goproxy/utils/conncrypt"
	"github.com/snail007/goproxy/utils/socks"
	"github.com/snail007/goproxy/utils/ss"
)

type SPS struct {
	outPool        utils.OutConn
	cfg            SPSArgs
	domainResolver utils.DomainResolver
	basicAuth      utils.BasicAuth
	serverChannels []*utils.ServerChannel
	userConns      utils.ConcurrentMap
	log            *logger.Logger
	cipher         *ss.Cipher
}

func NewSPS() Service {
	return &SPS{
		outPool:        utils.OutConn{},
		cfg:            SPSArgs{},
		basicAuth:      utils.BasicAuth{},
		serverChannels: []*utils.ServerChannel{},
		userConns:      utils.NewConcurrentMap(),
	}
}
func (s *SPS) CheckArgs() (err error) {
	if *s.cfg.Parent == "" {
		err = fmt.Errorf("parent required for %s %s", s.cfg.Protocol(), *s.cfg.Local)
		return
	}
	if *s.cfg.ParentType == "" {
		err = fmt.Errorf("parent type unkown,use -T <tls|tcp|kcp>")
		return
	}
	if *s.cfg.ParentType == TYPE_TLS || *s.cfg.LocalType == TYPE_TLS {
		s.cfg.CertBytes, s.cfg.KeyBytes, err = utils.TlsBytes(*s.cfg.CertFile, *s.cfg.KeyFile)
		if err != nil {
			return
		}
		if *s.cfg.CaCertFile != "" {
			s.cfg.CaCertBytes, err = ioutil.ReadFile(*s.cfg.CaCertFile)
			if err != nil {
				err = fmt.Errorf("read ca file error,ERR:%s", err)
				return
			}
		}
	}
	return
}
func (s *SPS) InitService() (err error) {
	s.InitOutConnPool()
	if *s.cfg.DNSAddress != "" {
		(*s).domainResolver = utils.NewDomainResolver(*s.cfg.DNSAddress, *s.cfg.DNSTTL, s.log)
	}
	err = s.InitBasicAuth()
	if *s.cfg.SSMethod != "" && *s.cfg.SSKey != "" {
		s.cipher, err = ss.NewCipher(*s.cfg.SSMethod, *s.cfg.SSKey)
		if err != nil {
			s.log.Printf("error generating cipher : %s", err)
			return
		}
	}
	return
}
func (s *SPS) InitOutConnPool() {
	if *s.cfg.ParentType == TYPE_TLS || *s.cfg.ParentType == TYPE_TCP || *s.cfg.ParentType == TYPE_KCP {
		//dur int, isTLS bool, certBytes, keyBytes []byte,
		//parent string, timeout int, InitialCap int, MaxCap int
		s.outPool = utils.NewOutConn(
			0,
			*s.cfg.ParentType,
			s.cfg.KCP,
			s.cfg.CertBytes, s.cfg.KeyBytes, s.cfg.CaCertBytes,
			*s.cfg.Parent,
			*s.cfg.Timeout,
		)
	}
}

func (s *SPS) StopService() {
	defer func() {
		e := recover()
		if e != nil {
			s.log.Printf("stop sps service crashed,%s", e)
		} else {
			s.log.Printf("service sps stoped")
		}
	}()
	for _, sc := range s.serverChannels {
		if sc.Listener != nil && *sc.Listener != nil {
			(*sc.Listener).Close()
		}
		if sc.UDPListener != nil {
			(*sc.UDPListener).Close()
		}
	}
	for _, c := range s.userConns.Items() {
		(*c.(*net.Conn)).Close()
	}
}
func (s *SPS) Start(args interface{}, log *logger.Logger) (err error) {
	s.log = log
	s.cfg = args.(SPSArgs)
	if err = s.CheckArgs(); err != nil {
		return
	}
	if err = s.InitService(); err != nil {
		return
	}
	s.log.Printf("use %s %s parent %s", *s.cfg.ParentType, *s.cfg.ParentServiceType, *s.cfg.Parent)
	for _, addr := range strings.Split(*s.cfg.Local, ",") {
		if addr != "" {
			host, port, _ := net.SplitHostPort(*s.cfg.Local)
			p, _ := strconv.Atoi(port)
			sc := utils.NewServerChannel(host, p, s.log)
			if *s.cfg.LocalType == TYPE_TCP {
				err = sc.ListenTCP(s.callback)
			} else if *s.cfg.LocalType == TYPE_TLS {
				err = sc.ListenTls(s.cfg.CertBytes, s.cfg.KeyBytes, s.cfg.CaCertBytes, s.callback)
			} else if *s.cfg.LocalType == TYPE_KCP {
				err = sc.ListenKCP(s.cfg.KCP, s.callback, s.log)
			}
			if err != nil {
				return
			}
			s.log.Printf("%s http(s)+socks proxy on %s", s.cfg.Protocol(), (*sc.Listener).Addr())
			s.serverChannels = append(s.serverChannels, &sc)
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
			s.log.Printf("%s conn handler crashed with err : %s \nstack: %s", s.cfg.Protocol(), err, string(debug.Stack()))
		}
	}()
	if *s.cfg.LocalCompress {
		inConn = utils.NewCompConn(inConn)
	}
	if *s.cfg.LocalKey != "" {
		inConn = conncrypt.New(inConn, &conncrypt.Config{
			Password: *s.cfg.LocalKey,
		})
	}
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
		s.log.Printf("connect to %s parent %s fail, ERR:%s from %s", *s.cfg.ParentType, *s.cfg.Parent, err, inConn.RemoteAddr())
		utils.CloseConn(&inConn)
	}
}
func (s *SPS) OutToTCP(inConn *net.Conn) (err error) {
	bInConn := utils.NewBufferedConn(*inConn)
	//important
	//action read will regist read event to system,
	//when data arrived , system call process
	//so that we can get buffered bytes count
	//otherwise Buffered() always return 0
	bInConn.ReadByte()
	bInConn.UnreadByte()

	n := 8
	if n > bInConn.Buffered() {
		n = bInConn.Buffered()
	}
	h, err := bInConn.Peek(n)
	if err != nil {
		s.log.Printf("peek error %s ", err)
		(*inConn).Close()
		return
	}

	*inConn = bInConn
	address := ""
	var auth socks.Auth
	var forwardBytes []byte
	//fmt.Printf("%v", header)
	if utils.IsSocks5(h) {
		if *s.cfg.DisableSocks5 {
			(*inConn).Close()
			return
		}
		//socks5 server
		var serverConn *socks.ServerConn
		if s.IsBasicAuth() {
			serverConn = socks.NewServerConn(inConn, time.Millisecond*time.Duration(*s.cfg.Timeout), &s.basicAuth, "", nil)
		} else {
			serverConn = socks.NewServerConn(inConn, time.Millisecond*time.Duration(*s.cfg.Timeout), nil, "", nil)
		}
		if err = serverConn.Handshake(); err != nil {
			return
		}
		address = serverConn.Target()
		auth = serverConn.AuthData()
	} else if utils.IsHTTP(h) {
		if *s.cfg.DisableHTTP {
			(*inConn).Close()
			return
		}
		//http
		var request utils.HTTPRequest
		(*inConn).SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
		if s.IsBasicAuth() {
			request, err = utils.NewHTTPRequest(inConn, 1024, true, &s.basicAuth, s.log)
		} else {
			request, err = utils.NewHTTPRequest(inConn, 1024, false, nil, s.log)
		}
		(*inConn).SetDeadline(time.Time{})
		if err != nil {
			s.log.Printf("new http request fail,ERR: %s", err)
			utils.CloseConn(inConn)
			return
		}
		if len(h) >= 7 && strings.ToLower(string(h[:7])) == "connect" {
			//https
			request.HTTPSReply()
			//s.log.Printf("https reply: %s", request.Host)
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
		//ss
		if *s.cfg.DisableSS {
			(*inConn).Close()
			return
		}
		(*inConn).SetDeadline(time.Now().Add(time.Second * 5))
		ssConn := ss.NewConn(*inConn, s.cipher.Copy())
		address, err = ss.GetRequest(ssConn)
		(*inConn).SetDeadline(time.Time{})
		if err != nil {
			return
		}
		// ensure the host does not contain some illegal characters, NUL may panic on Win32
		if strings.ContainsRune(address, 0x00) {
			err = errors.New("invalid domain name")
			return
		}
		*inConn = ssConn
	}
	if err != nil || address == "" {
		s.log.Printf("unknown request from: %s,%s", (*inConn).RemoteAddr(), string(h))
		(*inConn).Close()
		utils.CloseConn(inConn)
		err = errors.New("unknown request")
		return
	}
	//connect to parent
	var outConn net.Conn
	outConn, err = s.outPool.Get()
	if err != nil {
		s.log.Printf("connect to %s , err:%s", *s.cfg.Parent, err)
		utils.CloseConn(inConn)
		return
	}
	if *s.cfg.ParentCompress {
		outConn = utils.NewCompConn(outConn)
	}
	if *s.cfg.ParentKey != "" {
		outConn = conncrypt.New(outConn, &conncrypt.Config{
			Password: *s.cfg.ParentKey,
		})
	}
	//ask parent for connect to target address
	if *s.cfg.ParentServiceType == "http" {
		//http parent
		pb := new(bytes.Buffer)
		pb.Write([]byte(fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost:%s\r\nProxy-Connection: Keep-Alive\r\n", address, strings.TrimRight(address,":80"))))
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
			s.log.Printf("write CONNECT to %s , err:%s", *s.cfg.Parent, err)
			utils.CloseConn(inConn)
			utils.CloseConn(&outConn)
			return
		}
		reply := make([]byte, 1024)
		outConn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
		_, err = outConn.Read(reply)
		outConn.SetDeadline(time.Time{})
		if err != nil {
			s.log.Printf("read reply from %s , err:%s", *s.cfg.Parent, err)
			utils.CloseConn(inConn)
			utils.CloseConn(&outConn)
			return
		}
		//s.log.Printf("reply: %s", string(reply[:n]))
	} else {
		s.log.Printf("connect %s", address)
		//socks client
		var clientConn *socks.ClientConn
		if *s.cfg.ParentAuth != "" {
			a := strings.Split(*s.cfg.ParentAuth, ":")
			if len(a) != 2 {
				err = fmt.Errorf("parent auth data format error")
				return
			}
			clientConn = socks.NewClientConn(&outConn, "tcp", address, time.Millisecond*time.Duration(*s.cfg.Timeout), &socks.Auth{User: a[0], Password: a[1]}, nil)
		} else {
			if !s.IsBasicAuth() && auth.Password != "" && auth.User != "" {
				clientConn = socks.NewClientConn(&outConn, "tcp", address, time.Millisecond*time.Duration(*s.cfg.Timeout), &auth, nil)
			} else {
				clientConn = socks.NewClientConn(&outConn, "tcp", address, time.Millisecond*time.Duration(*s.cfg.Timeout), nil, nil)
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
		s.log.Printf("conn %s - %s released", inAddr, outAddr)
		s.userConns.Remove(inAddr)
	}, s.log)
	s.log.Printf("conn %s - %s connected", inAddr, outAddr)
	if c, ok := s.userConns.Get(inAddr); ok {
		(*c.(*net.Conn)).Close()
	}
	s.userConns.Set(inAddr, inConn)
	return
}
func (s *SPS) InitBasicAuth() (err error) {
	if *s.cfg.DNSAddress != "" {
		s.basicAuth = utils.NewBasicAuth(&(*s).domainResolver, s.log)
	} else {
		s.basicAuth = utils.NewBasicAuth(nil, s.log)
	}
	if *s.cfg.AuthURL != "" {
		s.basicAuth.SetAuthURL(*s.cfg.AuthURL, *s.cfg.AuthURLOkCode, *s.cfg.AuthURLTimeout, *s.cfg.AuthURLRetry)
		s.log.Printf("auth from %s", *s.cfg.AuthURL)
	}
	if *s.cfg.AuthFile != "" {
		var n = 0
		n, err = s.basicAuth.AddFromFile(*s.cfg.AuthFile)
		if err != nil {
			err = fmt.Errorf("auth-file ERR:%s", err)
			return
		}
		s.log.Printf("auth data added from file %d , total:%d", n, s.basicAuth.Total())
	}
	if len(*s.cfg.Auth) > 0 {
		n := s.basicAuth.Add(*s.cfg.Auth)
		s.log.Printf("auth data added %d, total:%d", n, s.basicAuth.Total())
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
		s.log.Printf("dns error %s , ERR:%s", address, err)
	}
	return ip
}

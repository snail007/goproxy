package socks

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	logger "log"
	"net"
	"runtime/debug"
	"strings"
	"time"

	"github.com/visenze/goproxy/services"
	"github.com/visenze/goproxy/services/kcpcfg"
	"github.com/visenze/goproxy/utils"
	"github.com/visenze/goproxy/utils/conncrypt"
	"github.com/visenze/goproxy/utils/socks"
	"golang.org/x/crypto/ssh"
)

type SocksArgs struct {
	Parent         *string
	ParentType     *string
	Local          *string
	LocalType      *string
	CertFile       *string
	KeyFile        *string
	CaCertFile     *string
	CaCertBytes    []byte
	CertBytes      []byte
	KeyBytes       []byte
	SSHKeyFile     *string
	SSHKeyFileSalt *string
	SSHPassword    *string
	SSHUser        *string
	SSHKeyBytes    []byte
	SSHAuthMethod  ssh.AuthMethod
	Timeout        *int
	Always         *bool
	Interval       *int
	Blocked        *string
	Direct         *string
	AuthFile       *string
	Auth           *[]string
	AuthURL        *string
	AuthURLOkCode  *int
	AuthURLTimeout *int
	AuthURLRetry   *int
	KCP            kcpcfg.KCPConfigArgs
	LocalIPS       *[]string
	DNSAddress     *string
	DNSTTL         *int
	LocalKey       *string
	ParentKey      *string
	LocalCompress  *bool
	ParentCompress *bool
}
type Socks struct {
	cfg                   SocksArgs
	checker               utils.Checker
	basicAuth             utils.BasicAuth
	sshClient             *ssh.Client
	lockChn               chan bool
	udpSC                 utils.ServerChannel
	sc                    *utils.ServerChannel
	domainResolver        utils.DomainResolver
	isStop                bool
	userConns             utils.ConcurrentMap
	log                   *logger.Logger
	udpRelatedPacketConns utils.ConcurrentMap
	udpLocalKey           []byte
	udpParentKey          []byte
}

func NewSocks() services.Service {
	return &Socks{
		cfg:                   SocksArgs{},
		checker:               utils.Checker{},
		basicAuth:             utils.BasicAuth{},
		lockChn:               make(chan bool, 1),
		isStop:                false,
		userConns:             utils.NewConcurrentMap(),
		udpRelatedPacketConns: utils.NewConcurrentMap(),
	}
}

func (s *Socks) CheckArgs() (err error) {

	if *s.cfg.LocalType == "tls" || (*s.cfg.Parent != "" && *s.cfg.ParentType == "tls") {
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
	if *s.cfg.Parent != "" {
		if *s.cfg.ParentType == "" {
			err = fmt.Errorf("parent type unkown,use -T <tls|tcp|ssh|kcp>")
			return
		}
		if *s.cfg.ParentType == "ssh" {
			if *s.cfg.SSHUser == "" {
				err = fmt.Errorf("ssh user required")
				return
			}
			if *s.cfg.SSHKeyFile == "" && *s.cfg.SSHPassword == "" {
				err = fmt.Errorf("ssh password or key required")
				return
			}
			if *s.cfg.SSHPassword != "" {
				s.cfg.SSHAuthMethod = ssh.Password(*s.cfg.SSHPassword)
			} else {
				var SSHSigner ssh.Signer
				s.cfg.SSHKeyBytes, err = ioutil.ReadFile(*s.cfg.SSHKeyFile)
				if err != nil {
					err = fmt.Errorf("read key file ERR: %s", err)
					return
				}
				if *s.cfg.SSHKeyFileSalt != "" {
					SSHSigner, err = ssh.ParsePrivateKeyWithPassphrase(s.cfg.SSHKeyBytes, []byte(*s.cfg.SSHKeyFileSalt))
				} else {
					SSHSigner, err = ssh.ParsePrivateKey(s.cfg.SSHKeyBytes)
				}
				if err != nil {
					err = fmt.Errorf("parse ssh private key fail,ERR: %s", err)
					return
				}
				s.cfg.SSHAuthMethod = ssh.PublicKeys(SSHSigner)
			}
		}
	}
	s.udpLocalKey = s.LocalUDPKey()
	s.udpParentKey = s.ParentUDPKey()
	//s.log.Printf("udpLocalKey : %v , udpParentKey : %v", s.udpLocalKey, s.udpParentKey)
	return
}
func (s *Socks) InitService() (err error) {
	s.InitBasicAuth()
	if *s.cfg.DNSAddress != "" {
		(*s).domainResolver = utils.NewDomainResolver(*s.cfg.DNSAddress, *s.cfg.DNSTTL, s.log)
	}
	s.checker = utils.NewChecker(*s.cfg.Timeout, int64(*s.cfg.Interval), *s.cfg.Blocked, *s.cfg.Direct, s.log)
	if *s.cfg.ParentType == "ssh" {
		e := s.ConnectSSH()
		if e != nil {
			err = fmt.Errorf("init service fail, ERR: %s", e)
			return
		}
		go func() {
			//循环检查ssh网络连通性
			for {
				if s.isStop {
					return
				}
				conn, err := utils.ConnectHost(s.Resolve(*s.cfg.Parent), *s.cfg.Timeout*2)
				if err == nil {
					conn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
					_, err = conn.Write([]byte{0})
					conn.SetDeadline(time.Time{})
				}
				if err != nil {
					if s.sshClient != nil {
						s.sshClient.Close()
					}
					s.log.Printf("ssh offline, retrying...")
					s.ConnectSSH()
				} else {
					conn.Close()
				}
				time.Sleep(time.Second * 3)
			}
		}()
	}
	if *s.cfg.ParentType == "ssh" {
		s.log.Printf("warn: socks udp not suppored for ssh")
	}
	return
}
func (s *Socks) StopService() {
	defer func() {
		e := recover()
		if e != nil {
			s.log.Printf("stop socks service crashed,%s", e)
		} else {
			s.log.Printf("service socks stoped")
		}
	}()
	s.isStop = true
	s.checker.Stop()
	if s.sshClient != nil {
		s.sshClient.Close()
	}
	if s.udpSC.UDPListener != nil {
		s.udpSC.UDPListener.Close()
	}
	if s.sc != nil && (*s.sc).Listener != nil {
		(*(*s.sc).Listener).Close()
	}
	for _, c := range s.userConns.Items() {
		(*c.(*net.Conn)).Close()
	}
	for _, c := range s.udpRelatedPacketConns.Items() {
		(*c.(*net.UDPConn)).Close()
	}
}
func (s *Socks) Start(args interface{}, log *logger.Logger) (err error) {
	s.log = log
	//start()
	s.cfg = args.(SocksArgs)
	if err = s.CheckArgs(); err != nil {
		return
	}
	if err = s.InitService(); err != nil {
		s.InitService()
	}
	if *s.cfg.Parent != "" {
		s.log.Printf("use %s parent %s", *s.cfg.ParentType, *s.cfg.Parent)
	}
	sc := utils.NewServerChannelHost(*s.cfg.Local, s.log)
	if *s.cfg.LocalType == "tcp" {
		err = sc.ListenTCP(s.socksConnCallback)
	} else if *s.cfg.LocalType == "tls" {
		err = sc.ListenTls(s.cfg.CertBytes, s.cfg.KeyBytes, nil, s.socksConnCallback)
	} else if *s.cfg.LocalType == "kcp" {
		err = sc.ListenKCP(s.cfg.KCP, s.socksConnCallback, s.log)
	}
	if err != nil {
		return
	}
	s.sc = &sc
	s.log.Printf("%s socks proxy on %s", *s.cfg.LocalType, (*sc.Listener).Addr())
	return
}
func (s *Socks) Clean() {
	s.StopService()
}

func (s *Socks) socksConnCallback(inConn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			s.log.Printf("socks conn handler crashed with err : %s \nstack: %s", err, string(debug.Stack()))
			inConn.Close()
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
	//协商开始

	//method select request
	inConn.SetReadDeadline(time.Now().Add(time.Second * 3))
	methodReq, err := socks.NewMethodsRequest(inConn)
	inConn.SetReadDeadline(time.Time{})
	if err != nil {
		methodReq.Reply(socks.Method_NONE_ACCEPTABLE)
		utils.CloseConn(&inConn)
		if err != io.EOF {
			s.log.Printf("new methods request fail,ERR: %s", err)
		}
		return
	}

	if !s.IsBasicAuth() {
		if !methodReq.Select(socks.Method_NO_AUTH) {
			methodReq.Reply(socks.Method_NONE_ACCEPTABLE)
			utils.CloseConn(&inConn)
			s.log.Printf("none method found : Method_NO_AUTH")
			return
		}
		//method select reply
		err = methodReq.Reply(socks.Method_NO_AUTH)
		if err != nil {
			s.log.Printf("reply answer data fail,ERR: %s", err)
			utils.CloseConn(&inConn)
			return
		}
		// s.log.Printf("% x", methodReq.Bytes())
	} else {
		//auth
		if !methodReq.Select(socks.Method_USER_PASS) {
			methodReq.Reply(socks.Method_NONE_ACCEPTABLE)
			utils.CloseConn(&inConn)
			s.log.Printf("none method found : Method_USER_PASS")
			return
		}
		//method reply need auth
		err = methodReq.Reply(socks.Method_USER_PASS)
		if err != nil {
			s.log.Printf("reply answer data fail,ERR: %s", err)
			utils.CloseConn(&inConn)
			return
		}
		//read auth
		buf := make([]byte, 500)
		inConn.SetReadDeadline(time.Now().Add(time.Second * 3))
		n, err := inConn.Read(buf)
		inConn.SetReadDeadline(time.Time{})
		if err != nil {
			utils.CloseConn(&inConn)
			return
		}
		r := buf[:n]
		user := string(r[2 : r[1]+2])
		pass := string(r[2+r[1]+1:])
		//s.log.Printf("user:%s,pass:%s", user, pass)
		//auth
		_addr := strings.Split(inConn.RemoteAddr().String(), ":")
		if s.basicAuth.CheckUserPass(user, pass, _addr[0], "") {
			inConn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
			inConn.Write([]byte{0x01, 0x00})
			inConn.SetDeadline(time.Time{})

		} else {
			inConn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
			inConn.Write([]byte{0x01, 0x01})
			inConn.SetDeadline(time.Time{})

			utils.CloseConn(&inConn)
			return
		}
	}

	//request detail
	request, err := socks.NewRequest(inConn)
	if err != nil {
		s.log.Printf("read request data fail,ERR: %s", err)
		utils.CloseConn(&inConn)
		return
	}
	//协商结束

	switch request.CMD() {
	case socks.CMD_BIND:
		//bind 不支持
		request.TCPReply(socks.REP_UNKNOWN)
		utils.CloseConn(&inConn)
		return
	case socks.CMD_CONNECT:
		//tcp
		s.proxyTCP(&inConn, methodReq, request)
	case socks.CMD_ASSOCIATE:
		//udp
		s.proxyUDP(&inConn, methodReq, request)
	}

}

func (s *Socks) proxyTCP(inConn *net.Conn, methodReq socks.MethodsRequest, request socks.Request) {
	var outConn net.Conn
	var err interface{}
	useProxy := true
	tryCount := 0
	maxTryCount := 5
	//防止死循环
	if s.IsDeadLoop((*inConn).LocalAddr().String(), request.Host()) {
		utils.CloseConn(inConn)
		s.log.Printf("dead loop detected , %s", request.Host())
		utils.CloseConn(inConn)
		return
	}
	for {
		if s.isStop {
			return
		}
		if *s.cfg.Always {
			outConn, err = s.getOutConn(methodReq.Bytes(), request.Bytes(), request.Addr(), true)
		} else {
			if *s.cfg.Parent != "" {
				host, _, _ := net.SplitHostPort(request.Addr())
				useProxy := false
				if utils.IsIternalIP(host, *s.cfg.Always) {
					useProxy = false
				} else {
					var isInMap bool
					useProxy, isInMap, _, _ = s.checker.IsBlocked(request.Addr())
					if !isInMap {
						s.checker.Add(request.Addr(), s.Resolve(request.Addr()))
					}
				}
				if useProxy {
					outConn, err = s.getOutConn(methodReq.Bytes(), request.Bytes(), request.Addr(), true)
				} else {
					outConn, err = utils.ConnectHost(s.Resolve(request.Addr()), *s.cfg.Timeout)
				}
			} else {
				outConn, err = utils.ConnectHost(s.Resolve(request.Addr()), *s.cfg.Timeout)
				useProxy = false
			}
		}
		tryCount++
		if err == nil || tryCount > maxTryCount || *s.cfg.Parent == "" {
			break
		} else {
			s.log.Printf("get out conn fail,%s,retrying...", err)
			time.Sleep(time.Second * 2)
		}
	}
	if err != nil {
		s.log.Printf("get out conn fail,%s", err)
		request.TCPReply(socks.REP_NETWOR_UNREACHABLE)
		return
	}

	s.log.Printf("use proxy %v : %s", useProxy, request.Addr())

	request.TCPReply(socks.REP_SUCCESS)
	inAddr := (*inConn).RemoteAddr().String()
	//inLocalAddr := (*inConn).LocalAddr().String()

	s.log.Printf("conn %s - %s connected", inAddr, request.Addr())
	utils.IoBind(*inConn, outConn, func(err interface{}) {
		s.log.Printf("conn %s - %s released", inAddr, request.Addr())
		s.userConns.Remove(inAddr)
	}, s.log)
	if c, ok := s.userConns.Get(inAddr); ok {
		(*c.(*net.Conn)).Close()
		s.userConns.Remove(inAddr)
	}
	s.userConns.Set(inAddr, inConn)
}
func (s *Socks) getOutConn(methodBytes, reqBytes []byte, host string, handshake bool) (outConn net.Conn, err interface{}) {
	switch *s.cfg.ParentType {
	case "kcp":
		fallthrough
	case "tls":
		fallthrough
	case "tcp":
		if *s.cfg.ParentType == "tls" {
			var _outConn tls.Conn
			_outConn, err = utils.TlsConnectHost(s.Resolve(*s.cfg.Parent), *s.cfg.Timeout, s.cfg.CertBytes, s.cfg.KeyBytes, nil)
			outConn = net.Conn(&_outConn)
		} else if *s.cfg.ParentType == "kcp" {
			outConn, err = utils.ConnectKCPHost(s.Resolve(*s.cfg.Parent), s.cfg.KCP)
		} else {
			outConn, err = utils.ConnectHost(s.Resolve(*s.cfg.Parent), *s.cfg.Timeout)
		}
		if err != nil {
			err = fmt.Errorf("connect fail,%s", err)
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
		if !handshake {
			return
		}
		var buf = make([]byte, 1024)
		//var n int
		outConn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
		_, err = outConn.Write(methodBytes)
		outConn.SetDeadline(time.Time{})
		if err != nil {
			err = fmt.Errorf("write method fail,%s", err)
			return
		}
		outConn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
		_, err = outConn.Read(buf)
		outConn.SetDeadline(time.Time{})
		if err != nil {
			err = fmt.Errorf("read method reply fail,%s", err)
			return
		}
		//resp := buf[:n]
		//s.log.Printf("resp:%v", resp)
		outConn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
		_, err = outConn.Write(reqBytes)
		outConn.SetDeadline(time.Time{})
		if err != nil {
			err = fmt.Errorf("write req detail fail,%s", err)
			return
		}
		outConn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
		_, err = outConn.Read(buf)
		outConn.SetDeadline(time.Time{})
		if err != nil {
			err = fmt.Errorf("read req reply fail,%s", err)
			return
		}
		//result := buf[:n]
		//s.log.Printf("result:%v", result)

	case "ssh":
		maxTryCount := 1
		tryCount := 0
	RETRY:
		if tryCount >= maxTryCount || s.isStop {
			return
		}
		wait := make(chan bool, 1)
		go func() {
			defer func() {
				if err == nil {
					err = recover()
				}
				wait <- true
			}()
			outConn, err = s.sshClient.Dial("tcp", host)
		}()
		select {
		case <-wait:
		case <-time.After(time.Millisecond * time.Duration(*s.cfg.Timeout) * 2):
			err = fmt.Errorf("ssh dial %s timeout", host)
			s.sshClient.Close()
		}
		if err != nil {
			s.log.Printf("connect ssh fail, ERR: %s, retrying...", err)
			e := s.ConnectSSH()
			if e == nil {
				tryCount++
				time.Sleep(time.Second * 3)
				goto RETRY
			} else {
				err = e
			}
		}
	}

	return
}
func (s *Socks) ConnectSSH() (err error) {
	select {
	case s.lockChn <- true:
	default:
		err = fmt.Errorf("can not connect at same time")
		return
	}
	config := ssh.ClientConfig{
		Timeout: time.Duration(*s.cfg.Timeout) * time.Millisecond,
		User:    *s.cfg.SSHUser,
		Auth:    []ssh.AuthMethod{s.cfg.SSHAuthMethod},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	if s.sshClient != nil {
		s.sshClient.Close()
	}
	s.sshClient, err = ssh.Dial("tcp", s.Resolve(*s.cfg.Parent), &config)
	<-s.lockChn
	return
}
func (s *Socks) InitBasicAuth() (err error) {
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
func (s *Socks) IsBasicAuth() bool {
	return *s.cfg.AuthFile != "" || len(*s.cfg.Auth) > 0 || *s.cfg.AuthURL != ""
}
func (s *Socks) IsDeadLoop(inLocalAddr string, host string) bool {
	inIP, inPort, err := net.SplitHostPort(inLocalAddr)
	if err != nil {
		return false
	}
	outDomain, outPort, err := net.SplitHostPort(host)
	if err != nil {
		return false
	}
	if inPort == outPort {
		var outIPs []net.IP
		if *s.cfg.DNSAddress != "" {
			outIPs = []net.IP{net.ParseIP(s.Resolve(outDomain))}
		} else {
			outIPs, err = net.LookupIP(outDomain)
		}
		if err == nil {
			for _, ip := range outIPs {
				if ip.String() == inIP {
					return true
				}
			}
		}
		interfaceIPs, err := utils.GetAllInterfaceAddr()
		for _, ip := range *s.cfg.LocalIPS {
			interfaceIPs = append(interfaceIPs, net.ParseIP(ip).To4())
		}
		if err == nil {
			for _, localIP := range interfaceIPs {
				for _, outIP := range outIPs {
					if localIP.Equal(outIP) {
						return true
					}
				}
			}
		}
	}
	return false
}
func (s *Socks) Resolve(address string) string {
	if *s.cfg.DNSAddress == "" {
		return address
	}
	ip, err := s.domainResolver.Resolve(address)
	if err != nil {
		s.log.Printf("dns error %s , ERR:%s", address, err)
		return address
	}
	return ip
}

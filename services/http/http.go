package http

import (
	"fmt"
	"io"
	"io/ioutil"
	logger "log"
	"net"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/snail007/goproxy/services"
	"github.com/snail007/goproxy/services/kcpcfg"
	"github.com/snail007/goproxy/utils"
	"github.com/snail007/goproxy/utils/conncrypt"

	"golang.org/x/crypto/ssh"
)

type HTTPArgs struct {
	Parent              *string
	CertFile            *string
	KeyFile             *string
	CaCertFile          *string
	CaCertBytes         []byte
	CertBytes           []byte
	KeyBytes            []byte
	Local               *string
	Always              *bool
	HTTPTimeout         *int
	Interval            *int
	Blocked             *string
	Direct              *string
	AuthFile            *string
	Auth                *[]string
	AuthURL             *string
	AuthURLOkCode       *int
	AuthURLTimeout      *int
	AuthURLRetry        *int
	ParentType          *string
	LocalType           *string
	Timeout             *int
	CheckParentInterval *int
	SSHKeyFile          *string
	SSHKeyFileSalt      *string
	SSHPassword         *string
	SSHUser             *string
	SSHKeyBytes         []byte
	SSHAuthMethod       ssh.AuthMethod
	KCP                 kcpcfg.KCPConfigArgs
	LocalIPS            *[]string
	DNSAddress          *string
	DNSTTL              *int
	LocalKey            *string
	ParentKey           *string
	LocalCompress       *bool
	ParentCompress      *bool
}
type HTTP struct {
	outPool        utils.OutConn
	cfg            HTTPArgs
	checker        utils.Checker
	basicAuth      utils.BasicAuth
	sshClient      *ssh.Client
	lockChn        chan bool
	domainResolver utils.DomainResolver
	isStop         bool
	serverChannels []*utils.ServerChannel
	userConns      utils.ConcurrentMap
	log            *logger.Logger
}

func NewHTTP() services.Service {
	return &HTTP{
		outPool:        utils.OutConn{},
		cfg:            HTTPArgs{},
		checker:        utils.Checker{},
		basicAuth:      utils.BasicAuth{},
		lockChn:        make(chan bool, 1),
		isStop:         false,
		serverChannels: []*utils.ServerChannel{},
		userConns:      utils.NewConcurrentMap(),
	}
}
func (s *HTTP) CheckArgs() (err error) {
	if *s.cfg.Parent != "" && *s.cfg.ParentType == "" {
		err = fmt.Errorf("parent type unkown,use -T <tls|tcp|ssh|kcp>")
		return
	}
	if *s.cfg.ParentType == "tls" || *s.cfg.LocalType == "tls" {
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
	return
}
func (s *HTTP) InitService() (err error) {
	s.InitBasicAuth()
	if *s.cfg.Parent != "" {
		s.checker = utils.NewChecker(*s.cfg.HTTPTimeout, int64(*s.cfg.Interval), *s.cfg.Blocked, *s.cfg.Direct, s.log)
	}
	if *s.cfg.DNSAddress != "" {
		(*s).domainResolver = utils.NewDomainResolver(*s.cfg.DNSAddress, *s.cfg.DNSTTL, s.log)
	}
	if *s.cfg.ParentType == "ssh" {
		err = s.ConnectSSH()
		if err != nil {
			err = fmt.Errorf("init service fail, ERR: %s", err)
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
						if s.sshClient.Conn != nil {
							s.sshClient.Conn.Close()
						}
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
	return
}
func (s *HTTP) StopService() {
	defer func() {
		e := recover()
		if e != nil {
			s.log.Printf("stop http(s) service crashed,%s", e)
		} else {
			s.log.Printf("service http(s) stoped")
		}
	}()
	s.isStop = true
	s.checker.Stop()
	if s.sshClient != nil {
		s.sshClient.Close()
	}
	for _, sc := range s.serverChannels {
		if sc.Listener != nil && *sc.Listener != nil {
			(*sc.Listener).Close()
		}
		if sc.UDPListener != nil {
			(*sc.UDPListener).Close()
		}
	}
}
func (s *HTTP) Start(args interface{}, log *logger.Logger) (err error) {
	s.log = log
	s.cfg = args.(HTTPArgs)
	if err = s.CheckArgs(); err != nil {
		return
	}

	if err = s.InitService(); err != nil {
		return
	}

	if *s.cfg.Parent != "" {
		s.log.Printf("use %s parent %s", *s.cfg.ParentType, *s.cfg.Parent)
		s.InitOutConnPool()
	}

	for _, addr := range strings.Split(*s.cfg.Local, ",") {
		if addr != "" {
			host, port, _ := net.SplitHostPort(addr)
			p, _ := strconv.Atoi(port)
			sc := utils.NewServerChannel(host, p, s.log)
			if *s.cfg.LocalType == "tcp" {
				err = sc.ListenTCP(s.callback)
			} else if *s.cfg.LocalType == "tls" {
				err = sc.ListenTls(s.cfg.CertBytes, s.cfg.KeyBytes, s.cfg.CaCertBytes, s.callback)
			} else if *s.cfg.LocalType == "kcp" {
				err = sc.ListenKCP(s.cfg.KCP, s.callback, s.log)
			}
			if err != nil {
				return
			}
			s.log.Printf("%s http(s) proxy on %s", *s.cfg.LocalType, (*sc.Listener).Addr())
			s.serverChannels = append(s.serverChannels, &sc)
		}
	}
	return
}

func (s *HTTP) Clean() {
	s.StopService()
}
func (s *HTTP) callback(inConn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			s.log.Printf("http(s) conn handler crashed with err : %s \nstack: %s", err, string(debug.Stack()))
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
	var err interface{}
	var req utils.HTTPRequest
	req, err = utils.NewHTTPRequest(&inConn, 4096, s.IsBasicAuth(), &s.basicAuth, s.log)
	if err != nil {
		if err != io.EOF {
			s.log.Printf("decoder error , from %s, ERR:%s", inConn.RemoteAddr(), err)
		}
		utils.CloseConn(&inConn)
		return
	}
	address := req.Host
	host, _, _ := net.SplitHostPort(address)
	useProxy := false
	if !utils.IsIternalIP(host, *s.cfg.Always) {
		useProxy = true
		if *s.cfg.Parent == "" {
			useProxy = false
		} else if *s.cfg.Always {
			useProxy = true
		} else {
			k := s.Resolve(address)
			s.checker.Add(address, k)
			//var n, m uint
			useProxy, _, _ = s.checker.IsBlocked(k)
			//s.log.Printf("blocked ? : %v, %s , fail:%d ,success:%d", useProxy, address, n, m)
		}
	}

	s.log.Printf("use proxy : %v, %s", useProxy, address)

	err = s.OutToTCP(useProxy, address, &inConn, &req)
	if err != nil {
		if *s.cfg.Parent == "" {
			s.log.Printf("connect to %s fail, ERR:%s", address, err)
		} else {
			s.log.Printf("connect to %s parent %s fail", *s.cfg.ParentType, *s.cfg.Parent)
		}
		utils.CloseConn(&inConn)
	}
}
func (s *HTTP) OutToTCP(useProxy bool, address string, inConn *net.Conn, req *utils.HTTPRequest) (err interface{}) {
	inAddr := (*inConn).RemoteAddr().String()
	inLocalAddr := (*inConn).LocalAddr().String()
	//防止死循环
	if s.IsDeadLoop(inLocalAddr, req.Host) {
		utils.CloseConn(inConn)
		err = fmt.Errorf("dead loop detected , %s", req.Host)
		return
	}
	var outConn net.Conn
	tryCount := 0
	maxTryCount := 5
	for {
		if s.isStop {
			return
		}
		if useProxy {
			if *s.cfg.ParentType == "ssh" {
				outConn, err = s.getSSHConn(address)
			} else {
				// s.log.Printf("%v", s.outPool)
				outConn, err = s.outPool.Get()
			}
		} else {
			outConn, err = utils.ConnectHost(s.Resolve(address), *s.cfg.Timeout)
		}
		tryCount++
		if err == nil || tryCount > maxTryCount {
			break
		} else {
			s.log.Printf("connect to %s , err:%s,retrying...", *s.cfg.Parent, err)
			time.Sleep(time.Second * 2)
		}
	}
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

	outAddr := outConn.RemoteAddr().String()
	//outLocalAddr := outConn.LocalAddr().String()
	if req.IsHTTPS() && (!useProxy || *s.cfg.ParentType == "ssh") {
		//https无上级或者上级非代理,proxy需要响应connect请求,并直连目标
		err = req.HTTPSReply()
	} else {
		//https或者http,上级是代理,proxy需要转发
		outConn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
		//直连目标或上级非代理,清理HTTP头部的代理头信息
		if !useProxy || *s.cfg.ParentType == "ssh" {
			_, err = outConn.Write(utils.RemoveProxyHeaders(req.HeadBuf))
		} else {
			_, err = outConn.Write(req.HeadBuf)
		}
		outConn.SetDeadline(time.Time{})
		if err != nil {
			s.log.Printf("write to %s , err:%s", *s.cfg.Parent, err)
			utils.CloseConn(inConn)
			return
		}
	}

	utils.IoBind((*inConn), outConn, func(err interface{}) {
		s.log.Printf("conn %s - %s released [%s]", inAddr, outAddr, req.Host)
		s.userConns.Remove(inAddr)
	}, s.log)
	s.log.Printf("conn %s - %s connected [%s]", inAddr, outAddr, req.Host)
	if c, ok := s.userConns.Get(inAddr); ok {
		(*c.(*net.Conn)).Close()
	}
	s.userConns.Set(inAddr, inConn)
	return
}

func (s *HTTP) getSSHConn(host string) (outConn net.Conn, err interface{}) {
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
	case <-time.After(time.Second * 5):
		err = fmt.Errorf("ssh dial %s timeout", host)
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
	return
}
func (s *HTTP) ConnectSSH() (err error) {
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
func (s *HTTP) InitOutConnPool() {
	if *s.cfg.ParentType == "tls" || *s.cfg.ParentType == "tcp" || *s.cfg.ParentType == "kcp" {
		//dur int, isTLS bool, certBytes, keyBytes []byte,
		//parent string, timeout int, InitialCap int, MaxCap int
		s.outPool = utils.NewOutConn(
			*s.cfg.CheckParentInterval,
			*s.cfg.ParentType,
			s.cfg.KCP,
			s.cfg.CertBytes, s.cfg.KeyBytes, s.cfg.CaCertBytes,
			s.Resolve(*s.cfg.Parent),
			*s.cfg.Timeout,
		)
	}
}
func (s *HTTP) InitBasicAuth() (err error) {
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
func (s *HTTP) IsBasicAuth() bool {
	return *s.cfg.AuthFile != "" || len(*s.cfg.Auth) > 0 || *s.cfg.AuthURL != ""
}
func (s *HTTP) IsDeadLoop(inLocalAddr string, host string) bool {
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
func (s *HTTP) Resolve(address string) string {
	if *s.cfg.DNSAddress == "" {
		return address
	}
	ip, err := s.domainResolver.Resolve(address)
	if err != nil {
		s.log.Printf("dns error %s , ERR:%s", address, err)
	}
	return ip
}

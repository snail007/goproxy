package services

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"proxy/utils"
	"runtime/debug"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
)

type HTTP struct {
	outPool   utils.OutPool
	cfg       HTTPArgs
	checker   utils.Checker
	basicAuth utils.BasicAuth
	sshClient *ssh.Client
	lockChn   chan bool
}

func NewHTTP() Service {
	return &HTTP{
		outPool:   utils.OutPool{},
		cfg:       HTTPArgs{},
		checker:   utils.Checker{},
		basicAuth: utils.BasicAuth{},
		lockChn:   make(chan bool, 1),
	}
}
func (s *HTTP) CheckArgs() {
	var err error
	if *s.cfg.Parent != "" && *s.cfg.ParentType == "" {
		log.Fatalf("parent type unkown,use -T <tls|tcp|ssh>")
	}
	if *s.cfg.ParentType == "tls" || *s.cfg.LocalType == "tls" {
		s.cfg.CertBytes, s.cfg.KeyBytes = utils.TlsBytes(*s.cfg.CertFile, *s.cfg.KeyFile)
	}
	if *s.cfg.ParentType == "ssh" {
		if *s.cfg.SSHUser == "" {
			log.Fatalf("ssh user required")
		}
		if *s.cfg.SSHKeyFile == "" && *s.cfg.SSHPassword == "" {
			log.Fatalf("ssh password or key required")
		}

		if *s.cfg.SSHPassword != "" {
			s.cfg.SSHAuthMethod = ssh.Password(*s.cfg.SSHPassword)
		} else {
			var SSHSigner ssh.Signer
			s.cfg.SSHKeyBytes, err = ioutil.ReadFile(*s.cfg.SSHKeyFile)
			if err != nil {
				log.Fatalf("read key file ERR: %s", err)
			}
			if *s.cfg.SSHKeyFileSalt != "" {
				SSHSigner, err = ssh.ParsePrivateKeyWithPassphrase(s.cfg.SSHKeyBytes, []byte(*s.cfg.SSHKeyFileSalt))
			} else {
				SSHSigner, err = ssh.ParsePrivateKey(s.cfg.SSHKeyBytes)
			}
			if err != nil {
				log.Fatalf("parse ssh private key fail,ERR: %s", err)
			}
			s.cfg.SSHAuthMethod = ssh.PublicKeys(SSHSigner)
		}
	}
}
func (s *HTTP) InitService() {
	s.InitBasicAuth()
	if *s.cfg.Parent != "" {
		s.checker = utils.NewChecker(*s.cfg.HTTPTimeout, int64(*s.cfg.Interval), *s.cfg.Blocked, *s.cfg.Direct)
	}
	if *s.cfg.ParentType == "ssh" {
		err := s.ConnectSSH()
		if err != nil {
			log.Fatalf("init service fail, ERR: %s", err)
		}
		go func() {
			//循环检查ssh网络连通性
			for {
				conn, err := utils.ConnectHost(*s.cfg.Parent, *s.cfg.Timeout*2)
				if err == nil {
					_, err = conn.Write([]byte{0})
				}
				if err != nil {
					if s.sshClient != nil {
						s.sshClient.Close()
						if s.sshClient.Conn != nil {
							s.sshClient.Conn.Close()
						}
					}
					log.Printf("ssh offline, retrying...")
					s.ConnectSSH()
				} else {
					conn.Close()
				}
				time.Sleep(time.Second * 3)
			}
		}()
	}
}
func (s *HTTP) StopService() {
	if s.outPool.Pool != nil {
		s.outPool.Pool.ReleaseAll()
	}
}
func (s *HTTP) Start(args interface{}) (err error) {
	s.cfg = args.(HTTPArgs)
	s.CheckArgs()
	if *s.cfg.Parent != "" {
		log.Printf("use %s parent %s", *s.cfg.ParentType, *s.cfg.Parent)
		s.InitOutConnPool()
	}
	s.InitService()
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
	log.Printf("%s http(s) proxy on %s", *s.cfg.LocalType, (*sc.Listener).Addr())
	return
}

func (s *HTTP) Clean() {
	s.StopService()
}
func (s *HTTP) callback(inConn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("http(s) conn handler crashed with err : %s \nstack: %s", err, string(debug.Stack()))
		}
	}()
	var err interface{}
	var req utils.HTTPRequest
	req, err = utils.NewHTTPRequest(&inConn, 4096, s.IsBasicAuth(), &s.basicAuth)
	if err != nil {
		if err != io.EOF {
			log.Printf("decoder error , form %s, ERR:%s", err, inConn.RemoteAddr())
		}
		utils.CloseConn(&inConn)
		return
	}
	address := req.Host

	useProxy := true
	if *s.cfg.Parent == "" {
		useProxy = false
	} else if *s.cfg.Always {
		useProxy = true
	} else {
		if req.IsHTTPS() {
			s.checker.Add(address, true, req.Method, "", nil)
		} else {
			s.checker.Add(address, false, req.Method, req.URL, req.HeadBuf)
		}
		//var n, m uint
		useProxy, _, _ = s.checker.IsBlocked(req.Host)
		//log.Printf("blocked ? : %v, %s , fail:%d ,success:%d", useProxy, address, n, m)
	}
	log.Printf("use proxy : %v, %s", useProxy, address)

	err = s.OutToTCP(useProxy, address, &inConn, &req)

	if err != nil {
		if *s.cfg.Parent == "" {
			log.Printf("connect to %s fail, ERR:%s", address, err)
		} else {
			log.Printf("connect to %s parent %s fail", *s.cfg.ParentType, *s.cfg.Parent)
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
	var _outConn interface{}
	tryCount := 0
	maxTryCount := 5
	for {
		if useProxy {
			if *s.cfg.ParentType == "ssh" {
				outConn, err = s.getSSHConn(address)
			} else {
				//log.Printf("%v", s.outPool)
				_outConn, err = s.outPool.Pool.Get()
				if err == nil {
					outConn = _outConn.(net.Conn)
				}
			}
		} else {
			outConn, err = utils.ConnectHost(address, *s.cfg.Timeout)
		}
		tryCount++
		if err == nil || tryCount > maxTryCount {
			break
		} else {
			log.Printf("connect to %s , err:%s,retrying...", *s.cfg.Parent, err)
			time.Sleep(time.Second * 2)
		}
	}
	if err != nil {
		log.Printf("connect to %s , err:%s", *s.cfg.Parent, err)
		utils.CloseConn(inConn)
		return
	}

	outAddr := outConn.RemoteAddr().String()
	//outLocalAddr := outConn.LocalAddr().String()

	if req.IsHTTPS() && (!useProxy || *s.cfg.ParentType == "ssh") {
		//https无上级或者上级非代理,proxy需要响应connect请求,并直连目标
		err = req.HTTPSReply()
	} else {
		//https或者http,上级是代理,proxy需要转发
		_, err = outConn.Write(req.HeadBuf)
		if err != nil {
			log.Printf("write to %s , err:%s", *s.cfg.Parent, err)
			utils.CloseConn(inConn)
			return
		}
	}

	utils.IoBind((*inConn), outConn, func(err interface{}) {
		log.Printf("conn %s - %s released [%s]", inAddr, outAddr, req.Host)
	})
	log.Printf("conn %s - %s connected [%s]", inAddr, outAddr, req.Host)

	return
}

func (s *HTTP) getSSHConn(host string) (outConn net.Conn, err interface{}) {
	maxTryCount := 1
	tryCount := 0
RETRY:
	if tryCount >= maxTryCount {
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
		log.Printf("connect ssh fail, ERR: %s, retrying...", err)
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
	s.sshClient, err = ssh.Dial("tcp", *s.cfg.Parent, &config)
	<-s.lockChn
	return
}
func (s *HTTP) InitOutConnPool() {
	if *s.cfg.ParentType == TYPE_TLS || *s.cfg.ParentType == TYPE_TCP || *s.cfg.ParentType == TYPE_KCP {
		//dur int, isTLS bool, certBytes, keyBytes []byte,
		//parent string, timeout int, InitialCap int, MaxCap int
		s.outPool = utils.NewOutPool(
			*s.cfg.CheckParentInterval,
			*s.cfg.ParentType,
			*s.cfg.KCPMethod,
			*s.cfg.KCPKey,
			s.cfg.CertBytes, s.cfg.KeyBytes,
			*s.cfg.Parent,
			*s.cfg.Timeout,
			*s.cfg.PoolSize,
			*s.cfg.PoolSize*2,
		)
	}
}
func (s *HTTP) InitBasicAuth() (err error) {
	s.basicAuth = utils.NewBasicAuth()
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
func (s *HTTP) IsBasicAuth() bool {
	return *s.cfg.AuthFile != "" || len(*s.cfg.Auth) > 0
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
		outIPs, err = net.LookupIP(outDomain)
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

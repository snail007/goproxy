package socks

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	logger "log"
	"net"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	server "github.com/snail007/goproxy/core/cs/server"
	"github.com/snail007/goproxy/core/lib/kcpcfg"
	"github.com/snail007/goproxy/services"
	"github.com/snail007/goproxy/utils"
	"github.com/snail007/goproxy/utils/conncrypt"
	"github.com/snail007/goproxy/utils/datasize"
	"github.com/snail007/goproxy/utils/dnsx"
	"github.com/snail007/goproxy/utils/iolimiter"
	"github.com/snail007/goproxy/utils/lb"
	"github.com/snail007/goproxy/utils/mapx"
	"github.com/snail007/goproxy/utils/socks"

	"golang.org/x/crypto/ssh"
)

type SocksArgs struct {
	Parent                *[]string
	ParentType            *string
	Local                 *string
	LocalType             *string
	CertFile              *string
	KeyFile               *string
	CaCertFile            *string
	CaCertBytes           []byte
	CertBytes             []byte
	KeyBytes              []byte
	SSHKeyFile            *string
	SSHKeyFileSalt        *string
	SSHPassword           *string
	SSHUser               *string
	SSHKeyBytes           []byte
	SSHAuthMethod         ssh.AuthMethod
	Timeout               *int
	Always                *bool
	Interval              *int
	Blocked               *string
	Direct                *string
	ParentAuth            *string
	AuthFile              *string
	Auth                  *[]string
	AuthURL               *string
	AuthURLOkCode         *int
	AuthURLTimeout        *int
	AuthURLRetry          *int
	KCP                   kcpcfg.KCPConfigArgs
	LocalIPS              *[]string
	DNSAddress            *string
	DNSTTL                *int
	LocalKey              *string
	ParentKey             *string
	LocalCompress         *bool
	ParentCompress        *bool
	Intelligent           *string
	LoadBalanceMethod     *string
	LoadBalanceTimeout    *int
	LoadBalanceRetryTime  *int
	LoadBalanceHashTarget *bool
	LoadBalanceOnlyHA     *bool

	RateLimit      *string
	RateLimitBytes float64
	BindListen     *bool
	Debug          *bool
}
type Socks struct {
	cfg                   SocksArgs
	checker               utils.Checker
	basicAuth             utils.BasicAuth
	sshClient             *ssh.Client
	lockChn               chan bool
	udpSC                 server.ServerChannel
	sc                    *server.ServerChannel
	domainResolver        dnsx.DomainResolver
	isStop                bool
	userConns             mapx.ConcurrentMap
	log                   *logger.Logger
	lb                    *lb.Group
	udpRelatedPacketConns mapx.ConcurrentMap
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
		userConns:             mapx.NewConcurrentMap(),
		udpRelatedPacketConns: mapx.NewConcurrentMap(),
	}
}

func (s *Socks) CheckArgs() (err error) {

	if *s.cfg.LocalType == "tls" || (len(*s.cfg.Parent) > 0 && *s.cfg.ParentType == "tls") {
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

	if len(*s.cfg.Parent) == 1 && (*s.cfg.Parent)[0] == "" {
		(*s.cfg.Parent) = []string{}
	}

	if len(*s.cfg.Parent) > 0 {
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
	if *s.cfg.RateLimit != "0" && *s.cfg.RateLimit != "" {
		var size uint64
		size, err = datasize.Parse(*s.cfg.RateLimit)
		if err != nil {
			err = fmt.Errorf("parse rate limit size error,ERR:%s", err)
			return
		}
		s.cfg.RateLimitBytes = float64(size)
	}
	s.udpLocalKey = s.LocalUDPKey()
	s.udpParentKey = s.ParentUDPKey()
	return
}
func (s *Socks) InitService() (err error) {
	s.InitBasicAuth()
	if *s.cfg.DNSAddress != "" {
		(*s).domainResolver = dnsx.NewDomainResolver(*s.cfg.DNSAddress, *s.cfg.DNSTTL, s.log)
	}
	if len(*s.cfg.Parent) > 0 {
		s.checker = utils.NewChecker(*s.cfg.Timeout, int64(*s.cfg.Interval), *s.cfg.Blocked, *s.cfg.Direct, s.log, *s.cfg.Intelligent)
		s.InitLB()
	}
	if *s.cfg.ParentType == "ssh" {
		e := s.ConnectSSH(s.Resolve(s.lb.Select("", *s.cfg.LoadBalanceOnlyHA)))
		if e != nil {
			err = fmt.Errorf("init service fail, ERR: %s", e)
			return
		}
		go func() {
			defer func() {
				if e := recover(); e != nil {
					fmt.Printf("crashed, err: %s\nstack:\n%s", e, string(debug.Stack()))
				}
			}()
			//循环检查ssh网络连通性
			for {
				if s.isStop {
					return
				}
				conn, err := utils.ConnectHost(s.Resolve(s.lb.Select("", *s.cfg.LoadBalanceOnlyHA)), *s.cfg.Timeout*2)
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
					s.ConnectSSH(s.Resolve(s.lb.Select("", *s.cfg.LoadBalanceOnlyHA)))
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
			s.log.Printf("service socks stopped")
		}
		s.basicAuth = utils.BasicAuth{}
		s.cfg = SocksArgs{}
		s.checker = utils.Checker{}
		s.domainResolver = dnsx.DomainResolver{}
		s.lb = nil
		s.lockChn = nil
		s.log = nil
		s.sc = nil
		s.sshClient = nil
		s.udpLocalKey = nil
		s.udpParentKey = nil
		s.udpRelatedPacketConns = nil
		s.udpSC = server.ServerChannel{}
		s.userConns = nil
		s = nil
	}()
	s.isStop = true
	if len(*s.cfg.Parent) > 0 {
		s.checker.Stop()
	}
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
	if s.lb != nil {
		s.lb.Stop()
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
	if len(*s.cfg.Parent) > 0 {
		s.log.Printf("use %s parent %v [ %s ]", *s.cfg.ParentType, *s.cfg.Parent, strings.ToUpper(*s.cfg.LoadBalanceMethod))
	}
	sc := server.NewServerChannelHost(*s.cfg.Local, s.log)
	if *s.cfg.LocalType == "tcp" {
		err = sc.ListenTCP(s.socksConnCallback)
	} else if *s.cfg.LocalType == "tls" {
		err = sc.ListenTLS(s.cfg.CertBytes, s.cfg.KeyBytes, s.cfg.CaCertBytes, s.socksConnCallback)
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

	//socks5 server
	var serverConn *socks.ServerConn
	udpIP, _, _ := net.SplitHostPort(inConn.LocalAddr().String())
	if s.IsBasicAuth() {
		serverConn = socks.NewServerConn(&inConn, time.Millisecond*time.Duration(*s.cfg.Timeout), &s.basicAuth, true, udpIP, nil)
	} else {
		serverConn = socks.NewServerConn(&inConn, time.Millisecond*time.Duration(*s.cfg.Timeout), nil, true, udpIP, nil)
	}
	if err := serverConn.Handshake(); err != nil {
		if !strings.HasSuffix(err.Error(), "EOF") {
			s.log.Printf("handshake fail, ERR: %s", err)
		}
		inConn.Close()
		return
	}
	if serverConn.IsUDP() {
		s.proxyUDP(&inConn, serverConn)
	} else if serverConn.IsTCP() {
		s.proxyTCP(&inConn, serverConn)
	}
}

func (s *Socks) proxyTCP(inConn *net.Conn, serverConn *socks.ServerConn) {
	var outConn net.Conn
	var err interface{}
	lbAddr := ""
	useProxy := true
	tryCount := 0
	maxTryCount := 5
	//防止死循环
	if s.IsDeadLoop((*inConn).LocalAddr().String(), serverConn.Host()) {
		utils.CloseConn(inConn)
		s.log.Printf("dead loop detected , %s", serverConn.Host())
		utils.CloseConn(inConn)
		return
	}
	for {
		if s.isStop {
			return
		}

		if *s.cfg.Always {
			selectAddr := (*inConn).RemoteAddr().String()
			if utils.LBMethod(*s.cfg.LoadBalanceMethod) == lb.SELECT_HASH && *s.cfg.LoadBalanceHashTarget {
				selectAddr = serverConn.Target()
			}
			lbAddr = s.lb.Select(selectAddr, *s.cfg.LoadBalanceOnlyHA)
			//lbAddr = s.lb.Select((*inConn).RemoteAddr().String())
			outConn, err = s.GetParentConn(lbAddr, serverConn)
			if err != nil {
				s.log.Printf("connect to parent fail, %s", err)
				return
			}
			//handshake
			//socks client
			_, err = s.HandshakeSocksParent(&outConn, "tcp", serverConn.Target(), serverConn.AuthData(), false)
			if err != nil {
				if err != io.EOF {
					s.log.Printf("handshake fail, %s", err)
				}
				return
			}
		} else {
			if len(*s.cfg.Parent) > 0 {
				host, _, _ := net.SplitHostPort(serverConn.Target())
				useProxy := false
				if utils.IsInternalIP(host, *s.cfg.Always) {
					useProxy = false
				} else {
					var isInMap bool
					useProxy, isInMap, _, _ = s.checker.IsBlocked(serverConn.Target())
					if !isInMap {
						s.checker.Add(serverConn.Target(), s.Resolve(serverConn.Target()))
					}
				}
				if useProxy {
					selectAddr := (*inConn).RemoteAddr().String()
					if utils.LBMethod(*s.cfg.LoadBalanceMethod) == lb.SELECT_HASH && *s.cfg.LoadBalanceHashTarget {
						selectAddr = serverConn.Target()
					}
					lbAddr = s.lb.Select(selectAddr, *s.cfg.LoadBalanceOnlyHA)
					//lbAddr = s.lb.Select((*inConn).RemoteAddr().String())
					outConn, err = s.GetParentConn(lbAddr, serverConn)
					if err != nil {
						s.log.Printf("connect to parent fail, %s", err)
						return
					}
					//handshake
					//socks client
					_, err = s.HandshakeSocksParent(&outConn, "tcp", serverConn.Target(), serverConn.AuthData(), false)
					if err != nil {
						s.log.Printf("handshake fail, %s", err)
						return
					}
				} else {
					outConn, err = s.GetDirectConn(s.Resolve(serverConn.Target()), (*inConn).LocalAddr().String())
				}
			} else {
				outConn, err = s.GetDirectConn(s.Resolve(serverConn.Target()), (*inConn).LocalAddr().String())
				useProxy = false
			}
		}
		tryCount++
		if err == nil || tryCount > maxTryCount || len(*s.cfg.Parent) == 0 {
			break
		} else {
			s.log.Printf("get out conn fail,%s,retrying...", err)
			time.Sleep(time.Second * 2)
		}
	}
	if err != nil {
		s.log.Printf("get out conn fail,%s", err)
		return
	}

	s.log.Printf("use proxy %v : %s", useProxy, serverConn.Target())

	inAddr := (*inConn).RemoteAddr().String()
	//outRemoteAddr := outConn.RemoteAddr().String()
	//inLocalAddr := (*inConn).LocalAddr().String()

	if s.cfg.RateLimitBytes > 0 {
		outConn = iolimiter.NewReaderConn(outConn, s.cfg.RateLimitBytes)
	}

	utils.IoBind(*inConn, outConn, func(err interface{}) {
		s.log.Printf("conn %s - %s released", inAddr, serverConn.Target())
		s.userConns.Remove(inAddr)
		if len(*s.cfg.Parent) > 0 {
			s.lb.DecreaseConns(lbAddr)
		}
	}, s.log)
	if c, ok := s.userConns.Get(inAddr); ok {
		(*c.(*net.Conn)).Close()
		s.userConns.Remove(inAddr)
	}
	s.userConns.Set(inAddr, inConn)
	if len(*s.cfg.Parent) > 0 {
		s.lb.IncreasConns(lbAddr)
	}
	s.log.Printf("conn %s - %s connected", inAddr, serverConn.Target())
}
func (s *Socks) GetParentConn(parentAddress string, serverConn *socks.ServerConn) (outConn net.Conn, err interface{}) {
	switch *s.cfg.ParentType {
	case "kcp", "tls", "tcp":
		if *s.cfg.ParentType == "tls" {
			var _conn tls.Conn
			_conn, err = utils.TlsConnectHost(parentAddress, *s.cfg.Timeout, s.cfg.CertBytes, s.cfg.KeyBytes, s.cfg.CaCertBytes)
			if err == nil {
				outConn = net.Conn(&_conn)
			}
		} else if *s.cfg.ParentType == "kcp" {
			outConn, err = utils.ConnectKCPHost(parentAddress, s.cfg.KCP)
		} else {
			outConn, err = utils.ConnectHost(parentAddress, *s.cfg.Timeout)
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
			outConn, err = s.sshClient.Dial("tcp", serverConn.Target())
		}()
		select {
		case <-wait:
		case <-time.After(time.Millisecond * time.Duration(*s.cfg.Timeout) * 2):
			err = fmt.Errorf("ssh dial %s timeout", serverConn.Target())
			s.sshClient.Close()
		}
		if err != nil {
			s.log.Printf("connect ssh fail, ERR: %s, retrying...", err)
			e := s.ConnectSSH(parentAddress)
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
func (s *Socks) ConnectSSH(lbAddr string) (err error) {
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
	s.sshClient, err = ssh.Dial("tcp", s.Resolve(lbAddr), &config)
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
func (s *Socks) InitLB() {
	configs := lb.BackendsConfig{}
	for _, addr := range *s.cfg.Parent {
		_addrInfo := strings.Split(addr, "@")
		_addr := _addrInfo[0]
		weight := 1
		if len(_addrInfo) == 2 {
			weight, _ = strconv.Atoi(_addrInfo[1])
		}
		configs = append(configs, &lb.BackendConfig{
			Address:       _addr,
			Weight:        weight,
			ActiveAfter:   1,
			InactiveAfter: 2,
			Timeout:       time.Duration(*s.cfg.LoadBalanceTimeout) * time.Millisecond,
			RetryTime:     time.Duration(*s.cfg.LoadBalanceRetryTime) * time.Millisecond,
		})
	}
	LB := lb.NewGroup(utils.LBMethod(*s.cfg.LoadBalanceMethod), configs, &s.domainResolver, s.log, *s.cfg.Debug)
	s.lb = &LB
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
			outIPs, err = utils.LookupIP(outDomain)
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
func (s *Socks) GetDirectConn(address string, localAddr string) (conn net.Conn, err error) {
	if !*s.cfg.BindListen {
		return utils.ConnectHost(address, *s.cfg.Timeout)
	}
	ip, _, _ := net.SplitHostPort(localAddr)
	if utils.IsInternalIP(ip, false) {
		return utils.ConnectHost(address, *s.cfg.Timeout)
	}
	local, _ := net.ResolveTCPAddr("tcp", ip+":0")
	d := net.Dialer{
		Timeout:   time.Millisecond * time.Duration(*s.cfg.Timeout),
		LocalAddr: local,
	}
	conn, err = d.Dial("tcp", address)
	return
}
func (s *Socks) HandshakeSocksParent(outconn *net.Conn, network, dstAddr string, auth socks.Auth, fromSS bool) (client *socks.ClientConn, err error) {
	if *s.cfg.ParentAuth != "" {
		a := strings.Split(*s.cfg.ParentAuth, ":")
		if len(a) != 2 {
			err = fmt.Errorf("parent auth data format error")
			return
		}
		client = socks.NewClientConn(outconn, network, dstAddr, time.Millisecond*time.Duration(*s.cfg.Timeout), &socks.Auth{User: a[0], Password: a[1]}, nil)
	} else {
		if !fromSS && !s.IsBasicAuth() && auth.Password != "" && auth.User != "" {
			client = socks.NewClientConn(outconn, network, dstAddr, time.Millisecond*time.Duration(*s.cfg.Timeout), &auth, nil)
		} else {
			client = socks.NewClientConn(outconn, network, dstAddr, time.Millisecond*time.Duration(*s.cfg.Timeout), nil, nil)
		}
	}
	err = client.Handshake()
	return
}

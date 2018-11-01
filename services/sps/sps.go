package sps

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
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

	"github.com/snail007/goproxy/core/cs/server"
	"github.com/snail007/goproxy/core/lib/kcpcfg"
	"github.com/snail007/goproxy/services"
	"github.com/snail007/goproxy/utils"
	"github.com/snail007/goproxy/utils/conncrypt"
	"github.com/snail007/goproxy/utils/datasize"
	"github.com/snail007/goproxy/utils/dnsx"
	"github.com/snail007/goproxy/utils/iolimiter"
	"github.com/snail007/goproxy/utils/jumper"
	"github.com/snail007/goproxy/utils/lb"
	"github.com/snail007/goproxy/utils/mapx"
	"github.com/snail007/goproxy/utils/sni"
	"github.com/snail007/goproxy/utils/socks"
	"github.com/snail007/goproxy/utils/ss"
)

type SPSArgs struct {
	Parent                *[]string
	CertFile              *string
	KeyFile               *string
	CaCertFile            *string
	CaCertBytes           []byte
	CertBytes             []byte
	KeyBytes              []byte
	Local                 *string
	ParentType            *string
	LocalType             *string
	Timeout               *int
	KCP                   kcpcfg.KCPConfigArgs
	ParentServiceType     *string
	DNSAddress            *string
	DNSTTL                *int
	AuthFile              *string
	Auth                  *[]string
	AuthURL               *string
	AuthURLOkCode         *int
	AuthURLTimeout        *int
	AuthURLRetry          *int
	LocalIPS              *[]string
	ParentAuth            *string
	LocalKey              *string
	ParentKey             *string
	LocalCompress         *bool
	ParentCompress        *bool
	SSMethod              *string
	SSKey                 *string
	ParentSSMethod        *string
	ParentSSKey           *string
	DisableHTTP           *bool
	DisableSocks5         *bool
	DisableSS             *bool
	LoadBalanceMethod     *string
	LoadBalanceTimeout    *int
	LoadBalanceRetryTime  *int
	LoadBalanceHashTarget *bool
	LoadBalanceOnlyHA     *bool

	RateLimit      *string
	RateLimitBytes float64
	Debug          *bool
	Jumper         *string
}
type SPS struct {
	cfg                   SPSArgs
	domainResolver        dnsx.DomainResolver
	basicAuth             utils.BasicAuth
	serverChannels        []*server.ServerChannel
	userConns             mapx.ConcurrentMap
	log                   *logger.Logger
	localCipher           *ss.Cipher
	parentCipher          *ss.Cipher
	udpRelatedPacketConns mapx.ConcurrentMap
	lb                    *lb.Group
	udpLocalKey           []byte
	udpParentKey          []byte
	jumper                *jumper.Jumper
}

func NewSPS() services.Service {
	return &SPS{
		cfg:                   SPSArgs{},
		basicAuth:             utils.BasicAuth{},
		serverChannels:        []*server.ServerChannel{},
		userConns:             mapx.NewConcurrentMap(),
		udpRelatedPacketConns: mapx.NewConcurrentMap(),
	}
}
func (s *SPS) CheckArgs() (err error) {

	if len(*s.cfg.Parent) == 1 && (*s.cfg.Parent)[0] == "" {
		(*s.cfg.Parent) = []string{}
	}

	if len(*s.cfg.Parent) == 0 {
		err = fmt.Errorf("parent required for %s %s", *s.cfg.LocalType, *s.cfg.Local)
		return
	}
	if *s.cfg.ParentType == "" {
		err = fmt.Errorf("parent type unkown,use -T <tls|tcp|kcp>")
		return
	}
	if *s.cfg.ParentType == "ss" && (*s.cfg.ParentSSKey == "" || *s.cfg.ParentSSMethod == "") {
		err = fmt.Errorf("ss parent need a ss key, set it by : -J <sskey>")
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
	if *s.cfg.Jumper != "" {
		if *s.cfg.ParentType != "tls" && *s.cfg.ParentType != "tcp" {
			err = fmt.Errorf("jumper only worked of -T is tls or tcp")
			return
		}
		var j jumper.Jumper
		j, err = jumper.New(*s.cfg.Jumper, time.Millisecond*time.Duration(*s.cfg.Timeout))
		if err != nil {
			err = fmt.Errorf("parse jumper fail, err %s", err)
			return
		}
		s.jumper = &j
	}
	return
}
func (s *SPS) InitService() (err error) {

	if *s.cfg.DNSAddress != "" {
		s.domainResolver = dnsx.NewDomainResolver(*s.cfg.DNSAddress, *s.cfg.DNSTTL, s.log)
	}

	if len(*s.cfg.Parent) > 0 {
		s.InitLB()
	}

	err = s.InitBasicAuth()
	if *s.cfg.SSMethod != "" && *s.cfg.SSKey != "" {
		s.localCipher, err = ss.NewCipher(*s.cfg.SSMethod, *s.cfg.SSKey)
		if err != nil {
			s.log.Printf("error generating cipher : %s", err)
			return
		}
	}
	if *s.cfg.ParentServiceType == "ss" {
		s.parentCipher, err = ss.NewCipher(*s.cfg.ParentSSMethod, *s.cfg.ParentSSKey)
		if err != nil {
			s.log.Printf("error generating cipher : %s", err)
			return
		}
	}
	return
}

func (s *SPS) StopService() {
	defer func() {
		e := recover()
		if e != nil {
			s.log.Printf("stop sps service crashed,%s", e)
		} else {
			s.log.Printf("service sps stopped")
		}
		s.basicAuth = utils.BasicAuth{}
		s.cfg = SPSArgs{}
		s.domainResolver = dnsx.DomainResolver{}
		s.lb = nil
		s.localCipher = nil
		s.jumper = nil
		s.log = nil
		s.parentCipher = nil
		s.serverChannels = nil
		s.udpLocalKey = nil
		s.udpParentKey = nil
		s.udpRelatedPacketConns = nil
		s.userConns = nil
		s = nil
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
		if _, ok := c.(*net.Conn); ok {
			(*c.(*net.Conn)).Close()
		}
		if _, ok := c.(**net.Conn); ok {
			(*(*c.(**net.Conn))).Close()
		}
	}
	if s.lb != nil {
		s.lb.Stop()
	}
	for _, c := range s.udpRelatedPacketConns.Items() {
		(*c.(*net.UDPConn)).Close()
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

	s.log.Printf("use %s %s parent %v [ %s ]", *s.cfg.ParentType, *s.cfg.ParentServiceType, *s.cfg.Parent, strings.ToUpper(*s.cfg.LoadBalanceMethod))
	for _, addr := range strings.Split(*s.cfg.Local, ",") {
		if addr != "" {
			host, port, _ := net.SplitHostPort(addr)
			p, _ := strconv.Atoi(port)
			sc := server.NewServerChannel(host, p, s.log)
			s.serverChannels = append(s.serverChannels, &sc)
			if *s.cfg.LocalType == "tcp" {
				err = sc.ListenTCP(s.callback)
			} else if *s.cfg.LocalType == "tls" {
				err = sc.ListenTLS(s.cfg.CertBytes, s.cfg.KeyBytes, s.cfg.CaCertBytes, s.callback)
			} else if *s.cfg.LocalType == "kcp" {
				err = sc.ListenKCP(s.cfg.KCP, s.callback, s.log)
			}
			if *s.cfg.ParentServiceType == "socks" {
				err = s.RunSSUDP(addr)
			} else {
				s.log.Println("warn : udp only for socks parent ")
			}
			if err != nil {
				return
			}
			s.log.Printf("%s http(s)+socks+ss proxy on %s", *s.cfg.LocalType, (*sc.Listener).Addr())
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
			s.log.Printf("%s conn handler crashed with err : %s \nstack: %s", *s.cfg.LocalType, err, string(debug.Stack()))
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
	lbAddr := ""
	switch *s.cfg.ParentType {
	case "kcp", "tcp", "tls":
		lbAddr, err = s.OutToTCP(&inConn)
	default:
		err = fmt.Errorf("unkown parent type %s", *s.cfg.ParentType)
	}
	if err != nil {
		s.log.Printf("connect to %s parent %s fail, ERR:%s from %s", *s.cfg.ParentType, lbAddr, err, inConn.RemoteAddr())
		utils.CloseConn(&inConn)
	}
}
func (s *SPS) OutToTCP(inConn *net.Conn) (lbAddr string, err error) {
	enableUDP := *s.cfg.ParentServiceType == "socks"
	udpIP, _, _ := net.SplitHostPort((*inConn).LocalAddr().String())
	if len(*s.cfg.LocalIPS) > 0 {
		udpIP = (*s.cfg.LocalIPS)[0]
	}
	bInConn := utils.NewBufferedConn(*inConn)
	//important
	//action read will regist read event to system,
	//when data arrived , system call process
	//so that we can get buffered bytes count
	//otherwise Buffered() always return 0
	bInConn.ReadByte()
	bInConn.UnreadByte()

	n := 2048
	if n > bInConn.Buffered() {
		n = bInConn.Buffered()
	}
	h, err := bInConn.Peek(n)
	if err != nil {
		s.log.Printf("peek error %s ", err)
		(*inConn).Close()
		return
	}
	isSNI, _ := sni.ServerNameFromBytes(h)
	*inConn = bInConn
	address := ""
	var auth = socks.Auth{}
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
			serverConn = socks.NewServerConn(inConn, time.Millisecond*time.Duration(*s.cfg.Timeout), &s.basicAuth, enableUDP, udpIP, nil)
		} else {
			serverConn = socks.NewServerConn(inConn, time.Millisecond*time.Duration(*s.cfg.Timeout), nil, enableUDP, udpIP, nil)
		}
		if err = serverConn.Handshake(); err != nil {
			return
		}
		address = serverConn.Target()
		auth = serverConn.AuthData()
		if serverConn.IsUDP() {
			s.proxyUDP(inConn, serverConn)
			return
		}
	} else if utils.IsHTTP(h) || isSNI != "" {
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
		ssConn := ss.NewConn(*inConn, s.localCipher.Copy())
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
	selectAddr := (*inConn).RemoteAddr().String()
	if utils.LBMethod(*s.cfg.LoadBalanceMethod) == lb.SELECT_HASH && *s.cfg.LoadBalanceHashTarget {
		selectAddr = address
	}
	lbAddr = s.lb.Select(selectAddr, *s.cfg.LoadBalanceOnlyHA)
	outConn, err = s.GetParentConn(lbAddr)
	if err != nil {
		s.log.Printf("connect to %s , err:%s", lbAddr, err)
		utils.CloseConn(inConn)
		return
	}

	if *s.cfg.ParentAuth != "" || *s.cfg.ParentSSKey != "" || s.IsBasicAuth() {
		forwardBytes = utils.RemoveProxyHeaders(forwardBytes)
	}

	//ask parent for connect to target address
	if *s.cfg.ParentServiceType == "http" {
		//http parent
		isHTTPS := false

		pb := new(bytes.Buffer)
		if len(forwardBytes) == 0 {
			isHTTPS = true
			pb.Write([]byte(fmt.Sprintf("CONNECT %s HTTP/1.1\r\n", address)))
		}
		pb.WriteString(fmt.Sprintf("Host: %s\r\n", address))
		pb.WriteString(fmt.Sprintf("Proxy-Host: %s\r\n", address))
		pb.WriteString("Proxy-Connection: Keep-Alive\r\n")
		pb.WriteString("Connection: Keep-Alive\r\n")

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
			pb.Write([]byte(fmt.Sprintf("Proxy-Authorization: Basic %s\r\n", base64.StdEncoding.EncodeToString([]byte(u)))))
		}

		if isHTTPS {
			pb.Write([]byte("\r\n"))
		} else {
			forwardBytes = utils.InsertProxyHeaders(forwardBytes, string(pb.Bytes()))
			pb.Reset()
			pb.Write(forwardBytes)
			forwardBytes = nil
		}

		outConn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
		_, err = outConn.Write(pb.Bytes())
		outConn.SetDeadline(time.Time{})
		if err != nil {
			s.log.Printf("write CONNECT to %s , err:%s", lbAddr, err)
			utils.CloseConn(inConn)
			utils.CloseConn(&outConn)
			return
		}

		if isHTTPS {
			reply := make([]byte, 1024)
			outConn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
			_, err = outConn.Read(reply)
			outConn.SetDeadline(time.Time{})
			if err != nil {
				s.log.Printf("read reply from %s , err:%s", lbAddr, err)
				utils.CloseConn(inConn)
				utils.CloseConn(&outConn)
				return
			}
			//s.log.Printf("reply: %s", string(reply[:n]))
		}
	} else if *s.cfg.ParentServiceType == "socks" {
		s.log.Printf("connect %s", address)

		//socks client
		_, err = s.HandshakeSocksParent(&outConn, "tcp", address, auth, false)
		if err != nil {
			s.log.Printf("handshake fail, %s", err)
			return
		}

	} else if *s.cfg.ParentServiceType == "ss" {
		ra, e := ss.RawAddr(address)
		if e != nil {
			err = fmt.Errorf("build ss raw addr fail, err: %s", e)
			return
		}

		outConn, err = ss.DialWithRawAddr(&outConn, ra, "", s.parentCipher.Copy())
		if err != nil {
			err = fmt.Errorf("dial ss parent fail, err : %s", err)
			return
		}
	}

	//forward client data to target,if necessary.
	if len(forwardBytes) > 0 {
		outConn.Write(forwardBytes)
	}

	if s.cfg.RateLimitBytes > 0 {
		outConn = iolimiter.NewReaderConn(outConn, s.cfg.RateLimitBytes)
	}

	//bind
	inAddr := (*inConn).RemoteAddr().String()
	outAddr := outConn.RemoteAddr().String()
	utils.IoBind((*inConn), outConn, func(err interface{}) {
		s.log.Printf("conn %s - %s released [%s]", inAddr, outAddr, address)
		s.userConns.Remove(inAddr)
		s.lb.DecreaseConns(lbAddr)
	}, s.log)
	s.log.Printf("conn %s - %s connected [%s]", inAddr, outAddr, address)

	s.lb.IncreasConns(lbAddr)

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
func (s *SPS) InitLB() {
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
		return address
	}
	return ip
}
func (s *SPS) GetParentConn(address string) (conn net.Conn, err error) {
	if *s.cfg.ParentType == "tls" {
		if s.jumper == nil {
			var _conn tls.Conn
			_conn, err = utils.TlsConnectHost(address, *s.cfg.Timeout, s.cfg.CertBytes, s.cfg.KeyBytes, s.cfg.CaCertBytes)
			if err == nil {
				conn = net.Conn(&_conn)
			}
		} else {
			conf, err := utils.TlsConfig(s.cfg.CertBytes, s.cfg.KeyBytes, s.cfg.CaCertBytes)
			if err != nil {
				return nil, err
			}
			var _c net.Conn
			_c, err = s.jumper.Dial(address, time.Millisecond*time.Duration(*s.cfg.Timeout))
			if err == nil {
				conn = net.Conn(tls.Client(_c, conf))
			}
		}

	} else if *s.cfg.ParentType == "kcp" {
		conn, err = utils.ConnectKCPHost(address, s.cfg.KCP)
	} else {
		if s.jumper == nil {
			conn, err = utils.ConnectHost(address, *s.cfg.Timeout)
		} else {
			conn, err = s.jumper.Dial(address, time.Millisecond*time.Duration(*s.cfg.Timeout))
		}
	}
	if err == nil {
		if *s.cfg.ParentCompress {
			conn = utils.NewCompConn(conn)
		}
		if *s.cfg.ParentKey != "" {
			conn = conncrypt.New(conn, &conncrypt.Config{
				Password: *s.cfg.ParentKey,
			})
		}
	}
	return
}
func (s *SPS) HandshakeSocksParent(outconn *net.Conn, network, dstAddr string, auth socks.Auth, fromSS bool) (client *socks.ClientConn, err error) {
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
func (s *SPS) ParentUDPKey() (key []byte) {
	switch *s.cfg.ParentType {
	case "tcp":
		if *s.cfg.ParentKey != "" {
			v := fmt.Sprintf("%x", md5.Sum([]byte(*s.cfg.ParentKey)))
			return []byte(v)[:24]
		}
	case "tls":
		return s.cfg.KeyBytes[:24]
	case "kcp":
		v := fmt.Sprintf("%x", md5.Sum([]byte(*s.cfg.KCP.Key)))
		return []byte(v)[:24]
	}
	return
}
func (s *SPS) LocalUDPKey() (key []byte) {
	switch *s.cfg.LocalType {
	case "tcp":
		if *s.cfg.LocalKey != "" {
			v := fmt.Sprintf("%x", md5.Sum([]byte(*s.cfg.LocalKey)))
			return []byte(v)[:24]
		}
	case "tls":
		return s.cfg.KeyBytes[:24]
	case "kcp":
		v := fmt.Sprintf("%x", md5.Sum([]byte(*s.cfg.KCP.Key)))
		return []byte(v)[:24]
	}
	return
}

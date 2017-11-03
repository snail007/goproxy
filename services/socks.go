package services

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"proxy/utils"
	"proxy/utils/aes"
	"proxy/utils/socks"
	"runtime/debug"
	"time"

	"golang.org/x/crypto/ssh"
)

type Socks struct {
	cfg       SocksArgs
	checker   utils.Checker
	basicAuth utils.BasicAuth
	sshClient *ssh.Client
	lockChn   chan bool
	udpSC     utils.ServerChannel
}

func NewSocks() Service {
	return &Socks{
		cfg:       SocksArgs{},
		checker:   utils.Checker{},
		basicAuth: utils.BasicAuth{},
		lockChn:   make(chan bool, 1),
	}
}

func (s *Socks) CheckArgs() {
	var err error
	if *s.cfg.LocalType == "tls" {
		s.cfg.CertBytes, s.cfg.KeyBytes = utils.TlsBytes(*s.cfg.CertFile, *s.cfg.KeyFile)
	}
	if *s.cfg.Parent != "" {
		if *s.cfg.ParentType == "" {
			log.Fatalf("parent type unkown,use -T <tls|tcp|ssh>")
		}
		if *s.cfg.ParentType == "tls" {
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

}
func (s *Socks) InitService() {
	s.InitBasicAuth()
	s.checker = utils.NewChecker(*s.cfg.Timeout, int64(*s.cfg.Interval), *s.cfg.Blocked, *s.cfg.Direct)
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
	if *s.cfg.ParentType == "ssh" {
		log.Println("warn: socks udp not suppored for ssh")
	} else {

		s.udpSC = utils.NewServerChannelHost(*s.cfg.UDPLocal)
		err := s.udpSC.ListenUDP(s.udpCallback)
		if err != nil {
			log.Fatalf("init udp service fail, ERR: %s", err)
		}
		log.Printf("udp socks proxy on %s", s.udpSC.UDPListener.LocalAddr())
	}
}
func (s *Socks) StopService() {
	if s.sshClient != nil {
		s.sshClient.Close()
	}
	if s.udpSC.UDPListener != nil {
		s.udpSC.UDPListener.Close()
	}
}
func (s *Socks) Start(args interface{}) (err error) {
	//start()
	s.cfg = args.(SocksArgs)
	s.CheckArgs()
	s.InitService()
	if *s.cfg.Parent != "" {
		log.Printf("use %s parent %s", *s.cfg.ParentType, *s.cfg.Parent)
	}
	sc := utils.NewServerChannelHost(*s.cfg.Local)
	if *s.cfg.LocalType == TYPE_TCP {
		err = sc.ListenTCP(s.socksConnCallback)
	} else if *s.cfg.LocalType == TYPE_TLS {
		err = sc.ListenTls(s.cfg.CertBytes, s.cfg.KeyBytes, s.socksConnCallback)
	} else if *s.cfg.LocalType == TYPE_KCP {
		err = sc.ListenKCP(*s.cfg.KCPMethod, *s.cfg.KCPKey, s.socksConnCallback)
	}
	if err != nil {
		return
	}
	log.Printf("%s socks proxy on %s", *s.cfg.LocalType, (*sc.Listener).Addr())
	return
}
func (s *Socks) Clean() {
	s.StopService()
}
func (s *Socks) UDPKey() []byte {
	return s.cfg.KeyBytes[:32]
}
func (s *Socks) udpCallback(b []byte, localAddr, srcAddr *net.UDPAddr) {
	rawB := b
	var err error
	if *s.cfg.LocalType == "tls" {
		//decode b
		rawB, err = goaes.Decrypt(s.UDPKey(), b)
		if err != nil {
			log.Printf("decrypt udp packet fail from %s", srcAddr.String())
			return
		}
	}
	p, err := socks.ParseUDPPacket(rawB)
	log.Printf("udp revecived:%v", len(p.Data()))
	if err != nil {
		log.Printf("parse udp packet fail, ERR:%s", err)
		return
	}
	//防止死循环
	if s.IsDeadLoop((*localAddr).String(), p.Host()) {
		log.Printf("dead loop detected , %s", p.Host())
		return
	}
	//log.Printf("##########udp to -> %s:%s###########", p.Host(), p.Port())
	if *s.cfg.Parent != "" {
		//有上级代理,转发给上级
		if *s.cfg.ParentType == "tls" {
			//encode b
			rawB, err = goaes.Encrypt(s.UDPKey(), rawB)
			if err != nil {
				log.Printf("encrypt udp data fail to %s", *s.cfg.Parent)
				return
			}
		}
		parent := *s.cfg.UDPParent
		if parent == "" {
			parent = *s.cfg.Parent
		}
		dstAddr, err := net.ResolveUDPAddr("udp", parent)
		if err != nil {
			log.Printf("can't resolve address: %s", err)
			return
		}
		clientSrcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
		conn, err := net.DialUDP("udp", clientSrcAddr, dstAddr)
		if err != nil {
			log.Printf("connect to udp %s fail,ERR:%s", dstAddr.String(), err)
			return
		}
		conn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout*5)))
		_, err = conn.Write(rawB)
		log.Printf("udp request:%v", len(rawB))
		if err != nil {
			log.Printf("send udp packet to %s fail,ERR:%s", dstAddr.String(), err)
			conn.Close()
			return
		}

		//log.Printf("send udp packet to %s success", dstAddr.String())
		buf := make([]byte, 10*1024)
		length, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("read udp response from %s fail ,ERR:%s", dstAddr.String(), err)
			conn.Close()
			return
		}
		respBody := buf[0:length]
		log.Printf("udp response:%v", len(respBody))
		//log.Printf("revecived udp packet from %s", dstAddr.String())
		if *s.cfg.ParentType == "tls" {
			//decode b
			respBody, err = goaes.Decrypt(s.UDPKey(), respBody)
			if err != nil {
				log.Printf("encrypt udp data fail to %s", *s.cfg.Parent)
				conn.Close()
				return
			}
		}
		if *s.cfg.LocalType == "tls" {
			d, err := goaes.Encrypt(s.UDPKey(), respBody)
			if err != nil {
				log.Printf("encrypt udp data fail from %s", dstAddr.String())
				conn.Close()
				return
			}
			s.udpSC.UDPListener.WriteToUDP(d, srcAddr)
			log.Printf("udp reply:%v", len(d))
		} else {
			s.udpSC.UDPListener.WriteToUDP(respBody, srcAddr)
			log.Printf("udp reply:%v", len(respBody))
		}

	} else {
		//本地代理
		dstAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(p.Host(), p.Port()))
		if err != nil {
			log.Printf("can't resolve address: %s", err)
			return
		}
		clientSrcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
		conn, err := net.DialUDP("udp", clientSrcAddr, dstAddr)
		if err != nil {
			log.Printf("connect to udp %s fail,ERR:%s", dstAddr.String(), err)
			return
		}
		conn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout*3)))
		_, err = conn.Write(p.Data())
		log.Printf("udp send:%v", len(p.Data()))
		if err != nil {
			log.Printf("send udp packet to %s fail,ERR:%s", dstAddr.String(), err)
			conn.Close()
			return
		}
		//log.Printf("send udp packet to %s success", dstAddr.String())
		buf := make([]byte, 10*1024)
		length, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("read udp response from %s fail ,ERR:%s", dstAddr.String(), err)
			conn.Close()
			return
		}
		respBody := buf[0:length]
		//封装来自真实服务器的数据,返回给访问者
		respPacket := p.NewReply(respBody)
		//log.Printf("revecived udp packet from %s", dstAddr.String())
		if *s.cfg.LocalType == "tls" {
			d, err := goaes.Encrypt(s.UDPKey(), respPacket)
			if err != nil {
				log.Printf("encrypt udp data fail from %s", dstAddr.String())
				conn.Close()
				return
			}
			s.udpSC.UDPListener.WriteToUDP(d, srcAddr)
		} else {
			s.udpSC.UDPListener.WriteToUDP(respPacket, srcAddr)
		}
		log.Printf("udp reply:%v", len(respPacket))
	}

}
func (s *Socks) socksConnCallback(inConn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("socks conn handler crashed with err : %s \nstack: %s", err, string(debug.Stack()))
		}
		utils.CloseConn(&inConn)
	}()
	//协商开始

	//method select request
	inConn.SetReadDeadline(time.Now().Add(time.Second * 3))
	methodReq, err := socks.NewMethodsRequest(inConn)
	inConn.SetReadDeadline(time.Time{})
	if err != nil {
		methodReq.Reply(socks.Method_NONE_ACCEPTABLE)
		utils.CloseConn(&inConn)
		log.Printf("new methods request fail,ERR: %s", err)
		return
	}

	if !s.IsBasicAuth() {
		if !methodReq.Select(socks.Method_NO_AUTH) {
			methodReq.Reply(socks.Method_NONE_ACCEPTABLE)
			utils.CloseConn(&inConn)
			log.Printf("none method found : Method_NO_AUTH")
			return
		}
		//method select reply
		err = methodReq.Reply(socks.Method_NO_AUTH)
		if err != nil {
			log.Printf("reply answer data fail,ERR: %s", err)
			utils.CloseConn(&inConn)
			return
		}
		// log.Printf("% x", methodReq.Bytes())
	} else {
		//auth
		if !methodReq.Select(socks.Method_USER_PASS) {
			methodReq.Reply(socks.Method_NONE_ACCEPTABLE)
			utils.CloseConn(&inConn)
			log.Printf("none method found : Method_USER_PASS")
			return
		}
		//method reply need auth
		err = methodReq.Reply(socks.Method_USER_PASS)
		if err != nil {
			log.Printf("reply answer data fail,ERR: %s", err)
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
		//log.Printf("user:%s,pass:%s", user, pass)
		//auth
		if s.basicAuth.CheckUserPass(user, pass) {
			inConn.Write([]byte{0x01, 0x00})
		} else {
			inConn.Write([]byte{0x01, 0x01})
			utils.CloseConn(&inConn)
			return
		}
	}

	//request detail
	request, err := socks.NewRequest(inConn)
	if err != nil {
		log.Printf("read request data fail,ERR: %s", err)
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
func (s *Socks) proxyUDP(inConn *net.Conn, methodReq socks.MethodsRequest, request socks.Request) {
	if *s.cfg.ParentType == "ssh" {
		utils.CloseConn(inConn)
		return
	}
	host, _, _ := net.SplitHostPort((*inConn).LocalAddr().String())
	_, port, _ := net.SplitHostPort(s.udpSC.UDPListener.LocalAddr().String())
	log.Printf("proxy udp on %s", net.JoinHostPort(host, port))
	request.UDPReply(socks.REP_SUCCESS, net.JoinHostPort(host, port))
}
func (s *Socks) proxyTCP(inConn *net.Conn, methodReq socks.MethodsRequest, request socks.Request) {
	var outConn net.Conn
	defer utils.CloseConn(&outConn)
	var err interface{}
	useProxy := true
	tryCount := 0
	maxTryCount := 5
	//防止死循环
	if s.IsDeadLoop((*inConn).LocalAddr().String(), request.Host()) {
		utils.CloseConn(inConn)
		log.Printf("dead loop detected , %s", request.Host())
		utils.CloseConn(inConn)
		return
	}
	for {
		if *s.cfg.Always {
			outConn, err = s.getOutConn(methodReq.Bytes(), request.Bytes(), request.Addr())
		} else {
			if *s.cfg.Parent != "" {
				s.checker.Add(request.Addr(), true, "", "", nil)
				useProxy, _, _ = s.checker.IsBlocked(request.Addr())
				if useProxy {
					outConn, err = s.getOutConn(methodReq.Bytes(), request.Bytes(), request.Addr())
				} else {
					outConn, err = utils.ConnectHost(request.Addr(), *s.cfg.Timeout)
				}
			} else {
				outConn, err = utils.ConnectHost(request.Addr(), *s.cfg.Timeout)
				useProxy = false
			}
		}
		tryCount++
		if err == nil || tryCount > maxTryCount {
			break
		} else {
			log.Printf("get out conn fail,%s,retrying...", err)
			time.Sleep(time.Second * 2)
		}
	}
	if err != nil {
		log.Printf("get out conn fail,%s", err)
		request.TCPReply(socks.REP_NETWOR_UNREACHABLE)
		return
	}
	log.Printf("use proxy %v : %s", useProxy, request.Addr())

	request.TCPReply(socks.REP_SUCCESS)
	inAddr := (*inConn).RemoteAddr().String()
	//inLocalAddr := (*inConn).LocalAddr().String()

	log.Printf("conn %s - %s connected", inAddr, request.Addr())
	utils.IoBind(*inConn, outConn, func(err error) {
		log.Printf("conn %s - %s released %s", inAddr, request.Addr(), err)
		utils.CloseConn(inConn)
		utils.CloseConn(&outConn)
	})
}
func (s *Socks) getOutConn(methodBytes, reqBytes []byte, host string) (outConn net.Conn, err interface{}) {
	switch *s.cfg.ParentType {
	case "kcp":
		fallthrough
	case "tls":
		fallthrough
	case "tcp":
		if *s.cfg.ParentType == "tls" {
			var _outConn tls.Conn
			_outConn, err = utils.TlsConnectHost(*s.cfg.Parent, *s.cfg.Timeout, s.cfg.CertBytes, s.cfg.KeyBytes)
			outConn = net.Conn(&_outConn)
		} else if *s.cfg.ParentType == "kcp" {
			outConn, err = utils.ConnectKCPHost(*s.cfg.Parent, *s.cfg.KCPMethod, *s.cfg.KCPKey)
		} else {
			outConn, err = utils.ConnectHost(*s.cfg.Parent, *s.cfg.Timeout)
		}
		if err != nil {
			err = fmt.Errorf("connect fail,%s", err)
			return
		}
		var buf = make([]byte, 1024)
		//var n int
		_, err = outConn.Write(methodBytes)
		if err != nil {
			err = fmt.Errorf("write method fail,%s", err)
			return
		}
		_, err = outConn.Read(buf)
		if err != nil {
			err = fmt.Errorf("read method reply fail,%s", err)
			return
		}
		//resp := buf[:n]
		//log.Printf("resp:%v", resp)

		_, err = outConn.Write(reqBytes)
		if err != nil {
			err = fmt.Errorf("write req detail fail,%s", err)
			return
		}
		// _, err = outConn.Read(buf)
		// if err != nil {
		// 	err = fmt.Errorf("read req reply fail,%s", err)
		// 	return
		// }
		//result := buf[:n]
		//log.Printf("result:%v", result)

	case "ssh":
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
		case <-time.After(time.Millisecond * time.Duration(*s.cfg.Timeout) * 2):
			err = fmt.Errorf("ssh dial %s timeout", host)
			s.sshClient.Close()
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
	s.sshClient, err = ssh.Dial("tcp", *s.cfg.Parent, &config)
	<-s.lockChn
	return
}
func (s *Socks) InitBasicAuth() (err error) {
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
func (s *Socks) IsBasicAuth() bool {
	return *s.cfg.AuthFile != "" || len(*s.cfg.Auth) > 0
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

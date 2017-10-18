package services

import (
	"crypto/tls"
	"fmt"
	"io"
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
		log.Println(*s.cfg.CertFile, *s.cfg.KeyFile)
		s.cfg.CertBytes, s.cfg.KeyBytes = utils.TlsBytes(*s.cfg.CertFile, *s.cfg.KeyFile)
	}
	if *s.cfg.Parent != "" {
		if *s.cfg.ParentType == "" {
			log.Fatalf("parent type unkown,use -T <tls|tcp|ssh>")
		}
		if *s.cfg.ParentType == "tls" {
			log.Println(*s.cfg.CertFile, *s.cfg.KeyFile)
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
	s.checker = utils.NewChecker(*s.cfg.Timeout, int64(*s.cfg.Interval), *s.cfg.Blocked, *s.cfg.Direct)
	if *s.cfg.ParentType == "ssh" {
		err := s.ConnectSSH()
		if err != nil {
			log.Fatalf("init service fail, ERR: %s", err)
		}
	}
	if *s.cfg.ParentType == "ssh" {
		log.Println("warn: socks udp not suppored for ssh")
	} else {
		_, port, _ := net.SplitHostPort(*s.cfg.Local)
		s.udpSC = utils.NewServerChannelHost(":" + port)
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
	} else {
		err = sc.ListenTls(s.cfg.CertBytes, s.cfg.KeyBytes, s.socksConnCallback)
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
	newB := b
	var err error
	if *s.cfg.LocalType == "tls" {
		//decode b
		newB, err = goaes.Decrypt(s.UDPKey(), b)
		if err != nil {
			log.Printf("decrypt udp packet fail from %s", srcAddr.String())
			return
		}
	}
	p, err := socks.ParseUDPPacket(newB)
	log.Printf("udp revecived:%v", len(p.Data()))
	if err != nil {
		log.Printf("parse udp packet fail, ERR:%s", err)
		return
	}
	//log.Printf("##########udp to -> %s:%s###########", p.Host(), p.Port())
	if *s.cfg.Parent != "" {
		//有上级代理,转发给上级
		if *s.cfg.ParentType == "tls" {
			//encode b
			newB, err = goaes.Encrypt(s.UDPKey(), newB)
			if err != nil {
				log.Printf("encrypt udp data fail to %s", *s.cfg.Parent)
				return
			}
		}
		dstAddr, err := net.ResolveUDPAddr("udp", *s.cfg.Parent)
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
		conn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout*2)))
		_, err = conn.Write(newB)
		log.Printf("udp request:%v", len(newB))
		if err != nil {
			log.Printf("send udp packet to %s fail,ERR:%s", dstAddr.String(), err)
			return
		}

		//log.Printf("send udp packet to %s success", dstAddr.String())
		buf := make([]byte, 1024)
		length, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("read udp response from %s fail ,ERR:%s", dstAddr.String(), err)
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
				return
			}
		}
		if *s.cfg.LocalType == "tls" {
			d, err := goaes.Encrypt(s.UDPKey(), respBody)
			if err != nil {
				log.Printf("encrypt udp data fail from %s", dstAddr.String())
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
		conn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout*2)))
		_, err = conn.Write(p.Data())
		log.Printf("udp send:%v", len(p.Data()))
		if err != nil {
			log.Printf("send udp packet to %s fail,ERR:%s", dstAddr.String(), err)
			return
		}
		log.Printf("send udp packet to %s success", dstAddr.String())
		buf := make([]byte, 1024)
		length, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("read udp response from %s fail ,ERR:%s", dstAddr.String(), err)
			return
		}
		respBody := buf[0:length]
		//log.Printf("revecived udp packet from %s", dstAddr.String())
		if *s.cfg.LocalType == "tls" {
			d, err := goaes.Encrypt(s.UDPKey(), respBody)
			if err != nil {
				log.Printf("encrypt udp data fail from %s", dstAddr.String())
				return
			}
			s.udpSC.UDPListener.WriteToUDP(d, srcAddr)
		} else {
			s.udpSC.UDPListener.WriteToUDP(respBody, srcAddr)
		}
		log.Printf("udp reply:%v", len(respBody))
	}

}
func (s *Socks) socksConnCallback(inConn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("socks conn handler crashed with err : %s \nstack: %s", err, string(debug.Stack()))
		}
		utils.CloseConn(&inConn)
	}()

	//method select request
	methodReq, err := socks.NewMethodsRequest(inConn)
	if err != nil || !methodReq.Select(socks.Method_NO_AUTH) {
		methodReq.Reply(socks.Method_NONE_ACCEPTABLE)
		utils.CloseConn(&inConn)
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

	//request detail
	request, err := socks.NewRequest(inConn)
	if err != nil {
		log.Printf("read request data fail,ERR: %s", err)
		utils.CloseConn(&inConn)
		return
	}

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
		return
	}
	host, _, _ := net.SplitHostPort((*inConn).LocalAddr().String())
	_, port, _ := net.SplitHostPort(s.udpSC.UDPListener.LocalAddr().String())
	// log.Printf("proxy udp on %s", net.JoinHostPort(host, port))
	request.UDPReply(socks.REP_SUCCESS, net.JoinHostPort(host, port))
	// log.Printf("%v", request.NewReply(socks.REP_SUCCESS, net.JoinHostPort(host, port)))
}
func (s *Socks) proxyTCP(inConn *net.Conn, methodReq socks.MethodsRequest, request socks.Request) {
	var outConn net.Conn
	defer utils.CloseConn(&outConn)
	var err error
	useProxy := true
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
	inLocalAddr := (*inConn).LocalAddr().String()

	log.Printf("conn %s - %s connected [%s]", inAddr, inLocalAddr, request.Addr())
	var bind = func() (err interface{}) {
		defer func() {
			if err == nil {
				if err = recover(); err != nil {
					log.Printf("bind crashed %s", err)
				}
			}
		}()
		go func() {
			defer func() {
				if err == nil {
					if err = recover(); err != nil {
						log.Printf("bind crashed %s", err)
					}
				}
			}()
			_, err = io.Copy(outConn, (*inConn))
		}()
		_, err = io.Copy((*inConn), outConn)
		return
	}
	bind()
	log.Printf("conn %s - %s released [%s]", inAddr, inLocalAddr, request.Addr())
	utils.CloseConn(inConn)
	utils.CloseConn(&outConn)
}
func (s *Socks) getOutConn(methodBytes, reqBytes []byte, host string) (outConn net.Conn, err error) {
	switch *s.cfg.ParentType {
	case "tls":
		fallthrough
	case "tcp":
		if *s.cfg.ParentType == "tls" {
			var _outConn tls.Conn
			_outConn, err = utils.TlsConnectHost(*s.cfg.Parent, *s.cfg.Timeout, s.cfg.CertBytes, s.cfg.KeyBytes)
			outConn = net.Conn(&_outConn)
		} else {
			outConn, err = utils.ConnectHost(*s.cfg.Parent, *s.cfg.Timeout)
		}
		if err != nil {
			return
		}
		var buf = make([]byte, 1024)
		//var n int
		_, err = outConn.Write(methodBytes)
		if err != nil {
			return
		}
		_, err = outConn.Read(buf)
		if err != nil {
			return
		}
		//resp := buf[:n]
		//log.Printf("resp:%v", resp)

		outConn.Write(reqBytes)
		_, err = outConn.Read(buf)
		if err != nil {
			return
		}
		//result := buf[:n]
		//log.Printf("result:%v", result)

	case "ssh":
		maxTryCount := 1
		tryCount := 0
	RETRY:
		if tryCount >= maxTryCount {
			return
		}
		outConn, err = s.sshClient.Dial("tcp", host)
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

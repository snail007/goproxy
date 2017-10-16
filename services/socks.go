package services

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"proxy/utils"
	"time"

	"golang.org/x/crypto/ssh"
)

type Socks struct {
	cfg       SocksArgs
	checker   utils.Checker
	basicAuth utils.BasicAuth
	sshClient *ssh.Client
}

func NewSocks() Service {
	return &Socks{
		cfg:       SocksArgs{},
		checker:   utils.Checker{},
		basicAuth: utils.BasicAuth{},
	}
}

func (s *Socks) CheckArgs() {
	var err error
	if *s.cfg.Parent != "" {
		if *s.cfg.ParentType == "" {
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

}
func (s *Socks) InitService() {
	s.checker = utils.NewChecker(*s.cfg.Timeout, int64(*s.cfg.Interval), *s.cfg.Blocked, *s.cfg.Direct)
	if *s.cfg.ParentType == "ssh" {
		err := s.ConnectSSH()
		if err != nil {
			log.Fatalf("init service fail, ERR: %s", err)
		}
	}
}
func (s *Socks) StopService() {
	if s.sshClient != nil {
		s.sshClient.Close()
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
		err = sc.ListenTCP(s.callback)
	} else {
		err = sc.ListenTls(s.cfg.CertBytes, s.cfg.KeyBytes, s.callback)
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
func (s *Socks) callback(inConn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			//log.Printf("socks conn handler crashed with err : %s \nstack: %s", err, string(debug.Stack()))
		}
		utils.CloseConn(&inConn)
	}()
	var outConn net.Conn
	defer utils.CloseConn(&outConn)

	var b [1024]byte
	n, err := inConn.Read(b[:])
	if err != nil {
		if err != io.EOF {
			log.Printf("read request data fail,ERR: %s", err)
		}
		return
	}

	var reqBytes = b[:n]
	//log.Printf("% x", b[:n])

	//reply
	n, err = inConn.Write([]byte{0x05, 0x00})
	if err != nil {
		log.Printf("reply answer data fail,ERR: %s", err)
		return
	}

	//read answer
	n, err = inConn.Read(b[:])
	if err != nil {
		log.Printf("read answer data fail,ERR: %s", err)
		return
	}
	var headBytes = b[:n]
	// log.Printf("% x", b[:n])
	var addr string
	switch b[3] {
	case 0x01:
		sip := sockIP{}
		if err := binary.Read(bytes.NewReader(b[4:n]), binary.BigEndian, &sip); err != nil {
			log.Printf("read ip fail,ERR: %s", err)
			return
		}
		addr = sip.toAddr()
	case 0x03:
		host := string(b[5 : n-2])
		var port uint16
		err = binary.Read(bytes.NewReader(b[n-2:n]), binary.BigEndian, &port)
		if err != nil {
			log.Printf("read domain fail,ERR: %s", err)
			return
		}
		addr = fmt.Sprintf("%s:%d", host, port)
	}
	useProxy := true
	if *s.cfg.Always {
		outConn, err = s.getOutConn(reqBytes, headBytes, addr)
	} else {
		if *s.cfg.Parent != "" {
			s.checker.Add(addr, true, "", "", nil)
			useProxy, _, _ = s.checker.IsBlocked(addr)
			if useProxy {
				outConn, err = s.getOutConn(reqBytes, headBytes, addr)
			} else {
				outConn, err = utils.ConnectHost(addr, *s.cfg.Timeout)
			}
		} else {
			outConn, err = utils.ConnectHost(addr, *s.cfg.Timeout)
		}
	}
	if err != nil {
		log.Printf("get out conn fail,%s", err)
		inConn.Write([]byte{0x05, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		return
	}
	log.Printf("use proxy %v : %s", useProxy, addr)

	inConn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	inAddr := inConn.RemoteAddr().String()
	inLocalAddr := inConn.LocalAddr().String()

	log.Printf("conn %s - %s connected [%s]", inAddr, inLocalAddr, addr)
	// utils.IoBind(outConn, inConn, func(err error) {
	// 	log.Printf("conn %s - %s released [%s]", inAddr, inLocalAddr, addr)

	// }, func(i int, b bool) {}, 0)
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
			_, err = io.Copy(outConn, inConn)
		}()
		_, err = io.Copy(inConn, outConn)
		return
	}
	bind()
	log.Printf("conn %s - %s released [%s]", inAddr, inLocalAddr, addr)
	utils.CloseConn(&inConn)
	utils.CloseConn(&outConn)
}
func (s *Socks) getOutConn(reqBytes, headBytes []byte, host string) (outConn net.Conn, err error) {
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
		_, err = outConn.Write(reqBytes)
		if err != nil {
			return
		}
		_, err = outConn.Read(buf)
		if err != nil {
			return
		}
		//resp := buf[:n]
		//log.Printf("resp:%v", resp)

		outConn.Write(headBytes)
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
			s.sshClient.Close()
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
	config := ssh.ClientConfig{
		User: *s.cfg.SSHUser,
		Auth: []ssh.AuthMethod{s.cfg.SSHAuthMethod},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	s.sshClient, err = ssh.Dial("tcp", *s.cfg.Parent, &config)
	return
}

type sockIP struct {
	A, B, C, D byte
	PORT       uint16
}

func (ip sockIP) toAddr() string {
	return fmt.Sprintf("%d.%d.%d.%d:%d", ip.A, ip.B, ip.C, ip.D, ip.PORT)
}

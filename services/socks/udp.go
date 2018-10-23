package socks

import (
	"crypto/md5"
	"fmt"
	"net"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/snail007/goproxy/utils"
	goaes "github.com/snail007/goproxy/utils/aes"
	"github.com/snail007/goproxy/utils/socks"
)

func (s *Socks) ParentUDPKey() (key []byte) {
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
func (s *Socks) LocalUDPKey() (key []byte) {
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
func (s *Socks) proxyUDP(inConn *net.Conn, serverConn *socks.ServerConn) {
	defer func() {
		if e := recover(); e != nil {
			s.log.Printf("udp local->out io copy crashed:\n%s\n%s", e, string(debug.Stack()))
		}
	}()
	if *s.cfg.ParentType == "ssh" {
		utils.CloseConn(inConn)
		return
	}
	inconnRemoteAddr := (*inConn).RemoteAddr().String()
	localAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	udpListener := serverConn.UDPConnListener
	srcIP, _, _ := net.SplitHostPort((*inConn).RemoteAddr().String())
	s.log.Printf("proxy udp on %s , for %s", udpListener.LocalAddr(), inconnRemoteAddr)

	s.userConns.Set(inconnRemoteAddr, inConn)
	var (
		outUDPConn       *net.UDPConn
		outconn          net.Conn
		outconnLocalAddr string
		isClosedErr      = func(err error) bool {
			return err != nil && strings.Contains(err.Error(), "use of closed network connection")
		}
		destAddr *net.UDPAddr
	)
	var clean = func(msg, err string) {
		raddr := ""
		if outUDPConn != nil {
			raddr = outUDPConn.RemoteAddr().String()
			outUDPConn.Close()
		}
		if msg != "" {
			if raddr != "" {
				s.log.Printf("%s , %s , %s -> %s", msg, err, inconnRemoteAddr, raddr)
			} else {
				s.log.Printf("%s , %s , from : %s", msg, err, inconnRemoteAddr)
			}
		}
		(*inConn).Close()
		udpListener.Close()
		s.userConns.Remove(inconnRemoteAddr)
		if outconn != nil {
			outconn.Close()
		}
		if outconnLocalAddr != "" {
			s.userConns.Remove(outconnLocalAddr)
		}
	}
	defer clean("", "")
	go func() {
		defer func() {
			if e := recover(); e != nil {
				s.log.Printf("udp related client tcp conn read crashed:\n%s\n%s", e, string(debug.Stack()))
			}
		}()
		buf := make([]byte, 1)
		(*inConn).SetReadDeadline(time.Time{})
		if _, err := (*inConn).Read(buf); err != nil {
			clean("udp related tcp conn disconnected with read", err.Error())
		}
	}()
	go func() {
		defer func() {
			if e := recover(); e != nil {
				s.log.Printf("udp related client tcp conn write crashed:\n%s\n%s", e, string(debug.Stack()))
			}
		}()
		for {
			(*inConn).SetWriteDeadline(time.Now().Add(time.Second * 5))
			if _, err := (*inConn).Write([]byte{0x00}); err != nil {
				clean("udp related tcp conn disconnected with write", err.Error())
				return
			}
			(*inConn).SetWriteDeadline(time.Time{})
			time.Sleep(time.Second * 5)
		}
	}()
	useProxy := true
	if len(*s.cfg.Parent) > 0 {
		dstHost, _, _ := net.SplitHostPort(serverConn.Target())
		if utils.IsInternalIP(dstHost, *s.cfg.Always) {
			useProxy = false
		} else {
			var isInMap bool
			useProxy, isInMap, _, _ = s.checker.IsBlocked(serverConn.Target())
			if !isInMap {
				s.checker.Add(serverConn.Target(), s.Resolve(serverConn.Target()))
			}
		}
	} else {
		useProxy = false
	}
	if useProxy {
		//parent proxy
		lbAddr := s.lb.Select((*inConn).RemoteAddr().String(), *s.cfg.LoadBalanceOnlyHA)
		outconn, err := s.GetParentConn(lbAddr, serverConn)
		//outconn, err := s.GetParentConn(nil, nil, "", false)
		if err != nil {
			clean("connnect fail", fmt.Sprintf("%s", err))
			return
		}

		client, err := s.HandshakeSocksParent(&outconn, "udp", serverConn.Target(), serverConn.AuthData(), false)

		if err != nil {
			clean("handshake fail", fmt.Sprintf("%s", err))
			return
		}
		//outconnRemoteAddr := outconn.RemoteAddr().String()
		outconnLocalAddr = outconn.LocalAddr().String()
		s.userConns.Set(outconnLocalAddr, &outconn)
		go func() {
			defer func() {
				if e := recover(); e != nil {
					s.log.Printf("udp related parent tcp conn read crashed:\n%s\n%s", e, string(debug.Stack()))
				}
			}()
			buf := make([]byte, 1)
			outconn.SetReadDeadline(time.Time{})
			if _, err := outconn.Read(buf); err != nil {
				clean("udp parent tcp conn disconnected", fmt.Sprintf("%s", err))
			}
		}()
		//forward to parent udp
		//s.log.Printf("parent udp address %s", client.UDPAddr)
		destAddr, _ = net.ResolveUDPAddr("udp", client.UDPAddr)
	}
	s.log.Printf("use proxy %v : udp %s", useProxy, serverConn.Target())
	//relay
	buf := utils.LeakyBuffer.Get()
	defer utils.LeakyBuffer.Put(buf)
	for {
		n, srcAddr, err := udpListener.ReadFromUDP(buf)
		if err != nil {
			s.log.Printf("udp listener read fail, %s", err.Error())
			if isClosedErr(err) {
				return
			}
			continue
		}
		srcIP0, _, _ := net.SplitHostPort(srcAddr.String())
		//IP not match drop it
		if srcIP != srcIP0 {
			continue
		}
		p := socks.NewPacketUDP()
		//convert data to raw
		if len(s.udpLocalKey) > 0 {
			var v []byte
			v, err = goaes.Decrypt(s.udpLocalKey, buf[:n])
			if err == nil {
				err = p.Parse(v)
			}
		} else {
			err = p.Parse(buf[:n])
		}
		//err = p.Parse(buf[:n])
		if err != nil {
			s.log.Printf("udp listener parse packet fail, %s", err.Error())
			continue
		}

		port, _ := strconv.Atoi(p.Port())

		if v, ok := s.udpRelatedPacketConns.Get(srcAddr.String()); !ok {
			if destAddr == nil {
				destAddr = &net.UDPAddr{IP: net.ParseIP(p.Host()), Port: port}
			}
			outUDPConn, err = net.DialUDP("udp", localAddr, destAddr)
			if err != nil {
				s.log.Printf("create out udp conn fail , %s , from : %s", err, srcAddr)
				continue
			}
			s.udpRelatedPacketConns.Set(srcAddr.String(), outUDPConn)
			go func() {
				defer func() {
					if e := recover(); e != nil {
						s.log.Printf("udp out->local io copy crashed:\n%s\n%s", e, string(debug.Stack()))
					}
				}()
				defer s.udpRelatedPacketConns.Remove(srcAddr.String())
				//out->local io copy
				buf := utils.LeakyBuffer.Get()
				defer utils.LeakyBuffer.Put(buf)
				for {
					n, err := outUDPConn.Read(buf)
					if err != nil {
						s.log.Printf("read out udp data fail , %s , from : %s", err, srcAddr)
						if isClosedErr(err) {
							return
						}
						continue
					}

					//var dlen = n
					if useProxy {
						//forward to local
						var v []byte
						//convert parent data to raw
						if len(s.udpParentKey) > 0 {
							v, err = goaes.Decrypt(s.udpParentKey, buf[:n])
							if err != nil {
								s.log.Printf("udp outconn parse packet fail, %s", err.Error())
								continue
							}
						} else {
							v = buf[:n]
						}
						//now v is raw, try convert v to local
						if len(s.udpLocalKey) > 0 {
							v, _ = goaes.Encrypt(s.udpLocalKey, v)
						}
						_, err = udpListener.WriteTo(v, srcAddr)
						// _, err = udpListener.WriteTo(buf[:n], srcAddr)
					} else {
						rp := socks.NewPacketUDP()
						rp.Build(destAddr.String(), buf[:n])
						v := rp.Bytes()
						//dlen = len(v)
						//rp.Bytes() v is raw, try convert to local
						if len(s.udpLocalKey) > 0 {
							v, _ = goaes.Encrypt(s.udpLocalKey, v)
						}
						_, err = udpListener.WriteTo(v, srcAddr)
					}

					if err != nil {
						s.udpRelatedPacketConns.Remove(srcAddr.String())
						s.log.Printf("write out data to local fail , %s , from : %s", err, srcAddr)
						if isClosedErr(err) {
							return
						}
						continue
					} else {
						//s.log.Printf("send udp data to local success , len %d, for : %s", dlen, srcAddr)
					}
				}
			}()
		} else {
			outUDPConn = v.(*net.UDPConn)
		}
		//local->out io copy
		if useProxy {
			//forward to parent
			//p is raw, now convert it to parent
			var v []byte
			if len(s.udpParentKey) > 0 {
				v, _ = goaes.Encrypt(s.udpParentKey, p.Bytes())
			} else {
				v = p.Bytes()
			}
			_, err = outUDPConn.Write(v)
			// _, err = outUDPConn.Write(p.Bytes())
		} else {
			_, err = outUDPConn.Write(p.Data())
		}
		if err != nil {
			if isClosedErr(err) {
				return
			}
			s.log.Printf("send out udp data fail , %s , from : %s", err, srcAddr)
			continue
		} else {
			//s.log.Printf("send udp data to remote success , len %d, for : %s", len(p.Data()), srcAddr)
		}
	}

}

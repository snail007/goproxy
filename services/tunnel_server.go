package services

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"io"
	"log"
	"net"
	"proxy/utils"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

type TunnelServer struct {
	cfg    TunnelServerArgs
	udpChn chan UDPItem
	sc     utils.ServerChannel
}

func NewTunnelServer() Service {
	return &TunnelServer{
		cfg:    TunnelServerArgs{},
		udpChn: make(chan UDPItem, 50000),
	}
}

type UDPItem struct {
	packet    *[]byte
	localAddr *net.UDPAddr
	srcAddr   *net.UDPAddr
}

func (s *TunnelServer) InitService() {
	s.UDPConnDeamon()
}
func (s *TunnelServer) Check() {
	if *s.cfg.Parent != "" {
		log.Printf("use tls parent %s", *s.cfg.Parent)
	} else {
		log.Fatalf("parent required")
	}
	if s.cfg.CertBytes == nil || s.cfg.KeyBytes == nil {
		log.Fatalf("cert and key file required")
	}
}
func (s *TunnelServer) StopService() {
}
func (s *TunnelServer) Start(args interface{}) (err error) {
	s.cfg = args.(TunnelServerArgs)
	s.Check()
	s.InitService()
	host, port, _ := net.SplitHostPort(*s.cfg.Local)
	p, _ := strconv.Atoi(port)
	s.sc = utils.NewServerChannel(host, p)

	if *s.cfg.IsUDP {
		err = s.sc.ListenUDP(func(packet []byte, localAddr, srcAddr *net.UDPAddr) {
			s.udpChn <- UDPItem{
				packet:    &packet,
				localAddr: localAddr,
				srcAddr:   srcAddr,
			}
		})
		if err != nil {
			return
		}
		log.Printf("proxy on udp tunnel server mode %s", (*s.sc.UDPListener).LocalAddr())
	} else {
		err = s.sc.ListenTCP(func(inConn net.Conn) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("tserver conn handler crashed with err : %s \nstack: %s", err, string(debug.Stack()))
				}
			}()
			var outConn net.Conn
			for {
				outConn, err = s.GetOutConn()
				if err != nil {
					utils.CloseConn(&outConn)
					log.Printf("connect to %s fail, err: %s, retrying...", *s.cfg.Parent, err)
					time.Sleep(time.Second * 3)
					continue
				} else {
					break
				}
			}

			utils.IoBind(inConn, outConn, func(isSrcErr bool, err error) {
				utils.CloseConn(&outConn)
				utils.CloseConn(&inConn)
				log.Printf("%s conn %s - %s - %s - %s released", *s.cfg.Key, inConn.RemoteAddr(), inConn.LocalAddr(), outConn.LocalAddr(), outConn.RemoteAddr())
			}, func(i int, b bool) {}, 0)

			log.Printf("%s conn %s - %s - %s - %s created", *s.cfg.Key, inConn.RemoteAddr(), inConn.LocalAddr(), outConn.LocalAddr(), outConn.RemoteAddr())
		})
		if err != nil {
			return
		}
		log.Printf("proxy on tunnel server mode %s", (*s.sc.Listener).Addr())
	}
	return
}
func (s *TunnelServer) Clean() {
	s.StopService()
}
func (s *TunnelServer) GetOutConn() (outConn net.Conn, err error) {
	outConn, err = s.GetConn()
	if err != nil {
		log.Printf("connection err: %s", err)
		return
	}
	keyBytes := []byte(*s.cfg.Key)
	keyLength := uint16(len(keyBytes))
	pkg := new(bytes.Buffer)
	binary.Write(pkg, binary.LittleEndian, CONN_SERVER)
	binary.Write(pkg, binary.LittleEndian, keyLength)
	binary.Write(pkg, binary.LittleEndian, keyBytes)
	_, err = outConn.Write(pkg.Bytes())
	if err != nil {
		log.Printf("write connection data err: %s ,retrying...", err)
		utils.CloseConn(&outConn)
		return
	}
	return
}
func (s *TunnelServer) GetConn() (conn net.Conn, err error) {
	var _conn tls.Conn
	_conn, err = utils.TlsConnectHost(*s.cfg.Parent, *s.cfg.Timeout, s.cfg.CertBytes, s.cfg.KeyBytes)
	if err == nil {
		conn = net.Conn(&_conn)
	}
	return
}
func (s *TunnelServer) UDPConnDeamon() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("udp conn deamon crashed with err : %s \nstack: %s", err, string(debug.Stack()))
			}
		}()
		var outConn net.Conn
		var cmdChn = make(chan bool, 1)

		var err error
		for {
			item := <-s.udpChn
		RETRY:
			if outConn == nil {
				for {
					outConn, err = s.GetOutConn()
					if err != nil {
						cmdChn <- true
						outConn = nil
						utils.CloseConn(&outConn)
						log.Printf("connect to %s fail, err: %s, retrying...", *s.cfg.Parent, err)
						time.Sleep(time.Second * 3)
						continue
					} else {
						go func(outConn net.Conn) {
							go func() {
								<-cmdChn
								outConn.Close()
							}()
							for {
								srcAddrFromConn, body, err := utils.ReadUDPPacket(&outConn)
								if err == io.EOF || err == io.ErrUnexpectedEOF {
									log.Printf("udp connection deamon exited, %s -> %s", outConn.LocalAddr(), outConn.RemoteAddr())
									break
								}
								if err != nil {
									log.Printf("parse revecived udp packet fail, err: %s", err)
									continue
								}
								//log.Printf("udp packet revecived over parent , local:%s", srcAddrFromConn)
								_srcAddr := strings.Split(srcAddrFromConn, ":")
								if len(_srcAddr) != 2 {
									log.Printf("parse revecived udp packet fail, addr error : %s", srcAddrFromConn)
									continue
								}
								port, _ := strconv.Atoi(_srcAddr[1])
								dstAddr := &net.UDPAddr{IP: net.ParseIP(_srcAddr[0]), Port: port}
								_, err = s.sc.UDPListener.WriteToUDP(body, dstAddr)
								if err != nil {
									log.Printf("udp response to local %s fail,ERR:%s", srcAddrFromConn, err)
									continue
								}
								//log.Printf("udp response to local %s success", srcAddrFromConn)
							}
						}(outConn)
						break
					}
				}
			}
			writer := bufio.NewWriter(outConn)
			writer.Write(utils.UDPPacket(item.srcAddr.String(), *item.packet))
			err := writer.Flush()
			if err != nil {
				outConn = nil
				log.Printf("write udp packet to %s fail ,flush err:%s", *s.cfg.Parent, err)
				goto RETRY
			}
			//log.Printf("write packet %v", *item.packet)
		}
	}()
}

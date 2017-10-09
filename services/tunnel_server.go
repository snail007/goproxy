package services

import (
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

type TunnelServerManager struct {
	cfg    TunnelServerArgs
	udpChn chan UDPItem
	sc     utils.ServerChannel
}

func NewTunnelServerManager() Service {
	return &TunnelServerManager{
		cfg:    TunnelServerArgs{},
		udpChn: make(chan UDPItem, 50000),
	}
}
func (s *TunnelServerManager) Start(args interface{}) (err error) {
	s.cfg = args.(TunnelServerArgs)
	if *s.cfg.Parent != "" {
		log.Printf("use tls parent %s", *s.cfg.Parent)
	} else {
		log.Fatalf("parent required")
	}
	//log.Printf("route:%v", *s.cfg.Route)
	for _, info := range *s.cfg.Route {
		_routeInfo := strings.Split(info, "@")
		server := NewTunnelServer()
		local := _routeInfo[0]
		remote := _routeInfo[1]
		err = server.Start(TunnelServerArgs{
			Args:    s.cfg.Args,
			Local:   &local,
			IsUDP:   s.cfg.IsUDP,
			Remote:  &remote,
			Key:     s.cfg.Key,
			Timeout: s.cfg.Timeout,
		})
		if err != nil {
			return
		}
	}
	return
}
func (s *TunnelServerManager) Clean() {

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
	if *s.cfg.Remote == "" {
		log.Fatalf("remote required")
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
			var ID string
			for {
				outConn, ID, err = s.GetOutConn("")
				if err != nil {
					utils.CloseConn(&outConn)
					log.Printf("connect to %s fail, err: %s, retrying...", *s.cfg.Parent, err)
					time.Sleep(time.Second * 3)
					continue
				} else {
					break
				}
			}
			// hb := utils.NewHeartbeatReadWriter(&outConn, 3, func(err error, hb *utils.HeartbeatReadWriter) {
			// 	log.Printf("%s conn %s to bridge released", *s.cfg.Key, ID)
			// 	hb.Close()
			// })
			// utils.IoBind(inConn, &hb, func(err error) {
			utils.IoBind(inConn, outConn, func(err error) {
				utils.CloseConn(&outConn)
				utils.CloseConn(&inConn)
				log.Printf("%s conn %s released", *s.cfg.Key, ID)
			}, func(i int, b bool) {}, 0)

			log.Printf("%s conn %s created", *s.cfg.Key, ID)
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
func (s *TunnelServer) GetOutConn(id string) (outConn net.Conn, ID string, err error) {
	outConn, err = s.GetConn()
	if err != nil {
		log.Printf("connection err: %s", err)
		return
	}
	keyBytes := []byte(*s.cfg.Key)
	keyLength := uint16(len(keyBytes))
	ID = utils.Uniqueid()
	IDBytes := []byte(ID)
	if id != "" {
		ID = id
		IDBytes = []byte(id)
	}
	IDLength := uint16(len(IDBytes))
	remoteAddr := []byte("tcp:" + *s.cfg.Remote)
	if *s.cfg.IsUDP {
		remoteAddr = []byte("udp:" + *s.cfg.Remote)
	}
	remoteAddrLength := uint16(len(remoteAddr))
	pkg := new(bytes.Buffer)
	binary.Write(pkg, binary.LittleEndian, CONN_SERVER)
	binary.Write(pkg, binary.LittleEndian, keyLength)
	binary.Write(pkg, binary.LittleEndian, keyBytes)
	binary.Write(pkg, binary.LittleEndian, IDLength)
	binary.Write(pkg, binary.LittleEndian, IDBytes)
	binary.Write(pkg, binary.LittleEndian, remoteAddrLength)
	binary.Write(pkg, binary.LittleEndian, remoteAddr)
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
		// var hb utils.HeartbeatReadWriter
		var ID string
		var cmdChn = make(chan bool, 1)

		var err error
		for {
			item := <-s.udpChn
		RETRY:
			if outConn == nil {
				for {
					outConn, ID, err = s.GetOutConn("")
					if err != nil {
						cmdChn <- true
						outConn = nil
						utils.CloseConn(&outConn)
						log.Printf("connect to %s fail, err: %s, retrying...", *s.cfg.Parent, err)
						time.Sleep(time.Second * 3)
						continue
					} else {
						// hb = utils.NewHeartbeatReadWriter(&outConn, 3, func(err error, hb *utils.HeartbeatReadWriter) {
						// 	log.Printf("%s conn %s to bridge released", *s.cfg.Key, ID)
						// 	hb.Close()
						// })
						// go func(outConn net.Conn, hb utils.HeartbeatReadWriter, ID string) {
						go func(outConn net.Conn, ID string) {
							go func() {
								<-cmdChn
								outConn.Close()
							}()
							for {
								//srcAddrFromConn, body, err := utils.ReadUDPPacket(&hb)
								srcAddrFromConn, body, err := utils.ReadUDPPacket(outConn)
								if err == io.EOF || err == io.ErrUnexpectedEOF {
									log.Printf("UDP deamon connection %s exited", ID)
									break
								}
								if err != nil {
									log.Printf("parse revecived udp packet fail, err: %s ,%v", err, body)
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
								//log.Printf("udp response to local %s success , %v", srcAddrFromConn, body)
							}
							// }(outConn, hb, ID)
						}(outConn, ID)
						break
					}
				}
			}
			outConn.SetWriteDeadline(time.Now().Add(time.Second))
			// _, err = hb.Write(utils.UDPPacket(item.srcAddr.String(), *item.packet))
			_, err = outConn.Write(utils.UDPPacket(item.srcAddr.String(), *item.packet))
			// writer := bufio.NewWriter(outConn)
			// writer.Write(utils.UDPPacket(item.srcAddr.String(), *item.packet))
			// err := writer.Flush()
			outConn.SetWriteDeadline(time.Time{})
			if err != nil {
				utils.CloseConn(&outConn)
				outConn = nil
				log.Printf("write udp packet to %s fail ,flush err:%s ,retrying...", *s.cfg.Parent, err)
				goto RETRY
			}
			//log.Printf("write packet %v", *item.packet)
		}
	}()
}

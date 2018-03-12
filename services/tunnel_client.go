package services

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"snail007/proxy/utils"
	"time"
)

type TunnelClient struct {
	cfg TunnelClientArgs
	// cm       utils.ConnManager
	ctrlConn net.Conn
}

func NewTunnelClient() Service {
	return &TunnelClient{
		cfg: TunnelClientArgs{},
		// cm:  utils.NewConnManager(),
	}
}

func (s *TunnelClient) InitService() {
	// s.InitHeartbeatDeamon()
}

// func (s *TunnelClient) InitHeartbeatDeamon() {
// 	log.Printf("heartbeat started")
// 	go func() {
// 		var heartbeatConn net.Conn
// 		var ID = *s.cfg.Key
// 		for {

// 			//close all connection
// 			s.cm.RemoveAll()
// 			if s.ctrlConn != nil {
// 				s.ctrlConn.Close()
// 			}
// 			utils.CloseConn(&heartbeatConn)
// 			heartbeatConn, err := s.GetInConn(CONN_CLIENT_HEARBEAT, ID)
// 			if err != nil {
// 				log.Printf("heartbeat connection err: %s, retrying...", err)
// 				time.Sleep(time.Second * 3)
// 				utils.CloseConn(&heartbeatConn)
// 				continue
// 			}
// 			log.Printf("heartbeat connection created,id:%s", ID)
// 			writeDie := make(chan bool)
// 			readDie := make(chan bool)
// 			go func() {
// 				for {
// 					heartbeatConn.SetWriteDeadline(time.Now().Add(time.Second * 3))
// 					_, err = heartbeatConn.Write([]byte{0x00})
// 					heartbeatConn.SetWriteDeadline(time.Time{})
// 					if err != nil {
// 						log.Printf("heartbeat connection write err %s", err)
// 						break
// 					}
// 					time.Sleep(time.Second * 3)
// 				}
// 				close(writeDie)
// 			}()
// 			go func() {
// 				for {
// 					signal := make([]byte, 1)
// 					heartbeatConn.SetReadDeadline(time.Now().Add(time.Second * 6))
// 					_, err := heartbeatConn.Read(signal)
// 					heartbeatConn.SetReadDeadline(time.Time{})
// 					if err != nil {
// 						log.Printf("heartbeat connection read err: %s", err)
// 						break
// 					} else {
// 						//log.Printf("heartbeat from bridge")
// 					}
// 				}
// 				close(readDie)
// 			}()
// 			select {
// 			case <-readDie:
// 			case <-writeDie:
// 			}
// 		}
// 	}()
// }
func (s *TunnelClient) CheckArgs() {
	if *s.cfg.Parent != "" {
		log.Printf("use tls parent %s", *s.cfg.Parent)
	} else {
		log.Fatalf("parent required")
	}
	if *s.cfg.CertFile == "" || *s.cfg.KeyFile == "" {
		log.Fatalf("cert and key file required")
	}
	s.cfg.CertBytes, s.cfg.KeyBytes = utils.TlsBytes(*s.cfg.CertFile, *s.cfg.KeyFile)
}
func (s *TunnelClient) StopService() {
	// s.cm.RemoveAll()
}
func (s *TunnelClient) Start(args interface{}) (err error) {
	s.cfg = args.(TunnelClientArgs)
	s.CheckArgs()
	s.InitService()
	log.Printf("proxy on tunnel client mode")

	for {
		//close all conn
		// s.cm.Remove(*s.cfg.Key)
		if s.ctrlConn != nil {
			s.ctrlConn.Close()
		}

		s.ctrlConn, err = s.GetInConn(CONN_CLIENT_CONTROL, *s.cfg.Key)
		if err != nil {
			log.Printf("control connection err: %s, retrying...", err)
			time.Sleep(time.Second * 3)
			if s.ctrlConn != nil {
				s.ctrlConn.Close()
			}
			continue
		}
		for {
			var ID, clientLocalAddr, serverID string
			err = utils.ReadPacketData(s.ctrlConn, &ID, &clientLocalAddr, &serverID)
			if err != nil {
				if s.ctrlConn != nil {
					s.ctrlConn.Close()
				}
				log.Printf("read connection signal err: %s, retrying...", err)
				break
			}
			log.Printf("signal revecived:%s %s %s", serverID, ID, clientLocalAddr)
			protocol := clientLocalAddr[:3]
			localAddr := clientLocalAddr[4:]
			if protocol == "udp" {
				go s.ServeUDP(localAddr, ID, serverID)
			} else {
				go s.ServeConn(localAddr, ID, serverID)
			}
		}
	}
}
func (s *TunnelClient) Clean() {
	s.StopService()
}
func (s *TunnelClient) GetInConn(typ uint8, data ...string) (outConn net.Conn, err error) {
	outConn, err = s.GetConn()
	if err != nil {
		err = fmt.Errorf("connection err: %s", err)
		return
	}
	_, err = outConn.Write(utils.BuildPacket(typ, data...))
	if err != nil {
		err = fmt.Errorf("write connection data err: %s ,retrying...", err)
		utils.CloseConn(&outConn)
		return
	}
	return
}
func (s *TunnelClient) GetConn() (conn net.Conn, err error) {
	var _conn tls.Conn
	_conn, err = utils.TlsConnectHost(*s.cfg.Parent, *s.cfg.Timeout, s.cfg.CertBytes, s.cfg.KeyBytes, nil)
	if err == nil {
		conn = net.Conn(&_conn)
	}
	return
}
func (s *TunnelClient) ServeUDP(localAddr, ID, serverID string) {
	var inConn net.Conn
	var err error
	// for {
	for {
		// s.cm.RemoveOne(*s.cfg.Key, ID)
		inConn, err = s.GetInConn(CONN_CLIENT, *s.cfg.Key, ID, serverID)
		if err != nil {
			utils.CloseConn(&inConn)
			log.Printf("connection err: %s, retrying...", err)
			time.Sleep(time.Second * 3)
			continue
		} else {
			break
		}
	}
	// s.cm.Add(*s.cfg.Key, ID, &inConn)
	log.Printf("conn %s created", ID)

	for {
		srcAddr, body, err := utils.ReadUDPPacket(inConn)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			log.Printf("connection %s released", ID)
			utils.CloseConn(&inConn)
			break
		} else if err != nil {
			log.Printf("udp packet revecived fail, err: %s", err)
		} else {
			//log.Printf("udp packet revecived:%s,%v", srcAddr, body)
			go s.processUDPPacket(&inConn, srcAddr, localAddr, body)
		}

	}
	// }
}
func (s *TunnelClient) processUDPPacket(inConn *net.Conn, srcAddr, localAddr string, body []byte) {
	dstAddr, err := net.ResolveUDPAddr("udp", localAddr)
	if err != nil {
		log.Printf("can't resolve address: %s", err)
		utils.CloseConn(inConn)
		return
	}
	clientSrcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	conn, err := net.DialUDP("udp", clientSrcAddr, dstAddr)
	if err != nil {
		log.Printf("connect to udp %s fail,ERR:%s", dstAddr.String(), err)
		return
	}
	conn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
	_, err = conn.Write(body)
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
	//log.Printf("revecived udp packet from %s , %v", dstAddr.String(), respBody)
	bs := utils.UDPPacket(srcAddr, respBody)
	_, err = (*inConn).Write(bs)
	if err != nil {
		log.Printf("send udp response fail ,ERR:%s", err)
		utils.CloseConn(inConn)
		return
	}
	//log.Printf("send udp response success ,from:%s ,%d ,%v", dstAddr.String(), len(bs), bs)
}
func (s *TunnelClient) ServeConn(localAddr, ID, serverID string) {
	var inConn, outConn net.Conn
	var err error
	for {
		inConn, err = s.GetInConn(CONN_CLIENT, *s.cfg.Key, ID, serverID)
		if err != nil {
			utils.CloseConn(&inConn)
			log.Printf("connection err: %s, retrying...", err)
			time.Sleep(time.Second * 3)
			continue
		} else {
			break
		}
	}

	i := 0
	for {
		i++
		outConn, err = utils.ConnectHost(localAddr, *s.cfg.Timeout)
		if err == nil || i == 3 {
			break
		} else {
			if i == 3 {
				log.Printf("connect to %s err: %s, retrying...", localAddr, err)
				time.Sleep(2 * time.Second)
				continue
			}
		}
	}
	if err != nil {
		utils.CloseConn(&inConn)
		utils.CloseConn(&outConn)
		log.Printf("build connection error, err: %s", err)
		return
	}
	utils.IoBind(inConn, outConn, func(err interface{}) {
		log.Printf("conn %s released", ID)
		// s.cm.RemoveOne(*s.cfg.Key, ID)
	})
	// s.cm.Add(*s.cfg.Key, ID, &inConn)
	log.Printf("conn %s created", ID)
}

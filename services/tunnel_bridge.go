package services

import (
	"bufio"
	"log"
	"net"
	"snail007/proxy/utils"
	"strconv"
	"time"
)

type ServerConn struct {
	//ClientLocalAddr string //tcp:2.2.22:333@ID
	Conn *net.Conn
}
type TunnelBridge struct {
	cfg                TunnelBridgeArgs
	serverConns        utils.ConcurrentMap
	clientControlConns utils.ConcurrentMap
	// cmServer           utils.ConnManager
	// cmClient           utils.ConnManager
}

func NewTunnelBridge() Service {
	return &TunnelBridge{
		cfg:                TunnelBridgeArgs{},
		serverConns:        utils.NewConcurrentMap(),
		clientControlConns: utils.NewConcurrentMap(),
		// cmServer:           utils.NewConnManager(),
		// cmClient:           utils.NewConnManager(),
	}
}

func (s *TunnelBridge) InitService() {

}
func (s *TunnelBridge) CheckArgs() {
	if *s.cfg.CertFile == "" || *s.cfg.KeyFile == "" {
		log.Fatalf("cert and key file required")
	}
	s.cfg.CertBytes, s.cfg.KeyBytes = utils.TlsBytes(*s.cfg.CertFile, *s.cfg.KeyFile)
}
func (s *TunnelBridge) StopService() {

}
func (s *TunnelBridge) Start(args interface{}) (err error) {
	s.cfg = args.(TunnelBridgeArgs)
	s.CheckArgs()
	s.InitService()
	host, port, _ := net.SplitHostPort(*s.cfg.Local)
	p, _ := strconv.Atoi(port)
	sc := utils.NewServerChannel(host, p)

	err = sc.ListenTls(s.cfg.CertBytes, s.cfg.KeyBytes, nil, func(inConn net.Conn) {
		//log.Printf("connection from %s ", inConn.RemoteAddr())

		reader := bufio.NewReader(inConn)
		var err error
		var connType uint8
		err = utils.ReadPacket(reader, &connType)
		if err != nil {
			log.Printf("read error,ERR:%s", err)
			return
		}
		switch connType {
		case CONN_SERVER:
			var key, ID, clientLocalAddr, serverID string
			err = utils.ReadPacketData(reader, &key, &ID, &clientLocalAddr, &serverID)
			if err != nil {
				log.Printf("read error,ERR:%s", err)
				return
			}
			packet := utils.BuildPacketData(ID, clientLocalAddr, serverID)
			log.Printf("server connection, key: %s , id: %s %s %s", key, ID, clientLocalAddr, serverID)

			//addr := clientLocalAddr + "@" + ID
			s.serverConns.Set(ID, ServerConn{
				Conn: &inConn,
			})
			for {
				item, ok := s.clientControlConns.Get(key)
				if !ok {
					log.Printf("client %s control conn not exists", key)
					time.Sleep(time.Second * 3)
					continue
				}
				(*item.(*net.Conn)).SetWriteDeadline(time.Now().Add(time.Second * 3))
				_, err := (*item.(*net.Conn)).Write(packet)
				(*item.(*net.Conn)).SetWriteDeadline(time.Time{})
				if err != nil {
					log.Printf("%s client control conn write signal fail, err: %s, retrying...", key, err)
					time.Sleep(time.Second * 3)
					continue
				} else {
					// s.cmServer.Add(serverID, ID, &inConn)
					break
				}
			}
		case CONN_CLIENT:
			var key, ID, serverID string
			err = utils.ReadPacketData(reader, &key, &ID, &serverID)
			if err != nil {
				log.Printf("read error,ERR:%s", err)
				return
			}
			log.Printf("client connection , key: %s , id: %s, server id:%s", key, ID, serverID)

			serverConnItem, ok := s.serverConns.Get(ID)
			if !ok {
				inConn.Close()
				log.Printf("server conn %s exists", ID)
				return
			}
			serverConn := serverConnItem.(ServerConn).Conn
			utils.IoBind(*serverConn, inConn, func(err interface{}) {
				s.serverConns.Remove(ID)
				// s.cmClient.RemoveOne(key, ID)
				// s.cmServer.RemoveOne(serverID, ID)
				log.Printf("conn %s released", ID)
			})
			// s.cmClient.Add(key, ID, &inConn)
			log.Printf("conn %s created", ID)

		case CONN_CLIENT_CONTROL:
			var key string
			err = utils.ReadPacketData(reader, &key)
			if err != nil {
				log.Printf("read error,ERR:%s", err)
				return
			}
			log.Printf("client control connection, key: %s", key)
			if s.clientControlConns.Has(key) {
				item, _ := s.clientControlConns.Get(key)
				(*item.(*net.Conn)).Close()
			}
			s.clientControlConns.Set(key, &inConn)
			log.Printf("set client %s control conn", key)

			// case CONN_SERVER_HEARBEAT:
			// 	var serverID string
			// 	err = utils.ReadPacketData(reader, &serverID)
			// 	if err != nil {
			// 		log.Printf("read error,ERR:%s", err)
			// 		return
			// 	}
			// 	log.Printf("server heartbeat connection, id: %s", serverID)
			// 	writeDie := make(chan bool)
			// 	readDie := make(chan bool)
			// 	go func() {
			// 		for {
			// 			inConn.SetWriteDeadline(time.Now().Add(time.Second * 3))
			// 			_, err = inConn.Write([]byte{0x00})
			// 			inConn.SetWriteDeadline(time.Time{})
			// 			if err != nil {
			// 				log.Printf("server heartbeat connection write err %s", err)
			// 				break
			// 			}
			// 			time.Sleep(time.Second * 3)
			// 		}
			// 		close(writeDie)
			// 	}()
			// 	go func() {
			// 		for {
			// 			signal := make([]byte, 1)
			// 			inConn.SetReadDeadline(time.Now().Add(time.Second * 6))
			// 			_, err := inConn.Read(signal)
			// 			inConn.SetReadDeadline(time.Time{})
			// 			if err != nil {
			// 				log.Printf("server heartbeat connection read err: %s", err)
			// 				break
			// 			} else {
			// 				// log.Printf("heartbeat from server ,id:%s", serverID)
			// 			}
			// 		}
			// 		close(readDie)
			// 	}()
			// 	select {
			// 	case <-readDie:
			// 	case <-writeDie:
			// 	}
			// 	utils.CloseConn(&inConn)
			// 	s.cmServer.Remove(serverID)
			// 	log.Printf("server heartbeat conn %s released", serverID)
			// case CONN_CLIENT_HEARBEAT:
			// 	var clientID string
			// 	err = utils.ReadPacketData(reader, &clientID)
			// 	if err != nil {
			// 		log.Printf("read error,ERR:%s", err)
			// 		return
			// 	}
			// 	log.Printf("client heartbeat connection, id: %s", clientID)
			// 	writeDie := make(chan bool)
			// 	readDie := make(chan bool)
			// 	go func() {
			// 		for {
			// 			inConn.SetWriteDeadline(time.Now().Add(time.Second * 3))
			// 			_, err = inConn.Write([]byte{0x00})
			// 			inConn.SetWriteDeadline(time.Time{})
			// 			if err != nil {
			// 				log.Printf("client heartbeat connection write err %s", err)
			// 				break
			// 			}
			// 			time.Sleep(time.Second * 3)
			// 		}
			// 		close(writeDie)
			// 	}()
			// 	go func() {
			// 		for {
			// 			signal := make([]byte, 1)
			// 			inConn.SetReadDeadline(time.Now().Add(time.Second * 6))
			// 			_, err := inConn.Read(signal)
			// 			inConn.SetReadDeadline(time.Time{})
			// 			if err != nil {
			// 				log.Printf("client control connection read err: %s", err)
			// 				break
			// 			} else {
			// 				// log.Printf("heartbeat from client ,id:%s", clientID)
			// 			}
			// 		}
			// 		close(readDie)
			// 	}()
			// 	select {
			// 	case <-readDie:
			// 	case <-writeDie:
			// 	}
			// 	utils.CloseConn(&inConn)
			// 	s.cmClient.Remove(clientID)
			// 	if s.clientControlConns.Has(clientID) {
			// 		item, _ := s.clientControlConns.Get(clientID)
			// 		(*item.(*net.Conn)).Close()
			// 	}
			// 	s.clientControlConns.Remove(clientID)
			// 	log.Printf("client heartbeat conn %s released", clientID)
		}
	})
	if err != nil {
		return
	}
	log.Printf("proxy on tunnel bridge mode %s", (*sc.Listener).Addr())
	return
}
func (s *TunnelBridge) Clean() {
	s.StopService()
}

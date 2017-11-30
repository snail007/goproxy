package services

import (
	"bufio"
	"log"
	"net"
	"proxy/utils"
	"strconv"
	"time"

	"github.com/xtaci/smux"
)

type MuxServerConn struct {
	//ClientLocalAddr string //tcp:2.2.22:333@ID
	Conn *net.Conn
}
type MuxBridge struct {
	cfg                MuxBridgeArgs
	serverConns        utils.ConcurrentMap
	clientControlConns utils.ConcurrentMap
}

func NewMuxBridge() Service {
	return &MuxBridge{
		cfg:                MuxBridgeArgs{},
		serverConns:        utils.NewConcurrentMap(),
		clientControlConns: utils.NewConcurrentMap(),
	}
}

func (s *MuxBridge) InitService() {

}
func (s *MuxBridge) CheckArgs() {
	if *s.cfg.CertFile == "" || *s.cfg.KeyFile == "" {
		log.Fatalf("cert and key file required")
	}
	s.cfg.CertBytes, s.cfg.KeyBytes = utils.TlsBytes(*s.cfg.CertFile, *s.cfg.KeyFile)
}
func (s *MuxBridge) StopService() {

}
func (s *MuxBridge) Start(args interface{}) (err error) {
	s.cfg = args.(MuxBridgeArgs)
	s.CheckArgs()
	s.InitService()
	host, port, _ := net.SplitHostPort(*s.cfg.Local)
	p, _ := strconv.Atoi(port)
	sc := utils.NewServerChannel(host, p)

	err = sc.ListenTls(s.cfg.CertBytes, s.cfg.KeyBytes, func(inConn net.Conn) {
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
			session, err := smux.Server(inConn, nil)
			if err != nil {
				utils.CloseConn(&inConn)
				log.Printf("server underlayer connection error,ERR:%s", err)
				return
			}
			conn, err := session.AcceptStream()
			if err != nil {
				session.Close()
				utils.CloseConn(&inConn)
				return
			}
			log.Printf("server connection %s", conn.RemoteAddr())
			//s.callback(conn)
		}
		s.callback(inConn)
	})
	if err != nil {
		return
	}
	log.Printf("proxy on mux bridge mode %s", (*sc.Listener).Addr())
	return
}
func (s *MuxBridge) Clean() {
	s.StopService()
}
func (s *MuxBridge) callback(inConn net.Conn) {
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
		s.serverConns.Set(ID, MuxServerConn{
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
		serverConn := serverConnItem.(MuxServerConn).Conn
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
	}
}

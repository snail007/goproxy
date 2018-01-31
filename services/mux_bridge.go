package services

import (
	"bufio"
	"io"
	"log"
	"net"
	"proxy/utils"
	"strconv"
	"time"

	"github.com/xtaci/smux"
)

type MuxBridge struct {
	cfg                MuxBridgeArgs
	clientControlConns utils.ConcurrentMap
}

func NewMuxBridge() Service {
	return &MuxBridge{
		cfg:                MuxBridgeArgs{},
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
		var key string
		err = utils.ReadPacket(reader, &connType, &key)
		if err != nil {
			log.Printf("read error,ERR:%s", err)
			return
		}
		switch connType {
		case CONN_SERVER:
			var serverID string
			err = utils.ReadPacketData(reader, &serverID)
			if err != nil {
				log.Printf("read error,ERR:%s", err)
				return
			}
			log.Printf("server connection %s %s connected", serverID, key)
			session, err := smux.Server(inConn, nil)
			if err != nil {
				utils.CloseConn(&inConn)
				log.Printf("server session error,ERR:%s", err)
				return
			}
			for {
				stream, err := session.AcceptStream()
				if err != nil {
					session.Close()
					utils.CloseConn(&inConn)
					return
				}
				go s.callback(stream, serverID, key)
			}
		case CONN_CLIENT:

			log.Printf("client connection %s connected", key)
			session, err := smux.Client(inConn, nil)
			if err != nil {
				utils.CloseConn(&inConn)
				log.Printf("client session error,ERR:%s", err)
				return
			}
			s.clientControlConns.Set(key, session)
			go func() {
				for {
					if session.IsClosed() {
						s.clientControlConns.Remove(key)
						break
					}
					time.Sleep(time.Second * 5)
				}
			}()
			//log.Printf("set client session,key: %s", key)
		}

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
func (s *MuxBridge) callback(inConn net.Conn, serverID, key string) {
	try := 20
	for {
		try--
		if try == 0 {
			break
		}
		session, ok := s.clientControlConns.Get(key)
		if !ok {
			log.Printf("client %s session not exists for server stream %s", key, serverID)
			time.Sleep(time.Second * 3)
			continue
		}
		stream, err := session.(*smux.Session).OpenStream()
		if err != nil {
			log.Printf("%s client session open stream %s fail, err: %s, retrying...", key, serverID, err)
			time.Sleep(time.Second * 3)
			continue
		} else {
			log.Printf("%s server %s stream created", key, serverID)
			die1 := make(chan bool, 1)
			die2 := make(chan bool, 1)
			go func() {
				io.Copy(stream, inConn)
				die1 <- true
			}()
			go func() {
				io.Copy(inConn, stream)
				die2 <- true
			}()
			select {
			case <-die1:
			case <-die2:
			}
			stream.Close()
			inConn.Close()
			log.Printf("%s server %s stream released", key, serverID)
			break
		}
	}

}

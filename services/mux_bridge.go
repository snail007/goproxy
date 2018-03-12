package services

import (
	"bufio"
	"io"
	"log"
	"math/rand"
	"net"
	"snail007/proxy/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/xtaci/smux"
)

type MuxBridge struct {
	cfg                MuxBridgeArgs
	clientControlConns utils.ConcurrentMap
	router             utils.ClientKeyRouter
	l                  *sync.Mutex
}

func NewMuxBridge() Service {
	b := &MuxBridge{
		cfg:                MuxBridgeArgs{},
		clientControlConns: utils.NewConcurrentMap(),
		l:                  &sync.Mutex{},
	}
	b.router = utils.NewClientKeyRouter(&b.clientControlConns, 50000)
	return b
}

func (s *MuxBridge) InitService() {

}
func (s *MuxBridge) CheckArgs() {
	if *s.cfg.CertFile == "" || *s.cfg.KeyFile == "" {
		log.Fatalf("cert and key file required")
	}
	if *s.cfg.LocalType == TYPE_TLS {
		s.cfg.CertBytes, s.cfg.KeyBytes = utils.TlsBytes(*s.cfg.CertFile, *s.cfg.KeyFile)
	}
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
	if *s.cfg.LocalType == TYPE_TCP {
		err = sc.ListenTCP(s.handler)
	} else if *s.cfg.LocalType == TYPE_TLS {
		err = sc.ListenTls(s.cfg.CertBytes, s.cfg.KeyBytes, s.handler)
	} else if *s.cfg.LocalType == TYPE_KCP {
		err = sc.ListenKCP(s.cfg.KCP, s.handler)
	}
	if err != nil {
		return
	}
	log.Printf("%s bridge on %s", *s.cfg.LocalType, (*sc.Listener).Addr())
	return
}
func (s *MuxBridge) Clean() {
	s.StopService()
}
func (s *MuxBridge) handler(inConn net.Conn) {
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
		keyInfo := strings.Split(key, "-")
		if len(keyInfo) != 2 {
			utils.CloseConn(&inConn)
			log.Printf("client key format error,key:%s", key)
			return
		}
		groupKey := keyInfo[0]
		index := keyInfo[1]
		if !s.clientControlConns.Has(groupKey) {
			item := utils.NewConcurrentMap()
			s.clientControlConns.Set(groupKey, &item)
		}
		_group, _ := s.clientControlConns.Get(groupKey)
		group := _group.(*utils.ConcurrentMap)
		s.l.Lock()
		defer s.l.Unlock()
		group.Set(index, session)
		// s.clientControlConns.Set(key, session)
		go func() {
			for {
				if session.IsClosed() {
					s.l.Lock()
					defer s.l.Unlock()
					if sess, ok := group.Get(index); ok && sess.(*smux.Session).IsClosed() {
						group.Remove(index)
					}
					break
				}
				time.Sleep(time.Second * 5)
			}
		}()
		//log.Printf("set client session,key: %s", key)
	}

}
func (s *MuxBridge) callback(inConn net.Conn, serverID, key string) {
	try := 20
	for {
		try--
		if try == 0 {
			break
		}
		if key == "*" {
			key = s.router.GetKey()
		}
		_group, ok := s.clientControlConns.Get(key)
		if !ok {
			log.Printf("client %s session not exists for server stream %s, retrying...", key, serverID)
			time.Sleep(time.Second * 3)
			continue
		}
		group := _group.(*utils.ConcurrentMap)
		keys := group.Keys()
		keysLen := len(keys)
		i := 0
		if keysLen > 0 {
			i = rand.Intn(keysLen)
		} else {
			log.Printf("client %s session empty for server stream %s, retrying...", key, serverID)
			time.Sleep(time.Second * 3)
			continue
		}
		index := keys[i]
		log.Printf("select client : %s-%s", key, index)
		session, _ := group.Get(index)
		stream, err := session.(*smux.Session).OpenStream()
		if err != nil {
			log.Printf("%s client session open stream %s fail, err: %s, retrying...", key, serverID, err)
			time.Sleep(time.Second * 3)
			continue
		} else {
			log.Printf("stream %s -> %s created", serverID, key)
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

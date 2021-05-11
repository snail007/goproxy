package services

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"proxy/utils"
	"strconv"
	"sync"
	"time"
)

type BridgeItem struct {
	ServerChn     chan *net.Conn
	ClientChn     chan *net.Conn
	ClientControl *net.Conn
	Once          *sync.Once
	Key           string
}
type TunnelBridge struct {
	cfg TunnelBridgeArgs
	br  utils.ConcurrentMap
}

func NewTunnelBridge() Service {
	return &TunnelBridge{
		cfg: TunnelBridgeArgs{},
		br:  utils.NewConcurrentMap(),
	}
}

func (s *TunnelBridge) InitService() {

}
func (s *TunnelBridge) Check() {
	if s.cfg.CertBytes == nil || s.cfg.KeyBytes == nil {
		log.Fatalf("cert and key file required")
	}

}
func (s *TunnelBridge) StopService() {

}
func (s *TunnelBridge) Start(args interface{}) (err error) {
	s.cfg = args.(TunnelBridgeArgs)
	s.Check()
	s.InitService()
	host, port, _ := net.SplitHostPort(*s.cfg.Local)
	p, _ := strconv.Atoi(port)
	sc := utils.NewServerChannel(host, p)

	err = sc.ListenTls(s.cfg.CertBytes, s.cfg.KeyBytes, func(inConn net.Conn) {
		reader := bufio.NewReader(inConn)
		var connType uint8
		err = binary.Read(reader, binary.LittleEndian, &connType)
		if err != nil {
			utils.CloseConn(&inConn)
			return
		}
		var key string
		var connTypeStrMap = map[uint8]string{CONN_SERVER: "server", CONN_CLIENT: "client", CONN_CONTROL: "client"}
		if connType == CONN_SERVER || connType == CONN_CLIENT || connType == CONN_CONTROL {
			var keyLength uint16
			err = binary.Read(reader, binary.LittleEndian, &keyLength)
			if err != nil {
				return
			}
			_key := make([]byte, keyLength)
			n, err := reader.Read(_key)
			if err != nil {
				return
			}
			if n != int(keyLength) {
				return
			}
			key = string(_key)
			log.Printf("connection from %s , key: %s", connTypeStrMap[connType], key)
		}
		switch connType {
		case CONN_SERVER:
			s.ServerConn(&inConn, key)
		case CONN_CLIENT:
			s.ClientConn(&inConn, key)
		case CONN_CONTROL:
			s.ClientControlConn(&inConn, key)
		default:
			log.Printf("unkown conn type %d", connType)
			utils.CloseConn(&inConn)
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
func (s *TunnelBridge) ClientConn(inConn *net.Conn, key string) {
	chn, _ := s.ConnChn(key, CONN_CLIENT)
	chn <- inConn
}
func (s *TunnelBridge) ServerConn(inConn *net.Conn, key string) {
	chn, _ := s.ConnChn(key, CONN_SERVER)
	chn <- inConn
}
func (s *TunnelBridge) ClientControlConn(inConn *net.Conn, key string) {
	_, item := s.ConnChn(key, CONN_CLIENT)
	utils.CloseConn(item.ClientControl)
	if item.ClientControl != nil {
		*item.ClientControl = *inConn
	} else {
		item.ClientControl = inConn
	}
	log.Printf("set client control conn,remote: %s", (*inConn).RemoteAddr())
}
func (s *TunnelBridge) ConnChn(key string, typ uint8) (chn chan *net.Conn, item *BridgeItem) {
	s.br.SetIfAbsent(key, &BridgeItem{
		ServerChn: make(chan *net.Conn, 10000),
		ClientChn: make(chan *net.Conn, 10000),
		Once:      &sync.Once{},
		Key:       key,
	})
	_item, _ := s.br.Get(key)
	item = _item.(*BridgeItem)
	item.Once.Do(func() {
		s.ChnDeamon(item)
	})
	if typ == CONN_CLIENT {
		chn = item.ClientChn
	} else {
		chn = item.ServerChn
	}
	return
}
func (s *TunnelBridge) ChnDeamon(item *BridgeItem) {
	go func() {
		log.Printf("%s conn chan deamon started", item.Key)
		for {
			var clientConn *net.Conn
			var serverConn *net.Conn
			serverConn = <-item.ServerChn
			log.Printf("%s server conn picked up", item.Key)
		OUT:
			for {
				_item, _ := s.br.Get(item.Key)
				Item := _item.(*BridgeItem)
				var err error
				if Item.ClientControl != nil && *Item.ClientControl != nil {
					_, err = (*Item.ClientControl).Write([]byte{'0'})
				} else {
					err = fmt.Errorf("client control conn not exists")
				}
				if err != nil {
					log.Printf("%s client control conn write signal fail, err: %s, retrying...", item.Key, err)
					utils.CloseConn(Item.ClientControl)
					*Item.ClientControl = nil
					Item.ClientControl = nil
					time.Sleep(time.Second * 3)
					continue
				} else {
					select {
					case clientConn = <-item.ClientChn:
						log.Printf("%s client conn picked up", item.Key)
						break OUT
					case <-time.After(time.Second * time.Duration(*s.cfg.Timeout*5)):
						log.Printf("%s client conn picked timeout, retrying...", item.Key)
					}
				}
			}

			utils.IoBind(*serverConn, *clientConn, func(isSrcErr bool, err error) {
				utils.CloseConn(serverConn)
				utils.CloseConn(clientConn)
				log.Printf("%s conn %s - %s - %s - %s released", item.Key, (*serverConn).RemoteAddr(), (*serverConn).LocalAddr(), (*clientConn).LocalAddr(), (*clientConn).RemoteAddr())
			}, func(i int, b bool) {}, 0)
			log.Printf("%s conn %s - %s - %s - %s created", item.Key, (*serverConn).RemoteAddr(), (*serverConn).LocalAddr(), (*clientConn).LocalAddr(), (*clientConn).RemoteAddr())
		}
	}()
}

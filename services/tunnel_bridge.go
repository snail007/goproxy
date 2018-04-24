package services

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"snail007/proxy/utils"
	"strconv"
	"time"

	"github.com/xtaci/smux"
)

type ServerConn struct {
	//ClientLocalAddr string //tcp:2.2.22:333@ID
	Conn *net.Conn
}
type TunnelBridge struct {
	cfg                TunnelBridgeArgs
	serverConns        utils.ConcurrentMap
	clientControlConns utils.ConcurrentMap
	isStop             bool
}

func NewTunnelBridge() Service {
	return &TunnelBridge{
		cfg:                TunnelBridgeArgs{},
		serverConns:        utils.NewConcurrentMap(),
		clientControlConns: utils.NewConcurrentMap(),
		isStop:             false,
	}
}

func (s *TunnelBridge) InitService() (err error) {
	return
}
func (s *TunnelBridge) CheckArgs() (err error) {
	if *s.cfg.CertFile == "" || *s.cfg.KeyFile == "" {
		err = fmt.Errorf("cert and key file required")
		return
	}
	s.cfg.CertBytes, s.cfg.KeyBytes, err = utils.TlsBytes(*s.cfg.CertFile, *s.cfg.KeyFile)
	return
}
func (s *TunnelBridge) StopService() {
	defer func() {
		e := recover()
		if e != nil {
			log.Printf("stop tbridge service crashed,%s", e)
		} else {
			log.Printf("service tbridge stoped")
		}
	}()
	s.isStop = true
	for _, sess := range s.clientControlConns.Items() {
		(*sess.(*net.Conn)).Close()
	}
	for _, sess := range s.serverConns.Items() {
		(*sess.(ServerConn).Conn).Close()
	}
}
func (s *TunnelBridge) Start(args interface{}) (err error) {
	s.cfg = args.(TunnelBridgeArgs)
	if err = s.CheckArgs(); err != nil {
		return
	}
	if err = s.InitService(); err != nil {
		return
	}
	host, port, _ := net.SplitHostPort(*s.cfg.Local)
	p, _ := strconv.Atoi(port)
	sc := utils.NewServerChannel(host, p)

	err = sc.ListenTls(s.cfg.CertBytes, s.cfg.KeyBytes, nil, s.callback)
	if err != nil {
		return
	}
	log.Printf("proxy on tunnel bridge mode %s", (*sc.Listener).Addr())
	return
}
func (s *TunnelBridge) Clean() {
	s.StopService()
}
func (s *TunnelBridge) callback(inConn net.Conn) {
	var err error
	//log.Printf("connection from %s ", inConn.RemoteAddr())
	sess, err := smux.Server(inConn, &smux.Config{
		KeepAliveInterval: 10 * time.Second,
		KeepAliveTimeout:  time.Duration(*s.cfg.Timeout) * time.Second,
		MaxFrameSize:      4096,
		MaxReceiveBuffer:  4194304,
	})
	if err != nil {
		log.Printf("new mux server conn error,ERR:%s", err)
		return
	}
	inConn, err = sess.AcceptStream()
	if err != nil {
		log.Printf("mux server conn accept error,ERR:%s", err)
		return
	}

	var buf = make([]byte, 1024)
	n, _ := inConn.Read(buf)
	reader := bytes.NewReader(buf[:n])
	//reader := bufio.NewReader(inConn)

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
			if s.isStop {
				return
			}
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
	}
}

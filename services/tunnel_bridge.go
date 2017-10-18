package services

import (
	"bufio"
	"encoding/binary"
	"log"
	"net"
	"proxy/utils"
	"strconv"
	"time"
)

type ServerConn struct {
	ClientLocalAddr string //tcp:2.2.22:333@ID
	Conn            *net.Conn
}
type TunnelBridge struct {
	cfg                TunnelBridgeArgs
	serverConns        utils.ConcurrentMap
	clientControlConns utils.ConcurrentMap
}

func NewTunnelBridge() Service {
	return &TunnelBridge{
		cfg:                TunnelBridgeArgs{},
		serverConns:        utils.NewConcurrentMap(),
		clientControlConns: utils.NewConcurrentMap(),
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

	err = sc.ListenTls(s.cfg.CertBytes, s.cfg.KeyBytes, func(inConn net.Conn) {
		//log.Printf("connection from %s ", inConn.RemoteAddr())

		reader := bufio.NewReader(inConn)
		var connType uint8
		err = binary.Read(reader, binary.LittleEndian, &connType)
		if err != nil {
			utils.CloseConn(&inConn)
			return
		}
		//log.Printf("conn type %d", connType)

		var key, clientLocalAddr, ID string
		var connTypeStrMap = map[uint8]string{CONN_SERVER: "server", CONN_CLIENT: "client", CONN_CONTROL: "client"}
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

		if connType != CONN_CONTROL {
			var IDLength uint16
			err = binary.Read(reader, binary.LittleEndian, &IDLength)
			if err != nil {
				return
			}
			_id := make([]byte, IDLength)
			n, err := reader.Read(_id)
			if err != nil {
				return
			}
			if n != int(IDLength) {
				return
			}
			ID = string(_id)

			if connType == CONN_SERVER {
				var addrLength uint16
				err = binary.Read(reader, binary.LittleEndian, &addrLength)
				if err != nil {
					return
				}
				_addr := make([]byte, addrLength)
				n, err = reader.Read(_addr)
				if err != nil {
					return
				}
				if n != int(addrLength) {
					return
				}
				clientLocalAddr = string(_addr)
			}
		}
		log.Printf("connection from %s , key: %s , id: %s", connTypeStrMap[connType], key, ID)

		switch connType {
		case CONN_SERVER:
			addr := clientLocalAddr + "@" + ID
			s.serverConns.Set(ID, ServerConn{
				Conn:            &inConn,
				ClientLocalAddr: addr,
			})
			for {
				item, ok := s.clientControlConns.Get(key)
				if !ok {
					log.Printf("client %s control conn not exists", key)
					time.Sleep(time.Second * 3)
					continue
				}
				(*item.(*net.Conn)).SetWriteDeadline(time.Now().Add(time.Second * 3))
				_, err := (*item.(*net.Conn)).Write([]byte(addr))
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
			serverConnItem, ok := s.serverConns.Get(ID)
			if !ok {
				inConn.Close()
				log.Printf("server conn %s exists", ID)
				return
			}
			serverConn := serverConnItem.(ServerConn).Conn
			utils.IoBind(*serverConn, inConn, func(err error) {

				(*serverConn).Close()
				utils.CloseConn(&inConn)
				s.serverConns.Remove(ID)
				log.Printf("conn %s released", ID)
			}, func(i int, b bool) {}, 0)
			log.Printf("conn %s created", ID)

		case CONN_CONTROL:
			if s.clientControlConns.Has(key) {
				item, _ := s.clientControlConns.Get(key)
				(*item.(*net.Conn)).Close()
			}
			s.clientControlConns.Set(key, &inConn)
			log.Printf("set client %s control conn", key)
			go func() {
				for {
					var b = make([]byte, 1)
					_, err = inConn.Read(b)
					if err != nil {
						inConn.Close()
						s.serverConns.Remove(ID)
						log.Printf("%s control conn from client released", key)
						break
					} else {
						//log.Printf("%s heartbeat from client", key)
					}
				}
			}()
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

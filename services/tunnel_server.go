package services

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"runtime/debug"
	"github.com/snail007/goproxy/utils"
	"strconv"
	"strings"
	"time"

	"github.com/xtaci/smux"
)

type TunnelServer struct {
	cfg       TunnelServerArgs
	udpChn    chan UDPItem
	sc        utils.ServerChannel
	isStop    bool
	udpConn   *net.Conn
	userConns utils.ConcurrentMap
}

type TunnelServerManager struct {
	cfg      TunnelServerArgs
	udpChn   chan UDPItem
	serverID string
	servers  []*Service
}

func NewTunnelServerManager() Service {
	return &TunnelServerManager{
		cfg:      TunnelServerArgs{},
		udpChn:   make(chan UDPItem, 50000),
		serverID: utils.Uniqueid(),
		servers:  []*Service{},
	}
}
func (s *TunnelServerManager) Start(args interface{}) (err error) {
	s.cfg = args.(TunnelServerArgs)
	if err = s.CheckArgs(); err != nil {
		return
	}
	if *s.cfg.Parent != "" {
		log.Printf("use tls parent %s", *s.cfg.Parent)
	} else {
		err = fmt.Errorf("parent required")
		return
	}

	if err = s.InitService(); err != nil {
		return
	}

	log.Printf("server id: %s", s.serverID)
	//log.Printf("route:%v", *s.cfg.Route)
	for _, _info := range *s.cfg.Route {
		IsUDP := *s.cfg.IsUDP
		if strings.HasPrefix(_info, "udp://") {
			IsUDP = true
		}
		info := strings.TrimPrefix(_info, "udp://")
		info = strings.TrimPrefix(info, "tcp://")
		_routeInfo := strings.Split(info, "@")
		server := NewTunnelServer()
		local := _routeInfo[0]
		remote := _routeInfo[1]
		KEY := *s.cfg.Key
		if strings.HasPrefix(remote, "[") {
			KEY = remote[1:strings.LastIndex(remote, "]")]
			remote = remote[strings.LastIndex(remote, "]")+1:]
		}
		if strings.HasPrefix(remote, ":") {
			remote = fmt.Sprintf("127.0.0.1%s", remote)
		}
		err = server.Start(TunnelServerArgs{
			CertBytes: s.cfg.CertBytes,
			KeyBytes:  s.cfg.KeyBytes,
			Parent:    s.cfg.Parent,
			CertFile:  s.cfg.CertFile,
			KeyFile:   s.cfg.KeyFile,
			Local:     &local,
			IsUDP:     &IsUDP,
			Remote:    &remote,
			Key:       &KEY,
			Timeout:   s.cfg.Timeout,
			Mgr:       s,
		})

		if err != nil {
			return
		}
		s.servers = append(s.servers, &server)
	}
	return
}
func (s *TunnelServerManager) Clean() {
	s.StopService()
}
func (s *TunnelServerManager) StopService() {
	for _, server := range s.servers {
		(*server).Clean()
	}
}
func (s *TunnelServerManager) CheckArgs() (err error) {
	if *s.cfg.CertFile == "" || *s.cfg.KeyFile == "" {
		err = fmt.Errorf("cert and key file required")
		return
	}
	s.cfg.CertBytes, s.cfg.KeyBytes, err = utils.TlsBytes(*s.cfg.CertFile, *s.cfg.KeyFile)
	return
}
func (s *TunnelServerManager) InitService() (err error) {
	return
}

func (s *TunnelServerManager) GetOutConn(typ uint8) (outConn net.Conn, ID string, err error) {
	outConn, err = s.GetConn()
	if err != nil {
		log.Printf("connection err: %s", err)
		return
	}
	ID = s.serverID
	_, err = outConn.Write(utils.BuildPacket(typ, s.serverID))
	if err != nil {
		log.Printf("write connection data err: %s ,retrying...", err)
		utils.CloseConn(&outConn)
		return
	}
	return
}
func (s *TunnelServerManager) GetConn() (conn net.Conn, err error) {
	var _conn tls.Conn
	_conn, err = utils.TlsConnectHost(*s.cfg.Parent, *s.cfg.Timeout, s.cfg.CertBytes, s.cfg.KeyBytes, nil)
	if err == nil {
		conn = net.Conn(&_conn)
	}
	return
}
func NewTunnelServer() Service {
	return &TunnelServer{
		cfg:       TunnelServerArgs{},
		udpChn:    make(chan UDPItem, 50000),
		isStop:    false,
		userConns: utils.NewConcurrentMap(),
	}
}

type UDPItem struct {
	packet    *[]byte
	localAddr *net.UDPAddr
	srcAddr   *net.UDPAddr
}

func (s *TunnelServer) StopService() {
	defer func() {
		e := recover()
		if e != nil {
			log.Printf("stop server service crashed,%s", e)
		} else {
			log.Printf("service server stoped")
		}
	}()
	s.isStop = true

	if s.sc.Listener != nil {
		(*s.sc.Listener).Close()
	}
	if s.sc.UDPListener != nil {
		(*s.sc.UDPListener).Close()
	}
	if s.udpConn != nil {
		(*s.udpConn).Close()
	}
	for _, c := range s.userConns.Items() {
		(*c.(*net.Conn)).Close()
	}
}
func (s *TunnelServer) InitService() (err error) {
	s.UDPConnDeamon()
	return
}
func (s *TunnelServer) CheckArgs() (err error) {
	if *s.cfg.Remote == "" {
		err = fmt.Errorf("remote required")
		return
	}
	return
}

func (s *TunnelServer) Start(args interface{}) (err error) {
	s.cfg = args.(TunnelServerArgs)
	if err = s.CheckArgs(); err != nil {
		return
	}
	if err = s.InitService(); err != nil {
		return
	}
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
				if s.isStop {
					return
				}
				outConn, ID, err = s.GetOutConn(CONN_SERVER)
				if err != nil {
					utils.CloseConn(&outConn)
					log.Printf("connect to %s fail, err: %s, retrying...", *s.cfg.Parent, err)
					time.Sleep(time.Second * 3)
					continue
				} else {
					break
				}
			}
			inAddr := inConn.RemoteAddr().String()
			utils.IoBind(inConn, outConn, func(err interface{}) {
				s.userConns.Remove(inAddr)
				log.Printf("%s conn %s released", *s.cfg.Key, ID)
			})
			if c, ok := s.userConns.Get(inAddr); ok {
				(*c.(*net.Conn)).Close()
			}
			s.userConns.Set(inAddr, &inConn)
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

}
func (s *TunnelServer) GetOutConn(typ uint8) (outConn net.Conn, ID string, err error) {
	outConn, err = s.GetConn()
	if err != nil {
		log.Printf("connection err: %s", err)
		return
	}
	remoteAddr := "tcp:" + *s.cfg.Remote
	if *s.cfg.IsUDP {
		remoteAddr = "udp:" + *s.cfg.Remote
	}
	ID = utils.Uniqueid()
	_, err = outConn.Write(utils.BuildPacket(typ, *s.cfg.Key, ID, remoteAddr, s.cfg.Mgr.serverID))
	if err != nil {
		log.Printf("write connection data err: %s ,retrying...", err)
		utils.CloseConn(&outConn)
		return
	}
	return
}
func (s *TunnelServer) GetConn() (conn net.Conn, err error) {
	var _conn tls.Conn
	_conn, err = utils.TlsConnectHost(*s.cfg.Parent, *s.cfg.Timeout, s.cfg.CertBytes, s.cfg.KeyBytes, nil)
	if err == nil {
		conn = net.Conn(&_conn)
		c, e := smux.Client(conn, &smux.Config{
			KeepAliveInterval: 10 * time.Second,
			KeepAliveTimeout:  time.Duration(*s.cfg.Timeout) * time.Second,
			MaxFrameSize:      4096,
			MaxReceiveBuffer:  4194304,
		})
		if e != nil {
			log.Printf("new mux client conn error,ERR:%s", e)
			err = e
			return
		}
		conn, e = c.OpenStream()
		if e != nil {
			log.Printf("mux client conn open stream error,ERR:%s", e)
			err = e
			return
		}
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
		// var cmdChn = make(chan bool, 1000)
		var err error
		for {
			if s.isStop {
				return
			}
			item := <-s.udpChn
		RETRY:
			if s.isStop {
				return
			}
			if outConn == nil {
				for {
					if s.isStop {
						return
					}
					outConn, ID, err = s.GetOutConn(CONN_SERVER)
					if err != nil {
						// cmdChn <- true
						outConn = nil
						utils.CloseConn(&outConn)
						log.Printf("connect to %s fail, err: %s, retrying...", *s.cfg.Parent, err)
						time.Sleep(time.Second * 3)
						continue
					} else {
						go func(outConn net.Conn, ID string) {
							if s.udpConn != nil {
								(*s.udpConn).Close()
							}
							s.udpConn = &outConn
							for {
								if s.isStop {
									return
								}
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
						}(outConn, ID)
						break
					}
				}
			}
			outConn.SetWriteDeadline(time.Now().Add(time.Second))
			_, err = outConn.Write(utils.UDPPacket(item.srcAddr.String(), *item.packet))
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

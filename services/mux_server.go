package services

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"runtime/debug"
	"github.com/snail007/goproxy/utils"
	"strconv"
	"strings"
	"time"

	"github.com/golang/snappy"
	"github.com/xtaci/smux"
)

type MuxServer struct {
	cfg      MuxServerArgs
	udpChn   chan MuxUDPItem
	sc       utils.ServerChannel
	sessions utils.ConcurrentMap
	lockChn  chan bool
	isStop   bool
	udpConn  *net.Conn
}

type MuxServerManager struct {
	cfg      MuxServerArgs
	udpChn   chan MuxUDPItem
	serverID string
	servers  []*Service
}

func NewMuxServerManager() Service {
	return &MuxServerManager{
		cfg:      MuxServerArgs{},
		udpChn:   make(chan MuxUDPItem, 50000),
		serverID: utils.Uniqueid(),
		servers:  []*Service{},
	}
}

func (s *MuxServerManager) Start(args interface{}) (err error) {
	s.cfg = args.(MuxServerArgs)
	if err = s.CheckArgs(); err != nil {
		return
	}
	if *s.cfg.Parent != "" {
		log.Printf("use %s parent %s", *s.cfg.ParentType, *s.cfg.Parent)
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
		if _info == "" {
			continue
		}
		IsUDP := *s.cfg.IsUDP
		if strings.HasPrefix(_info, "udp://") {
			IsUDP = true
		}
		info := strings.TrimPrefix(_info, "udp://")
		info = strings.TrimPrefix(info, "tcp://")
		_routeInfo := strings.Split(info, "@")
		server := NewMuxServer()
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
		err = server.Start(MuxServerArgs{
			CertBytes:    s.cfg.CertBytes,
			KeyBytes:     s.cfg.KeyBytes,
			Parent:       s.cfg.Parent,
			CertFile:     s.cfg.CertFile,
			KeyFile:      s.cfg.KeyFile,
			Local:        &local,
			IsUDP:        &IsUDP,
			Remote:       &remote,
			Key:          &KEY,
			Timeout:      s.cfg.Timeout,
			Mgr:          s,
			IsCompress:   s.cfg.IsCompress,
			SessionCount: s.cfg.SessionCount,
			KCP:          s.cfg.KCP,
			ParentType:   s.cfg.ParentType,
		})

		if err != nil {
			return
		}
		s.servers = append(s.servers, &server)
	}
	return
}
func (s *MuxServerManager) Clean() {
	s.StopService()
}
func (s *MuxServerManager) StopService() {
	for _, server := range s.servers {
		(*server).Clean()
	}
}
func (s *MuxServerManager) CheckArgs() (err error) {
	if *s.cfg.CertFile == "" || *s.cfg.KeyFile == "" {
		err = fmt.Errorf("cert and key file required")
		return
	}
	if *s.cfg.ParentType == "tls" {
		s.cfg.CertBytes, s.cfg.KeyBytes, err = utils.TlsBytes(*s.cfg.CertFile, *s.cfg.KeyFile)
		if err != nil {
			return
		}
	}
	return
}
func (s *MuxServerManager) InitService() (err error) {
	return
}

func NewMuxServer() Service {
	return &MuxServer{
		cfg:      MuxServerArgs{},
		udpChn:   make(chan MuxUDPItem, 50000),
		lockChn:  make(chan bool, 1),
		sessions: utils.NewConcurrentMap(),
		isStop:   false,
	}
}

type MuxUDPItem struct {
	packet    *[]byte
	localAddr *net.UDPAddr
	srcAddr   *net.UDPAddr
}

func (s *MuxServer) StopService() {
	defer func() {
		e := recover()
		if e != nil {
			log.Printf("stop server service crashed,%s", e)
		} else {
			log.Printf("service server stoped")
		}
	}()
	s.isStop = true
	for _, sess := range s.sessions.Items() {
		sess.(*smux.Session).Close()
	}
	if s.sc.Listener != nil {
		(*s.sc.Listener).Close()
	}
	if s.sc.UDPListener != nil {
		(*s.sc.UDPListener).Close()
	}
	if s.udpConn != nil {
		(*s.udpConn).Close()
	}
}
func (s *MuxServer) InitService() (err error) {
	s.UDPConnDeamon()
	return
}
func (s *MuxServer) CheckArgs() (err error) {
	if *s.cfg.Remote == "" {
		err = fmt.Errorf("remote required")
		return
	}
	return
}

func (s *MuxServer) Start(args interface{}) (err error) {
	s.cfg = args.(MuxServerArgs)
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
			s.udpChn <- MuxUDPItem{
				packet:    &packet,
				localAddr: localAddr,
				srcAddr:   srcAddr,
			}
		})
		if err != nil {
			return
		}
		log.Printf("server on %s", (*s.sc.UDPListener).LocalAddr())
	} else {
		err = s.sc.ListenTCP(func(inConn net.Conn) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("connection handler crashed with err : %s \nstack: %s", err, string(debug.Stack()))
				}
			}()
			var outConn net.Conn
			var ID string
			for {
				if s.isStop {
					return
				}
				outConn, ID, err = s.GetOutConn()
				if err != nil {
					utils.CloseConn(&outConn)
					log.Printf("connect to %s fail, err: %s, retrying...", *s.cfg.Parent, err)
					time.Sleep(time.Second * 3)
					continue
				} else {
					break
				}
			}
			log.Printf("%s stream %s created", *s.cfg.Key, ID)
			if *s.cfg.IsCompress {
				die1 := make(chan bool, 1)
				die2 := make(chan bool, 1)
				go func() {
					io.Copy(inConn, snappy.NewReader(outConn))
					die1 <- true
				}()
				go func() {
					io.Copy(snappy.NewWriter(outConn), inConn)
					die2 <- true
				}()
				select {
				case <-die1:
				case <-die2:
				}
				outConn.Close()
				inConn.Close()
				log.Printf("%s stream %s released", *s.cfg.Key, ID)
			} else {
				utils.IoBind(inConn, outConn, func(err interface{}) {
					log.Printf("%s stream %s released", *s.cfg.Key, ID)
				})
			}
		})
		if err != nil {
			return
		}
		log.Printf("server on %s", (*s.sc.Listener).Addr())
	}
	return
}
func (s *MuxServer) Clean() {
	s.StopService()
}
func (s *MuxServer) GetOutConn() (outConn net.Conn, ID string, err error) {
	i := 1
	if *s.cfg.SessionCount > 0 {
		i = rand.Intn(*s.cfg.SessionCount)
	}
	outConn, err = s.GetConn(fmt.Sprintf("%d", i))
	if err != nil {
		log.Printf("connection err: %s", err)
		return
	}
	remoteAddr := "tcp:" + *s.cfg.Remote
	if *s.cfg.IsUDP {
		remoteAddr = "udp:" + *s.cfg.Remote
	}
	ID = utils.Uniqueid()
	outConn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
	_, err = outConn.Write(utils.BuildPacketData(ID, remoteAddr, s.cfg.Mgr.serverID))
	outConn.SetDeadline(time.Time{})
	if err != nil {
		log.Printf("write stream data err: %s ,retrying...", err)
		utils.CloseConn(&outConn)
		return
	}
	return
}
func (s *MuxServer) GetConn(index string) (conn net.Conn, err error) {
	select {
	case s.lockChn <- true:
	default:
		err = fmt.Errorf("can not connect at same time")
		return
	}
	defer func() {
		<-s.lockChn
	}()
	var session *smux.Session
	_session, ok := s.sessions.Get(index)
	if !ok {
		var c net.Conn
		c, err = s.getParentConn()
		if err != nil {
			return
		}
		c.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
		_, err = c.Write(utils.BuildPacket(CONN_SERVER, *s.cfg.Key, s.cfg.Mgr.serverID))
		c.SetDeadline(time.Time{})
		if err != nil {
			c.Close()
			return
		}
		if err == nil {
			session, err = smux.Client(c, nil)
			if err != nil {
				return
			}
		}
		if _sess, ok := s.sessions.Get(index); ok {
			_sess.(*smux.Session).Close()
		}
		s.sessions.Set(index, session)
		log.Printf("session[%s] created", index)
		go func() {
			for {
				if s.isStop {
					return
				}
				if session.IsClosed() {
					s.sessions.Remove(index)
					break
				}
				time.Sleep(time.Second * 5)
			}
		}()
	} else {
		session = _session.(*smux.Session)
	}
	conn, err = session.OpenStream()
	if err != nil {
		session.Close()
		s.sessions.Remove(index)
	}
	return
}
func (s *MuxServer) getParentConn() (conn net.Conn, err error) {
	if *s.cfg.ParentType == "tls" {
		var _conn tls.Conn
		_conn, err = utils.TlsConnectHost(*s.cfg.Parent, *s.cfg.Timeout, s.cfg.CertBytes, s.cfg.KeyBytes, nil)
		if err == nil {
			conn = net.Conn(&_conn)
		}
	} else if *s.cfg.ParentType == "kcp" {
		conn, err = utils.ConnectKCPHost(*s.cfg.Parent, s.cfg.KCP)
	} else {
		conn, err = utils.ConnectHost(*s.cfg.Parent, *s.cfg.Timeout)
	}
	return
}
func (s *MuxServer) UDPConnDeamon() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("udp conn deamon crashed with err : %s \nstack: %s", err, string(debug.Stack()))
			}
		}()
		var outConn net.Conn
		var ID string
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
					outConn, ID, err = s.GetOutConn()
					if err != nil {
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
								outConn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
								srcAddrFromConn, body, err := utils.ReadUDPPacket(outConn)
								outConn.SetDeadline(time.Time{})
								if err != nil {
									log.Printf("parse revecived udp packet fail, err: %s ,%v", err, body)
									log.Printf("UDP deamon connection %s exited", ID)
									break
								}
								//log.Printf("udp packet revecived over parent , local:%s", srcAddrFromConn)
								_srcAddr := strings.Split(srcAddrFromConn, ":")
								if len(_srcAddr) != 2 {
									log.Printf("parse revecived udp packet fail, addr error : %s", srcAddrFromConn)
									continue
								}
								port, _ := strconv.Atoi(_srcAddr[1])
								dstAddr := &net.UDPAddr{IP: net.ParseIP(_srcAddr[0]), Port: port}
								s.sc.UDPListener.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
								_, err = s.sc.UDPListener.WriteToUDP(body, dstAddr)
								s.sc.UDPListener.SetDeadline(time.Time{})
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

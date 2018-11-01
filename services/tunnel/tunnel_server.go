package tunnel

import (
	"crypto/tls"
	"fmt"
	"io"
	logger "log"
	"net"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/snail007/goproxy/core/cs/server"
	"github.com/snail007/goproxy/services"
	"github.com/snail007/goproxy/utils"
	"github.com/snail007/goproxy/utils/jumper"
	"github.com/snail007/goproxy/utils/mapx"

	//"github.com/xtaci/smux"
	smux "github.com/hashicorp/yamux"
)

type TunnelServerArgs struct {
	Parent    *string
	CertFile  *string
	KeyFile   *string
	CertBytes []byte
	KeyBytes  []byte
	Local     *string
	IsUDP     *bool
	Key       *string
	Remote    *string
	Timeout   *int
	Route     *[]string
	Mgr       *TunnelServerManager
	Jumper    *string
}
type TunnelServer struct {
	cfg       TunnelServerArgs
	sc        server.ServerChannel
	isStop    bool
	udpConn   *net.Conn
	userConns mapx.ConcurrentMap
	log       *logger.Logger
	jumper    *jumper.Jumper
	udpConns  mapx.ConcurrentMap
}

type TunnelServerManager struct {
	cfg      TunnelServerArgs
	serverID string
	servers  []*services.Service
	log      *logger.Logger
}

func NewTunnelServerManager() services.Service {
	return &TunnelServerManager{
		cfg:      TunnelServerArgs{},
		serverID: utils.Uniqueid(),
		servers:  []*services.Service{},
	}
}
func (s *TunnelServerManager) Start(args interface{}, log *logger.Logger) (err error) {
	s.log = log
	s.cfg = args.(TunnelServerArgs)
	if err = s.CheckArgs(); err != nil {
		return
	}
	if *s.cfg.Parent != "" {
		s.log.Printf("use tls parent %s", *s.cfg.Parent)
	} else {
		err = fmt.Errorf("parent required")
		return
	}

	if err = s.InitService(); err != nil {
		return
	}

	s.log.Printf("server id: %s", s.serverID)
	//s.log.Printf("route:%v", *s.cfg.Route)
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
			Jumper:    s.cfg.Jumper,
		}, log)

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
	s.cfg = TunnelServerArgs{}
	s.log = nil
	s.serverID = ""
	s.servers = nil
	s = nil
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

func NewTunnelServer() services.Service {
	return &TunnelServer{
		cfg:       TunnelServerArgs{},
		isStop:    false,
		userConns: mapx.NewConcurrentMap(),
		udpConns:  mapx.NewConcurrentMap(),
	}
}

type TunnelUDPPacketItem struct {
	packet    *[]byte
	localAddr *net.UDPAddr
	srcAddr   *net.UDPAddr
}
type TunnelUDPConnItem struct {
	conn      *net.Conn
	isActive  bool
	touchtime int64
	srcAddr   *net.UDPAddr
	localAddr *net.UDPAddr
	connid    string
}

func (s *TunnelServer) StopService() {
	defer func() {
		e := recover()
		if e != nil {
			s.log.Printf("stop server service crashed,%s", e)
		} else {
			s.log.Printf("service server stopped")
		}
		s.cfg = TunnelServerArgs{}
		s.jumper = nil
		s.log = nil
		s.sc = server.ServerChannel{}
		s.udpConn = nil
		s.udpConns = nil
		s.userConns = nil
		s = nil
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
	s.UDPGCDeamon()
	return
}
func (s *TunnelServer) CheckArgs() (err error) {
	if *s.cfg.Remote == "" {
		err = fmt.Errorf("remote required")
		return
	}
	if *s.cfg.Jumper != "" {
		var j jumper.Jumper
		j, err = jumper.New(*s.cfg.Jumper, time.Millisecond*time.Duration(*s.cfg.Timeout))
		if err != nil {
			err = fmt.Errorf("parse jumper fail, err %s", err)
			return
		}
		s.jumper = &j
	}
	return
}

func (s *TunnelServer) Start(args interface{}, log *logger.Logger) (err error) {
	s.log = log
	s.cfg = args.(TunnelServerArgs)
	if err = s.CheckArgs(); err != nil {
		return
	}
	if err = s.InitService(); err != nil {
		return
	}
	host, port, _ := net.SplitHostPort(*s.cfg.Local)
	p, _ := strconv.Atoi(port)
	s.sc = server.NewServerChannel(host, p, s.log)
	if *s.cfg.IsUDP {
		err = s.sc.ListenUDP(func(listener *net.UDPConn, packet []byte, localAddr, srcAddr *net.UDPAddr) {
			s.UDPSend(packet, localAddr, srcAddr)
		})
		if err != nil {
			return
		}
		s.log.Printf("proxy on udp tunnel server mode %s", (*s.sc.UDPListener).LocalAddr())
	} else {
		err = s.sc.ListenTCP(func(inConn net.Conn) {
			defer func() {
				if err := recover(); err != nil {
					s.log.Printf("tserver conn handler crashed with err : %s \nstack: %s", err, string(debug.Stack()))
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
					s.log.Printf("connect to %s fail, err: %s, retrying...", *s.cfg.Parent, err)
					time.Sleep(time.Second * 3)
					continue
				} else {
					break
				}
			}
			inAddr := inConn.RemoteAddr().String()
			utils.IoBind(inConn, outConn, func(err interface{}) {
				s.userConns.Remove(inAddr)
				s.log.Printf("%s conn %s released", *s.cfg.Key, ID)
			}, s.log)
			if c, ok := s.userConns.Get(inAddr); ok {
				(*c.(*net.Conn)).Close()
			}
			s.userConns.Set(inAddr, &inConn)
			s.log.Printf("%s conn %s created", *s.cfg.Key, ID)
		})
		if err != nil {
			return
		}
		s.log.Printf("proxy on tunnel server mode %s", (*s.sc.Listener).Addr())
	}
	return
}
func (s *TunnelServer) Clean() {

}
func (s *TunnelServer) GetOutConn(typ uint8) (outConn net.Conn, ID string, err error) {
	outConn, err = s.GetConn()
	if err != nil {
		s.log.Printf("connection err: %s", err)
		return
	}
	remoteAddr := "tcp:" + *s.cfg.Remote
	if *s.cfg.IsUDP {
		remoteAddr = "udp:" + *s.cfg.Remote
	}
	ID = utils.Uniqueid()
	_, err = outConn.Write(utils.BuildPacket(typ, *s.cfg.Key, ID, remoteAddr, s.cfg.Mgr.serverID))
	if err != nil {
		s.log.Printf("write connection data err: %s ,retrying...", err)
		utils.CloseConn(&outConn)
		return
	}
	return
}
func (s *TunnelServer) GetConn() (conn net.Conn, err error) {
	var dconn net.Conn
	if s.jumper == nil {
		var _conn tls.Conn
		_conn, err = utils.TlsConnectHost(*s.cfg.Parent, *s.cfg.Timeout, s.cfg.CertBytes, s.cfg.KeyBytes, nil)
		if err == nil {
			dconn = net.Conn(&_conn)
		}
	} else {
		conf, err := utils.TlsConfig(s.cfg.CertBytes, s.cfg.KeyBytes, nil)
		if err != nil {
			return nil, err
		}
		var _c net.Conn
		_c, err = s.jumper.Dial(*s.cfg.Parent, time.Millisecond*time.Duration(*s.cfg.Timeout))
		if err == nil {
			dconn = net.Conn(tls.Client(_c, conf))
		}
	}
	if err == nil {
		sess, e := smux.Client(dconn, &smux.Config{
			AcceptBacklog:          256,
			EnableKeepAlive:        true,
			KeepAliveInterval:      9 * time.Second,
			ConnectionWriteTimeout: 3 * time.Second,
			MaxStreamWindowSize:    512 * 1024,
			LogOutput:              os.Stderr,
		})
		if e != nil {
			s.log.Printf("new mux client conn error,ERR:%s", e)
			err = e
			dconn.Close()
			return
		}
		conn, e = sess.OpenStream()
		if e != nil {
			s.log.Printf("mux client conn open stream error,ERR:%s", e)
			err = e
			dconn.Close()
			return
		}
		go func() {
			defer func() {
				_ = recover()
			}()
			timer := time.NewTicker(time.Second * 3)
			for {
				<-timer.C
				if sess.NumStreams() == 0 {
					sess.Close()
					timer.Stop()
					return
				}
			}
		}()
	}
	return
}
func (s *TunnelServer) UDPGCDeamon() {
	gctime := int64(30)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("crashed, err: %s\nstack:\n%s", e, string(debug.Stack()))
			}
		}()
		if s.isStop {
			return
		}
		timer := time.NewTicker(time.Second)
		for {
			<-timer.C
			gcKeys := []string{}
			s.udpConns.IterCb(func(key string, v interface{}) {
				if time.Now().Unix()-v.(*TunnelUDPConnItem).touchtime > gctime {
					(*(v.(*TunnelUDPConnItem).conn)).Close()
					gcKeys = append(gcKeys, key)
					s.log.Printf("gc udp conn %s", v.(*TunnelUDPConnItem).connid)
				}
			})
			for _, k := range gcKeys {
				s.udpConns.Remove(k)
			}
			gcKeys = nil
		}
	}()
}
func (s *TunnelServer) UDPSend(data []byte, localAddr, srcAddr *net.UDPAddr) {
	var (
		uc      *TunnelUDPConnItem
		key     = srcAddr.String()
		ID      string
		err     error
		outconn net.Conn
	)
	v, ok := s.udpConns.Get(key)
	if !ok {
		outconn, ID, err = s.GetOutConn(CONN_SERVER)
		if err != nil {
			s.log.Printf("connect to %s fail, err: %s", *s.cfg.Parent, err)
			return
		}
		uc = &TunnelUDPConnItem{
			conn:      &outconn,
			srcAddr:   srcAddr,
			localAddr: localAddr,
			connid:    ID,
		}
		s.udpConns.Set(key, uc)
		s.UDPRevecive(key, ID)
	} else {
		uc = v.(*TunnelUDPConnItem)
	}
	go func() {
		defer func() {
			if e := recover(); e != nil {
				(*uc.conn).Close()
				s.udpConns.Remove(key)
				s.log.Printf("udp sender crashed with error : %s", e)
			}
		}()
		uc.touchtime = time.Now().Unix()
		(*uc.conn).SetWriteDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
		_, err = (*uc.conn).Write(utils.UDPPacket(srcAddr.String(), data))
		(*uc.conn).SetWriteDeadline(time.Time{})
		if err != nil {
			s.log.Printf("write udp packet to %s fail ,flush err:%s ", *s.cfg.Parent, err)
		}
	}()
}
func (s *TunnelServer) UDPRevecive(key, ID string) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("crashed, err: %s\nstack:\n%s", e, string(debug.Stack()))
			}
		}()
		s.log.Printf("udp conn %s connected", ID)
		var uc *TunnelUDPConnItem
		defer func() {
			if uc != nil {
				(*uc.conn).Close()
			}
			s.udpConns.Remove(key)
			s.log.Printf("udp conn %s released", ID)
		}()
		v, ok := s.udpConns.Get(key)
		if !ok {
			s.log.Printf("[warn] udp conn not exists for %s, connid : %s", key, ID)
			return
		}
		uc = v.(*TunnelUDPConnItem)
		for {
			_, body, err := utils.ReadUDPPacket(*uc.conn)
			if err != nil {
				if strings.Contains(err.Error(), "n != int(") {
					continue
				}
				if err != io.EOF {
					s.log.Printf("udp conn read udp packet fail , err: %s ", err)
				}
				return
			}
			uc.touchtime = time.Now().Unix()
			go func() {
				defer func() {
					if e := recover(); e != nil {
						fmt.Printf("crashed, err: %s\nstack:\n%s", e, string(debug.Stack()))
					}
				}()
				s.sc.UDPListener.WriteToUDP(body, uc.srcAddr)
			}()
		}
	}()
}

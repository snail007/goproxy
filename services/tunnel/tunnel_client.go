package tunnel

import (
	"crypto/tls"
	"fmt"
	"io"
	logger "log"
	"net"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/snail007/goproxy/services"
	"github.com/snail007/goproxy/utils"
	"github.com/snail007/goproxy/utils/jumper"
	"github.com/snail007/goproxy/utils/mapx"

	//"github.com/xtaci/smux"
	smux "github.com/hashicorp/yamux"
)

type TunnelClientArgs struct {
	Parent    *string
	CertFile  *string
	KeyFile   *string
	CertBytes []byte
	KeyBytes  []byte
	Key       *string
	Timeout   *int
	Jumper    *string
}
type ClientUDPConnItem struct {
	conn      *net.Conn
	isActive  bool
	touchtime int64
	srcAddr   *net.UDPAddr
	localAddr *net.UDPAddr
	udpConn   *net.UDPConn
	connid    string
}
type TunnelClient struct {
	cfg       TunnelClientArgs
	ctrlConn  *net.Conn
	isStop    bool
	userConns mapx.ConcurrentMap
	log       *logger.Logger
	jumper    *jumper.Jumper
	udpConns  mapx.ConcurrentMap
}

func NewTunnelClient() services.Service {
	return &TunnelClient{
		cfg:       TunnelClientArgs{},
		userConns: mapx.NewConcurrentMap(),
		isStop:    false,
		udpConns:  mapx.NewConcurrentMap(),
	}
}

func (s *TunnelClient) InitService() (err error) {
	s.UDPGCDeamon()
	return
}

func (s *TunnelClient) CheckArgs() (err error) {
	if *s.cfg.Parent != "" {
		s.log.Printf("use tls parent %s", *s.cfg.Parent)
	} else {
		err = fmt.Errorf("parent required")
		return
	}
	if *s.cfg.CertFile == "" || *s.cfg.KeyFile == "" {
		err = fmt.Errorf("cert and key file required")
		return
	}
	s.cfg.CertBytes, s.cfg.KeyBytes, err = utils.TlsBytes(*s.cfg.CertFile, *s.cfg.KeyFile)
	if err != nil {
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
func (s *TunnelClient) StopService() {
	defer func() {
		e := recover()
		if e != nil {
			s.log.Printf("stop tclient service crashed,%s", e)
		} else {
			s.log.Printf("service tclient stopped")
		}
		s.cfg = TunnelClientArgs{}
		s.ctrlConn = nil
		s.jumper = nil
		s.log = nil
		s.udpConns = nil
		s.userConns = nil
		s = nil
	}()
	s.isStop = true
	if s.ctrlConn != nil {
		(*s.ctrlConn).Close()
	}
	for _, c := range s.userConns.Items() {
		(*c.(*net.Conn)).Close()
	}
}
func (s *TunnelClient) Start(args interface{}, log *logger.Logger) (err error) {
	s.log = log
	s.cfg = args.(TunnelClientArgs)
	if err = s.CheckArgs(); err != nil {
		return
	}
	if err = s.InitService(); err != nil {
		return
	}
	s.log.Printf("proxy on tunnel client mode")

	for {
		if s.isStop {
			return
		}
		if s.ctrlConn != nil {
			(*s.ctrlConn).Close()
		}
		var c net.Conn
		c, err = s.GetInConn(CONN_CLIENT_CONTROL, *s.cfg.Key)
		if err != nil {
			s.log.Printf("control connection err: %s, retrying...", err)
			time.Sleep(time.Second * 3)
			continue
		}
		s.ctrlConn = &c
		for {
			if s.isStop {
				return
			}
			var ID, clientLocalAddr, serverID string
			err = utils.ReadPacketData(*s.ctrlConn, &ID, &clientLocalAddr, &serverID)
			if err != nil {
				if s.ctrlConn != nil {
					(*s.ctrlConn).Close()
				}
				s.log.Printf("read connection signal err: %s, retrying...", err)
				break
			}
			//s.log.Printf("signal revecived:%s %s %s", serverID, ID, clientLocalAddr)
			protocol := clientLocalAddr[:3]
			localAddr := clientLocalAddr[4:]
			if protocol == "udp" {
				go func() {
					defer func() {
						if e := recover(); e != nil {
							fmt.Printf("crashed, err: %s\nstack:\n%s",e, string(debug.Stack()))
						}
					}()
					s.ServeUDP(localAddr, ID, serverID)
				}()
			} else {
				go func() {
					defer func() {
						if e := recover(); e != nil {
							fmt.Printf("crashed, err: %s\nstack:\n%s",e, string(debug.Stack()))
						}
					}()
					s.ServeConn(localAddr, ID, serverID)
				}()
			}
		}
	}
}
func (s *TunnelClient) Clean() {
	s.StopService()
}
func (s *TunnelClient) GetInConn(typ uint8, data ...string) (outConn net.Conn, err error) {
	outConn, err = s.GetConn()
	if err != nil {
		err = fmt.Errorf("connection err: %s", err)
		return
	}
	_, err = outConn.Write(utils.BuildPacket(typ, data...))
	if err != nil {
		err = fmt.Errorf("write connection data err: %s ,retrying...", err)
		utils.CloseConn(&outConn)
		return
	}
	return
}
func (s *TunnelClient) GetConn() (conn net.Conn, err error) {
	if s.jumper == nil {
		var _conn tls.Conn
		_conn, err = utils.TlsConnectHost(*s.cfg.Parent, *s.cfg.Timeout, s.cfg.CertBytes, s.cfg.KeyBytes, nil)
		if err == nil {
			conn = net.Conn(&_conn)
		}
	} else {
		conf, err := utils.TlsConfig(s.cfg.CertBytes, s.cfg.KeyBytes, nil)
		if err != nil {
			return nil, err
		}
		var _c net.Conn
		_c, err = s.jumper.Dial(*s.cfg.Parent, time.Millisecond*time.Duration(*s.cfg.Timeout))
		if err == nil {
			conn = net.Conn(tls.Client(_c, conf))
		}
	}
	if err == nil {
		sess, e := smux.Client(conn, &smux.Config{
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
			return
		}
		conn, e = sess.OpenStream()
		if e != nil {
			s.log.Printf("mux client conn open stream error,ERR:%s", e)
			err = e
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
func (s *TunnelClient) ServeUDP(localAddr, ID, serverID string) {
	var inConn net.Conn
	var err error
	// for {
	for {
		if s.isStop {
			if inConn != nil {
				inConn.Close()
			}
			return
		}
		// s.cm.RemoveOne(*s.cfg.Key, ID)
		inConn, err = s.GetInConn(CONN_CLIENT, *s.cfg.Key, ID, serverID)
		if err != nil {
			utils.CloseConn(&inConn)
			s.log.Printf("connection err: %s, retrying...", err)
			time.Sleep(time.Second * 3)
			continue
		} else {
			break
		}
	}
	// s.cm.Add(*s.cfg.Key, ID, &inConn)
	s.log.Printf("conn %s created", ID)
	var item *ClientUDPConnItem
	var body []byte
	srcAddr := ""
	defer func() {
		if item != nil {
			(*(*item).conn).Close()
			(*item).udpConn.Close()
			s.udpConns.Remove(srcAddr)
			inConn.Close()
		}
	}()
	for {
		if s.isStop {
			return
		}
		srcAddr, body, err = utils.ReadUDPPacket(inConn)
		if err != nil {
			if strings.Contains(err.Error(), "n != int(") {
				continue
			}
			if !utils.IsNetDeadlineErr(err) && err != io.EOF {
				s.log.Printf("udp packet revecived from bridge fail, err: %s", err)
			}
			return
		}
		if v, ok := s.udpConns.Get(srcAddr); !ok {
			_srcAddr, _ := net.ResolveUDPAddr("udp", srcAddr)
			zeroAddr, _ := net.ResolveUDPAddr("udp", ":")
			_localAddr, _ := net.ResolveUDPAddr("udp", localAddr)
			c, err := net.DialUDP("udp", zeroAddr, _localAddr)
			if err != nil {
				s.log.Printf("create local udp conn fail, err : %s", err)
				inConn.Close()
				return
			}
			item = &ClientUDPConnItem{
				conn:      &inConn,
				srcAddr:   _srcAddr,
				localAddr: _localAddr,
				udpConn:   c,
				connid:    ID,
			}
			s.udpConns.Set(srcAddr, item)
			s.UDPRevecive(srcAddr, ID)
		} else {
			item = v.(*ClientUDPConnItem)
		}
		(*item).touchtime = time.Now().Unix()
		go func() {
			defer func() {
				if e := recover(); e != nil {
					fmt.Printf("crashed, err: %s\nstack:\n%s",e, string(debug.Stack()))
				}
			}()
			(*item).udpConn.Write(body)
		}()
	}
}
func (s *TunnelClient) UDPRevecive(key, ID string) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("crashed, err: %s\nstack:\n%s",e, string(debug.Stack()))
			}
		}()
		s.log.Printf("udp conn %s connected", ID)
		v, ok := s.udpConns.Get(key)
		if !ok {
			s.log.Printf("[warn] udp conn not exists for %s, connid : %s", key, ID)
			return
		}
		cui := v.(*ClientUDPConnItem)
		buf := utils.LeakyBuffer.Get()
		defer func() {
			utils.LeakyBuffer.Put(buf)
			(*cui.conn).Close()
			cui.udpConn.Close()
			s.udpConns.Remove(key)
			s.log.Printf("udp conn %s released", ID)
		}()
		for {
			n, err := cui.udpConn.Read(buf)
			if err != nil {
				if !utils.IsNetClosedErr(err) {
					s.log.Printf("udp conn read udp packet fail , err: %s ", err)
				}
				return
			}
			cui.touchtime = time.Now().Unix()
			go func() {
				defer func() {
					if e := recover(); e != nil {
						fmt.Printf("crashed, err: %s\nstack:\n%s",e, string(debug.Stack()))
					}
				}()
				(*cui.conn).SetWriteDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
				_, err = (*cui.conn).Write(utils.UDPPacket(cui.srcAddr.String(), buf[:n]))
				(*cui.conn).SetWriteDeadline(time.Time{})
				if err != nil {
					cui.udpConn.Close()
					return
				}
			}()
		}
	}()
}
func (s *TunnelClient) UDPGCDeamon() {
	gctime := int64(30)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("crashed, err: %s\nstack:\n%s",e, string(debug.Stack()))
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
				if time.Now().Unix()-v.(*ClientUDPConnItem).touchtime > gctime {
					(*(v.(*ClientUDPConnItem).conn)).Close()
					(v.(*ClientUDPConnItem).udpConn).Close()
					gcKeys = append(gcKeys, key)
					s.log.Printf("gc udp conn %s", v.(*ClientUDPConnItem).connid)
				}
			})
			for _, k := range gcKeys {
				s.udpConns.Remove(k)
			}
			gcKeys = nil
		}
	}()
}
func (s *TunnelClient) ServeConn(localAddr, ID, serverID string) {
	var inConn, outConn net.Conn
	var err error
	for {
		if s.isStop {
			return
		}
		inConn, err = s.GetInConn(CONN_CLIENT, *s.cfg.Key, ID, serverID)
		if err != nil {
			utils.CloseConn(&inConn)
			s.log.Printf("connection err: %s, retrying...", err)
			time.Sleep(time.Second * 3)
			continue
		} else {
			break
		}
	}

	i := 0
	for {
		if s.isStop {
			return
		}
		i++
		outConn, err = utils.ConnectHost(localAddr, *s.cfg.Timeout)
		if err == nil || i == 3 {
			break
		} else {
			if i == 3 {
				s.log.Printf("connect to %s err: %s, retrying...", localAddr, err)
				time.Sleep(2 * time.Second)
				continue
			}
		}
	}
	if err != nil {
		utils.CloseConn(&inConn)
		utils.CloseConn(&outConn)
		s.log.Printf("build connection error, err: %s", err)
		return
	}
	inAddr := inConn.RemoteAddr().String()
	utils.IoBind(inConn, outConn, func(err interface{}) {
		s.log.Printf("conn %s released", ID)
		s.userConns.Remove(inAddr)
	}, s.log)
	if c, ok := s.userConns.Get(inAddr); ok {
		(*c.(*net.Conn)).Close()
	}
	s.userConns.Set(inAddr, &inConn)
	s.log.Printf("conn %s created", ID)
}

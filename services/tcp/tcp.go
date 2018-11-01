package tcp

import (
	"crypto/tls"
	"fmt"
	logger "log"
	"net"
	"runtime/debug"
	"strings"
	"time"

	"github.com/snail007/goproxy/core/cs/server"
	"github.com/snail007/goproxy/core/lib/kcpcfg"
	"github.com/snail007/goproxy/services"
	"github.com/snail007/goproxy/utils"
	"github.com/snail007/goproxy/utils/jumper"
	"github.com/snail007/goproxy/utils/mapx"

	"strconv"
)

type TCPArgs struct {
	Parent              *string
	CertFile            *string
	KeyFile             *string
	CertBytes           []byte
	KeyBytes            []byte
	Local               *string
	ParentType          *string
	LocalType           *string
	Timeout             *int
	CheckParentInterval *int
	KCP                 kcpcfg.KCPConfigArgs
	Jumper              *string
}
type UDPConnItem struct {
	conn      *net.Conn
	isActive  bool
	touchtime int64
	srcAddr   *net.UDPAddr
	localAddr *net.UDPAddr
	udpConn   *net.UDPConn
	connid    string
}
type TCP struct {
	cfg       TCPArgs
	sc        *server.ServerChannel
	isStop    bool
	userConns mapx.ConcurrentMap
	log       *logger.Logger
	jumper    *jumper.Jumper
	udpConns  mapx.ConcurrentMap
}

func NewTCP() services.Service {
	return &TCP{
		cfg:       TCPArgs{},
		isStop:    false,
		userConns: mapx.NewConcurrentMap(),
		udpConns:  mapx.NewConcurrentMap(),
	}
}
func (s *TCP) CheckArgs() (err error) {
	if len(*s.cfg.Parent) == 0 {
		err = fmt.Errorf("parent required for %s %s", *s.cfg.LocalType, *s.cfg.Local)
		return
	}
	if *s.cfg.ParentType == "" {
		err = fmt.Errorf("parent type unkown,use -T <tls|tcp|kcp|udp>")
		return
	}
	if *s.cfg.ParentType == "tls" || *s.cfg.LocalType == "tls" {
		s.cfg.CertBytes, s.cfg.KeyBytes, err = utils.TlsBytes(*s.cfg.CertFile, *s.cfg.KeyFile)
		if err != nil {
			return
		}
	}
	if *s.cfg.Jumper != "" {
		if *s.cfg.ParentType != "tls" && *s.cfg.ParentType != "tcp" {
			err = fmt.Errorf("jumper only worked of -T is tls or tcp")
			return
		}
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
func (s *TCP) InitService() (err error) {
	s.UDPGCDeamon()
	return
}
func (s *TCP) StopService() {
	defer func() {
		e := recover()
		if e != nil {
			s.log.Printf("stop tcp service crashed,%s", e)
		} else {
			s.log.Printf("service tcp stopped")
		}
		s.cfg = TCPArgs{}
		s.jumper = nil
		s.log = nil
		s.sc = nil
		s.userConns = nil
		s = nil
	}()
	s.isStop = true
	if s.sc.Listener != nil && *s.sc.Listener != nil {
		(*s.sc.Listener).Close()
	}
	if s.sc.UDPListener != nil {
		(*s.sc.UDPListener).Close()
	}
	for _, c := range s.userConns.Items() {
		(*c.(*net.Conn)).Close()
	}
}
func (s *TCP) Start(args interface{}, log *logger.Logger) (err error) {
	s.log = log
	s.cfg = args.(TCPArgs)
	if err = s.CheckArgs(); err != nil {
		return
	}
	if err = s.InitService(); err != nil {
		return
	}
	s.log.Printf("use %s parent %v", *s.cfg.ParentType, *s.cfg.Parent)
	host, port, _ := net.SplitHostPort(*s.cfg.Local)
	p, _ := strconv.Atoi(port)
	sc := server.NewServerChannel(host, p, s.log)

	if *s.cfg.LocalType == "tcp" {
		err = sc.ListenTCP(s.callback)
	} else if *s.cfg.LocalType == "tls" {
		err = sc.ListenTLS(s.cfg.CertBytes, s.cfg.KeyBytes, nil, s.callback)
	} else if *s.cfg.LocalType == "kcp" {
		err = sc.ListenKCP(s.cfg.KCP, s.callback, s.log)
	}
	if err != nil {
		return
	}
	s.log.Printf("%s proxy on %s", *s.cfg.LocalType, (*sc.Listener).Addr())
	s.sc = &sc
	return
}

func (s *TCP) Clean() {
	s.StopService()
}
func (s *TCP) callback(inConn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			s.log.Printf("%s conn handler crashed with err : %s \nstack: %s", *s.cfg.LocalType, err, string(debug.Stack()))
		}
	}()
	var err error
	lbAddr := ""
	switch *s.cfg.ParentType {
	case "kcp", "tcp", "tls":
		err = s.OutToTCP(&inConn)
	case "udp":
		s.OutToUDP(&inConn)
	default:
		err = fmt.Errorf("unkown parent type %s", *s.cfg.ParentType)
	}
	if err != nil {
		if !utils.IsNetClosedErr(err) {
			s.log.Printf("connect to %s parent %s fail, ERR:%s", *s.cfg.ParentType, lbAddr, err)
		}
		utils.CloseConn(&inConn)
	}
}
func (s *TCP) OutToTCP(inConn *net.Conn) (err error) {

	var outConn net.Conn
	outConn, err = s.GetParentConn()
	if err != nil {
		s.log.Printf("connect to %s , err:%s", *s.cfg.Parent, err)
		utils.CloseConn(inConn)
		return
	}
	inAddr := (*inConn).RemoteAddr().String()
	//inLocalAddr := (*inConn).LocalAddr().String()
	outAddr := outConn.RemoteAddr().String()
	//outLocalAddr := outConn.LocalAddr().String()
	utils.IoBind((*inConn), outConn, func(err interface{}) {
		s.log.Printf("conn %s - %s released", inAddr, outAddr)
		s.userConns.Remove(inAddr)
	}, s.log)
	s.log.Printf("conn %s - %s connected", inAddr, outAddr)
	if c, ok := s.userConns.Get(inAddr); ok {
		(*c.(*net.Conn)).Close()
	}
	s.userConns.Set(inAddr, inConn)
	return
}
func (s *TCP) OutToUDP(inConn *net.Conn) (err error) {
	var item *UDPConnItem
	var body []byte
	srcAddr := ""
	defer func() {
		if item != nil {
			(*item.conn).Close()
			item.udpConn.Close()
			s.udpConns.Remove(srcAddr)
			(*inConn).Close()
		}
	}()
	for {
		if s.isStop {
			return
		}
		var srcAddr string
		srcAddr, body, err = utils.ReadUDPPacket(*inConn)
		if err != nil {
			if strings.Contains(err.Error(), "n != int(") {
				continue
			}
			// if !utils.IsNetDeadlineErr(err) && err != io.EOF && !utils.IsNetClosedErr(err) {
			// 	s.log.Printf("udp packet revecived from client fail, err: %s", err)
			// }
			return
		}
		localAddr := *s.cfg.Parent
		if v, ok := s.udpConns.Get(srcAddr); !ok {
			_srcAddr, _ := net.ResolveUDPAddr("udp", srcAddr)
			zeroAddr, _ := net.ResolveUDPAddr("udp", ":")
			_localAddr, _ := net.ResolveUDPAddr("udp", localAddr)
			var c *net.UDPConn
			c, err = net.DialUDP("udp", zeroAddr, _localAddr)
			if err != nil {
				s.log.Printf("create local udp conn fail, err : %s", err)
				(*inConn).Close()
				return
			}
			item = &UDPConnItem{
				conn:      inConn,
				srcAddr:   _srcAddr,
				localAddr: _localAddr,
				udpConn:   c,
			}
			s.udpConns.Set(srcAddr, item)
			s.UDPRevecive(srcAddr)
		} else {
			item = v.(*UDPConnItem)
		}
		item.touchtime = time.Now().Unix()
		go item.udpConn.Write(body)
	}
}
func (s *TCP) UDPRevecive(key string) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("crashed, err: %s\nstack:%s", e, string(debug.Stack()))
			}
		}()
		s.log.Printf("udp conn %s connected", key)
		v, ok := s.udpConns.Get(key)
		if !ok {
			s.log.Printf("[warn] udp conn not exists for %s", key)
			return
		}
		cui := v.(*UDPConnItem)
		buf := utils.LeakyBuffer.Get()
		defer func() {
			utils.LeakyBuffer.Put(buf)
			(*cui.conn).Close()
			cui.udpConn.Close()
			s.udpConns.Remove(key)
			s.log.Printf("udp conn %s released", key)
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
						fmt.Printf("crashed, err: %s\nstack:%s", e, string(debug.Stack()))
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
func (s *TCP) UDPGCDeamon() {
	gctime := int64(30)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("crashed, err: %s\nstack:%s", e, string(debug.Stack()))
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
				if time.Now().Unix()-v.(*UDPConnItem).touchtime > gctime {
					(*(v.(*UDPConnItem).conn)).Close()
					(v.(*UDPConnItem).udpConn).Close()
					gcKeys = append(gcKeys, key)
					s.log.Printf("gc udp conn %s", key)
				}
			})
			for _, k := range gcKeys {
				s.udpConns.Remove(k)
			}
			gcKeys = nil
		}
	}()
}
func (s *TCP) GetParentConn() (conn net.Conn, err error) {
	if *s.cfg.ParentType == "tls" {
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

	} else if *s.cfg.ParentType == "kcp" {
		conn, err = utils.ConnectKCPHost(*s.cfg.Parent, s.cfg.KCP)
	} else {
		if s.jumper == nil {
			conn, err = utils.ConnectHost(*s.cfg.Parent, *s.cfg.Timeout)
		} else {
			conn, err = s.jumper.Dial(*s.cfg.Parent, time.Millisecond*time.Duration(*s.cfg.Timeout))
		}
	}
	return
}

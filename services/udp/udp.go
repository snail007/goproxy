package udp

import (
	"crypto/tls"
	"fmt"
	"io"
	logger "log"
	"net"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/snail007/goproxy/core/cs/server"
	"github.com/snail007/goproxy/services"
	"github.com/snail007/goproxy/utils"
	"github.com/snail007/goproxy/utils/mapx"
)

type UDPArgs struct {
	Parent              *string
	CertFile            *string
	KeyFile             *string
	CertBytes           []byte
	KeyBytes            []byte
	Local               *string
	ParentType          *string
	Timeout             *int
	CheckParentInterval *int
}
type UDP struct {
	p                mapx.ConcurrentMap
	cfg              UDPArgs
	sc               *server.ServerChannel
	isStop           bool
	log              *logger.Logger
	outUDPConnCtxMap mapx.ConcurrentMap
	udpConns         mapx.ConcurrentMap
	dstAddr          *net.UDPAddr
}
type UDPConnItem struct {
	conn      *net.Conn
	touchtime int64
	srcAddr   *net.UDPAddr
	localAddr *net.UDPAddr
	connid    string
}
type outUDPConnCtx struct {
	localAddr *net.UDPAddr
	srcAddr   *net.UDPAddr
	udpconn   *net.UDPConn
	touchtime int64
}

func NewUDP() services.Service {
	return &UDP{
		p:                mapx.NewConcurrentMap(),
		isStop:           false,
		outUDPConnCtxMap: mapx.NewConcurrentMap(),
		udpConns:         mapx.NewConcurrentMap(),
	}
}
func (s *UDP) CheckArgs() (err error) {
	if len(*s.cfg.Parent) == 0 {
		err = fmt.Errorf("parent required for udp %s", *s.cfg.Local)
		return
	}
	if *s.cfg.ParentType == "" {
		err = fmt.Errorf("parent type unkown,use -T <udp|tls|tcp>")
		return
	}
	if *s.cfg.ParentType == "tls" {
		s.cfg.CertBytes, s.cfg.KeyBytes, err = utils.TlsBytes(*s.cfg.CertFile, *s.cfg.KeyFile)
		if err != nil {
			return
		}
	}

	s.dstAddr, err = net.ResolveUDPAddr("udp", *s.cfg.Parent)
	if err != nil {
		s.log.Printf("resolve udp addr %s fail  fail,ERR:%s", *s.cfg.Parent, err)
		return
	}
	return
}
func (s *UDP) InitService() (err error) {
	s.OutToUDPGCDeamon()
	s.UDPGCDeamon()
	return
}
func (s *UDP) StopService() {
	defer func() {
		e := recover()
		if e != nil {
			s.log.Printf("stop udp service crashed,%s", e)
		} else {
			s.log.Printf("service udp stopped")
		}
		s.cfg = UDPArgs{}
		s.log = nil
		s.p = nil
		s.sc = nil
		s = nil
	}()
	s.isStop = true
	if s.sc.Listener != nil && *s.sc.Listener != nil {
		(*s.sc.Listener).Close()
	}
	if s.sc.UDPListener != nil {
		(*s.sc.UDPListener).Close()
	}
}
func (s *UDP) Start(args interface{}, log *logger.Logger) (err error) {
	s.log = log
	s.cfg = args.(UDPArgs)
	if err = s.CheckArgs(); err != nil {
		return
	}
	s.log.Printf("use %s parent %s", *s.cfg.ParentType, *s.cfg.Parent)
	if err = s.InitService(); err != nil {
		return
	}
	host, port, _ := net.SplitHostPort(*s.cfg.Local)
	p, _ := strconv.Atoi(port)
	sc := server.NewServerChannel(host, p, s.log)
	s.sc = &sc
	err = sc.ListenUDP(s.callback)
	if err != nil {
		return
	}
	s.log.Printf("udp proxy on %s", (*sc.UDPListener).LocalAddr())
	return
}

func (s *UDP) Clean() {
	s.StopService()
}
func (s *UDP) callback(listener *net.UDPConn, packet []byte, localAddr, srcAddr *net.UDPAddr) {
	defer func() {
		if err := recover(); err != nil {
			s.log.Printf("udp conn handler crashed with err : %s \nstack: %s", err, string(debug.Stack()))
		}
	}()
	switch *s.cfg.ParentType {
	case "tcp", "tls":
		s.OutToTCP(packet, localAddr, srcAddr)
	case "udp":
		s.OutToUDP(packet, localAddr, srcAddr)
	default:
		s.log.Printf("unkown parent type %s", *s.cfg.ParentType)
	}
}
func (s *UDP) GetConn(connKey string) (conn net.Conn, isNew bool, err error) {
	isNew = !s.p.Has(connKey)
	var _conn interface{}
	if isNew {
		_conn, err = s.GetParentConn()
		if err != nil {
			return nil, false, err
		}
		s.p.Set(connKey, _conn)
	} else {
		_conn, _ = s.p.Get(connKey)
	}
	conn = _conn.(net.Conn)
	return
}
func (s *UDP) OutToTCP(data []byte, localAddr, srcAddr *net.UDPAddr) (err error) {
	s.UDPSend(data, localAddr, srcAddr)
	return
}
func (s *UDP) OutToUDPGCDeamon() {
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
			s.outUDPConnCtxMap.IterCb(func(key string, v interface{}) {
				if time.Now().Unix()-v.(*outUDPConnCtx).touchtime > gctime {
					(*(v.(*outUDPConnCtx).udpconn)).Close()
					gcKeys = append(gcKeys, key)
					s.log.Printf("gc udp conn %s <--> %s", (*v.(*outUDPConnCtx)).srcAddr, (*v.(*outUDPConnCtx)).localAddr)
				}
			})
			for _, k := range gcKeys {
				s.outUDPConnCtxMap.Remove(k)
			}
			gcKeys = nil
		}
	}()
}
func (s *UDP) OutToUDP(packet []byte, localAddr, srcAddr *net.UDPAddr) {
	var ouc *outUDPConnCtx
	if v, ok := s.outUDPConnCtxMap.Get(srcAddr.String()); !ok {
		clientSrcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
		conn, err := net.DialUDP("udp", clientSrcAddr, s.dstAddr)
		if err != nil {
			s.log.Printf("connect to udp %s fail,ERR:%s", s.dstAddr.String(), err)

		}
		ouc = &outUDPConnCtx{
			localAddr: localAddr,
			srcAddr:   srcAddr,
			udpconn:   conn,
		}
		s.outUDPConnCtxMap.Set(srcAddr.String(), ouc)
		go func() {
			defer func() {
				if e := recover(); e != nil {
					fmt.Printf("crashed, err: %s\nstack:\n%s", e, string(debug.Stack()))
				}
			}()
			s.log.Printf("udp conn %s <--> %s connected", srcAddr.String(), localAddr.String())
			buf := utils.LeakyBuffer.Get()
			defer func() {
				utils.LeakyBuffer.Put(buf)
				s.outUDPConnCtxMap.Remove(srcAddr.String())
				s.log.Printf("udp conn %s <--> %s released", srcAddr.String(), localAddr.String())
			}()
			for {
				n, err := ouc.udpconn.Read(buf)
				if err != nil {
					if !utils.IsNetClosedErr(err) {
						s.log.Printf("udp conn read udp packet fail , err: %s ", err)
					}
					return
				}
				ouc.touchtime = time.Now().Unix()
				go func() {
					defer func() {
						if e := recover(); e != nil {
							fmt.Printf("crashed, err: %s\nstack:\n%s", e, string(debug.Stack()))
						}
					}()
					(*(s.sc).UDPListener).SetWriteDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
					_, err = (*(s.sc).UDPListener).WriteTo(buf[:n], srcAddr)
					(*(s.sc).UDPListener).SetWriteDeadline(time.Time{})
				}()
			}
		}()
	} else {
		ouc = v.(*outUDPConnCtx)
	}
	go func() {
		ouc.touchtime = time.Now().Unix()
		ouc.udpconn.SetWriteDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
		ouc.udpconn.Write(packet)
		ouc.udpconn.SetWriteDeadline(time.Time{})
	}()
	return
}
func (s *UDP) GetParentConn() (conn net.Conn, err error) {
	if *s.cfg.ParentType == "tls" {
		var _conn tls.Conn
		_conn, err = utils.TlsConnectHost(*s.cfg.Parent, *s.cfg.Timeout, s.cfg.CertBytes, s.cfg.KeyBytes, nil)
		if err == nil {
			conn = net.Conn(&_conn)
		}
	} else {
		conn, err = utils.ConnectHost(*s.cfg.Parent, *s.cfg.Timeout)
	}
	return
}
func (s *UDP) UDPGCDeamon() {
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
				if time.Now().Unix()-v.(*UDPConnItem).touchtime > gctime {
					(*(v.(*UDPConnItem).conn)).Close()
					gcKeys = append(gcKeys, key)
					s.log.Printf("gc udp conn %s", v.(*UDPConnItem).connid)
				}
			})
			for _, k := range gcKeys {
				s.udpConns.Remove(k)
			}
			gcKeys = nil
		}
	}()
}
func (s *UDP) UDPSend(data []byte, localAddr, srcAddr *net.UDPAddr) {
	var (
		uc      *UDPConnItem
		key     = srcAddr.String()
		err     error
		outconn net.Conn
	)
	v, ok := s.udpConns.Get(key)
	if !ok {
		for {
			outconn, err = s.GetParentConn()
			if err != nil && strings.Contains(err.Error(), "can not connect at same time") {
				time.Sleep(time.Millisecond * 500)
				continue
			} else {
				break
			}
		}
		if err != nil {
			s.log.Printf("connect to %s fail, err: %s", *s.cfg.Parent, err)
			return
		}
		uc = &UDPConnItem{
			conn:      &outconn,
			srcAddr:   srcAddr,
			localAddr: localAddr,
		}
		s.udpConns.Set(key, uc)
		s.UDPRevecive(key)
	} else {
		uc = v.(*UDPConnItem)
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
		_, err = (*uc.conn).Write(utils.UDPPacket(fmt.Sprintf("%s", srcAddr.String()), data))
		(*uc.conn).SetWriteDeadline(time.Time{})
		if err != nil {
			s.log.Printf("write udp packet to %s fail ,flush err:%s ", *s.cfg.Parent, err)
		}
	}()
}
func (s *UDP) UDPRevecive(key string) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("crashed, err: %s\nstack:\n%s", e, string(debug.Stack()))
			}
		}()
		s.log.Printf("udp conn %s connected", key)
		var uc *UDPConnItem
		defer func() {
			if uc != nil {
				(*uc.conn).Close()
			}
			s.udpConns.Remove(key)
			s.log.Printf("udp conn %s released", key)
		}()
		v, ok := s.udpConns.Get(key)
		if !ok {
			s.log.Printf("[warn] udp conn not exists for %s", key)
			return
		}
		uc = v.(*UDPConnItem)
		for {
			_, body, err := utils.ReadUDPPacket(*uc.conn)
			if err != nil {
				if strings.Contains(err.Error(), "n != int(") {
					continue
				}
				if err != io.EOF && !utils.IsNetClosedErr(err) {
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

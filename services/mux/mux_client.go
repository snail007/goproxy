package mux

import (
	"crypto/tls"
	"fmt"
	"io"
	logger "log"
	"net"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/golang/snappy"
	clienttransport "github.com/snail007/goproxy/core/cs/client"
	"github.com/snail007/goproxy/core/lib/kcpcfg"
	encryptconn "github.com/snail007/goproxy/core/lib/transport/encrypt"
	"github.com/snail007/goproxy/services"
	"github.com/snail007/goproxy/utils"
	"github.com/snail007/goproxy/utils/jumper"
	"github.com/snail007/goproxy/utils/mapx"
	//"github.com/xtaci/smux"
	smux "github.com/hashicorp/yamux"
)

type MuxClientArgs struct {
	Parent       *string
	ParentType   *string
	CertFile     *string
	KeyFile      *string
	CertBytes    []byte
	KeyBytes     []byte
	Key          *string
	Timeout      *int
	IsCompress   *bool
	SessionCount *int
	KCP          kcpcfg.KCPConfigArgs
	Jumper       *string
	TCPSMethod   *string
	TCPSPassword *string
	TOUMethod    *string
	TOUPassword  *string
}
type ClientUDPConnItem struct {
	conn      *smux.Stream
	isActive  bool
	touchtime int64
	srcAddr   *net.UDPAddr
	localAddr *net.UDPAddr
	udpConn   *net.UDPConn
	connid    string
}
type MuxClient struct {
	cfg      MuxClientArgs
	isStop   bool
	sessions mapx.ConcurrentMap
	log      *logger.Logger
	jumper   *jumper.Jumper
	udpConns mapx.ConcurrentMap
}

func NewMuxClient() services.Service {
	return &MuxClient{
		cfg:      MuxClientArgs{},
		isStop:   false,
		sessions: mapx.NewConcurrentMap(),
		udpConns: mapx.NewConcurrentMap(),
	}
}

func (s *MuxClient) InitService() (err error) {
	s.UDPGCDeamon()
	return
}

func (s *MuxClient) CheckArgs() (err error) {
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
	if *s.cfg.ParentType == "tls" {
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
func (s *MuxClient) StopService() {
	defer func() {
		e := recover()
		if e != nil {
			s.log.Printf("stop client service crashed,%s", e)
		} else {
			s.log.Printf("service client stopped")
		}
		s.cfg = MuxClientArgs{}
		s.jumper = nil
		s.log = nil
		s.sessions = nil
		s.udpConns = nil
		s = nil
	}()
	s.isStop = true
	for _, sess := range s.sessions.Items() {
		sess.(*smux.Session).Close()
	}
}
func (s *MuxClient) Start(args interface{}, log *logger.Logger) (err error) {
	s.log = log
	s.cfg = args.(MuxClientArgs)
	if err = s.CheckArgs(); err != nil {
		return
	}
	if err = s.InitService(); err != nil {
		return
	}
	s.log.Printf("client started")
	count := 1
	if *s.cfg.SessionCount > 0 {
		count = *s.cfg.SessionCount
	}
	for i := 1; i <= count; i++ {
		key := fmt.Sprintf("worker[%d]", i)
		s.log.Printf("session %s started", key)
		go func(i int) {
			defer func() {
				e := recover()
				if e != nil {
					s.log.Printf("session worker crashed: %s\nstack:%s", e, string(debug.Stack()))
				}
			}()
			for {
				if s.isStop {
					return
				}
				conn, err := s.getParentConn()
				if err != nil {
					s.log.Printf("connection err: %s, retrying...", err)
					time.Sleep(time.Second * 3)
					continue
				}
				conn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
				g := sync.WaitGroup{}
				g.Add(1)
				go func() {
					defer func() {
						_ = recover()
						g.Done()
					}()
					_, err = conn.Write(utils.BuildPacket(CONN_CLIENT, fmt.Sprintf("%s-%d", *s.cfg.Key, i)))
				}()
				g.Wait()
				conn.SetDeadline(time.Time{})
				if err != nil {
					conn.Close()
					s.log.Printf("connection err: %s, retrying...", err)
					time.Sleep(time.Second * 3)
					continue
				}
				session, err := smux.Server(conn, nil)
				if err != nil {
					s.log.Printf("session err: %s, retrying...", err)
					conn.Close()
					time.Sleep(time.Second * 3)
					continue
				}
				if _sess, ok := s.sessions.Get(key); ok {
					_sess.(*smux.Session).Close()
				}
				s.sessions.Set(key, session)
				for {
					if s.isStop {
						return
					}
					stream, err := session.AcceptStream()
					if err != nil {
						s.log.Printf("accept stream err: %s, retrying...", err)
						session.Close()
						time.Sleep(time.Second * 3)
						break
					}
					go func() {
						defer func() {
							e := recover()
							if e != nil {
								s.log.Printf("stream handler crashed: %s", e)
							}
						}()
						var ID, clientLocalAddr, serverID string
						stream.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
						err = utils.ReadPacketData(stream, &ID, &clientLocalAddr, &serverID)
						stream.SetDeadline(time.Time{})
						if err != nil {
							s.log.Printf("read stream signal err: %s", err)
							stream.Close()
							return
						}
						//s.log.Printf("worker[%d] signal revecived,server %s stream %s %s", i, serverID, ID, clientLocalAddr)
						protocol := clientLocalAddr[:3]
						localAddr := clientLocalAddr[4:]
						if protocol == "udp" {
							s.ServeUDP(stream, localAddr, ID)
						} else {
							s.ServeConn(stream, localAddr, ID)
						}
					}()
				}
			}
		}(i)
	}
	return
}
func (s *MuxClient) Clean() {
	s.StopService()
}
func (s *MuxClient) getParentConn() (conn net.Conn, err error) {
	if *s.cfg.ParentType == "tls" {
		if s.jumper == nil {
			var _conn tls.Conn
			_conn, err = clienttransport.TlsConnectHost(*s.cfg.Parent, *s.cfg.Timeout, s.cfg.CertBytes, s.cfg.KeyBytes, nil)
			if err == nil {
				conn = net.Conn(&_conn)
			}
		} else {
			conf, e := utils.TlsConfig(s.cfg.CertBytes, s.cfg.KeyBytes, nil)
			if e != nil {
				return nil, e
			}
			var _c net.Conn
			_c, err = s.jumper.Dial(*s.cfg.Parent, time.Millisecond*time.Duration(*s.cfg.Timeout))
			if err == nil {
				conn = net.Conn(tls.Client(_c, conf))
			}
		}

	} else if *s.cfg.ParentType == "kcp" {
		conn, err = clienttransport.KCPConnectHost(*s.cfg.Parent, s.cfg.KCP)
	} else if *s.cfg.ParentType == "tcps" {
		if s.jumper == nil {
			conn, err = clienttransport.TCPSConnectHost(*s.cfg.Parent, *s.cfg.TCPSMethod, *s.cfg.TCPSPassword, false, *s.cfg.Timeout)
		} else {
			conn, err = s.jumper.Dial(*s.cfg.Parent, time.Millisecond*time.Duration(*s.cfg.Timeout))
			if err == nil {
				conn, err = encryptconn.NewConn(conn, *s.cfg.TCPSMethod, *s.cfg.TCPSPassword)
			}
		}

	} else if *s.cfg.ParentType == "tou" {
		conn, err = clienttransport.TOUConnectHost(*s.cfg.Parent, *s.cfg.TCPSMethod, *s.cfg.TCPSPassword, false, *s.cfg.Timeout)
	} else {
		if s.jumper == nil {
			conn, err = clienttransport.TCPConnectHost(*s.cfg.Parent, *s.cfg.Timeout)
		} else {
			conn, err = s.jumper.Dial(*s.cfg.Parent, time.Millisecond*time.Duration(*s.cfg.Timeout))
		}
	}
	return
}
func (s *MuxClient) ServeUDP(inConn *smux.Stream, localAddr, ID string) {
	var item *ClientUDPConnItem
	var body []byte
	var err error
	srcAddr := ""
	defer func() {
		if item != nil {
			(*item).conn.Close()
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
				conn:      inConn,
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
			defer func() { _ = recover() }()
			(*item).udpConn.Write(body)
		}()
	}
}
func (s *MuxClient) UDPRevecive(key, ID string) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("crashed, err: %s\nstack:%s", e, string(debug.Stack()))
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
			cui.conn.Close()
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
						fmt.Printf("crashed, err: %s\nstack:%s", e, string(debug.Stack()))
					}
				}()
				cui.conn.SetWriteDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
				_, err = cui.conn.Write(utils.UDPPacket(cui.srcAddr.String(), buf[:n]))
				cui.conn.SetWriteDeadline(time.Time{})
				if err != nil {
					cui.udpConn.Close()
					return
				}
			}()
		}
	}()
}
func (s *MuxClient) UDPGCDeamon() {
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
func (s *MuxClient) ServeConn(inConn *smux.Stream, localAddr, ID string) {
	var err error
	var outConn net.Conn
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
		inConn.Close()
		utils.CloseConn(&outConn)
		s.log.Printf("build connection error, err: %s", err)
		return
	}

	s.log.Printf("stream %s created", ID)
	if *s.cfg.IsCompress {
		die1 := make(chan bool, 1)
		die2 := make(chan bool, 1)
		go func() {
			defer func() {
				if e := recover(); e != nil {
					fmt.Printf("crashed, err: %s\nstack:%s", e, string(debug.Stack()))
				}
			}()
			io.Copy(outConn, snappy.NewReader(inConn))
			die1 <- true
		}()
		go func() {
			defer func() {
				if e := recover(); e != nil {
					fmt.Printf("crashed, err: %s\nstack:%s", e, string(debug.Stack()))
				}
			}()
			io.Copy(snappy.NewWriter(inConn), outConn)
			die2 <- true
		}()
		select {
		case <-die1:
		case <-die2:
		}
		outConn.Close()
		inConn.Close()
		s.log.Printf("%s stream %s released", *s.cfg.Key, ID)
	} else {
		utils.IoBind(inConn, outConn, func(err interface{}) {
			s.log.Printf("stream %s released", ID)
		}, s.log)
	}
}

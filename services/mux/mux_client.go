package mux

import (
	"crypto/tls"
	"fmt"
	"io"
	logger "log"
	"net"
	"time"

	"github.com/golang/snappy"
	"github.com/snail007/goproxy/services"
	"github.com/snail007/goproxy/services/kcpcfg"
	"github.com/snail007/goproxy/utils"
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
}
type MuxClient struct {
	cfg      MuxClientArgs
	isStop   bool
	sessions utils.ConcurrentMap
	log      *logger.Logger
}

func NewMuxClient() services.Service {
	return &MuxClient{
		cfg:      MuxClientArgs{},
		isStop:   false,
		sessions: utils.NewConcurrentMap(),
	}
}

func (s *MuxClient) InitService() (err error) {
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
	return
}
func (s *MuxClient) StopService() {
	defer func() {
		e := recover()
		if e != nil {
			s.log.Printf("stop client service crashed,%s", e)
		} else {
			s.log.Printf("service client stoped")
		}
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
					s.log.Printf("session worker crashed: %s", e)
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
				_, err = conn.Write(utils.BuildPacket(CONN_CLIENT, fmt.Sprintf("%s-%d", *s.cfg.Key, i)))
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
						s.log.Printf("worker[%d] signal revecived,server %s stream %s %s", i, serverID, ID, clientLocalAddr)
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
func (s *MuxClient) ServeUDP(inConn *smux.Stream, localAddr, ID string) {

	for {
		if s.isStop {
			return
		}
		inConn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
		srcAddr, body, err := utils.ReadUDPPacket(inConn)
		inConn.SetDeadline(time.Time{})
		if err != nil {
			s.log.Printf("udp packet revecived fail, err: %s", err)
			s.log.Printf("connection %s released", ID)
			inConn.Close()
			break
		} else {
			//s.log.Printf("udp packet revecived:%s,%v", srcAddr, body)
			go func() {
				defer func() {
					if e := recover(); e != nil {
						s.log.Printf("client processUDPPacket crashed,err: %s", e)
					}
				}()
				s.processUDPPacket(inConn, srcAddr, localAddr, body)
			}()

		}

	}
	// }
}
func (s *MuxClient) processUDPPacket(inConn *smux.Stream, srcAddr, localAddr string, body []byte) {
	dstAddr, err := net.ResolveUDPAddr("udp", localAddr)
	if err != nil {
		s.log.Printf("can't resolve address: %s", err)
		inConn.Close()
		return
	}
	clientSrcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	conn, err := net.DialUDP("udp", clientSrcAddr, dstAddr)
	if err != nil {
		s.log.Printf("connect to udp %s fail,ERR:%s", dstAddr.String(), err)
		return
	}
	conn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
	_, err = conn.Write(body)
	conn.SetDeadline(time.Time{})
	if err != nil {
		s.log.Printf("send udp packet to %s fail,ERR:%s", dstAddr.String(), err)
		return
	}
	//s.log.Printf("send udp packet to %s success", dstAddr.String())
	buf := make([]byte, 1024)
	conn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
	length, _, err := conn.ReadFromUDP(buf)
	conn.SetDeadline(time.Time{})
	if err != nil {
		s.log.Printf("read udp response from %s fail ,ERR:%s", dstAddr.String(), err)
		return
	}
	respBody := buf[0:length]
	//s.log.Printf("revecived udp packet from %s , %v", dstAddr.String(), respBody)
	bs := utils.UDPPacket(srcAddr, respBody)
	(*inConn).SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
	_, err = (*inConn).Write(bs)
	(*inConn).SetDeadline(time.Time{})
	if err != nil {
		s.log.Printf("send udp response fail ,ERR:%s", err)
		inConn.Close()
		return
	}
	//s.log.Printf("send udp response success ,from:%s ,%d ,%v", dstAddr.String(), len(bs), bs)
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
			io.Copy(outConn, snappy.NewReader(inConn))
			die1 <- true
		}()
		go func() {
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

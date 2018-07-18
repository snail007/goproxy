package tunnel

import (
	"crypto/tls"
	"fmt"
	"io"
	logger "log"
	"net"
	"os"
	"time"

	"github.com/snail007/goproxy/services"
	"github.com/snail007/goproxy/utils"
	//"github.com/xtaci/smux"
	smux "github.com/hashicorp/yamux"
)

const (
	CONN_SERVER_MUX = uint8(6)
	CONN_CLIENT_MUX = uint8(7)
)

type TunnelClientArgs struct {
	Parent    *string
	CertFile  *string
	KeyFile   *string
	CertBytes []byte
	KeyBytes  []byte
	Key       *string
	Timeout   *int
	Mux       *bool
}
type TunnelClient struct {
	cfg       TunnelClientArgs
	ctrlConn  net.Conn
	isStop    bool
	userConns utils.ConcurrentMap
	log       *logger.Logger
}

func NewTunnelClient() services.Service {
	return &TunnelClient{
		cfg:       TunnelClientArgs{},
		userConns: utils.NewConcurrentMap(),
		isStop:    false,
	}
}

func (s *TunnelClient) InitService() (err error) {
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
	}()
	s.isStop = true
	if s.ctrlConn != nil {
		s.ctrlConn.Close()
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
			s.ctrlConn.Close()
		}

		s.ctrlConn, err = s.GetInConn(CONN_CLIENT_CONTROL, *s.cfg.Key)
		if err != nil {
			s.log.Printf("control connection err: %s, retrying...", err)
			time.Sleep(time.Second * 3)
			if s.ctrlConn != nil {
				s.ctrlConn.Close()
			}
			continue
		}
		for {
			if s.isStop {
				return
			}
			var ID, clientLocalAddr, serverID string
			err = utils.ReadPacketData(s.ctrlConn, &ID, &clientLocalAddr, &serverID)
			if err != nil {
				if s.ctrlConn != nil {
					s.ctrlConn.Close()
				}
				s.log.Printf("read connection signal err: %s, retrying...", err)
				break
			}
			s.log.Printf("signal revecived:%s %s %s", serverID, ID, clientLocalAddr)
			protocol := clientLocalAddr[:3]
			localAddr := clientLocalAddr[4:]
			if protocol == "udp" {
				go s.ServeUDP(localAddr, ID, serverID)
			} else {
				go s.ServeConn(localAddr, ID, serverID)
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
	var _conn tls.Conn
	_conn, err = utils.TlsConnectHost(*s.cfg.Parent, *s.cfg.Timeout, s.cfg.CertBytes, s.cfg.KeyBytes, nil)
	if err == nil {
		conn = net.Conn(&_conn)
		c, e := smux.Client(conn, &smux.Config{
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
		conn, e = c.OpenStream()
		if e != nil {
			s.log.Printf("mux client conn open stream error,ERR:%s", e)
			err = e
			return
		}
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

	for {
		if s.isStop {
			return
		}
		srcAddr, body, err := utils.ReadUDPPacket(inConn)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			s.log.Printf("connection %s released", ID)
			utils.CloseConn(&inConn)
			break
		} else if err != nil {
			s.log.Printf("udp packet revecived fail, err: %s", err)
		} else {
			//s.log.Printf("udp packet revecived:%s,%v", srcAddr, body)
			go s.processUDPPacket(&inConn, srcAddr, localAddr, body)
		}

	}
	// }
}
func (s *TunnelClient) processUDPPacket(inConn *net.Conn, srcAddr, localAddr string, body []byte) {
	dstAddr, err := net.ResolveUDPAddr("udp", localAddr)
	if err != nil {
		s.log.Printf("can't resolve address: %s", err)
		utils.CloseConn(inConn)
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
	if err != nil {
		s.log.Printf("send udp packet to %s fail,ERR:%s", dstAddr.String(), err)
		return
	}
	//s.log.Printf("send udp packet to %s success", dstAddr.String())
	buf := make([]byte, 1024)
	length, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		s.log.Printf("read udp response from %s fail ,ERR:%s", dstAddr.String(), err)
		return
	}
	respBody := buf[0:length]
	//s.log.Printf("revecived udp packet from %s , %v", dstAddr.String(), respBody)
	bs := utils.UDPPacket(srcAddr, respBody)
	_, err = (*inConn).Write(bs)
	if err != nil {
		s.log.Printf("send udp response fail ,ERR:%s", err)
		utils.CloseConn(inConn)
		return
	}
	//s.log.Printf("send udp response success ,from:%s ,%d ,%v", dstAddr.String(), len(bs), bs)
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

package udp

import (
	"bufio"
	"fmt"
	"hash/crc32"
	"io"
	logger "log"
	"net"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/snail007/goproxy/services"
	"github.com/snail007/goproxy/services/kcpcfg"
	"github.com/snail007/goproxy/utils"
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
	p       utils.ConcurrentMap
	outPool utils.OutConn
	cfg     UDPArgs
	sc      *utils.ServerChannel
	isStop  bool
	log     *logger.Logger
}

func NewUDP() services.Service {
	return &UDP{
		outPool: utils.OutConn{},
		p:       utils.NewConcurrentMap(),
		isStop:  false,
	}
}
func (s *UDP) CheckArgs() (err error) {
	if *s.cfg.Parent == "" {
		err = fmt.Errorf("parent required for udp %s", *s.cfg.Local)
		return
	}
	if *s.cfg.ParentType == "" {
		err = fmt.Errorf("parent type unkown,use -T <tls|tcp>")
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
func (s *UDP) InitService() (err error) {
	if *s.cfg.ParentType != "udp" {
		s.InitOutConnPool()
	}
	return
}
func (s *UDP) StopService() {
	defer func() {
		e := recover()
		if e != nil {
			s.log.Printf("stop udp service crashed,%s", e)
		} else {
			s.log.Printf("service udp stoped")
		}
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
	sc := utils.NewServerChannel(host, p, s.log)
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
func (s *UDP) callback(packet []byte, localAddr, srcAddr *net.UDPAddr) {
	defer func() {
		if err := recover(); err != nil {
			s.log.Printf("udp conn handler crashed with err : %s \nstack: %s", err, string(debug.Stack()))
		}
	}()
	var err error
	switch *s.cfg.ParentType {
	case "tcp":
		fallthrough
	case "tls":
		err = s.OutToTCP(packet, localAddr, srcAddr)
	case "udp":
		err = s.OutToUDP(packet, localAddr, srcAddr)
	default:
		err = fmt.Errorf("unkown parent type %s", *s.cfg.ParentType)
	}
	if err != nil {
		s.log.Printf("connect to %s parent %s fail, ERR:%s", *s.cfg.ParentType, *s.cfg.Parent, err)
	}
}
func (s *UDP) GetConn(connKey string) (conn net.Conn, isNew bool, err error) {
	isNew = !s.p.Has(connKey)
	var _conn interface{}
	if isNew {
		_conn, err = s.outPool.Get()
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
func (s *UDP) OutToTCP(packet []byte, localAddr, srcAddr *net.UDPAddr) (err error) {
	numLocal := crc32.ChecksumIEEE([]byte(localAddr.String()))
	numSrc := crc32.ChecksumIEEE([]byte(srcAddr.String()))
	mod := uint32(10)
	if mod == 0 {
		mod = 10
	}
	connKey := uint64((numLocal/10)*10 + numSrc%mod)
	conn, isNew, err := s.GetConn(fmt.Sprintf("%d", connKey))
	if err != nil {
		s.log.Printf("upd get conn to %s parent %s fail, ERR:%s", *s.cfg.ParentType, *s.cfg.Parent, err)
		return
	}
	if isNew {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					s.log.Printf("udp conn handler out to tcp crashed with err : %s \nstack: %s", err, string(debug.Stack()))
				}
			}()
			s.log.Printf("conn %d created , local: %s", connKey, srcAddr.String())
			for {
				if s.isStop {
					conn.Close()
					return
				}
				srcAddrFromConn, body, err := utils.ReadUDPPacket(bufio.NewReader(conn))
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					//s.log.Printf("connection %d released", connKey)
					s.p.Remove(fmt.Sprintf("%d", connKey))
					break
				}
				if err != nil {
					s.log.Printf("parse revecived udp packet fail, err: %s", err)
					continue
				}
				//s.log.Printf("udp packet revecived over parent , local:%s", srcAddrFromConn)
				_srcAddr := strings.Split(srcAddrFromConn, ":")
				if len(_srcAddr) != 2 {
					s.log.Printf("parse revecived udp packet fail, addr error : %s", srcAddrFromConn)
					continue
				}
				port, _ := strconv.Atoi(_srcAddr[1])
				dstAddr := &net.UDPAddr{IP: net.ParseIP(_srcAddr[0]), Port: port}
				_, err = s.sc.UDPListener.WriteToUDP(body, dstAddr)
				if err != nil {
					s.log.Printf("udp response to local %s fail,ERR:%s", srcAddr, err)
					continue
				}
				//s.log.Printf("udp response to local %s success", srcAddr)
			}
		}()
	}
	//s.log.Printf("select conn %d , local: %s", connKey, srcAddr.String())
	writer := bufio.NewWriter(conn)
	//fmt.Println(conn, writer)
	writer.Write(utils.UDPPacket(srcAddr.String(), packet))
	err = writer.Flush()
	if err != nil {
		s.log.Printf("write udp packet to %s fail ,flush err:%s", *s.cfg.Parent, err)
		return
	}
	//s.log.Printf("write packet %v", packet)
	return
}
func (s *UDP) OutToUDP(packet []byte, localAddr, srcAddr *net.UDPAddr) (err error) {
	//s.log.Printf("udp packet revecived:%s,%v", srcAddr, packet)
	dstAddr, err := net.ResolveUDPAddr("udp", *s.cfg.Parent)
	if err != nil {
		s.log.Printf("resolve udp addr %s fail  fail,ERR:%s", dstAddr.String(), err)
		return
	}
	clientSrcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	conn, err := net.DialUDP("udp", clientSrcAddr, dstAddr)
	if err != nil {
		s.log.Printf("connect to udp %s fail,ERR:%s", dstAddr.String(), err)
		return
	}
	conn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
	_, err = conn.Write(packet)
	if err != nil {
		s.log.Printf("send udp packet to %s fail,ERR:%s", dstAddr.String(), err)
		return
	}
	//s.log.Printf("send udp packet to %s success", dstAddr.String())
	buf := make([]byte, 512)
	len, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		s.log.Printf("read udp response from %s fail ,ERR:%s", dstAddr.String(), err)
		return
	}
	//s.log.Printf("revecived udp packet from %s , %v", dstAddr.String(), respBody)
	_, err = s.sc.UDPListener.WriteToUDP(buf[0:len], srcAddr)
	if err != nil {
		s.log.Printf("send udp response to cluster fail ,ERR:%s", err)
		return
	}
	//s.log.Printf("send udp response to cluster success ,from:%s", dstAddr.String())
	return
}
func (s *UDP) InitOutConnPool() {
	if *s.cfg.ParentType == "tls" || *s.cfg.ParentType == "tcp" {
		//dur int, isTLS bool, certBytes, keyBytes []byte,
		//parent string, timeout int, InitialCap int, MaxCap int
		s.outPool = utils.NewOutConn(
			*s.cfg.CheckParentInterval,
			*s.cfg.ParentType,
			kcpcfg.KCPConfigArgs{},
			s.cfg.CertBytes, s.cfg.KeyBytes, nil,
			*s.cfg.Parent,
			*s.cfg.Timeout,
		)
	}
}

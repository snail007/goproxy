package services

import (
	"fmt"
	"io"
	"log"
	"net"
	"proxy/utils"
	"runtime/debug"
	"time"

	"strconv"
)

type TCP struct {
	outPool utils.OutPool
	cfg     TCPArgs
}

func NewTCP() Service {
	return &TCP{
		outPool: utils.OutPool{},
		cfg:     TCPArgs{},
	}
}
func (s *TCP) InitService() {
	s.InitOutConnPool()
}
func (s *TCP) StopService() {
	if s.outPool.Pool != nil {
		s.outPool.Pool.ReleaseAll()
	}
}
func (s *TCP) Start(args interface{}) (err error) {
	s.cfg = args.(TCPArgs)
	if *s.cfg.Parent != "" {
		log.Printf("use %s parent %s", *s.cfg.ParentType, *s.cfg.Parent)
	} else {
		log.Fatalf("parent required for %s %s", s.cfg.Protocol(), *s.cfg.Local)
	}

	s.InitService()

	host, port, _ := net.SplitHostPort(*s.cfg.Local)
	p, _ := strconv.Atoi(port)
	sc := utils.NewServerChannel(host, p)
	if !*s.cfg.IsTLS {
		err = sc.ListenTCP(s.callback)
	} else {
		err = sc.ListenTls(s.cfg.CertBytes, s.cfg.KeyBytes, s.callback)
	}
	if err != nil {
		return
	}
	log.Printf("%s proxy on %s", s.cfg.Protocol(), (*sc.Listener).Addr())
	return
}

func (s *TCP) Clean() {
	s.StopService()
}
func (s *TCP) callback(inConn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%s conn handler crashed with err : %s \nstack: %s", s.cfg.Protocol(), err, string(debug.Stack()))
		}
	}()
	var err error
	switch *s.cfg.ParentType {
	case TYPE_TCP:
		fallthrough
	case TYPE_TLS:
		err = s.OutToTCP(&inConn)
	case TYPE_UDP:
		err = s.OutToUDP(&inConn)
	default:
		err = fmt.Errorf("unkown parent type %s", *s.cfg.ParentType)
	}
	if err != nil {
		log.Printf("connect to %s parent %s fail, ERR:%s", *s.cfg.ParentType, *s.cfg.Parent, err)
		utils.CloseConn(&inConn)
	}
}
func (s *TCP) OutToTCP(inConn *net.Conn) (err error) {
	var outConn net.Conn
	var _outConn interface{}
	_outConn, err = s.outPool.Pool.Get()
	if err == nil {
		outConn = _outConn.(net.Conn)
	}
	if err != nil {
		log.Printf("connect to %s , err:%s", *s.cfg.Parent, err)
		utils.CloseConn(inConn)
		return
	}
	inAddr := (*inConn).RemoteAddr().String()
	inLocalAddr := (*inConn).LocalAddr().String()
	outAddr := outConn.RemoteAddr().String()
	outLocalAddr := outConn.LocalAddr().String()
	utils.IoBind((*inConn), outConn, func(isSrcErr bool, err error) {
		log.Printf("conn %s - %s - %s -%s released", inAddr, inLocalAddr, outLocalAddr, outAddr)
		utils.CloseConn(inConn)
		utils.CloseConn(&outConn)
	}, func(n int, d bool) {}, 0)
	log.Printf("conn %s - %s - %s -%s connected", inAddr, inLocalAddr, outLocalAddr, outAddr)
	return
}
func (s *TCP) OutToUDP(inConn *net.Conn) (err error) {
	log.Printf("conn created , remote : %s ", (*inConn).RemoteAddr())
	for {
		srcAddr, body, err := utils.ReadUDPPacket(inConn)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			//log.Printf("connection %s released", srcAddr)
			utils.CloseConn(inConn)
			break
		}
		//log.Debugf("udp packet revecived:%s,%v", srcAddr, body)
		dstAddr, err := net.ResolveUDPAddr("udp", *s.cfg.Parent)
		if err != nil {
			log.Printf("can't resolve address: %s", err)
			utils.CloseConn(inConn)
			break
		}
		clientSrcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
		conn, err := net.DialUDP("udp", clientSrcAddr, dstAddr)
		if err != nil {
			log.Printf("connect to udp %s fail,ERR:%s", dstAddr.String(), err)
			continue
		}
		conn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
		_, err = conn.Write(body)
		if err != nil {
			log.Printf("send udp packet to %s fail,ERR:%s", dstAddr.String(), err)
			continue
		}
		//log.Debugf("send udp packet to %s success", dstAddr.String())
		buf := make([]byte, 512)
		len, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("read udp response from %s fail ,ERR:%s", dstAddr.String(), err)
			continue
		}
		respBody := buf[0:len]
		//log.Debugf("revecived udp packet from %s , %v", dstAddr.String(), respBody)
		_, err = (*inConn).Write(utils.UDPPacket(srcAddr, respBody))
		if err != nil {
			log.Printf("send udp response fail ,ERR:%s", err)
			utils.CloseConn(inConn)
			break
		}
		//log.Printf("send udp response success ,from:%s", dstAddr.String())
	}
	return

}
func (s *TCP) InitOutConnPool() {
	if *s.cfg.ParentType == TYPE_TLS || *s.cfg.ParentType == TYPE_TCP {
		//dur int, isTLS bool, certBytes, keyBytes []byte,
		//parent string, timeout int, InitialCap int, MaxCap int
		s.outPool = utils.NewOutPool(
			*s.cfg.CheckParentInterval,
			*s.cfg.ParentType == TYPE_TLS,
			s.cfg.CertBytes, s.cfg.KeyBytes,
			*s.cfg.Parent,
			*s.cfg.Timeout,
			*s.cfg.PoolSize,
			*s.cfg.PoolSize*2,
		)
	}
}

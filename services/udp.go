package services

import (
	"bufio"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net"
	"proxy/utils"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

type UDP struct {
	p       utils.ConcurrentMap
	outPool utils.OutPool
	cfg     UDPArgs
	sc      *utils.ServerChannel
}

func NewUDP() Service {
	return &UDP{
		outPool: utils.OutPool{},
		p:       utils.NewConcurrentMap(),
	}
}
func (s *UDP) InitService() {
	if *s.cfg.ParentType != TYPE_UDP {
		s.InitOutConnPool()
	}
}
func (s *UDP) StopService() {
	if s.outPool.Pool != nil {
		s.outPool.Pool.ReleaseAll()
	}
}
func (s *UDP) Start(args interface{}) (err error) {
	s.cfg = args.(UDPArgs)
	if *s.cfg.Parent != "" {
		log.Printf("use %s parent %s", *s.cfg.ParentType, *s.cfg.Parent)
	} else {
		log.Fatalf("parent required for udp %s", *s.cfg.Local)
	}

	s.InitService()

	host, port, _ := net.SplitHostPort(*s.cfg.Local)
	p, _ := strconv.Atoi(port)
	sc := utils.NewServerChannel(host, p)
	s.sc = &sc
	err = sc.ListenUDP(s.callback)
	if err != nil {
		return
	}
	log.Printf("udp proxy on %s", (*sc.UDPListener).LocalAddr())
	return
}

func (s *UDP) Clean() {
	s.StopService()
}
func (s *UDP) callback(packet []byte, localAddr, srcAddr *net.UDPAddr) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("udp conn handler crashed with err : %s \nstack: %s", err, string(debug.Stack()))
		}
	}()
	var err error
	switch *s.cfg.ParentType {
	case TYPE_TCP:
		fallthrough
	case TYPE_TLS:
		err = s.OutToTCP(packet, localAddr, srcAddr)
	case TYPE_UDP:
		err = s.OutToUDP(packet, localAddr, srcAddr)
	default:
		err = fmt.Errorf("unkown parent type %s", *s.cfg.ParentType)
	}
	if err != nil {
		log.Printf("connect to %s parent %s fail, ERR:%s", *s.cfg.ParentType, *s.cfg.Parent, err)
	}
}
func (s *UDP) GetConn(connKey string) (conn net.Conn, isNew bool, err error) {
	isNew = !s.p.Has(connKey)
	var _conn interface{}
	if isNew {
		_conn, err = s.outPool.Pool.Get()
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
	mod := uint32(*s.cfg.PoolSize)
	if mod == 0 {
		mod = 10
	}
	connKey := uint64((numLocal/10)*10 + numSrc%mod)
	conn, isNew, err := s.GetConn(fmt.Sprintf("%d", connKey))
	if err != nil {
		log.Printf("upd get conn to %s parent %s fail, ERR:%s", *s.cfg.ParentType, *s.cfg.Parent, err)
		return
	}
	if isNew {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("udp conn handler out to tcp crashed with err : %s \nstack: %s", err, string(debug.Stack()))
				}
			}()
			log.Printf("conn %d created , local: %s", connKey, srcAddr.String())
			for {
				srcAddrFromConn, body, err := utils.ReadUDPPacket(&conn)
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					//log.Printf("connection %d released", connKey)
					s.p.Remove(fmt.Sprintf("%d", connKey))
					break
				}
				if err != nil {
					log.Printf("parse revecived udp packet fail, err: %s", err)
					continue
				}
				//log.Printf("udp packet revecived over parent , local:%s", srcAddrFromConn)
				_srcAddr := strings.Split(srcAddrFromConn, ":")
				if len(_srcAddr) != 2 {
					log.Printf("parse revecived udp packet fail, addr error : %s", srcAddrFromConn)
					continue
				}
				port, _ := strconv.Atoi(_srcAddr[1])
				dstAddr := &net.UDPAddr{IP: net.ParseIP(_srcAddr[0]), Port: port}
				_, err = s.sc.UDPListener.WriteToUDP(body, dstAddr)
				if err != nil {
					log.Printf("udp response to local %s fail,ERR:%s", srcAddr, err)
					continue
				}
				//log.Printf("udp response to local %s success", srcAddr)
			}
		}()
	}
	//log.Printf("select conn %d , local: %s", connKey, srcAddr.String())
	writer := bufio.NewWriter(conn)
	//fmt.Println(conn, writer)
	writer.Write(utils.UDPPacket(srcAddr.String(), packet))
	err = writer.Flush()
	if err != nil {
		log.Printf("write udp packet to %s fail ,flush err:%s", *s.cfg.Parent, err)
		return
	}
	//log.Printf("write packet %v", packet)
	return
}
func (s *UDP) OutToUDP(packet []byte, localAddr, srcAddr *net.UDPAddr) (err error) {
	//log.Printf("udp packet revecived:%s,%v", srcAddr, packet)
	dstAddr, err := net.ResolveUDPAddr("udp", *s.cfg.Parent)
	if err != nil {
		log.Printf("resolve udp addr %s fail  fail,ERR:%s", dstAddr.String(), err)
		return
	}
	clientSrcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	conn, err := net.DialUDP("udp", clientSrcAddr, dstAddr)
	if err != nil {
		log.Printf("connect to udp %s fail,ERR:%s", dstAddr.String(), err)
		return
	}
	conn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
	_, err = conn.Write(packet)
	if err != nil {
		log.Printf("send udp packet to %s fail,ERR:%s", dstAddr.String(), err)
		return
	}
	//log.Printf("send udp packet to %s success", dstAddr.String())
	buf := make([]byte, 512)
	len, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		log.Printf("read udp response from %s fail ,ERR:%s", dstAddr.String(), err)
		return
	}
	//log.Printf("revecived udp packet from %s , %v", dstAddr.String(), respBody)
	_, err = s.sc.UDPListener.WriteToUDP(buf[0:len], srcAddr)
	if err != nil {
		log.Printf("send udp response to cluster fail ,ERR:%s", err)
		return
	}
	//log.Printf("send udp response to cluster success ,from:%s", dstAddr.String())
	return
}
func (s *UDP) InitOutConnPool() {
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

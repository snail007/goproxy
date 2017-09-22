package services

import (
	"fmt"
	"log"
	"net"
	"proxy/utils"
	"runtime/debug"
	"strconv"
)

type TLS struct {
	outPool utils.OutPool
	cfg     TLSArgs
}

func NewTLS() Service {
	return &TLS{}
}
func (s *TLS) InitService() {
	s.InitOutConnPool()
}
func (s *TLS) StopService() {
	if s.outPool.Pool != nil {
		s.outPool.Pool.ReleaseAll()
	}
}
func (s *TLS) Start(args interface{}) (err error) {
	s.cfg = args.(TLSArgs)
	if *s.cfg.Parent != "" {
		log.Printf("use %s parent %s", *s.cfg.ParentType, *s.cfg.Parent)
	} else {
		log.Fatalf("parent required for tls")
	}

	s.InitService()

	host, port, _ := net.SplitHostPort(*s.cfg.Local)
	p, _ := strconv.Atoi(port)
	sc := utils.NewServerChannel(host, p)
	err = sc.ListenTls(s.cfg.CertBytes, s.cfg.KeyBytes, func(inConn net.Conn) {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("tls conn handler crashed with err : %s \nstack: %s", err, string(debug.Stack()))
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
		}()
	})
	if err != nil {
		return
	}
	log.Printf("tls proxy on %s", (*sc.Listener).Addr())
	return
}

func (s *TLS) Clean() {
	s.StopService()
}

func (s *TLS) OutToTCP(inConn *net.Conn) (err error) {
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
	utils.IoBind((*inConn), outConn, func(err error) {
		log.Printf("conn %s - %s - %s -%s released", inAddr, inLocalAddr, outLocalAddr, outAddr)
		utils.CloseConn(inConn)
		utils.CloseConn(&outConn)
	}, func(n int, d bool) {}, 0)
	log.Printf("conn %s - %s - %s -%s connected", inAddr, inLocalAddr, outLocalAddr, outAddr)
	return
}
func (s *TLS) OutToUDP(inConn *net.Conn) (err error) {
	return
}
func (s *TLS) InitOutConnPool() {
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

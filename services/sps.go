package services

import (
	"fmt"
	"log"
	"net"
	"runtime/debug"
	"snail007/proxy/utils"
	"snail007/proxy/utils/socks"
	"strconv"
)

type SPS struct {
	outPool utils.OutPool
	cfg     SPSArgs
}

func NewSPS() Service {
	return &SPS{
		outPool: utils.OutPool{},
		cfg:     SPSArgs{},
	}
}
func (s *SPS) CheckArgs() {
	if *s.cfg.Parent == "" {
		log.Fatalf("parent required for %s %s", s.cfg.Protocol(), *s.cfg.Local)
	}
	if *s.cfg.ParentType == "" {
		log.Fatalf("parent type unkown,use -T <tls|tcp>")
	}
	if *s.cfg.ParentType == TYPE_TLS || *s.cfg.LocalType == TYPE_TLS {
		s.cfg.CertBytes, s.cfg.KeyBytes = utils.TlsBytes(*s.cfg.CertFile, *s.cfg.KeyFile)
	}
}
func (s *SPS) InitService() {

}
func (s *SPS) StopService() {
	if s.outPool.Pool != nil {
		s.outPool.Pool.ReleaseAll()
	}
}
func (s *SPS) Start(args interface{}) (err error) {
	s.cfg = args.(SPSArgs)
	s.CheckArgs()
	log.Printf("use %s %s parent %s", *s.cfg.ParentType, *s.cfg.ParentServiceType, *s.cfg.Parent)
	s.InitService()

	host, port, _ := net.SplitHostPort(*s.cfg.Local)
	p, _ := strconv.Atoi(port)
	sc := utils.NewServerChannel(host, p)

	if *s.cfg.LocalType == TYPE_TCP {
		err = sc.ListenTCP(s.callback)
	} else if *s.cfg.LocalType == TYPE_TLS {
		err = sc.ListenTls(s.cfg.CertBytes, s.cfg.KeyBytes, s.callback)
	}
	if err != nil {
		return
	}
	log.Printf("%s http(s)+socks proxy on %s", s.cfg.Protocol(), (*sc.Listener).Addr())
	return
}

func (s *SPS) Clean() {
	s.StopService()
}
func (s *SPS) callback(inConn net.Conn) {
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
	default:
		err = fmt.Errorf("unkown parent type %s", *s.cfg.ParentType)
	}
	if err != nil {
		log.Printf("connect to %s parent %s fail, ERR:%s", *s.cfg.ParentType, *s.cfg.Parent, err)
		utils.CloseConn(&inConn)
	}
}
func (s *SPS) OutToTCP(inConn *net.Conn) (err error) {
	buf := make([]byte, 1024)
	n, err := (*inConn).Read(buf)
	header := buf[:n]
	if err != nil {
		log.Printf("ERR:%s", err)
		utils.CloseConn(inConn)
		return
	}
	fmt.Printf("%v", header[0])

	if header[0] == socks.VERSION_V5 {
		req, e := socks.NewMethodsRequest(*inConn, header)
		if e != nil {
			log.Printf("ERR:%s", e)
			utils.CloseConn(inConn)
			err = e.(error)
			return
		}
		fmt.Printf("address:%v", req.Version())
	}
	return
}

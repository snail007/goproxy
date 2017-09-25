package services

import "log"

type Tunnel struct {
	cfg TunnelArgs
}

func NewTunnel() Service {
	return &Tunnel{
		cfg: TunnelArgs{},
	}
}

func (s *Tunnel) InitService() {

}
func (s *Tunnel) Check() {
	if *s.cfg.Parent != "" {
		log.Printf("use tls parent %s", *s.cfg.Parent)
	} else {
		log.Fatalf("parent required")
	}
	if s.cfg.CertBytes == nil || s.cfg.KeyBytes == nil {
		log.Fatalf("cert and key file required")
	}
}
func (s *Tunnel) StopService() {

}
func (s *Tunnel) Start(args interface{}) (err error) {
	s.cfg = args.(TunnelArgs)
	s.Check()
	s.InitService()
	return
}
func (s *Tunnel) Clean() {
	s.StopService()
}

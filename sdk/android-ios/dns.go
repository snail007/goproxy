package proxy

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	logger "log"
	"net"
	"runtime/debug"
	"time"

	"golang.org/x/net/proxy"

	"github.com/miekg/dns"
	gocache "github.com/pmylund/go-cache"
	"github.com/snail007/goproxy/core/lib/kcpcfg"
	services "github.com/snail007/goproxy/services"
)

type DNSArgs struct {
	ParentServiceType *string
	ParentType        *string
	Parent            *string
	ParentAuth        *string
	ParentKey         *string
	ParentCompress    *bool
	KCP               kcpcfg.KCPConfigArgs
	CertFile          *string
	KeyFile           *string
	CaCertFile        *string
	Local             *string
	Timeout           *int
	RemoteDNSAddress  *string
	DNSTTL            *int
	CacheFile         *string
	LocalSocks5Port   *string
}
type DNS struct {
	cfg        DNSArgs
	log        *logger.Logger
	cache      *gocache.Cache
	exitSig    chan bool
	serviceKey string
	dialer     proxy.Dialer
}

func NewDNS() services.Service {
	return &DNS{
		cfg:        DNSArgs{},
		exitSig:    make(chan bool, 1),
		serviceKey: "dns-service-" + fmt.Sprintf("%d", time.Now().UnixNano()),
	}
}
func (s *DNS) CheckArgs() (err error) {
	return
}
func (s *DNS) InitService() (err error) {
	s.cache = gocache.New(time.Second*time.Duration(*s.cfg.DNSTTL), time.Second*60)
	s.cache.LoadFile(*s.cfg.CacheFile)
	go func() {
		for {
			select {
			case <-s.exitSig:
				return
			case <-time.After(time.Second * 300):
				s.cache.DeleteExpired()
				s.cache.SaveFile(*s.cfg.CacheFile)
			}
		}
	}()
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("crashed, err: %s\nstack:%s", e, string(debug.Stack()))
			}
		}()
		for {
			select {
			case <-s.exitSig:
				return
			case <-time.After(time.Second * 60):
				err := s.cache.SaveFile(*s.cfg.CacheFile)
				if err == nil {
					//s.log.Printf("cache saved: %s", *s.cfg.CacheFile)
				} else {
					s.log.Printf("cache save failed: %s, %s", *s.cfg.CacheFile, err)
				}
			}
		}
	}()
	s.dialer, err = proxy.SOCKS5("tcp", *s.cfg.Parent,
		nil,
		&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 2 * time.Second,
		},
	)
	if err != nil {
		return
	}

	sdkArgs := fmt.Sprintf("sps -S %s -T %s -P %s -C %s -K %s -i %d -p 127.0.0.1:%s --disable-http",
		*s.cfg.ParentServiceType,
		*s.cfg.ParentType,
		*s.cfg.Parent,
		*s.cfg.CertFile,
		*s.cfg.KeyFile,
		*s.cfg.Timeout,
		*s.cfg.LocalSocks5Port,
	)
	if *s.cfg.ParentKey != "" {
		sdkArgs += " -Z " + *s.cfg.ParentKey
	}
	if *s.cfg.ParentAuth != "" {
		sdkArgs += " -A " + *s.cfg.ParentAuth
	}
	if *s.cfg.CaCertFile != "" {
		sdkArgs += " --ca " + *s.cfg.CaCertFile
	}
	if *s.cfg.ParentCompress {
		sdkArgs += " -M"
	}
	s.log.Printf("start sps with : %s", sdkArgs)
	errStr := Start(s.serviceKey, sdkArgs)
	if errStr != "" {
		err = fmt.Errorf("start sps service fail,%s", errStr)
	}
	return
}
func (s *DNS) StopService() {
	defer func() {
		e := recover()
		if e != nil {
			s.log.Printf("stop dns service crashed,%s", e)
		} else {
			s.log.Printf("service dns stopped")
		}
	}()
	Stop(s.serviceKey)
	s.cache.Flush()
	s.exitSig <- true
}
func (s *DNS) Start(args interface{}, log *logger.Logger) (err error) {
	s.log = log
	s.cfg = args.(DNSArgs)
	if err = s.CheckArgs(); err != nil {
		return
	}
	if err = s.InitService(); err != nil {
		return
	}
	dns.HandleFunc(".", s.callback)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("crashed, err: %s\nstack:%s", e, string(debug.Stack()))
			}
		}()
		log.Printf("dns server on udp %s", *s.cfg.Local)
		err := dns.ListenAndServe(*s.cfg.Local, "udp", nil)
		if err != nil {
			log.Printf("dns listen error: %s", err)
		}
	}()
	return
}

func (s *DNS) Clean() {
	s.StopService()
}
func (s *DNS) callback(w dns.ResponseWriter, req *dns.Msg) {
	defer func() {
		if err := recover(); err != nil {
			s.log.Printf("dns handler crashed with err : %s \nstack: %s", err, string(debug.Stack()))
		}
	}()
	var (
		key       string
		m         *dns.Msg
		err       error
		data      []byte
		id        uint16
		query     []string
		questions []dns.Question
	)
	if req.MsgHdr.Response == true {
		return
	}
	query = make([]string, len(req.Question))
	for i, q := range req.Question {
		if q.Qtype != dns.TypeAAAA {
			questions = append(questions, q)
		}
		query[i] = fmt.Sprintf("(%s %s %s)", q.Name, dns.ClassToString[q.Qclass], dns.TypeToString[q.Qtype])
	}

	if len(questions) == 0 {
		return
	}

	req.Question = questions
	id = req.Id
	req.Id = 0
	key = s.toMd5(req.String())
	req.Id = id
	if reply, ok := s.cache.Get(key); ok {
		data, _ = reply.([]byte)
	}
	if data != nil && len(data) > 0 {
		m = &dns.Msg{}
		m.Unpack(data)
		m.Id = id
		err = w.WriteMsg(m)
		s.log.Printf("id: %5d cache: HIT %v", id, query)
		return

	} else {
		s.log.Printf("id: %5d cache: MISS %v", id, query)
	}

	s.log.Printf("id: %5d resolve: %v %s", id, query, *s.cfg.RemoteDNSAddress)

	rawConn, err := s.dialer.Dial("tcp", *s.cfg.RemoteDNSAddress)
	if err != nil {
		s.log.Printf("dail to %s fail,%s", *s.cfg.RemoteDNSAddress, err)
		return
	}
	defer rawConn.Close()
	co := new(dns.Conn)
	co.Conn = rawConn
	defer co.Close()
	if err = co.WriteMsg(req); err != nil {
		s.log.Printf("write dns query fail,%s", err)
		return
	}
	m, err = co.ReadMsg()
	if err == nil && m.Id != req.Id {
		s.log.Printf("id: %5d mismath", id)
		return
	}
	if err != nil || len(m.Answer) == 0 {
		s.log.Printf("dns query fail,%s", err)
		return
	}
	data, err = m.Pack()
	if err != nil {
		s.log.Printf("dns query fail,%s", err)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		s.log.Printf("dns query fail,%s", err)
		return
	}
	m.Id = 0
	data, _ = m.Pack()
	ttl := 0
	if len(m.Answer) > 0 {
		if *s.cfg.DNSTTL > 0 {
			ttl = *s.cfg.DNSTTL
		} else {
			ttl = int(m.Answer[0].Header().Ttl)
			if ttl < 0 {
				ttl = *s.cfg.DNSTTL
			}
		}
	}
	s.cache.Set(key, data, time.Second*time.Duration(ttl))
	m.Id = id
	s.log.Printf("id: %5d cache: CACHED %v TTL %v", id, query, ttl)
}
func (s *DNS) toMd5(data string) string {
	m := md5.New()
	m.Write([]byte(data))
	return hex.EncodeToString(m.Sum(nil))
}

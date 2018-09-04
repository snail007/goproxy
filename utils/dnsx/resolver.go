package dnsx

import (
	"fmt"
	logger "log"
	"net"
	"strings"
	"time"

	"github.com/snail007/goproxy/utils/mapx"
	dns "github.com/miekg/dns"
)

type DomainResolver struct {
	ttl         int
	dnsAddrress string
	data        mapx.ConcurrentMap
	log         *logger.Logger
}
type DomainResolverItem struct {
	ip        string
	domain    string
	expiredAt int64
}

func NewDomainResolver(dnsAddrress string, ttl int, log *logger.Logger) DomainResolver {
	return DomainResolver{
		ttl:         ttl,
		dnsAddrress: dnsAddrress,
		data:        mapx.NewConcurrentMap(),
		log:         log,
	}
}
func (a *DomainResolver) DnsAddress() (address string) {
	address = a.dnsAddrress
	return
}
func (a *DomainResolver) MustResolve(address string) (ip string) {
	ip, _ = a.Resolve(address)
	return
}
func (a *DomainResolver) Resolve(address string) (ip string, err error) {
	domain := address
	port := ""
	fromCache := "false"
	defer func() {
		if port != "" {
			ip = net.JoinHostPort(ip, port)
		}
		a.log.Printf("dns:%s->%s,cache:%s", address, ip, fromCache)
		//a.PrintData()
	}()
	if strings.Contains(domain, ":") {
		domain, port, err = net.SplitHostPort(domain)
		if err != nil {
			return
		}
	}
	if net.ParseIP(domain) != nil {
		ip = domain
		fromCache = "ip ignore"
		return
	}
	item, ok := a.data.Get(domain)
	if ok {
		//log.Println("find ", domain)
		if (*item.(*DomainResolverItem)).expiredAt > time.Now().Unix() {
			ip = (*item.(*DomainResolverItem)).ip
			fromCache = "true"
			//log.Println("from cache ", domain)
			return
		}
	} else {
		item = &DomainResolverItem{
			domain: domain,
		}

	}
	c := new(dns.Client)
	c.DialTimeout = time.Millisecond * 5000
	c.ReadTimeout = time.Millisecond * 5000
	c.WriteTimeout = time.Millisecond * 5000
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	m.RecursionDesired = true
	r, _, err := c.Exchange(m, a.dnsAddrress)
	if r == nil {
		return
	}
	if r.Rcode != dns.RcodeSuccess {
		err = fmt.Errorf(" *** invalid answer name %s after A query for %s", domain, a.dnsAddrress)
		return
	}
	for _, answer := range r.Answer {
		if answer.Header().Rrtype == dns.TypeA {
			info := strings.Fields(answer.String())
			if len(info) >= 5 {
				ip = info[4]
				_item := item.(*DomainResolverItem)
				(*_item).expiredAt = time.Now().Unix() + int64(a.ttl)
				(*_item).ip = ip
				a.data.Set(domain, item)
				return
			}
		}
	}
	return
}
func (a *DomainResolver) PrintData() {
	for k, item := range a.data.Items() {
		d := item.(*DomainResolverItem)
		a.log.Printf("%s:ip[%s],domain[%s],expired at[%d]\n", k, (*d).ip, (*d).domain, (*d).expiredAt)
	}
}

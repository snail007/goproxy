package lb

import (
	"crypto/md5"
	"log"
	"net"
	"sync"

	"github.com/snail007/goproxy/utils/dnsx"
)

const (
	SELECT_ROUNDROBIN = iota
	SELECT_LEASTCONN
	SELECT_HASH
	SELECT_WEITHT
	SELECT_LEASTTIME
)

type Selector interface {
	Select(srcAddr string) (addr string)
	SelectBackend(srcAddr string) (b *Backend)
	IncreasConns(addr string)
	DecreaseConns(addr string)
	Stop()
	Reset(configs BackendsConfig, dr *dnsx.DomainResolver, log *log.Logger)
	IsActive() bool
	ActiveCount() (count int)
	Backends() (bs []*Backend)
}

type Group struct {
	selector *Selector
	log      *log.Logger
	dr       *dnsx.DomainResolver
	lock     *sync.Mutex
	last     *Backend
	debug    bool
	bks      []*Backend
}

func NewGroup(selectType int, configs BackendsConfig, dr *dnsx.DomainResolver, log *log.Logger, debug bool) Group {
	bks := []*Backend{}
	for _, c := range configs {
		b, _ := NewBackend(*c, dr, log)
		bks = append(bks, b)
	}
	if len(bks) > 1 {
		for _, b := range bks {
			b.StartHeartCheck()
		}
	}
	var s Selector
	switch selectType {
	case SELECT_ROUNDROBIN:
		s = NewRoundRobin(bks, log, debug)
	case SELECT_LEASTCONN:
		s = NewLeastConn(bks, log, debug)
	case SELECT_HASH:
		s = NewHash(bks, log, debug)
	case SELECT_WEITHT:
		s = NewWeight(bks, log, debug)
	case SELECT_LEASTTIME:
		s = NewLeastTime(bks, log, debug)
	}
	return Group{
		selector: &s,
		log:      log,
		dr:       dr,
		lock:     &sync.Mutex{},
		debug:    debug,
		bks:      bks,
	}
}
func (g *Group) Select(srcAddr string, onlyHa bool) (addr string) {
	if len(g.bks) == 1 {
		return g.bks[0].Address
	}
	if onlyHa {
		g.lock.Lock()
		defer g.lock.Unlock()
		if g.last != nil && (g.last.Active || g.last.ConnectUsedMillisecond == 0) {
			if g.debug {
				g.log.Printf("############ choosed %s from lastest ############", g.last.Address)
				printDebug(true, g.log, nil, srcAddr, (*g.selector).Backends())
			}
			return g.last.Address
		}
		g.last = (*g.selector).SelectBackend(srcAddr)
		if !g.last.Active && g.last.ConnectUsedMillisecond > 0 {
			g.log.Printf("###warn### lb selected empty , return default , for : %s", srcAddr)
		}
		return g.last.Address
	}
	b := (*g.selector).SelectBackend(srcAddr)
	return b.Address

}
func (g *Group) IncreasConns(addr string) {
	(*g.selector).IncreasConns(addr)
}
func (g *Group) DecreaseConns(addr string) {
	(*g.selector).DecreaseConns(addr)
}
func (g *Group) Stop() {
	if g.selector != nil {
		(*g.selector).Stop()
	}
}
func (g *Group) IsActive() bool {
	return (*g.selector).IsActive()
}
func (g *Group) ActiveCount() (count int) {
	return (*g.selector).ActiveCount()
}
func (g *Group) Reset(addrs []string) {
	bks := (*g.selector).Backends()
	if len(bks) == 0 {
		return
	}
	cfg := bks[0].BackendConfig
	configs := BackendsConfig{}
	for _, addr := range addrs {
		c := cfg
		c.Address = addr
		configs = append(configs, &c)
	}
	(*g.selector).Reset(configs, g.dr, g.log)
	g.bks = (*g.selector).Backends()
}
func (g *Group) Backends() []*Backend {
	return (*g.selector).Backends()
}

//########################RoundRobin##########################
type RoundRobin struct {
	sync.Mutex
	backendIndex int
	backends     Backends
	log          *log.Logger
	debug        bool
}

func NewRoundRobin(backends Backends, log *log.Logger, debug bool) Selector {
	return &RoundRobin{
		backends: backends,
		log:      log,
		debug:    debug,
	}

}
func (r *RoundRobin) Select(srcAddr string) (addr string) {
	return r.SelectBackend(srcAddr).Address
}
func (r *RoundRobin) SelectBackend(srcAddr string) (b *Backend) {
	r.Lock()
	defer r.Unlock()
	defer func() {
		printDebug(r.debug, r.log, b, srcAddr, r.backends)
	}()
	if len(r.backends) == 0 {
		return
	}
	if len(r.backends) == 1 {
		return r.backends[0]
	}
RETRY:
	found := false
	for _, b := range r.backends {
		if b.Active {
			found = true
			break
		}
	}
	if !found {
		return r.backends[0]
	}
	r.backendIndex++
	if r.backendIndex > len(r.backends)-1 {
		r.backendIndex = 0
	}
	if !r.backends[r.backendIndex].Active {
		goto RETRY
	}
	return r.backends[r.backendIndex]
}
func (r *RoundRobin) IncreasConns(addr string) {

}
func (r *RoundRobin) DecreaseConns(addr string) {

}
func (r *RoundRobin) Stop() {
	for _, b := range r.backends {
		b.StopHeartCheck()
	}
}
func (r *RoundRobin) Backends() []*Backend {
	return r.backends
}
func (r *RoundRobin) IsActive() bool {
	for _, b := range r.backends {
		if b.Active {
			return true
		}
	}
	return false
}
func (r *RoundRobin) ActiveCount() (count int) {
	for _, b := range r.backends {
		if b.Active {
			count++
		}
	}
	return
}
func (r *RoundRobin) Reset(configs BackendsConfig, dr *dnsx.DomainResolver, log *log.Logger) {
	r.Lock()
	defer r.Unlock()
	r.Stop()
	bks := []*Backend{}
	for _, c := range configs {
		b, _ := NewBackend(*c, dr, log)
		bks = append(bks, b)
	}
	if len(bks) > 1 {
		for _, b := range bks {
			b.StartHeartCheck()
		}
	}
	r.backends = bks
}

//########################LeastConn##########################

type LeastConn struct {
	sync.Mutex
	backends Backends
	log      *log.Logger
	debug    bool
}

func NewLeastConn(backends []*Backend, log *log.Logger, debug bool) Selector {
	lc := LeastConn{
		backends: backends,
		log:      log,
		debug:    debug,
	}
	return &lc
}

func (lc *LeastConn) Select(srcAddr string) (addr string) {
	return lc.SelectBackend(srcAddr).Address
}
func (lc *LeastConn) SelectBackend(srcAddr string) (b *Backend) {
	lc.Lock()
	defer lc.Unlock()
	defer func() {
		printDebug(lc.debug, lc.log, b, srcAddr, lc.backends)
	}()
	if len(lc.backends) == 0 {
		return
	}
	if len(lc.backends) == 1 {
		return lc.backends[0]
	}
	found := false
	for _, b := range lc.backends {
		if b.Active {
			found = true
			break
		}
	}
	if !found {
		return lc.backends[0]
	}
	min := lc.backends[0].Connections
	index := 0
	for i, b := range lc.backends {
		if b.Active {
			min = b.Connections
			index = i
			break
		}
	}
	for i, b := range lc.backends {
		if b.Active && b.Connections <= min {
			min = b.Connections
			index = i
		}
	}
	return lc.backends[index]
}
func (lc *LeastConn) IncreasConns(addr string) {
	for _, a := range lc.backends {
		if a.Address == addr {
			a.IncreasConns()
			return
		}
	}
}
func (lc *LeastConn) DecreaseConns(addr string) {
	for _, a := range lc.backends {
		if a.Address == addr {
			a.DecreaseConns()
			return
		}
	}
}
func (lc *LeastConn) Stop() {
	for _, b := range lc.backends {
		b.StopHeartCheck()
	}
}
func (lc *LeastConn) IsActive() bool {
	for _, b := range lc.backends {
		if b.Active {
			return true
		}
	}
	return false
}
func (lc *LeastConn) ActiveCount() (count int) {
	for _, b := range lc.backends {
		if b.Active {
			count++
		}
	}
	return
}
func (lc *LeastConn) Reset(configs BackendsConfig, dr *dnsx.DomainResolver, log *log.Logger) {
	lc.Lock()
	defer lc.Unlock()
	lc.Stop()
	bks := []*Backend{}
	for _, c := range configs {
		b, _ := NewBackend(*c, dr, log)
		bks = append(bks, b)
	}
	if len(bks) > 1 {
		for _, b := range bks {
			b.StartHeartCheck()
		}
	}
	lc.backends = bks
}
func (lc *LeastConn) Backends() []*Backend {
	return lc.backends
}

//########################Hash##########################
type Hash struct {
	sync.Mutex
	backends Backends
	log      *log.Logger
	debug    bool
}

func NewHash(backends Backends, log *log.Logger, debug bool) Selector {
	return &Hash{
		backends: backends,
		log:      log,
		debug:    debug,
	}
}
func (h *Hash) Select(srcAddr string) (addr string) {
	return h.SelectBackend(srcAddr).Address
}
func (h *Hash) SelectBackend(srcAddr string) (b *Backend) {
	h.Lock()
	defer h.Unlock()
	defer func() {
		printDebug(h.debug, h.log, b, srcAddr, h.backends)
	}()
	if len(h.backends) == 0 {
		return
	}
	if len(h.backends) == 1 {
		return h.backends[0]
	}
	i := 0
	host, _, err := net.SplitHostPort(srcAddr)
	if err != nil {
		return
	}
	//porti, _ := strconv.Atoi(port)
	//i += porti
	for _, b := range md5.Sum([]byte(host)) {
		i += int(b)
	}
RETRY:
	found := false
	for _, b := range h.backends {
		if b.Active {
			found = true
			break
		}
	}
	if !found {
		return h.backends[0]
	}
	k := i % len(h.backends)
	if !h.backends[k].Active {
		i++
		goto RETRY
	}
	return h.backends[k]
}
func (h *Hash) IncreasConns(addr string) {

}
func (h *Hash) DecreaseConns(addr string) {

}
func (h *Hash) Stop() {
	for _, b := range h.backends {
		b.StopHeartCheck()
	}
}
func (h *Hash) IsActive() bool {
	for _, b := range h.backends {
		if b.Active {
			return true
		}
	}
	return false
}
func (h *Hash) ActiveCount() (count int) {
	for _, b := range h.backends {
		if b.Active {
			count++
		}
	}
	return
}
func (h *Hash) Reset(configs BackendsConfig, dr *dnsx.DomainResolver, log *log.Logger) {
	h.Lock()
	defer h.Unlock()
	h.Stop()
	bks := []*Backend{}
	for _, c := range configs {
		b, _ := NewBackend(*c, dr, log)
		bks = append(bks, b)
	}
	if len(bks) > 1 {
		for _, b := range bks {
			b.StartHeartCheck()
		}
	}
	h.backends = bks
}
func (h *Hash) Backends() []*Backend {
	return h.backends
}

//########################Weight##########################
type Weight struct {
	sync.Mutex
	backends Backends
	log      *log.Logger
	debug    bool
}

func NewWeight(backends Backends, log *log.Logger, debug bool) Selector {
	return &Weight{
		backends: backends,
		log:      log,
		debug:    debug,
	}
}
func (w *Weight) Select(srcAddr string) (addr string) {
	return w.SelectBackend(srcAddr).Address
}
func (w *Weight) SelectBackend(srcAddr string) (b *Backend) {
	w.Lock()
	defer w.Unlock()
	defer func() {
		printDebug(w.debug, w.log, b, srcAddr, w.backends)
	}()
	if len(w.backends) == 0 {
		return
	}
	if len(w.backends) == 1 {
		return w.backends[0]
	}

	found := false
	for _, b := range w.backends {
		if b.Active {
			found = true
			break
		}
	}
	if !found {
		return w.backends[0]
	}

	min := w.backends[0].Connections / w.backends[0].Weight
	index := 0
	for i, b := range w.backends {
		if b.Active {
			min = b.Connections / b.Weight
			index = i
			break
		}
	}
	for i, b := range w.backends {
		if b.Active && b.Connections/b.Weight <= min {
			min = b.Connections
			index = i
		}
	}
	return w.backends[index]
}
func (w *Weight) IncreasConns(addr string) {
	w.Lock()
	defer w.Unlock()
	for _, a := range w.backends {
		if a.Address == addr {
			a.IncreasConns()
			return
		}
	}
}
func (w *Weight) DecreaseConns(addr string) {
	w.Lock()
	defer w.Unlock()
	for _, a := range w.backends {
		if a.Address == addr {
			a.DecreaseConns()
			return
		}
	}
}
func (w *Weight) Stop() {
	for _, b := range w.backends {
		b.StopHeartCheck()
	}
}
func (w *Weight) IsActive() bool {
	for _, b := range w.backends {
		if b.Active {
			return true
		}
	}
	return false
}
func (w *Weight) ActiveCount() (count int) {
	for _, b := range w.backends {
		if b.Active {
			count++
		}
	}
	return
}
func (w *Weight) Reset(configs BackendsConfig, dr *dnsx.DomainResolver, log *log.Logger) {
	w.Lock()
	defer w.Unlock()
	w.Stop()
	bks := []*Backend{}
	for _, c := range configs {
		b, _ := NewBackend(*c, dr, log)
		bks = append(bks, b)
	}
	if len(bks) > 1 {
		for _, b := range bks {
			b.StartHeartCheck()
		}
	}
	w.backends = bks
}
func (w *Weight) Backends() []*Backend {
	return w.backends
}

//########################LeastTime##########################

type LeastTime struct {
	sync.Mutex
	backends Backends
	log      *log.Logger
	debug    bool
}

func NewLeastTime(backends []*Backend, log *log.Logger, debug bool) Selector {
	lt := LeastTime{
		backends: backends,
		log:      log,
		debug:    debug,
	}
	return &lt
}

func (lt *LeastTime) Select(srcAddr string) (addr string) {
	return lt.SelectBackend(srcAddr).Address
}
func (lt *LeastTime) SelectBackend(srcAddr string) (b *Backend) {
	lt.Lock()
	defer lt.Unlock()
	defer func() {
		printDebug(lt.debug, lt.log, b, srcAddr, lt.backends)
	}()
	if len(lt.backends) == 0 {
		return
	}
	if len(lt.backends) == 1 {
		return lt.backends[0]
	}
	found := false
	for _, b := range lt.backends {
		if b.Active {
			found = true
			break
		}
	}
	if !found {
		return lt.backends[0]
	}
	min := lt.backends[0].ConnectUsedMillisecond
	index := 0
	for i, b := range lt.backends {
		if b.Active {
			min = b.ConnectUsedMillisecond
			index = i
			break
		}
	}
	for i, b := range lt.backends {
		if b.Active && b.ConnectUsedMillisecond > 0 && b.ConnectUsedMillisecond <= min {
			min = b.ConnectUsedMillisecond
			index = i
		}
	}
	return lt.backends[index]
}
func (lt *LeastTime) IncreasConns(addr string) {

}
func (lt *LeastTime) DecreaseConns(addr string) {

}
func (lt *LeastTime) Stop() {
	for _, b := range lt.backends {
		b.StopHeartCheck()
	}
}
func (lt *LeastTime) IsActive() bool {
	for _, b := range lt.backends {
		if b.Active {
			return true
		}
	}
	return false
}
func (lt *LeastTime) ActiveCount() (count int) {
	for _, b := range lt.backends {
		if b.Active {
			count++
		}
	}
	return
}
func (lt *LeastTime) Reset(configs BackendsConfig, dr *dnsx.DomainResolver, log *log.Logger) {
	lt.Lock()
	defer lt.Unlock()
	lt.Stop()
	bks := []*Backend{}
	for _, c := range configs {
		b, _ := NewBackend(*c, dr, log)
		bks = append(bks, b)
	}
	if len(bks) > 1 {
		for _, b := range bks {
			b.StartHeartCheck()
		}
	}
	lt.backends = bks
}
func (lt *LeastTime) Backends() []*Backend {
	return lt.backends
}
func printDebug(isDebug bool, log *log.Logger, selected *Backend, srcAddr string, backends []*Backend) {
	if isDebug {
		log.Printf("############ LB start ############\n")
		if selected != nil {
			log.Printf("choosed %s for %s\n", selected.Address, srcAddr)
		}
		for _, v := range backends {
			log.Printf("addr:%s,conns:%d,time:%d,weight:%d,active:%v\n", v.Address, v.Connections, v.ConnectUsedMillisecond, v.Weight, v.Active)
		}
		log.Printf("############ LB end ############\n")
	}
}

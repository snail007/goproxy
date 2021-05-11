package utils

import (
	"log"
	"sync"
	"time"
)

//ConnPool to use
type ConnPool interface {
	Get() (conn interface{}, err error)
	Put(conn interface{})
	ReleaseAll()
	Len() (length int)
}
type poolConfig struct {
	Factory    func() (interface{}, error)
	IsActive   func(interface{}) bool
	Release    func(interface{})
	InitialCap int
	MaxCap     int
}

func NewConnPool(poolConfig poolConfig) (pool ConnPool, err error) {
	p := netPool{
		config: poolConfig,
		conns:  make(chan interface{}, poolConfig.MaxCap),
		lock:   &sync.Mutex{},
	}
	//log.Printf("pool MaxCap:%d", poolConfig.MaxCap)
	if poolConfig.MaxCap > 0 {
		err = p.initAutoFill(false)
		if err == nil {
			p.initAutoFill(true)
		}
	}
	return &p, nil
}

type netPool struct {
	conns  chan interface{}
	lock   *sync.Mutex
	config poolConfig
}

func (p *netPool) initAutoFill(async bool) (err error) {
	var worker = func() (err error) {
		for {
			//log.Printf("pool fill: %v , len: %d", p.Len() <= p.config.InitialCap/2, p.Len())
			if p.Len() <= p.config.InitialCap/2 {
				p.lock.Lock()
				errN := 0
				for i := 0; i < p.config.InitialCap; i++ {
					c, err := p.config.Factory()
					if err != nil {
						errN++
						if async {
							continue
						} else {
							p.lock.Unlock()
							return err
						}
					}
					select {
					case p.conns <- c:
					default:
						p.config.Release(c)
						break
					}
					if p.Len() >= p.config.InitialCap {
						break
					}
				}
				if errN > 0 {
					log.Printf("fill conn pool fail , ERRN:%d", errN)
				}
				p.lock.Unlock()
			}
			if !async {
				return
			}
			time.Sleep(time.Second * 2)
		}
	}
	if async {
		go worker()
	} else {
		err = worker()
	}
	return

}

func (p *netPool) Get() (conn interface{}, err error) {
	// defer func() {
	// 	log.Printf("pool len : %d", p.Len())
	// }()
	p.lock.Lock()
	defer p.lock.Unlock()
	// for {
	select {
	case conn = <-p.conns:
		if p.config.IsActive(conn) {
			return
		}
		p.config.Release(conn)
	default:
		conn, err = p.config.Factory()
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
	// }
	return
}

func (p *netPool) Put(conn interface{}) {
	if conn == nil {
		return
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	if !p.config.IsActive(conn) {
		p.config.Release(conn)
	}
	select {
	case p.conns <- conn:
	default:
		p.config.Release(conn)
	}
}
func (p *netPool) ReleaseAll() {
	p.lock.Lock()
	defer p.lock.Unlock()
	close(p.conns)
	for c := range p.conns {
		p.config.Release(c)
	}
	p.conns = make(chan interface{}, p.config.InitialCap)

}
func (p *netPool) Len() (length int) {
	return len(p.conns)
}

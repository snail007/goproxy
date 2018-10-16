package lb

import (
	"errors"
	"fmt"
	"log"
	"net"
	"runtime/debug"
	"sync"
	"time"

	"github.com/snail007/goproxy/utils/dnsx"
)

// BackendConfig it's the configuration loaded
type BackendConfig struct {
	Address string

	ActiveAfter   int
	InactiveAfter int
	Weight        int

	Timeout   time.Duration
	RetryTime time.Duration

	IsMuxCheck  bool
	ConnFactory func(address string, timeout time.Duration) (net.Conn, error)
}
type BackendsConfig []*BackendConfig

// BackendControl keep the control data
type BackendControl struct {
	Failed bool // The last request failed
	Active bool

	InactiveTries int
	ActiveTries   int

	Connections int

	ConnectUsedMillisecond int

	isStop bool
}

// Backend structure
type Backend struct {
	BackendConfig
	BackendControl
	sync.RWMutex
	log *log.Logger
	dr  *dnsx.DomainResolver
}

type Backends []*Backend

func NewBackend(backendConfig BackendConfig, dr *dnsx.DomainResolver, log *log.Logger) (*Backend, error) {

	if backendConfig.Address == "" {
		return nil, errors.New("Address rquired")
	}
	if backendConfig.ActiveAfter == 0 {
		backendConfig.ActiveAfter = 2
	}
	if backendConfig.InactiveAfter == 0 {
		backendConfig.InactiveAfter = 3
	}
	if backendConfig.Weight == 0 {
		backendConfig.Weight = 1
	}
	if backendConfig.Timeout == 0 {
		backendConfig.Timeout = time.Millisecond * 1500
	}
	if backendConfig.RetryTime == 0 {
		backendConfig.RetryTime = time.Millisecond * 2000
	}
	return &Backend{
		dr:            dr,
		log:           log,
		BackendConfig: backendConfig,
		BackendControl: BackendControl{
			Failed:                 true,
			Active:                 false,
			InactiveTries:          0,
			ActiveTries:            0,
			Connections:            0,
			ConnectUsedMillisecond: 0,
			isStop:                 false,
		},
	}, nil
}
func (b *Backend) StopHeartCheck() {
	b.isStop = true
}

func (b *Backend) IncreasConns() {
	b.RWMutex.Lock()
	b.Connections++
	b.RWMutex.Unlock()
}

func (b *Backend) DecreaseConns() {
	b.RWMutex.Lock()
	b.Connections--
	b.RWMutex.Unlock()
}

func (b *Backend) StartHeartCheck() {
	if b.IsMuxCheck {
		b.startMuxHeartCheck()
	} else {
		b.startTCPHeartCheck()
	}
}
func (b *Backend) startMuxHeartCheck() {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("crashed, err: %s\nstack:\n%s",e, string(debug.Stack()))
			}
		}()
		for {
			if b.isStop {
				return
			}
			var c net.Conn
			var err error
			start := time.Now().UnixNano() / int64(time.Microsecond)
			c, err = b.getConn()
			b.ConnectUsedMillisecond = int(time.Now().UnixNano()/int64(time.Microsecond) - start)
			if err != nil {
				b.Active = false
				time.Sleep(time.Second * 2)
				continue
			} else {
				b.Active = true
			}
			for {
				buf := make([]byte, 1)
				c.Read(buf)
				buf = nil
				break
			}
			b.Active = false
		}
	}()
}

// Monitoring the backend
func (b *Backend) startTCPHeartCheck() {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("crashed, err: %s\nstack:\n%s",e, string(debug.Stack()))
			}
		}()
		for {
			if b.isStop {
				return
			}
			var c net.Conn
			var err error
			start := time.Now().UnixNano() / int64(time.Microsecond)
			c, err = b.getConn()
			b.ConnectUsedMillisecond = int(time.Now().UnixNano()/int64(time.Microsecond) - start)
			if err == nil {
				c.Close()
			}
			if err != nil {
				b.RWMutex.Lock()
				// Max tries before consider inactive
				if b.InactiveTries >= b.InactiveAfter {
					//b.log.Printf("Backend inactive [%s]", b.Address)
					b.Active = false
					b.ActiveTries = 0
				} else {
					// Ok that guy it's out of the game
					b.Failed = true
					b.InactiveTries++
					//b.log.Printf("Error to check address [%s] tries [%d]", b.Address, b.InactiveTries)
				}
				b.RWMutex.Unlock()
			} else {

				// Ok, let's keep working boys
				b.RWMutex.Lock()
				if b.ActiveTries >= b.ActiveAfter {
					if b.Failed {
						//log.Printf("Backend active [%s]", b.Address)
					}
					b.Failed = false
					b.Active = true
					b.InactiveTries = 0
				} else {
					b.ActiveTries++
				}
				b.RWMutex.Unlock()
			}
			time.Sleep(b.RetryTime)
		}
	}()
}
func (b *Backend) getConn() (conn net.Conn, err error) {
	address := b.Address
	if b.dr != nil && b.dr.DnsAddress() != "" {
		address, err = b.dr.Resolve(b.Address)
		if err != nil {
			b.log.Printf("dns error %s , ERR:%s", b.Address, err)
		}
	}
	if b.ConnFactory != nil {
		return b.ConnFactory(address, b.Timeout)
	}
	return net.DialTimeout("tcp", address, b.Timeout)
}

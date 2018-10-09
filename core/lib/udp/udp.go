package udputils

import (
	"fmt"
	logger "log"
	"net"
	"runtime/debug"
	"strings"
	"time"

	bufx "github.com/snail007/goproxy/core/lib/buf"
	mapx "github.com/snail007/goproxy/core/lib/mapx"
)

type CreateOutUDPConnFn func(listener *net.UDPConn, srcAddr *net.UDPAddr, packet []byte) (outconn *net.UDPConn, err error)
type CleanFn func(srcAddr string)
type BeforeSendFn func(listener *net.UDPConn, srcAddr *net.UDPAddr, b []byte) (sendB []byte, err error)
type BeforeReplyFn func(listener *net.UDPConn, srcAddr *net.UDPAddr, outconn *net.UDPConn, b []byte) (replyB []byte, err error)

type IOBinder struct {
	outConns           mapx.ConcurrentMap
	listener           *net.UDPConn
	createOutUDPConnFn CreateOutUDPConnFn
	log                *logger.Logger
	timeout            time.Duration
	cleanFn            CleanFn
	inTCPConn          *net.Conn
	outTCPConn         *net.Conn
	beforeSendFn       BeforeSendFn
	beforeReplyFn      BeforeReplyFn
}

func NewIOBinder(listener *net.UDPConn, log *logger.Logger) *IOBinder {
	return &IOBinder{
		listener: listener,
		outConns: mapx.NewConcurrentMap(),
		log:      log,
	}
}
func (s *IOBinder) Factory(fn CreateOutUDPConnFn) *IOBinder {
	s.createOutUDPConnFn = fn
	return s
}
func (s *IOBinder) AfterReadFromClient(fn BeforeSendFn) *IOBinder {
	s.beforeSendFn = fn
	return s
}
func (s *IOBinder) AfterReadFromServer(fn BeforeReplyFn) *IOBinder {
	s.beforeReplyFn = fn
	return s
}
func (s *IOBinder) Timeout(timeout time.Duration) *IOBinder {
	s.timeout = timeout
	return s
}
func (s *IOBinder) Clean(fn CleanFn) *IOBinder {
	s.cleanFn = fn
	return s
}
func (s *IOBinder) AliveWithServeConn(srcAddr string, inTCPConn *net.Conn) *IOBinder {
	s.inTCPConn = inTCPConn
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("crashed, err: %s\nstack:\n%s",e, string(debug.Stack()))
			}
		}()
		buf := make([]byte, 1)
		(*inTCPConn).SetReadDeadline(time.Time{})
		if _, err := (*inTCPConn).Read(buf); err != nil {
			s.log.Printf("udp related tcp conn of client disconnected with read , %s", err.Error())
			s.clean(srcAddr)
		}
	}()
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("crashed, err: %s\nstack:\n%s",e, string(debug.Stack()))
			}
		}()
		for {
			(*inTCPConn).SetWriteDeadline(time.Now().Add(time.Second * 5))
			if _, err := (*inTCPConn).Write([]byte{0x00}); err != nil {
				s.log.Printf("udp related tcp conn of client disconnected with write , %s", err.Error())
				s.clean(srcAddr)
				return
			}
			(*inTCPConn).SetWriteDeadline(time.Time{})
			time.Sleep(time.Second * 5)
		}
	}()
	return s
}
func (s *IOBinder) AliveWithClientConn(srcAddr string, outTCPConn *net.Conn) *IOBinder {
	s.outTCPConn = outTCPConn
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("crashed, err: %s\nstack:\n%s",e, string(debug.Stack()))
			}
		}()
		buf := make([]byte, 1)
		(*outTCPConn).SetReadDeadline(time.Time{})
		if _, err := (*outTCPConn).Read(buf); err != nil {
			s.log.Printf("udp related tcp conn to parent disconnected with read , %s", err.Error())
			s.clean(srcAddr)
		}
	}()
	return s
}
func (s *IOBinder) Run() (err error) {
	var (
		isClosedErr = func(err error) bool {
			return err != nil && strings.Contains(err.Error(), "use of closed network connection")
		}
		isTimeoutErr = func(err error) bool {
			if err == nil {
				return false
			}
			e, ok := err.(net.Error)
			return ok && e.Timeout()
		}
		isRefusedErr = func(err error) bool {
			return err != nil && strings.Contains(err.Error(), "connection refused")
		}
	)
	for {
		buf := bufx.Get()
		defer bufx.Put(buf)
		n, srcAddr, err := s.listener.ReadFromUDP(buf)
		if err != nil {
			s.log.Printf("read from client error %s", err)
			if isClosedErr(err) {
				return err
			}
			continue
		}
		var data []byte
		if s.beforeSendFn != nil {
			data, err = s.beforeSendFn(s.listener, srcAddr, buf[:n])
			if err != nil {
				s.log.Printf("beforeSend retured an error , %s", err)
				continue
			}
		} else {
			data = buf[:n]
		}
		inconnRemoteAddr := srcAddr.String()
		var outconn *net.UDPConn
		if v, ok := s.outConns.Get(inconnRemoteAddr); !ok {
			outconn, err = s.createOutUDPConnFn(s.listener, srcAddr, data)
			if err != nil {
				s.log.Printf("connnect fail %s", err)
				return err
			}
			go func() {
				defer func() {
					if e := recover(); e != nil {
						fmt.Printf("crashed, err: %s\nstack:\n%s",e, string(debug.Stack()))
					}
				}()
				defer func() {
					s.clean(srcAddr.String())
				}()
				buf := bufx.Get()
				defer bufx.Put(buf)
				for {
					if s.timeout > 0 {
						outconn.SetReadDeadline(time.Now().Add(s.timeout))
					}
					n, srcAddr, err := outconn.ReadFromUDP(buf)
					if err != nil {
						s.log.Printf("read from remote error %s", err)
						if isClosedErr(err) || isTimeoutErr(err) || isRefusedErr(err) {
							return
						}
						continue
					}
					data := buf[:n]
					if s.beforeReplyFn != nil {
						data, err = s.beforeReplyFn(s.listener, srcAddr, outconn, buf[:n])
						if err != nil {
							s.log.Printf("beforeReply retured an error , %s", err)
							continue
						}
					}
					_, err = s.listener.WriteTo(data, srcAddr)
					if err != nil {
						s.log.Printf("write to remote error %s", err)
						if isClosedErr(err) {
							return
						}
						continue
					}
				}
			}()
		} else {
			outconn = v.(*net.UDPConn)
		}

		s.log.Printf("use decrpyted data , %v", data)

		_, err = outconn.Write(data)

		if err != nil {
			s.log.Printf("write to remote error %s", err)
			if isClosedErr(err) {
				return err
			}
		}
	}
}
func (s *IOBinder) clean(srcAddr string) *IOBinder {
	if v, ok := s.outConns.Get(srcAddr); ok {
		(*v.(*net.UDPConn)).Close()
		s.outConns.Remove(srcAddr)
	}
	if s.inTCPConn != nil {
		(*s.inTCPConn).Close()
	}
	if s.outTCPConn != nil {
		(*s.outTCPConn).Close()
	}
	if s.cleanFn != nil {
		s.cleanFn(srcAddr)
	}
	return s
}

func (s *IOBinder) Close() {
	for _, c := range s.outConns.Items() {
		(*c.(*net.UDPConn)).Close()
	}
}

package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"

	logger "log"
	"net"
	"runtime/debug"
	"strconv"

	tou "github.com/snail007/goproxy/core/dst"
	compressconn "github.com/snail007/goproxy/core/lib/transport"
	transportc "github.com/snail007/goproxy/core/lib/transport"
	encryptconn "github.com/snail007/goproxy/core/lib/transport/encrypt"

	"github.com/snail007/goproxy/core/lib/kcpcfg"

	kcp "github.com/xtaci/kcp-go"
)

func init() {

}

type ServerChannel struct {
	ip               string
	port             int
	Listener         *net.Listener
	UDPListener      *net.UDPConn
	errAcceptHandler func(err error)
	log              *logger.Logger
	TOUServer        *tou.Mux
}

func NewServerChannel(ip string, port int, log *logger.Logger) ServerChannel {
	return ServerChannel{
		ip:   ip,
		port: port,
		log:  log,
		errAcceptHandler: func(err error) {
			log.Printf("accept error , ERR:%s", err)
		},
	}
}
func NewServerChannelHost(host string, log *logger.Logger) ServerChannel {
	h, port, _ := net.SplitHostPort(host)
	p, _ := strconv.Atoi(port)
	return ServerChannel{
		ip:   h,
		port: p,
		log:  log,
		errAcceptHandler: func(err error) {
			log.Printf("accept error , ERR:%s", err)
		},
	}
}
func (s *ServerChannel) SetErrAcceptHandler(fn func(err error)) {
	s.errAcceptHandler = fn
}
func (s *ServerChannel) ListenSingleTLS(certBytes, keyBytes, caCertBytes []byte, fn func(conn net.Conn)) (err error) {
	return s._ListenTLS(certBytes, keyBytes, caCertBytes, fn, true)

}
func (s *ServerChannel) ListenTLS(certBytes, keyBytes, caCertBytes []byte, fn func(conn net.Conn)) (err error) {
	return s._ListenTLS(certBytes, keyBytes, caCertBytes, fn, false)
}
func (s *ServerChannel) _ListenTLS(certBytes, keyBytes, caCertBytes []byte, fn func(conn net.Conn), single bool) (err error) {
	s.Listener, err = s.listenTLS(s.ip, s.port, certBytes, keyBytes, caCertBytes, single)
	if err == nil {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					s.log.Printf("ListenTLS crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
				}
			}()
			for {
				var conn net.Conn
				conn, err = (*s.Listener).Accept()
				if err == nil {
					go func() {
						defer func() {
							if e := recover(); e != nil {
								s.log.Printf("tls connection handler crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
							}
						}()
						fn(conn)
					}()
				} else {
					s.errAcceptHandler(err)
					(*s.Listener).Close()
					break
				}
			}
		}()
	}
	return
}
func (s *ServerChannel) listenTLS(ip string, port int, certBytes, keyBytes, caCertBytes []byte, single bool) (ln *net.Listener, err error) {
	var cert tls.Certificate
	cert, err = tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		return
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	if !single {
		clientCertPool := x509.NewCertPool()
		caBytes := certBytes
		if caCertBytes != nil {
			caBytes = caCertBytes
		}
		ok := clientCertPool.AppendCertsFromPEM(caBytes)
		if !ok {
			err = errors.New("failed to parse root certificate")
		}
		config.ClientCAs = clientCertPool
		config.ClientAuth = tls.RequireAndVerifyClientCert
	}
	_ln, err := tls.Listen("tcp", net.JoinHostPort(ip, fmt.Sprintf("%d", port)), config)
	if err == nil {
		ln = &_ln
	}
	return
}
func (s *ServerChannel) ListenTCPS(method, password string, compress bool, fn func(conn net.Conn)) (err error) {
	_, err = encryptconn.NewCipher(method, password)
	if err != nil {
		return
	}
	return s.ListenTCP(func(c net.Conn) {
		if compress {
			c = transportc.NewCompConn(c)
		}
		c, _ = encryptconn.NewConn(c, method, password)
		fn(c)
	})
}
func (s *ServerChannel) ListenTCP(fn func(conn net.Conn)) (err error) {
	var l net.Listener
	l, err = net.Listen("tcp", net.JoinHostPort(s.ip, fmt.Sprintf("%d", s.port)))
	if err == nil {
		s.Listener = &l
		go func() {
			defer func() {
				if e := recover(); e != nil {
					s.log.Printf("ListenTCP crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
				}
			}()
			for {
				var conn net.Conn
				conn, err = (*s.Listener).Accept()
				if err == nil {
					go func() {
						defer func() {
							if e := recover(); e != nil {
								s.log.Printf("tcp connection handler crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
							}
						}()
						fn(conn)
					}()
				} else {
					s.errAcceptHandler(err)
					(*s.Listener).Close()
					break
				}
			}
		}()
	}
	return
}
func (s *ServerChannel) ListenUDP(fn func(listener *net.UDPConn, packet []byte, localAddr, srcAddr *net.UDPAddr)) (err error) {
	addr := &net.UDPAddr{IP: net.ParseIP(s.ip), Port: s.port}
	l, err := net.ListenUDP("udp", addr)
	if err == nil {
		s.UDPListener = l
		go func() {
			defer func() {
				if e := recover(); e != nil {
					s.log.Printf("ListenUDP crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
				}
			}()
			for {
				var buf = make([]byte, 2048)
				n, srcAddr, err := (*s.UDPListener).ReadFromUDP(buf)
				if err == nil {
					packet := buf[0:n]
					go func() {
						defer func() {
							if e := recover(); e != nil {
								s.log.Printf("udp data handler crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
							}
						}()
						fn(s.UDPListener, packet, addr, srcAddr)
					}()
				} else {
					s.errAcceptHandler(err)
					(*s.UDPListener).Close()
					break
				}
			}
		}()
	}
	return
}
func (s *ServerChannel) ListenKCP(config kcpcfg.KCPConfigArgs, fn func(conn net.Conn), log *logger.Logger) (err error) {
	lis, err := kcp.ListenWithOptions(net.JoinHostPort(s.ip, fmt.Sprintf("%d", s.port)), config.Block, *config.DataShard, *config.ParityShard)
	if err == nil {
		if err = lis.SetDSCP(*config.DSCP); err != nil {
			log.Println("SetDSCP:", err)
			return
		}
		if err = lis.SetReadBuffer(*config.SockBuf); err != nil {
			log.Println("SetReadBuffer:", err)
			return
		}
		if err = lis.SetWriteBuffer(*config.SockBuf); err != nil {
			log.Println("SetWriteBuffer:", err)
			return
		}
		s.Listener = new(net.Listener)
		*s.Listener = lis
		go func() {
			defer func() {
				if e := recover(); e != nil {
					s.log.Printf("ListenKCP crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
				}
			}()
			for {
				conn, err := lis.AcceptKCP()
				if err == nil {
					go func() {
						defer func() {
							if e := recover(); e != nil {
								s.log.Printf("kcp connection handler crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
							}
						}()
						conn.SetStreamMode(true)
						conn.SetWriteDelay(true)
						conn.SetNoDelay(*config.NoDelay, *config.Interval, *config.Resend, *config.NoCongestion)
						conn.SetMtu(*config.MTU)
						conn.SetWindowSize(*config.SndWnd, *config.RcvWnd)
						conn.SetACKNoDelay(*config.AckNodelay)
						if *config.NoComp {
							fn(conn)
						} else {
							cconn := transportc.NewCompStream(conn)
							fn(cconn)
						}
					}()
				} else {
					s.errAcceptHandler(err)
					(*s.Listener).Close()
					break
				}
			}
		}()
	}
	return
}

func (s *ServerChannel) ListenTOU(method, password string, compress bool, fn func(conn net.Conn)) (err error) {
	addr := &net.UDPAddr{IP: net.ParseIP(s.ip), Port: s.port}
	s.UDPListener, err = net.ListenUDP("udp", addr)
	if err != nil {
		s.log.Println(err)
		return
	}
	s.TOUServer = tou.NewMux(s.UDPListener, 0)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				s.log.Printf("ListenRUDP crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
			}
		}()
		for {
			var conn net.Conn
			conn, err = (*s.TOUServer).Accept()
			if err == nil {
				go func() {
					defer func() {
						if e := recover(); e != nil {
							s.log.Printf("tcp connection handler crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
						}
					}()
					if compress {
						conn = compressconn.NewCompConn(conn)
					}
					conn, err = encryptconn.NewConn(conn, method, password)
					if err != nil {
						conn.Close()
						s.log.Println(err)
						return
					}
					fn(conn)
				}()
			} else {
				s.errAcceptHandler(err)
				s.TOUServer.Close()
				s.UDPListener.Close()
				break
			}
		}
	}()

	return
}
func (s *ServerChannel) Close() {
	defer func() {
		if e := recover(); e != nil {
			s.log.Printf("close crashed :\n%s\n%s", e, string(debug.Stack()))
		}
	}()
	if s.Listener != nil && *s.Listener != nil {
		(*s.Listener).Close()
	}
	if s.TOUServer != nil {
		s.TOUServer.Close()
	}
	if s.UDPListener != nil {
		s.UDPListener.Close()
	}
}
func (s *ServerChannel) Addr() string {
	defer func() {
		if e := recover(); e != nil {
			s.log.Printf("close crashed :\n%s\n%s", e, string(debug.Stack()))
		}
	}()
	if s.Listener != nil && *s.Listener != nil {
		return (*s.Listener).Addr().String()
	}

	if s.UDPListener != nil {
		return s.UDPListener.LocalAddr().String()
	}
	return ""
}

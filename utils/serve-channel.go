package utils

import (
	"fmt"
	"github.com/snail007/goproxy/services/kcpcfg"
	"log"
	"net"
	"runtime/debug"
	"strconv"

	kcp "github.com/xtaci/kcp-go"
)

type ServerChannel struct {
	ip               string
	port             int
	Listener         *net.Listener
	UDPListener      *net.UDPConn
	errAcceptHandler func(err error)
}

func NewServerChannel(ip string, port int) ServerChannel {
	return ServerChannel{
		ip:   ip,
		port: port,
		errAcceptHandler: func(err error) {
			log.Printf("accept error , ERR:%s", err)
		},
	}
}
func NewServerChannelHost(host string) ServerChannel {
	h, port, _ := net.SplitHostPort(host)
	p, _ := strconv.Atoi(port)
	return ServerChannel{
		ip:   h,
		port: p,
		errAcceptHandler: func(err error) {
			log.Printf("accept error , ERR:%s", err)
		},
	}
}
func (sc *ServerChannel) SetErrAcceptHandler(fn func(err error)) {
	sc.errAcceptHandler = fn
}
func (sc *ServerChannel) ListenTls(certBytes, keyBytes, caCertBytes []byte, fn func(conn net.Conn)) (err error) {
	sc.Listener, err = ListenTls(sc.ip, sc.port, certBytes, keyBytes, caCertBytes)
	if err == nil {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					log.Printf("ListenTls crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
				}
			}()
			for {
				var conn net.Conn
				conn, err = (*sc.Listener).Accept()
				if err == nil {
					go func() {
						defer func() {
							if e := recover(); e != nil {
								log.Printf("tls connection handler crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
							}
						}()
						fn(conn)
					}()
				} else {
					sc.errAcceptHandler(err)
					(*sc.Listener).Close()
					break
				}
			}
		}()
	}
	return
}

func (sc *ServerChannel) ListenTCP(fn func(conn net.Conn)) (err error) {
	var l net.Listener
	l, err = net.Listen("tcp", fmt.Sprintf("%s:%d", sc.ip, sc.port))
	if err == nil {
		sc.Listener = &l
		go func() {
			defer func() {
				if e := recover(); e != nil {
					log.Printf("ListenTCP crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
				}
			}()
			for {
				var conn net.Conn
				conn, err = (*sc.Listener).Accept()
				if err == nil {
					go func() {
						defer func() {
							if e := recover(); e != nil {
								log.Printf("tcp connection handler crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
							}
						}()
						fn(conn)
					}()
				} else {
					sc.errAcceptHandler(err)
					break
				}
			}
		}()
	}
	return
}
func (sc *ServerChannel) ListenUDP(fn func(packet []byte, localAddr, srcAddr *net.UDPAddr)) (err error) {
	addr := &net.UDPAddr{IP: net.ParseIP(sc.ip), Port: sc.port}
	l, err := net.ListenUDP("udp", addr)
	if err == nil {
		sc.UDPListener = l
		go func() {
			defer func() {
				if e := recover(); e != nil {
					log.Printf("ListenUDP crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
				}
			}()
			for {
				var buf = make([]byte, 2048)
				n, srcAddr, err := (*sc.UDPListener).ReadFromUDP(buf)
				if err == nil {
					packet := buf[0:n]
					go func() {
						defer func() {
							if e := recover(); e != nil {
								log.Printf("udp data handler crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
							}
						}()
						fn(packet, addr, srcAddr)
					}()
				} else {
					sc.errAcceptHandler(err)
					break
				}
			}
		}()
	}
	return
}
func (sc *ServerChannel) ListenKCP(config kcpcfg.KCPConfigArgs, fn func(conn net.Conn)) (err error) {
	lis, err := kcp.ListenWithOptions(fmt.Sprintf("%s:%d", sc.ip, sc.port), config.Block, *config.DataShard, *config.ParityShard)
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
		sc.Listener = new(net.Listener)
		*sc.Listener = lis
		go func() {
			defer func() {
				if e := recover(); e != nil {
					log.Printf("ListenKCP crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
				}
			}()
			for {
				//var conn net.Conn
				conn, err := lis.AcceptKCP()
				if err == nil {
					go func() {
						defer func() {
							if e := recover(); e != nil {
								log.Printf("kcp connection handler crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
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
							cconn := NewCompStream(conn)
							fn(cconn)
						}
					}()
				} else {
					sc.errAcceptHandler(err)
					break
				}
			}
		}()
	}
	return
}

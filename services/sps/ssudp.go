package sps

import (
	"bytes"
	"fmt"
	"net"
	"runtime/debug"
	"time"

	"github.com/snail007/goproxy/utils"
	goaes "github.com/snail007/goproxy/utils/aes"
	"github.com/snail007/goproxy/utils/socks"
)

func (s *SPS) RunSSUDP(addr string) (err error) {
	a, _ := net.ResolveUDPAddr("udp", addr)
	listener, err := net.ListenUDP("udp", a)
	if err != nil {
		s.log.Printf("ss udp bind error %s", err)
		return
	}
	s.log.Printf("ss udp on %s", listener.LocalAddr())
	s.udpRelatedPacketConns.Set(addr, listener)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				s.log.Printf("udp local->out io copy crashed:\n%s\n%s", e, string(debug.Stack()))
			}
		}()
		buf := utils.LeakyBuffer.Get()
		defer utils.LeakyBuffer.Put(buf)
		for {
			n, srcAddr, err := listener.ReadFrom(buf)
			if err != nil {
				s.log.Printf("read from client error %s", err)
				if utils.IsNetClosedErr(err) {
					return
				}
				continue
			}
			var (
				inconnRemoteAddr = srcAddr.String()
				outUDPConn       *net.UDPConn
				outconn          net.Conn
				outconnLocalAddr string
				destAddr         *net.UDPAddr
				clean            = func(msg, err string) {
					raddr := ""
					if outUDPConn != nil {
						raddr = outUDPConn.RemoteAddr().String()
						outUDPConn.Close()
					}
					if msg != "" {
						if raddr != "" {
							s.log.Printf("%s , %s , %s -> %s", msg, err, inconnRemoteAddr, raddr)
						} else {
							s.log.Printf("%s , %s , from : %s", msg, err, inconnRemoteAddr)
						}
					}
					s.userConns.Remove(inconnRemoteAddr)
					if outconn != nil {
						outconn.Close()
					}
					if outconnLocalAddr != "" {
						s.userConns.Remove(outconnLocalAddr)
					}
				}
			)
			defer clean("", "")

			raw := new(bytes.Buffer)
			raw.Write([]byte{0x00, 0x00, 0x00})
			raw.Write(s.localCipher.Decrypt(buf[:n]))
			socksPacket := socks.NewPacketUDP()
			err = socksPacket.Parse(raw.Bytes())
			raw = nil
			if err != nil {
				s.log.Printf("udp parse error %s", err)
				return
			}

			if v, ok := s.udpRelatedPacketConns.Get(inconnRemoteAddr); !ok {
				//socks client
				lbAddr := s.lb.Select(inconnRemoteAddr, *s.cfg.LoadBalanceOnlyHA)
				outconn, err := s.GetParentConn(lbAddr)
				if err != nil {
					clean("connnect fail", fmt.Sprintf("%s", err))
					return
				}

				client, err := s.HandshakeSocksParent(&outconn, "udp", socksPacket.Addr(), socks.Auth{}, true)
				if err != nil {
					clean("handshake fail", fmt.Sprintf("%s", err))
					return
				}

				outconnLocalAddr = outconn.LocalAddr().String()
				s.userConns.Set(outconnLocalAddr, &outconn)
				go func() {
					defer func() {
						if e := recover(); e != nil {
							s.log.Printf("udp related parent tcp conn read crashed:\n%s\n%s", e, string(debug.Stack()))
						}
					}()
					buf := make([]byte, 1)
					outconn.SetReadDeadline(time.Time{})
					if _, err := outconn.Read(buf); err != nil {
						clean("udp parent tcp conn disconnected", fmt.Sprintf("%s", err))
					}
				}()
				destAddr, _ = net.ResolveUDPAddr("udp", client.UDPAddr)
				localZeroAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
				outUDPConn, err = net.DialUDP("udp", localZeroAddr, destAddr)
				if err != nil {
					s.log.Printf("create out udp conn fail , %s , from : %s", err, srcAddr)
					return
				}
				s.udpRelatedPacketConns.Set(srcAddr.String(), outUDPConn)
				utils.UDPCopy(listener, outUDPConn, srcAddr, time.Second*5, func(data []byte) []byte {
					//forward to local
					var v []byte
					//convert parent data to raw
					if len(s.udpParentKey) > 0 {
						v, err = goaes.Decrypt(s.udpParentKey, data)
						if err != nil {
							s.log.Printf("udp outconn parse packet fail, %s", err.Error())
							return []byte{}
						}
					} else {
						v = data
					}
					return s.localCipher.Encrypt(v[3:])
				}, func(err interface{}) {
					s.udpRelatedPacketConns.Remove(srcAddr.String())
					if err != nil {
						s.log.Printf("udp out->local io copy crashed:\n%s\n%s", err, string(debug.Stack()))
					}
				})
			} else {
				outUDPConn = v.(*net.UDPConn)
			}
			//forward to parent
			//p is raw, now convert it to parent
			var v []byte
			if len(s.udpParentKey) > 0 {
				v, _ = goaes.Encrypt(s.udpParentKey, socksPacket.Bytes())
			} else {
				v = socksPacket.Bytes()
			}
			_, err = outUDPConn.Write(v)
			socksPacket = socks.PacketUDP{}
			if err != nil {
				if utils.IsNetClosedErr(err) {
					return
				}
				s.log.Printf("send out udp data fail , %s , from : %s", err, srcAddr)
			}
		}
	}()
	return
}

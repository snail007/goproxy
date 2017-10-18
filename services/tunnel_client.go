package services

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"proxy/utils"
	"strings"
	"time"
)

type TunnelClient struct {
	cfg TunnelClientArgs
}

func NewTunnelClient() Service {
	return &TunnelClient{
		cfg: TunnelClientArgs{},
	}
}

func (s *TunnelClient) InitService() {
}
func (s *TunnelClient) CheckArgs() {
	if *s.cfg.Parent != "" {
		log.Printf("use tls parent %s", *s.cfg.Parent)
	} else {
		log.Fatalf("parent required")
	}
	if *s.cfg.CertFile == "" || *s.cfg.KeyFile == "" {
		log.Fatalf("cert and key file required")
	}
	s.cfg.CertBytes, s.cfg.KeyBytes = utils.TlsBytes(*s.cfg.CertFile, *s.cfg.KeyFile)
}
func (s *TunnelClient) StopService() {
}
func (s *TunnelClient) Start(args interface{}) (err error) {
	s.cfg = args.(TunnelClientArgs)
	s.CheckArgs()
	s.InitService()
	log.Printf("proxy on tunnel client mode")
	for {
		ctrlConn, err := s.GetInConn(CONN_CONTROL, "")
		if err != nil {
			log.Printf("control connection err: %s, retrying...", err)
			time.Sleep(time.Second * 3)
			utils.CloseConn(&ctrlConn)
			continue
		}
		go func() {
			for {
				ctrlConn.SetWriteDeadline(time.Now().Add(time.Second * 3))
				_, err = ctrlConn.Write([]byte{0x00})
				ctrlConn.SetWriteDeadline(time.Time{})
				if err != nil {
					utils.CloseConn(&ctrlConn)
					log.Printf("ctrlConn err %s", err)
					break
				}
				time.Sleep(time.Second * 3)
			}
		}()
		for {
			signal := make([]byte, 50)
			n, err := ctrlConn.Read(signal)
			if err != nil {
				utils.CloseConn(&ctrlConn)
				log.Printf("read connection signal err: %s, retrying...", err)
				break
			}
			addr := string(signal[:n])
			log.Printf("signal revecived:%s", addr)
			protocol := addr[:3]
			atIndex := strings.Index(addr, "@")
			ID := addr[atIndex+1:]
			localAddr := addr[4:atIndex]
			if protocol == "udp" {
				go s.ServeUDP(localAddr, ID)
			} else {
				go s.ServeConn(localAddr, ID)
			}
		}
	}
}
func (s *TunnelClient) Clean() {
	s.StopService()
}
func (s *TunnelClient) GetInConn(typ uint8, ID string) (outConn net.Conn, err error) {
	outConn, err = s.GetConn()
	if err != nil {
		err = fmt.Errorf("connection err: %s", err)
		return
	}
	keyBytes := []byte(*s.cfg.Key)
	keyLength := uint16(len(keyBytes))
	pkg := new(bytes.Buffer)
	binary.Write(pkg, binary.LittleEndian, typ)
	binary.Write(pkg, binary.LittleEndian, keyLength)
	binary.Write(pkg, binary.LittleEndian, keyBytes)
	if ID != "" {
		IDBytes := []byte(ID)
		IDLength := uint16(len(IDBytes))
		binary.Write(pkg, binary.LittleEndian, IDLength)
		binary.Write(pkg, binary.LittleEndian, IDBytes)
	}
	_, err = outConn.Write(pkg.Bytes())
	if err != nil {
		err = fmt.Errorf("write connection data err: %s ,retrying...", err)
		utils.CloseConn(&outConn)
		return
	}
	return
}
func (s *TunnelClient) GetConn() (conn net.Conn, err error) {
	var _conn tls.Conn
	_conn, err = utils.TlsConnectHost(*s.cfg.Parent, *s.cfg.Timeout, s.cfg.CertBytes, s.cfg.KeyBytes)
	if err == nil {
		conn = net.Conn(&_conn)
	}
	return
}
func (s *TunnelClient) ServeUDP(localAddr, ID string) {
	var inConn net.Conn
	var err error
	// for {
	for {
		inConn, err = s.GetInConn(CONN_CLIENT, ID)
		if err != nil {
			utils.CloseConn(&inConn)
			log.Printf("connection err: %s, retrying...", err)
			time.Sleep(time.Second * 3)
			continue
		} else {
			break
		}
	}
	log.Printf("conn %s created", ID)
	// hw := utils.NewHeartbeatReadWriter(&inConn, 3, func(err error, hw *utils.HeartbeatReadWriter) {
	// 	log.Printf("hw err %s", err)
	// 	hw.Close()
	// })
	for {
		// srcAddr, body, err := utils.ReadUDPPacket(&hw)
		srcAddr, body, err := utils.ReadUDPPacket(inConn)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			log.Printf("connection %s released", ID)
			utils.CloseConn(&inConn)
			break
		} else if err != nil {
			log.Printf("udp packet revecived fail, err: %s", err)
		} else {
			//log.Printf("udp packet revecived:%s,%v", srcAddr, body)
			go s.processUDPPacket(&inConn, srcAddr, localAddr, body)
		}

	}
	// }
}
func (s *TunnelClient) processUDPPacket(inConn *net.Conn, srcAddr, localAddr string, body []byte) {
	dstAddr, err := net.ResolveUDPAddr("udp", localAddr)
	if err != nil {
		log.Printf("can't resolve address: %s", err)
		utils.CloseConn(inConn)
		return
	}
	clientSrcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	conn, err := net.DialUDP("udp", clientSrcAddr, dstAddr)
	if err != nil {
		log.Printf("connect to udp %s fail,ERR:%s", dstAddr.String(), err)
		return
	}
	conn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(*s.cfg.Timeout)))
	_, err = conn.Write(body)
	if err != nil {
		log.Printf("send udp packet to %s fail,ERR:%s", dstAddr.String(), err)
		return
	}
	//log.Printf("send udp packet to %s success", dstAddr.String())
	buf := make([]byte, 1024)
	length, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		log.Printf("read udp response from %s fail ,ERR:%s", dstAddr.String(), err)
		return
	}
	respBody := buf[0:length]
	//log.Printf("revecived udp packet from %s , %v", dstAddr.String(), respBody)
	bs := utils.UDPPacket(srcAddr, respBody)
	_, err = (*inConn).Write(bs)
	if err != nil {
		log.Printf("send udp response fail ,ERR:%s", err)
		utils.CloseConn(inConn)
		return
	}
	//log.Printf("send udp response success ,from:%s ,%d ,%v", dstAddr.String(), len(bs), bs)
}
func (s *TunnelClient) ServeConn(localAddr, ID string) {
	var inConn, outConn net.Conn
	var err error
	for {
		inConn, err = s.GetInConn(CONN_CLIENT, ID)
		if err != nil {
			utils.CloseConn(&inConn)
			log.Printf("connection err: %s, retrying...", err)
			time.Sleep(time.Second * 3)
			continue
		} else {
			break
		}
	}

	i := 0
	for {
		i++
		outConn, err = utils.ConnectHost(localAddr, *s.cfg.Timeout)
		if err == nil || i == 3 {
			break
		} else {
			if i == 3 {
				log.Printf("connect to %s err: %s, retrying...", localAddr, err)
				time.Sleep(2 * time.Second)
				continue
			}
		}
	}
	if err != nil {
		utils.CloseConn(&inConn)
		utils.CloseConn(&outConn)
		log.Printf("build connection error, err: %s", err)
		return
	}
	utils.IoBind(inConn, outConn, func(err error) {
		log.Printf("conn %s released", ID)
		utils.CloseConn(&inConn)
		utils.CloseConn(&outConn)
	}, func(i int, b bool) {}, 0)
	log.Printf("conn %s created", ID)
}

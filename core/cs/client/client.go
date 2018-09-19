package client

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/snail007/goproxy/core/dst"
	"github.com/snail007/goproxy/core/lib/kcpcfg"
	compressconn "github.com/snail007/goproxy/core/lib/transport"
	encryptconn "github.com/snail007/goproxy/core/lib/transport/encrypt"
	kcp "github.com/xtaci/kcp-go"
)

func TlsConnectHost(host string, timeout int, certBytes, keyBytes, caCertBytes []byte) (conn tls.Conn, err error) {
	h := strings.Split(host, ":")
	port, _ := strconv.Atoi(h[1])
	return TlsConnect(h[0], port, timeout, certBytes, keyBytes, caCertBytes)
}

func TlsConnect(host string, port, timeout int, certBytes, keyBytes, caCertBytes []byte) (conn tls.Conn, err error) {
	conf, err := getRequestTlsConfig(certBytes, keyBytes, caCertBytes)
	if err != nil {
		return
	}
	_conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Duration(timeout)*time.Millisecond)
	if err != nil {
		return
	}
	return *tls.Client(_conn, conf), err
}
func TlsConfig(certBytes, keyBytes, caCertBytes []byte) (conf *tls.Config, err error) {
	return getRequestTlsConfig(certBytes, keyBytes, caCertBytes)
}
func getRequestTlsConfig(certBytes, keyBytes, caCertBytes []byte) (conf *tls.Config, err error) {

	var cert tls.Certificate
	cert, err = tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		return
	}
	serverCertPool := x509.NewCertPool()
	caBytes := certBytes
	if caCertBytes != nil {
		caBytes = caCertBytes

	}
	ok := serverCertPool.AppendCertsFromPEM(caBytes)
	if !ok {
		err = errors.New("failed to parse root certificate")
	}
	block, _ := pem.Decode(caBytes)
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	x509Cert, _ := x509.ParseCertificate(block.Bytes)
	if x509Cert == nil {
		panic("failed to parse block")
	}
	conf = &tls.Config{
		RootCAs:            serverCertPool,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		ServerName:         x509Cert.Subject.CommonName,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			opts := x509.VerifyOptions{
				Roots: serverCertPool,
			}
			for _, rawCert := range rawCerts {
				cert, _ := x509.ParseCertificate(rawCert)
				_, err := cert.Verify(opts)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	return
}

func TCPConnectHost(hostAndPort string, timeout int) (conn net.Conn, err error) {
	conn, err = net.DialTimeout("tcp", hostAndPort, time.Duration(timeout)*time.Millisecond)
	return
}

func TCPSConnectHost(hostAndPort string, method, password string, compress bool, timeout int) (conn net.Conn, err error) {
	conn, err = net.DialTimeout("tcp", hostAndPort, time.Duration(timeout)*time.Millisecond)
	if err != nil {
		return
	}
	if compress {
		conn = compressconn.NewCompConn(conn)
	}
	conn, err = encryptconn.NewConn(conn, method, password)
	return
}

func TOUConnectHost(hostAndPort string, method, password string, compress bool, timeout int) (conn net.Conn, err error) {
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{})
	if err != nil {
		panic(err)
	}
	// Create a DST mux around the packet connection with the default max
	// packet size.
	mux := dst.NewMux(udpConn, 0)
	conn, err = mux.Dial("dst", hostAndPort)
	if compress {
		conn = compressconn.NewCompConn(conn)
	}
	conn, err = encryptconn.NewConn(conn, method, password)
	return
}
func KCPConnectHost(hostAndPort string, config kcpcfg.KCPConfigArgs) (conn net.Conn, err error) {
	kcpconn, err := kcp.DialWithOptions(hostAndPort, config.Block, *config.DataShard, *config.ParityShard)
	if err != nil {
		return
	}
	kcpconn.SetStreamMode(true)
	kcpconn.SetWriteDelay(true)
	kcpconn.SetNoDelay(*config.NoDelay, *config.Interval, *config.Resend, *config.NoCongestion)
	kcpconn.SetMtu(*config.MTU)
	kcpconn.SetWindowSize(*config.SndWnd, *config.RcvWnd)
	kcpconn.SetACKNoDelay(*config.AckNodelay)
	if *config.NoComp {
		return kcpconn, err
	}
	return compressconn.NewCompStream(kcpconn), err
}

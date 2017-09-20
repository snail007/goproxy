package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func IoBind(dst io.ReadWriter, src io.ReadWriter, fn func(err error), cfn func(count int, isPositive bool), bytesPreSec float64) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				log.Printf("IoBind crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
			}
		}()
		errchn := make(chan error, 2)
		go func() {
			defer func() {
				if e := recover(); e != nil {
					log.Printf("IoBind crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
				}
			}()
			var err error
			if bytesPreSec > 0 {
				newreader := NewReader(src)
				newreader.SetRateLimit(bytesPreSec)
				_, err = ioCopy(dst, newreader, func(c int) {
					cfn(c, false)
				})

			} else {
				_, err = ioCopy(dst, src, func(c int) {
					cfn(c, false)
				})
			}
			errchn <- err
		}()
		go func() {
			defer func() {
				if e := recover(); e != nil {
					log.Printf("IoBind crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
				}
			}()
			var err error
			if bytesPreSec > 0 {
				newReader := NewReader(dst)
				newReader.SetRateLimit(bytesPreSec)
				_, err = ioCopy(src, newReader, func(c int) {
					cfn(c, true)
				})
			} else {
				_, err = ioCopy(src, dst, func(c int) {
					cfn(c, true)
				})
			}
			errchn <- err
		}()
		fn(<-errchn)
	}()
}
func ioCopy(dst io.Writer, src io.Reader, fn ...func(count int)) (written int64, err error) {
	buf := make([]byte, 32*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
				if len(fn) == 1 {
					fn[0](nw)
				}
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}
func TlsConnectHost(host string, timeout int, certBytes, keyBytes []byte) (conn tls.Conn, err error) {
	h := strings.Split(host, ":")
	port, _ := strconv.Atoi(h[1])
	return TlsConnect(h[0], port, timeout, certBytes, keyBytes)
}

func TlsConnect(host string, port, timeout int, certBytes, keyBytes []byte) (conn tls.Conn, err error) {
	conf, err := getRequestTlsConfig(certBytes, keyBytes)
	if err != nil {
		return
	}
	_conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Duration(timeout)*time.Millisecond)
	if err != nil {
		return
	}
	return *tls.Client(_conn, conf), err
}
func getRequestTlsConfig(certBytes, keyBytes []byte) (conf *tls.Config, err error) {
	var cert tls.Certificate
	cert, err = tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		return
	}
	serverCertPool := x509.NewCertPool()
	ok := serverCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		err = errors.New("failed to parse root certificate")
	}
	conf = &tls.Config{
		RootCAs:            serverCertPool,
		Certificates:       []tls.Certificate{cert},
		ServerName:         "proxy",
		InsecureSkipVerify: false,
	}
	return
}

func ConnectHost(hostAndPort string, timeout int) (conn net.Conn, err error) {
	conn, err = net.DialTimeout("tcp", hostAndPort, time.Duration(timeout)*time.Millisecond)
	return
}
func ListenTls(ip string, port int, certBytes, keyBytes []byte) (ln *net.Listener, err error) {
	var cert tls.Certificate
	cert, err = tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		return
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		err = errors.New("failed to parse root certificate")
	}
	config := &tls.Config{
		ClientCAs:    clientCertPool,
		ServerName:   "proxy",
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	_ln, err := tls.Listen("tcp", fmt.Sprintf("%s:%d", ip, port), config)
	if err == nil {
		ln = &_ln
	}
	return
}
func PathExists(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
func HTTPGet(URL string, timeout int) (err error) {
	tr := &http.Transport{}
	var resp *http.Response
	var client *http.Client
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		tr.CloseIdleConnections()
	}()
	client = &http.Client{Timeout: time.Millisecond * time.Duration(timeout), Transport: tr}
	resp, err = client.Get(URL)
	if err != nil {
		return
	}
	return
}

func initOutPool(isTLS bool, certBytes, keyBytes []byte, address string, timeout int, InitialCap int, MaxCap int) {
	var err error
	outPool, err = NewConnPool(poolConfig{
		IsActive: func(conn interface{}) bool { return true },
		Release: func(conn interface{}) {
			if conn != nil {
				conn.(net.Conn).SetDeadline(time.Now().Add(time.Millisecond))
				conn.(net.Conn).Close()
				// log.Println("conn released")
			}
		},
		InitialCap: InitialCap,
		MaxCap:     MaxCap,
		Factory: func() (conn interface{}, err error) {
			conn, err = getConn(isTLS, certBytes, keyBytes, address, timeout)
			return
		},
	})
	if err != nil {
		log.Fatalf("init conn pool fail ,%s", err)
	} else {
		log.Printf("init conn pool success")
		initPoolDeamon(isTLS, certBytes, keyBytes, address, timeout)
	}
}
func getConn(isTLS bool, certBytes, keyBytes []byte, address string, timeout int) (conn interface{}, err error) {
	if isTLS {
		var _conn tls.Conn
		_conn, err = TlsConnectHost(address, timeout, certBytes, keyBytes)
		if err == nil {
			conn = net.Conn(&_conn)
		}
	} else {
		conn, err = ConnectHost(address, timeout)
	}
	return
}
func initPoolDeamon(isTLS bool, certBytes, keyBytes []byte, address string, timeout int) {
	go func() {
		dur := cfg.GetInt("check-proxy-interval")
		if dur <= 0 {
			return
		}
		log.Printf("pool deamon started")
		for {
			time.Sleep(time.Second * time.Duration(dur))
			conn, err := getConn(isTLS, certBytes, keyBytes, address, timeout)
			if err != nil {
				log.Printf("pool deamon err %s , release pool", err)
				outPool.ReleaseAll()
			} else {
				conn.(net.Conn).SetDeadline(time.Now().Add(time.Millisecond))
				conn.(net.Conn).Close()
			}
		}
	}()
}
func IsBasicAuth() bool {
	return cfg.GetString("auth-file") != "" || len(cfg.GetStringSlice("auth")) > 0
}
func InitBasicAuth() (err error) {
	basicAuth = NewBasicAuth()
	if cfg.GetString("auth-file") != "" {
		n, err := basicAuth.AddFromFile(cfg.GetString("auth-file"))
		if err != nil {
			return fmt.Errorf("auth-file ERR:%s", err)
		}
		log.Printf("auth data added from file %d , total:%d", n, basicAuth.Total())
	}
	if len(cfg.GetStringSlice("auth")) > 0 {
		n := basicAuth.Add(cfg.GetStringSlice("auth"))
		log.Printf("auth data added %d, total:%d", n, basicAuth.Total())
	}
	return
}
func clean() {
	//block main()
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		for _ = range signalChan {
			if outPool != nil {
				fmt.Println("\nReceived an interrupt, stopping services...")
				outPool.ReleaseAll()
				//time.Sleep(time.Second * 10)
				// fmt.Println("done")
			}
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
func HTTPProxyDecoder(inConn *net.Conn) (useProxy bool, request *HTTPRequest, err error) {
	var req HTTPRequest
	req, err = NewHTTPRequest(inConn, 4096)
	if err != nil {
		//log.Printf("NewHTTPRequest ERR:%s", err)
		return
	}
	useProxy = false
	if checker.data != nil {
		useProxy, _, _ = checker.IsBlocked(req.Host)
	}
	request = &req
	return
}

// type sockaddr struct {
// 	family uint16
// 	data   [14]byte
// }

// const SO_ORIGINAL_DST = 80

// realServerAddress returns an intercepted connection's original destination.
// func realServerAddress(conn *net.Conn) (string, error) {
// 	tcpConn, ok := (*conn).(*net.TCPConn)
// 	if !ok {
// 		return "", errors.New("not a TCPConn")
// 	}

// 	file, err := tcpConn.File()
// 	if err != nil {
// 		return "", err
// 	}

// 	// To avoid potential problems from making the socket non-blocking.
// 	tcpConn.Close()
// 	*conn, err = net.FileConn(file)
// 	if err != nil {
// 		return "", err
// 	}

// 	defer file.Close()
// 	fd := file.Fd()

// 	var addr sockaddr
// 	size := uint32(unsafe.Sizeof(addr))
// 	err = getsockopt(int(fd), syscall.SOL_IP, SO_ORIGINAL_DST, uintptr(unsafe.Pointer(&addr)), &size)
// 	if err != nil {
// 		return "", err
// 	}

// 	var ip net.IP
// 	switch addr.family {
// 	case syscall.AF_INET:
// 		ip = addr.data[2:6]
// 	default:
// 		return "", errors.New("unrecognized address family")
// 	}

// 	port := int(addr.data[0])<<8 + int(addr.data[1])

// 	return net.JoinHostPort(ip.String(), strconv.Itoa(port)), nil
// }

// func getsockopt(s int, level int, name int, val uintptr, vallen *uint32) (err error) {
// 	_, _, e1 := syscall.Syscall6(syscall.SYS_GETSOCKOPT, uintptr(s), uintptr(level), uintptr(name), uintptr(val), uintptr(unsafe.Pointer(vallen)), 0)
// 	if e1 != 0 {
// 		err = e1
// 	}
// 	return
// }

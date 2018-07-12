package utils

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	logger "log"
	"math/rand"
	"net"
	"net/http"
	"os"

	"github.com/snail007/goproxy/services/kcpcfg"

	"golang.org/x/crypto/pbkdf2"

	"strconv"
	"strings"
	"time"
	"context"

	"github.com/snail007/goproxy/utils/id"

	kcp "github.com/xtaci/kcp-go"
	
)

func IoBind(dst io.ReadWriteCloser, src io.ReadWriteCloser, fn func(err interface{}), log *logger.Logger) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("bind crashed %s", err)
			}
		}()
		e1 := make(chan interface{}, 1)
		e2 := make(chan interface{}, 1)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("bind crashed %s", err)
				}
			}()
			//_, err := io.Copy(dst, src)
			err := ioCopy(dst, src)
			e1 <- err
		}()
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("bind crashed %s", err)
				}
			}()
			//_, err := io.Copy(src, dst)
			err := ioCopy(src, dst)
			e2 <- err
		}()
		var err interface{}
		select {
		case err = <-e1:
			//log.Printf("e1")
		case err = <-e2:
			//log.Printf("e2")
		}
		src.Close()
		dst.Close()
		if fn != nil {
			fn(err)
		}
	}()
}
func ioCopy(dst io.ReadWriter, src io.ReadWriter) (err error) {
	buf := LeakyBuffer.Get()
	defer LeakyBuffer.Put(buf)
	n := 0
	for {
		n, err = src.Read(buf)
		if n > 0 {
			if _, e := dst.Write(buf[0:n]); e != nil {
				return e
			}
		}
		if err != nil {
			return
		}
	}
}
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

func ConnectHost(hostAndPort string, timeout int) (conn net.Conn, err error) {
	conn, err = net.DialTimeout("tcp", hostAndPort, time.Duration(timeout)*time.Millisecond)
	return
}
func ConnectKCPHost(hostAndPort string, config kcpcfg.KCPConfigArgs) (conn net.Conn, err error) {
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
	return NewCompStream(kcpconn), err
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

func CloseConn(conn *net.Conn) {
	defer func() {
		_ = recover()
	}()
	if conn != nil && *conn != nil {
		(*conn).SetDeadline(time.Now().Add(time.Millisecond))
		(*conn).Close()
	}
}
func GetAllInterfaceAddr() ([]net.IP, error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	addresses := []net.IP{}
	for _, iface := range ifaces {

		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		// if iface.Flags&net.FlagLoopback != 0 {
		// 	continue // loopback interface
		// }
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// if ip == nil || ip.IsLoopback() {
			// 	continue
			// }
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			addresses = append(addresses, ip)
		}
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("no address Found, net.InterfaceAddrs: %v", addresses)
	}
	//only need first
	return addresses, nil
}
func UDPPacket(srcAddr string, packet []byte) []byte {
	addrBytes := []byte(srcAddr)
	addrLength := uint16(len(addrBytes))
	bodyLength := uint16(len(packet))
	//log.Printf("build packet : addr len %d, body len %d", addrLength, bodyLength)
	pkg := new(bytes.Buffer)
	binary.Write(pkg, binary.LittleEndian, addrLength)
	binary.Write(pkg, binary.LittleEndian, addrBytes)
	binary.Write(pkg, binary.LittleEndian, bodyLength)
	binary.Write(pkg, binary.LittleEndian, packet)
	return pkg.Bytes()
}
func ReadUDPPacket(_reader io.Reader) (srcAddr string, packet []byte, err error) {
	reader := bufio.NewReader(_reader)
	var addrLength uint16
	var bodyLength uint16
	err = binary.Read(reader, binary.LittleEndian, &addrLength)
	if err != nil {
		return
	}
	_srcAddr := make([]byte, addrLength)
	n, err := reader.Read(_srcAddr)
	if err != nil {
		return
	}
	if n != int(addrLength) {
		err = fmt.Errorf("n != int(addrLength), %d,%d", n, addrLength)
		return
	}
	srcAddr = string(_srcAddr)

	err = binary.Read(reader, binary.LittleEndian, &bodyLength)
	if err != nil {

		return
	}
	packet = make([]byte, bodyLength)
	n, err = reader.Read(packet)
	if err != nil {
		return
	}
	if n != int(bodyLength) {
		err = fmt.Errorf("n != int(bodyLength), %d,%d", n, bodyLength)
		return
	}
	return
}
func Uniqueid() string {
	return xid.New().String()
	// var src = rand.NewSource(time.Now().UnixNano())
	// s := fmt.Sprintf("%d", src.Int63())
	// return s[len(s)-5:len(s)-1] + fmt.Sprintf("%d", uint64(time.Now().UnixNano()))[8:]
}
func RandString(strlen int) string {
	codes := "QWERTYUIOPLKJHGFDSAZXCVBNMabcdefghijklmnopqrstuvwxyz0123456789"
	codeLen := len(codes)
	data := make([]byte, strlen)
	rand.Seed(time.Now().UnixNano() + rand.Int63() + rand.Int63() + rand.Int63() + rand.Int63())
	for i := 0; i < strlen; i++ {
		idx := rand.Intn(codeLen)
		data[i] = byte(codes[idx])
	}
	return string(data)
}
func RandInt(strLen int) int64 {
	codes := "123456789"
	codeLen := len(codes)
	data := make([]byte, strLen)
	rand.Seed(time.Now().UnixNano() + rand.Int63() + rand.Int63() + rand.Int63() + rand.Int63())
	for i := 0; i < strLen; i++ {
		idx := rand.Intn(codeLen)
		data[i] = byte(codes[idx])
	}
	i, _ := strconv.ParseInt(string(data), 10, 64)
	return i
}
func ReadData(r io.Reader) (data string, err error) {
	var len uint16
	err = binary.Read(r, binary.LittleEndian, &len)
	if err != nil {
		return
	}
	var n int
	_data := make([]byte, len)
	n, err = r.Read(_data)
	if err != nil {
		return
	}
	if n != int(len) {
		err = fmt.Errorf("error data len")
		return
	}
	data = string(_data)
	return
}
func ReadPacketData(r io.Reader, data ...*string) (err error) {
	for _, d := range data {
		*d, err = ReadData(r)
		if err != nil {
			return
		}
	}
	return
}
func ReadPacket(r io.Reader, typ *uint8, data ...*string) (err error) {
	var connType uint8
	err = binary.Read(r, binary.LittleEndian, &connType)
	if err != nil {
		return
	}
	*typ = connType
	for _, d := range data {
		*d, err = ReadData(r)
		if err != nil {
			return
		}
	}
	return
}
func BuildPacket(typ uint8, data ...string) []byte {
	pkg := new(bytes.Buffer)
	binary.Write(pkg, binary.LittleEndian, typ)
	for _, d := range data {
		bytes := []byte(d)
		binary.Write(pkg, binary.LittleEndian, uint16(len(bytes)))
		binary.Write(pkg, binary.LittleEndian, bytes)
	}
	return pkg.Bytes()
}
func BuildPacketData(data ...string) []byte {
	pkg := new(bytes.Buffer)
	for _, d := range data {
		bytes := []byte(d)
		binary.Write(pkg, binary.LittleEndian, uint16(len(bytes)))
		binary.Write(pkg, binary.LittleEndian, bytes)
	}
	return pkg.Bytes()
}
func SubStr(str string, start, end int) string {
	if len(str) == 0 {
		return ""
	}
	if end >= len(str) {
		end = len(str) - 1
	}
	return str[start:end]
}
func SubBytes(bytes []byte, start, end int) []byte {
	if len(bytes) == 0 {
		return []byte{}
	}
	if end >= len(bytes) {
		end = len(bytes) - 1
	}
	return bytes[start:end]
}
func TlsBytes(cert, key string) (certBytes, keyBytes []byte, err error) {
	certBytes, err = ioutil.ReadFile(cert)
	if err != nil {
		err = fmt.Errorf("err : %s", err)
		return
	}
	keyBytes, err = ioutil.ReadFile(key)
	if err != nil {
		err = fmt.Errorf("err : %s", err)
		return
	}
	return
}
func GetKCPBlock(method, key string) (block kcp.BlockCrypt) {
	pass := pbkdf2.Key([]byte(key), []byte(key), 4096, 32, sha1.New)
	switch method {
	case "sm4":
		block, _ = kcp.NewSM4BlockCrypt(pass[:16])
	case "tea":
		block, _ = kcp.NewTEABlockCrypt(pass[:16])
	case "xor":
		block, _ = kcp.NewSimpleXORBlockCrypt(pass)
	case "none":
		block, _ = kcp.NewNoneBlockCrypt(pass)
	case "aes-128":
		block, _ = kcp.NewAESBlockCrypt(pass[:16])
	case "aes-192":
		block, _ = kcp.NewAESBlockCrypt(pass[:24])
	case "blowfish":
		block, _ = kcp.NewBlowfishBlockCrypt(pass)
	case "twofish":
		block, _ = kcp.NewTwofishBlockCrypt(pass)
	case "cast5":
		block, _ = kcp.NewCast5BlockCrypt(pass[:16])
	case "3des":
		block, _ = kcp.NewTripleDESBlockCrypt(pass[:24])
	case "xtea":
		block, _ = kcp.NewXTEABlockCrypt(pass[:16])
	case "salsa20":
		block, _ = kcp.NewSalsa20BlockCrypt(pass)
	default:
		block, _ = kcp.NewAESBlockCrypt(pass)
	}
	return
}
func HttpGet(URL string, timeout int, host ...string) (body []byte, code int, err error) {
	var tr *http.Transport
	var client *http.Client
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	if strings.Contains(URL, "https://") {
		tr = &http.Transport{TLSClientConfig: conf}
		client = &http.Client{Timeout: time.Millisecond * time.Duration(timeout), Transport: tr}
	} else {
		tr = &http.Transport{}
		client = &http.Client{Timeout: time.Millisecond * time.Duration(timeout), Transport: tr}
	}
	defer tr.CloseIdleConnections()

	//resp, err := client.Get(URL)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return
	}
	if len(host) == 1 && host[0] != "" {
		req.Host = host[0]
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	code = resp.StatusCode
	body, err = ioutil.ReadAll(resp.Body)
	return
}
func IsIternalIP(domainOrIP string, always bool) bool {
	var outIPs []net.IP
	var err error
	var isDomain bool
	if net.ParseIP(domainOrIP) == nil {
		isDomain = true
	}
	if always && isDomain {
		return false
	}

	if isDomain {
		outIPs, err = MyLookupIP(domainOrIP)
	} else {
		outIPs = []net.IP{net.ParseIP(domainOrIP)}
	}

	if err != nil {
		return false
	}
	for _, ip := range outIPs {
		if ip.IsLoopback() {
			return true
		}
		if ip.To4().Mask(net.IPv4Mask(255, 0, 0, 0)).String() == "10.0.0.0" {
			return true
		}
		if ip.To4().Mask(net.IPv4Mask(255, 0, 0, 0)).String() == "192.168.0.0" {
			return true
		}
		if ip.To4().Mask(net.IPv4Mask(255, 0, 0, 0)).String() == "172.0.0.0" {
			i, _ := strconv.Atoi(strings.Split(ip.To4().String(), ".")[1])
			return i >= 16 && i <= 31
		}
	}
	return false
}
func IsHTTP(head []byte) bool {
	keys := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}
	for _, key := range keys {
		if bytes.HasPrefix(head, []byte(key)) || bytes.HasPrefix(head, []byte(strings.ToLower(key))) {
			return true
		}
	}
	return false
}
func IsSocks5(head []byte) bool {
	if len(head) < 3 {
		return false
	}
	if head[0] == uint8(0x05) && 0 < int(head[1]) && int(head[1]) < 255 {
		if len(head) == 2+int(head[1]) {
			return true
		}
	}
	return false
}
func RemoveProxyHeaders(head []byte) []byte {
	newLines := [][]byte{}
	var keys = map[string]bool{}
	lines := bytes.Split(head, []byte("\r\n"))
	IsBody := false
	i := -1
	for _, line := range lines {
		i++
		if len(line) == 0 || IsBody {
			newLines = append(newLines, line)
			IsBody = true
		} else {
			hline := bytes.SplitN(line, []byte(":"), 2)
			if i == 0 && IsHTTP(head) {
				newLines = append(newLines, line)
				continue
			}
			if len(hline) != 2 {
				continue
			}
			k := strings.ToUpper(string(hline[0]))
			if _, ok := keys[k]; ok || strings.HasPrefix(k, "PROXY-") {
				continue
			}
			keys[k] = true
			newLines = append(newLines, line)
		}
	}
	return bytes.Join(newLines, []byte("\r\n"))
}
func InsertProxyHeaders(head []byte, headers string) []byte {
	return bytes.Replace(head, []byte("\r\n"), []byte("\r\n"+headers), 1)
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


/*
net.LookupIP may cause  deadlock in windows
https://github.com/golang/go/issues/24178
*/
func MyLookupIP(host string) ([]net.IP, error) {

	ctx ,cancel := context.WithTimeout(context.Background(),time.Second *time.Duration(3))
	defer func() {
		cancel()
		//ctx.Done()
	}()
	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}
	ips := make([]net.IP, len(addrs))
	for i, ia := range addrs {
		ips[i] = ia.IP
	}
	return ips, nil
}

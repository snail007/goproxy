package utils

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"

	"golang.org/x/crypto/pbkdf2"

	"runtime/debug"
	"strconv"
	"strings"
	"time"

	kcp "github.com/xtaci/kcp-go"
)

func IoBind(dst io.ReadWriter, src io.ReadWriter, fn func(err error)) {
	go func() {
		e1 := make(chan error, 1)
		e2 := make(chan error, 1)
		go func() {
			defer func() {
				if e := recover(); e != nil {
					log.Printf("IoBind crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
				}
			}()

			_, e := io.Copy(dst, src)
			e1 <- e
		}()
		go func() {
			defer func() {
				if e := recover(); e != nil {
					log.Printf("IoBind crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
				}
			}()

			_, e := io.Copy(src, dst)
			e2 <- e
		}()
		var err error
		select {
		case err = <-e1:
			//log.Printf("e1")
		case err = <-e2:
			//log.Printf("e2")
		}
		fn(err)
	}()
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
func ConnectKCPHost(hostAndPort, method, key string) (conn net.Conn, err error) {
	conn, err = kcp.DialWithOptions(hostAndPort, GetKCPBlock(method, key), 10, 3)
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

func CloseConn(conn *net.Conn) {
	defer func() {
		_ = recover()
	}()
	if conn != nil && *conn != nil {
		(*conn).SetDeadline(time.Now().Add(time.Millisecond))
		(*conn).Close()
	}
}
func Keygen() (err error) {
	cmd := exec.Command("sh", "-c", "openssl genrsa -out proxy.key 2048")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("err:%s", err)
		return
	}
	fmt.Println(string(out))
	cmd = exec.Command("sh", "-c", `openssl req -new -key proxy.key -x509 -days 3650 -out proxy.crt -subj /C=CN/ST=BJ/O="Localhost Ltd"/CN=proxy`)
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Printf("err:%s", err)
		return
	}
	fmt.Println(string(out))
	return
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
	var src = rand.NewSource(time.Now().UnixNano())
	s := fmt.Sprintf("%d", src.Int63())
	return s[len(s)-5:len(s)-1] + fmt.Sprintf("%d", uint64(time.Now().UnixNano()))[8:]
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
func TlsBytes(cert, key string) (certBytes, keyBytes []byte) {
	certBytes, err := ioutil.ReadFile(cert)
	if err != nil {
		log.Fatalf("err : %s", err)
		return
	}
	keyBytes, err = ioutil.ReadFile(key)
	if err != nil {
		log.Fatalf("err : %s", err)
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

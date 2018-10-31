package jumper

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

type Jumper struct {
	proxyURL *url.URL
	timeout  time.Duration
}
type socks5Dialer struct {
	timeout time.Duration
}

func (s socks5Dialer) Dial(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, s.timeout)
}

func New(proxyURL string, timeout time.Duration) (j Jumper, err error) {
	u, e := url.Parse(proxyURL)
	if e != nil {
		err = e
		return
	}
	j = Jumper{
		proxyURL: u,
		timeout:  timeout,
	}
	return
}
func (j *Jumper) Dial(address string, timeout time.Duration) (net.Conn, error) {
	switch j.proxyURL.Scheme {
	case "https":
		return j.dialHTTPS(address, timeout)
	case "socks5":
		return j.dialSOCKS5(address, timeout)
	default:
		return nil, fmt.Errorf("unkown scheme of %s", j.proxyURL.String())
	}
}
func (j *Jumper) dialHTTPS(address string, timeout time.Duration) (conn net.Conn, err error) {
	conn, err = net.DialTimeout("tcp", j.proxyURL.Host, timeout)
	if err != nil {
		return
	}
	pb := new(bytes.Buffer)
	pb.Write([]byte(fmt.Sprintf("CONNECT %s HTTP/1.1\r\n", address)))
	pb.WriteString(fmt.Sprintf("Host: %s\r\n", address))
	pb.WriteString(fmt.Sprintf("Proxy-Host: %s\r\n", address))
	pb.WriteString("Proxy-Connection: Keep-Alive\r\n")
	pb.WriteString("Connection: Keep-Alive\r\n")
	if j.proxyURL.User != nil {
		p, _ := j.proxyURL.User.Password()
		u := fmt.Sprintf("%s:%s", j.proxyURL.User.Username(), p)
		pb.Write([]byte(fmt.Sprintf("Proxy-Authorization: Basic %s\r\n", base64.StdEncoding.EncodeToString([]byte(u)))))
	}
	pb.Write([]byte("\r\n"))
	_, err = conn.Write(pb.Bytes())
	if err != nil {
		conn.Close()
		conn = nil
		err = fmt.Errorf("error connecting to proxy: %s", err)
		return
	}
	reply := make([]byte, 1024)
	conn.SetDeadline(time.Now().Add(timeout))
	n, e := conn.Read(reply)
	conn.SetDeadline(time.Time{})
	if e != nil {
		err = fmt.Errorf("error read reply from proxy: %s", e)
		conn.Close()
		conn = nil
		return
	}
	if bytes.Index(reply[:n], []byte("200")) == -1 {
		err = fmt.Errorf("error greeting to proxy, response: %s", string(reply[:n]))
		conn.Close()
		conn = nil
		return
	}
	return
}
func (j *Jumper) dialSOCKS5(address string, timeout time.Duration) (conn net.Conn, err error) {
	auth := &proxy.Auth{}
	if j.proxyURL.User != nil {
		auth.User = j.proxyURL.User.Username()
		auth.Password, _ = j.proxyURL.User.Password()
	} else {
		auth = nil
	}
	dialSocksProxy, e := proxy.SOCKS5("tcp", j.proxyURL.Host, auth, socks5Dialer{timeout: timeout})
	if e != nil {
		err = fmt.Errorf("error connecting to proxy: %s", e)
		return
	}
	return dialSocksProxy.Dial("tcp", address)
}

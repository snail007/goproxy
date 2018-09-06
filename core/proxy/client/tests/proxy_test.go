package tests

import (
	"io/ioutil"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	proxyclient "github.com/snail007/goproxy/core/proxy/client"
	sdk "github.com/snail007/goproxy/sdk/android-ios"
)

func TestSocks5(t *testing.T) {
	estr := sdk.Start("s1", "socks -p :8185 --log test.log")
	if estr != "" {
		t.Fatal(estr)
	}
	p, e := proxyclient.SOCKS5(time.Second, nil)
	if e != nil {
		t.Error(e)
	} else {
		c, e := net.Dial("tcp", "127.0.0.1:8185")
		if e != nil {
			t.Fatal(e)
		}
		e = p.DialConn(&c, "tcp", "www.baidu.com:80")
		if e != nil {
			t.Fatal(e)
		}
		_, e = c.Write([]byte("Get / http/1.1\r\nHost: www.baidu.com\r\n"))
		if e != nil {
			t.Fatal(e)
		}
		b, e := ioutil.ReadAll(c)
		if e != nil {
			t.Fatal(e)
		}
		if !strings.HasPrefix(string(b), "HTTP") {
			t.Fatalf("request baidu fail:%s", string(b))
		}
	}
	sdk.Stop("s1")
	os.Remove("test.log")
}

func TestSocks5Auth(t *testing.T) {
	estr := sdk.Start("s1", "socks -p :8185 -a u:p --log test.log")
	if estr != "" {
		t.Fatal(estr)
	}
	p, e := proxyclient.SOCKS5(time.Second, &proxyclient.Auth{User: "u", Password: "p"})
	if e != nil {
		t.Error(e)
	} else {
		c, e := net.Dial("tcp", "127.0.0.1:8185")
		if e != nil {
			t.Fatal(e)
		}
		e = p.DialConn(&c, "tcp", "www.baidu.com:80")
		if e != nil {
			t.Fatal(e)
		}
		_, e = c.Write([]byte("Get / http/1.1\r\nHost: www.baidu.com\r\n"))
		if e != nil {
			t.Fatal(e)
		}
		b, e := ioutil.ReadAll(c)
		if e != nil {
			t.Fatal(e)
		}
		if !strings.HasPrefix(string(b), "HTTP") {
			t.Fatalf("request baidu fail:%s", string(b))
		}
	}
	sdk.Stop("s1")
	os.Remove("test.log")
}

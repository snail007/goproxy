package tests

import (
	"log"
	"net"
	"os"
	"testing"

	ctransport "github.com/snail007/goproxy/core/cs/client"
	stransport "github.com/snail007/goproxy/core/cs/server"
)

func TestTCPS(t *testing.T) {
	l := log.New(os.Stderr, "", log.LstdFlags)
	s := stransport.NewServerChannelHost(":", l)
	err := s.ListenTCPS("aes-256-cfb", "password", true, func(inconn net.Conn) {
		buf := make([]byte, 2048)
		_, err := inconn.Read(buf)
		if err != nil {
			t.Error(err)
			return
		}
		_, err = inconn.Write([]byte("okay"))
		if err != nil {
			t.Error(err)
			return
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	client, err := ctransport.TCPSConnectHost((*s.Listener).Addr().String(), "aes-256-cfb", "password", true, 1000)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	_, err = client.Write([]byte("test"))
	if err != nil {
		t.Fatal(err)
	}
	b := make([]byte, 20)
	n, err := client.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	if string(b[:n]) != "okay" {
		t.Fatalf("client revecive okay excepted,revecived : %s", string(b[:n]))
	}
}

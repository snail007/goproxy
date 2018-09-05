// Package proxy provides support for a variety of protocols to proxy network
// data.
package client

import (
	"net"
	"time"

	socks5c "github.com/snail007/goproxy/core/lib/socks5"
	socks5 "github.com/snail007/goproxy/core/proxy/client/socks5"
)

// A Dialer is a means to establish a connection.
type Dialer interface {
	// Dial connects to the given address via the proxy.
	DialConn(conn *net.Conn, network, addr string) (err error)
}

// Auth contains authentication parameters that specific Dialers may require.
type Auth struct {
	User, Password string
}

func SOCKS5(timeout time.Duration, auth *Auth) (Dialer, error) {
	var a *socks5c.UsernamePassword
	if auth != nil {
		a = &socks5c.UsernamePassword{auth.User, auth.Password}
	}
	d := socks5.NewDialer(a, timeout)
	return d, nil
}

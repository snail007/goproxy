package services

import "golang.org/x/crypto/ssh"

// tcp := app.Command("tcp", "proxy on tcp mode")
// t := tcp.Flag("tcp-timeout", "tcp timeout milliseconds when connect to real server or parent proxy").Default("2000").Int()

const (
	TYPE_TCP     = "tcp"
	TYPE_UDP     = "udp"
	TYPE_HTTP    = "http"
	TYPE_TLS     = "tls"
	CONN_CONTROL = uint8(1)
	CONN_SERVER  = uint8(2)
	CONN_CLIENT  = uint8(3)
)

type TunnelServerArgs struct {
	Parent    *string
	CertFile  *string
	KeyFile   *string
	CertBytes []byte
	KeyBytes  []byte
	Local     *string
	IsUDP     *bool
	Key       *string
	Remote    *string
	Timeout   *int
	Route     *[]string
}
type TunnelClientArgs struct {
	Parent    *string
	CertFile  *string
	KeyFile   *string
	CertBytes []byte
	KeyBytes  []byte
	Key       *string
	Timeout   *int
}
type TunnelBridgeArgs struct {
	Parent    *string
	CertFile  *string
	KeyFile   *string
	CertBytes []byte
	KeyBytes  []byte
	Local     *string
	Timeout   *int
}
type TCPArgs struct {
	Parent              *string
	CertFile            *string
	KeyFile             *string
	CertBytes           []byte
	KeyBytes            []byte
	Local               *string
	ParentType          *string
	IsTLS               *bool
	Timeout             *int
	PoolSize            *int
	CheckParentInterval *int
}

type HTTPArgs struct {
	Parent              *string
	CertFile            *string
	KeyFile             *string
	CertBytes           []byte
	KeyBytes            []byte
	Local               *string
	Always              *bool
	HTTPTimeout         *int
	Interval            *int
	Blocked             *string
	Direct              *string
	AuthFile            *string
	Auth                *[]string
	ParentType          *string
	LocalType           *string
	Timeout             *int
	PoolSize            *int
	CheckParentInterval *int
	SSHKeyFile          *string
	SSHKeyFileSalt      *string
	SSHPassword         *string
	SSHUser             *string
	SSHKeyBytes         []byte
	SSHAuthMethod       ssh.AuthMethod
}
type UDPArgs struct {
	Parent              *string
	CertFile            *string
	KeyFile             *string
	CertBytes           []byte
	KeyBytes            []byte
	Local               *string
	ParentType          *string
	Timeout             *int
	PoolSize            *int
	CheckParentInterval *int
}
type SocksArgs struct {
	Parent         *string
	ParentType     *string
	Local          *string
	LocalType      *string
	CertFile       *string
	KeyFile        *string
	CertBytes      []byte
	KeyBytes       []byte
	SSHKeyFile     *string
	SSHKeyFileSalt *string
	SSHPassword    *string
	SSHUser        *string
	SSHKeyBytes    []byte
	SSHAuthMethod  ssh.AuthMethod
	Timeout        *int
	Always         *bool
	Interval       *int
	Blocked        *string
	Direct         *string
	AuthFile       *string
	Auth           *[]string
}

func (a *TCPArgs) Protocol() string {
	if *a.IsTLS {
		return "tls"
	}
	return "tcp"
}

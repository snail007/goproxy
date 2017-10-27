package services

import "golang.org/x/crypto/ssh"

// tcp := app.Command("tcp", "proxy on tcp mode")
// t := tcp.Flag("tcp-timeout", "tcp timeout milliseconds when connect to real server or parent proxy").Default("2000").Int()

const (
	TYPE_TCP             = "tcp"
	TYPE_UDP             = "udp"
	TYPE_HTTP            = "http"
	TYPE_TLS             = "tls"
	TYPE_KCP             = "kcp"
	CONN_CLIENT_CONTROL  = uint8(1)
	CONN_CLIENT_HEARBEAT = uint8(2)
	CONN_SERVER_HEARBEAT = uint8(3)
	CONN_SERVER          = uint8(4)
	CONN_CLIENT          = uint8(5)
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
	Mgr       *TunnelServerManager
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
	LocalType           *string
	Timeout             *int
	PoolSize            *int
	CheckParentInterval *int
	KCPMethod           *string
	KCPKey              *string
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
	KCPMethod           *string
	KCPKey              *string
	LocalIPS            *[]string
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
	KCPMethod      *string
	KCPKey         *string
	UDPParent      *string
	UDPLocal       *string
	LocalIPS       *[]string
}

func (a *TCPArgs) Protocol() string {
	switch *a.LocalType {
	case TYPE_TLS:
		return TYPE_TLS
	case TYPE_TCP:
		return TYPE_TCP
	case TYPE_KCP:
		return TYPE_KCP
	}
	return "unknown"
}

package services

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

type Args struct {
	Local     *string
	Parent    *string
	CertBytes []byte
	KeyBytes  []byte
}
type TunnelServerArgs struct {
	Args
	IsUDP   *bool
	Key     *string
	Timeout *int
}
type TunnelClientArgs struct {
	Args
	IsUDP   *bool
	Key     *string
	Timeout *int
}
type TunnelBridgeArgs struct {
	Args
	Timeout *int
}
type TCPArgs struct {
	Args
	ParentType          *string
	IsTLS               *bool
	Timeout             *int
	PoolSize            *int
	CheckParentInterval *int
}

type HTTPArgs struct {
	Args
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
}
type UDPArgs struct {
	Args
	ParentType          *string
	Timeout             *int
	PoolSize            *int
	CheckParentInterval *int
}

func (a *TCPArgs) Protocol() string {
	if *a.IsTLS {
		return "tls"
	}
	return "tcp"
}

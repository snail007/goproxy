package services

// tcp := app.Command("tcp", "proxy on tcp mode")
// t := tcp.Flag("tcp-timeout", "tcp timeout milliseconds when connect to real server or parent proxy").Default("2000").Int()

const (
	TYPE_TCP  = "tcp"
	TYPE_UDP  = "udp"
	TYPE_HTTP = "http"
	TYPE_TLS  = "tls"
)

type Args struct {
	Local               *string
	Parent              *string
	CertBytes           []byte
	KeyBytes            []byte
	PoolSize            *int
	CheckParentInterval *int
}
type TCPArgs struct {
	Args
	Timeout    *int
	ParentType *string
}
type TLSArgs struct {
	Args
	Timeout    *int
	ParentType *string
}
type HTTPArgs struct {
	Args
	Always      *bool
	HTTPTimeout *int
	Timeout     *int
	Interval    *int
	Blocked     *string
	Direct      *string
	AuthFile    *string
	Auth        *[]string
	ParentType  *string
	LocalType   *string
}
type UDPArgs struct {
	Args
}

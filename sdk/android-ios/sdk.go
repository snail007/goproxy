package proxy

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	logger "log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/snail007/goproxy/core/lib/kcpcfg"
	encryptconn "github.com/snail007/goproxy/core/lib/transport/encrypt"
	"github.com/snail007/goproxy/services"
	httpx "github.com/snail007/goproxy/services/http"
	keygenx "github.com/snail007/goproxy/services/keygen"
	mux "github.com/snail007/goproxy/services/mux"
	socksx "github.com/snail007/goproxy/services/socks"
	spsx "github.com/snail007/goproxy/services/sps"
	tcpx "github.com/snail007/goproxy/services/tcp"
	tunnelx "github.com/snail007/goproxy/services/tunnel"
	udpx "github.com/snail007/goproxy/services/udp"

	kcp "github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var SDK_VERSION = "No Version Provided"

var (
	app *kingpin.Application
)

type LogCallback interface {
	Write(line string)
}
type logCallback interface {
	Write(line string)
}
type logWriter struct {
	callback LogCallback
}

func (s *logWriter) Write(p []byte) (n int, err error) {
	s.callback.Write(string(p))
	return
}

func Start(serviceID, serviceArgsStr string) (errStr string) {
	return StartWithLog(serviceID, serviceArgsStr, nil)
}

//Start
//serviceID : is service identify id,different service's id should be difference
//serviceArgsStr: is the whole command line args string
//such as :
//1."http -t tcp -p :8989"
//2."socks -t tcp -p :8989"
//and so on.
//if an error occured , errStr will be the error reason
//if start success, errStr is empty.
func StartWithLog(serviceID, serviceArgsStr string, loggerCallback LogCallback) (errStr string) {
	//define  args
	tcpArgs := tcpx.TCPArgs{}
	httpArgs := httpx.HTTPArgs{}
	tunnelServerArgs := tunnelx.TunnelServerArgs{}
	tunnelClientArgs := tunnelx.TunnelClientArgs{}
	tunnelBridgeArgs := tunnelx.TunnelBridgeArgs{}
	muxServerArgs := mux.MuxServerArgs{}
	muxClientArgs := mux.MuxClientArgs{}
	muxBridgeArgs := mux.MuxBridgeArgs{}
	udpArgs := udpx.UDPArgs{}
	socksArgs := socksx.SocksArgs{}
	spsArgs := spsx.SPSArgs{}
	dnsArgs := DNSArgs{}
	keygenArgs := keygenx.KeygenArgs{}
	kcpArgs := kcpcfg.KCPConfigArgs{}
	//build srvice args
	app = kingpin.New("proxy", "happy with proxy")
	app.Author("snail").Version(SDK_VERSION)
	debug := app.Flag("debug", "debug log output").Default("false").Bool()
	logfile := app.Flag("log", "log file path").Default("").String()
	nolog := app.Flag("nolog", "turn off logging").Default("false").Bool()
	kcpArgs.Key = app.Flag("kcp-key", "pre-shared secret between client and server").Default("secrect").String()
	kcpArgs.Crypt = app.Flag("kcp-method", "encrypt/decrypt method, can be: aes, aes-128, aes-192, salsa20, blowfish, twofish, cast5, 3des, tea, xtea, xor, sm4, none").Default("aes").Enum("aes", "aes-128", "aes-192", "salsa20", "blowfish", "twofish", "cast5", "3des", "tea", "xtea", "xor", "sm4", "none")
	kcpArgs.Mode = app.Flag("kcp-mode", "profiles: fast3, fast2, fast, normal, manual").Default("fast3").Enum("fast3", "fast2", "fast", "normal", "manual")
	kcpArgs.MTU = app.Flag("kcp-mtu", "set maximum transmission unit for UDP packets").Default("1350").Int()
	kcpArgs.SndWnd = app.Flag("kcp-sndwnd", "set send window size(num of packets)").Default("1024").Int()
	kcpArgs.RcvWnd = app.Flag("kcp-rcvwnd", "set receive window size(num of packets)").Default("1024").Int()
	kcpArgs.DataShard = app.Flag("kcp-ds", "set reed-solomon erasure coding - datashard").Default("10").Int()
	kcpArgs.ParityShard = app.Flag("kcp-ps", "set reed-solomon erasure coding - parityshard").Default("3").Int()
	kcpArgs.DSCP = app.Flag("kcp-dscp", "set DSCP(6bit)").Default("0").Int()
	kcpArgs.NoComp = app.Flag("kcp-nocomp", "disable compression").Default("false").Bool()
	kcpArgs.AckNodelay = app.Flag("kcp-acknodelay", "be carefull! flush ack immediately when a packet is received").Default("true").Bool()
	kcpArgs.NoDelay = app.Flag("kcp-nodelay", "be carefull!").Default("0").Int()
	kcpArgs.Interval = app.Flag("kcp-interval", "be carefull!").Default("50").Int()
	kcpArgs.Resend = app.Flag("kcp-resend", "be carefull!").Default("0").Int()
	kcpArgs.NoCongestion = app.Flag("kcp-nc", "be carefull! no congestion").Default("0").Int()
	kcpArgs.SockBuf = app.Flag("kcp-sockbuf", "be carefull!").Default("4194304").Int()
	kcpArgs.KeepAlive = app.Flag("kcp-keepalive", "be carefull!").Default("10").Int()

	//########http#########
	http := app.Command("http", "proxy on http mode")
	httpArgs.Parent = http.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').Strings()
	httpArgs.CaCertFile = http.Flag("ca", "ca cert file for tls").Default("").String()
	httpArgs.CertFile = http.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	httpArgs.KeyFile = http.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	httpArgs.LocalType = http.Flag("local-type", "local protocol type <tls|tcp|kcp>").Default("tcp").Short('t').Enum("tls", "tcp", "kcp")
	httpArgs.ParentType = http.Flag("parent-type", "parent protocol type <tls|tcp|ssh|kcp>").Short('T').Enum("tls", "tcp", "ssh", "kcp")
	httpArgs.Always = http.Flag("always", "always use parent proxy").Default("false").Bool()
	httpArgs.Timeout = http.Flag("timeout", "tcp timeout milliseconds when connect to real server or parent proxy").Default("2000").Int()
	httpArgs.HTTPTimeout = http.Flag("http-timeout", "check domain if blocked , http request timeout milliseconds when connect to host").Default("3000").Int()
	httpArgs.Interval = http.Flag("interval", "check domain if blocked every interval seconds").Default("10").Int()
	httpArgs.Blocked = http.Flag("blocked", "blocked domain file , one domain each line").Default("blocked").Short('b').String()
	httpArgs.Direct = http.Flag("direct", "direct domain file , one domain each line").Default("direct").Short('d').String()
	httpArgs.AuthFile = http.Flag("auth-file", "http basic auth file,\"username:password\" each line in file").Short('F').String()
	httpArgs.Auth = http.Flag("auth", "http basic auth username and password, mutiple user repeat -a ,such as: -a user1:pass1 -a user2:pass2").Short('a').Strings()
	httpArgs.CheckParentInterval = http.Flag("check-parent-interval", "check if proxy is okay every interval seconds,zero: means no check").Short('I').Default("3").Int()
	httpArgs.Local = http.Flag("local", "local ip:port to listen,multiple address use comma split,such as: 0.0.0.0:80,0.0.0.0:443").Short('p').Default(":33080").String()
	httpArgs.SSHUser = http.Flag("ssh-user", "user for ssh").Short('u').Default("").String()
	httpArgs.SSHKeyFile = http.Flag("ssh-key", "private key file for ssh").Short('S').Default("").String()
	httpArgs.SSHKeyFileSalt = http.Flag("ssh-keysalt", "salt of ssh private key").Short('s').Default("").String()
	httpArgs.SSHPassword = http.Flag("ssh-password", "password for ssh").Short('A').Default("").String()
	httpArgs.LocalIPS = http.Flag("local-bind-ips", "if your host behind a nat,set your public ip here avoid dead loop").Short('g').Strings()
	httpArgs.AuthURL = http.Flag("auth-url", "http basic auth username and password will send to this url,response http code equal to 'auth-code' means ok,others means fail.").Default("").String()
	httpArgs.AuthURLTimeout = http.Flag("auth-timeout", "access 'auth-url' timeout milliseconds").Default("3000").Int()
	httpArgs.AuthURLOkCode = http.Flag("auth-code", "access 'auth-url' success http code").Default("204").Int()
	httpArgs.AuthURLRetry = http.Flag("auth-retry", "access 'auth-url' fail and retry count").Default("1").Int()
	httpArgs.DNSAddress = http.Flag("dns-address", "if set this, proxy will use this dns for resolve doamin").Short('q').Default("").String()
	httpArgs.DNSTTL = http.Flag("dns-ttl", "caching seconds of dns query result").Short('e').Default("300").Int()
	httpArgs.LocalKey = http.Flag("local-key", "the password for auto encrypt/decrypt local connection data").Short('z').Default("").String()
	httpArgs.ParentKey = http.Flag("parent-key", "the password for auto encrypt/decrypt parent connection data").Short('Z').Default("").String()
	httpArgs.LocalCompress = http.Flag("local-compress", "auto compress/decompress data on local connection").Short('m').Default("false").Bool()
	httpArgs.ParentCompress = http.Flag("parent-compress", "auto compress/decompress data on parent connection").Short('M').Default("false").Bool()
	httpArgs.Intelligent = http.Flag("intelligent", "settting intelligent HTTP, SOCKS5 proxy mode, can be <intelligent|direct|parent>").Default("intelligent").Enum("intelligent", "direct", "parent")
	httpArgs.LoadBalanceMethod = http.Flag("lb-method", "load balance method when use multiple parent,can be <roundrobin|leastconn|leasttime|hash|weight>").Default("roundrobin").Enum("roundrobin", "weight", "leastconn", "leasttime", "hash")
	httpArgs.LoadBalanceTimeout = http.Flag("lb-timeout", "tcp milliseconds timeout of connecting to parent").Default("500").Int()
	httpArgs.LoadBalanceRetryTime = http.Flag("lb-retrytime", "sleep time milliseconds after checking").Default("1000").Int()
	httpArgs.LoadBalanceHashTarget = http.Flag("lb-hashtarget", "use target address to choose parent for LB").Default("false").Bool()
	httpArgs.LoadBalanceOnlyHA = http.Flag("lb-onlyha", "use only `high availability mode` to choose parent for LB").Default("false").Bool()
	httpArgs.RateLimit = http.Flag("rate-limit", "rate limit (bytes/second) of each connection, such as: 100K 1.5M . 0 means no limitation").Short('l').Default("0").String()
	httpArgs.BindListen = http.Flag("bind-listen", "using listener binding IP when connect to target").Short('B').Default("false").Bool()
	httpArgs.Jumper = http.Flag("jumper", "https or socks5 proxies used when connecting to parent, only worked of -T is tls or tcp, format is https://username:password@host:port https://host:port or socks5://username:password@host:port socks5://host:port").Short('J').Default("").String()
	httpArgs.Debug = debug
	//########tcp#########
	tcp := app.Command("tcp", "proxy on tcp mode")
	tcpArgs.Parent = tcp.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	tcpArgs.CertFile = tcp.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	tcpArgs.KeyFile = tcp.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	tcpArgs.Timeout = tcp.Flag("timeout", "tcp timeout milliseconds when connect to real server or parent proxy").Short('e').Default("2000").Int()
	tcpArgs.ParentType = tcp.Flag("parent-type", "parent protocol type <tls|tcp|kcp|udp>").Short('T').Enum("tls", "tcp", "udp", "kcp")
	tcpArgs.LocalType = tcp.Flag("local-type", "local protocol type <tls|tcp|kcp>").Default("tcp").Short('t').Enum("tls", "tcp", "kcp")
	tcpArgs.CheckParentInterval = tcp.Flag("check-parent-interval", "check if proxy is okay every interval seconds,zero: means no check").Short('I').Default("3").Int()
	tcpArgs.Local = tcp.Flag("local", "local ip:port to listen").Short('p').Default(":33080").String()
	tcpArgs.Jumper = tcp.Flag("jumper", "https or socks5 proxies used when connecting to parent, only worked of -T is tls or tcp, format is https://username:password@host:port https://host:port or socks5://username:password@host:port socks5://host:port").Short('J').Default("").String()

	//########udp#########
	udp := app.Command("udp", "proxy on udp mode")
	udpArgs.Parent = udp.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	udpArgs.CertFile = udp.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	udpArgs.KeyFile = udp.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	udpArgs.Timeout = udp.Flag("timeout", "tcp timeout milliseconds when connect to parent proxy").Short('t').Default("2000").Int()
	udpArgs.ParentType = udp.Flag("parent-type", "parent protocol type <tls|tcp|udp>").Short('T').Enum("tls", "tcp", "udp")
	udpArgs.CheckParentInterval = udp.Flag("check-parent-interval", "check if proxy is okay every interval seconds,zero: means no check").Short('I').Default("3").Int()
	udpArgs.Local = udp.Flag("local", "local ip:port to listen").Short('p').Default(":33080").String()

	//########mux-server#########
	muxServer := app.Command("server", "proxy on mux server mode")
	muxServerArgs.Parent = muxServer.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	muxServerArgs.ParentType = muxServer.Flag("parent-type", "parent protocol type <tls|tcp|tcps|kcp|tou>").Default("tls").Short('T').Enum("tls", "tcp", "tcps", "kcp", "tou")
	muxServerArgs.CertFile = muxServer.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	muxServerArgs.KeyFile = muxServer.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	muxServerArgs.Timeout = muxServer.Flag("timeout", "tcp timeout with milliseconds").Short('i').Default("2000").Int()
	muxServerArgs.IsUDP = muxServer.Flag("udp", "proxy on udp mux server mode").Default("false").Bool()
	muxServerArgs.Key = muxServer.Flag("k", "client key").Default("default").String()
	muxServerArgs.Route = muxServer.Flag("route", "local route to client's network, such as: PROTOCOL://LOCAL_IP:LOCAL_PORT@[CLIENT_KEY]CLIENT_LOCAL_HOST:CLIENT_LOCAL_PORT").Short('r').Default("").Strings()
	muxServerArgs.IsCompress = muxServer.Flag("c", "compress data when tcp|tls mode").Default("false").Bool()
	muxServerArgs.SessionCount = muxServer.Flag("session-count", "session count which connect to bridge").Short('n').Default("10").Int()
	muxServerArgs.Jumper = muxServer.Flag("jumper", "https or socks5 proxies used when connecting to parent, only worked of -T is tls or tcp, format is https://username:password@host:port https://host:port or socks5://username:password@host:port socks5://host:port").Short('J').Default("").String()
	muxServerArgs.TCPSMethod = muxServer.Flag("tcps-method", "method of parent tcps's encrpyt/decrypt, these below are supported :\n"+strings.Join(encryptconn.GetCipherMethods(), ",")).Default("aes-192-cfb").String()
	muxServerArgs.TCPSPassword = muxServer.Flag("tcps-password", "password of parent tcps's encrpyt/decrypt").Default("snail007's_goproxy").String()
	muxServerArgs.TOUMethod = muxServer.Flag("tou-method", "method of parent tou's encrpyt/decrypt, these below are supported :\n"+strings.Join(encryptconn.GetCipherMethods(), ",")).Default("aes-192-cfb").String()
	muxServerArgs.TOUPassword = muxServer.Flag("tou-password", "password of parent tou's encrpyt/decrypt").Default("snail007's_goproxy").String()

	//########mux-client#########
	muxClient := app.Command("client", "proxy on mux client mode")
	muxClientArgs.Parent = muxClient.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	muxClientArgs.ParentType = muxClient.Flag("parent-type", "parent protocol type <tls|tcp|tcps|kcp|tou>").Default("tls").Short('T').Enum("tls", "tcp", "tcps", "kcp", "tou")
	muxClientArgs.CertFile = muxClient.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	muxClientArgs.KeyFile = muxClient.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	muxClientArgs.Timeout = muxClient.Flag("timeout", "tcp timeout with milliseconds").Short('i').Default("2000").Int()
	muxClientArgs.Key = muxClient.Flag("k", "key same with server").Default("default").String()
	muxClientArgs.IsCompress = muxClient.Flag("c", "compress data when tcp|tls mode").Default("false").Bool()
	muxClientArgs.SessionCount = muxClient.Flag("session-count", "session count which connect to bridge").Short('n').Default("10").Int()
	muxClientArgs.Jumper = muxClient.Flag("jumper", "https or socks5 proxies used when connecting to parent, only worked of -T is tls or tcp, format is https://username:password@host:port https://host:port or socks5://username:password@host:port socks5://host:port").Short('J').Default("").String()
	muxClientArgs.TCPSMethod = muxClient.Flag("tcps-method", "method of parent tcps's encrpyt/decrypt, these below are supported :\n"+strings.Join(encryptconn.GetCipherMethods(), ",")).Default("aes-192-cfb").String()
	muxClientArgs.TCPSPassword = muxClient.Flag("tcps-password", "password of parent tcps's encrpyt/decrypt").Default("snail007's_goproxy").String()
	muxClientArgs.TOUMethod = muxClient.Flag("tou-method", "method of parent tou's encrpyt/decrypt, these below are supported :\n"+strings.Join(encryptconn.GetCipherMethods(), ",")).Default("aes-192-cfb").String()
	muxClientArgs.TOUPassword = muxClient.Flag("tou-password", "password of parent tou's encrpyt/decrypt").Default("snail007's_goproxy").String()

	//########mux-bridge#########
	muxBridge := app.Command("bridge", "proxy on mux bridge mode")
	muxBridgeArgs.CertFile = muxBridge.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	muxBridgeArgs.KeyFile = muxBridge.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	muxBridgeArgs.Timeout = muxBridge.Flag("timeout", "tcp timeout with milliseconds").Short('i').Default("2000").Int()
	muxBridgeArgs.Local = muxBridge.Flag("local", "local ip:port to listen").Short('p').Default(":33080").String()
	muxBridgeArgs.LocalType = muxBridge.Flag("local-type", "local protocol type <tls|tcp|tcps|kcp|tou>").Default("tls").Short('t').Enum("tls", "tcp", "tcps", "kcp", "tou")
	muxBridgeArgs.TCPSMethod = muxBridge.Flag("tcps-method", "method of local tcps's encrpyt/decrypt, these below are supported :\n"+strings.Join(encryptconn.GetCipherMethods(), ",")).Default("aes-192-cfb").String()
	muxBridgeArgs.TCPSPassword = muxBridge.Flag("tcps-password", "password of local tcps's encrpyt/decrypt").Default("snail007's_goproxy").String()
	muxBridgeArgs.TOUMethod = muxBridge.Flag("tou-method", "method of local tou's encrpyt/decrypt, these below are supported :\n"+strings.Join(encryptconn.GetCipherMethods(), ",")).Default("aes-192-cfb").String()
	muxBridgeArgs.TOUPassword = muxBridge.Flag("tou-password", "password of local tou's encrpyt/decrypt").Default("snail007's_goproxy").String()

	//########tunnel-server#########
	tunnelServer := app.Command("tserver", "proxy on tunnel server mode")
	tunnelServerArgs.Parent = tunnelServer.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	tunnelServerArgs.CertFile = tunnelServer.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	tunnelServerArgs.KeyFile = tunnelServer.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	tunnelServerArgs.Timeout = tunnelServer.Flag("timeout", "tcp timeout with milliseconds").Short('t').Default("2000").Int()
	tunnelServerArgs.IsUDP = tunnelServer.Flag("udp", "proxy on udp tunnel server mode").Default("false").Bool()
	tunnelServerArgs.Key = tunnelServer.Flag("k", "client key").Default("default").String()
	tunnelServerArgs.Route = tunnelServer.Flag("route", "local route to client's network, such as: PROTOCOL://LOCAL_IP:LOCAL_PORT@[CLIENT_KEY]CLIENT_LOCAL_HOST:CLIENT_LOCAL_PORT").Short('r').Default("").Strings()
	tunnelServerArgs.Jumper = tunnelServer.Flag("jumper", "https or socks5 proxies used when connecting to parent, only worked of -T is tls or tcp, format is https://username:password@host:port https://host:port or socks5://username:password@host:port socks5://host:port").Short('J').Default("").String()

	//########tunnel-client#########
	tunnelClient := app.Command("tclient", "proxy on tunnel client mode")
	tunnelClientArgs.Parent = tunnelClient.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	tunnelClientArgs.CertFile = tunnelClient.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	tunnelClientArgs.KeyFile = tunnelClient.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	tunnelClientArgs.Timeout = tunnelClient.Flag("timeout", "tcp timeout with milliseconds").Short('t').Default("2000").Int()
	tunnelClientArgs.Key = tunnelClient.Flag("k", "key same with server").Default("default").String()
	tunnelClientArgs.Jumper = tunnelClient.Flag("jumper", "https or socks5 proxies used when connecting to parent, only worked of -T is tls or tcp, format is https://username:password@host:port https://host:port or socks5://username:password@host:port socks5://host:port").Short('J').Default("").String()

	//########tunnel-bridge#########
	tunnelBridge := app.Command("tbridge", "proxy on tunnel bridge mode")
	tunnelBridgeArgs.CertFile = tunnelBridge.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	tunnelBridgeArgs.KeyFile = tunnelBridge.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	tunnelBridgeArgs.Timeout = tunnelBridge.Flag("timeout", "tcp timeout with milliseconds").Short('t').Default("2000").Int()
	tunnelBridgeArgs.Local = tunnelBridge.Flag("local", "local ip:port to listen").Short('p').Default(":33080").String()

	//########ssh#########
	socks := app.Command("socks", "proxy on ssh mode")
	socksArgs.Parent = socks.Flag("parent", "parent ssh address, such as: \"23.32.32.19:22\"").Default("").Short('P').Strings()
	socksArgs.ParentType = socks.Flag("parent-type", "parent protocol type <tls|tcp|kcp|ssh>").Default("tcp").Short('T').Enum("tls", "tcp", "kcp", "ssh")
	socksArgs.LocalType = socks.Flag("local-type", "local protocol type <tls|tcp|kcp>").Default("tcp").Short('t').Enum("tls", "tcp", "kcp")
	socksArgs.Local = socks.Flag("local", "local ip:port to listen").Short('p').Default(":33080").String()
	socksArgs.CertFile = socks.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	socksArgs.CaCertFile = socks.Flag("ca", "ca cert file for tls").Default("").String()
	socksArgs.KeyFile = socks.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	socksArgs.SSHUser = socks.Flag("ssh-user", "user for ssh").Short('u').Default("").String()
	socksArgs.SSHKeyFile = socks.Flag("ssh-key", "private key file for ssh").Short('S').Default("").String()
	socksArgs.SSHKeyFileSalt = socks.Flag("ssh-keysalt", "salt of ssh private key").Short('s').Default("").String()
	socksArgs.SSHPassword = socks.Flag("ssh-password", "password for ssh").Short('D').Default("").String()
	socksArgs.Always = socks.Flag("always", "always use parent proxy").Default("false").Bool()
	socksArgs.Timeout = socks.Flag("timeout", "tcp timeout milliseconds when connect to real server or parent proxy").Default("5000").Int()
	socksArgs.Interval = socks.Flag("interval", "check domain if blocked every interval seconds").Default("10").Int()
	socksArgs.Blocked = socks.Flag("blocked", "blocked domain file , one domain each line").Default("blocked").Short('b').String()
	socksArgs.Direct = socks.Flag("direct", "direct domain file , one domain each line").Default("direct").Short('d').String()
	socksArgs.AuthFile = socks.Flag("auth-file", "http basic auth file,\"username:password\" each line in file").Short('F').String()
	socksArgs.Auth = socks.Flag("auth", "socks auth username and password, mutiple user repeat -a ,such as: -a user1:pass1 -a user2:pass2").Short('a').Strings()
	socksArgs.LocalIPS = socks.Flag("local-bind-ips", "if your host behind a nat,set your public ip here avoid dead loop").Short('g').Strings()
	socksArgs.AuthURL = socks.Flag("auth-url", "auth username and password will send to this url,response http code equal to 'auth-code' means ok,others means fail.").Default("").String()
	socksArgs.AuthURLTimeout = socks.Flag("auth-timeout", "access 'auth-url' timeout milliseconds").Default("3000").Int()
	socksArgs.AuthURLOkCode = socks.Flag("auth-code", "access 'auth-url' success http code").Default("204").Int()
	socksArgs.AuthURLRetry = socks.Flag("auth-retry", "access 'auth-url' fail and retry count").Default("0").Int()
	socksArgs.ParentAuth = socks.Flag("parent-auth", "parent socks auth username and password, such as: -A user1:pass1").Short('A').String()
	socksArgs.DNSAddress = socks.Flag("dns-address", "if set this, proxy will use this dns for resolve doamin").Short('q').Default("").String()
	socksArgs.DNSTTL = socks.Flag("dns-ttl", "caching seconds of dns query result").Short('e').Default("300").Int()
	socksArgs.LocalKey = socks.Flag("local-key", "the password for auto encrypt/decrypt local connection data").Short('z').Default("").String()
	socksArgs.ParentKey = socks.Flag("parent-key", "the password for auto encrypt/decrypt parent connection data").Short('Z').Default("").String()
	socksArgs.LocalCompress = socks.Flag("local-compress", "auto compress/decompress data on local connection").Short('m').Default("false").Bool()
	socksArgs.ParentCompress = socks.Flag("parent-compress", "auto compress/decompress data on parent connection").Short('M').Default("false").Bool()
	socksArgs.Intelligent = socks.Flag("intelligent", "settting intelligent HTTP, SOCKS5 proxy mode, can be <intelligent|direct|parent>").Default("intelligent").Enum("intelligent", "direct", "parent")
	socksArgs.LoadBalanceMethod = socks.Flag("lb-method", "load balance method when use multiple parent,can be <roundrobin|leastconn|leasttime|hash|weight>").Default("roundrobin").Enum("roundrobin", "weight", "leastconn", "leasttime", "hash")
	socksArgs.LoadBalanceTimeout = socks.Flag("lb-timeout", "tcp milliseconds timeout of connecting to parent").Default("500").Int()
	socksArgs.LoadBalanceRetryTime = socks.Flag("lb-retrytime", "sleep time milliseconds after checking").Default("1000").Int()
	socksArgs.LoadBalanceHashTarget = socks.Flag("lb-hashtarget", "use target address to choose parent for LB").Default("false").Bool()
	socksArgs.LoadBalanceOnlyHA = socks.Flag("lb-onlyha", "use only `high availability mode` to choose parent for LB").Default("false").Bool()
	socksArgs.RateLimit = socks.Flag("rate-limit", "rate limit (bytes/second) of each connection, such as: 100K 1.5M . 0 means no limitation").Short('l').Default("0").String()
	socksArgs.BindListen = socks.Flag("bind-listen", "using listener binding IP when connect to target").Short('B').Default("false").Bool()
	socksArgs.Debug = debug

	//########socks+http(s)#########
	sps := app.Command("sps", "proxy on socks+http(s) mode")
	spsArgs.Parent = sps.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').Strings()
	spsArgs.CertFile = sps.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	spsArgs.KeyFile = sps.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	spsArgs.CaCertFile = sps.Flag("ca", "ca cert file for tls").Default("").String()
	spsArgs.Timeout = sps.Flag("timeout", "tcp timeout milliseconds when connect to real server or parent proxy").Short('i').Default("2000").Int()
	spsArgs.ParentType = sps.Flag("parent-type", "parent protocol type <tls|tcp|kcp>").Short('T').Enum("tls", "tcp", "kcp")
	spsArgs.LocalType = sps.Flag("local-type", "local protocol type <tls|tcp|kcp>").Default("tcp").Short('t').Enum("tls", "tcp", "kcp")
	spsArgs.Local = sps.Flag("local", "local ip:port to listen,multiple address use comma split,such as: 0.0.0.0:80,0.0.0.0:443").Short('p').Default(":33080").String()
	spsArgs.ParentServiceType = sps.Flag("parent-service-type", "parent service type <http|socks|ss>").Short('S').Enum("http", "socks", "ss")
	spsArgs.DNSAddress = sps.Flag("dns-address", "if set this, proxy will use this dns for resolve doamin").Short('q').Default("").String()
	spsArgs.DNSTTL = sps.Flag("dns-ttl", "caching seconds of dns query result").Short('e').Default("300").Int()
	spsArgs.AuthFile = sps.Flag("auth-file", "http basic auth file,\"username:password\" each line in file").Short('F').String()
	spsArgs.Auth = sps.Flag("auth", "socks auth username and password, mutiple user repeat -a ,such as: -a user1:pass1 -a user2:pass2").Short('a').Strings()
	spsArgs.LocalIPS = sps.Flag("local-bind-ips", "if your host behind a nat,set your public ip here avoid dead loop").Short('g').Strings()
	spsArgs.AuthURL = sps.Flag("auth-url", "auth username and password will send to this url,response http code equal to 'auth-code' means ok,others means fail.").Default("").String()
	spsArgs.AuthURLTimeout = sps.Flag("auth-timeout", "access 'auth-url' timeout milliseconds").Default("3000").Int()
	spsArgs.AuthURLOkCode = sps.Flag("auth-code", "access 'auth-url' success http code").Default("204").Int()
	spsArgs.AuthURLRetry = sps.Flag("auth-retry", "access 'auth-url' fail and retry count").Default("0").Int()
	spsArgs.ParentAuth = sps.Flag("parent-auth", "parent socks auth username and password, such as: -A user1:pass1").Short('A').String()
	spsArgs.LocalKey = sps.Flag("local-key", "the password for auto encrypt/decrypt local connection data").Short('z').Default("").String()
	spsArgs.ParentKey = sps.Flag("parent-key", "the password for auto encrypt/decrypt parent connection data").Short('Z').Default("").String()
	spsArgs.LocalCompress = sps.Flag("local-compress", "auto compress/decompress data on local connection").Short('m').Default("false").Bool()
	spsArgs.ParentCompress = sps.Flag("parent-compress", "auto compress/decompress data on parent connection").Short('M').Default("false").Bool()
	spsArgs.SSMethod = sps.Flag("ss-method", "the following methods are supported: aes-128-cfb, aes-192-cfb, aes-256-cfb, bf-cfb, cast5-cfb, des-cfb, rc4-md5, rc4-md5-6, chacha20, salsa20, rc4, table, des-cfb, chacha20-ietf; if you use ss client , \"-t tcp\" is required").Short('h').Default("aes-256-cfb").String()
	spsArgs.SSKey = sps.Flag("ss-key", "if you use ss client , \"-t tcp\" is required").Short('j').Default("sspassword").String()
	spsArgs.ParentSSMethod = sps.Flag("parent-ss-method", "the following methods are supported: aes-128-cfb, aes-192-cfb, aes-256-cfb, bf-cfb, cast5-cfb, des-cfb, rc4-md5, rc4-md5-6, chacha20, salsa20, rc4, table, des-cfb, chacha20-ietf; if you use ss server as parent, \"-T tcp\" is required").Short('H').Default("aes-256-cfb").String()
	spsArgs.ParentSSKey = sps.Flag("parent-ss-key", "if you use ss server as parent, \"-T tcp\" is required").Short('J').Default("sspassword").String()
	spsArgs.DisableHTTP = sps.Flag("disable-http", "disable http(s) proxy").Default("false").Bool()
	spsArgs.DisableSocks5 = sps.Flag("disable-socks", "disable socks proxy").Default("false").Bool()
	spsArgs.DisableSS = sps.Flag("disable-ss", "disable ss proxy").Default("false").Bool()
	spsArgs.LoadBalanceMethod = sps.Flag("lb-method", "load balance method when use multiple parent,can be <roundrobin|leastconn|leasttime|hash|weight>").Default("hash").Enum("roundrobin", "weight", "leastconn", "leasttime", "hash")
	spsArgs.LoadBalanceTimeout = sps.Flag("lb-timeout", "tcp milliseconds timeout of connecting to parent").Default("500").Int()
	spsArgs.LoadBalanceRetryTime = sps.Flag("lb-retrytime", "sleep time milliseconds after checking").Default("1000").Int()
	spsArgs.LoadBalanceHashTarget = sps.Flag("lb-hashtarget", "use target address to choose parent for LB").Default("false").Bool()
	spsArgs.LoadBalanceOnlyHA = sps.Flag("lb-onlyha", "use only `high availability mode` to choose parent for LB").Default("false").Bool()
	spsArgs.RateLimit = sps.Flag("rate-limit", "rate limit (bytes/second) of each connection, such as: 100K 1.5M . 0 means no limitation").Short('l').Default("0").String()
	spsArgs.Jumper = sps.Flag("jumper", "https or socks5 proxies used when connecting to parent, only worked of -T is tls or tcp, format is https://username:password@host:port https://host:port or socks5://username:password@host:port socks5://host:port").Default("").String()
	spsArgs.Debug = debug

	//########dns#########
	dns := app.Command("dns", "proxy on dns server mode")
	dnsArgs.Parent = dns.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	dnsArgs.CertFile = dns.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	dnsArgs.KeyFile = dns.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	dnsArgs.CaCertFile = dns.Flag("ca", "ca cert file for tls").Default("").String()
	dnsArgs.Timeout = dns.Flag("timeout", "tcp timeout milliseconds when connect to real server or parent proxy").Short('i').Default("2000").Int()
	dnsArgs.ParentType = dns.Flag("parent-type", "parent protocol type <tls|tcp|kcp>").Short('T').Enum("tls", "tcp", "kcp")
	dnsArgs.Local = dns.Flag("local", "local ip:port to listen,multiple address use comma split,such as: 0.0.0.0:80,0.0.0.0:443").Short('p').Default(":53").String()
	dnsArgs.ParentServiceType = dns.Flag("parent-service-type", "parent service type <http|socks>").Short('S').Enum("http", "socks")
	dnsArgs.RemoteDNSAddress = dns.Flag("dns-address", "remote dns for resolve doamin").Short('q').Default("8.8.8.8:53").String()
	dnsArgs.DNSTTL = dns.Flag("dns-ttl", "caching seconds of dns query result").Short('e').Default("300").Int()
	dnsArgs.ParentAuth = dns.Flag("parent-auth", "parent socks auth username and password, such as: -A user1:pass1").Short('A').String()
	dnsArgs.ParentKey = dns.Flag("parent-key", "the password for auto encrypt/decrypt parent connection data").Short('Z').Default("").String()
	dnsArgs.ParentCompress = dns.Flag("parent-compress", "auto compress/decompress data on parent connection").Short('M').Default("false").Bool()
	dnsArgs.CacheFile = dns.Flag("cache-file", "dns result cached file").Short('f').Default(filepath.Join(path.Dir(os.Args[0]), "cache.dat")).String()
	dnsArgs.LocalSocks5Port = dns.Flag("socks-port", "local socks5 port").Short('s').Default("65501").String()

	//########keygen#########
	keygen := app.Command("keygen", "create certificate for proxy")
	keygenArgs.CommonName = keygen.Flag("cn", "common name").Short('n').Default("").String()
	keygenArgs.CaName = keygen.Flag("ca", "ca name").Short('C').Default("").String()
	keygenArgs.CertName = keygen.Flag("cert", "cert name of sign to create").Short('c').Default("").String()
	keygenArgs.SignDays = keygen.Flag("days", "days of sign").Short('d').Default("365").Int()
	keygenArgs.Sign = keygen.Flag("sign", "cert is to signin").Short('s').Default("false").Bool()

	//parse args
	_args := strings.Fields(strings.Trim(serviceArgsStr, " "))
	args := []string{}
	for _, a := range _args {
		args = append(args, strings.Trim(a, "\""))
	}
	serviceName, err := app.Parse(args)
	if err != nil {
		return fmt.Sprintf("parse args fail,err: %s", err)
	}
	//set kcp config

	switch *kcpArgs.Mode {
	case "normal":
		*kcpArgs.NoDelay, *kcpArgs.Interval, *kcpArgs.Resend, *kcpArgs.NoCongestion = 0, 40, 2, 1
	case "fast":
		*kcpArgs.NoDelay, *kcpArgs.Interval, *kcpArgs.Resend, *kcpArgs.NoCongestion = 0, 30, 2, 1
	case "fast2":
		*kcpArgs.NoDelay, *kcpArgs.Interval, *kcpArgs.Resend, *kcpArgs.NoCongestion = 1, 20, 2, 1
	case "fast3":
		*kcpArgs.NoDelay, *kcpArgs.Interval, *kcpArgs.Resend, *kcpArgs.NoCongestion = 1, 10, 2, 1
	}
	pass := pbkdf2.Key([]byte(*kcpArgs.Key), []byte("snail007-goproxy"), 4096, 32, sha1.New)

	switch *kcpArgs.Crypt {
	case "sm4":
		kcpArgs.Block, _ = kcp.NewSM4BlockCrypt(pass[:16])
	case "tea":
		kcpArgs.Block, _ = kcp.NewTEABlockCrypt(pass[:16])
	case "xor":
		kcpArgs.Block, _ = kcp.NewSimpleXORBlockCrypt(pass)
	case "none":
		kcpArgs.Block, _ = kcp.NewNoneBlockCrypt(pass)
	case "aes-128":
		kcpArgs.Block, _ = kcp.NewAESBlockCrypt(pass[:16])
	case "aes-192":
		kcpArgs.Block, _ = kcp.NewAESBlockCrypt(pass[:24])
	case "blowfish":
		kcpArgs.Block, _ = kcp.NewBlowfishBlockCrypt(pass)
	case "twofish":
		kcpArgs.Block, _ = kcp.NewTwofishBlockCrypt(pass)
	case "cast5":
		kcpArgs.Block, _ = kcp.NewCast5BlockCrypt(pass[:16])
	case "3des":
		kcpArgs.Block, _ = kcp.NewTripleDESBlockCrypt(pass[:24])
	case "xtea":
		kcpArgs.Block, _ = kcp.NewXTEABlockCrypt(pass[:16])
	case "salsa20":
		kcpArgs.Block, _ = kcp.NewSalsa20BlockCrypt(pass)
	default:
		*kcpArgs.Crypt = "aes"
		kcpArgs.Block, _ = kcp.NewAESBlockCrypt(pass)
	}
	//attach kcp config
	tcpArgs.KCP = kcpArgs
	httpArgs.KCP = kcpArgs
	socksArgs.KCP = kcpArgs
	spsArgs.KCP = kcpArgs
	muxBridgeArgs.KCP = kcpArgs
	muxServerArgs.KCP = kcpArgs
	muxClientArgs.KCP = kcpArgs
	dnsArgs.KCP = kcpArgs

	log := logger.New(os.Stderr, "", logger.Ldate|logger.Ltime)
	flags := logger.Ldate
	if *debug {
		flags |= logger.Lshortfile | logger.Lmicroseconds
	} else {
		flags |= logger.Ltime
	}
	log.SetFlags(flags)

	if loggerCallback == nil {
		if *nolog {
			log.SetOutput(ioutil.Discard)
		} else if *logfile != "" {
			f, e := os.OpenFile(*logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
			if e != nil {
				log.Fatal(e)
			}
			log.SetOutput(f)
		}
	} else {
		log.SetOutput(&logWriter{
			callback: loggerCallback,
		})
	}

	//regist services and run service
	switch serviceName {
	case "http":
		services.Regist(serviceID, httpx.NewHTTP(), httpArgs, log)
	case "tcp":
		services.Regist(serviceID, tcpx.NewTCP(), tcpArgs, log)
	case "udp":
		services.Regist(serviceID, udpx.NewUDP(), udpArgs, log)
	case "tserver":
		services.Regist(serviceID, tunnelx.NewTunnelServerManager(), tunnelServerArgs, log)
	case "tclient":
		services.Regist(serviceID, tunnelx.NewTunnelClient(), tunnelClientArgs, log)
	case "tbridge":
		services.Regist(serviceID, tunnelx.NewTunnelBridge(), tunnelBridgeArgs, log)
	case "server":
		services.Regist(serviceID, mux.NewMuxServerManager(), muxServerArgs, log)
	case "client":
		services.Regist(serviceID, mux.NewMuxClient(), muxClientArgs, log)
	case "bridge":
		services.Regist(serviceID, mux.NewMuxBridge(), muxBridgeArgs, log)
	case "socks":
		services.Regist(serviceID, socksx.NewSocks(), socksArgs, log)
	case "sps":
		services.Regist(serviceID, spsx.NewSPS(), spsArgs, log)
	case "dns":
		services.Regist(serviceID, NewDNS(), dnsArgs, log)
	}
	_, err = services.Run(serviceID, nil)
	if err != nil {
		return fmt.Sprintf("run service [%s:%s] fail, ERR:%s", serviceID, serviceName, err)
	}
	return
}

func Stop(serviceID string) {
	services.Stop(serviceID)
}

func Version() string {
	return SDK_VERSION
}

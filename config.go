package main

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	logger "log"
	"os"
	"os/exec"
	"runtime/pprof"
	"time"

	"github.com/visenze/goproxy/services"
	"github.com/visenze/goproxy/services/kcpcfg"

	httpx "github.com/visenze/goproxy/services/http"
	keygenx "github.com/visenze/goproxy/services/keygen"
	mux "github.com/visenze/goproxy/services/mux"
	socksx "github.com/visenze/goproxy/services/socks"
	spsx "github.com/visenze/goproxy/services/sps"
	tcpx "github.com/visenze/goproxy/services/tcp"
	tunnel "github.com/visenze/goproxy/services/tunnel"

	kcp "github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app                                                                                                       *kingpin.Application
	service                                                                                                   *services.ServiceItem
	cmd                                                                                                       *exec.Cmd
	cpuProfilingFile, memProfilingFile, blockProfilingFile, goroutineProfilingFile, threadcreateProfilingFile *os.File
	isDebug                                                                                                   bool
)

func initConfig() (err error) {
	//define  args
	tcpArgs := tcpx.TCPArgs{}
	httpArgs := httpx.HTTPArgs{}
	tunnelServerArgs := tunnel.TunnelServerArgs{}
	tunnelClientArgs := tunnel.TunnelClientArgs{}
	tunnelBridgeArgs := tunnel.TunnelBridgeArgs{}
	muxServerArgs := mux.MuxServerArgs{}
	muxClientArgs := mux.MuxClientArgs{}
	muxBridgeArgs := mux.MuxBridgeArgs{}
	socksArgs := socksx.SocksArgs{}
	spsArgs := spsx.SPSArgs{}
	keygenArgs := keygenx.KeygenArgs{}
	kcpArgs := kcpcfg.KCPConfigArgs{}

	//build srvice args
	app = kingpin.New("proxy", "happy with proxy")
	app.Author("snail").Version(APP_VERSION)
	debug := app.Flag("debug", "debug log output").Default("false").Bool()
	daemon := app.Flag("daemon", "run proxy in background").Default("false").Bool()
	forever := app.Flag("forever", "run proxy in forever,fail and retry").Default("false").Bool()
	logfile := app.Flag("log", "log file path").Default("").String()
	kcpArgs.Key = app.Flag("kcp-key", "pre-shared secret between client and server").Default("secrect").String()
	kcpArgs.Crypt = app.Flag("kcp-method", "encrypt/decrypt method, can be: aes, aes-128, aes-192, salsa20, blowfish, twofish, cast5, 3des, tea, xtea, xor, sm4, none").Default("aes").Enum("aes", "aes-128", "aes-192", "salsa20", "blowfish", "twofish", "cast5", "3des", "tea", "xtea", "xor", "sm4", "none")
	kcpArgs.Mode = app.Flag("kcp-mode", "profiles: fast3, fast2, fast, normal, manual").Default("fast").Enum("fast3", "fast2", "fast", "normal", "manual")
	kcpArgs.MTU = app.Flag("kcp-mtu", "set maximum transmission unit for UDP packets").Default("450").Int()
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
	httpArgs.Parent = http.Flag("parent", "parent address, such as: \"209.205.219.26:3000\"").Default("").Short('P').String()
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

	//########mux-server#########
	muxServer := app.Command("server", "proxy on mux server mode")
	muxServerArgs.Parent = muxServer.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	muxServerArgs.ParentType = muxServer.Flag("parent-type", "parent protocol type <tls|tcp|kcp>").Default("tls").Short('T').Enum("tls", "tcp", "kcp")
	muxServerArgs.CertFile = muxServer.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	muxServerArgs.KeyFile = muxServer.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	muxServerArgs.Timeout = muxServer.Flag("timeout", "tcp timeout with milliseconds").Short('i').Default("2000").Int()
	muxServerArgs.IsUDP = muxServer.Flag("udp", "proxy on udp mux server mode").Default("false").Bool()
	muxServerArgs.Key = muxServer.Flag("k", "client key").Default("default").String()
	muxServerArgs.Route = muxServer.Flag("route", "local route to client's network, such as: PROTOCOL://LOCAL_IP:LOCAL_PORT@[CLIENT_KEY]CLIENT_LOCAL_HOST:CLIENT_LOCAL_PORT").Short('r').Default("").Strings()
	muxServerArgs.IsCompress = muxServer.Flag("c", "compress data when tcp|tls mode").Default("false").Bool()
	muxServerArgs.SessionCount = muxServer.Flag("session-count", "session count which connect to bridge").Short('n').Default("10").Int()

	//########mux-client#########
	muxClient := app.Command("client", "proxy on mux client mode")
	muxClientArgs.Parent = muxClient.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	muxClientArgs.ParentType = muxClient.Flag("parent-type", "parent protocol type <tls|tcp|kcp>").Default("tls").Short('T').Enum("tls", "tcp", "kcp")
	muxClientArgs.CertFile = muxClient.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	muxClientArgs.KeyFile = muxClient.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	muxClientArgs.Timeout = muxClient.Flag("timeout", "tcp timeout with milliseconds").Short('i').Default("2000").Int()
	muxClientArgs.Key = muxClient.Flag("k", "key same with server").Default("default").String()
	muxClientArgs.IsCompress = muxClient.Flag("c", "compress data when tcp|tls mode").Default("false").Bool()
	muxClientArgs.SessionCount = muxClient.Flag("session-count", "session count which connect to bridge").Short('n').Default("10").Int()

	//########mux-bridge#########
	muxBridge := app.Command("bridge", "proxy on mux bridge mode")
	muxBridgeArgs.CertFile = muxBridge.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	muxBridgeArgs.KeyFile = muxBridge.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	muxBridgeArgs.Timeout = muxBridge.Flag("timeout", "tcp timeout with milliseconds").Short('i').Default("2000").Int()
	muxBridgeArgs.Local = muxBridge.Flag("local", "local ip:port to listen").Short('p').Default(":33080").String()
	muxBridgeArgs.LocalType = muxBridge.Flag("local-type", "local protocol type <tls|tcp|kcp>").Default("tls").Short('t').Enum("tls", "tcp", "kcp")

	//########tunnel-server#########
	tunnelServer := app.Command("tserver", "proxy on tunnel server mode")
	tunnelServerArgs.Parent = tunnelServer.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	tunnelServerArgs.CertFile = tunnelServer.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	tunnelServerArgs.KeyFile = tunnelServer.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	tunnelServerArgs.Timeout = tunnelServer.Flag("timeout", "tcp timeout with milliseconds").Short('t').Default("2000").Int()
	tunnelServerArgs.IsUDP = tunnelServer.Flag("udp", "proxy on udp tunnel server mode").Default("false").Bool()
	tunnelServerArgs.Key = tunnelServer.Flag("k", "client key").Default("default").String()
	tunnelServerArgs.Route = tunnelServer.Flag("route", "local route to client's network, such as: PROTOCOL://LOCAL_IP:LOCAL_PORT@[CLIENT_KEY]CLIENT_LOCAL_HOST:CLIENT_LOCAL_PORT").Short('r').Default("").Strings()

	//########tunnel-client#########
	tunnelClient := app.Command("tclient", "proxy on tunnel client mode")
	tunnelClientArgs.Parent = tunnelClient.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	tunnelClientArgs.CertFile = tunnelClient.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	tunnelClientArgs.KeyFile = tunnelClient.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	tunnelClientArgs.Timeout = tunnelClient.Flag("timeout", "tcp timeout with milliseconds").Short('t').Default("2000").Int()
	tunnelClientArgs.Key = tunnelClient.Flag("k", "key same with server").Default("default").String()

	//########tunnel-bridge#########
	tunnelBridge := app.Command("tbridge", "proxy on tunnel bridge mode")
	tunnelBridgeArgs.CertFile = tunnelBridge.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	tunnelBridgeArgs.KeyFile = tunnelBridge.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	tunnelBridgeArgs.Timeout = tunnelBridge.Flag("timeout", "tcp timeout with milliseconds").Short('t').Default("2000").Int()
	tunnelBridgeArgs.Local = tunnelBridge.Flag("local", "local ip:port to listen").Short('p').Default(":33080").String()

	//########ssh#########
	socks := app.Command("socks", "proxy on ssh mode")
	socksArgs.Parent = socks.Flag("parent", "parent ssh address, such as: \"23.32.32.19:22\"").Default("").Short('P').String()
	socksArgs.ParentType = socks.Flag("parent-type", "parent protocol type <tls|tcp|kcp|ssh>").Default("tcp").Short('T').Enum("tls", "tcp", "kcp", "ssh")
	socksArgs.LocalType = socks.Flag("local-type", "local protocol type <tls|tcp|kcp>").Default("tcp").Short('t').Enum("tls", "tcp", "kcp")
	socksArgs.Local = socks.Flag("local", "local ip:port to listen").Short('p').Default(":33080").String()
	socksArgs.CertFile = socks.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	socksArgs.CaCertFile = socks.Flag("ca", "ca cert file for tls").Default("").String()
	socksArgs.KeyFile = socks.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	socksArgs.SSHUser = socks.Flag("ssh-user", "user for ssh").Short('u').Default("").String()
	socksArgs.SSHKeyFile = socks.Flag("ssh-key", "private key file for ssh").Short('S').Default("").String()
	socksArgs.SSHKeyFileSalt = socks.Flag("ssh-keysalt", "salt of ssh private key").Short('s').Default("").String()
	socksArgs.SSHPassword = socks.Flag("ssh-password", "password for ssh").Short('A').Default("").String()
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
	socksArgs.DNSAddress = socks.Flag("dns-address", "if set this, proxy will use this dns for resolve doamin").Short('q').Default("").String()
	socksArgs.DNSTTL = socks.Flag("dns-ttl", "caching seconds of dns query result").Short('e').Default("300").Int()
	socksArgs.LocalKey = socks.Flag("local-key", "the password for auto encrypt/decrypt local connection data").Short('z').Default("").String()
	socksArgs.ParentKey = socks.Flag("parent-key", "the password for auto encrypt/decrypt parent connection data").Short('Z').Default("").String()
	socksArgs.LocalCompress = socks.Flag("local-compress", "auto compress/decompress data on local connection").Short('m').Default("false").Bool()
	socksArgs.ParentCompress = socks.Flag("parent-compress", "auto compress/decompress data on parent connection").Short('M').Default("false").Bool()

	//########socks+http(s)#########
	sps := app.Command("sps", "proxy on socks+http(s) mode")
	spsArgs.Parent = sps.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	spsArgs.CertFile = sps.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	spsArgs.KeyFile = sps.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	spsArgs.CaCertFile = sps.Flag("ca", "ca cert file for tls").Default("").String()
	spsArgs.Timeout = sps.Flag("timeout", "tcp timeout milliseconds when connect to real server or parent proxy").Short('i').Default("2000").Int()
	spsArgs.ParentType = sps.Flag("parent-type", "parent protocol type <tls|tcp|kcp>").Short('T').Enum("tls", "tcp", "kcp")
	spsArgs.LocalType = sps.Flag("local-type", "local protocol type <tls|tcp|kcp>").Default("tcp").Short('t').Enum("tls", "tcp", "kcp")
	spsArgs.Local = sps.Flag("local", "local ip:port to listen,multiple address use comma split,such as: 0.0.0.0:80,0.0.0.0:443").Short('p').Default(":33080").String()
	spsArgs.ParentServiceType = sps.Flag("parent-service-type", "parent service type <http|socks>").Short('S').Enum("http", "socks")
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
	spsArgs.DisableHTTP = sps.Flag("disable-http", "disable http(s) proxy").Default("false").Bool()
	spsArgs.DisableSocks5 = sps.Flag("disable-socks", "disable socks proxy").Default("false").Bool()

	//########keygen#########
	keygen := app.Command("keygen", "create certificate for proxy")
	keygenArgs.CommonName = keygen.Flag("cn", "common name").Short('n').Default("").String()
	keygenArgs.CaName = keygen.Flag("ca", "ca name").Short('C').Default("").String()
	keygenArgs.CertName = keygen.Flag("cert", "cert name of sign to create").Short('c').Default("").String()
	keygenArgs.SignDays = keygen.Flag("days", "days of sign").Short('d').Default("365").Int()
	keygenArgs.Sign = keygen.Flag("sign", "cert is to signin").Short('s').Default("false").Bool()

	//parse args
	serviceName := kingpin.MustParse(app.Parse(os.Args[1:]))

	isDebug = *debug

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

	log := logger.New(os.Stderr, "", logger.Ldate|logger.Ltime)

	flags := logger.Ldate
	if *debug {
		flags |= logger.Lshortfile | logger.Lmicroseconds
		cpuProfilingFile, _ = os.Create("cpu.prof")
		memProfilingFile, _ = os.Create("memory.prof")
		blockProfilingFile, _ = os.Create("block.prof")
		goroutineProfilingFile, _ = os.Create("goroutine.prof")
		threadcreateProfilingFile, _ = os.Create("threadcreate.prof")
		pprof.StartCPUProfile(cpuProfilingFile)
	} else {
		flags |= logger.Ltime
	}
	log.SetFlags(flags)

	if *logfile != "" {
		f, e := os.OpenFile(*logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if e != nil {
			log.Fatal(e)
		}
		log.SetOutput(f)
	}
	if *daemon {
		args := []string{}
		for _, arg := range os.Args[1:] {
			if arg != "--daemon" {
				args = append(args, arg)
			}
		}
		cmd = exec.Command(os.Args[0], args...)
		cmd.Start()
		f := ""
		if *forever {
			f = "forever "
		}
		log.Printf("%s%s [PID] %d running...\n", f, os.Args[0], cmd.Process.Pid)
		os.Exit(0)
	}
	if *forever {
		args := []string{}
		for _, arg := range os.Args[1:] {
			if arg != "--forever" {
				args = append(args, arg)
			}
		}
		go func() {
			for {
				if cmd != nil {
					cmd.Process.Kill()
				}
				cmd = exec.Command(os.Args[0], args...)
				cmdReaderStderr, err := cmd.StderrPipe()
				if err != nil {
					log.Printf("ERR:%s,restarting...\n", err)
					continue
				}
				cmdReader, err := cmd.StdoutPipe()
				if err != nil {
					log.Printf("ERR:%s,restarting...\n", err)
					continue
				}
				scanner := bufio.NewScanner(cmdReader)
				scannerStdErr := bufio.NewScanner(cmdReaderStderr)
				go func() {
					for scanner.Scan() {
						fmt.Println(scanner.Text())
					}
				}()
				go func() {
					for scannerStdErr.Scan() {
						fmt.Println(scannerStdErr.Text())
					}
				}()
				if err := cmd.Start(); err != nil {
					log.Printf("ERR:%s,restarting...\n", err)
					continue
				}
				pid := cmd.Process.Pid
				log.Printf("worker %s [PID] %d running...\n", os.Args[0], pid)
				if err := cmd.Wait(); err != nil {
					log.Printf("ERR:%s,restarting...", err)
					continue
				}
				log.Printf("worker %s [PID] %d unexpected exited, restarting...\n", os.Args[0], pid)
				time.Sleep(time.Second * 5)
			}
		}()
		return
	}
	if *logfile == "" {
		poster()
		if *debug {
			log.Println("[profiling] cpu profiling save to file : cpu.prof")
			log.Println("[profiling] memory profiling save to file : memory.prof")
			log.Println("[profiling] block profiling save to file : block.prof")
			log.Println("[profiling] goroutine profiling save to file : goroutine.prof")
			log.Println("[profiling] threadcreate profiling save to file : threadcreate.prof")
		}
	}
	//regist services and run service
	//regist services and run service
	switch serviceName {
	case "http":
		services.Regist(serviceName, httpx.NewHTTP(), httpArgs, log)
	case "tcp":
		services.Regist(serviceName, tcpx.NewTCP(), tcpArgs, log)
	case "tserver":
		services.Regist(serviceName, tunnel.NewTunnelServerManager(), tunnelServerArgs, log)
	case "tclient":
		services.Regist(serviceName, tunnel.NewTunnelClient(), tunnelClientArgs, log)
	case "tbridge":
		services.Regist(serviceName, tunnel.NewTunnelBridge(), tunnelBridgeArgs, log)
	case "server":
		services.Regist(serviceName, mux.NewMuxServerManager(), muxServerArgs, log)
	case "client":
		services.Regist(serviceName, mux.NewMuxClient(), muxClientArgs, log)
	case "bridge":
		services.Regist(serviceName, mux.NewMuxBridge(), muxBridgeArgs, log)
	case "socks":
		services.Regist(serviceName, socksx.NewSocks(), socksArgs, log)
	case "sps":
		services.Regist(serviceName, spsx.NewSPS(), spsArgs, log)
	case "keygen":
		services.Regist(serviceName, keygenx.NewKeygen(), keygenArgs, log)
	}

	service, err = services.Run(serviceName, nil)
	if err != nil {
		log.Fatalf("run service [%s] fail, ERR:%s", serviceName, err)
	}
	return
}

func poster() {
	fmt.Printf(`
		########  ########   #######  ##     ## ##    ## 
		##     ## ##     ## ##     ##  ##   ##   ##  ##  
		##     ## ##     ## ##     ##   ## ##     ####   
		########  ########  ##     ##    ###       ##    
		##        ##   ##   ##     ##   ## ##      ##    
		##        ##    ##  ##     ##  ##   ##     ##    
		##        ##     ##  #######  ##     ##    ##    
		
		v%s`+" by snail , blog : http://www.host900.com/\n\n", APP_VERSION)
}
func saveProfiling() {
	goroutine := pprof.Lookup("goroutine")
	goroutine.WriteTo(goroutineProfilingFile, 1)
	heap := pprof.Lookup("heap")
	heap.WriteTo(memProfilingFile, 1)
	block := pprof.Lookup("block")
	block.WriteTo(blockProfilingFile, 1)
	threadcreate := pprof.Lookup("threadcreate")
	threadcreate.WriteTo(threadcreateProfilingFile, 1)
	pprof.StopCPUProfile()
}

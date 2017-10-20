package main

import (
	"fmt"
	"log"
	"os"
	"proxy/services"
	"proxy/utils"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app     *kingpin.Application
	service *services.ServiceItem
)

func initConfig() (err error) {
	//keygen
	if len(os.Args) > 1 {
		if os.Args[1] == "keygen" {
			utils.Keygen()
			os.Exit(0)
		}
	}

	//define  args
	tcpArgs := services.TCPArgs{}
	httpArgs := services.HTTPArgs{}
	tunnelServerArgs := services.TunnelServerArgs{}
	tunnelClientArgs := services.TunnelClientArgs{}
	tunnelBridgeArgs := services.TunnelBridgeArgs{}
	udpArgs := services.UDPArgs{}
	socksArgs := services.SocksArgs{}
	//build srvice args
	app = kingpin.New("proxy", "happy with proxy")
	app.Author("snail").Version(APP_VERSION)
	debug := app.Flag("debug", "debug log output").Default("false").Bool()
	//########http#########
	http := app.Command("http", "proxy on http mode")
	httpArgs.Parent = http.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	httpArgs.CertFile = http.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	httpArgs.KeyFile = http.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	httpArgs.LocalType = http.Flag("local-type", "local protocol type <tls|tcp>").Default("tcp").Short('t').Enum("tls", "tcp")
	httpArgs.ParentType = http.Flag("parent-type", "parent protocol type <tls|tcp|ssh>").Short('T').Enum("tls", "tcp", "ssh")
	httpArgs.Always = http.Flag("always", "always use parent proxy").Default("false").Bool()
	httpArgs.Timeout = http.Flag("timeout", "tcp timeout milliseconds when connect to real server or parent proxy").Default("2000").Int()
	httpArgs.HTTPTimeout = http.Flag("http-timeout", "check domain if blocked , http request timeout milliseconds when connect to host").Default("3000").Int()
	httpArgs.Interval = http.Flag("interval", "check domain if blocked every interval seconds").Default("10").Int()
	httpArgs.Blocked = http.Flag("blocked", "blocked domain file , one domain each line").Default("blocked").Short('b').String()
	httpArgs.Direct = http.Flag("direct", "direct domain file , one domain each line").Default("direct").Short('d').String()
	httpArgs.AuthFile = http.Flag("auth-file", "http basic auth file,\"username:password\" each line in file").Short('F').String()
	httpArgs.Auth = http.Flag("auth", "http basic auth username and password, mutiple user repeat -a ,such as: -a user1:pass1 -a user2:pass2").Short('a').Strings()
	httpArgs.PoolSize = http.Flag("pool-size", "conn pool size , which connect to parent proxy, zero: means turn off pool").Short('L').Default("0").Int()
	httpArgs.CheckParentInterval = http.Flag("check-parent-interval", "check if proxy is okay every interval seconds,zero: means no check").Short('I').Default("3").Int()
	httpArgs.Local = http.Flag("local", "local ip:port to listen").Short('p').Default(":33080").String()
	httpArgs.SSHUser = http.Flag("ssh-user", "user for ssh").Short('u').Default("").String()
	httpArgs.SSHKeyFile = http.Flag("ssh-key", "private key file for ssh").Short('S').Default("").String()
	httpArgs.SSHKeyFileSalt = http.Flag("ssh-keysalt", "salt of ssh private key").Short('s').Default("").String()
	httpArgs.SSHPassword = http.Flag("ssh-password", "password for ssh").Short('A').Default("").String()

	//########tcp#########
	tcp := app.Command("tcp", "proxy on tcp mode")
	tcpArgs.Parent = tcp.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	tcpArgs.CertFile = tcp.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	tcpArgs.KeyFile = tcp.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	tcpArgs.Timeout = tcp.Flag("timeout", "tcp timeout milliseconds when connect to real server or parent proxy").Short('t').Default("2000").Int()
	tcpArgs.ParentType = tcp.Flag("parent-type", "parent protocol type <tls|tcp|udp>").Short('T').Enum("tls", "tcp", "udp")
	tcpArgs.IsTLS = tcp.Flag("tls", "proxy on tls mode").Default("false").Bool()
	tcpArgs.PoolSize = tcp.Flag("pool-size", "conn pool size , which connect to parent proxy, zero: means turn off pool").Short('L').Default("0").Int()
	tcpArgs.CheckParentInterval = tcp.Flag("check-parent-interval", "check if proxy is okay every interval seconds,zero: means no check").Short('I').Default("3").Int()
	tcpArgs.Local = tcp.Flag("local", "local ip:port to listen").Short('p').Default(":33080").String()

	//########udp#########
	udp := app.Command("udp", "proxy on udp mode")
	udpArgs.Parent = udp.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	udpArgs.CertFile = udp.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	udpArgs.KeyFile = udp.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	udpArgs.Timeout = udp.Flag("timeout", "tcp timeout milliseconds when connect to parent proxy").Short('t').Default("2000").Int()
	udpArgs.ParentType = udp.Flag("parent-type", "parent protocol type <tls|tcp|udp>").Short('T').Enum("tls", "tcp", "udp")
	udpArgs.PoolSize = udp.Flag("pool-size", "conn pool size , which connect to parent proxy, zero: means turn off pool").Short('L').Default("0").Int()
	udpArgs.CheckParentInterval = udp.Flag("check-parent-interval", "check if proxy is okay every interval seconds,zero: means no check").Short('I').Default("3").Int()
	udpArgs.Local = udp.Flag("local", "local ip:port to listen").Short('p').Default(":33080").String()

	//########tunnel-server#########
	tunnelServer := app.Command("tserver", "proxy on tunnel server mode")
	tunnelServerArgs.Parent = tunnelServer.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	tunnelServerArgs.CertFile = tunnelServer.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	tunnelServerArgs.KeyFile = tunnelServer.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	tunnelServerArgs.Timeout = tunnelServer.Flag("timeout", "tcp timeout with milliseconds").Short('t').Default("2000").Int()
	tunnelServerArgs.IsUDP = tunnelServer.Flag("udp", "proxy on udp tunnel server mode").Default("false").Bool()
	tunnelServerArgs.Key = tunnelServer.Flag("k", "client key").Default("default").String()
	tunnelServerArgs.Route = tunnelServer.Flag("route", "local route to client's network, such as :PROTOCOL://LOCAL_IP:LOCAL_PORT@[CLIENT_KEY]CLIENT_LOCAL_HOST:CLIENT_LOCAL_PORT").Short('r').Default("").Strings()

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
	socksArgs.ParentType = socks.Flag("parent-type", "parent protocol type <tls|tcp|ssh>").Default("tcp").Short('T').Enum("tls", "tcp", "ssh")
	socksArgs.LocalType = socks.Flag("local-type", "local protocol type <tls|tcp>").Default("tcp").Short('t').Enum("tls", "tcp")
	socksArgs.Local = socks.Flag("local", "local ip:port to listen").Short('p').Default(":33080").String()
	socksArgs.CertFile = socks.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
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
	//parse args
	serviceName := kingpin.MustParse(app.Parse(os.Args[1:]))
	flags := log.Ldate
	if *debug {
		flags |= log.Lshortfile | log.Lmicroseconds
	} else {
		flags |= log.Ltime
	}
	log.SetFlags(flags)
	poster()
	//regist services and run service
	services.Regist("http", services.NewHTTP(), httpArgs)
	services.Regist("tcp", services.NewTCP(), tcpArgs)
	services.Regist("udp", services.NewUDP(), udpArgs)
	services.Regist("tserver", services.NewTunnelServerManager(), tunnelServerArgs)
	services.Regist("tclient", services.NewTunnelClient(), tunnelClientArgs)
	services.Regist("tbridge", services.NewTunnelBridge(), tunnelBridgeArgs)
	services.Regist("socks", services.NewSocks(), socksArgs)
	service, err = services.Run(serviceName)
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

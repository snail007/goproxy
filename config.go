package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"proxy/services"
	"proxy/utils"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app     *kingpin.Application
	service services.ServiceItem
)

func initConfig() (err error) {
	args := services.Args{}
	//define  args
	tcpArgs := services.TCPArgs{}
	httpArgs := services.HTTPArgs{}
	tlsArgs := services.TLSArgs{}
	udpArgs := services.UDPArgs{}

	//build srvice args
	app = kingpin.New("proxy", "happy with proxy")
	app.Author("snail").Version(APP_VERSION)
	args.Parent = app.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	args.Local = app.Flag("local", "local ip:port to listen").Short('p').Default(":33080").String()
	certTLS := app.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	keyTLS := app.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	args.PoolSize = app.Flag("pool-size", "conn pool size , which connect to parent proxy, zero: means turn off pool").Default("50").Int()
	args.CheckParentInterval = app.Flag("check-parent-interval", "check if proxy is okay every interval seconds,zero: means no check").Default("3").Int()

	//########http#########
	http := app.Command("http", "proxy on http mode")
	httpArgs.LocalType = http.Flag("local-type", "parent protocol type <tls|tcp>").Default("tcp").Short('t').Enum("tls", "tcp")
	httpArgs.ParentType = http.Flag("parent-type", "parent protocol type <tls|tcp>").Short('T').Enum("tls", "tcp")
	httpArgs.Always = http.Flag("always", "always use parent proxy").Default("false").Bool()
	httpArgs.Timeout = http.Flag("timeout", "tcp timeout milliseconds when connect to real server or parent proxy").Default("2000").Int()
	httpArgs.HTTPTimeout = http.Flag("http-timeout", "check domain if blocked , http request timeout milliseconds when connect to host").Default("3000").Int()
	httpArgs.Interval = http.Flag("interval", "check domain if blocked every interval seconds").Default("10").Int()
	httpArgs.Blocked = http.Flag("blocked", "blocked domain file , one domain each line").Default("blocked").Short('b').String()
	httpArgs.Direct = http.Flag("direct", "direct domain file , one domain each line").Default("direct").Short('d').String()
	httpArgs.AuthFile = http.Flag("auth-file", "http basic auth file,\"username:password\" each line in file").Short('F').String()
	httpArgs.Auth = http.Flag("auth", "http basic auth username and password, mutiple user repeat -a ,such as: -a user1:pass1 -a user2:pass2").Short('a').Strings()

	//########tcp#########
	tcp := app.Command("tcp", "proxy on tcp mode")
	tcpArgs.Timeout = tcp.Flag("timeout", "tcp timeout milliseconds when connect to real server or parent proxy").Default("2000").Int()
	tcpArgs.ParentType = tcp.Flag("parent-type", "parent protocol type <tls|tcp|udp>").Short('T').Enum("tls", "tcp", "udp")

	//########tls#########
	tls := app.Command("tls", "proxy on tls mode")
	tlsArgs.Timeout = tls.Flag("timeout", "tcp timeout milliseconds when connect to real server or parent proxy").Default("2000").Int()
	tlsArgs.ParentType = tls.Flag("parent-type", "parent protocol type <tls|tcp|udp>").Short('T').Enum("tls", "tcp", "udp")

	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *certTLS != "" && *keyTLS != "" {
		args.CertBytes, args.KeyBytes = tlsBytes(*certTLS, *keyTLS)
	}
	httpArgs.Args = args
	tcpArgs.Args = args
	tlsArgs.Args = args
	udpArgs.Args = args

	//keygen
	if len(os.Args) > 1 {
		if os.Args[1] == "keygen" {
			utils.Keygen()
			os.Exit(0)
		}
	}
	//regist services and run service
	serviceName := kingpin.MustParse(app.Parse(os.Args[1:]))
	services.Regist("http", services.NewHTTP(), httpArgs)
	services.Regist("tcp", services.NewTCP(), tcpArgs)
	services.Regist("tls", services.NewTLS(), tlsArgs)
	services.Regist("udp", services.NewUDP(), udpArgs)
	service, err = services.Run(serviceName)
	if err != nil {
		log.Fatalf("run service [%s] fail, ERR:%s", service, err)
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
func tlsBytes(cert, key string) (certBytes, keyBytes []byte) {
	certBytes, err := ioutil.ReadFile(cert)
	if err != nil {
		log.Fatalf("err : %s", err)
		return
	}
	keyBytes, err = ioutil.ReadFile(key)
	if err != nil {
		log.Fatalf("err : %s", err)
		return
	}
	return
}

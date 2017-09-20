package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg = viper.New()
)

func initConfig() (err error) {
	//define command line args

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	configFile := pflag.StringP("config", "c", "", "config file path")

	pflag.BoolP("parent-tls", "X", false, "parent proxy is tls")
	pflag.BoolP("local-tls", "x", false, "local proxy is tls")
	pflag.BoolP("parent-tcp", "W", false, "parent proxy is tcp")
	pflag.BoolP("local-tcp", "w", true, "local proxy is tcp")
	pflag.BoolP("parent-udp", "U", false, "parent is udp")
	pflag.BoolP("local-udp", "u", false, "local proxy is udp")
	version := pflag.BoolP("version", "v", false, "show version")
	pflag.BoolP("local-http", "z", false, "proxy on http")
	pflag.Bool("always", false, "always use parent proxy")

	pflag.Int("check-proxy-interval", 3, "check if proxy is okay every interval seconds")
	pflag.IntP("port", "p", 33080, "local port to listen")
	pflag.IntP("check-timeout", "t", 3000, "chekc domain blocked , http request timeout milliseconds when connect to host")
	pflag.IntP("tcp-timeout", "T", 2000, "tcp timeout milliseconds when connect to real server or parent proxy")
	pflag.IntP("check-interval", "I", 10, "check domain if blocked every interval seconds")
	pflag.IntP("pool-size", "s", 50, "conn pool size , which connect to parent proxy, zero: means turn off pool")

	pflag.StringP("parent", "P", "", "parent proxy address")
	pflag.StringP("ip", "i", "0.0.0.0", "local ip to bind")
	pflag.StringP("cert", "f", "proxy.crt", "cert file for tls")
	pflag.StringP("key", "k", "proxy.key", "key file for tls")
	pflag.StringP("blocked", "b", "blocked", "blocked domain file , one domain each line")
	pflag.StringP("direct", "d", "direct", "direct domain file , one domain each line")
	pflag.StringP("auth-file", "F", "", "http basic auth file,\"username:password\" each line in file")
	pflag.StringSliceP("auth", "a", []string{}, "http basic auth username and password,such as: \"user1:pass1,user2:pass2\"")

	pflag.Parse()

	cfg.BindPFlag("parent-tls", pflag.Lookup("parent-tls"))
	cfg.BindPFlag("local-tls", pflag.Lookup("local-tls"))
	cfg.BindPFlag("parent-udp", pflag.Lookup("parent-udp"))
	cfg.BindPFlag("local-udp", pflag.Lookup("local-udp"))
	cfg.BindPFlag("parent-tcp", pflag.Lookup("parent-tcp"))
	cfg.BindPFlag("local-tcp", pflag.Lookup("local-tcp"))
	cfg.BindPFlag("local-http", pflag.Lookup("local-http"))
	cfg.BindPFlag("always", pflag.Lookup("always"))
	cfg.BindPFlag("check-proxy-interval", pflag.Lookup("check-proxy-interval"))
	cfg.BindPFlag("port", pflag.Lookup("port"))
	cfg.BindPFlag("check-timeout", pflag.Lookup("check-timeout"))
	cfg.BindPFlag("tcp-timeout", pflag.Lookup("tcp-timeout"))
	cfg.BindPFlag("check-interval", pflag.Lookup("check-interval"))
	cfg.BindPFlag("pool-size", pflag.Lookup("pool-size"))
	cfg.BindPFlag("parent", pflag.Lookup("parent"))
	cfg.BindPFlag("ip", pflag.Lookup("ip"))
	cfg.BindPFlag("cert", pflag.Lookup("cert"))
	cfg.BindPFlag("key", pflag.Lookup("key"))
	cfg.BindPFlag("blocked", pflag.Lookup("blocked"))
	cfg.BindPFlag("direct", pflag.Lookup("direct"))
	cfg.BindPFlag("auth", pflag.Lookup("auth"))
	cfg.BindPFlag("auth-file", pflag.Lookup("auth-file"))

	//version
	if *version {
		fmt.Printf("proxy v%s\n", APP_VERSION)
		os.Exit(0)
	}

	//keygen
	if len(pflag.Args()) > 0 {
		if pflag.Arg(0) == "keygen" {
			keygen()
			os.Exit(0)
		}
	}

	poster()

	if *configFile != "" {
		cfg.SetConfigFile(*configFile)
	} else {
		cfg.SetConfigName("proxy")
		cfg.AddConfigPath("/etc/proxy/")
		cfg.AddConfigPath("$HOME/.proxy")
		cfg.AddConfigPath(".proxy")
		cfg.AddConfigPath(".")
	}

	err = cfg.ReadInConfig()
	file := cfg.ConfigFileUsed()
	if err != nil && !strings.Contains(err.Error(), "Not") {
		log.Fatalf("parse config fail, ERR:%s", err)
	} else if file != "" {
		log.Printf("use config file : %s", file)
	}
	err = nil
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

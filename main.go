package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"bitbucket.org/snail/certverify/check"
	"bitbucket.org/snail/proxy/services"
)

const APP_VERSION = "6.0"

func main() {
	isForever := false
	for _, v := range os.Args[1:] {
		if v == "--forever" {
			isForever = true
		}
	}
	if !isForever {
		check.Init("proxy")
	}
	err := initConfig()
	if err != nil {
		log.Fatalf("err : %s", err)
	}
	if service != nil && service.S != nil {
		Clean(&service.S)
	} else {
		Clean(nil)
	}
}
func Clean(s *services.Service) {
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		for _ = range signalChan {
			log.Println("Received an interrupt, stopping services...")
			if s != nil && *s != nil {
				(*s).Clean()
			}
			if cmd != nil {
				log.Printf("clean process %d", cmd.Process.Pid)
				cmd.Process.Kill()
			}
			if isDebug {
				saveProfiling()
			}
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

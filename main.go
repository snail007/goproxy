package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/snail007/goproxy/services"
)

const APP_VERSION = "5.2"

func main() {
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

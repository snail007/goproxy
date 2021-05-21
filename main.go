package main

import (
	"fmt"
	"github.com/snail007/goproxy/services"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const APP_VERSION = "3.0"

func main() {
	err := initConfig()
	if err != nil {
		log.Fatalf("err : %s", err)
	}
	Clean(&service.S)
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
			fmt.Println("\nReceived an interrupt, stopping services...")
			(*s).Clean()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

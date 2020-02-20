package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nevercase/go-gpt/conf"
	"github.com/nevercase/go-gpt/service/logic"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Fatal(err)
	}
	s := logic.NewService(conf.GetConfig())
	signalHandler(s)
}

func signalHandler(s *logic.Service) {
	var (
		ch = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-ch
		log.Printf("get a signal %s, stop the go-gpt serivce \n", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			s.Close()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

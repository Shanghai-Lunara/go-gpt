package main

import (
	"flag"
	"go-gpt/conf"
	"go-gpt/service/logic"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Fatal(err)
	}
	log.Println(conf.GetConfig())
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

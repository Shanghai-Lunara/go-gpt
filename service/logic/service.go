package logic

import (
	"context"
	"log"
	"os"

	"github.com/Shanghai-Lunara/go-gpt/conf"
)

type Service struct {
	C           *conf.Config
	GitHub      *GitHub
	HttpService *HttpService
	Output      *os.File
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewService(c *conf.Config) *Service {
	file, err := os.OpenFile(c.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panic("OpenFile err: ", err)
	}
	log.SetOutput(file)
	ctx, cancel := context.WithCancel(context.Background())
	s := &Service{
		C:      c,
		ctx:    ctx,
		Output: file,
		cancel: cancel,
	}
	s.GitHub = s.NewGitHub()
	s.HttpService = s.InitHttpServer()
	return s
}

func (s *Service) Close() {
	s.HttpService.ShutDown()
	s.cancel()
	if err := s.Output.Close(); err != nil {
		log.Println("log file close err:", err)
	}
}

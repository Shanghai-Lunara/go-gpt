package logic

import (
	"context"
	"go-gpt/conf"
)

type Service struct {
	C           *conf.Config
	GitHub      *GitHub
	HttpService *HttpService
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewService(c *conf.Config) *Service {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Service{
		C:      c,
		ctx:    ctx,
		cancel: cancel,
	}
	s.GitHub = s.NewGitHub()
	s.HttpService = s.InitHttpServer()
	return s
}

func (s *Service) Close() {
	s.HttpService.ShutDown()
	s.cancel()
}

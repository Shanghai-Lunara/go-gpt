package logic

import (
	"go-gpt/conf"
	"golang.org/x/net/context"
)

type Service struct {
	C           *conf.Config
	GitHub      *GitHub
	HttpService *HttpService
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewService(c *conf.Config) *Service {
	s := &Service{
		C:           c,
		HttpService: InitHttpServer(c),
	}
	s.GitHub = s.NewGitHub()
	return s
}

func (s *Service) Close() {
	s.HttpService.ShutDown()
	s.cancel()
}

package logic

import (
	"go-gpt/conf"
	"golang.org/x/net/context"
)

type Service struct {
	C      *conf.Config
	GitHub *GitHub
	ctx    context.Context
	cancel context.CancelFunc
}

func NewService(c *conf.Config) *Service {
	s := &Service{
		C: c,
	}
	s.GitHub = s.NewGitHub()
	return s
}

func (s *Service) Close() {}

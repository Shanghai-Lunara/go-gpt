package logic

import (
	"context"
	"os"

	"github.com/Shanghai-Lunara/go-gpt/conf"
	"github.com/Shanghai-Lunara/go-gpt/pkg/operator"
	"k8s.io/klog"
)

type Service struct {
	C           *conf.Config
	Project     operator.Project
	HttpService *HttpService
	Output      *os.File
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewService(c *conf.Config) *Service {
	file, err := os.OpenFile(c.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	klog.SetOutput(file)
	ctx, cancel := context.WithCancel(context.Background())
	s := &Service{
		C:       c,
		Project: operator.NewProject(c.Projects, ctx),
		ctx:     ctx,
		Output:  file,
		cancel:  cancel,
	}
	s.HttpService = s.InitHttpServer()
	return s
}

func (s *Service) Close() {
	s.HttpService.ShutDown()
	s.cancel()
	if err := s.Output.Close(); err != nil {
		klog.V(2).Infof("log file close err:", err)
	}
}

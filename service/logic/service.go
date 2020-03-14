package logic

import (
	"context"
	"os"

	"github.com/Shanghai-Lunara/go-gpt/conf"
	"k8s.io/klog"
)

type Service struct {
	c           *conf.Config
	httpService *HttpService
	writer      *os.File
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
		c:           c,
		httpService: InitHttpServer(c, file, ctx),
		ctx:         ctx,
		writer:      file,
		cancel:      cancel,
	}
	return s
}

func (s *Service) Close() {
	s.httpService.ShutDown()
	s.cancel()
	if err := s.writer.Close(); err != nil {
		klog.V(2).Infof("log file close err:%v", err)
	}
}

package logic

import (
	"context"
	"fmt"
	"k8s.io/klog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HttpService struct {
	server  *http.Server
	router  Router
}

func (s *Service) InitHttpServer() *HttpService {
	h := &HttpService{
		router: NewRouter(s.Project),
	}
	router := gin.New()
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: s.Output}), gin.RecoveryWithWriter(s.Output))
	router.LoadHTMLGlob(fmt.Sprintf("%s/*", s.C.Http.TemplatesPath))
	router.GET(RouteGetGitAll, func(c *gin.Context) {
		h.router.GetGitAll(c)
	})
	router.GET(RouteGitGenerate, func(c *gin.Context) {
		h.router.GitGenerate(c)
	})
	h.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.C.Http.IP, s.C.Http.Port),
		Handler: router,
	}
	go func() {
		if err := h.server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				klog.V(2).Info("Server closed under request")
			} else {
				klog.V(0).Info("Server closed unexpected err:", err)
			}
		}
	}()
	return h
}

func (h *HttpService) ShutDown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.server.Shutdown(ctx); err != nil {
		klog.V(3).Infof("http.Server shutdown err:", err)
	}
}

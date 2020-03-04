package logic

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Shanghai-Lunara/go-gpt/pkg/operator"
	"github.com/gin-gonic/gin"
)

type HttpService struct {
	server  *http.Server
	project *operator.Project
}

func (s *Service) InitHttpServer() *HttpService {
	router := gin.New()
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: s.Output}), gin.RecoveryWithWriter(s.Output))
	router.LoadHTMLGlob(fmt.Sprintf("%s/*", s.C.Http.TemplatesPath))
	router.GET("/", func(c *gin.Context) {
		//all, err := s.GitHub.handleAll()
		//if err != nil {
		//	log.Println("handleAll err:", err)
		//	all = "{}"
		//}
		//c.HTML(http.StatusOK, "all.html", map[string]string{"all": all})
	})
	router.GET("/gen/:name/:branch/:command", func(c *gin.Context) {
		log.Println("params:", c.Params)
		//command := &Command{
		//	projectName: c.Param("name"),
		//	branchName:  c.Param("branch"),
		//	command:     c.Param("command"),
		//}
		//if err := s.GitHub.handleCommand(command); err != nil {
		//	log.Println("s.GitHub.handleCommand err:", err)
		//}
		c.String(http.StatusOK, "success")
	})
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.C.Http.IP, s.C.Http.Port),
		Handler: router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Println("Server closed under request")
			} else {
				log.Fatal("Server closed unexpected err:", err)
			}
		}
	}()
	return &HttpService{
		server: server,
	}
}

func (h *HttpService) ShutDown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.server.Shutdown(ctx); err != nil {
		log.Println("http.Server shutdown err:", err)
	}
}

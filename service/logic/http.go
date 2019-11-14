package logic

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type HttpService struct {
	server *http.Server
}

func (s *Service) InitHttpServer() *HttpService {
	router := gin.Default()
	router.LoadHTMLGlob(fmt.Sprintf("%s/*", s.C.TemplatesPath))
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})
	router.GET("/all", func(c *gin.Context) {
		all, err := s.GitHub.handleAll()
		if err != nil {
			log.Println("handleAll err:", err)
		}
		c.HTML(http.StatusOK, "all.html", map[string]string{"all": all})
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
				log.Fatal("Server closed unexpect")
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

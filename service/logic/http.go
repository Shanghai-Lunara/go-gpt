package logic

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-gpt/conf"
	"log"
	"net/http"
	"time"
)

type HttpService struct {
	server *http.Server
}

func InitHttpServer(c *conf.Config) *HttpService {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", c.Http.IP, c.Http.Port),
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

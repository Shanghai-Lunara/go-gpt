package logic

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"k8s.io/klog"
)

type HttpService struct {
	server *http.Server
	router Router
}

func header() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k, v := range c.Request.Header {
			_ = v
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//  header types
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		c.Next()
	}
}

func (s *Service) InitHttpServer() *HttpService {
	h := &HttpService{
		router: NewRouter(s.Project),
	}
	router := gin.New()
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: s.Output}), gin.RecoveryWithWriter(s.Output))
	router.Use(header())
	router.GET(RouteGetGitAll, func(c *gin.Context) {
		res, err := h.router.GetGitAll()
		if err != nil {
			res = GetQuickErrorResponse(CodeUnknownError)
		}
		c.JSON(http.StatusOK, res)
	})
	router.GET(RouteGitGenerate, func(c *gin.Context) {
		p := &GitGenerateParam{
			ProjectName: c.Param("projectName"),
			BranchName:  c.Param("branchName"),
		}
		klog.V(0).Infof("r:%s p:%v", RouteGitGenerate, p)
		res, err := h.router.GitGenerate(p)
		if err != nil {
			res = GetQuickErrorResponse(CodeUnknownError)
		}
		c.JSON(http.StatusOK, res)
	})
	router.GET(RouteSetGitBranchSvnTag, func(c *gin.Context) {
		p := &SetGitBranchSvnTagParam{
			ProjectName: c.Param("projectName"),
			BranchName:  c.Param("branchName"),
			SvnTag:      c.Param("svnTag"),
		}
		klog.V(0).Infof("r:%s p:%v", RouteSetGitBranchSvnTag, p)
		res, err := h.router.SetGitBranchSvnTag(p)
		if err != nil {
			res = GetQuickErrorResponse(CodeUnknownError)
		}
		c.JSON(http.StatusOK, res)
	})
	router.GET(RouteSvnCommit, func(c *gin.Context) {
		p := &SvnCommitParam{
			ProjectName: c.Param("projectName"),
			BranchName:  c.Param("branchName"),
			SvnMessage:  c.Param("svnMsg"),
		}
		klog.V(0).Infof("r:%s p:%v", RouteSvnCommit, p)
		res, err := h.router.SvnCommit(p)
		if err != nil {
			res = GetQuickErrorResponse(CodeUnknownError)
		}
		c.JSON(http.StatusOK, res)
	})
	router.GET(RouteSvnLog, func(c *gin.Context) {
		i, err := strconv.Atoi(c.Param("logNumber"))
		if err != nil {
			c.JSON(http.StatusOK, GetQuickErrorResponse(CodeUnknownError))
			return
		}
		p := &SvnLogParam{
			ProjectName: c.Param("projectName"),
			LogNumber:   i,
		}
		klog.V(0).Infof("r:%s p:%v", RouteSvnLog, p)
		res, err := h.router.SvnLog(p)
		if err != nil {
			res = GetQuickErrorResponse(CodeUnknownError)
		}
		c.JSON(http.StatusOK, res)
	})
	h.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.C.Http.IP, s.C.Http.Port),
		Handler: router,
	}
	go func() {
		if err := h.server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				klog.V(0).Info("Server closed under request")
			} else {
				klog.V(2).Info("Server closed unexpected err:", err)
			}
		}
	}()
	return h
}

func (h *HttpService) Response(c *gin.Context, res interface{}, err error) {
	if err != nil {
		res = GetQuickErrorResponse(CodeUnknownError)
	}
	c.JSON(http.StatusOK, res)
}

func (h *HttpService) ShutDown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.server.Shutdown(ctx); err != nil {
		klog.V(2).Infof("http.Server shutdown err:", err)
	}
}

package logic

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"k8s.io/klog"
)

type HttpService struct {
	server *http.Server
	router Router
}

func (s *Service) InitHttpServer() *HttpService {
	h := &HttpService{
		router: NewRouter(s.Project),
	}
	router := gin.New()
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: s.Output}), gin.RecoveryWithWriter(s.Output))
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

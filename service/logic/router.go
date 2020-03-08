// Package classification API.
//
// The purpose of this service is to provide an application
// that is using plain go code to define an API
//
//      Host: localhost
//      Version: 0.0.1
//
// swagger:meta
package logic

import (
	"net/http"

	"github.com/Shanghai-Lunara/go-gpt/pkg/operator"
	"github.com/gin-gonic/gin"
	"k8s.io/klog"
)

type Router interface {
	GetGitAll(c *gin.Context)
	GitGenerate(c *gin.Context)
	SetGitBranchSvnTag(projectName, branchName, svnTag string) (res Response, err error)
	SvnCommit(projectName, branchName, svnMessage string) (res Response, err error)
	SvnLog(projectName string) (res Response, err error)
}

const (
	RouteGetGitAll          = "/git/all"
	RouteGitGenerate        = "/git/gen/:projectName/:branchName"
	RouteSetGitBranchSvnTag = "/git/set/:projectName/:branchName/:svnTag"
	RouteSvnCommit          = "/git/gen/:projectName/:branchName/:svnMsg"
	RouteSvnLog             = "/git/log/:projectName"
)

type router struct {
	project operator.Project
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type SwaggerResponse struct {
	// The BaseResponse code
	//
	// Required: true
	// Example: 100001
	Code int `json:"code"`
	// The BaseResponse message
	//
	// Required: true
	// Example: success
	Message string `json:"message"`
}

// CommonResponse
// swagger:response CommonResponse
type CommonResponse struct {
	// CommonResponse
	// in: body
	Body struct {
		SwaggerResponse
	}
}

// GitAllResponse
// swagger:response GitAllResponse
type GitAllResponse struct {
	// The all gits' full info
	// in: body
	Body struct {
		SwaggerResponse
		// The set of all gits
		//
		// Required: true
		// An optional field name to which this validation applies
		Gits map[string]operator.GitInfo `json:"gits"`
	}
}

// swagger:route GET /git/all git all
//
// get all gits' info
//
// This will return all gits' info
//
//     Responses:
//       200: GitAllResponse
func (r *router) GetGitAll(c *gin.Context) {
	res, err := r.project.GetAllGitInfo()
	if err != nil {
		klog.V(2).Info("GetGitAll err:", err)
	}
	c.JSON(http.StatusOK, GetQuickResponse(res))
}

// swagger:parameters genSpecificGit
type GitGenerateParam struct {
	// ProjectName
	//
	// Required: true
	// in: path
	ProjectName string `json:"project_name"`
	// BranchName
	//
	// Required: true
	// in: path
	BranchName string `json:"branch_name"`
}

// swagger:route GET /git/gen/{projectName}/{branchName} git gen genSpecificGit
//
// It would generate code and commit to git with specific projectName and branchName
//
// generate and commit
//
//     Responses:
//       200: CommonResponse
func (r *router) GitGenerate(c *gin.Context) {
	gg := &GitGenerateParam {
		ProjectName: c.Param("projectName"),
		BranchName: c.Param("branchName"),
	}
	klog.V(2).Info("gg:", gg)
	if err := r.project.GitGenerate(gg.ProjectName, gg.BranchName); err != nil {
		klog.V(2).Info("GitGenerate err:", err)
	}
	c.JSON(http.StatusOK, GetQuickResponse(map[string]interface{}{}))
}

// swagger:parameters SetParam
type SetGitBranchSvnTagParam struct {
	// ProjectName
	//
	// Required: true
	// in: path
	ProjectName string `json:"project_name"`
	// BranchName
	//
	// Required: true
	// in: path
	BranchName string `json:"branch_name"`
	// SvnTag
	//
	// Required: true
	// in: path
	SvnTag string `json:"svn_tag"`
}

// swagger:route GET /git/set/{projectName}/{branchName}/{svnTag} git set SetParam
//
// It would set a git branch with the specific tag
//
// set
//
//     Responses:
//       200: CommonResponse
func (r *router) SetGitBranchSvnTag(projectName, branchName, svnTag string) (res Response, err error) {
	return res, nil
}

func (r *router) SvnCommit(projectName, branchName, svnMessage string) (res Response, err error) {
	return res, nil
}

func (r *router) SvnLog(projectName string) (res Response, err error) {
	return res, nil
}

func NewRouter(p operator.Project) Router {
	var r Router = &router{
		project: p,
	}
	return r
}

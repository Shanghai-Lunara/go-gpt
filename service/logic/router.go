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
	"github.com/Shanghai-Lunara/go-gpt/pkg/operator"
	"k8s.io/klog"
)

type Router interface {
	GetGitAll() (res HttpResponse, err error)
	GitGenerate(param *GitGenerateParam) (res HttpResponse, err error)
	SetGitBranchSvnTag(param *SetGitBranchSvnTagParam) (res HttpResponse, err error)
	SvnCommit(param *SvnCommitParam) (res HttpResponse, err error)
	SvnLog(param *SvnLogParam) (res HttpResponse, err error)
}

const (
	RouteGetGitAll          = "/git/all"
	RouteGitGenerate        = "/git/gen/:projectName/:branchName"
	RouteSetGitBranchSvnTag = "/git/set/:projectName/:branchName/:svnTag"
	RouteSvnCommit          = "/svn/commit/:projectName/:branchName/:svnMsg"
	RouteSvnLog             = "/svn/log/:projectName/:logNumber"
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
func (r *router) GetGitAll() (res HttpResponse, err error) {
	ret, err := r.project.GetAllGitInfo()
	if err != nil {
		klog.V(2).Infof("GetGitAll err:", err)
		return res, err
	}
	return GetQuickResponse(ret), nil
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
// It would generate code and commit to git with the specific projectName and branchName
//
// generate and commit
//
//     Responses:
//       200: CommonResponse
func (r *router) GitGenerate(param *GitGenerateParam) (res HttpResponse, err error) {
	if err := r.project.GitGenerate(param.ProjectName, param.BranchName); err != nil {
		klog.V(2).Infof("GitGenerate cmd:%v err:%v", *param, err)
		return res, err
	}
	return GetQuickResponse(map[string]interface{}{}), nil
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
func (r *router) SetGitBranchSvnTag(param *SetGitBranchSvnTagParam) (res HttpResponse, err error) {
	err = r.project.GitSetBranchSvnTag(param.ProjectName, param.BranchName, param.SvnTag)
	if err != nil {
		klog.V(2).Infof("SetGitBranchSvnTag cmd:%v err:%v", *param, err)
		return res, err
	}
	return GetQuickResponse(map[string]interface{}{}), nil
}

// swagger:parameters SetSvnTag
type SvnCommitParam struct {
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
	SvnMessage string `json:"svn_message"`
}

// swagger:route GET /svn/commit/{projectName}/{branchName}/{svnMessage} svn commit SetSvnTag
//
// It would sync project files from the specific git.branch and commit to svn server
//
// scn commit
//
//     Responses:
//       200: CommonResponse
func (r *router) SvnCommit(param *SvnCommitParam) (res HttpResponse, err error) {
	err = r.project.SvnCommit(param.ProjectName, param.BranchName, param.SvnMessage)
	if err != nil {
		klog.V(2).Infof("SvnCommit cmd:%v err:%v", *param, err)
		return res, err
	}
	return GetQuickResponse(map[string]interface{}{}), nil
}

// swagger:parameters SvnLog
type SvnLogParam struct {
	// ProjectName
	//
	// Required: true
	// in: path
	ProjectName string `json:"project_name"`
	// LogNumber
	//
	// Required: true
	// in: path
	LogNumber int `json:"log_number"`
}

// SvnLogResponse
// swagger:response SvnLogResponse
type SvnLogResponse struct {
	// The svn logs
	// in: body
	Body struct {
		SwaggerResponse
		// The set of svn logs
		//
		// Required: true
		// An optional field name to which this validation applies
		Logentrys []operator.Logentry `json:"logentrys"`
	}
}

// swagger:route GET /svn/log/{projectName}/{logNumber} svn log SvnLog
//
// It would pull svn logs from the remote svn server with the specific number
//
// svn log
//
//     Responses:
//       200: SvnLogResponse
func (r *router) SvnLog(param *SvnLogParam) (res HttpResponse, err error) {
	ret, err := r.project.SvnLog(param.ProjectName, param.LogNumber)
	if err != nil {
		klog.V(2).Infof("SvnLog cmd:%v err:%v", *param, err)
		return res, err
	}
	return GetQuickResponse(ret), nil
}

func NewRouter(p operator.Project) Router {
	var r Router = &router{
		project: p,
	}
	return r
}

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
	GitGenerate(param *GitGenerateParam) (res HttpResponse, err error) // async
	SetGitBranchSvnTag(param *SetGitBranchSvnTagParam) (res HttpResponse, err error)
	SvnCommit(param *SvnCommitParam) (res HttpResponse, err error) // async
	SvnLog(param *SvnLogParam) (res HttpResponse, err error)
	FtpLog(param *FtpLogParam) (res HttpResponse, err error)
	FtpReadFile(param *FtpReadFileParam) (res HttpResponse, err error)
	FtpWriteFile(param *FtpWriteFileParam) (res HttpResponse, err error)
	FtpCompress(param *FtpCompressParam) (res HttpResponse, err error)
	TaskAll(param *TaskAllParam) (res HttpResponse, err error)
	OssEnvs(param *OssEnvsParam) (res HttpResponse, err error)
	OssContent(param *OssContentParam) (res HttpResponse, err error)
	OssUpdate(param *OssUpdateParam) (res HttpResponse, err error)
}

const (
	RouteGetGitAll          = "/git/all"
	RouteGitGenerate        = "/git/gen/:projectName/:branchName"
	RouteSetGitBranchSvnTag = "/git/set/:projectName/:branchName/:svnTag"
	RouteSvnCommit          = "/svn/commit/:projectName/:branchName/:svnMsg"
	RouteSvnLog             = "/svn/log/:projectName/:logNumber"
	RouteFtpLog             = "/ftp/log/:projectName/:filter"
	RouteFtpReadFile        = "/ftp/read/:projectName/:fileName"
	RouteFtpWriteFile       = "/ftp/write"
	RouteFtpCompress        = "/ftp/compress/:projectName/:branchName/:zipType/:zipFlags"
	RouteTaskAll            = "/task/all/:projectName"
	RouteOssEnvs            = "/oss/envs/:projectName"
	RouteOssContent         = "/oss/content/:projectName/:env"
	RouteOssUpdate          = "/oss/update"
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
	c := &operator.Command{
		ProjectName: param.ProjectName,
		BranchName:  param.BranchName,
		Command:     operator.TaskCmdGitGen,
	}
	if err := r.project.AsyncTask(c); err != nil {
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
	c := &operator.Command{
		ProjectName: param.ProjectName,
		BranchName:  param.BranchName,
		Command:     operator.TaskCmdSvnCommit,
		Message:     param.SvnMessage,
	}
	if err = r.project.AsyncTask(c); err != nil {
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

// swagger:parameters FtpLog
type FtpLogParam struct {
	// ProjectName
	//
	// Required: true
	// in: path
	ProjectName string `json:"project_name"`
	// Filter
	//
	// Required: true
	// in: path
	Filter string `json:"filter"`
}

// FtpLogResponse
// swagger:response FtpLogResponse
type FtpLogResponse struct {
	// The svn logs
	// in: body
	Body struct {
		SwaggerResponse
		// The list of ftp files
		//
		// Required: true
		// An optional field name to which this validation applies
		Entries []operator.Entry `json:"entries"`
	}
}

// swagger:route GET /ftp/log/{projectName}/{filter} ftp log FtpLog
//
// It would pull ftp files from the remote ftp server with the specific filter
//
// ftp log
//
//     Responses:
//       200: FtpLogResponse
func (r *router) FtpLog(param *FtpLogParam) (res HttpResponse, err error) {
	ret, err := r.project.FtpLog(param.ProjectName, param.Filter)
	if err != nil {
		klog.V(2).Infof("FtpLog cmd:%v err:%v", *param, err)
		return res, err
	}
	return GetQuickResponse(ret), nil
}

// swagger:parameters FtpReadFile
type FtpReadFileParam struct {
	// ProjectName
	//
	// Required: true
	// in: path
	ProjectName string `json:"project_name"`
	// FileName
	//
	// Required: true
	// in: path
	FileName string `json:"file_name"`
}

// FtpLogResponse
// swagger:response FtpReadFileResponse
type FtpReadFileResponse struct {
	// The svn logs
	// in: body
	Body struct {
		SwaggerResponse
		// The file content
		//
		// Required: true
		// An optional field name to which this validation applies
		Content string `json:"content"`
	}
}

// swagger:route GET /ftp/read/{projectName}/{fileName} ftp read FtpReadFile
//
// It would get the content of the specific file from the remote ftp server by the filter
//
// ftp read
//
//     Responses:
//       200: FtpReadFileResponse
func (r *router) FtpReadFile(param *FtpReadFileParam) (res HttpResponse, err error) {
	ret, err := r.project.FtpReadFile(param.ProjectName, param.FileName)
	if err != nil {
		klog.V(2).Infof("FtpReadFile cmd:%v err:%v", *param, err)
		return res, err
	}
	return GetQuickResponse(string(ret)), nil
}

// swagger:parameters FtpWriteFile
type FtpWriteFileParam struct {
	// ProjectName
	//
	// Required: true
	// in: path
	ProjectName string `json:"projectName"`
	// FileName
	//
	// Required: true
	// in: path
	FileName string `json:"fileName"`
	// Content
	//
	// Required: true
	// in: path
	Content string `json:"content"`
}

// swagger:route POST /ftp/write ftp write FtpWriteFile
//
// It would overwrite the specific file on the FTP server with the provided content
//
// ftp write
//
//     Responses:
//       200: CommonResponse
func (r *router) FtpWriteFile(param *FtpWriteFileParam) (res HttpResponse, err error) {
	err = r.project.FtpWriteFile(param.ProjectName, param.FileName, param.Content)
	if err != nil {
		klog.V(2).Infof("FtpWriteFile cmd:%v err:%v", *param, err)
		return res, err
	}
	return GetQuickResponse(map[string]interface{}{}), nil
}

// swagger:parameters FtpCompress
type FtpCompressParam struct {
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
	// ZipType
	//
	// Required: true
	// in: path
	ZipType string `json:"zip_type"`
	// ZipFlags
	//
	// Required: true
	// in: path
	ZipFlags string `json:"zip_flags"`
}

// swagger:route GET /ftp/compress/{projectName}/{branchName}/{zipType}/{zipFlags} ftp compress FtpCompress
//
// It would compress the project into zip with the specific flags and upload to the ftp server
//
// ftp compress
//
//     Responses:
//       200: CommonResponse
func (r *router) FtpCompress(param *FtpCompressParam) (res HttpResponse, err error) {
	c := &operator.Command{
		ProjectName: param.ProjectName,
		BranchName:  param.BranchName,
		Command:     operator.TaskCmdFtpUpload,
		ZipType:     param.ZipType,
		ZipFlags:    param.ZipFlags,
	}
	if err = r.project.AsyncTask(c); err != nil {
		klog.V(2).Infof("FtpCompress cmd:%v err:%v", *param, err)
		return res, err
	}
	return GetQuickResponse(map[string]interface{}{}), nil
}

// swagger:parameters TaskAll
type TaskAllParam struct {
	// ProjectName
	//
	// Required: true
	// in: path
	ProjectName string `json:"project_name"`
}

// FtpLogResponse
// swagger:response TaskAllResponse
type TaskAllResponse struct {
	// The svn logs
	// in: body
	Body struct {
		SwaggerResponse
		// all tasks
		//
		// Required: true
		// An optional field name to which this validation applies
		Tasks map[int]operator.Task `json:"tasks"`
	}
}

// swagger:route GET /task/all/{projectName} task all TaskAll
//
// It would get all the tasks of the specific project
//
// task all
//
//     Responses:
//       200: TaskAllResponse
func (r *router) TaskAll(param *TaskAllParam) (res HttpResponse, err error) {
	ret, err := r.project.TaskAll(param.ProjectName)
	if err != nil {
		klog.V(2).Infof("TaskAll cmd:%v err:%v", *param, err)
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

// swagger:parameters OssEnvs
type OssEnvsParam struct {
	// ProjectName
	//
	// Required: true
	// in: path
	ProjectName string `json:"project_name"`
}

// OssEnvsResponse
// swagger:response OssEnvsResponse
type OssEnvsResponse struct {
	// The all gits' full info
	// in: body
	Body struct {
		SwaggerResponse
		// envs
		//
		// Required: true
		// An optional field name to which this validation applies
		Envs map[string]string `json:"envs"`
	}
}

// swagger:route GET /oss/envs/{projectName} oss envs OssEnvs
//
// get oss envs
//
// This will return the specific project's oss environments
//
//     Responses:
//       200: OssEnvsResponse
func (r *router) OssEnvs(param *OssEnvsParam) (res HttpResponse, err error) {
	ret, err := r.project.OssEnvs(param.ProjectName)
	if err != nil {
		klog.V(2).Infof("OssEnvs err:", err)
		return res, err
	}
	return GetQuickResponse(ret), nil
}

// swagger:parameters OssContent
type OssContentParam struct {
	// ProjectName
	//
	// Required: true
	// in: path
	ProjectName string `json:"project_name"`
	// Env
	//
	// Required: true
	// in: path
	Env string `json:"env"`
}

// OssContentResponse
// swagger:response OssContentResponse
type OssContentResponse struct {
	// The all gits' full info
	// in: body
	Body struct {
		SwaggerResponse
		// content
		//
		// Required: true
		// An optional field name to which this validation applies
		NoticeContent operator.NoticeContent `json:"notice_content"`
	}
}

// swagger:route GET /oss/content/{projectName}/{env} oss envs OssContent
//
// get oss envs
//
// This will return the specific object's content in oss
//
//     Responses:
//       200: OssContentResponse
func (r *router) OssContent(param *OssContentParam) (res HttpResponse, err error) {
	ret, err := r.project.OssContent(param.ProjectName, param.Env)
	if err != nil {
		klog.V(2).Infof("OssContent err:", err)
		return res, err
	}
	return GetQuickResponse(ret), nil
}

// swagger:parameters OssUpdate
type OssUpdateParam struct {
	// ProjectName
	//
	// Required: true
	// in: path
	ProjectName string `json:"projectName"`
	// Env
	//
	// Required: true
	// in: path
	Env string `json:"env"`
	// Title
	//
	// Required: true
	// in: Title
	Title string `json:"title"`
	// Time
	//
	// Required: true
	// in: Time
	Time string `json:"time"`
	// Content
	//
	// Required: true
	// in: Content
	Content string `json:"content"`
}

// swagger:route POST /oss/update oss update OssUpdate
//
// It would overwrite the specific file on the oss server with the provided content
//
// ftp write
//
//     Responses:
//       200: CommonResponse
func (r *router) OssUpdate(param *OssUpdateParam) (res HttpResponse, err error) {
	nc := operator.NoticeContent{
		Title:   param.Title,
		Content: param.Content,
		Time:    param.Time,
	}
	err = r.project.OssUpdateContent(param.ProjectName, param.Env, nc)
	if err != nil {
		klog.V(2).Infof("OssUpdate cmd:%v err:%v", *param, err)
		return res, err
	}
	return GetQuickResponse(map[string]interface{}{}), nil
}

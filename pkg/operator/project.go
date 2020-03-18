package operator

import (
	"context"
	"errors"
	"fmt"
	"time"

	"k8s.io/klog"
)

type Project interface {
	Add(projectName string, p *project)
	GetProject(projectName string) (p *project, err error)
	GetGitInfo(projectName string) (gi GitInfo, err error)
	GetAllGitInfo() (res map[string]GitInfo, err error)
	GitGenerate(projectName, branchName string) error // needed async
	GitSetBranchSvnTag(projectName, branchName, svnTag string) error
	SvnCommit(projectName, branchName, svnMessage string) error // needed async
	SvnLog(projectName string, showNumber int) (res []Logentry, err error)
	FtpLog(projectName, filter string) (res []Entry, err error)
	FtpReadFile(projectName, fileName string) (res []byte, err error)
	FtpWriteFile(projectName, fileName, content string) error
	FtpCompress(projectName, branchName, zipType, zipFlags string) error // needed async
	AsyncTask(c *Command) error
	TaskAll(projectName string) (res map[int]Task, err error)
}

const (
	errNotExistedProject = "the project: %s is not existed"
)

const (
	ZipTypeAll   = "ser"
	ZipTypePatch = "pat"
)

const (
	introduceTemplate = "introduce_%s.txt"
	serverTemplate    = "HelixServer_%s.zip"
	serverMd5Template = "HelixServer_%s.zip.txt"
	versionTemplate   = "%s%s"
	patchTemplate     = "patch_%s"
)

type project struct {
	git GitOperator
	svn SvnOperator
	ftp FtpOperator

	worker Worker
	tasks  *TaskHub

	ctx    context.Context
	cancel context.CancelFunc
}

type projects struct {
	projects map[string]*project
}

func (ph *projects) Add(projectName string, p *project) {
	ph.projects[projectName] = p
}

func (ph *projects) GetProject(projectName string) (p *project, err error) {
	if t, ok := ph.projects[projectName]; ok {
		return t, nil
	} else {
		return p, errors.New(fmt.Sprintf(errNotExistedProject, projectName))
	}
}

func (ph *projects) GetGitInfo(projectName string) (gi GitInfo, err error) {
	p, err := ph.GetProject(projectName)
	if err != nil {
		return gi, err
	}
	return p.git.GetGitInfo(), nil
}

func (ph *projects) GetAllGitInfo() (res map[string]GitInfo, err error) {
	res = make(map[string]GitInfo, 0)
	for k, v := range ph.projects {
		_ = v
		if gi, err := ph.GetGitInfo(k); err != nil {
			return res, err
		} else {
			res[k] = gi
		}
	}
	return res, nil
}

func (ph *projects) GitGenerate(projectName, branchName string) error {
	p, err := ph.GetProject(projectName)
	if err != nil {
		return err
	}
	c := &GitCmd{
		cmd:        cmdGitGenerate,
		branchName: branchName,
	}
	return p.git.SendCommand(c)
}

func (ph *projects) GitSetBranchSvnTag(projectName, branchName, svnTag string) error {
	p, err := ph.GetProject(projectName)
	if err != nil {
		return err
	}
	return p.git.SetSvnTag(branchName, svnTag)
}

func (ph *projects) SvnCommit(projectName, branchName, svnMessage string) error {
	p, err := ph.GetProject(projectName)
	if err != nil {
		return err
	}
	p.svn.Lock()
	defer p.svn.Unlock()
	if err := p.git.SvnSync(branchName); err != nil {
		return err
	}
	return p.svn.Commit(svnMessage)
}

func (ph *projects) SvnLog(projectName string, showNumber int) (res []Logentry, err error) {
	p, err := ph.GetProject(projectName)
	if err != nil {
		return res, err
	}
	return p.svn.Log(showNumber)
}

func (ph *projects) FtpLog(projectName, filter string) (res []Entry, err error) {
	p, err := ph.GetProject(projectName)
	if err != nil {
		return res, err
	}
	return p.ftp.List(filter)
}

func (ph *projects) FtpReadFile(projectName, fileName string) (res []byte, err error) {
	p, err := ph.GetProject(projectName)
	if err != nil {
		return res, err
	}
	return p.ftp.ReadFileContent(fileName)
}

func (ph *projects) FtpWriteFile(projectName, fileName, content string) error {
	p, err := ph.GetProject(projectName)
	if err != nil {
		return err
	}
	return p.ftp.WriteFileContent(fileName, []byte(content))
}

func (ph *projects) FtpCompress(projectName, branchName, zipType, zipFlags string) error {
	p, err := ph.GetProject(projectName)
	if err != nil {
		return err
	}
	p.git.RLock()
	defer p.git.RUnlock()
	version, err := p.ftp.GetNextVersion()
	if err != nil {
		return err
	}
	if err := p.git.FtpCompress(branchName, zipType, version, zipFlags); err != nil {
		return err
	}
	defer func() {
		if err := p.git.Revert(); err != nil {
			klog.V(2).Info(err)
		}
	}()
	v := fmt.Sprintf(versionTemplate, time.Now().Format("20060102"), version)
	introName := fmt.Sprintf(introduceTemplate, v)
	if err := p.ftp.UploadFile(fmt.Sprintf("%s/%s", p.git.Conf().WorkDir, introName), introName); err != nil {
		return err
	}
	switch zipType {
	case ZipTypeAll:
	case ZipTypePatch:
		v = fmt.Sprintf(patchTemplate, v)
	}
	zipName := fmt.Sprintf(serverTemplate, v)
	klog.Info(zipName)
	if err := p.ftp.UploadFile(fmt.Sprintf("%s/%s", p.git.Conf().WorkDir, zipName), zipName); err != nil {
		return err
	}
	zipMd5Name := fmt.Sprintf(serverMd5Template, v)
	klog.Info(zipMd5Name)
	if err := p.ftp.UploadFile(fmt.Sprintf("%s/%s", p.git.Conf().WorkDir, zipMd5Name), zipMd5Name); err != nil {
		return err
	}
	return nil
}

func (ph *projects) AsyncTask(c *Command) error {
	p, err := ph.GetProject(c.ProjectName)
	if err != nil {
		return err
	}
	p.worker.Add(p.tasks.NewTask(c))
	return nil
}

func (ph *projects) TaskAll(projectName string) (res map[int]Task, err error) {
	p, err := ph.GetProject(projectName)
	if err != nil {
		return res, err
	}
	return p.tasks.GetAll(), nil
}

func NewProject(conf []ProjectConfig, ctx context.Context) Project {
	var ph Project = &projects{
		projects: make(map[string]*project, 0),
	}
	for _, v := range conf {
		p := &project{
			git:    NewGitOperator(&v, ctx),
			svn:    NewSvnOperator(&v, ctx),
			ftp:    NewFtpOperator(v.Ftp),
			worker: NewWorker(ctx.Done(), ph),
			tasks:  NewTaskHub(),
			ctx:    ctx,
		}
		ph.Add(v.ProjectName, p)
	}
	return ph
}

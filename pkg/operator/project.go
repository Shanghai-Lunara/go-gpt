package operator

import (
	"context"
	"errors"
	"fmt"
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
}

const (
	errNotExistedProject = "the project: %s is not existed"
)

type project struct {
	git GitOperator
	svn SvnOperator
	ftp FtpOperator

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
	c := &Command{
		projectName: projectName,
		branchName:  branchName,
		command:     cmdGitGenerate,
		message:     "",
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

func NewProject(conf []ProjectConfig, ctx context.Context) Project {
	var ph Project = &projects{
		projects: make(map[string]*project, 0),
	}
	for _, v := range conf {
		p := &project{
			git: NewGitOperator(&v, ctx),
			svn: NewSvnOperator(&v, ctx),
			ftp: NewFtpOperator(v.Ftp),
		}
		ph.Add(v.ProjectName, p)
	}
	return ph
}

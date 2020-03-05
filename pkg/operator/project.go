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
	GitGenerate(projectName, branchName string) error
	SvnCommit(projectName, branchName, svnMessage string) error
	SvnLog(projectName string, number int) (res []Logentry, err error)
}

const (
	errNotExistedProject = "the project: %s is not existed"
)

type project struct {
	Git GitOperator
	Svn SvnOperator

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
	return p.Git.GetGitInfo(), nil
}

func (ph *projects) GitGenerate(projectName, branchName string) error {
	p, err := ph.GetProject(projectName)
	if err != nil {
		return err
	}
	c := &Command{
		projectName: projectName,
		branchName:  branchName,
		command:     "",
		message:     "",
	}
	return p.Git.HandleCommand(c)
}

func (ph *projects) SvnCommit(projectName, branchName, svnMessage string) error {
	p, err := ph.GetProject(projectName)
	if err != nil {
		return err
	}
	p.Svn.Lock()
	defer p.Svn.Unlock()
	if err := p.Git.SvnSync(branchName); err != nil {
		return err
	}
	return p.Svn.Commit(svnMessage)
}

func (ph *projects) SvnLog(projectName string, number int) (res []Logentry, err error) {
	p, err := ph.GetProject(projectName)
	if err != nil {
		return res, err
	}
	return p.Svn.Log(number)
}

func NewProject(conf []ProjectConfig, ctx context.Context) *Project {
	var ph Project = &projects{
		projects: make(map[string]*project, 0),
	}
	for _, v := range conf {
		p := &project{
			Git: NewGitOperator(&v, ctx),
			Svn: NewSvnOperator(&v, ctx),
		}
		ph.Add(v.ProjectName, p)
	}
	return &ph
}

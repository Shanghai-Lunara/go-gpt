package logic

import (
	"context"

	"github.com/Shanghai-Lunara/go-gpt/conf"
)

type ProjectOperator interface {
	GitGenerate() error
	SvnCommit() error
	SvnStatus(number int)
}

type project struct {
	Git *Git
	Svn *SvnOperator

	ctx    context.Context
	cancel context.CancelFunc
}

type ProjectHub struct {
	Projects map[string]*project
}

func NewProjectHub(conf *conf.Config, ctx context.Context) *ProjectHub {
	ph := &ProjectHub{
		Projects: make(map[string]*project, 0),
	}
	for _, v := range conf.Projects {
		p := &project{
			Svn: NewSvnOperator(&v, ctx),
		}
	}
	return ph
}

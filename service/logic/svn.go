package logic

import (
	"context"
	"sync"

	"github.com/Shanghai-Lunara/go-gpt/conf"
)

type SvnHub struct {
	mu   sync.RWMutex
	Svns map[string]*Svn
}

type Svn struct {
	mu sync.RWMutex

	ProjectName string `json:"project_name"`

	Username string `json:"username"`
	Password string `json:"password"`

	WorkDir string `json:"work_dir"`
	Url     string `json:"url"`
	Port    int    `json:"port"`

	ctx context.Context
}

func NewSvnHub(c *conf.Config, ctx context.Context) *SvnHub {
	sh := &SvnHub{
		Svns: make(map[string]*Svn, 0),
	}
	for _, v := range c.Projects {
		svn := &Svn{
			ProjectName: v.ProjectName,
			Username:    v.Svn.Username,
			Password:    v.Svn.Password,
			WorkDir:     v.Svn.WorkDir,
			Url:         v.Svn.Url,
			Port:        v.Svn.Port,
			ctx:         ctx,
		}
		sh.Svns[v.ProjectName] = svn
	}
	return sh
}

func (s *Svn) CheckOut() {

}

func (s *Svn) Update() {

}

func (s *Svn) Add() {

}

func (s *Svn) Commit() {

}

func (s *Svn) Log() {

}

func (s *Svn) LoopChan() {

}

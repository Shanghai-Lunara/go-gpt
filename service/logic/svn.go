package logic

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"sync"

	"github.com/Shanghai-Lunara/go-gpt/conf"
)

type SvnOperator interface {
	CheckOut() error
	Update() error
	Status() error
	AddAll() error
	Clean() error
	Commit(message string) error
	Log(number int) (res []string, err error)
} 

const SvnUrl = "svn://%s@%s:%d/%s"

const (
	scriptName = "svn.sh"

	cmdCheckOut = "checkout"
	cmdUpdate   = "update"
	cmdStatus = "status"
	cmdAddAll = "addAll"
	cmdClean = "clean"
	cmdCommit = "commit"
)

type SvnHub struct {
	mu   sync.RWMutex
	Svns map[string]*SvnOperator
}

type Svn struct {
	mu sync.RWMutex

	ScriptPath  string `json:"script_path"`
	ProjectName string `json:"project_name"`

	Username string `json:"username"`
	Password string `json:"password"`

	WorkDir   string `json:"work_dir"`
	Url       string `json:"url"`
	Port      int    `json:"port"`
	RemoteDir string `json:"remote_dir"`
	SvnUrl    string `json:"svn_url"`

	ctx context.Context
}

func NewSvnHub(c *conf.Config, ctx context.Context) *SvnHub {
	sh := &SvnHub{
		Svns: make(map[string]*SvnOperator, 0),
	}
	for _, v := range c.Projects {
		var svn SvnOperator = &Svn{
			ScriptPath:  fmt.Sprintf("%s%s", v.ScriptsPath, scriptName),
			ProjectName: v.ProjectName,
			Username:    v.Svn.Username,
			Password:    v.Svn.Password,
			WorkDir:     v.Svn.WorkDir,
			Url:         v.Svn.Url,
			Port:        v.Svn.Port,
			RemoteDir:   v.Svn.RemoteDir,
			SvnUrl:      fmt.Sprintf(SvnUrl, v.Svn.Username, v.Svn.Url, v.Svn.Port, v.Svn.RemoteDir),
			ctx:         ctx,
		}
		sh.Svns[v.ProjectName] = &svn
		if err := svn.CheckOut(); err != nil {
			log.Println(err)
		}
		if err := svn.Update(); err != nil {
			log.Println(err)
		}
		if err := svn.Status(); err != nil {
			log.Println(err)
		}
		if err := svn.Clean(); err != nil {
			log.Println(err)
		}
		if err := svn.Status(); err != nil {
			log.Println(err)
		}
	}
	return sh
}

func (s *Svn) CheckOut() error {
	_, err := s.ExecuteWithArgs(cmdCheckOut, s.SvnUrl)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svn) Update() error {
	_, err := s.ExecuteWithArgs(cmdUpdate)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svn) Status() error {
	_, err := s.ExecuteWithArgs(cmdStatus)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svn) AddAll() error {
	_, err := s.ExecuteWithArgs(cmdAddAll)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svn) Clean() error {
	_, err := s.ExecuteWithArgs(cmdClean)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svn) Commit(message string) error {
	_, err := s.ExecuteWithArgs(cmdCommit, message)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svn) Log(number int) (res []string, err error) {

	return res, err
}

func (s *Svn) ExecuteWithArgs(args ...string) (res []byte,err error) {
	t := append([]string{s.ScriptPath, s.Username, s.Password, s.WorkDir, s.RemoteDir}, args...)
	out, err := exec.Command("sh", t...).Output()
	if err != nil {
		return out, errors.New(fmt.Sprintf("Svn %s exec.Command err:%v\n", args[0], err))
	}
	log.Printf("Svn Command `%s` output:\n%s\n", args[0], string(out))
	return out, err
}

func (s *Svn) LoopChan() {

}

func (s *Svn) handle() {

}

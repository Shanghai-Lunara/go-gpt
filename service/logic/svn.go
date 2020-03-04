package logic

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/Shanghai-Lunara/go-gpt/conf"
)

type SvnOperator interface {
	CheckOut() error
	Update() error
	Status() error
	AddAll() error
	Clean() error
	Commit(handleFunc func() error, message string) error
	Log(number int) (res []Logentry, err error)
	ExecuteWithArgs(args ...string) (res []byte, err error)
	Timer()
	Listener(ch chan *Command)
}

const SvnUrl = "svn://%s@%s:%d/%s"

const (
	scriptName = "svn.sh"

	cmdCheckOut = "checkout"
	cmdUpdate   = "update"
	cmdStatus   = "status"
	cmdAddAll   = "addAll"
	cmdClean    = "clean"
	cmdCommit   = "commit"
	cmdLog      = "log"
)

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

func NewSvnOperator(v *conf.Project, ctx context.Context) *SvnOperator {
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
	if err := svn.CheckOut(); err != nil {
		log.Println(err)
	}
	if err := svn.Update(); err != nil {
		log.Println(err)
	}
	return &svn
}

func (s *Svn) CheckOut() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.ExecuteWithArgs(cmdCheckOut, s.SvnUrl)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svn) Update() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.ExecuteWithArgs(cmdUpdate)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svn) Status() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, err := s.ExecuteWithArgs(cmdStatus)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svn) AddAll() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.ExecuteWithArgs(cmdAddAll)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svn) Clean() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.ExecuteWithArgs(cmdClean)
	if err != nil {
		return err
	}
	return nil
}

func (s *Svn) Commit(handleFunc func() error, message string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := handleFunc(); err != nil {
		return errors.New(fmt.Sprintf("Svn commit handleFunc err:%v", err))
	}
	_, err := s.ExecuteWithArgs(cmdCommit, message)
	if err != nil {
		return err
	}
	return nil
}

type LogResponse struct {
	XMLName   xml.Name   `xml:"log"`
	Logentrys []Logentry `xml:"logentry" json:"logentrys"`
}

type Logentry struct {
	Revision string    `xml:"revision,attr" json:"revision,omitempty"`
	Author   string    `xml:"author" json:"author,omitempty"`
	DateTime time.Time `xml:"date" json:"date_time,omitempty"`
	Msg      string    `xml:"msg" json:"msg,omitempty"`
	Paths    []Path    `xml:"paths>path" json:"paths,omitempty"`
}

type Path struct {
	Action   string `xml:"action,attr" json:"action,omitempty"`
	PropMods string `xml:"prop-mods,attr" json:"prop_mods,omitempty"`
	TextMods string `xml:"text-mods,attr" json:"text_mods,omitempty"`
	Kind     string `xml:"kind,attr" json:"kind,omitempty"`
	Value    string `xml:",chardata" json:"value,omitempty"`
}

func (s *Svn) Log(number int) (res []Logentry, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out, err := s.ExecuteWithArgs(cmdLog, strconv.Itoa(number))
	if err != nil {
		return res, err
	}
	rest := LogResponse{}
	if err := xml.Unmarshal(out, &rest); err != nil {
		return res, err
	}
	return rest.Logentrys, nil
}

func (s *Svn) ExecuteWithArgs(args ...string) (res []byte, err error) {
	t := append([]string{s.ScriptPath, s.Username, s.Password, s.WorkDir, s.RemoteDir}, args...)
	out, err := exec.Command("sh", t...).Output()
	if err != nil {
		return out, errors.New(fmt.Sprintf("Svn %s exec.Command err:%v\n", args[0], err))
	}
	log.Printf("Svn Command `%s` output:\n%s\n", args[0], string(out))
	return out, nil
}

func (s *Svn) Timer() {
	tick := time.NewTicker(time.Second * 10)
	defer tick.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-tick.C:
			if err := s.Update(); err != nil {
				log.Println(err)
			}
		}
	}
}

func (s *Svn) Listener(ch chan *Command) {
	defer close(ch)
	for {
		select {
		case c, isClose := <-ch:
			if !isClose {
				return
			}
			log.Println("cmd:", *c)
			if _, err := s.ExecuteWithArgs(c.command, c.message); err != nil {
				log.Println("Listener exc err:", err)
			}
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Svn) handle() {

}

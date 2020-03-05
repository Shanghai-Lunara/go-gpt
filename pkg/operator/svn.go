package operator

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
)

type SvnOperator interface {
	Lock()
	Unlock()
	ExecuteWithArgs(args ...string) (res []byte, err error)
	CheckOut() error
	Update() error
	Status() error
	AddAll() error
	Clean() error
	Commit(svnMessage string) error
	Log(number int) (res []Logentry, err error)
	Timer()
	Listener(ch chan *Command)
}

const svnUrl = "svn://%s@%s:%d/%s"

const (
	svnScriptName = "svn.sh"

	cmdCheckOut = "checkout"
	cmdUpdate   = "update"
	cmdStatus   = "status"
	cmdAddAll   = "addAll"
	cmdClean    = "clean"
	cmdCommit   = "commit"
	cmdLog      = "log"
)

type svn struct {
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

func (s *svn) Lock() {
	s.mu.Lock()
}

func (s *svn) Unlock() {
	s.mu.Unlock()
}

func (s *svn) ExecuteWithArgs(args ...string) (res []byte, err error) {
	t := append([]string{s.ScriptPath, s.Username, s.Password, s.WorkDir, s.RemoteDir}, args...)
	out, err := exec.Command("sh", t...).Output()
	if err != nil {
		return out, errors.New(fmt.Sprintf("Svn %s exec.Command err:%v\n", args[0], err))
	}
	log.Printf("Svn Command `%s` output:\n%s\n", args[0], string(out))
	return out, nil
}

func (s *svn) CheckOut() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.ExecuteWithArgs(cmdCheckOut, s.SvnUrl)
	if err != nil {
		return err
	}
	return nil
}

func (s *svn) Update() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.ExecuteWithArgs(cmdUpdate)
	if err != nil {
		return err
	}
	return nil
}

func (s *svn) Status() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, err := s.ExecuteWithArgs(cmdStatus)
	if err != nil {
		return err
	}
	return nil
}

func (s *svn) AddAll() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.ExecuteWithArgs(cmdAddAll)
	if err != nil {
		return err
	}
	return nil
}

func (s *svn) Clean() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.ExecuteWithArgs(cmdClean)
	if err != nil {
		return err
	}
	return nil
}

func (s *svn) Commit(message string) error {
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

func (s *svn) Log(number int) (res []Logentry, err error) {
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

func (s *svn) Timer() {
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

func (s *svn) Listener(ch chan *Command) {
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

func NewSvnOperator(v *ProjectConfig, ctx context.Context) SvnOperator {
	var svn SvnOperator = &svn{
		ScriptPath:  fmt.Sprintf("%s%s", v.ScriptsPath, svnScriptName),
		ProjectName: v.ProjectName,
		Username:    v.Svn.Username,
		Password:    v.Svn.Password,
		WorkDir:     v.Svn.WorkDir,
		Url:         v.Svn.Url,
		Port:        v.Svn.Port,
		RemoteDir:   v.Svn.RemoteDir,
		SvnUrl:      fmt.Sprintf(svnUrl, v.Svn.Username, v.Svn.Url, v.Svn.Port, v.Svn.RemoteDir),
		ctx:         ctx,
	}
	if err := svn.CheckOut(); err != nil {
		log.Println(err)
	}
	if err := svn.Update(); err != nil {
		log.Println(err)
	}
	go svn.Timer()
	return svn
}

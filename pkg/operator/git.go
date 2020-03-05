package operator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type GitOperator interface {
	GetBranchFullName(name string) string
	GetBranchShortName(name string) string
	ExecuteWithArgs(args ...string) (res []byte, err error)
	FetchAll() error
	ShowAll(lock bool) error
	CheckOutBranch(name string) error
	Generate(name string) error
	Commit(name string) error
	Push(name string) error
	Update(name string) error
	Common(name string) error
	SetSvnTag(name, tag string) error
	SvnSync(name string) error
	ChangeTaskCount(incr int32)
	LoopChan()
	GetCurrentTask() string
	HandleCommand(c *Command) error
}

const (
	errGitExec                = "Git %s exec.Command err:%v"
	errGitBranchWasNotExisted = "the branch name `%s` of the git is not existed"
	errSvnTagWasNull          = "the svnTag of the branch name `%s` was null"
	errWriteToChannelTimeout  = "the command `%v` writes to channel time out"

	execOutputTemplate = "Git Command `%s` output:\n%s\n"
)

const (
	remoteBranchPrefix = "remotes/origin/"
)

const (
	gitScriptName = "git.sh"

	cmdGitCheckOut = "checkout"
	cmdGitFetchAll = "fetch"
	cmdGitShowAll  = "showAll"
	cmdGitGenerate = "generate"
	cmdGitCommit   = "commit"
	cmdGitPush     = "push"
	cmdGitUpdate   = "update"
	cmdSvnSync     = "svnSync"
)

const (
	gitInActive = iota
	gitActive
)

type git struct {
	mu             sync.RWMutex
	ScriptPath     string             `json:"script_path"`
	Path           string             `json:"path"`
	Name           string             `json:"name"`
	RemoteBranches map[string]*Branch `json:"remote_branches"`
	ListBranches   []string           `json:"list_branches"`
	TaskCount      int32              `json:"task_count"`
	TaskChan       chan *Command
	CurrentTask    *Command
	ctx            context.Context
}

func (g *git) GetBranchFullName(name string) string {
	return fmt.Sprintf("%s%s", remoteBranchPrefix, name)
}

func (g *git) GetBranchShortName(name string) string {
	return strings.Replace(name, remoteBranchPrefix, "", -1)
}

func (g *git) ExecuteWithArgs(args ...string) (res []byte, err error) {
	t := append([]string{g.ScriptPath, g.Path}, args...)
	out, err := exec.Command("sh", t...).Output()
	if err != nil {
		return out, errors.New(fmt.Sprintf(errGitExec, args[0], err))
	}
	log.Printf(execOutputTemplate, args[0], string(out))
	return out, nil
}

func (g *git) FetchAll() (err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	_, err = g.ExecuteWithArgs(cmdGitFetchAll)
	if err != nil {
		return err
	}
	return nil
}

func (g *git) ShowAll(lock bool) (err error) {
	if lock == true {
		g.mu.Lock()
		defer g.mu.Unlock()
	}
	out, err := g.ExecuteWithArgs(cmdGitShowAll)
	if err != nil {
		return err
	}
	arr := strings.Split(string(out), "\n")
	if len(arr) == 0 {
		return errors.New(fmt.Sprintf("ShowAll len == 0 name:%s out:%v\n", g.Name, out))
	}
	tmp := make(map[string]*Branch, 0)
	g.ListBranches = make([]string, 0)
	var activeBranch string
	for _, v := range arr {
		if v == "" {
			continue
		}
		v = strings.Replace(v, " ", "", -1)
		s := strings.Replace(v, "*", "", -1)
		activeMatched, err := regexp.Match(`\*`, []byte(v))
		if err != nil {
			return errors.New(fmt.Sprintf("ShowAll regexp.Match active name:%s err:%v\n", g.Name, err))
		}
		if activeMatched == true {
			activeBranch = s
		}
		if matched, err := regexp.Match(`remotes`, []byte(v)); err != nil {
			return errors.New(fmt.Sprintf("ShowAll regexp.Match local/remote name:%s err:%v\n", g.Name, err))
		} else {
			if matched == false {
				continue
			}
			s = g.GetBranchShortName(s)
			if t, ok := g.RemoteBranches[s]; ok {
				if s == activeBranch {
					t.Active = gitActive
				}
				tmp[s] = t
			} else {
				b := &Branch{
					Name:   s,
					Active: gitInActive,
					SvnTag: "",
				}
				if s == activeBranch {
					b.Active = gitActive
				}
				tmp[s] = b
			}
			g.ListBranches = append(g.ListBranches, s)
		}
	}
	g.RemoteBranches = tmp
	return nil
}

func (g *git) CheckOutBranch(name string) (err error) {
	_, ok := g.RemoteBranches[name]
	if !ok {
		return errors.New(fmt.Sprintf(errGitBranchWasNotExisted, name))
	}
	_, err = g.ExecuteWithArgs(cmdGitCheckOut, name, g.GetBranchFullName(name))
	if err != nil {
		return err
	}
	return nil
}

func (g *git) Generate(name string) (err error) {
	_, err = g.ExecuteWithArgs(cmdGitGenerate, name)
	if err != nil {
		return err
	}
	return nil
}

func (g *git) Commit(name string) (err error) {
	_, err = g.ExecuteWithArgs(cmdGitCommit, name)
	if err != nil {
		return err
	}
	return nil
}

func (g *git) Push(name string) (err error) {
	_, err = g.ExecuteWithArgs(cmdGitPush, name)
	if err != nil {
		return err
	}
	return nil
}

func (g *git) Common(name string) (err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if err = g.ShowAll(false); err != nil {
		return err
	}
	if err = g.CheckOutBranch(name); err != nil {
		return err
	}
	if err = g.Generate(name); err != nil {
		return err
	}
	if err = g.Commit(name); err != nil {
		return err
	}
	if err = g.Push(name); err != nil {
		return err
	}
	return nil
}

func (g *git) Update(name string) (err error) {
	_, err = g.ExecuteWithArgs(cmdGitUpdate, name)
	if err != nil {
		return err
	}
	return nil
}

func (g *git) SetSvnTag(name, tag string) (err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	t, ok := g.RemoteBranches[name]
	if !ok {
		return errors.New(fmt.Sprintf(errGitBranchWasNotExisted, name))
	}
	t.SvnTag = tag
	return nil
}

func (g *git) SvnSync(name string) (err error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	t, ok := g.RemoteBranches[name]
	if !ok {
		return errors.New(fmt.Sprintf(errGitBranchWasNotExisted, name))
	}
	if t.SvnTag == "" {
		return errors.New(fmt.Sprintf(errSvnTagWasNull, name))
	}
	_, err = g.ExecuteWithArgs(cmdSvnSync, name, t.SvnTag)
	if err != nil {
		return err
	}
	return nil
}

func (g *git) ChangeTaskCount(incr int32) {
	atomic.AddInt32(&g.TaskCount, incr)
}


func (g *git) LoopChan() {
	defer close(g.TaskChan)
	tick := time.NewTicker(time.Second * 10)
	defer tick.Stop()
	for {
		select {
		case <-g.ctx.Done():
			return
		case c := <-g.TaskChan:
			g.CurrentTask = c
			if err := g.Common(c.branchName); err != nil {
				log.Printf("LoopChan Common err:%v\n", err)
			}
			if err := g.Update(c.branchName); err != nil {
				log.Printf("LoopChan Update err:%v\n", err)
			}
			g.ChangeTaskCount(-1)
			g.CurrentTask = nil
		case <-tick.C:
			go func() {
				if err := g.FetchAll(); err != nil {
					log.Println("LoopChan FetchAll err:", err)
				}
				if err := g.ShowAll(true); err != nil {
					log.Println("LoopChan show err:", err)
				}
			}()
		}
	}
}

func (g *git) GetCurrentTask() string {
	if g.CurrentTask == nil {
		return "N/A"
	} else {
		return fmt.Sprintf("%s-%s-%s", g.CurrentTask.projectName, g.CurrentTask.branchName, g.CurrentTask.command)
	}
}

func (g *git) HandleCommand(c *Command) (err error) {
	g.ChangeTaskCount(1)
	tick := time.NewTicker(time.Second * 1)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			g.ChangeTaskCount(-1)
			return errors.New(fmt.Sprintf(errWriteToChannelTimeout, c))
		case g.TaskChan <- c:
			return nil
		}
	}
}

func (g *git) GetGitInfo() GitInfo {
	g.mu.RLock()
	defer g.mu.RUnlock()
	gi := GitInfo{
		Name:         g.Name,
		ListBranches: make([]Branch, 0),
		TaskCount:    g.TaskCount,
		CurrentTask:  g.GetCurrentTask(),
	}
	for _, v := range g.ListBranches {
		if t, ok := g.RemoteBranches[v]; ok {
			gi.ListBranches = append(gi.ListBranches, *t)
		}
	}
	return gi
}

func NewGitOperator(v *ProjectConfig, ctx context.Context) GitOperator {
	var git GitOperator = &git{
		ScriptPath:     fmt.Sprintf("%s%s", v.ScriptsPath, gitScriptName),
		Path:           v.Git.WorkDir,
		Name:           v.ProjectName,
		RemoteBranches: make(map[string]*Branch, 0),
		ListBranches:   make([]string, 0),
		TaskCount:      0,
		TaskChan:       make(chan *Command, 1024),
		ctx:            ctx,
	}
	go git.LoopChan()
	return git
}

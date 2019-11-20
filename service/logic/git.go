package logic

import (
	"context"
	"encoding/json"
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

type GitHub struct {
	mu         sync.RWMutex
	ScriptPath string
	Gits       map[string]*Git
}

type Git struct {
	mu             sync.RWMutex
	ScriptPath     string         `json:"script_path"`
	Path           string         `json:"path"`
	Name           string         `json:"name"`
	ActiveBranch   string         `json:"active_branch"`
	LocalBranches  map[string]int `json:"local_branches"`
	RemoteBranches map[string]int `json:"remote_branches"`
	ListBranches   []string       `json:"list_branches"`
	TaskCount      int32          `json:"task_count"`
	TaskChan       chan *Command
	ctx            context.Context
}

func (g *Git) ShowAll(lock bool) (err error) {
	if lock == true {
		g.mu.Lock()
		defer g.mu.Unlock()
	}
	if out, err := exec.Command("sh", g.ScriptPath, g.Path, "all").Output(); err != nil {
		return errors.New(fmt.Sprintf("ShowAll exec.Command name:%s err:%v\n", g.Name, err))
	} else {
		arr := strings.Split(string(out), "\n")
		if len(arr) > 0 {
			g.ActiveBranch = ""
			g.LocalBranches = make(map[string]int, 0)
			g.RemoteBranches = make(map[string]int, 0)
			g.ListBranches = make([]string, 0)
		} else {
			return errors.New(fmt.Sprintf("ShowAll exec.Command name:%s len 0 out:%v\n", g.Name, out))
		}
		for _, v := range arr {
			if v == "" {
				continue
			}
			v = strings.Replace(v, " ", "", -1)
			s := strings.Replace(v, "*", "", -1)
			if matched, err := regexp.Match(`\*`, []byte(v)); err != nil {
				return errors.New(fmt.Sprintf("ShowAll regexp.Match active name:%s err:%v\n", g.Name, err))
			} else {
				if matched == true {
					g.ActiveBranch = s
				}
			}
			if matched, err := regexp.Match(`remotes`, []byte(v)); err != nil {
				return errors.New(fmt.Sprintf("ShowAll regexp.Match local/remote name:%s err:%v\n", g.Name, err))
			} else {
				if matched == true {
					g.RemoteBranches[s] = 1
					g.ListBranches = append(g.ListBranches, s)
				} else {
					g.LocalBranches[s] = 1
				}
			}
		}
	}
	return nil
}

func (g *Git) ShowActive() {
	log.Printf("%s active:%s\n", g.Name, g.ActiveBranch)
}

func (g *Git) CheckOutBranch(name string) (err error) {
	fullName := fmt.Sprintf("remotes/origin/%s", name)
	if _, ok := g.RemoteBranches[fullName]; ok {
		if out, err := exec.Command("sh", g.ScriptPath, g.Path, "checkout", name, fullName).Output(); err != nil {
			return errors.New(fmt.Sprintf("CheckOutBranch exec.Command name:%s err:%v\n", g.Name, err))
		} else {
			log.Println("out:", string(out))
		}
	} else {
		return errors.New(fmt.Sprintf("CheckOutBranch exec.Command name:%s to:`%s` wasn't exist\n", g.Name, name))
	}
	return nil
}

func (g *Git) Generator(name string) (err error) {
	if out, err := exec.Command("sh", g.ScriptPath, g.Path, "generator", name).Output(); err != nil {
		return errors.New(fmt.Sprintf("Generator exec.Command name:%s err:%v\n", g.Name, err))
	} else {
		log.Println("out:", string(out))
	}
	return nil
}

func (g *Git) Commit(name string) (err error) {
	if out, err := exec.Command("sh", g.ScriptPath, g.Path, "commit", name).Output(); err != nil {
		return errors.New(fmt.Sprintf("Commit exec.Command name:%s err:%v\n", g.Name, err))
	} else {
		log.Println("out:", string(out))
	}
	return nil
}

func (g *Git) Push(name string) (err error) {
	if out, err := exec.Command("sh", g.ScriptPath, g.Path, "push", name).Output(); err != nil {
		return errors.New(fmt.Sprintf("Push exec.Command name:%s err:%v\n", g.Name, err))
	} else {
		log.Println("out:", string(out))
	}
	return nil
}

func (g *Git) ChangeTaskCount(incr int32) {
	atomic.AddInt32(&g.TaskCount, incr)
}

func (g *Git) Common(name string) (err error) {
	if err = g.ShowAll(true); err != nil {
		return err
	}
	if err = g.CheckOutBranch(name); err != nil {
		return err
	}
	if err = g.Generator(name); err != nil {
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

func (g *Git) LoopChan() {
	defer close(g.TaskChan)
	for {
		select {
		case <-g.ctx.Done():
			return
		case c := <-g.TaskChan:
			if err := g.Common(c.branchName); err != nil {
				log.Printf("LoopChan Common err:(%v)\n", err)
			}
			g.ChangeTaskCount(-1)
		}
	}
}

func (s *Service) NewGitHub() *GitHub {
	g := &GitHub{
		ScriptPath: fmt.Sprintf("%s%s", s.C.ScriptsPath, "git.sh"),
		Gits:       make(map[string]*Git, 0),
	}
	for _, v := range s.C.Projects {
		git := &Git{
			ScriptPath:     g.ScriptPath,
			Path:           v[1],
			Name:           v[0],
			ActiveBranch:   "",
			LocalBranches:  make(map[string]int, 0),
			RemoteBranches: make(map[string]int, 0),
			ListBranches:   make([]string, 0),
			TaskCount:      0,
			TaskChan:       make(chan *Command, 1024),
			ctx:            s.ctx,
		}
		g.Gits[git.Name] = git
		go git.LoopChan()
	}
	return g
}

func (gh *GitHub) handleAll() (res string, err error) {
	t := make([]GitResponse, 0)
	for _, v := range gh.Gits {
		if err = v.ShowAll(true); err != nil {
			return "", nil
		}
		tmp := GitResponse{
			Name:         v.Name,
			ActiveBranch: fmt.Sprintf("remotes/origin/%s", v.ActiveBranch),
			ListBranches: v.ListBranches,
			TaskCount:    v.TaskCount,
		}
		t = append(t, tmp)
	}
	if ret, err := json.Marshal(t); err != nil {
		return "", err
	} else {
		return string(ret), nil
	}
}

func (gh *GitHub) handleCommand(c *Command) (err error) {
	if t, ok := gh.Gits[c.projectName]; ok {
		t.ChangeTaskCount(1)
		tick := time.NewTicker(time.Second * 1)
		defer tick.Stop()
		for {
			select {
			case <-tick.C:
				return
			case t.TaskChan <- c:
				return
			}
		}
	} else {
		log.Printf("no project:%v\n", c)
	}
	return nil
}

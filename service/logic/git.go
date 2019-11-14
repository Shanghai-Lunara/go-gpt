package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type GitHub struct {
	ScriptPath string
	Gits       map[string]*Git
}

type Git struct {
	ScriptPath     string         `json:"script_path"`
	Path           string         `json:"path"`
	Name           string         `json:"name"`
	ActiveBranch   string         `json:"active_branch"`
	LocalBranches  map[string]int `json:"local_branches"`
	RemoteBranches map[string]int `json:"remote_branches"`
	ListBranches   []string       `json:"list_branches"`
}

func (g *Git) ShowAll() (err error) {
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
	if err = g.ShowAll(); err != nil {
		return errors.New(fmt.Sprintf("CheckOutBranch err:(%v)\n", err))
	}
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
		}
		g.Gits[git.Name] = git
	}
	return g
}

func (g *GitHub) handleAll() (res string, err error) {
	t := make([]Git, 0)
	for _, v := range g.Gits {
		if err = v.ShowAll(); err != nil {
			return "", nil
		}
		tmp := *v
		tmp.ActiveBranch = fmt.Sprintf("remotes/origin/%s", tmp.ActiveBranch)
		tmp.LocalBranches = make(map[string]int, 0)
		tmp.RemoteBranches = make(map[string]int, 0)
		t = append(t, tmp)
	}
	if ret, err := json.Marshal(t); err != nil {
		return "", err
	} else {
		return string(ret), nil
	}
}

package logic

import (
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
	ScriptPath   string
	Path         string
	Name         string
	ActiveBranch string
	Branches     []string
}

func (g *Git) ShowAll() (err error) {
	if out, err := exec.Command("sh", g.ScriptPath, g.Path, "all").Output(); err != nil {
		return errors.New(fmt.Sprintf("ShowAll exec.Command name:%s err:%v\n", g.Name, err))
	} else {
		arr := strings.Split(string(out), "\n")
		if len(arr) > 0 {
			g.ActiveBranch = ""
			g.Branches = make([]string, 0)
		} else {
			return errors.New(fmt.Sprintf("ShowAll exec.Command name:%s len 0 out:%v\n", g.Name, out))
		}
		for _, v := range arr {
			if v == "" {
				continue
			}
			v = strings.Replace(v, " ", "", -1)
			if matched, err := regexp.Match(`\*`, []byte(v)); err != nil {
				return errors.New(fmt.Sprintf("ShowAll regexp.Match name:%s err:%v\n", g.Name, err))
			} else {
				if matched == true {
					g.ActiveBranch = strings.Replace(v, "*", "", -1)
				}
			}
			g.Branches = append(g.Branches, strings.Replace(v, "*", "", -1))
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
	exist := false
	for _, v := range g.Branches {
		if v == name {
			exist = true
			if out, err := exec.Command("sh", g.ScriptPath, g.Path, "change", name).Output(); err != nil {
				return errors.New(fmt.Sprintf("CheckOutBranch exec.Command name:%s err:%v\n", g.Name, err))
			} else {
				log.Println("out:", string(out))
			}
		}
	}
	if exist == false {
		return errors.New(fmt.Sprintf("CheckOutBranch exec.Command name:%s to:`%s` wasn't exist\n", g.Name, name))
	}
	return nil
}

func (g *Git) Generator(name string) (err error) {
	if err = g.CheckOutBranch(name); err != nil {
		return errors.New(fmt.Sprintf("Generator CheckOutBranch exec.Command name:%s to:`%s` wasn't exist\n", g.Name, name))
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
			ScriptPath:   g.ScriptPath,
			Path:         v[1],
			Name:         v[0],
			ActiveBranch: "",
			Branches:     make([]string, 0),
		}
		g.Gits[git.Name] = git
	}
	for _, v := range g.Gits {
		if err := v.CheckOutBranch("develop"); err != nil {
			log.Println(err)
		}
	}
	return g
}

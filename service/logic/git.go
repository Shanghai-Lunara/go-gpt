package logic

import (
	"fmt"
	"go-gpt/conf"
	"log"
	"os/exec"
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
}

func (g *Git) ShowAll() {
	if out, err := exec.Command("sh", g.ScriptPath, g.Path, "all").Output(); err != nil {
		log.Printf("ShowAll name:%s err:%v\n", g.Name, err)
	} else {
		fmt.Println(string(out))
	}
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
		}
		g.Gits[git.Name] = git
	}
	return g
}

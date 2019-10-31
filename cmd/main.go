package main

import (
	"flag"
	"go-gpt/conf"
	"log"
	"os/exec"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Fatal(err)
	}
	log.Println(conf.GetConfig())
	if out, err := exec.Command("sh", "../scripts/git.sh").Output(); err != nil {
		log.Println(err)
	} else {
		log.Println(string(out))
	}
}

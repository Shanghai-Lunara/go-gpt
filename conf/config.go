package conf

import (
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ProjectsConfig struct {
	Projects [][]string `yaml:"Projects,flow"`
}

type Config struct {
	LogFile        string `yaml:"LogFile"`
	ScriptsPath    string `yaml:"ScriptsPath"`
	ProjectsConfig `yaml:",inline"`
	ConfigPath     string
}

var (
	conf       *Config
	configPath string
)

func init() {
	flag.StringVar(&configPath, "configPath", "gpt.yml", "configuration file path")
}

func Init() (err error) {
	var (
		data []byte
	)
	if data, err = ioutil.ReadFile(configPath); err != nil {
		return errors.New(fmt.Sprintf("ioutil.ReadFile err:%v", err))
	}
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return errors.New(fmt.Sprintf("yaml.Unmarshal err:%v", err))
	}
	conf.ConfigPath = configPath
	//log.Println("a:", conf)
	return nil
}

func GetConfig() *Config {
	return conf
}

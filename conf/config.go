package conf

import (
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type HttpConfig struct {
	IP            string `yaml:"ip"`
	Port          int    `yaml:"port"`
	TemplatesPath string `yaml:"templates_path"`
}

type Config struct {
	Http       HttpConfig `yaml:"HttpService"`
	LogFile    string     `yaml:"LogFile"`
	Projects   []Project  `yaml:"Projects"`
	ConfigPath string
}

type Project struct {
	ProjectName string    `yaml:"project_name"`
	ScriptsPath string    `yaml:"scripts_path"`
	Git         GitConfig `yaml:"git"`
	Svn         SvnConfig `yaml:"svn"`
}

type GitConfig struct {
	WorkDir string `yaml:"work_dir"`
}

type SvnConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	WorkDir   string `yaml:"work_dir"`
	Url       string `yaml:"url"`
	Port      int    `yaml:"port"`
	RemoteDir string `yaml:"remote_dir"`
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
	return nil
}

func GetConfig() *Config {
	return conf
}

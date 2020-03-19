package conf

import (
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	"github.com/Shanghai-Lunara/go-gpt/pkg/operator"
)

type HttpConfig struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

type Config struct {
	Http       HttpConfig               `yaml:"HttpService"`
	LogFile    string                   `yaml:"LogFile"`
	Projects   []operator.ProjectConfig `yaml:"Projects"`
	ConfigPath string
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

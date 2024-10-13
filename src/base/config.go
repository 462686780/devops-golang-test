package base

import (
	"fmt"
	"io/ioutil"
	"os"

	"time"

	"gopkg.in/yaml.v3"
)

var (
	Config = &GlobalSetting{
		RunMode: "production",
		Logger: LoggerSetting{
			Level: 5,
			Name:  "statefulset.log",
			Dir:   "D:\\statefulset\\logs",
		},
		Listen: Listen{
			Addr:         "0.0.0.0",
			Port:         8899,
			ReadTimeout:  35,
			WriteTimeout: 35,
			IdleTimeout:  35,
		},
	}
)

type GlobalSetting struct {
	Server  ServerSetting `yaml:"server"`
	Logger  LoggerSetting `yaml:"log"`
	Listen  Listen        `yaml:"listen"`
	RunMode string        `yaml:"run_mode"`
}

type ServerSetting struct {
	Addr              string        `yaml:"addr"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
	ExitTimeout       time.Duration `yaml:"exit_timeout"`
	MaxHeaderBytes    int           `yaml:"max_header_bytes"`
}

type LoggerSetting struct {
	Level int    `yaml:"level"`
	Name  string `yaml:"name"`
	Dir   string `yaml:"dir"`
}

type Listen struct {
	Addr         string        `yaml:"addr"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

func LoadConfig(cfg string) (*GlobalSetting, error) {
	if len(cfg) == 0 {
		return Config, nil
	}
	if stat, err := os.Stat(cfg); err != nil || stat.IsDir() {
		return nil, fmt.Errorf("%s is not file", cfg)
	}
	bs, err := ioutil.ReadFile(cfg)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(bs, Config); err != nil {
		return nil, err
	}
	return Config, nil
}

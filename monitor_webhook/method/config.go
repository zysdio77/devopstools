package method

import (
	"os"

	"gopkg.in/yaml.v3"
)

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type RouteConfig struct {
	UseDing  *bool    `yaml:"use_ding"`
	DingTalk string   `yaml:"dingtalk"`
	UseEmail *bool    `yaml:"use_email"`
	Email    []string `yaml:"email"`
}

type Config struct {
	SMTP   SMTPConfig             `yaml:"smtp"`
	Routes map[string]RouteConfig `yaml:"routes"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

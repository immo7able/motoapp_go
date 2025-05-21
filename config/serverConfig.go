package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Config struct {
	ServerAddress string `yaml:"server_address"`
	DatabaseUrl   string `yaml:"db_url"`
	JWT           struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
}

func NewConfig() *Config {
	content, err := os.ReadFile("config/application.yaml")
	if err != nil {
		panic("Failed to read config file: " + err.Error())
	}

	expanded := os.ExpandEnv(string(content))

	var cfg Config
	if err := yaml.NewDecoder(strings.NewReader(expanded)).Decode(&cfg); err != nil {
		panic("Failed to decode config: " + err.Error())
	}

	return &cfg
}

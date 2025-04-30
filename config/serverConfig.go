package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	ServerAddress string `yaml:"server_address"`
	DatabaseUrl   string `yaml:"db_url"`
	JWT           struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
}

func NewConfig() *Config {
	cfg := &Config{}

	bytes, err := ioutil.ReadFile("config/application.yaml")
	if err != nil {
		log.Fatalf("could not read config file: %v", err)
	}

	err = yaml.Unmarshal(bytes, cfg)
	if err != nil {
		log.Fatalf("could not unmarshal config: %v", err)
	}

	return cfg
}

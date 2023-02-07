package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	GetEndpoint    string `yaml:"getEndpoint"`
	PostEndpoint   string `yaml:"postEndpoint"`
	OpaqueEndpoint string `yaml:"opaqueEndpoint"`
	Sleep          int    `yaml:"sleep"`
	UserAgent      string `yaml:"userAgent"`
}

func NewConfig(path string) (*Config, error) {
	var conf Config
	file, err := os.ReadFile(path)
	if err != nil {
		return &conf, err
	}
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		log.Fatal(err)
	}
	return &conf, nil
}

func (config *Config) GetHost() string {
	return config.Host
}

func (config *Config) GetPort() int {
	return config.Port
}

func (config *Config) GetGetEndpoint() string {
	return config.GetEndpoint
}

func (config *Config) GetPostEndpoint() string {
	return config.PostEndpoint
}

func (config *Config) GetOpaqueEndpoint() string {
	return config.OpaqueEndpoint
}

func (config *Config) GetSleep() int {
	return config.Sleep
}

func (config *Config) GetUserAgent() string {
	return config.UserAgent
}

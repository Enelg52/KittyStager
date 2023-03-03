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
	Jitter         int    `yaml:"jitter"`
	UserAgent      string `yaml:"userAgent"`
	Protocol       string `yaml:"protocol"`
	Key            string `yaml:"key"`
	Cert           string `yaml:"cert"`
	LocalUpload    string `yaml:"localUpload"`
	WebUpload      string `yaml:"webUpload"`
	malPath        string `yaml:"malPath"`
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

func (config *Config) GetJitter() int {
	return config.Jitter
}

func (config *Config) SetSleep(i int) {
	config.Sleep = i
}

func (config *Config) GetUserAgent() string {
	return config.UserAgent
}

func (config *Config) GetProtocol() string {
	return config.Protocol
}

func (config *Config) GetKey() string {
	return config.Key
}
func (config *Config) GetCert() string {
	return config.Cert
}

func (config *Config) GetLocalUpload() string {
	return config.LocalUpload
}

func (config *Config) GetWebUpload() string {
	return config.WebUpload
}

func (config *Config) GetMalPath() string {
	return config.malPath
}

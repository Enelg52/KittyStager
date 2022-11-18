package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

// General is the struct that contains the config
type General struct {
	Conf Http `yaml:"Http"`
}

type Http struct {
	Host      string   `yaml:"host"`
	Port      int      `yaml:"port"`
	EndPoint  string   `yaml:"endpoint"`
	Sleep     int      `yaml:"sleep"`
	UserAgent string   `yaml:"userAgent"`
	MalPath   []string `yaml:"malPath,flow"`
	Reg1      string   `yaml:"reg1"`
	Reg2      string   `yaml:"reg2"`
	Auth1     string   `yaml:"auth1"`
	Auth2     string   `yaml:"auth2"`
}

func NewConfig(path string) (*General, error) {
	var conf General
	filename, err := filepath.Abs(path)
	if err != nil {
		return &conf, err
	}
	file, err := os.ReadFile(filename)
	if err != nil {
		return &conf, err
	}
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		log.Fatal(err)
	}
	return &conf, error(nil)
}

func (c *General) GetHost() string {
	return c.Conf.Host
}

func (c *General) GetPort() int {
	return c.Conf.Port
}

func (c *General) GetEndpoint() string {
	return c.Conf.EndPoint
}

func (c *General) GetMalPath() []string {
	return c.Conf.MalPath
}

func (c *General) GetMalPathWithId(i int) string {
	return c.Conf.MalPath[i]
}

func (c *General) GetUserAgent() string {
	return c.Conf.UserAgent
}

func (c *General) GetSleep() int {
	return c.Conf.Sleep
}

func (c *General) GetReg1() string {
	return c.Conf.Reg1
}

func (c *General) GetReg2() string {
	return c.Conf.Reg2
}

func (c *General) GetAuth1() string {
	return c.Conf.Auth1
}

func (c *General) GetAuth2() string {
	return c.Conf.Auth2
}

func (c *General) CheckConfig() error {
	if c.Conf.Host == "" || c.Conf.Port == 0 || c.Conf.EndPoint == "" || c.Conf.UserAgent == "" || c.Conf.MalPath == nil {
		return errors.New("please check your config file")
	}
	return nil
}

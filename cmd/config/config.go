package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"path/filepath"
)

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
}

func GetConfig() (General, error) {
	var c General
	filename, _ := filepath.Abs("../config/conf.yml")
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		log.Fatal(err)
	}
	return c, error(nil)
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

func (c *General) GetUserAgent() string {
	return c.Conf.UserAgent
}

func (c *General) GetSleep() int {
	return c.Conf.Sleep
}

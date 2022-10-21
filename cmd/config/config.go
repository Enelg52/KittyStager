package config

import (
	"errors"
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
	var conf General
	filename, _ := filepath.Abs("C:\\Users\\yanng\\go\\Project_go\\GoStager\\cmd\\config\\conf.yml")
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return conf, err
	}
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		log.Fatal(err)
	}
	return conf, error(nil)
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

func (c *General) CheckConfig() error {
	if c.Conf.Host == "" || c.Conf.Port == 0 || c.Conf.EndPoint == "" || c.Conf.UserAgent == "" || c.Conf.MalPath == nil {
		return errors.New("please check your config file")
	}
	return nil
}

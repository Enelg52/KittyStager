package api

import "encoding/json"

type Config struct {
	Host         string `json:"Host"`
	GetEndpoint  string `json:"GetEndpoint"`
	PostEndpoint string `json:"PostEndpoint"`
	RegEndpoint  string `json:"RegEndpoint"`
	UserAgent    string `json:"UserAgent"`
	Sleep        int    `json:"sleep"`
	Jitter       int    `json:"Jitter"`
	Name         string `json:"Name"`
}

func NewConfig(host, getEndpoint, postEndpoint, regEndpoint, userAgent, name string, sleep, jitter int) *Config {
	return &Config{host,
		getEndpoint,
		postEndpoint,
		regEndpoint,
		userAgent,
		sleep,
		jitter,
		name,
	}
}

func (config *Config) MarshallConfig() ([]byte, error) {
	t, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (config *Config) UnmarshallConfig(j []byte) error {
	err := json.Unmarshal(j, &config)
	if err != nil {
		return err
	}
	return nil
}

func (config *Config) SetName(name string) {
	config.Name = name
}

func (config *Config) SetSleep(sleep int) {
	config.Sleep = sleep
}

func (config *Config) getName() string {
	return config.Name
}

func (config *Config) getGetEndpoint() string {
	return config.GetEndpoint
}

func (config *Config) getPostEndpoint() string {
	return config.PostEndpoint
}

func (config *Config) getRegEndpoint() string {
	return config.RegEndpoint
}

func (config *Config) getUserAgent() string {
	return config.UserAgent
}

func (config *Config) getSleep() int {
	return config.Sleep
}

func (config *Config) getJitter() int {
	return config.Jitter
}

func (config *Config) getHost() string {
	return config.Host
}

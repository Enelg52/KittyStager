package test

import (
	"KittyStager/internal/api"
	"KittyStager/internal/config"
	"KittyStager/internal/kitten"
	"encoding/json"
	"fmt"
	"github.com/go-playground/assert/v2"
	"io"
	"net/http"
	"testing"
	"time"
)

var (
	c           http.Client
	backend     string
	kittens     map[string]*kitten.Kitten
	kittensList string
	conf        *config.Config
)

func apiBeforeAll() {
	host = "127.0.0.1"
	port = 8080
	getEndpoint = "getLegit/test"
	postEndpoint = "postLegit"
	opaqueEndpoint = "reg"
	sleepTime = 5
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/102.0"
	c = http.Client{Timeout: time.Duration(3) * time.Second}
	backend = "127.0.0.1:1337"
	kittensList = "kittensList"
	kittens = make(map[string]*kitten.Kitten)
	var err error
	conf, err = config.NewConfig("../config.yaml")
	if err != nil {
		return
	}
	go func() {
		err := api.Api(conf)
		if err != nil {
			return
		}
	}()
}

func TestApiBackConf(t *testing.T) {
	t.Parallel()
	// given
	// when
	target := fmt.Sprintf("http://%s/%s", backend, kittensList)
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	// then
	assert.Equal(t, resp.StatusCode, 200)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	err = json.Unmarshal(b, &conf)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	assert.Equal(t, conf.Host, host)
}

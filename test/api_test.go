package test

import (
	"KittyStager/internal/api"
	"KittyStager/internal/config"
	"KittyStager/internal/kitten"
	"KittyStager/malware"
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
	getEndpoint = "getLegit"
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

func TestApiFront(t *testing.T) {
	t.Parallel()
	// given

	// when
	target := fmt.Sprintf("http://%s:%d/%s", host, port, getEndpoint)
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

	// when
	target = fmt.Sprintf("http://%s:%d/%s", host, port, postEndpoint)
	req, err = http.NewRequest("POST", target, nil)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	resp, err = c.Do(req)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	// then
	assert.Equal(t, resp.StatusCode, 200)

	// when
	target = fmt.Sprintf("http://%s:%d/%s", host, port, opaqueEndpoint)
	req, err = http.NewRequest("POST", target, nil)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	resp, err = c.Do(req)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	// then
	assert.Equal(t, resp.StatusCode, 200)
}

func TestApiBackKittensList(t *testing.T) {
	t.Parallel()
	// given
	conf := malware.NewConfig("http://127.0.0.1:8080",
		"getLegit",
		"postLegit",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/102.0",
		"reg",
		"test",
		0,
	)
	// when
	err := malware.DoPwreg(username, password, *conf)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	_, err = malware.DoAuth(username, password, *conf)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
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
	err = json.Unmarshal(b, &kittens)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	//convert to json
	jTest, err := json.MarshalIndent(&kittens, "", " ")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	//convert to json
	jApi, err := json.MarshalIndent(&api.Kittens, "", " ")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	assert.Equal(t, jTest, jApi)
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

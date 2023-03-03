package test

import (
	"KittyStager/internal/config"
	"github.com/go-playground/assert/v2"
	"testing"
)

var (
	host           string
	port           int
	getEndpoint    string
	postEndpoint   string
	opaqueEndpoint string
	sleepTime      int
	userAgent      string
)

func configBeforeAll() {
	host = "127.0.0.1"
	port = 8080
	getEndpoint = "getLegit"
	postEndpoint = "postLegit"
	opaqueEndpoint = "reg"
	sleepTime = 2
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/102.0"
}

func TestConfig(t *testing.T) {
	t.Parallel()
	conf, err := config.NewConfig("../config.yaml")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	assert.Equal(t, conf.GetHost(), host)
	assert.Equal(t, conf.GetPort(), port)
	assert.Equal(t, conf.GetGetEndpoint(), getEndpoint)
	assert.Equal(t, conf.GetPostEndpoint(), postEndpoint)
	assert.Equal(t, conf.GetSleep(), sleepTime)
	assert.Equal(t, conf.GetUserAgent(), userAgent)
	assert.Equal(t, conf.GetOpaqueEndpoint(), opaqueEndpoint)
}

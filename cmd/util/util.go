package util

import (
	"KittyStager/cmd/config"
	"encoding/json"
	"fmt"
	color "github.com/logrusorgru/aurora"
	"os"
)

// GenerateConfig generate the config file for all the kitten
func GenerateConfig(conf config.General) error {
	malConf := NewMalConf()
	malConf.SetHost(fmt.Sprintf("http://%s:%d", conf.GetHost(), conf.GetPort()))
	malConf.SetEndPoint(conf.GetEndpoint())
	malConf.SetUserA(conf.GetUserAgent())
	malConf.SetSleep(conf.GetSleep())
	malConf.SetReg1(conf.GetReg1())
	malConf.SetReg2(conf.GetReg2())
	malConf.SetAuth1(conf.GetAuth1())
	malConf.SetAuth2(conf.GetAuth2())
	data, err := MalConfMarshalJSON(malConf)
	if err != nil {
		return err
	}
	//data := fmt.Sprintf("http://%s:%d,%s,%s,%d,%s,%s,%s,%s", conf.GetHost(), conf.GetPort(), conf.GetEndpoint(), conf.GetUserAgent(), conf.GetSleep(), conf.GetReg1(), conf.GetReg2(), conf.GetAuth1(), conf.GetAuth2())
	for x := range conf.GetMalPath() {
		err := os.WriteFile(conf.GetMalPathWithId(x)+"conf.txt", data, 0644)
		if err != nil {
			return err
		}
		fmt.Printf("%s %s\n", color.Green("[+] Generated conf file for"), color.Yellow(conf.GetMalPathWithId(x)))
	}
	return nil
}

// ErrPrint print the error
func ErrPrint(err error) {
	if err != nil {
		fmt.Printf("\n%s %s\n", color.Red("[-]"), color.Red(err.Error()))
	}
}

// PrintInit print the recon info when the kitten calls back
func PrintInit(recon []byte) error {
	initChecks, err := InitUnmarshalJSON(recon)
	if err != nil {
		return err
	}
	fmt.Printf("%s %s\n", color.Green("[+] Hostname:"), color.Yellow(initChecks.GetHostname()))
	fmt.Printf("%s %s\n", color.Green("[+] Domain:"), color.Yellow(initChecks.GetDomain()))
	fmt.Printf("%s %s\n", color.Green("[+] Username:"), color.Yellow(initChecks.GetUsername()))
	fmt.Printf("%s %s\n", color.Green("[+] IP:"), color.Yellow(initChecks.GetIp()))
	fmt.Print(color.Green("[+] To get more, use the recon command\n"))
	return nil
}

// PrintRecon print the recon info when the command recon is called
func PrintRecon(i InitialChecks) {
	fmt.Printf("\n%s %s\n", color.Green("[+] Kitten name:"), color.Yellow(i.GetKittenName()))
	fmt.Printf("%s %s\n", color.Green("[+] Hostname:"), color.Yellow(i.GetHostname()))
	fmt.Printf("%s %s\n", color.Green("[+] Username:"), color.Yellow(i.GetUsername()))
	fmt.Printf("%s %s\n", color.Green("[+] IP:"), color.Yellow(i.GetIp()))
	fmt.Printf("%s %s\n", color.Green("[+] Domain:"), color.Yellow(i.GetDomain()))
	fmt.Printf("%s %s\n", color.Green("[+] Process:"), color.Yellow(i.GetPName()))
	fmt.Printf("%s %s\n", color.Green("[+] Process path:"), color.Yellow(i.GetPath()))
	fmt.Printf("%s %s\n", color.Green("[+] Process pid:"), color.Yellow(fmt.Sprintf("%d", i.GetPid())))
	fmt.Print(color.Green("[+] Installed software : \n"))
	f := i.GetDir()
	s, _ := json.MarshalIndent(f, "", "\t")
	fmt.Printf("%s\n", color.Yellow(string(s)))
}

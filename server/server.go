package main

import (
	"KittyStager/internal/config"
	"KittyStager/server/api"
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

package main

import (
	"KittyStager/internal/api"
	"KittyStager/internal/config"
)

func main() {
	conf, _ := config.NewConfig("config.yaml")
	err := api.Api(conf)
	if err != nil {
		return
	}
}

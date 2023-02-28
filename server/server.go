package main

import (
	"KittyStager/internal/config"
	"KittyStager/kitten/malware"
	"KittyStager/server/api"
	"fmt"
	"os"
)

func main() {
	conf, _ := config.NewConfig("config.yaml")
	host := fmt.Sprintf("%s://%s:%d", conf.GetProtocol(), conf.GetHost(), conf.GetPort())
	malConf := malware.NewConfig(host,
		conf.GetEndpoint,
		conf.GetPostEndpoint(),
		conf.GetOpaqueEndpoint(),
		conf.GetUserAgent(),
		"",
		conf.GetSleep(),
		conf.GetJitter(),
	)
	c, err := malConf.MarshallConfig()
	if err != nil {
		return
	}
	err = os.WriteFile("kitten/basicKitten/conf.txt", c, 0644)
	if err != nil {
		return
	}
	fmt.Println("[*] Running...")
	err = api.Api(conf)
	if err != nil {
		return
	}
}

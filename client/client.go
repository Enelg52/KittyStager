package main

import (
	"KittyStager/client/cli"
	"fmt"
	"github.com/inancgumus/screen"
	color "github.com/logrusorgru/aurora"
	"log"
)

func main() {
	screen.Clear()
	screen.MoveTopLeft()
	fmt.Println(color.BrightCyan("                     _\n                    / )\n                   ( (\n     A.-.A  .-\"\"-.  ) )\n    / , , \\/      \\/ /\n   =\\  t  ;=    /   /\n     `--,'  .\"\"|   /\n         || |  \\\\ \\\n        ((,_|  ((,_\\\n"))
	fmt.Println(color.BrightCyan("KittyStager - A simple stager written in Go\n"))
	err := cli.Cli()
	if err != nil {
		log.Fatal(err)
	}
}

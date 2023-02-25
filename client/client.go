package main

import (
	"KittyStager/client/cli"
	"github.com/inancgumus/screen"
	"log"
)

func main() {
	screen.Clear()
	screen.MoveTopLeft()
	err := cli.Cli()
	if err != nil {
		log.Fatal(err)
	}
}

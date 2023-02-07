package main

import (
	"KittyStager/internal/client"
	"github.com/inancgumus/screen"
	"log"
)

func main() {
	screen.Clear()
	screen.MoveTopLeft()
	err := client.Cli()
	if err != nil {
		log.Fatal(err)
	}
}

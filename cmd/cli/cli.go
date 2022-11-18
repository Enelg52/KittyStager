package cli

import (
	"KittyStager/cmd/config"
	"KittyStager/cmd/http"
	"KittyStager/cmd/util"
	"fmt"
	"github.com/c-bata/go-prompt"
	"os"
	"strings"
)

// completerCli is the completer for the main menu
func completerCli(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "exit", Description: "Exit the program"},
		{Text: "config", Description: "Show config"},
		{Text: "target", Description: "Show targets"},
		{Text: "interact", Description: "Interact with a target"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// completerInteract is the completer for the interact menu
func completerInteract(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "exit", Description: "Exit the program"},
		{Text: "back", Description: "Go back to the main menu"},
		{Text: "target", Description: "Show targets"},
		{Text: "interact", Description: "Interact with a target"},
		{Text: "payload", Description: "Host a payload"},
		{Text: "sleep", Description: "Set sleep time"},
		{Text: "recon", Description: "Show recon information"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// Cli is the main menu
func Cli(conf config.General) {
	for {
		t := prompt.Input("KittyStager ❯ ", completerCli,
			prompt.OptionTitle("KittyStager 🐈 "),
			prompt.OptionPrefixTextColor(prompt.Blue),
			prompt.OptionPreviewSuggestionTextColor(prompt.Green),
			prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
			prompt.OptionSelectedSuggestionTextColor(prompt.Blue),
			prompt.OptionDescriptionBGColor(prompt.Blue),
			prompt.OptionSuggestionBGColor(prompt.DarkGray),
		)
		switch t {
		case "exit":
			fmt.Println("Bye bye!")
			return
		case "config":
			printConfig(conf)
		case "target":
			printTarget()
		case "interact":
			interact()
		default:
			fmt.Println("Unknown command")
		}
	}
}

// Interact is the interact menu
func Interact(kittenName string) error {
	in := fmt.Sprintf("KittyStager - %s❯ ", kittenName)
	for {
		t := prompt.Input(in, completerInteract,
			prompt.OptionPrefixTextColor(prompt.Blue),
			prompt.OptionPreviewSuggestionTextColor(prompt.Green),
			prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
			prompt.OptionSelectedSuggestionTextColor(prompt.Blue),
			prompt.OptionDescriptionBGColor(prompt.Blue),
			prompt.OptionSuggestionBGColor(prompt.DarkGray))
		input := strings.Split(t, " ")
		switch input[0] {
		case "exit":
			fmt.Println("Bye bye!")
			os.Exit(0)
		case "back":
			return nil
		case "target":
			printTarget()
		case "interact":
			interact()
		case "payload":
			payload(kittenName)
		case "sleep":
			sleep(input, kittenName)
		case "recon":
			initChecks := http.Targets[kittenName].GetInitChecks()
			util.PrintRecon(initChecks)
		default:
			fmt.Println("Unknown command")
		}
	}
}

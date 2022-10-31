package cli

import (
	"KittyStager/cmd/config"
	"github.com/c-bata/go-prompt"
)

var c config.General

func completerCli(d prompt.Document) []prompt.Suggest {
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
			return
		case "config":
			printConfig(conf)
		case "target":
			printTarget()
		case "interact":
			interact()
		}
	}
}

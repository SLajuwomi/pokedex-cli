package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var supportedCommands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Display the help menu",
		callback:    help,
	},
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func help() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")
	for k, v := range supportedCommands {
		fmt.Printf("%s: %s\n", k, v.description)
	}
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userString := scanner.Text()
		userStringSlice := cleanInput(userString)
		_, ok := supportedCommands[userStringSlice[0]]
		if ok {
			supportedCommands[userStringSlice[0]].callback()
		} else {
			fmt.Println("Unknown command")
		}
	}

}

func cleanInput(text string) []string {
	cleaned_text := strings.ToLower(strings.TrimSpace(text))
	final_slice := strings.Split(cleaned_text, " ")
	return final_slice
}

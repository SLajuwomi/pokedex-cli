package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/slajuwomi/pokedexcli/internal"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*internal.Config, bool) error
}

var supportedCommands map[string]cliCommand

func commandExit(cfg *internal.Config, prev bool) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func help(cfg *internal.Config, prev bool) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for k, v := range supportedCommands {
		fmt.Printf("%s: %s\n", k, v.description)
	}
	return nil
}

func init() {
	supportedCommands = map[string]cliCommand{
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
		"map": {
			name:        "map",
			description: "Display next 20 location areas in Pokemon world",
			callback:    internal.GetMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display previous 20 location areas in Pokemon world",
			callback:    internal.GetMap,
		},
	}
}
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var prev bool
	cfg := internal.Config{}
	cfg.Next = "https://pokeapi.co/api/v2/location-area/"
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userString := scanner.Text()
		userStringSlice := cleanInput(userString)
		if supportedCommands[userStringSlice[0]].name == "mapb" {
			prev = true
		} else {
			prev = false
		}
		_, ok := supportedCommands[userStringSlice[0]]
		if ok {
			supportedCommands[userStringSlice[0]].callback(&cfg, prev)
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

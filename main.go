package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/slajuwomi/pokedexcli/internal"
	"github.com/slajuwomi/pokedexcli/internal/pokecache"
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
		"explore": {
			name:        "explore",
			description: "Get Pokemon in location area",
			callback:    internal.Explore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon",
			callback:    internal.Catch,
		},
		"inspect": {
			name:        "inspect",
			description: "Get information about a caught Pokemon",
			callback:    internal.Inspect,
		},
	}
}
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var prev bool
	cfg := internal.Config{}
	cfg.Next = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	cfg.Cache = pokecache.NewCache(20 * time.Second)
	cfg.Pokedex = make(map[string]internal.PokemonInformation)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userString := scanner.Text()
		userStringSlice := cleanInput(userString)
		command, ok := supportedCommands[userStringSlice[0]]
		if command.name == "mapb" {
			prev = true
		} else {
			prev = false
		}
		if command.name == "explore" {
			location := userStringSlice[1:]
			if ok {
				cfg.Location = location[0]
			} else {
				fmt.Println("expected location")
			}
		}
		if command.name == "catch" || command.name == "inspect" {
			pokemon := userStringSlice[1:]
			if ok {
				cfg.Pokemon = strings.ToLower(pokemon[0])
			} else {
				fmt.Println("expected pokemon name")
			}
		}
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

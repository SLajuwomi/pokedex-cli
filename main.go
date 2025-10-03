package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, bool) error
}

type locationAreaList struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type config struct {
	Next     string
	Previous any
}

var supportedCommands map[string]cliCommand

func commandExit(cfg *config, prev bool) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func help(cfg *config, prev bool) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for k, v := range supportedCommands {
		fmt.Printf("%s: %s\n", k, v.description)
	}
	return nil
}

func getMap(cfg *config, prev bool) error {
	var url string
	if prev {
		_, ok := cfg.Previous.(string)
		if ok {
			url = cfg.Previous.(string)
		} else {
			fmt.Println("\"map\" has not been used yet, so no previous locations to view. please use \"map\" first!")
			return nil
		}
	} else {
		url = cfg.Next
	}
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and \n body: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	currentLocationAreaList := locationAreaList{}
	err = json.Unmarshal(body, &currentLocationAreaList)
	if err != nil {
		log.Fatal(err)
	}
	for _, resultStruct := range currentLocationAreaList.Results {
		fmt.Println(resultStruct.Name)
	}
	cfg.Next = currentLocationAreaList.Next
	cfg.Previous = currentLocationAreaList.Previous
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
			callback:    getMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display previous 20 location areas in Pokemon world",
			callback:    getMap,
		},
	}
}
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var prev bool
	cfg := config{}
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

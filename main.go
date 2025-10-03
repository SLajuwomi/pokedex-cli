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
	callback    func(*config) error
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

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func help(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for k, v := range supportedCommands {
		fmt.Printf("%s: %s\n", k, v.description)
	}
	return nil
}

func getMap(cfg *config) error {
	res, err := http.Get(cfg.Next)
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
			description: "Display Pokemon Areas",
			callback:    getMap,
		},
	}
}
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := config{}
	cfg.Next = "https://pokeapi.co/api/v2/location-area/"
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userString := scanner.Text()
		userStringSlice := cleanInput(userString)
		_, ok := supportedCommands[userStringSlice[0]]
		if ok {
			supportedCommands[userStringSlice[0]].callback(&cfg)
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

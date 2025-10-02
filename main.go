package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userString := scanner.Text()
		userStringSlice := cleanInput(userString)
		fmt.Printf("Your command was: %s\n", userStringSlice[0])
	}
}

func cleanInput(text string) []string {
	cleaned_text := strings.ToLower(strings.TrimSpace(text))
	final_slice := strings.Split(cleaned_text, " ")
	return final_slice 
}


package internal

import "fmt"

func Pokedex(cfg *Config, prev bool) error {
	fmt.Println("Your Pokedex: ")
	if len(cfg.Pokedex) == 0 {
		fmt.Println("You haven't caught any Pokemon yet. Use the catch command to catch some! ;)")
		return nil
	}
	for k := range cfg.Pokedex {
		fmt.Printf(" - %s\n", k)
	}
	return nil
}

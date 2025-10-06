package internal

import "fmt"

func Inspect(cfg *Config, prev bool) error {
	value, exists := cfg.Pokedex[cfg.Pokemon]
	if exists {
		fmt.Printf("Name: %s\n", value.Name)
		fmt.Printf("Height: %d\n", value.Height)
		fmt.Printf("Weight: %d\n", value.Weight)
		fmt.Println("Stats:")
		fmt.Printf("   -hp: %d\n", value.Stats[0].BaseStat)
		fmt.Printf("   -attack: %d\n", value.Stats[1].BaseStat)
		fmt.Printf("   -defense: %d\n", value.Stats[2].BaseStat)
		fmt.Printf("   -special-attack: %d\n", value.Stats[3].BaseStat)
		fmt.Printf("   -special-defense: %d\n", value.Stats[4].BaseStat)
		fmt.Printf("   -speed: %d\n", value.Stats[5].BaseStat)
		fmt.Println("Types:")
		for _, t := range value.Types {
			fmt.Printf("   - %s\n", t.Type.Name)
		}
	} else {
		fmt.Println("you have not caught that pokemon")
	}
	return nil
}

package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type LocationInformation struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func Explore(cfg *Config, prev bool) error {
	url := baseUrl + "/location-area/" + cfg.Location

	currentLocationInformation := LocationInformation{}

	cachedValue, exists := cfg.Cache.Get(url)
	if exists {
		fmt.Println("Exploring " + cfg.Location)
		err := json.Unmarshal(cachedValue, &currentLocationInformation)
		if err != nil {
			log.Fatal(err)
		}

		for _, pokemonEncounterStruct := range currentLocationInformation.PokemonEncounters {
			fmt.Println(" - " + pokemonEncounterStruct.Pokemon.Name)
		}
		// TODO: Remove
		fmt.Println("CACHE USED")
		return nil
	}

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	fmt.Println("Exploring " + cfg.Location)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and \n body: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	cfg.Cache.Add(url, body)
	err = json.Unmarshal(body, &currentLocationInformation)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Found Pokemon: ")
	for _, pokemonEncounterStruct := range currentLocationInformation.PokemonEncounters {
		fmt.Println(" - " + pokemonEncounterStruct.Pokemon.Name)
	}
	return nil
}

package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/slajuwomi/pokedexcli/internal/pokecache"
)

type LocationAreaList struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Config struct {
	Next     string
	Previous any
	Cache    *pokecache.Cache
}

func GetMap(cfg *Config, prev bool) error {
	var url string
	currentLocationAreaList := LocationAreaList{}

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
	_, exists := cfg.Cache.CacheMap[url]
	if exists {
		value, _ := cfg.Cache.Get(url)
		err := json.Unmarshal(value, &currentLocationAreaList)
		if err != nil {
			log.Fatal(err)
		}
		for _, resultStruct := range currentLocationAreaList.Results {
			fmt.Println(resultStruct.Name)
		}
		cfg.Next = currentLocationAreaList.Next
		cfg.Previous = currentLocationAreaList.Previous
		fmt.Println("cache used")
		return nil
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
	cfg.Cache.Add(url, body)
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

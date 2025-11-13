package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"math/rand"
)


type pokemonInfo struct{
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
	Height    int `json:"height"`
	Name          string `json:"name"`
	BaseExperience int `json:"base_experience"`
}

func catchPokemon(name string,pokemonMap map[string]pokemonInfo)error{
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s",name)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var pokemon pokemonInfo
	err = json.Unmarshal(data,&pokemon)
	if err != nil {
		return err
	}
	catchChance := 100 - (pokemon.BaseExperience / 10)
	random := rand.Intn(100)
	
	if random < catchChance{
		pokemonMap[pokemon.Name] = pokemon
		fmt.Printf("%s was caught!\n",pokemon.Name)
	}else{
		fmt.Printf("%s escaped!\n",pokemon.Name)
	}
	return nil

	


}
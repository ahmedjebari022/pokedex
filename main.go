package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/ahmedjebari022/pokedex/pokecache"
)


func main(){
	initCliCommands()
	exit := false
	scanner := bufio.NewScanner(os.Stdin)
	cache := pokecache.NewCache(5*time.Second)

	for !exit {
		fmt.Print("Pokedex >")
		if !scanner.Scan(){
			break
		}
		err := scanner.Err()
		if err != nil{
			log.Fatal(err)
		}
		input := cleanInput(scanner.Text())
		
		for key,ele := range supportedCommands{
			if input[0] == key{
				argument := ""
				if len(input) > 1{
					argument = input[1]
				}
				err := ele.callback(argument)
				if err == nil {
					fmt.Printf("Unknown command")
				}
				if err.Error() == "exit"{
					os.Exit(0)
				}else if err.Error() == "map"{
					mapCommand := supportedCommands["map"]
					locations, err := getLocations(&mapCommand,cache)
					if err != nil{
							break
						}
					for _, v := range locations.Results {
						fmt.Printf("%s \n",v.Name)
					}
				}else if err.Error() == "bmap" {
					bmapCommand := supportedCommands["bmap"]
					locations, err := getLocations(&bmapCommand,cache)
					if err != nil{
						fmt.Printf("%v\n",err)
						break
					}
					for _, v := range locations.Results {
						fmt.Printf("%s \n",v.Name)
					}	
				}else if err.Error() == "explore" {
				
					pokemons, err := getPokemonsFromLocation(argument,cache)
					if err != nil {
						fmt.Printf("%v\n",err)
					}
					for _,p := range pokemons.PokemonEncounters{
						fmt.Printf("- %s\n",p.Pokemon.Name)
					} 
				}
			}
		}
	}}






func commandExit(name string)error{
	fmt.Println("Closing the Pokedex... Goodbye!")
	return fmt.Errorf("exit")
}
func commandHelp(name string)error{
	fmt.Printf("\nWelcome to the Pokedex! \nUsage:\n\n")
	for _,v := range supportedCommands{
		fmt.Printf("%s: %s\n",v.name,v.description)
	}
	return fmt.Errorf("help")
}
func commandMap(name string)error{
	return fmt.Errorf("map")
}
func commandBMap(name string)error{
	return fmt.Errorf("bmap")
}
func commandExplore(name string)error{
	fmt.Printf("Exploring %s...\n",name)
	fmt.Printf("Found Pokemon:\n")
	return fmt.Errorf("explore")
}


func getLocations(mapCommand *cliCommand,cache *pokecache.Cache)(location,error){
	url := ""
	
	switch mapCommand.name {
		case "map" :
			url = mapCommand.conf.Next
		case "bmap" :
			if mapCommand.conf.Previous == ""{return location{},fmt.Errorf("you're on the first page")}
			url = mapCommand.conf.Previous
	}
	fmt.Println(url)

	if cachedLocations, ok := cache.Get(url); ok {
		var locations location
		err := json.Unmarshal(cachedLocations,&locations)
		if err != nil{
			return location{}, err
		}
		mapCommand.conf.Next = locations.Next
		mapCommand.conf.Previous = locations.Previous
		return locations,nil
	}
	res, err := http.Get(url)
	if err != nil{
		return location{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	
	cache.Add(url,data)
	
	

	if err != nil {
		return location{}, err
	}
	
	var locations location
	if err := json.Unmarshal(data, &locations); err != nil {
		return location{},err
	}
	mapCommand.conf.Next = locations.Next
	mapCommand.conf.Previous = locations.Previous
	return locations, nil
	

}

func getPokemonsFromLocation(area string,cache *pokecache.Cache)(pokemon,error){
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s",area)

	val, ok := cache.Get(url)
	if ok {
		var pokemons pokemon
		json.Unmarshal(val,&pokemons)
		return pokemons,nil
	}
	res, err := http.Get(url)
	if err != nil {
		return pokemon{},err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return pokemon{},err
	}
	cache.Add(url,data)
	var pokemons pokemon
	err = json.Unmarshal(data,&pokemons)
	if err != nil {
		return pokemon{},err
	}

	return pokemons,nil
}
type pokemon struct{
	PokemonEncounters []struct {
		Pokemon struct{
			Name string `json:"name"`
			URL string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`

}



type location struct{
	Count 	int `json:"count"`
	Next 	string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct{
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}


type cliCommand struct{
	name 		string
	description string
	callback 	func(string) error
	conf		*config
}

type config struct{
	Next 		string
	Previous 	string
}


var supportedCommands = make(map[string]cliCommand)


func initCliCommands(){

	initConfig := &config{
		Next :"https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		Previous: "",
	}

	supportedCommands["help"] = cliCommand{name:"help",description:"Display A Help Message",callback:commandHelp}
	supportedCommands["exit"] = cliCommand{name:"exit",description:"Exit the Pokedex",callback:commandExit}
	supportedCommands["map"] = cliCommand{name:"map",description:"Display the next 20 locations",callback:commandMap,conf:initConfig}
	supportedCommands["bmap"] = cliCommand{name:"bmap",description:"Display the previous 20 locations",callback:commandBMap,conf:initConfig}
	supportedCommands["explore"] = cliCommand{name:"explore",description:"Display pokemons of an area",callback:commandExplore} 
}


package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)


func main(){
	initCliCommands()
	exit := false
	scanner := bufio.NewScanner(os.Stdin)
	
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
			for _,i := range input{
				if i == key{
					err = ele.callback()
					if err == nil {
						fmt.Println("Unknown command")
					}
					
					if err.Error() == "exit"{
						os.Exit(0)
					}else if err.Error() == "map" {
						mapCommand := supportedCommands["map"]
						locations, err := getLocations(&mapCommand)
						if err != nil{
							break
						}
						for _, v := range locations.Results {
							fmt.Printf("%s \n",v.Name)
						}
					}else if err.Error() == "bmap" {
						bmapCommand := supportedCommands["bmap"]
						locations, err := getLocations(&bmapCommand)
						if err != nil{
							fmt.Printf("%v\n",err)
							break
						}
						for _, v := range locations.Results {
							fmt.Printf("%s \n",v.Name)
						}
					}
				}
			}
		}

				

	} 
}

func commandExit()error{
	fmt.Println("Closing the Pokedex... Goodbye!")
	return fmt.Errorf("exit")
}
func commandHelp()error{
	fmt.Printf("\nWelcome to the Pokedex! \nUsage:\n\n")
	for _,v := range supportedCommands{
		fmt.Printf("%s: %s\n",v.name,v.description)
	}
	return fmt.Errorf("help")
}
func commandMap()error{
	return fmt.Errorf("map")
}
func commandBMap()error{
	return fmt.Errorf("bmap")
}

func getLocations(mapCommand *cliCommand)(location,error){
	url := ""
	
	switch mapCommand.name {
		case "map" :
			url = mapCommand.conf.Next
		case "bmap" :
			if mapCommand.conf.Previous == ""{return location{},fmt.Errorf("you're on the first page")}
			url = mapCommand.conf.Previous
	}

	res, err := http.Get(url)
	if err != nil{
		return location{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
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
	callback 	func() error
	conf		*config
}

type config struct{
	Next 		string
	Previous 	string
}


var supportedCommands = make(map[string]cliCommand)


func initCliCommands(){

	initConfig := &config{
		Next :"https://pokeapi.co/api/v2/location-area/?limit=20&offset=0",
		Previous: "",
	}

	supportedCommands["help"] = cliCommand{name:"help",description:"Display A Help Message",callback:commandHelp}
	supportedCommands["exit"] = cliCommand{name:"exit",description:"Exit the Pokedex",callback:commandExit}
	supportedCommands["map"] = cliCommand{name:"map",description:"Display the next 20 locations",callback:commandMap,conf:initConfig}
	supportedCommands["bmap"] = cliCommand{name:"bmap",description:"Display the previous 20 locations",callback:commandBMap,conf:initConfig}

}


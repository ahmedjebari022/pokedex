package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	
)


func main(){
	initCliCommands()
	exit := false
	for !exit {
		fmt.Print("Pokedex >")
		
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
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

	fmt.Printf("Welcome to the Pokedex! \nUsage:\n\n")
	for _,v := range supportedCommands{
		fmt.Printf("%s: %s\n",v.name,v.description)
	}
	return fmt.Errorf("help")
}

type cliCommand struct{
	name 		string
	description string
	callback 	func() error
}

var supportedCommands = make(map[string]cliCommand)


func initCliCommands(){
	supportedCommands["exit"] = cliCommand{name:"exit",description:"Exit the Pokedex",callback:commandExit}
	supportedCommands["help"] = cliCommand{name:"help",description:"Display A Help Message",callback:commandHelp}
}

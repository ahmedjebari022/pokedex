package main

import(
	
	"strings"
)

func cleanInput(text string) []string{
	
	trimmedInput := strings.Fields(strings.ToLower(text))
	return trimmedInput


}
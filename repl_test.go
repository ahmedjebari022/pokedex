package main

import (
	
	"testing"
)

func TestCleanInput(t *testing.T){
	cases := []struct{
		input  string
		expected []string
	}{
		{input: " Hello world",expected: []string{"hello","world"}},
		{input: "Bonjour je suis nizard. ",expected: []string{"bonjour","je","suis","nizard."}},
		{input: "Charmander Bulbasaur PIKACHU",expected: []string{"charmander", "bulbasaur", "pikachu"}},	
	}
	

	for n,c := range cases{
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected){
			t.Errorf("failed at %d test",n)
			t.Fail()
		}
		for i := range actual{
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord{
				t.Errorf("failed at %d test , Expected: %s , Got: %s",n,expectedWord,word)
				t.Fail()
			}
			
		}

		
	}


}
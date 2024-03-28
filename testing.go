package main

import (
	"fmt"

	"github.com/kljensen/snowball"
)

func main_() {
	stemmed, err := snowball.Stem("I am Andrey. I am a pupil. I am 24. I live in Moscow.", "english", false)
	if err == nil {
		fmt.Println(stemmed) // Prints "accumul"
	}
}

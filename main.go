// main.go
package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/kljensen/snowball/english"
)

var commonWords = map[string]bool{
	"a":   true,
	"an":  true,
	"the": true,
	"of":  true,
	"in":  true,
	"on":  true,
	"at":  true,
}

func main() {
	sentence := flag.String("s", "", "sentence to normalize")
	flag.Parse()

	normalizedWords := make([]string, 0)

	for _, word := range strings.Fields(*sentence) {
		// Исключаем общие слова и приводим к нижнему регистру
		if !commonWords[strings.ToLower(word)] {
			// Делаем стемминг
			normalizedWords = append(normalizedWords, english.Stem(word, false))
		}
	}

	fmt.Println(strings.Join(normalizedWords, " "))
}

// main.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/kljensen/snowball"
)

func readWordsFromFile(filename string) (map[string]bool, error) {
	wordMap := make(map[string]bool)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		for _, word := range words {
			wordMap[word] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return wordMap, nil
}

func main() {
	sentence := flag.String("s", "", "sentence to normalize")
	flag.Parse()

	wordMap, err := readWordsFromFile("unused_english_words.txt")
	if err != nil {
		panic(err)
	}

	firstSlice := make([]string, 0)
	for _, i := range strings.Fields(*sentence) {
		firstSlice = append(firstSlice, strings.ToLower(i))
	}
	// fmt.Println(strings.Join(firstSlice, " "))

	for unusedWord, _ := range wordMap {
		for i := 0; i < len(firstSlice); i++ {
			if unusedWord == firstSlice[i] {
				// fmt.Println(unusedWord, firstSlice, i, len(firstSlice))
				if i == 0 && len(firstSlice) > 1 {
					firstSlice = firstSlice[i+1:]
				} else if i == 0 && len(firstSlice) <= 1 {
					firstSlice = make([]string, 0)
				} else if i == len(firstSlice)-1 {
					firstSlice = firstSlice[:i]
				} else {
					firstSlice = append(firstSlice[:i], firstSlice[i+1:]...)
				}
				i--
			}
		}
	}
	secondSlice := make([]string, 0)
	for _, i := range firstSlice {
		stemmed, err := snowball.Stem(i, "english", true)
		if err == nil {
			// fmt.Println(stemmed, " ")
			secondSlice = append(secondSlice, stemmed)
		}
	}
	// fmt.Println(strings.Join(secondSlice, " "))
	if len(secondSlice) > 1 {
		for i := 0; i < len(secondSlice)-1; i++ {
			for j := i + 1; j < len(secondSlice); j++ {
				if secondSlice[i] == secondSlice[j] {
					if j == len(secondSlice)-1 {
						secondSlice = secondSlice[:j]
					} else if j < len(secondSlice)-1 {
						secondSlice = append(secondSlice[:j], secondSlice[j+1:]...)
						j--
					}
				}
			}
		}
	}
	fmt.Println(strings.Join(secondSlice, " "))
}

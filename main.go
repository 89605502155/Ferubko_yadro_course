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

func readWordsFromFile(filename string, bufferSize int, numLines int) <-chan []string {
	wordCh := make(chan []string)

	go func() {
		defer close(wordCh)

		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Buffer(make([]byte, bufferSize), numLines*bufferSize)

		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
			if len(lines) == numLines {
				wordCh <- lines
				lines = nil
			}
		}

		if len(lines) > 0 {
			wordCh <- lines // Отправляем оставшиеся слова в канал перед закрытием
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}()

	return wordCh
}

func main() {
	sentence := flag.String("s", "", "sentence to normalize")
	flag.Parse()

	wordCh := readWordsFromFile("unused_english_words.txt", 1024, 4)

	firstSlice := make([]string, 0)
	for _, i := range strings.Fields(*sentence) {
		firstSlice = append(firstSlice, strings.ToLower(i))
	}
	// fmt.Println(strings.Join(firstSlice, " "))

	for lines := range wordCh {
		for _, unusedWord := range lines {
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
					// fmt.Println(unusedWord, firstSlice, i, len(firstSlice))
					// fmt.Println()

					// stemmed, err := snowball.Stem(word, "english", true)
					// if err == nil {
					// 	fmt.Print(stemmed, " ggg ")
					// }
					// normalizedWords = append(normalizedWords, english.Stem(word, true))
				}
			}
		}
	}
	// fmt.Println(strings.Join(firstSlice, " "))
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

package words

import (
	"bufio"
	"flag"
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

func deleteComa(s string) string {
	res := make([]rune, 0)
	for _, el := range s {
		if el >= rune('A') && el <= rune('z') {
			res = append(res, el)
		}
	}
	return string(res)
}

func Normalization() {
	var sentence string
	flag.StringVar(&sentence, "s", "", "sentence to normalize")
	flag.Parse()

	wordMap, err := readWordsFromFile("unused_english_words.txt")
	if err != nil {
		panic(err)
	}

	myWordsMap := map[string]bool{}
	for _, i := range strings.Fields(sentence) {

		w := strings.ToLower(i)
		w = deleteComa(w)
		if wordMap[w] {
			continue
		}
		stemmed, err := snowball.Stem(w, "english", true)
		if err == nil {
			myWordsMap[stemmed] = true
		}
	}
	// for i, _ := range myWordsMap {
	// 	fmt.Printf("%s ", i)
	// }
	// fmt.Println()
}

package words

import (
	"bufio"
	"os"
	"strings"

	"github.com/kljensen/snowball"
)

type StremmingStruct struct{}

func NewStrimming() *StremmingStruct {
	return &StremmingStruct{}
}

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
		if (el >= rune('A') && el <= rune('Z')) || (el >= rune('a') && el <= rune('z')) {
			res = append(res, el)
		}
	}
	return string(res)
}

func (s *StremmingStruct) Normalization(sentence string) (*map[string]bool, error) {

	wordMap, err := readWordsFromFile("unused_english_words.txt")
	if err != nil {
		panic(err)
	}
	words := strings.Split(sentence, " ")

	// Выводим полученный слайс слов
	myWordsMap := map[string]bool{}
	for _, i := range words {

		w := strings.ToLower(i)
		w_ := deleteComa(w)
		if wordMap[w_] {
			continue
		}
		stemmed, err := snowball.Stem(w_, "english", true)
		if err == nil {
			myWordsMap[stemmed] = true
		}
	}
	return &myWordsMap, nil
}

func (s *StremmingStruct) MergeMapToString(m1, m2 *map[string]bool) []string {
	for k := range *m2 {
		(*m1)[k] = true
	}
	resp := make([]string, 0)
	for k := range *m1 {
		resp = append(resp, k)
	}
	return resp
}

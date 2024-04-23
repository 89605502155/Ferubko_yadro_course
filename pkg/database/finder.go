package database

import (
	"encoding/json"
	"os"
	"slices"
	"strconv"

	"github.com/sirupsen/logrus"

	"xkcd/pkg/xkcd"
)

type DatabaseFind struct {
	Name string
}

func NewDatabaseFinder(name string) *DatabaseFind {
	return &DatabaseFind{
		Name: name,
	}
}

func (d *DatabaseFind) read() map[string]xkcd.ComicsInfo {
	fileContent, err := os.ReadFile(d.Name)
	if err != nil {
		logrus.Fatalf("Ошибка чтения файла: %v", err)
	}

	// Определяем map для разбора JSON
	data := make(map[string]xkcd.ComicsInfo)
	// fmt.Println(string(fileContent))
	// Парсим JSON
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		logrus.Fatalf("Ошибка декодинга: %v", err)
	}
	return data
}

func (d *DatabaseFind) makeTwoSlices(word string, data *map[string]xkcd.ComicsInfo, maxLenOfResult int) []int {
	result := make([]int, 0)
	lenght := make([]int, 0)
	var mi int
	for key, valeu := range *data {
		s := 0
		for _, v := range valeu.Keywords {
			if v == word {
				s++
			}
		}
		if len(result) == 0 {
			index, _ := strconv.Atoi(key)
			result = append(result, index)
			lenght = append(lenght, s)
			mi = s
		} else if len(result) < maxLenOfResult {
			if s < mi {
				mi = s
			}
			index, _ := strconv.Atoi(key)
			result = append(result, index)
			lenght = append(lenght, s)
		} else {
			if s > mi {
				miIndex := slices.Index(lenght, mi)
				result = append(result[:miIndex], result[miIndex+1:]...)
				lenght = append(lenght[:miIndex], lenght[miIndex+1:]...)
				index, _ := strconv.Atoi(key)
				result = append(result, index)
				lenght = append(lenght, s)
				mi = slices.Min(lenght)
			}
		}
	}
	return result
}

func (d *DatabaseFind) Find(input *map[string]bool) map[string][]int {
	data := d.read()
	res := make(map[string][]int)
	for word, _ := range *input {
		res[word] = d.makeTwoSlices(word, &data, 10)

	}
	return res
}

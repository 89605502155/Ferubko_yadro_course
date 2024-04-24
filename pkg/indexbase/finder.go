package indexbase

import (
	"encoding/json"
	"os"
	"slices"
	"sort"

	"github.com/sirupsen/logrus"
)

type IndexBaseFinde struct {
	Name string
}

func NewIndexFinde(name string) *IndexBaseFinde {
	return &IndexBaseFinde{
		Name: name,
	}
}

func (ind *IndexBaseFinde) read() map[string]IndexStatistics {
	fileContent, err := os.ReadFile(ind.Name)
	if err != nil {
		logrus.Fatalf("Ошибка чтения файла: %v", err)
	}
	data := make(map[string]IndexStatistics)
	// fmt.Println(string(fileContent))
	// Парсим JSON
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		logrus.Fatalf("Ошибка при разборе JSON: %v", err)
		return map[string]IndexStatistics{}
	}
	return data
}

func (ind *IndexBaseFinde) sorted(d IndexStatistics, limit int) []int {
	res := make([]int, 0)
	copySlice := make([]int, len(d.NumberComicsOfIndex))
	copy(copySlice, d.NumberComicsOfIndex)
	sort.Sort(sort.Reverse(sort.IntSlice(copySlice)))
	// fmt.Println(copySlice)
	// fmt.Println(d.NumberComicsOfIndex)
	// fmt.Println(d.ComicsIndex)
	// fmt.Println()
	var index int
	for i := 0; i < limit; i++ {
		if i > 0 {
			if copySlice[i-1] > copySlice[i] {
				index = slices.Index(d.NumberComicsOfIndex, copySlice[i])
				res = append(res, d.ComicsIndex[index])
			} else if copySlice[i] == copySlice[i-1] {
				for j := index; j < len(d.ComicsIndex); j++ {
					if j != index && d.NumberComicsOfIndex[j] == d.NumberComicsOfIndex[index] {
						index = j
						break
					}
				}
				// index = slices.Index(d.NumberComicsOfIndex, copySlice[i])
				res = append(res, d.ComicsIndex[index])
			}
		} else if i == 0 {
			index = slices.Index(d.NumberComicsOfIndex, copySlice[i])
			res = append(res, d.ComicsIndex[index])
		}
	}
	return res
}

func (ind *IndexBaseFinde) Find(input *map[string]bool, limit int) map[string][]int {
	data := ind.read()
	res := make(map[string][]int)
	for word, _ := range *input {
		if _, ok := data[word]; ok {
			if len(data[word].ComicsIndex) <= limit {
				res[word] = data[word].ComicsIndex
			} else {
				res[word] = ind.sorted(data[word], limit)
			}
		}

	}
	return res
}

package indexbase

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"

	"github.com/sirupsen/logrus"

	"xkcd/pkg/xkcd"
)

type IndexBaseStruct struct {
	Name string
}

func NewIndexBase(name string) *IndexBaseStruct {
	return &IndexBaseStruct{
		Name: name,
	}
}

type IndexStatistics struct {
	ComicsIndex         []int `json:"comics_index"`
	NumberComicsOfIndex []int `json:"number_comics_of_index"`
}

func (ind *IndexBaseStruct) ReadBase() *map[string]IndexStatistics {
	fileContent, err := os.ReadFile(ind.Name)
	if err != nil {
		ind.CreateEmptyDatabase()
		logrus.Fatalf("Ошибка чтения файла: %v", err)
	}
	data := make(map[string]IndexStatistics)
	fmt.Println(string(fileContent))
	// Парсим JSON
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		logrus.Fatalf("Ошибка при разборе JSON: %v", err)
		return &map[string]IndexStatistics{}
	}
	return &data
}

func (ind *IndexBaseStruct) CreateEmptyDatabase() {
	file, err := os.Create(ind.Name)
	if err != nil {
		logrus.Fatalf("you have error %s", err.Error())
		return
	}
	defer file.Close()
}

func (ind *IndexBaseStruct) BuildIndexFromDB(db *map[string]xkcd.ComicsInfo, indexBase *map[string]IndexStatistics) {
	for index, comic := range *db {
		for _, words := range comic.Keywords {
			if _, ok := (*indexBase)[words]; !ok {
				slic := IndexStatistics{}
				slic.NumberComicsOfIndex = make([]int, 0)
				slic.ComicsIndex = make([]int, 0)
				(*indexBase)[words] = slic
			}
			intIndex, err := strconv.Atoi(index)
			if err != nil {
				logrus.Fatalf("you have error %s", err.Error())
				return
			}
			indexInSlice := slices.Index((*indexBase)[words].ComicsIndex, intIndex)
			if indexInSlice >= 0 {
				(*indexBase)[words].NumberComicsOfIndex[indexInSlice]++
			} else {
				str := (*indexBase)[words]
				str.ComicsIndex = append(str.ComicsIndex, intIndex)
				str.NumberComicsOfIndex = append(str.NumberComicsOfIndex, 1)
				(*indexBase)[words] = str
			}
		}
	}
}

func (ind *IndexBaseStruct) SaveIndexToFile(indexBase *map[string]IndexStatistics) {
	jsonData, err := json.MarshalIndent(*indexBase, "", "    ")
	if err != nil {
		logrus.Fatalf("you have error %s", err.Error())
		return
	}

	file, err := os.OpenFile(ind.Name, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logrus.Fatalf("you have error %s", err.Error())
		return
	}
	defer file.Close()

	// Записываем строку JSON в файл
	writer := io.Writer(file)

	// Записываем строку JSON в файл
	_, err = fmt.Fprint(writer, string(jsonData))
	if err != nil {
		logrus.Fatalf("you have error %s", err.Error())
		return
	}
}

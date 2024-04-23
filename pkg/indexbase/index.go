package indexbase

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type IndexBaseStruct struct {
	Name string
}

func NewIndexBase(name string) *IndexBaseStruct {
	return &IndexBaseStruct{
		Name: name,
	}
}

func (ind *IndexBaseStruct) ReadBase() *map[string][]int {
	fileContent, err := os.ReadFile(ind.Name)
	if err != nil {
		ind.CreateEmptyDatabase()
		logrus.Fatalf("Ошибка чтения файла: %v", err)
	}
	data := make(map[string][]int)
	fmt.Println(string(fileContent))
	// Парсим JSON
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		logrus.Fatalf("Ошибка при разборе JSON: %v", err)
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

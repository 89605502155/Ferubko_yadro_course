package database

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"

	"xkcd/pkg/xkcd"
)

type DatabaseStruct struct {
	Name string
}

func NewDatabase(name string) *DatabaseStruct {
	return &DatabaseStruct{
		Name: name,
	}
}
func (d *DatabaseStruct) ReadDatabase() *map[string]xkcd.ComicsInfo {
	fileContent, err := os.ReadFile(d.Name)
	if err != nil {
		d.CreateEmptyDatabase()
		logrus.Debugf("Ошибка чтения файла: %v", err)
	}

	// Определяем map для разбора JSON
	data := make(map[string]xkcd.ComicsInfo)
	logrus.Debug(string(fileContent))
	// Парсим JSON
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		logrus.Debugf("Ошибка при разборе JSON: %v", err)
	}
	return &data
}
func (d *DatabaseStruct) CreateEmptyDatabase() {
	file, err := os.Create(d.Name)
	if err != nil {
		logrus.Fatalf("you have error %s", err.Error())
		return
	}
	defer file.Close()
}
func (d *DatabaseStruct) WriteAllOnDatabase(data *map[string]xkcd.ComicsInfo, printOnConsole bool) {
	jsonData, err := json.MarshalIndent(*data, "", "    ")
	if err != nil {
		logrus.Fatalf("you have error %s", err.Error())
		return
	}
	if printOnConsole {
		logrus.Println(string(jsonData))
	}

	file, err := os.OpenFile(d.Name, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logrus.Fatalf("you have error %s", err.Error())
		return
	}
	defer file.Close()

	writer := io.Writer(file)

	_, err = fmt.Fprint(writer, string(jsonData))
	if err != nil {
		logrus.Fatalf("you have error %s", err.Error())
		return
	}
}

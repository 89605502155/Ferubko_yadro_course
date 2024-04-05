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
		fmt.Println(string(jsonData))
	}

	file, err := os.OpenFile(d.Name, os.O_WRONLY|os.O_CREATE, 0644)
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

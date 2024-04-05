package database

import "xkcd/pkg/xkcd"

type Database interface {
	WriteAllOnDatabase(data *map[string]xkcd.ComicsInfo, printOnConsole bool)
	CreateEmptyDatabase()
}

type JsonDatabase struct {
	Database
}

func NewJsonDatabase(name string) *JsonDatabase {
	return &JsonDatabase{
		Database: NewDatabase(name),
	}
}

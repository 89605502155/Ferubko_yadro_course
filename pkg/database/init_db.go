package database

import "xkcd/pkg/xkcd"

type Database interface {
	WriteAllOnDatabase(data *map[string]xkcd.ComicsInfo, printOnConsole bool)
	CreateEmptyDatabase()
	ReadDatabase() *map[string]xkcd.ComicsInfo
}
type FindInDB interface {
	Find(input *map[string]bool, limit int) []int
}

type JsonDatabase struct {
	Database
	FindInDB
}

func NewJsonDatabase(name string) *JsonDatabase {
	return &JsonDatabase{
		Database: NewDatabase(name),
		FindInDB: NewDatabaseFinder(name),
	}
}

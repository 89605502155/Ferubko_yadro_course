package repository

import (
	"fmt"
	"reflect"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
	comicsTable  = "comics"
	indexesTable = "indexes"
)

type Config struct {
	DBName string `sql:"-"  exclude:"true"`
	// Mode        string `sql:"mode"`
	// JournalMode string `sql:"journal_mode"`
	// Cache       string `sql:"cache"`
}

func NewSQLiteDB(cfg Config) (*sqlx.DB, error) {
	t := reflect.TypeOf(cfg)
	v := reflect.ValueOf(cfg)
	paramString := cfg.DBName
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		sqlTag := field.Tag.Get("sql")
		excludeTag := field.Tag.Get("exclude")
		if excludeTag == "true" {
			continue
		}
		// Проверяем, что значение поля является строкой перед вызовом String()
		if value.Kind() == reflect.String && value.String() != "" && value.String() != " " {
			fieldString := fmt.Sprintf("%s=%s", sqlTag, value.String())
			if i == 0 {
				paramString += "?"
				paramString += fieldString
			} else {
				paramString += "&"
				paramString += fieldString
			}
		}
	}

	db, err := sqlx.Open("sqlite3", paramString)
	if err != nil {
		return nil, err
	}
	logrus.Info("db created")

	var tables []string
	err = db.Select(&tables, "SELECT name FROM sqlite_master WHERE type='table';")
	if err != nil {
		fmt.Println("Error querying database:", err)
		return nil, err
	}

	for _, tableName := range tables {
		fmt.Println("Table:", tableName)
	}
	time.Sleep(10 * time.Second)
	return db, nil
}

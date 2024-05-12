package repository

import (
	"fmt"
	"reflect"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
	comicsTable  = "comics"
	indexesTable = "indexes"
)

type Config struct {
	DBName string `sql:"-"  exclude:"true"`
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
	logrus.Info("db start")
	return db, nil
}

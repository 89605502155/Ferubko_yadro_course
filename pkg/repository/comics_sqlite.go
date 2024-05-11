package repository

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"xkcd/pkg/xkcd"
)

type ComicsSQLite struct {
	db *sqlx.DB
}

func NewComicsSQLite(db *sqlx.DB) *ComicsSQLite {
	return &ComicsSQLite{db: db}
}

func comicsInfoFields() string {
	obj := new(xkcd.ComicsInfo)
	t := reflect.TypeOf(obj)
	var result string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		result += ", "
		result += jsonTag
	}
	return result
}

func (c *ComicsSQLite) Generate(data map[string]xkcd.ComicsInfo) error {
	tx, err := c.db.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	insertQuery := fmt.Sprintf("INSERT INTO %s (comics_id %s) VALUES ", comicsTable, comicsInfoFields())
	values := ""
	for key, value := range data {
		keywordsStr := strings.Join(value.Keywords, " ")
		values += fmt.Sprintf("('%s', '%s','%s'),", key, value.Url, keywordsStr)
	}
	values = values[:len(values)-1]
	insertQuery += values
	_, err = c.db.Exec(insertQuery)
	if err != nil {
		tx.Rollback()
		logrus.Fatal(err)
	}
	logrus.Info("Inserted")
	return tx.Commit()
}

func (c *ComicsSQLite) Create(comics string, obj xkcd.ComicsInfo) error {
	tx, err := c.db.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	insertQuery := fmt.Sprintf("INSERT INTO %s (comics_id %s) VALUES ", comicsTable, comicsInfoFields())
	keywordsStr := strings.Join(obj.Keywords, " ")
	insertQuery += fmt.Sprintf("('%s', '%s','%s')", comics, obj.Url, keywordsStr)
	_, err = c.db.Exec(insertQuery)
	if err != nil {
		tx.Rollback()
		logrus.Fatal(err)
	}
	logrus.Info("Inserted")
	return tx.Commit()
}

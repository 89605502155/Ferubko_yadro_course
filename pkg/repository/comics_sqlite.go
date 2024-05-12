package repository

import (
	"fmt"
	"strings"
	"time"

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

func (c *ComicsSQLite) Generate(data map[string]xkcd.ComicsInfo) error {
	tx, err := c.db.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	logrus.Info("Doneyyy")
	time.Sleep(2 * time.Second)
	rows, err := c.db.Queryx("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		fmt.Println("Error querying database:", err)
		return err
	}
	defer rows.Close()

	// Обрабатываем результат запроса
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			fmt.Println("Error scanning row:", err)
			return err
		}
		fmt.Println("Table:", tableName)
	}

	// Проверяем ошибки после итерации по результату
	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating rows:", err)
		return err
	}
	logrus.Info("Done")
	time.Sleep(1 * time.Second)

	insertQuery := fmt.Sprintf("INSERT INTO %s (comics_id, url, keywords) VALUES ", comicsTable)
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

	insertQuery := fmt.Sprintf("INSERT INTO %s (comics_id, url, keywords) VALUES ", comicsTable)
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

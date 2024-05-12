package repository

import (
	"fmt"
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

func (c *ComicsSQLite) Generate(data map[string]xkcd.ComicsInfo) error {
	tx, err := c.db.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	insertQuery := fmt.Sprintf("INSERT INTO %s (comics_id, url, keywords) VALUES ", comicsTable)
	values := ""
	for key, value := range data {
		for _, v := range value.Keywords {
			values += fmt.Sprintf("('%s', '%s','%s'),", key, value.Url, v)
		}
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

func (c *ComicsSQLite) Get(word map[string]bool, limit int) ([]int, error) {
	// SELECT comics_id as c FROM comics WHERE keywords in ('follow','bring','bunch','question') GROUP BY c order by count(keywords) desc LIMIT 10;
	str := MapToString(word)
	var result []int
	queryString := fmt.Sprintf("SELECT comics_id as c FROM %s WHERE keywords IN (%s) GROUP BY c ORDER BY count(keywords) DESC  LIMIT %d", comicsTable, str, limit)
	err := c.db.Select(&result, queryString)
	return result, err
}

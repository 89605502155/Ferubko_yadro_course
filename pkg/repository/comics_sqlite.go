package repository

import (
	"fmt"

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

type IndividualComics struct {
	Key        string
	ComicsInfo xkcd.ComicsInfo
}
type FinderResponse struct {
	Key    string
	Number int
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
	if len(values) > 0 {
		values = values[:len(values)-1]
	}
	insertQuery += values
	_, err = c.db.Exec(insertQuery)
	if err != nil {
		tx.Rollback()
		logrus.Fatal(err)
	}
	logrus.Info("Inserted")
	return tx.Commit()
}

func (c *ComicsSQLite) Clear() error {
	tx, err := c.db.Begin()
	if err != nil {
		tx.Rollback()
		logrus.Fatal(err, "comics_sqlite, f")
		return err
	}
	_, err = c.db.Exec(fmt.Sprintf("DELETE FROM %s", comicsTable))
	if err != nil {
		tx.Rollback()
		logrus.Fatal(err, "comics_sqlite, s")
		return err
	}
	return tx.Commit()
}

func (c *ComicsSQLite) Get(word map[string]bool, limit int) ([]int, error) {
	str := MapToString(word)
	var result []int
	queryString := fmt.Sprintf("SELECT comics_id as c FROM %s WHERE keywords IN (%s) GROUP BY c ORDER BY count(keywords) DESC  LIMIT %d", comicsTable, str, limit)
	err := c.db.Select(&result, queryString)
	return result, err
}

type readComics struct {
	ComicsID string `db:"comics_id" binding:"required"`
	Url      string `db:"url" binding:"required"`
	Keywords string `db:"keywords" binding:"required"`
}

func (c *ComicsSQLite) GetAll() (map[string]xkcd.ComicsInfo, error) {
	var keys []string
	resoult := map[string]xkcd.ComicsInfo{}
	queryString := fmt.Sprintf("SELECT comics_id as c FROM %s ORDER BY c", comicsTable)
	err := c.db.Select(&keys, queryString)
	if err != nil {
		return map[string]xkcd.ComicsInfo{}, err
	}
	for _, key := range keys {
		resoult[key] = xkcd.ComicsInfo{
			Keywords: make([]string, 0),
		}
	}
	var res []readComics
	queryString2 := fmt.Sprintf("SELECT comics_id as c, url, keywords FROM %s ORDER BY c", comicsTable)
	err = c.db.Select(&res, queryString2)
	if err != nil {
		return map[string]xkcd.ComicsInfo{}, err
	}
	for i := 0; i < len(res); i++ {
		prom := resoult[res[i].ComicsID]
		prom.Keywords = append(prom.Keywords, res[i].Keywords)
		resoult[res[i].ComicsID] = prom
	}
	return resoult, err
}

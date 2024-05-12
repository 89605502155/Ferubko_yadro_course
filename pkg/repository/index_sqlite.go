package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"xkcd/pkg/indexbase"
)

type IndexSQLite struct {
	db *sqlx.DB
}

func NewIndexSQLite(db *sqlx.DB) *IndexSQLite {
	return &IndexSQLite{db: db}
}

func (i *IndexSQLite) Generate(indexBase map[string]indexbase.IndexStatistics) error {
	tx, err := i.db.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	insertQuery := fmt.Sprintf("INSERT INTO %s (word, comics_index, number_comics_of_index) VALUES ", indexesTable)
	values := ""
	for key, value := range indexBase {
		for i := 0; i < len(value.ComicsIndex); i++ {
			values += fmt.Sprintf("('%s', '%d','%d'),", key, value.ComicsIndex[i],
				value.NumberComicsOfIndex[i])
		}
	}
	logrus.Info(values, len(values))
	if len(values) > 0 {
		values = values[:len(values)-1]
	}
	insertQuery += values
	_, err = i.db.Exec(insertQuery)
	if err != nil {
		tx.Rollback()
		logrus.Fatal(err)
	}
	logrus.Info("Inserted")
	return tx.Commit()
}

func MapToString(m map[string]bool) string {
	s := ""
	for key := range m {
		s += fmt.Sprintf("'%s',", key)
	}
	return s[:len(s)-1]
}

func (i *IndexSQLite) Get(word map[string]bool, limit int) ([]int, error) {
	str := MapToString(word)
	var result []int
	queryString := fmt.Sprintf("SELECT comics_index as c FROM %s WHERE word IN (%s) GROUP BY c ORDER BY sum(number_comics_of_index) DESC  LIMIT %d", indexesTable, str, limit)
	err := i.db.Select(&result, queryString)
	return result, err
}

type indexTable struct {
	word                string `db:"word" binding:"required"`
	comicsIndex         int    `db:"comics_index" binding:"required"`
	numberComicsOfIndex int    `db:"number_comics_of_index" binding:"required"`
}

func (i *IndexSQLite) GetAll() (map[string]indexbase.IndexStatistics, error) {
	resoult := map[string]indexbase.IndexStatistics{}
	var res []indexTable
	var words []string
	queryString := fmt.Sprintf("SELECT word as w FROM %s ORDER BY w", indexesTable)
	err := i.db.Select(&words, queryString)
	if err != nil {
		return map[string]indexbase.IndexStatistics{}, err
	}
	for _, key := range words {
		resoult[key] = indexbase.IndexStatistics{
			ComicsIndex:         make([]int, 0),
			NumberComicsOfIndex: make([]int, 0),
		}
	}
	queryString2 := fmt.Sprintf("SELECT word as w, comics_index, number_comics_of_index FROM %s ORDER BY w", indexesTable)
	err = i.db.Select(&res, queryString2)
	if err != nil {
		return map[string]indexbase.IndexStatistics{}, err
	}
	for i := 0; i < len(res); i++ {
		pro := resoult[res[i].word]
		pro.ComicsIndex = append(pro.ComicsIndex, res[i].comicsIndex)
		pro.NumberComicsOfIndex = append(pro.NumberComicsOfIndex, res[i].numberComicsOfIndex)
		resoult[res[i].word] = pro
	}
	return resoult, err
}

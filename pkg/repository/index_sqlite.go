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
	values = values[:len(values)-1]
	insertQuery += values
	_, err = i.db.Exec(insertQuery)
	if err != nil {
		tx.Rollback()
		logrus.Fatal(err)
	}
	logrus.Info("Inserted")
	return tx.Commit()
}

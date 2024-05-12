package repository

import (
	"fmt"
	"strconv"
	"strings"

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

func intSliceToString(slice []int) string {
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = strconv.Itoa(v)
	}
	result := strings.Join(strSlice, " ")
	return result
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
		comicsInd := intSliceToString(value.ComicsIndex)
		numberCom := intSliceToString(value.NumberComicsOfIndex)
		values += fmt.Sprintf("('%s', '%s','%s'),", key, comicsInd, numberCom)
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

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
	DBName          string `sql:"-"  exclude:"true"`
	Mode            string `sql:"mode"`
	JournalMode     string `sql:"journal_mode"`
	Cache           string `sql:"cache"`
	User            string `sql:"_auth_user"`
	Password        string `sql:"_auth_pass"`
	CryptoAlgorithm string `sql:"_auth_crypt"`
	CriptoSalt      string `sql:"_auth_salt"`
}

func NewSQLiteDB(cfg Config) (*sqlx.DB, error) {
	t := reflect.TypeOf(cfg)
	v := reflect.ValueOf(cfg)
	paramString := fmt.Sprintf("%s", cfg.DBName)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		sqlTag := field.Tag.Get("sql")
		excludeTag := field.Tag.Get("exclude")
		if excludeTag == "true" {
			continue
		}
		if value.String() != "" && value.String() != " " {
			fieldString := fmt.Sprintf("%s=%s", sqlTag, value.String())
			if i == 0 {
				paramString += "?_auth&"
				paramString += fieldString
			} else {
				paramString += "&"
				paramString += fieldString
			}
		}
	}

	db, err := sqlx.Open("sqlite3", paramString) //example.db?journal_mode=wal&cache=shared&mode=rwc
	if err != nil {
		logrus.Fatal(err)
	}

	defer db.Close()
	return db, nil

}

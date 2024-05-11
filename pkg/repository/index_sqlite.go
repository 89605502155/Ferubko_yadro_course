package repository

import "github.com/jmoiron/sqlx"

type IndexSQLite struct {
	db *sqlx.DB
}

func NewIndexSQLite(db *sqlx.DB) *IndexSQLite {
	return &IndexSQLite{db: db}
}

func (i *IndexSQLite) Create() error {
	i.db.Create()
	return nil
}

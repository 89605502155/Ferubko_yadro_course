package repository

import "github.com/jmoiron/sqlx"

type UserSQLite struct {
	db *sqlx.DB
}

func NewUserSQLite(db *sqlx.DB) *ComicsSQLite {
	return &ComicsSQLite{db: db}
}

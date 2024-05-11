package repository

import (
	"github.com/jmoiron/sqlx"

	"xkcd/pkg/xkcd"
)

type Comics interface {
	Generate(data map[string]xkcd.ComicsInfo) error
	Create(comics string, obj xkcd.ComicsInfo) error
}
type Index interface {
	Create() error
}
type Repository struct {
	Comics
	Index
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Index:  NewIndexSQLite(db),
		Comics: NewComicsSQLite(db),
	}
}

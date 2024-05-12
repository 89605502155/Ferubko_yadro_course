package repository

import (
	"github.com/jmoiron/sqlx"

	"xkcd/pkg/indexbase"
	"xkcd/pkg/xkcd"
)

type Comics interface {
	Generate(data map[string]xkcd.ComicsInfo) error
	Get(word map[string]bool, limit int) ([]int, error)
	Create(comics string, obj xkcd.ComicsInfo) error
}
type Index interface {
	Generate(indexBase map[string]indexbase.IndexStatistics) error
	Get(word map[string]bool, limit int) ([]int, error)
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

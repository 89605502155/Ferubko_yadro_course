package repository

import (
	"github.com/jmoiron/sqlx"

	"xkcd/pkg/xkcd"
)

type Comics interface {
	Generate(data map[string]xkcd.ComicsInfo) error
	Get(word map[string]bool, limit int) ([]int, error)
	GetAll() (map[string]xkcd.ComicsInfo, error)
	Clear() error
}
type Index interface {
	Generate(indexBase map[string]IndexStatistics) error
	Get(word map[string]bool, limit int) ([]int, error)
	GetAll() (map[string]IndexStatistics, error)
	Clear() error
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

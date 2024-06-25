package repository

import (
	"github.com/jmoiron/sqlx"

	server "xkcd"
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

type Auth interface {
	GetUser(username string) (server.User, error)
	CreateUser(user server.User) error
}
type Repository struct {
	Comics
	Index
	Auth
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Index:  NewIndexSQLite(db),
		Comics: NewComicsSQLite(db),
		Auth:   NewUserSQLite(db),
	}
}

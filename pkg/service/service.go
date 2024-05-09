package service

import (
	"context"
	"time"

	"xkcd/pkg/database"
	"xkcd/pkg/indexbase"
	"xkcd/pkg/words"
	"xkcd/pkg/xkcd"
)

type Comics interface {
	Update() error
}
type Search interface {
	SearchInDB(input string) ([]int, time.Duration, error)
	SearchInIndex(input string) ([]int, time.Duration, error)
}

type Service struct {
	Comics
	Search
}

func NewService(db *database.JsonDatabase, index *indexbase.JsonIndex, n int,
	cl *xkcd.Client, ctx context.Context, stop context.CancelFunc,
	words *words.Words, serch_limit int) *Service {
	return &Service{
		Comics: NewComicsService(db, index, n, cl, ctx, stop),
		Search: NewSearchService(words, db, serch_limit, index),
	}
}

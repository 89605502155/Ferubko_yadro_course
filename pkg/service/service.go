package service

import (
	"context"
	"time"

	server "xkcd"
	"xkcd/pkg/repository"
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
type Auth interface{
	GenerateToken(username, password string, tokenTTL time.Duration) (string, error)
	ParseToken(accessToken string) (int, string, error)
	CreateUser(user server.User) error
}

type Service struct {
	Comics
	Search
	Auth
}

func NewService(n int, cl *xkcd.Client, ctx context.Context,
	stop context.CancelFunc, repo *repository.Repository,
	words *words.Words, serch_limit int) *Service {
	return &Service{
		Comics: NewComicsService(n, cl, ctx, stop, repo),
		Search: NewSearchService(words, serch_limit, repo),
		Auth: NewAuthService(repo),
	}
}

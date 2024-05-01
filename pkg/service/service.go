package service

import (
	"context"

	"xkcd/pkg/database"
	"xkcd/pkg/indexbase"
	"xkcd/pkg/xkcd"
)

type Comics interface {
	Update()
}

type Service struct {
	Comics
}

func NewService(db *database.JsonDatabase, index *indexbase.JsonIndex, n int,
	cl *xkcd.Client, ctx context.Context, stop context.CancelFunc) *Service {
	return &Service{
		Comics: NewComicsService(db,index,n,cl,ctx,stop),
	
	}
}

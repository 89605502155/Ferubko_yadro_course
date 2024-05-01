package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/viper"

	"xkcd/pkg/database"
	"xkcd/pkg/indexbase"
	"xkcd/pkg/worker"
	"xkcd/pkg/xkcd"
)

type ComicsService struct {
	db    *database.JsonDatabase
	index *indexbase.JsonIndex
	n     int
	cl    *xkcd.Client
	ctx   context.Context
	stop  context.CancelFunc
}

func NewComicsService(db *database.JsonDatabase, index *indexbase.JsonIndex, n int,
	cl *xkcd.Client, ctx context.Context, stop context.CancelFunc) *ComicsService {
	return &ComicsService{
		db:    db,
		index: index,
		n:     n,
		cl:    cl,
		ctx:   ctx,
		stop:  stop,
	}
}

func (s *ComicsService) Update() error {
	data := s.db.Database.ReadDatabase()

	worker.WorkerPool(s.cl, s.n, viper.GetInt("parallel"), data, s.ctx, s.stop)
	defer func(db *database.JsonDatabase, index *indexbase.JsonIndex) {
		fmt.Println("Gangut")
		db.Database.CreateEmptyDatabase()
		db.Database.WriteAllOnDatabase(data, false)
		indexes := index.IndexBase.ReadBase()

		index.IndexBase.BuildIndexFromDB(data, indexes)
		// fmt.Println(indexes)
		index.IndexBase.SaveIndexToFile(indexes)
		fmt.Println("Davu")
	}(s.db, s.index)
	if s.n == 123 {
		return errors.New("123")
	}
	return nil
}

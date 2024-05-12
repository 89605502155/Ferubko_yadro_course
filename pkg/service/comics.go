package service

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"xkcd/pkg/database"
	"xkcd/pkg/indexbase"
	"xkcd/pkg/repository"
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
	repo  *repository.Repository
}

func NewComicsService(db *database.JsonDatabase, index *indexbase.JsonIndex, n int,
	cl *xkcd.Client, ctx context.Context, stop context.CancelFunc,
	repo *repository.Repository) *ComicsService {
	return &ComicsService{
		db:    db,
		index: index,
		n:     n,
		cl:    cl,
		ctx:   ctx,
		stop:  stop,
		repo:  repo,
	}
}

func (s *ComicsService) Update() error {
	data := s.db.Database.ReadDatabase()

	worker.WorkerPool(s.cl, s.n, viper.GetInt("parallel"), data, s.ctx, s.stop)
	defer func(db *database.JsonDatabase, index *indexbase.JsonIndex,
		repo *repository.Repository) {
		logrus.Println("Gangut")
		db.Database.CreateEmptyDatabase()
		db.Database.WriteAllOnDatabase(data, false)
		indexes := index.IndexBase.ReadBase()

		index.IndexBase.BuildIndexFromDB(data, indexes)
		index.IndexBase.SaveIndexToFile(indexes)
		logrus.Println("Davu")
		repo.Comics.Generate(*data)
		repo.Index.Generate(*indexes)
	}(s.db, s.index, s.repo)
	if s.n == -123 {
		return errors.New("123")
	}
	return nil
}

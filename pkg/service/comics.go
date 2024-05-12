package service

import (
	"context"
	"errors"
	"slices"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"xkcd/pkg/indexbase"
	"xkcd/pkg/repository"
	"xkcd/pkg/worker"
	"xkcd/pkg/xkcd"
)

type ComicsService struct {
	n    int
	cl   *xkcd.Client
	ctx  context.Context
	stop context.CancelFunc
	repo *repository.Repository
}

func NewComicsService(n int, cl *xkcd.Client,
	ctx context.Context, stop context.CancelFunc,
	repo *repository.Repository) *ComicsService {
	return &ComicsService{
		n:    n,
		cl:   cl,
		ctx:  ctx,
		stop: stop,
		repo: repo,
	}
}

func (s *ComicsService) Update() error {
	data, _ := s.repo.Comics.GetAll()

	worker.WorkerPool(s.cl, s.n, viper.GetInt("parallel"), &data, s.ctx, s.stop)
	defer func(repo *repository.Repository) {
		logrus.Println("Gangut")
		indexes, _ := s.repo.Index.GetAll()
		s.BuildIndexFromDB(&data, &indexes)
		logrus.Println("Davu")
		logrus.Println(data)
		repo.Comics.Generate(data)
		repo.Index.Generate(indexes)
	}(s.repo)
	if s.n == -123 {
		return errors.New("123")
	}
	return nil
}

func (s *ComicsService) BuildIndexFromDB(db *map[string]xkcd.ComicsInfo,
	indexBase *map[string]indexbase.IndexStatistics) {
	for index, comic := range *db {
		for _, words := range comic.Keywords {
			if _, ok := (*indexBase)[words]; !ok {
				slic := indexbase.IndexStatistics{}
				slic.NumberComicsOfIndex = make([]int, 0)
				slic.ComicsIndex = make([]int, 0)
				(*indexBase)[words] = slic
			}
			intIndex, err := strconv.Atoi(index)
			if err != nil {
				logrus.Fatalf("you have error %s", err.Error())
				return
			}
			indexInSlice := slices.Index((*indexBase)[words].ComicsIndex, intIndex)
			if indexInSlice >= 0 {
				(*indexBase)[words].NumberComicsOfIndex[indexInSlice]++
			} else {
				str := (*indexBase)[words]
				str.ComicsIndex = append(str.ComicsIndex, intIndex)
				str.NumberComicsOfIndex = append(str.NumberComicsOfIndex, 1)
				(*indexBase)[words] = str
			}
		}
	}
}

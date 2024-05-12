package service

import (
	"time"

	"github.com/sirupsen/logrus"

	"xkcd/pkg/database"
	"xkcd/pkg/indexbase"
	"xkcd/pkg/repository"
	"xkcd/pkg/words"
)

type SearchService struct {
	words       *words.Words
	db          *database.JsonDatabase
	serch_limit int
	index       *indexbase.JsonIndex
	repo        *repository.Repository
}

func NewSearchService(words *words.Words, db *database.JsonDatabase, serch_limit int,
	index *indexbase.JsonIndex, repo *repository.Repository) *SearchService {
	return &SearchService{
		words:       words,
		db:          db,
		serch_limit: serch_limit,
		index:       index,
		repo:        repo,
	}
}

func (s *SearchService) SearchInDB(input string) ([]int, time.Duration, error) {
	inputData, err := s.words.Stremming.Normalization(input)
	if err != nil {
		logrus.Fatalf("you have error %s", err.Error())
		return nil, 0, err
	}
	logrus.Println("input string ", inputData)

	firstFind := s.db.FindInDB.Find(inputData, s.serch_limit)
	a1 := time.Now()
	fourthFind, err := s.repo.Comics.Get(*inputData, s.serch_limit)
	if err != nil {
		return nil, 0, err
	}
	a2 := time.Now()
	aDelta := a2.Sub(a1)
	logrus.Println("first find ", firstFind, aDelta)
	logrus.Println("fourth find ", fourthFind, aDelta)
	return fourthFind, aDelta, nil
}

func (s *SearchService) SearchInIndex(input string) ([]int, time.Duration, error) {
	inputData, err := s.words.Stremming.Normalization(input)
	if err != nil {
		logrus.Fatalf("you have error %s", err.Error())
		return nil, 0, err
	}

	secondFind := s.index.IndexFind.Find(inputData, s.serch_limit)
	b1 := time.Now()
	thirdFind, err := s.repo.Index.Get(*inputData, s.serch_limit)
	if err != nil {
		return nil, 0, err
	}
	b2 := time.Now()
	bDelta := b2.Sub(b1)
	logrus.Println("second find ", secondFind)
	logrus.Println("third find ", thirdFind)
	return thirdFind, bDelta, nil
}

package service

import (
	"time"

	"github.com/sirupsen/logrus"

	"xkcd/pkg/database"
	"xkcd/pkg/indexbase"
	"xkcd/pkg/words"
)

type SearchService struct {
	words       *words.Words
	db          *database.JsonDatabase
	serch_limit int
	index       *indexbase.JsonIndex
}

func NewSearchService(words *words.Words, db *database.JsonDatabase, serch_limit int,
	index *indexbase.JsonIndex) *SearchService {
	return &SearchService{
		words:       words,
		db:          db,
		serch_limit: serch_limit,
		index:       index,
	}
}

func (s *SearchService) SearchInDB(input string) ([]int, time.Duration, error) {
	inputData, err := s.words.Stremming.Normalization(input)
	if err != nil {
		logrus.Fatalf("you have error %s", err.Error())
		return nil, 0, err
	}
	logrus.Println("input string ", inputData)
	a1 := time.Now()
	firstFind := s.db.FindInDB.Find(inputData, s.serch_limit)
	a2 := time.Now()
	aDelta := a2.Sub(a1)
	logrus.Println("first find ", firstFind, aDelta)
	return firstFind, aDelta, nil
}

func (s *SearchService) SearchInIndex(input string) ([]int, time.Duration, error) {
	inputData, err := s.words.Stremming.Normalization(input)
	if err != nil {
		logrus.Fatalf("you have error %s", err.Error())
		return nil, 0, err
	}
	b1 := time.Now()
	secondFind := s.index.IndexFind.Find(inputData, s.serch_limit)
	b2 := time.Now()
	bDelta := b2.Sub(b1)
	logrus.Println("second find ", secondFind)
	return secondFind, bDelta, nil
}

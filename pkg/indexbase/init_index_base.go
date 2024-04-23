package indexbase

import "xkcd/pkg/xkcd"

type IndexBase interface {
	CreateEmptyDatabase()
	ReadBase() *map[string]IndexStatistics
	BuildIndexFromDB(db *map[string]xkcd.ComicsInfo, indexBase *map[string]IndexStatistics)
	SaveIndexToFile(indexBase *map[string]IndexStatistics)
}

type JsonIndex struct {
	IndexBase
}

func NewJsonIndex(name string) *JsonIndex {
	return &JsonIndex{
		IndexBase: NewIndexBase(name),
	}

}

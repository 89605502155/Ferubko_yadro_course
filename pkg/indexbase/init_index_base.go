package indexbase

import "xkcd/pkg/xkcd"

type IndexBase interface {
	CreateEmptyDatabase()
	ReadBase() *map[string]IndexStatistics
	BuildIndexFromDB(db *map[string]xkcd.ComicsInfo, indexBase *map[string]IndexStatistics)
	SaveIndexToFile(indexBase *map[string]IndexStatistics)
}

type IndexFind interface {
	Find(input *map[string]bool, limit int) map[string][]int
}

type JsonIndex struct {
	IndexBase
	IndexFind
}

func NewJsonIndex(name string) *JsonIndex {
	return &JsonIndex{
		IndexBase: NewIndexBase(name),
		IndexFind: NewIndexFinde(name),
	}

}

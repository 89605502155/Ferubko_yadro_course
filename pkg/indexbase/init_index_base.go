package indexbase

import "xkcd/pkg/xkcd"

type IndexBase interface {
	CreateEmptyDatabase()
	ReadBase() *map[string][]int
	BuildIndexFromDB(db *map[string]xkcd.ComicsInfo, index *map[string][]int)
	SaveIndexToFile(indexBase *map[string][]int)
}

type JsonIndex struct {
	IndexBase
}

func NewJsonIndex(name string) *JsonIndex {
	return &JsonIndex{
		IndexBase: NewIndexBase(name),
	}

}

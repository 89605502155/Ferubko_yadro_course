package indexbase

type IndexBase interface {
	CreateEmptyDatabase()
	ReadBase() *map[string][]int
}

type JsonIndex struct {
	IndexBase
}

func NewJsonIndex(name string) *JsonIndex {
	return &JsonIndex{
		IndexBase: NewIndexBase(name),
	}

}

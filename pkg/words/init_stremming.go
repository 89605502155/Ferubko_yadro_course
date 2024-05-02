package words

type Stremming interface {
	Normalization(sentence string) (*map[string]bool, error)
	MergeMapToString(m1, m2 *map[string]bool) []string
}

type Words struct {
	Stremming
}

func NewWordsStremming() *Words {
	return &Words{
		Stremming: NewStrimming(),
	}
}

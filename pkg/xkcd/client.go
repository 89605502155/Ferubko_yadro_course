package xkcd

import "xkcd/pkg/words"

type ClientInterface interface {
	GetComics(comicID int) (map[int]ComicsInfo, int, error)
}
type Client struct {
	ClientInterface
}

func NewClient(cl string, w *words.Words) *Client {
	return &Client{
		ClientInterface: NewHttpClient(cl, w),
	}
}

package xkcd

import "xkcd/pkg/words"

type ClientInterface interface {
	GetLatestComicsNumber() (int, error)
	GetComics(comicID int) (*map[string]ComicsInfo, error)
}
type Client struct {
	ClientInterface
}

func NewClient(cl string, w *words.Words) *Client {
	return &Client{
		ClientInterface: NewHttpClient(cl, w),
	}
}

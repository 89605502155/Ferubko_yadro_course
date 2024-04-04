package xkcd

type ClientInterface interface {
	GetLatestComicNumber() (int, error)
	GetComic(int) (*Comic, error)
}
type Client struct {
	ClientInterface
}

func NewClient(cl string) *Client {
	return &Client{
		ClientInterface: NewHttpClient(cl),
	}
}

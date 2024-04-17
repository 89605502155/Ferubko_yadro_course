package xkcd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"xkcd/pkg/words"
)

type Comics struct {
	Num        int      `json:"num"`
	Img        string   `json:"img"`
	SafeTitle  string   `json:"safe_title"`
	Transcript string   `json:"transcript"`
	Alt        string   `json:"alt"`
	Title      string   `json:"title"`
	Day        string   `json:"day"`
	Month      string   `json:"month"`
	Year       string   `json:"year"`
	Link       string   `json:"link"`
	News       string   `json:"news"`
	Errors     []string `json:"errors"`
}

type ComicsInfo struct {
	Url      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

type HttpClient struct {
	baseURL string
	w       *words.Words
}

func NewHttpClient(baseURL string, w *words.Words) *HttpClient {
	return &HttpClient{
		baseURL: baseURL,
		w:       w,
	}

}

func (c *HttpClient) GetLatestComicsNumber() (int, error) {
	url := fmt.Sprintf("%s/info.0.json", c.baseURL)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var comics Comics

	if err := json.NewDecoder(resp.Body).Decode(&comics); err != nil {
		return 0, err
	}

	return comics.Num, nil
}

func (c *HttpClient) GetComics(comicID int) (*map[int]ComicsInfo, error, int) {
	url := fmt.Sprintf("%s/%d/info.0.json", c.baseURL, comicID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err, -1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode), resp.StatusCode
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err, -2
	}

	var comics Comics
	if err := json.Unmarshal(body, &comics); err != nil {
		return nil, err, -3
	}
	var comicsInfo ComicsInfo
	ret := make(map[int]ComicsInfo)
	comicsInfo.Url = comics.Img
	map1, _ := c.w.Normalization(comics.Transcript + " " + comics.Alt)
	resp_ := make([]string, 0)
	for k := range *map1 {
		resp_ = append(resp_, k)
	}
	comicsInfo.Keywords = resp_
	ret[comics.Num] = comicsInfo

	return &ret, nil, 200
}

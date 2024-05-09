package xkcd

import (
	"encoding/json"
	"fmt"
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

	// a:=http.Client{
	// 	Timeout: 10*time.Second,
	// }

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

func (c *HttpClient) GetComics(comicID int) (map[int]ComicsInfo, int, error) {
	url := fmt.Sprintf("%s/%d/info.0.json", c.baseURL, comicID)

	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var comics Comics
	// var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&comics)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return nil, 0, nil
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

	return ret, 200, nil
}

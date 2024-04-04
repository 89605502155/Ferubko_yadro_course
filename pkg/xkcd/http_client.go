package xkcd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Comic struct {
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

type HttpClient struct {
	baseURL string
}

func NewHttpClient(baseURL string) *HttpClient {
	return &HttpClient{
		baseURL: baseURL,
	}

}

func (c *HttpClient) GetLatestComicNumber() (int, error) {
	// url_ := "https://xkcd.com/info.0.json"
	url := fmt.Sprintf("%s/info.0.json", c.baseURL)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var comicInfo Comic

	if err := json.NewDecoder(resp.Body).Decode(&comicInfo); err != nil {
		return 0, err
	}

	return comicInfo.Num, nil
}

func (c *HttpClient) GetComic(comicID int) (*Comic, error) {
	url := fmt.Sprintf("%s/%d/info.0.json", c.baseURL, comicID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var comic Comic
	if err := json.Unmarshal(body, &comic); err != nil {
		return nil, err
	}

	return &comic, nil
}

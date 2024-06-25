package xkcd

import (
	"reflect"
	"slices"
	"testing"

	"xkcd/pkg/words"
)

func TestNewHttpClient(t *testing.T) {
	testTable := []struct {
		cl       string
		w        *words.Words
		expected *HttpClient
	}{
		{
			cl: "https://xkcd.com",
			w:  words.NewWordsStremming(),
			expected: &HttpClient{
				baseURL: "https://xkcd.com",
				w:       words.NewWordsStremming(),
			},
		},
	}
	count := 0
	for _, test := range testTable {
		res := NewHttpClient(test.cl, test.w)
		if reflect.DeepEqual(*res, *test.expected) == false {
			t.Errorf("expected %v, got %v", test.expected, res)
		} else {
			count++
		}
	}
	t.Logf("good tests: %d from: %d", count, len(testTable))
}

type getCommics struct {
	comicsId int
	info     ComicsInfo
	codeInfo int
}

func TestGetComics(t *testing.T) {
	testTable := []struct {
		cl   string
		w    *words.Words
		data []getCommics
	}{
		{
			cl: "https://xkcd.com",
			w:  words.NewWordsStremming(),
			data: []getCommics{
				{
					1,
					ComicsInfo{
						Url:      "https://imgs.xkcd.com/comics/barrel_cropped_(1).jpg",
						Keywords: []string{"float", "oceanboy", "drift", "can", "be", "seenalt", "sit", "where", "dont", "barrel", "i", "distanc", "els", "boy", "wonder", "ill", "nextth"},
					},
					200,
				},
				{
					404,
					ComicsInfo{
						Url:      "",
						Keywords: []string{},
					},
					404,
				},
				{
					4404,
					ComicsInfo{
						Url:      "",
						Keywords: []string{},
					},
					404,
				},
				{
					0,
					ComicsInfo{
						Url:      "",
						Keywords: []string{},
					},
					404,
				},
			},
		},
	}
	count := 0
	for _, test := range testTable {
		client := NewHttpClient(test.cl, test.w)
		for i := 0; i < len(test.data); i++ {
			res, httpCode, err := client.GetComics(test.data[i].comicsId)
			slices.Sort(res[1].Keywords)
			slices.Sort(test.data[0].info.Keywords)
			if httpCode != 200 {
				if httpCode != test.data[i].codeInfo {
					t.Error("http code not 200 ", httpCode)
				} else {
					count++
					continue
				}

			} else if err != nil {
				t.Error("error not nil ", err)
			} else if reflect.DeepEqual(res[1], test.data[i].info) == false {
				t.Errorf("expected %v, got %v", test.data[i].info, res[1])
			} else {
				count++
			}
		}
	}
	t.Logf("good tests:  %d from:  %d", count, len(testTable[0].data))
}

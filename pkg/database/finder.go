package database

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"

	"xkcd/pkg/xkcd"
)

type DatabaseFind struct {
	Name string
}

func NewDatabaseFinder(name string) *DatabaseFind {
	return &DatabaseFind{
		Name: name,
	}
}

func (d *DatabaseFind) read() map[string]xkcd.ComicsInfo {
	fmt.Println("Чтение файла")
	fileContent, err := os.ReadFile(d.Name)
	if err != nil {
		logrus.Fatalf("Ошибка чтения файла: %v", err)
	}

	// Определяем map для разбора JSON
	data := make(map[string]xkcd.ComicsInfo)
	// fmt.Println(string(fileContent))
	// Парсим JSON
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		logrus.Fatalf("Ошибка декодинга: %v", err)
	}
	return data
}

type IndividualComics struct {
	Key        string
	ComicsInfo xkcd.ComicsInfo
}
type FinderResponse struct {
	Key    string
	Number int
}

func (d *DatabaseFind) Find(input *map[string]bool, limit int) []int {
	data := d.read()
	fmt.Println("read file")
	length := make([]int, 0)
	keySlice := make([]string, 0)
	var wg sync.WaitGroup
	comicsChan := make(chan IndividualComics, 1800)
	responseChan := make(chan FinderResponse, 1800)
	numWorkers := 10

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for comics := range comicsChan {
				fmt.Println(i, comics.Key)
				s := 0
				for _, v := range comics.ComicsInfo.Keywords {
					if (*input)[v] {
						s++
					}
				}
				// time.Sleep(time.Second)
				responseChan <- FinderResponse{
					Key:    comics.Key,
					Number: s,
				}
			}
		}()
	}

	for key, value := range data {
		comicsChan <- IndividualComics{
			Key:        key,
			ComicsInfo: value,
		}
	}
	close(comicsChan)

	go func() {
		wg.Wait()
		close(responseChan)
	}()

	for response := range responseChan {
		if response.Number > 0 {
			length = append(length, response.Number)
			keySlice = append(keySlice, response.Key)
		}
	}

	copySlice := make([]int, len(length))
	copy(copySlice, length)
	sort.Sort(sort.Reverse(sort.IntSlice(copySlice)))

	res := make([]int, 0)
	var index int
	// fmt.Println(copySlice)
	for i := 0; i < limit && i < len(copySlice); i++ {
		if i > 0 {
			if copySlice[i-1] > copySlice[i] {
				index = slices.Index(length, copySlice[i])
				k, _ := strconv.Atoi(keySlice[index])
				res = append(res, k)
			} else if copySlice[i] == copySlice[i-1] {
				for j := index; j < len(keySlice); j++ {
					if j != index && length[j] == length[index] {
						index = j
						break
					}
				}
				// index = slices.Index(d.NumberComicsOfIndex, copySlice[i])
				k, _ := strconv.Atoi(keySlice[index])
				res = append(res, k)
			}
		} else if i == 0 {
			index = slices.Index(length, copySlice[i])
			k, _ := strconv.Atoi(keySlice[index])
			res = append(res, k)
		}
	}
	// fmt.Println(copySlice[:31])
	return res
}

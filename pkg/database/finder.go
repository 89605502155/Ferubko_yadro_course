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

// func (d *DatabaseFind) makeTwoSlices(word string, data *map[string]xkcd.ComicsInfo, maxLenOfResult int) []int {
// 	result := make([]int, 0)
// 	lenght := make([]int, 0)
// 	var mi int
// 	for key, valeu := range *data {
// 		s := 0
// 		for _, v := range valeu.Keywords {
// 			if v == word {
// 				s++
// 			}
// 		}
// 		if len(result) == 0 && s > 0 {
// 			index, _ := strconv.Atoi(key)
// 			result = append(result, index)
// 			lenght = append(lenght, s)
// 			mi = s
// 		} else if len(result) < maxLenOfResult && s > 0 {
// 			if s < mi {
// 				mi = s
// 			}
// 			index, _ := strconv.Atoi(key)
// 			result = append(result, index)
// 			lenght = append(lenght, s)
// 		} else if len(result) >= maxLenOfResult && s > 0 {
// 			if s > mi {
// 				miIndex := slices.Index(lenght, mi)
// 				result = append(result[:miIndex], result[miIndex+1:]...)
// 				lenght = append(lenght[:miIndex], lenght[miIndex+1:]...)
// 				index, _ := strconv.Atoi(key)
// 				result = append(result, index)
// 				lenght = append(lenght, s)
// 				mi = slices.Min(lenght)
// 			}
// 		}
// 	}
// 	return result
// }

type IndividualComics struct {
	Key        string
	ComicsInfo xkcd.ComicsInfo
}
type FinderResponse struct {
	Key    string
	Number int
}

// func (d *DatabaseFind) Find(input *map[string]bool, limit int) []int {
// 	data := d.read()
// 	length := make([]int, 0)
// 	keyClice := make([]string, 0)
// 	// var mu sync.Mutex
// 	var wg sync.WaitGroup
// 	comicsChan := make(chan IndividualComics, 100)
// 	responseChan := make(chan FinderResponse, 100)
// 	numWorkers := 10

// 	for i := 0; i < numWorkers; i++ {
// 		wg.Add(1)
// 		fmt.Println(i, "jjj")
// 		go func() {
// 			defer wg.Done()
// 			// mu.Lock()
// 			comics := <-comicsChan
// 			// mu.Unlock()
// 			fmt.Println(comics.Key, "read")
// 			s := 0
// 			for _, v := range comics.ComicsInfo.Keywords {
// 				if (*input)[v] {
// 					s++
// 				}
// 			}
// 			fmt.Println("Calcuate", comics.Key)
// 			// mu.Lock()
// 			responseChan <- FinderResponse{
// 				Key:    comics.Key,
// 				Number: s,
// 			}
// 			// mu.Unlock()
// 			fmt.Println("Send r", comics.Key)
// 		}()
// 	}
// 	fmt.Println("Sparta")
// 	numberWorkers := 0
// 	for key, valeu := range data {
// 		fmt.Println(key, "start")
// 		if numberWorkers < numWorkers {
// 			// mu.Lock()
// 			comicsChan <- IndividualComics{
// 				Key:        key,
// 				ComicsInfo: valeu,
// 			}
// 			// mu.Unlock()
// 			fmt.Println(key)
// 			numberWorkers++
// 		} else {
// 			exit := false
// 			for {
// 				// mu.Lock()
// 				select {
// 				case response := <-responseChan:
// 					fmt.Println(key, "+")
// 					if response.Number > 0 {
// 						length = append(length, response.Number)
// 						keyClice = append(keyClice, response.Key)
// 					}
// 					// mu.Lock()
// 					comicsChan <- IndividualComics{
// 						Key:        key,
// 						ComicsInfo: valeu,
// 					}
// 					// mu.Unlock()
// 					fmt.Println(key, "-")
// 					exit = true
// 				default:
// 				}
// 				if exit {
// 					break
// 				}
// 				// mu.Unlock()
// 			}

// 		}

// 	}
// 	fmt.Println("Lotaring")
// 	copySlice := make([]int, len(length))
// 	copy(copySlice, length)
// 	sort.Sort(sort.Reverse(sort.IntSlice(copySlice)))

// 	res := make([]int, 0)
// 	var index int
// 	fmt.Println(copySlice)
// 	for i := 0; i < limit; i++ {
// 		if i > 0 {
// 			if copySlice[i-1] > copySlice[i] {
// 				index = slices.Index(length, copySlice[i])
// 				k, _ := strconv.Atoi(keyClice[index])
// 				res = append(res, k)
// 			} else if copySlice[i] == copySlice[i-1] {
// 				for j := index; j < len(keyClice); j++ {
// 					if j != index && length[j] == length[index] {
// 						index = j
// 						break
// 					}
// 				}
// 				// index = slices.Index(d.NumberComicsOfIndex, copySlice[i])
// 				k, _ := strconv.Atoi(keyClice[index])
// 				res = append(res, k)
// 			}
// 		} else if i == 0 {
// 			index = slices.Index(length, copySlice[i])
// 			k, _ := strconv.Atoi(keyClice[index])
// 			res = append(res, k)
// 		}
// 	}
// 	fmt.Println("Trosna")
// 	wg.Wait()

// 	// for word, _ := range *input {
// 	// 	res[word] = d.makeTwoSlices(word, &data, limit)

//		// }
//		return res
//	}
func (d *DatabaseFind) Find(input *map[string]bool, limit int) []int {
	data := d.read()
	length := make([]int, 0)
	keySlice := make([]string, 0)
	var wg sync.WaitGroup
	comicsChan := make(chan IndividualComics, 180)
	responseChan := make(chan FinderResponse, 180)
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

	// sort.SliceStable(length, func(i, j int) bool {
	// 	return length[i] > length[j]
	// })

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

	return res
}

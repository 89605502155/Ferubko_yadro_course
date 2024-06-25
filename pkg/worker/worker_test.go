package worker

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"testing"

	"github.com/sirupsen/logrus"

	"xkcd/pkg/words"
	"xkcd/pkg/xkcd"
)

const (
	source_url = "https://xkcd.com"
)

func readJson(filePath string) (map[string]xkcd.ComicsInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logrus.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		logrus.Println("Error reading file:", err)
		return nil, err
	}
	var expactAll map[string]xkcd.ComicsInfo
	err = json.Unmarshal(byteValue, &expactAll)
	if err != nil {
		logrus.Println("Error unmarshalling JSON:", err)
	}
	return expactAll, nil
}
func TestWorkerPool(t *testing.T) {
	testData := []struct {
		parallel      int
		numIterations int
	}{
		{5, 10}, {5, 2}, {100, 4}, {1, 409}, {2, 0}, {2, -1},
	}
	filePath := "./../../database.json"
	expactAll, _ := readJson(filePath)
	for i := 0; i < len(testData); i++ {
		words := words.NewWordsStremming()
		cl := xkcd.NewClient(source_url, words)
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		data := map[string]xkcd.ComicsInfo{}
		parallel := testData[i].parallel
		numIterations := testData[i].numIterations
		WorkerPool(cl, numIterations, parallel, &data, ctx, stop)
		expected := map[string]xkcd.ComicsInfo{}
		for j := 1; j < numIterations+1; j++ {
			a := strconv.Itoa(j)
			expected[a] = expactAll[a]
		}
		if reflect.DeepEqual(data, expected) {
			t.Log("Test True")
		} else {
			t.Log("Test False")
		}
	}

}

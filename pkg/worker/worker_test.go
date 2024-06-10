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

func TestWorkerPool(t *testing.T) {
	words := words.NewWordsStremming()

	cl := xkcd.NewClient(source_url, words)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	data := map[string]xkcd.ComicsInfo{}
	parallel := 5
	numIterations := 10
	filePath := "./../../database.json"
	file, err := os.Open(filePath)
	if err != nil {
		logrus.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		logrus.Println("Error reading file:", err)
		return
	}
	var expactAll map[string]xkcd.ComicsInfo
	err = json.Unmarshal(byteValue, &expactAll)
	if err != nil {
		logrus.Println("Error unmarshalling JSON:", err)
	}
	WorkerPool(cl, numIterations, parallel, &data, ctx, stop)
	expected := map[string]xkcd.ComicsInfo{}
	for i := 1; i < numIterations+1; i++ {
		a := strconv.Itoa(i)
		expected[a] = expactAll[a]
	}
	if reflect.DeepEqual(data, expected) {
		t.Log("Test True")
	} else {
		t.Log("Test False")
	}

}

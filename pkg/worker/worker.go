package worker

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	"xkcd/pkg/xkcd"
)

func WorkerPool(cl *xkcd.Client, numIter int, numWorkers int, data *map[string]xkcd.ComicsInfo,
	ctx context.Context, stop context.CancelFunc) {

	keyChan := make(chan int, numWorkers*30)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for key := range keyChan {
				res, _, err := cl.ClientInterface.GetComics(key)
				logrus.Println(i, key, err)
				if err != nil {
					logrus.Println(err)
					continue

				}
				mu.Lock()
				(*data)[fmt.Sprintf("%d", key)] = res[key]
				mu.Unlock()
			}
		}()
	}
	for key := 1; key < numIter; key++ {
		select {
		case <-ctx.Done():
			logrus.Println("Pressish-Eilay")
			stop()
			logrus.Println("Bagram")
			key = numIter
		default:
			if _, ok := (*data)[fmt.Sprintf("%d", key)]; ok {
				key++
				continue
			}
			keyChan <- key
		}
	}
	close(keyChan)

	go func() {
		wg.Wait()
	}()
}

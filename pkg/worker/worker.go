package worker

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"xkcd/pkg/xkcd"
)

func WorkerPool(cl *xkcd.Client, numIter int, numWorkers int, data *map[string]xkcd.ComicsInfo,
	ctx context.Context, stop context.CancelFunc) {

	keyChan := make(chan int, numWorkers)
	errChan := make(chan error, numWorkers)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				key := <-keyChan
				res, err, stCode := cl.GetComics(key)
				if err != nil {

					if stCode >= 500 || (stCode >= 300 && stCode < 400) {
						errChan <- nil
					}
					errChan <- err
					return

				} else if key >= numIter {
					errChan <- errors.New("very long base")
					return
				} else {
					mu.Lock()
					errChan <- nil
				}

				(*data)[fmt.Sprintf("%d", key)] = (*res)[key]
				mu.Unlock()

			}

		}()
	}
	key := 1
	for i := 0; i < numWorkers; i++ {
		fmt.Println(i)
		select {
		case <-ctx.Done():
			stop()
		default:
		}
		if _, ok := (*data)[fmt.Sprintf("%d", key)]; ok {
			key++
			i -= 1
			continue
		}
		// Пропускаем это значение, так как ключ уже существует

		keyChan <- key
		key++
	}
	for {
		fmt.Println("key ", key)
		select {
		case <-ctx.Done():
			stop()
		default:
		}
		if err := <-errChan; err != nil {
			break
		}
		for {
			if _, ok := (*data)[fmt.Sprintf("%d", key)]; !ok {
				break
			}
			key++
		}

		keyChan <- key
		key++
	}
	// Ожидаем завершения всех воркеров
	wg.Wait()
}

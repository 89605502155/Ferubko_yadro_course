package worker

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"xkcd/pkg/xkcd"
)

func WorkerPool(cl *xkcd.Client, numIter int, numWorkers int, data *map[string]xkcd.ComicsInfo,
	ctx context.Context, stop context.CancelFunc, exitChan chan bool, isWrite chan bool) {

	keyChan := make(chan int, numWorkers)
	errChan := make(chan error, numWorkers)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				// mu.Lock()
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

				(*data)[fmt.Sprintf("%d", key)] = res[key]
				mu.Unlock()

			}

		}()
	}

	key := 1
	for i := 0; i < numWorkers; i++ {
		select {
		case <-ctx.Done():
			// wg.Wait()
			exitChan <- true
			if <-isWrite {
				stop()
				fmt.Println("Betta")
				wg.Wait()
			}
		case <-exitChan:
			fmt.Println("Betta")
			stop()
			wg.Wait()
			return
		default:
		}
		if _, ok := (*data)[fmt.Sprintf("%d", key)]; ok {
			key++
			i -= 1
			continue
		}
		// Пропускаем это значение, так как ключ уже существует
		fmt.Println(i)
		keyChan <- key
		key++
	}
	for {

		select {
		case <-ctx.Done():
			fmt.Println("Alpha")
			// wg.Wait()
			// fmt.Println("Betta")

			exitChan <- true
			fmt.Println("Gamma")
			if <-isWrite {
				fmt.Println("Tissafern", key)
				stop()
				fmt.Println("Betta", key)
				wg.Wait()
			}
			fmt.Println("Rockrua")

		case <-exitChan:
			fmt.Println("Betta")
			stop()
			wg.Wait()
			return
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
		fmt.Println("key ", key)
		keyChan <- key
		key++
	}
	// Ожидаем завершения всех воркеров
	wg.Wait()
}

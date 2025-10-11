package main

import (
	"fmt"

	"sync"
	"time"
)

const (
	workers    = 8 //колво запускаемых горутин
	iterations = 300000 //столько раз каждая горутина увеличивает счетчик
)

func run(useMutex bool) (actual int64, elapsed time.Duration) {
	var (
		n  int64
		mu sync.Mutex
		wg sync.WaitGroup
	)

	start := time.Now()
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				if useMutex {
					mu.Lock()
					n++
					mu.Unlock()
				} else {
					n++
				}
			}
		}()
	}
	wg.Wait()
	return n, time.Since(start)
}

func main() {
	expected := int64(workers * iterations)

	actual1, t1 := run(false)
	fmt.Printf("БЕЗ мьютекса: факт=%d, ожидаемо=%d, время=%v\n", actual1, expected, t1)

	actual2, t2 := run(true)
	fmt.Printf("С мьютексом:   факт=%d, ожидаемо=%d, время=%v\n", actual2, expected, t2)
}
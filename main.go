package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("Hello")
	type workerResult struct {
		worker int
		result int
	}
	queue := make(chan workerResult)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for i := 0; i < 20; i++ {
			queue <- workerResult{1, i}
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 20; i++ {
			queue <- workerResult{2, i}
		}
		wg.Done()
	}()
	fmt.Println("Waiting")
	go func() {
		wg.Wait()
		close(queue)
	}()

	for i := range queue {
		//fmt.Println("Done:", done)
		fmt.Printf("Worker: %d, result: %d\n", i.worker, i.result)
	}
	fmt.Println("Done")
}

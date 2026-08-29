package main

import (
	"fmt"
	"sync"
	"time"
)

func wg_worker(id int) {
	fmt.Println("Worker ", id, " Starting")
	time.Sleep(time.Second)
	fmt.Println("Worker ", id, "Finished")
}

func main() {
	start := time.Now()

	// This WaitGroup is used to wait for all the goroutines launched here to finish.
	// Note: if a WaitGroup is explicitly passed into functions, it should be done by pointer.
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Go(func() {
			wg_worker(i)
		})
	}

	// Block until all the goroutines started by wg are done. A goroutine is done when the function it invokes returns
	wg.Wait()

	fmt.Println("Execution time:", time.Since(start))

}

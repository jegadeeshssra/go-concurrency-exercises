package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	start := time.Now()

	//  atomic integer type
	var ops atomic.Uint32
	var wg sync.WaitGroup

	for range 50 {
		wg.Go(func() {
			for range 1000 {
				ops.Add(1)
			}
		})
	}

	wg.Wait()

	// Here no goroutines are writing to ‘ops’, but using Load it’s safe to atomically read a value even while other goroutines are (atomically) updating it
	fmt.Println("Final Value - ", ops.Load())

	fmt.Println("Execution time:", time.Since(start))
}

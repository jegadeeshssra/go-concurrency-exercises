package main

import (
	"fmt"
	"sync"
	"time"
)

func rangeGen(start, stop int) chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := start; i < stop; i++ {
			time.Sleep(50 * time.Millisecond)
			out <- i
		}
	}()
	return out
}

func merge(int1 chan int, int2 chan int) chan int {
	combined := make(chan int)

	var wg sync.WaitGroup

	wg.Go(func() {
		for val1 := range int1 {
			combined <- val1
		}
	})

	wg.Go(func() {
		for val2 := range int2 {
			combined <- val2
		}
	})

	// U have to return the combined chan immediatly, so that chan can be ready to send and receive
	go func() {
		wg.Wait()
		close(combined)
	}()
	return combined
}

func main() {

	int1 := rangeGen(10, 14)
	int2 := rangeGen(14, 18)

	start := time.Now()
	combined := merge(int1, int2)
	for ele := range combined {
		fmt.Println(ele)
	}
	fmt.Println()
	fmt.Println("Took", time.Since(start))
}

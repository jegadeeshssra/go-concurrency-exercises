package main

import (
	"fmt"
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
	go func() {
		defer close(combined)
		for int1 != nil || int2 != nil {
			select {
			case val1, ok := <-int1:
				if !ok {
					int1 = nil
				} else {
					combined <- val1
				}
			case val2, ok := <-int2:
				if !ok {
					int2 = nil
				} else {
					combined <- val2
				}
			}
		}
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

package main

import (
	"fmt"
	"sync"
)

func generate() chan struct{} {
	out := make(chan struct{})
	go func() {
		for range 1000000 {
			out <- struct{}{}
		}
		close(out)
	}()
	return out
}

func main() {
	in := generate()

	var wg sync.WaitGroup
	wg.Add(4)

	var counter int

	count := func() {
		defer wg.Done()

		for range in {
			counter++
		}
	}

	for range 4 {
		go count()
	}

	wg.Wait()

	fmt.Println(counter)
}

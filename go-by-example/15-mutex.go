package main

import (
	"fmt"
	"sync"
	"time"
)

type Container struct {
	//Note that mutexes must not be copied, so if this struct is passed around, it should be done by pointer.
	mu      sync.Mutex
	counter map[string]int
}

func (c *Container) inc(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counter[name]++
}

func main() {
	start := time.Now()
	c := Container{
		counter: map[string]int{"a": 0, "b": 0},
	}

	var wg sync.WaitGroup

	increment := func(name string, num int) {
		for range num {
			c.inc(name)
		}
	}

	wg.Go(func() {
		increment("a", 10000)
	})
	wg.Go(func() {
		increment("b", 20000)
	})
	wg.Go(func() {
		increment("a", 30000)
	})

	wg.Wait()
	fmt.Println(c.counter)
	fmt.Println("Execution Time - ", time.Since(start))
}

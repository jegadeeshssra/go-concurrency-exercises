package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	// internal fields
	counter map[string]int
	mu      sync.Mutex
}

func (c *Counter) Add(key string) {
	c.mu.Lock()
	c.counter[key]++
	c.mu.Unlock()
}

func (c *Counter) Get(key string) int {
	return c.counter[key]
}

func main() {
	c := &Counter{counter: make(map[string]int)}
	c.counter["requests"] = 0

	var wg sync.WaitGroup
	increment := func() {
		defer wg.Done()
		for range 10000 {
			c.Add("requests")
		}
	}

	wg.Add(10)
	for range 10 {
		go increment()
	}

	wg.Wait()
	fmt.Println("Counter - ", c.Get("requests"))

}

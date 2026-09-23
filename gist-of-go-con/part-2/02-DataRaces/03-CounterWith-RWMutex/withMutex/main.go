package main

import (
	"fmt"
	"sync"
	"time"
)

type Counter struct {
	// internal fields
	counter map[string]int
	mu      sync.Mutex
}

func (c *Counter) Add(key string) {
	c.mu.Lock()
	c.counter[key]++
	time.Sleep(time.Millisecond)
	c.mu.Unlock()

}

func (c *Counter) Get(key string) int {
	c.mu.Lock()
	num := c.counter[key]
	time.Sleep(time.Millisecond)
	c.mu.Unlock()
	return num
}

func main() {
	c := &Counter{counter: make(map[string]int)}
	c.counter["requests"] = 0

	var wg sync.WaitGroup
	wg.Add(5)

	write := func() {
		defer wg.Done()
		for range 100 {
			c.Add("requests")

		}
	}
	read := func() {
		defer wg.Done()
		for range 100 {
			c.Get("requests")
		}
	}

	start := time.Now()
	go write()
	go read()
	go read()
	go read()
	go read()

	wg.Wait()
	fmt.Println("Time Elapsed -", time.Since(start))

}

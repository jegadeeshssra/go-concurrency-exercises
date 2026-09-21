package main

import (
	"fmt"
	"sync"
	"time"
)

type ConcurrentGroup struct {
	// internal fields
	functions []func(int)
	wg        sync.WaitGroup
}

func NewConcurrentGroup() *ConcurrentGroup {
	return &ConcurrentGroup{}
}

func (c *ConcurrentGroup) Add(fns ...func(int)) {
	c.functions = append(c.functions, fns...)
}

func (c *ConcurrentGroup) Run() {
	c.wg.Add(len(c.functions))
	for i, fn := range c.functions {
		go func() {
			defer c.wg.Done()
			fn(i + 1)
		}()
	}
	fmt.Println("Waiting for the tasks to get finished")
	c.wg.Wait()
	fmt.Println("Finished all the work")
}

func work(id int) {
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("work %d done\n", id)
}

func main() {
	group := NewConcurrentGroup()
	group.Add(work, work, work, work, work)
	start := time.Now()
	group.Run()
	fmt.Println("Time Elapsed -", time.Since(start))
}

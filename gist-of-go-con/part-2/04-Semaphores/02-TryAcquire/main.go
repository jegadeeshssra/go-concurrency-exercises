package main

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore chan struct{}

func NewSemaphore(limit int) Semaphore {
	if limit <= 0 {
		panic("expected limit > 0")
	}
	return make(Semaphore, limit)
}

func (s Semaphore) Acquire() {
	s <- struct{}{}
}

func (s Semaphore) TryAcquire() bool {
	select {
	case s <- struct{}{}:
		return true
	default:
		return false
	}
}

func (s Semaphore) Release() {
	// select for Non-blocking the program and idempotency
	select {
	case <-s:
	default:
		panic("released more than acquired")
	}
}

func Work(wg *sync.WaitGroup, sema Semaphore) {
	defer func() {
		sema.Release()
	}()
	defer wg.Done()
	time.Sleep(100 * time.Millisecond)
}

func main() {
	maxConn := 5
	maxCalls := 15
	sema := NewSemaphore(maxConn)

	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(maxCalls)

	for range maxCalls {
		if sema.TryAcquire() {
			fmt.Println("Goroutine is available")
			go Work(&wg, sema)
		} else {
			fmt.Println("Goroutines are busy")
			wg.Done()
		}
	}
	wg.Wait()
	fmt.Println("Time Elapsed -", time.Since(start))

}

package main

import (
	"fmt"
	"sync"
	"time"
)

type Service struct {
	lock sync.Mutex
}

func (s *Service) Call(wg *sync.WaitGroup) bool {
	defer wg.Done()
	if !s.lock.TryLock() {
		fmt.Println("Busy")
		return false
	}
	defer s.lock.Unlock()
	time.Sleep(100 * time.Microsecond)
	fmt.Println("Finished")
	return true
}

func main() {
	var wg sync.WaitGroup
	start := time.Now()
	service := &Service{}

	wg.Add(10)
	for range 10 {
		go service.Call(&wg)
	}
	wg.Wait()
	fmt.Println("Time Elapsed - ", time.Since(start))
}

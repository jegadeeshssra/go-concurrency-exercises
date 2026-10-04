package main

import (
	"fmt"
	"sync"
	"time"
)

type Rendezvous struct {
	// your fields
	counter int
	mu      sync.Mutex
	ch      chan struct{}
}

func NewRendezvous() *Rendezvous {
	return &Rendezvous{counter: 0, mu: sync.Mutex{}, ch: make(chan struct{})}
}

func (r *Rendezvous) Add(count int) {
	for range count {
		r.mu.Lock()
		r.counter++
		r.mu.Unlock()
	}
}

func (r *Rendezvous) Wait() {
	r.mu.Lock()
	r.counter--
	if r.counter == 0 {
		close(r.ch)
		r.mu.Unlock()
		return
	}
	r.mu.Unlock()
	<-r.ch
}

func worker1(r *Rendezvous, wg *sync.WaitGroup) {
	defer wg.Done()
	// phase 1
	time.Sleep(10 * time.Millisecond)
	fmt.Println("Worker 1 Phase-1 Completed ")

	r.Wait()

	// final work
	time.Sleep(30 * time.Millisecond)
	fmt.Println("Worker 1 work Completed ")
}

func worker2(r *Rendezvous, wg *sync.WaitGroup) {
	defer wg.Done()
	// phase 1
	time.Sleep(20 * time.Millisecond)
	fmt.Println("Worker 2 Phase-1 Completed ")

	r.Wait()

	// final work
	time.Sleep(20 * time.Millisecond)
	fmt.Println("Worker 2 work Completed ")
}

func worker3(r *Rendezvous, wg *sync.WaitGroup) {
	defer wg.Done()
	// phase 1
	time.Sleep(30 * time.Millisecond)
	fmt.Println("Worker 3 Phase-1 Completed ")

	r.Wait()

	// final work
	time.Sleep(10 * time.Millisecond)
	fmt.Println("Worker 3 work Completed ")
}

func main() {

	r := NewRendezvous()
	r.Add(3)

	var wg sync.WaitGroup
	wg.Add(3)
	go worker1(r, &wg)
	go worker2(r, &wg)
	go worker3(r, &wg)
	wg.Wait()

}

package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(100 * time.Millisecond)
	fmt.Println("worker finished")
}

func waiter(wg *sync.WaitGroup, num int) {
	wg.Wait()
	time.Sleep(10 * time.Millisecond)
	fmt.Println("Waiter ", num, " Finished")

}

func main() {
	var wg sync.WaitGroup
	go worker(&wg)

	for i := range 4 {
		go waiter(&wg, i+1)
	}

	wg.Wait()
	time.Sleep(50 * time.Millisecond)
	fmt.Println("All waiters finished their work")
}

package main

import (
	"fmt"
	"sync"
	"time"
)

func normalWork(id int) func() {
	return func() {
		time.Sleep(50 * time.Millisecond)
		fmt.Printf("work %d completed\n", id)
	}
}

func failingWork(id int) func() {
	return func() {
		time.Sleep(50 * time.Millisecond)
		panic(fmt.Sprintf("work %d failed", id))
	}
}

func runSafe(funcs ...func()) bool {
	var panicked bool = false
	var wg sync.WaitGroup
	wg.Add(len(funcs))

	catchPanic := func() {
		err := recover()
		if err != nil {
			panicked = true
		}
	}

	for _, fn := range funcs {
		go func() {
			defer catchPanic()
			defer wg.Done()
			fn()

		}()
	}
	wg.Wait()
	return panicked
}

func main() {
	err := runSafe(
		normalWork(1),
		failingWork(2),
		normalWork(3),
		failingWork(4),
		normalWork(5),
	)
	if err {
		fmt.Println("Goroutine got pannicked")
	} else {
		fmt.Println("All goroutines ran successfully")
	}
}

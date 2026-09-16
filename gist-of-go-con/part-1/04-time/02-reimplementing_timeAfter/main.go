package main

import (
	"fmt"
	"math/rand"
	"time"
)

func work() int {
	if rand.Intn(10) < 8 {
		time.Sleep(10 * time.Millisecond)
	} else {
		time.Sleep(200 * time.Millisecond)
	}

	return 42
}

func after(d time.Duration) <-chan time.Time {
	// implement
	out := make(chan time.Time)
	go func() {
		time.Sleep(d)
		out <- time.Now()
	}()
	return out
}

func WithTimeout(timeout time.Duration, fn func() int) (int, error) {

	done := make(chan int)

	go func() {
		val := fn()
		fmt.Println("Output -", val)
		done <- val
	}()

	select {
	case val := <-done:
		return val, nil
	case <-(after(timeout)):
		return 0, fmt.Errorf("Error - Took more than %v", timeout)
	}

}

func main() {
	timeout := 50 * time.Millisecond
	for range 10 {
		start := time.Now()
		if val, err := WithTimeout(timeout, work); err != nil {
			fmt.Printf("Took longer than %v. Error: %v\n", time.Since(start), err)
		} else {
			fmt.Printf("Took %v. Result: %v\n", time.Since(start), val)
		}
	}
}

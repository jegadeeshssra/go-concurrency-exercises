// - Goal

// Implement throttling and backpressure using a bounded level of concurrency, and understand the difference between:

// blocking until capacity becomes available
// immediately rejecting work when all handlers are busy
// using a buffered channel as a semaphore
// using select with default for non-blocking behavior

// Expected behavior: With a concurrency limit of 2, only requests that can acquire a slot immediately should be accepted; requests arriving while both slots are occupied should receive a "busy" error.

package main

import (
	"fmt"
	"time"
)

func processImage(id int) {
	fmt.Printf("processing image %d\n", id)
	time.Sleep(1 * time.Millisecond)
	fmt.Printf("finished image %d\n", id)
}

func throttle(n int) {
	sema := make(chan struct{}, 2)
	for i := range n {
		select {
		case sema <- struct{}{}:
			fmt.Println("Image ID -", i+1, " was accepted")
			go func() {
				processImage(i + 1)
				<-sema
			}()
		default:
			fmt.Println("Image ID -", i+1, " was REJECTED")
			time.Sleep(1 * time.Millisecond)
		}
	}
	sema <- struct{}{}
	sema <- struct{}{}
}

func main() {
	throttle(25)
}

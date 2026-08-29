package main

import (
	"fmt"
	"time"
)

func main() {

	// Part 1 - Basic rate limiting where 5 requests are handled every 200ms
	requests := make(chan int, 5)
	for i := range 5 {
		requests <- i
	}
	close(requests)

	// This limiter channel will receive a value every 200 milliseconds.
	// This is the regulator in our rate limiting scheme.
	limiter := time.Tick(200 * time.Millisecond)

	// By blocking on a receive from the limiter channel before serving each request,
	// we limit ourselves to 1 request every 200 milliseconds.
	for req := range requests { // loop through the 5 requests
		<-limiter // produces an event every 200ms
		fmt.Println("Incoming requests - ", req, time.Now())
	}

	// PART 2 - Handling bursts

	burstyLimiter := make(chan time.Time, 3)
	// Fills the limiter ,so that it handles the 3 requests simultaneouly
	for range 3 {
		burstyLimiter <- time.Now()
	}

	// To add value to the limiter every 200ms
	go func() {
		ticker := time.Tick(200 * time.Millisecond)
		// The ticker is not sending values directly.
		// It is sending values into ticker.C.
		// Your goroutine receives those values one at a time, and then tries to send each one into burstyLimiter
		for i := range ticker {
			burstyLimiter <- i
		}
	}()

	burstyRequests := make(chan int, 5)
	for i := range 5 {
		burstyRequests <- i
	}
	close(burstyRequests)

	for req := range burstyRequests {
		<-burstyLimiter
		fmt.Println("Request", req, time.Now())
	}

}

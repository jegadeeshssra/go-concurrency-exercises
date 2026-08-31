// - The producer should submit all three orders immediately, while the consumer starts processing only after a short delay.
// Requirements
// Use a buffered channel.
// The producer should be able to send all three orders without waiting for the consumer to start.
// The channel should have exactly enough capacity to hold the three orders.
// Print the channel's len() and cap() at an appropriate point to demonstrate that the buffer is being used.
// Use a sync.WaitGroup to ensure the program doesn't exit before both goroutines finish.

package main

import (
	"fmt"
	"time"
)

func producer(orders []string, out chan<- string) {
	for _, order := range orders {
		out <- order
	}
	close(out)
}

func consumer(out <-chan string, done chan struct{}) {
	for order := range out {
		fmt.Printf("processing %s\n", order)
		time.Sleep(200 * time.Millisecond)
		fmt.Printf("finished %s\n", order)
	}
	done <- struct{}{}
}

func main() {
	orders := []string{
		"order-101",
		"order-102",
		"order-103",
	}
	out := make(chan string, 3)
	done := make(chan struct{})
	go producer(orders, out)
	go consumer(out, done)
	<-done
	fmt.Println("PROCEESED ALL ORDERS")
}

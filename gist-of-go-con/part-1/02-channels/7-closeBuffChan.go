// The producer should send all numbers, close the channel, and finish before the consumer starts reading.

// Write a program that:

// Creates a buffered channel with capacity 5.
// Starts a producer goroutine.
// The producer sends all five numbers.
// The producer closes the channel immediately after the final send.
// The main goroutine waits for 100ms before beginning to receive.
// The main goroutine uses for range to read the channel.
// Every number must still be received even though the channel was already closed.

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

func main() {
	orders := []string{
		"order-101",
		"order-102",
		"order-103",
	}
	out := make(chan string, 3)
	go producer(orders, out)

	time.Sleep(time.Second)

	for order := range out {
		fmt.Printf("processing %s\n", order)
		time.Sleep(200 * time.Millisecond)
		fmt.Printf("finished %s\n", order)
	}

	fmt.Println("PROCEESED ALL ORDERS")
}

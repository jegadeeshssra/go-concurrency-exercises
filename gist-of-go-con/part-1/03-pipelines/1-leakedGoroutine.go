// FIX - Leaked Goroutine
// Write a function rangeGen(start, stop int) <-chan int that generates numbers from start up to, but not including, stop.

// Then write main() so that it:
// Creates a generator for 1 through 100.
// Prints generated numbers.
// Stops consuming the generator when it reaches 5.
// The program should terminate without leaving the generator goroutine blocked.

// "cancel" channel is used to signal the goroutine to stop and exit

package main

import "fmt"

func rangeGen(cancel <-chan struct{}, start, stop int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for start := range stop {
			select {
			case out <- start:
			case <-cancel:
				return
			}
		}
		fmt.Println("Sent ALL data")
	}()
	return out
}

func main() {

	cancel := make(chan struct{})
	generator := rangeGen(cancel, 1, 11)

	for i := range generator {
		fmt.Println("->", i)
		if i == 5 {
			cancel <- struct{}{}
			return
		}
	}
}

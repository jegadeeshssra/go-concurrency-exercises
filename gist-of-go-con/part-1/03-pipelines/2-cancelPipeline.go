// Create a worker that continuously receives numbers from an input channel, doubles them, and sends the results to an output channel.
// Write a program where the consumer reads only a few results and then cancels the worker.

// Requirements

// The worker must:
// Listen for input values.
// Process each value.
// Send the result to the output channel.
// Respond to cancellation while waiting for input.
// Respond to cancellation while trying to send the result.
// Close the output channel when it exits.

// Important: Do not assume that putting a select around the input receive automatically makes the entire worker cancellable.
// The worker must not remain stuck trying to send a result after cancellation has been requested.
// Goal: Understand that cancellation must protect channel operations that can block, not just the outer loop.

package main

import (
	"fmt"
	"time"
)

func gen(cancel <-chan struct{}, start int, stop int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for start := range stop {
			select {
			case out <- start:
				fmt.Println("Generated -", start)
			case <-cancel:
				fmt.Println("Received Cancel in GENERATOR")
				return
			}
		}
	}()
	return out
}

func multiply(val int) int {
	time.Sleep(100 * time.Millisecond)
	return val * val
}

func process(cancel <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for val := range in {
			select {
			case out <- multiply(val):
				fmt.Println("Processed -", val*val)
			case <-cancel:
				fmt.Println("Recieved Cancel in PROCESSOR")
				return
			}
		}
	}()
	return out
}

func main() {
	cancel := make(chan struct{})
	generator := gen(cancel, 1, 11)
	result := process(cancel, generator)
	for val := range result {
		fmt.Println("->", val)
		if val == 36 {
			close(cancel)
			fmt.Println("Send CANCEL Signal")
			time.Sleep(time.Second)
			return
		}
	}
}

//OUTPUT
// Generated - 0
// Processed - 0
// Generated - 1
// -> 0
// Processed - 1
// Generated - 2
// -> 1
// Processed - 4
// Generated - 3
// -> 4
// Processed - 9
// Generated - 4
// -> 9
// Processed - 16
// Generated - 5
// -> 16
// Processed - 25
// Generated - 6
// -> 25
// Processed - 36
// Generated - 7
// -> 36

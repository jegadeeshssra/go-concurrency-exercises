// - Suppose the pipeline occasionally generates errors, but the error collector may be slower than the pipeline.

// Create an error channel with capacity 3.
// When a pipeline stage encounters an error:

// Send the error if the error channel has space.
// If the channel is full, do not block the normal pipeline processing.
// Instead, discard the error.

// GOAL - Understand how a buffered channel combined with select can prevent auxiliary error handling from blocking the primary pipeline.

package main

import "fmt"

type ErrResult struct {
	val int
	err error
}

// Centralized error collector to collect errors from all piplines
func collectErrors(errc chan ErrResult) chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := range errc {
			fmt.Println("Val and err -", i.val, " ", i.err)
		}
	}()
	return done
}

func generate(stop int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := range stop {
			out <- i + 1
		}
	}()
	return out
}

func fetchScore(val int) (int, error) {
	if val%3 == 0 {
		return val, fmt.Errorf("It is divisble by 3")
	}
	return val * val, nil
}

func calc(in <-chan int, errc chan ErrResult) chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := range in {
			val, err := fetchScore(i)
			if err == nil {
				out <- val
			} else {
				select {
				case errc <- ErrResult{val, err}:
				default:
					fmt.Println("Error Dropped - ", val, " , ", err)
				}
			}
		}
	}()
	return out
}

func main() {
	errc := make(chan ErrResult, 1)
	done := collectErrors(errc)
	gen := generate(6)
	results := calc(gen, errc)
	for res := range results {
		fmt.Println("Result Val -", res)
	}
	close(errc)
	<-done
}

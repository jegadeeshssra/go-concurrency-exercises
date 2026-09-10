// - Build a pipeline where multiple stages can report errors to a single error collector.
// Requirements
// Use one shared error channel.
// Pipeline stages must only send errors to it.
// The main goroutine must close the error channel after all pipeline stages have finished sending errors.
// Wait for the error collector to finish before exiting.
// Do not close errc from inside the pipeline stages.

// Goal: Separate normal pipeline results from centralized error reporting.

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
				errc <- ErrResult{val, err}
			}
		}
	}()
	return out
}

func main() {
	errc := make(chan ErrResult)
	done := collectErrors(errc)
	gen := generate(6)
	results := calc(gen, errc)
	for res := range results {
		fmt.Println("Result Val -", res)
	}
	close(errc)
	<-done
}

// This program is not to block the program for one error from fetchScore
// Keep successful results and their corresponding errors together so that one failed operation doesn't terminate the pipeline.

package main

import "fmt"

// To know which answer caused the error
type Result struct {
	ans int
	err error
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
	if val%5 == 0 {
		return 0, fmt.Errorf("It is divisble by 3")
	}
	return val * val, nil
}

func cal(in <-chan int) chan Result {
	out := make(chan Result)
	go func() {
		defer close(out)
		for i := range in {
			val, err := fetchScore(i)
			out <- Result{val, err}
		}
	}()
	return out
}

func main() {
	gen := generate(5)
	result := cal(gen)
	for res := range result {
		if res.err == nil {
			fmt.Println("Value - ", res.ans)
		} else {
			fmt.Println("error:", res.err)
		}
	}
}

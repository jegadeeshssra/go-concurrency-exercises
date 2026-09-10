// Goal: Generalize the composite-result pattern so that every pipeline stage doesn't need its own result structure.

package main

import "fmt"

type Res[T any] struct {
	val T
	err error
}

func (r *Res[T]) OK() bool {
	return r.err == nil
}
func (r *Res[T]) Val() T {
	return r.val
}
func (r *Res[T]) Err() error {
	return r.err
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

func cal(in <-chan int) chan Res[int] {
	out := make(chan Res[int])
	go func() {
		defer close(out)
		for i := range in {
			val, err := fetchScore(i)
			out <- Res[int]{val, err}
		}
	}()
	return out
}

func main() {
	gen := generate(5)
	results := cal(gen)
	for res := range results {
		if res.OK() {
			fmt.Printf("✓ %v\n", res.Val())
		} else {
			fmt.Printf("✗ %v\n", res.Err())
		}
	}
}

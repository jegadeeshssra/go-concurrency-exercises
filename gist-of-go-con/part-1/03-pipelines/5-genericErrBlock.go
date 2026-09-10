package main

import "fmt"

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

func calculate(in <-chan int) (chan int, chan error) {
	out := make(chan int)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		for i := range in {
			val, err := fetchScore(i)
			if err != nil {
				errc <- err
				return
			}
			out <- val
		}
		errc <- nil
	}()
	return out, errc
}

func main() {
	gen := generate(5)
	result, errc := calculate(gen)
	for val := range result {
		fmt.Println("Value - ", val)
	}
	if err := <-errc; err != nil {
		fmt.Println("error:", err)
	}
}

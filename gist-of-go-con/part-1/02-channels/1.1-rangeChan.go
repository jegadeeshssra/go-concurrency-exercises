package main

import "fmt"

func main() {
	numbers := make(chan int)

	go func() {
		for i := 1; i <= 10; i++ {
			numbers <- i
		}
		close(numbers)
	}()

	for {
		for n := range numbers {
			if n%2 == 0 {
				fmt.Println(n)
			}
		}
		break
	}
}

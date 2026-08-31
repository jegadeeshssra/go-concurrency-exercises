// Write a program where a worker calculates numbers and sends the results through one channel while also notifying the main program when it has finished.

package main

import "fmt"

func worker(out chan<- int, done chan<- struct{}) {
	for i := range 5 {
		out <- i
	}
	close(out)
	done <- struct{}{}
}

func main() {

	out := make(chan int)
	done := make(chan struct{})

	go worker(out, done)

	go func() {
		for num := range out {
			fmt.Println("Num -", num)
		}
	}()

	<-done
	fmt.Println("Program is finished")

}

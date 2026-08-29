package main

import "fmt"

func main() {
	jobs := make(chan int, 5)
	done := make(chan bool)

	go func() {
		for {
			j, ok := <-jobs
			if ok {
				fmt.Println("Received job - ", j)
			} else {
				fmt.Println("All jobs Received")
				done <- true
				return
			}
		}
	}()

	for i := range 3 {
		jobs <- i
		fmt.Println("Send the job - ", i)
	}

	close(jobs)

	<-done

	_, ok := <-jobs
	fmt.Println("Is any job available -", ok)
}

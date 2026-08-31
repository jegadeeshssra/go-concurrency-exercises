// Start one goroutine for each file.
// Do not use time.Sleep() in main() to wait.
// Do not use sync.WaitGroup.
// Create a done channel.
// Each worker must send a completion signal through the channel after finishing.
// main() must wait for all workers to report completion

package main

import (
	"fmt"
	"time"
)

func main() {

	files := []string{
		"users.csv",
		"orders.csv",
		"payments.csv",
		"products.csv",
	}

	done := make(chan struct{})
	for i, file := range files {
		go processFile(i+1, file, done)
	}
	for range len(files) {
		<-done
	}
	fmt.Println("All goroutines are finished")
}

func processFile(id int, name string, done chan<- struct{}) {
	fmt.Printf("Worker #%d processing %s\n", id, name)
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Worker #%d finished %s\n", id, name)
	done <- struct{}{}
}

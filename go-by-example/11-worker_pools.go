package main

import (
	"fmt"
	"time"
)

// Here’s the worker, of which we’ll run several concurrent instances.
// These workers will receive work on the jobs channel and send the corresponding results on results.
func worker(id int, jobs <-chan int, results chan<- int) {
	for i := range jobs {
		fmt.Println("Worker", id, "Started the job", i)
		time.Sleep(time.Second)
		fmt.Println("Worker", id, "Finished the job", i)
		results <- (i * i)
	}
}

func main() {
	const numJobs int = 5

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// This starts up 3 workers, initially blocked because there are no jobs yet.
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	for i := 1; i <= numJobs; i++ {
		fmt.Println("Results -", <-results)
	}

}

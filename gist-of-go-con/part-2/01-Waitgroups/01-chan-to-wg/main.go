package main

import (
	"fmt"
	"sync"
	"time"
)

// Only ptr of wg should be received.
func runTask(wg *sync.WaitGroup, id int) {
	time.Sleep(time.Duration(id*50) * time.Millisecond)
	fmt.Printf("task %d done\n", id)
	wg.Done()
}

func main() {

	var wg sync.WaitGroup

	for i := range 3 {
		wg.Add(1)
		go runTask(&wg, i+1)
		// wg.Go(runTask(i))
	}

	wg.Wait()
	fmt.Println("done")

}

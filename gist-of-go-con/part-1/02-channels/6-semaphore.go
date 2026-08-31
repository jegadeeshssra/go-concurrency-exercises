// Write a program that starts a goroutine for each image but allows at most 3 image-processing goroutines to run at the same time.
// Requirements
// Use a buffered channel as a semaphore.
// The semaphore capacity must be 3.
// A goroutine must acquire a token before processing.
// It must return the token after processing.
// No more than three images may be processed concurrently.
// Wait until all images have finished before printing:

package main

import (
	"fmt"
	"time"
)

func process(id int, image string, sema chan struct{}) {
	fmt.Println("Worker -", id, " Processing -", image)
	time.Sleep(200 * time.Millisecond)
	fmt.Println("Worker -", id, " Finished -", image)
	<-sema
}

func main() {
	images := []string{
		"img1.jpg",
		"img2.jpg",
		"img3.jpg",
		"img4.jpg",
		"img5.jpg",
		"img6.jpg",
		"img7.jpg",
		"img8.jpg",
		"img9.jpg",
		"img10.jpg",
	}
	sema := make(chan struct{}, 2)

	for i, image := range images {
		sema <- struct{}{}
		go process(i+1, image, sema)
	}

	// Till wait until all the buffer is free
	for range cap(sema) {
		sema <- struct{}{}
	}
	fmt.Println("done")

}

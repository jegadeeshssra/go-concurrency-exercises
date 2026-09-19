package main

import (
	"context"
	"fmt"
	"time"
)

func acquireResource() {
	fmt.Println("resource acquired")
}

func releaseResource() {
	fmt.Println("resource released")
}

func work() chan struct{} {
	acquireResource()
	done := make(chan struct{}, 1) // to prevent leaked goroutine
	go func() {
		time.Sleep(100 * time.Millisecond)
		done <- struct{}{}
		fmt.Println("Leaked goroutine got exited")
	}()

	return done
}

func worker(ctx context.Context) {
	// implement
	select {
	case <-work():
		fmt.Println("Finished the Work")
		return
	case <-ctx.Done():
		fmt.Println("Time Interrupt")
		return
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 99*time.Millisecond)
	defer cancel()

	releaseCancel := context.AfterFunc(ctx, releaseResource)
	if releaseCancel() {
		fmt.Println("Cancelled the cleanup Job")
	}
	worker(ctx)

	// Too late to cancel Cleanup function
	if releaseCancel() {
		fmt.Println("Too late to cancel cleanup")
	}

	time.Sleep(time.Second)
}

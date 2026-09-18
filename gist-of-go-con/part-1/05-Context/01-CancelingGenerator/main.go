package main

import (
	"context"
	"fmt"
)

func generate(ctx context.Context, start, stop int) <-chan int {
	// implement
	out := make(chan int)
	go func() {
		for val := start; val <= stop; val++ {
			select {
			case <-ctx.Done():
				close(out)
				return
			case out <- val:
				fmt.Println("Sent Value -", val)
			}
		}
	}()
	return out
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	generator := generate(ctx, 1, 100)

	for val := range generator {
		fmt.Println("Received Val -", val)
		if val == 10 {
			break
		}
	}

	fmt.Println("Received Required Values")

}

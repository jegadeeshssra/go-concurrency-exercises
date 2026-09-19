package main

import (
	"context"
	"fmt"
	"time"
)

func produce(ctx context.Context, start int, stop int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := start; i <= stop; i++ {
			select {
			case out <- i:
			case <-ctx.Done():
				fmt.Println("Producer Stopped")
				return
			}
		}
	}()
	return out
}

func process(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for {
			select {
			case n := <-in:
				time.Sleep(10 * time.Millisecond)
				select {
				case out <- (n * 2):
				case <-ctx.Done():
					fmt.Println("Processor Stopped")
					return
				}
			case <-ctx.Done():
				fmt.Println("Processor Stopped")
				return
			}
		}
	}()
	return out
}
func consume(ctx context.Context, in <-chan int) {
	go func() {
		for {
			select {
			case val := <-in:
				fmt.Println("Final Val -", val)
			case <-ctx.Done():
				fmt.Println("Consumer Stopped")
				return
			}
		}
	}()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	producer := produce(ctx, 1, 100)
	processer := process(ctx, producer)
	consume(ctx, processer)

	time.Sleep(150 * time.Millisecond)
}

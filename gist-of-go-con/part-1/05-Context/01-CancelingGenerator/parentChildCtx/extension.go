package main

import (
	"context"
	"fmt"
	"time"
)

func produce(ctx context.Context, start, stop int) <-chan int {
	// implement
	out := make(chan int)
	go func() {
		for val := start; val <= stop; val++ {
			select {
			case <-ctx.Done():
				close(out)
				fmt.Println("")
				return
			case out <- val:
				//fmt.Println("Sent Value -", val)
			}
		}
	}()
	return out
}

func consume(ctx context.Context, producer <-chan int, src string) {
	for {
		select {
		case val := <-producer:
			fmt.Println("Received ", src, " Val -", val)
		case <-ctx.Done():
			fmt.Println("Ctx Cancelled")
			return
		}
	}
}

func main() {
	ctx := context.Background()

	parentCtx, parentCancel := context.WithCancel(ctx)
	childCtx, childCancel := context.WithCancel(parentCtx)

	defer childCancel()

	parentProduce := produce(parentCtx, 1, 50)
	childProduce := produce(childCtx, 1, 50)

	go consume(parentCtx, parentProduce, "parent")
	go consume(childCtx, childProduce, "child")

	time.Sleep(time.Millisecond)
	parentCancel()
	time.Sleep(5 * time.Millisecond)

	fmt.Println("Received Required Values")

}

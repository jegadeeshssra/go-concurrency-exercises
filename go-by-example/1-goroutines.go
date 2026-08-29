package main

import (
	"fmt"
	"time"
)

func f(from string) {
	for i := range 3 {
		fmt.Printf("\n%s - %d", from, i)
	}
}

func main() {
	f("Blocking Call")

	go f("Concurrent Call")

	go func(val string) {
		for i := range 3 {
			fmt.Printf("\n%s - %d", val, i)
		}
	}("2nd Concurrent Call")

	time.Sleep(time.Second)
	fmt.Println("\nTHE END")
}

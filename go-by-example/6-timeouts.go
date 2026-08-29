package main

import (
	"fmt"
	"time"
)

func main() {

	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	go func() {
		time.Sleep(time.Second * 2)
		ch1 <- "Result 1"
	}()

	select {
	case msg := <-ch1:
		fmt.Println("msg - ", msg)
	case <-time.After(time.Second):
		fmt.Println("Timeout 1")
	}

	go func() {
		time.Sleep(time.Second * 2)
		ch2 <- "Result 2"
	}()

	select {
	case msg := <-ch2:
		fmt.Println("msg - ", msg)
	case <-time.After(time.Second * 3):
		fmt.Println("Timeout 2")
	}
}

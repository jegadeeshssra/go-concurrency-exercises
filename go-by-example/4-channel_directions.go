package main

import "fmt"

func ping(ch1 chan<- string, msg string) {
	ch1 <- msg
}

func pong(ch1 <-chan string, ch2 chan<- string) {
	msg := <-ch1
	ch2 <- msg
}

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go ping(ch1, "Secret-key")
	go pong(ch1, ch2)

	// Deadlock - because synchronous function calling will block until a value is send to the channel
	// pong(ch1, ch2)
	// ping(ch1, "Secret-key")

	fmt.Println("msg -", <-ch2)

}

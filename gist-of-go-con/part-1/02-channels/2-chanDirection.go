// - The program works, but both functions accept a bidirectional channel. This means submit() could accidentally receive from the channel, or receive() could accidentally send to or close it.
// - Change only the function signatures so that each function can perform only the operations it is supposed to perform.

package main

import (
	"fmt"
	"strings"
)

func main() {
	str := "one,two,,four"
	stream := make(chan string)

	go submit(str, stream)
	receive(stream)
}

func submit(str string, stream chan<- string) {
	words := strings.Split(str, ",")
	for _, word := range words {
		stream <- word
	}
	// This operation should not be allowed.
	// 	<-stream

	// The producer is responsible for closing the stream.
	close(stream)
}

func receive(stream <-chan string) {
	for word := range stream {
		if word != "" {
			fmt.Printf("%s ", word)
		}
	}
	// This operation should not be allowed.
	// close(stream)

	fmt.Println()
}

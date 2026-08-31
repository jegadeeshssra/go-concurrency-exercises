// - A program starts a goroutine that reads a list of log messages and sends them through a channel.
// The main goroutine prints only messages containing "ERROR".
// - Modify the program so that the reader knows when the goroutine has finished sending all log
// messages and can stop reading from the channel.

package main

import (
	"fmt"
	"strings"
)

func main() {
	logs := []string{
		"INFO: server started",
		"ERROR: database unavailable",
		"INFO: retrying",
		"ERROR: connection refused",
	}

	messages := make(chan string)

	go func() {
		for _, log := range logs {
			messages <- log
		}
		close(messages)
	}()

	for {
		message, ok := <-messages
		if !ok {
			break
		}

		if strings.HasPrefix(message, "ERROR") {
			fmt.Println(message)
		}
	}
}

package main

import "fmt"

func main() {
	message := make(chan string)
	signals := make(chan int)

	select {
	case msg := <-message:
		fmt.Println("Receive Message", msg)
	default:
		fmt.Println("No Receive Message")
	}

	msg := "secret"
	select {
	case message <- msg:
		fmt.Println("Send Message")
	default:
		fmt.Println("No Send Message")
	}

	select {
	case <-message:
		fmt.Println("Receive Message")
	case <-signals:
		fmt.Println("Receive Signal")
	default:
		fmt.Println("No Activity")
	}
}

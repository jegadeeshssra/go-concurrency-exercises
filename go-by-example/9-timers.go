package main

import (
	"fmt"
	"time"
)

func main() {
	// The <-timer1.C blocks on the timer’s channel C until it sends a value indicating that the timer fired.
	t1 := time.NewTimer(2 * time.Second)
	<-t1.C
	fmt.Println("Timer 1 fired")

	t2 := time.NewTimer(time.Microsecond)
	// done := make(chan bool)
	go func() {
		<-t2.C
		fmt.Println("Timer 2 fired")
	}()

	// If you just wanted to wait, you could have used time.Sleep.
	// One reason a timer may be useful is that you can cancel the timer before it fires.
	stop2 := t2.Stop()
	if stop2 {
		fmt.Println("Timer 2 Stopped from firing")
	}
}

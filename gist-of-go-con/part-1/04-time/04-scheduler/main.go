package main

import (
	"fmt"
	"time"
)

func work(at time.Time, cancel chan struct{}) {
	fmt.Printf("work started at %s\n", at.Format("15:04:05.000"))
	time.Sleep(100 * time.Millisecond)
	//fmt.Printf("work ended at %s\n", at.Format("15:04:05.000"))
	select {
	case <-time.After(100 * time.Millisecond):
		return
	case <-cancel:
		return
	}
}

func main() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	cancel := make(chan struct{})
	go func() {
		for {
			curTime := <-ticker.C
			go work(curTime, cancel)
		}
	}()

	time.Sleep(560 * time.Millisecond)
	cancel <- struct{}{}
	close(cancel)
}

package main

import (
	"fmt"
	"time"
)

type Timer struct {
	done   chan struct{}
	cancel chan struct{}
}

func (t *Timer) Stop() bool {
	fmt.Println("Initiating Cancel operation")
	select {
	case <-t.done:
		return false
	case t.cancel <- struct{}{}:
		close(t.cancel)
		return true
	}
}

func (t *Timer) Start(d time.Duration) {
	select {
	case <-time.After(d):
		t.done <- struct{}{}
		close(t.done)
	case <-t.cancel:
		return
	}
}

func afterFunc(d time.Duration, fn func()) *Timer {

	cancel := make(chan struct{})
	done := make(chan struct{})
	timer := &Timer{done, cancel}

	go timer.Start(d)

	go func(t *Timer) {
		select {
		case <-t.done:
			fn()
			return
		case <-t.cancel:
			return
		}
	}(timer)

	return timer
}

func work() {
	fmt.Println("Work Started")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Work Done")
}

func main() {
	duration := time.Duration(50 * time.Millisecond)
	timer := afterFunc(duration, work)

	time.Sleep(60 * time.Millisecond)
	if !timer.Stop() {
		fmt.Println("The function has already started to execute")
	}
}

package main

import (
	"fmt"
	"math/rand/v2"
	"sync/atomic"
	"time"
)

type readOp struct {
	key      int
	response chan int
}
type writeOp struct {
	key      int
	val      int
	response chan bool
}

func main() {

	// In order to count the no of readOps and writeOps
	var readOps uint32
	var writeOps uint32

	// COMMON channels to listen for reads and writes
	reads := make(chan readOp)
	writes := make(chan writeOp)

	// this func contains a map which can only be modified within this function by listening to the
	// incoming reads and writes requests.
	go func() {
		for {
			// To fill in with default values
			var state = make(map[int]int)
			select {
			case read := <-reads:
				// reads the value of the key from the map
				read.response <- state[read.key]
			case write := <-writes:
				// writes the value to a key in the map
				state[write.key] = write.val
				write.response <- true
			}
		}
	}()

	// 100 goroutines to read random keys of the map of the owning goroutine.
	for range 100 {
		go func() {
			for {
				read := readOp{
					key:      rand.IntN(5),
					response: make(chan int),
				}
				reads <- read
				<-read.response
				// Increments the read Ops count
				atomic.AddUint32(&readOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	for range 100 {
		go func() {
			for {
				write := writeOp{
					key:      rand.IntN(5),
					val:      rand.IntN(100),
					response: make(chan bool),
				}
				writes <- write
				<-write.response
				atomic.AddUint32(&writeOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	time.Sleep(time.Second)

	readOpsFinal := atomic.LoadUint32(&readOps)
	fmt.Println("readOps:", readOpsFinal)
	writeOpsFinal := atomic.LoadUint32(&writeOps)
	fmt.Println("writeOps:", writeOpsFinal)

}

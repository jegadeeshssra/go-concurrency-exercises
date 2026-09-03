package main

import (
	"fmt"
	"sync"
	"time"
)

func rangeGenerator(start int, stop int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := start; i <= stop; i++ {
			out <- i
		}
	}()
	return out
}

func mergeChan(sensors []<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, sensor := range sensors {
		wg.Go(func() {
			for val := range sensor {
				out <- val
			}
		})
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	start := time.Now()
	sensor1 := rangeGenerator(1, 5)
	sensor2 := rangeGenerator(101, 105)
	sensor3 := rangeGenerator(201, 205)
	sensor4 := rangeGenerator(301, 305)

	sensors := []<-chan int{
		sensor1,
		sensor2,
		sensor3,
		sensor4,
	}

	merged := mergeChan(sensors)
	for val := range merged {
		fmt.Println(val)
	}
	fmt.Println("Duration -", time.Since(start))

}

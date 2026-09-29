package main

import (
	"fmt"
	"sync"
)

type Inventory struct {
	stock uint
	mu    sync.Mutex
}

func NewInventory() *Inventory {
	return &Inventory{}
}

func (i *Inventory) Get() uint {
	i.mu.Lock()
	defer i.mu.Unlock()

	return i.stock
}

func (i *Inventory) CompareAndSet(old uint) bool {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.stock != old {
		return false
	}

	i.stock -= 1
	return true
}

func reserve(i *Inventory, wg *sync.WaitGroup) bool {
	defer wg.Done()
	for range 5 {
		old := i.Get()
		if old == 0 {
			fmt.Println("Not Enough Stock")
			return false
		}
		if i.CompareAndSet(old) {
			fmt.Println("Reserved")
			return true
		}
		fmt.Println("Retrying")
	}
	fmt.Println("Limit reached for Retrying")
	return false
}

func main() {
	i := NewInventory()
	i.stock = 10

	var wg sync.WaitGroup
	wg.Add(20)
	for range 20 {
		go reserve(i, &wg)
	}

	wg.Wait()
	fmt.Println("Stock Remaining -", i.stock)
}

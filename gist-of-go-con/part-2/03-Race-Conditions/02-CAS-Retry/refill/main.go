package main

import (
	"fmt"
	"sync"
)

const INVENTORY_LIMIT = 20

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

func (i *Inventory) CompareAndRefill(old uint) bool {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.stock != old {
		return false
	}

	i.stock += 1
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

func refill(i *Inventory, wg *sync.WaitGroup) bool {
	defer wg.Done()
	for range 5 {
		old := i.Get()
		if old >= INVENTORY_LIMIT {
			fmt.Println("Too much Stock")
			return false
		}
		if i.CompareAndRefill(old) {
			fmt.Println("Refilled")
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
	wg.Add(40)
	for range 20 {
		go reserve(i, &wg)
	}
	for range 10 {
		go refill(i, &wg)
	}

	wg.Wait()
	fmt.Println("Stock Remaining -", i.stock)
}

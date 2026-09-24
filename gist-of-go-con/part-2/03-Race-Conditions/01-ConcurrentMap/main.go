package main

import (
	"fmt"
	"maps"
	"sync"
)

type Accounts struct {
	bal map[string]int
	mu  sync.Mutex
}

func NewAccounts(bal map[string]int) *Accounts {
	return &Accounts{bal: maps.Clone(bal)}
}

func (acc *Accounts) Get(name string) int {
	acc.mu.Lock()
	defer acc.mu.Unlock()
	return acc.bal[name]
}

func (acc *Accounts) Print() {
	for key := range acc.bal {
		fmt.Println(key, " - ", acc.bal[key])
	}
}

// Data Race
func (acc *Accounts) Add(name string, amt int) {
	acc.mu.Lock()
	defer acc.mu.Unlock()
	acc.bal[name] += amt
}

// Race Condition
func (acc *Accounts) CompareAndSet(name string, oldAmt int, newAmt int) bool {
	acc.mu.Lock()
	defer acc.mu.Unlock()
	if acc.bal[name] != oldAmt {
		return false
	}
	acc.bal[name] = newAmt
	return true
}

// Race Condition
func transfer(wg *sync.WaitGroup, accounts *Accounts, from string, to string, amount int) bool {
	defer wg.Done()

	accounts.mu.Lock()
	defer accounts.mu.Unlock()

	if accounts.bal[from] < amount {
		return false
	}
	accounts.bal[to] += amount
	accounts.bal[from] -= amount
	return true
}

func main() {
	alice := "Alice"
	bob := "Bob"

	balance := make(map[string]int)
	acc := NewAccounts(balance)
	acc.bal[alice] = 100
	acc.bal[bob] = 150

	var wg sync.WaitGroup
	wg.Add(4)
	go transfer(&wg, acc, alice, bob, 50)
	go transfer(&wg, acc, alice, bob, 150)
	go transfer(&wg, acc, bob, alice, 100)
	go transfer(&wg, acc, bob, alice, 125)

	wg.Wait()

	acc.Print()
	fmt.Println("")
}

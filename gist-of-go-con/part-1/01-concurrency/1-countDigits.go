package main

import (
	"fmt"
	"strings"
	"sync"
	"unicode"
)

type State struct {
	mu    sync.Mutex
	count int
}

func (s *State) inc() {
	s.mu.Lock()
	s.count++
	defer s.mu.Unlock()
}

func countDigits(word string, state *State) {
	for _, char := range word {
		if unicode.IsDigit(char) {
			state.inc()
		}
	}
}

func main() {
	s := &State{
		count: 0,
	}
	sentence := "Gola90ng 1is1 4also5 24used fo4r s8ystem level programming"
	words := strings.Fields(sentence)

	var wg sync.WaitGroup

	for _, word := range words {
		wg.Go(func() {
			countDigits(word, s)
		})
	}

	wg.Wait()
	fmt.Println("Total Count -", s.count)
}

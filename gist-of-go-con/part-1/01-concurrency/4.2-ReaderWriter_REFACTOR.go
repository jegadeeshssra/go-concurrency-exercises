// Reader and worker.
package main

import (
	"fmt"
	"strings"
	"unicode"
)

// counter stores the number of digits in each word.
// The key is the word, and the value is the number of digits.
type counter map[string]int

// pair stores a word and the number of digits in it.
type pair struct {
	word  string
	count int
}

// solution start

func submitWords(next func() string) chan string {
	pending := make(chan string)
	go func() {
		for {
			word := next()
			if word == "" {
				close(pending)
				break
			}
			pending <- word
		}
	}()
	return pending
}

func countWords(pending <-chan string) chan pair {
	counted := make(chan pair)
	go func() {
		for {
			word, ok := <-pending
			if !ok {
				close(counted)
				break
			}
			count := countDigits(word)
			counted <- pair{word, count}
		}
	}()
	return counted
}

func fillStats(counted chan pair) counter {
	stats := counter{}
	for {
		p, ok := <-counted
		if !ok {
			break
		}
		stats[p.word] = p.count
	}
	return stats
}

// countDigitsInWords counts the number of digits in words,
// fetching the next word with next().
func countDigitsInWords(next func() string) counter {

	pending := submitWords(next)
	counted := countWords(pending)
	return fillStats(counted)

	// As a result, stats should contain words
	// and the number of digits in each.

}

// solution end

// countDigits returns the number of digits in a string.
func countDigits(str string) int {
	count := 0
	for _, char := range str {
		if unicode.IsDigit(char) {
			count++
		}
	}
	return count
}

// printStats prints the number of digits in words.
func printStats(stats counter) {
	for word, count := range stats {
		fmt.Printf("%s: %d\n", word, count)
	}
}

// wordGenerator returns a generator that yields words from a phrase.
func wordGenerator(phrase string) func() string {
	words := strings.Fields(phrase)
	idx := 0
	return func() string {
		if idx == len(words) {
			return ""
		}
		word := words[idx]
		idx++
		return word
	}
}

func main() {
	phrase := "0ne 1wo thr33 4068"
	next := wordGenerator(phrase)
	stats := countDigitsInWords(next)
	printStats(stats)
}

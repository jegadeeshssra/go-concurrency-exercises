// Do the following:
// Start a goroutine.
// In this goroutine, fetch words from the generator, count the digits in each, and write to the counted channel.
// In the outer function, read values from the channel and fill stats.

// OUTPUT :
// 0ne: 1
// 1wo: 1
// thr33: 2
// 4068: 4

// Generator with goroutines.
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

// countDigitsInWords counts the number of digits in words,
// fetching the next word with next().
func countDigitsInWords(next func() string) counter {
	counted := make(chan pair)

	go func() {
		// Fetch words from the generator,
		// count the number of digits in each,
		// and write it to the counted channel.
		for {
			word := next()
			if word == "" {
				close(counted)
				break
			}
			count := countDigits(word)
			counted <- pair{word, count}
		}
	}()

	// Read values from the counted channel
	// and fill stats.

	// As a result, stats should contain words
	// and the number of digits in each.
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

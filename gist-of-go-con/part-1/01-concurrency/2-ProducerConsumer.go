// Do the following:

// Start a goroutine.
// In this goroutine, loop through the words, count the digits in each, and write to the counted channel (producer).
// In the outer function, read values from the channel and fill the stats counter (consumer).

// Result channel.
package main

import (
	"fmt"
	"strings"
	"unicode"
)

// counter stores the number of digits in each word.
// The key is the word, and the value is the number of digits.
type counter map[string]int

// solution start

// countDigitsInWords counts the number of digits in the words of a phrase.
func countDigitsInWords(phrase string) counter {
	words := strings.Fields(phrase)
	counted := make(chan int)

	go func() {
		// Loop through the words,
		// count the number of digits in each,
		// and write it to the counted channel.
		for _, word := range words {
			count := countDigits(word)
			counted <- count
		}
		close(counted)

	}()

	// Read values from the counted channel
	// and fill stats.

	// As a result, stats should contain words
	// and the number of digits in each.

	stats := counter{}
	for _, word := range words {
		count := <-counted
		stats[word] = count
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

func main() {
	phrase := "0ne 1wo thr33 4068"
	stats := countDigitsInWords(phrase)
	printStats(stats)
}

package jumble

import (
	_ "embed"
	"math/rand"
	"strings"
)

//go:embed words.txt
var wordListRaw string

// wordsByLength maps word length to a slice of words
var wordsByLength map[int][]string

func init() {
	wordsByLength = make(map[int][]string)
	lines := strings.Split(strings.TrimSpace(wordListRaw), "\n")
	for _, w := range lines {
		w = strings.TrimSpace(strings.ToUpper(w))
		if len(w) >= 4 && len(w) <= 8 {
			wordsByLength[len(w)] = append(wordsByLength[len(w)], w)
		}
	}
}

// RandomWordOfLength returns a random word of the given length.
func RandomWordOfLength(length int) string {
	words := wordsByLength[length]
	if len(words) == 0 {
		return ""
	}
	return words[rand.Intn(len(words))]
}

// WordCount returns the number of words of a given length.
func WordCount(length int) int {
	return len(wordsByLength[length])
}

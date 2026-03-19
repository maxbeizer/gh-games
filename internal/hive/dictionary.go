package hive

import (
	_ "embed"
	"strings"
)

//go:embed words.txt
var dictRaw string

var dictionary map[string]bool

func init() {
	lines := strings.Split(strings.TrimSpace(dictRaw), "\n")
	dictionary = make(map[string]bool, len(lines))
	for _, w := range lines {
		w = strings.TrimSpace(strings.ToUpper(w))
		if len(w) >= 4 {
			dictionary[w] = true
		}
	}
}

// IsWord checks if a word is in the dictionary.
func IsWord(word string) bool {
	return dictionary[strings.ToUpper(word)]
}

// DictSize returns the number of words in the dictionary.
func DictSize() int {
	return len(dictionary)
}

// FindValidWords returns all dictionary words that can be made from the given
// letters where each word must contain the required center letter.
// Only letters in the allowed set may be used (letters can repeat).
func FindValidWords(letters [7]rune, center rune) []string {
	allowed := make(map[rune]bool, 7)
	for _, l := range letters {
		allowed[l] = true
	}

	var valid []string
	for word := range dictionary {
		if !strings.ContainsRune(word, center) {
			continue
		}
		ok := true
		for _, r := range word {
			if !allowed[r] {
				ok = false
				break
			}
		}
		if ok {
			valid = append(valid, word)
		}
	}
	return valid
}

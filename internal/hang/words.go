package hang

import (
	_ "embed"
	"math/rand"
	"strings"
)

//go:embed words.txt
var hangWordsRaw string

var wordList []string

func init() {
	lines := strings.Split(strings.TrimSpace(hangWordsRaw), "\n")
	wordList = make([]string, 0, len(lines))
	for _, w := range lines {
		w = strings.TrimSpace(strings.ToUpper(w))
		if len(w) >= 5 && len(w) <= 8 {
			wordList = append(wordList, w)
		}
	}
}

// RandomWord returns a random word from the hangman word list.
func RandomWord() string {
	return wordList[rand.Intn(len(wordList))]
}

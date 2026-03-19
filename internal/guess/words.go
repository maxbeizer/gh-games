package guess

import (
	_ "embed"
	"math/rand"
	"strings"
	"time"
)

//go:embed words.txt
var wordListRaw string

var ValidWords []string
var wordSet map[string]bool

func init() {
	lines := strings.Split(strings.TrimSpace(wordListRaw), "\n")
	ValidWords = make([]string, 0, len(lines))
	wordSet = make(map[string]bool, len(lines))
	for _, w := range lines {
		w = strings.TrimSpace(strings.ToUpper(w))
		if len(w) == 5 {
			ValidWords = append(ValidWords, w)
			wordSet[w] = true
		}
	}
}

func IsValidWord(word string) bool {
	return wordSet[strings.ToUpper(word)]
}

func RandomWord() string {
	return ValidWords[rand.Intn(len(ValidWords))]
}

func DailyWord() string {
	// Deterministic based on date
	t := time.Now()
	seed := int64(t.Year()*10000 + int(t.Month())*100 + t.Day())
	r := rand.New(rand.NewSource(seed))
	return ValidWords[r.Intn(len(ValidWords))]
}

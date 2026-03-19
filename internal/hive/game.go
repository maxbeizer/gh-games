package hive

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

// Game represents a Hive word puzzle game.
type Game struct {
	Letters  [7]rune         // all 7 letters
	Center   rune            // the required center letter (Letters[0])
	Found    []string        // words the player has found
	foundSet map[string]bool // dedup set for found words
	AllValid []string        // all possible valid words (the answer key)
	Score    int
	MaxScore int // score if all words found
}

// ScoreWord computes the score for a single word.
func ScoreWord(word string, isPangram bool) int {
	n := len([]rune(word))
	points := 0
	if n == 4 {
		points = 1
	} else if n > 4 {
		points = n
	}
	if isPangram {
		points += 7
	}
	return points
}

// NewGame generates a random puzzle with at least 20 valid words.
func NewGame() *Game {
	for {
		word := randomPangramCandidate()
		if word == "" {
			continue
		}
		letters := uniqueLetters(word)
		if len(letters) != 7 {
			continue
		}

		// Shuffle letters and pick center
		var arr [7]rune
		perm := rand.Perm(7)
		for i, p := range perm {
			arr[i] = letters[p]
		}

		g := newGameFromLetters(arr)
		if len(g.AllValid) >= 20 {
			return g
		}
	}
}

// NewGameWithLetters creates a game with specific letters (for testing).
func NewGameWithLetters(letters [7]rune) *Game {
	return newGameFromLetters(letters)
}

func newGameFromLetters(letters [7]rune) *Game {
	center := letters[0]
	allValid := FindValidWords(letters, center)
	sort.Strings(allValid)

	maxScore := 0
	for _, w := range allValid {
		maxScore += ScoreWord(w, isPangramWord(w, letters))
	}

	return &Game{
		Letters:  letters,
		Center:   center,
		Found:    []string{},
		foundSet: make(map[string]bool),
		AllValid: allValid,
		Score:    0,
		MaxScore: maxScore,
	}
}

// Guess validates and scores a player's guess.
func (g *Game) Guess(word string) (points int, pangram bool, err error) {
	word = strings.ToLower(strings.TrimSpace(word))
	runes := []rune(word)

	if len(runes) < 4 {
		return 0, false, fmt.Errorf("too short — words must be at least 4 letters")
	}

	if !containsRune(runes, g.Center) {
		return 0, false, fmt.Errorf("missing center letter '%c'", g.Center)
	}

	letterSet := make(map[rune]bool, 7)
	for _, l := range g.Letters {
		letterSet[l] = true
	}
	for _, r := range runes {
		if !letterSet[r] {
			return 0, false, fmt.Errorf("letter '%c' is not in the puzzle", r)
		}
	}

	if !IsWord(word) {
		return 0, false, fmt.Errorf("'%s' is not a valid word", word)
	}

	if g.foundSet[word] {
		return 0, false, fmt.Errorf("already found '%s'", word)
	}

	pangram = g.IsPangram(word)
	points = ScoreWord(word, pangram)

	g.Found = append(g.Found, word)
	g.foundSet[word] = true
	g.Score += points

	return points, pangram, nil
}

// IsPangram returns true if the word uses all 7 unique letters.
func (g *Game) IsPangram(word string) bool {
	return isPangramWord(word, g.Letters)
}

// Progress returns the ratio of found words to all valid words.
func (g *Game) Progress() float64 {
	if len(g.AllValid) == 0 {
		return 0
	}
	return float64(len(g.Found)) / float64(len(g.AllValid))
}

// Rank returns a label based on the player's score as a percentage of MaxScore.
func (g *Game) Rank() string {
	if g.MaxScore == 0 {
		return "Beginner"
	}
	pct := float64(g.Score) / float64(g.MaxScore) * 100

	switch {
	case pct >= 100:
		return "Queen Bee"
	case pct >= 85:
		return "Genius"
	case pct >= 70:
		return "Amazing"
	case pct >= 50:
		return "Great"
	case pct >= 25:
		return "Nice"
	case pct >= 10:
		return "Good"
	default:
		return "Beginner"
	}
}

// --- helpers ---

func containsRune(runes []rune, target rune) bool {
	for _, r := range runes {
		if r == target {
			return true
		}
	}
	return false
}

func isPangramWord(word string, letters [7]rune) bool {
	runes := []rune(word)
	for _, l := range letters {
		if !containsRune(runes, l) {
			return false
		}
	}
	return true
}

func uniqueLetters(word string) []rune {
	seen := make(map[rune]bool)
	var out []rune
	for _, r := range strings.ToLower(word) {
		if !seen[r] {
			seen[r] = true
			out = append(out, r)
		}
	}
	return out
}

// randomPangramCandidate picks a random valid dictionary word with exactly 7 unique letters.
func randomPangramCandidate() string {
	// Use FindValidWords with a broad set won't work here; we need a word
	// with exactly 7 unique letters. We'll try random letter combos via the
	// dictionary. Instead, pick random 7-letter sets and check for pangrams.
	// A simpler approach: generate random letters, find words, look for pangrams.
	vowels := []rune("aeiou")
	consonants := []rune("bcdfghjklmnpqrstvwxyz")

	// Pick 2-3 vowels and 4-5 consonants for a good mix
	numVowels := 2 + rand.Intn(2) // 2 or 3
	numConsonants := 7 - numVowels

	var letters [7]rune
	idx := 0

	// Pick unique vowels
	vPerm := rand.Perm(len(vowels))
	for i := 0; i < numVowels; i++ {
		letters[idx] = vowels[vPerm[i]]
		idx++
	}

	// Pick unique consonants
	cPerm := rand.Perm(len(consonants))
	for i := 0; i < numConsonants; i++ {
		letters[idx] = consonants[cPerm[i]]
		idx++
	}

	center := letters[0]
	words := FindValidWords(letters, center)

	// Look for a pangram among the found words
	for _, w := range words {
		if isPangramWord(w, letters) {
			return w
		}
	}
	return ""
}

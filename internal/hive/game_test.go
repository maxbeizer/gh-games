package hive

import (
	"testing"
	"time"
)

// testLetters is a fixed set of 7 letters for deterministic tests.
// Center letter (Letters[0]) is 'a'.
var testLetters = [7]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g'}

func TestNewGameWithLetters(t *testing.T) {
	g := NewGameWithLetters(testLetters)

	if g.Center != 'a' {
		t.Errorf("expected center='a', got='%c'", g.Center)
	}
	if g.Letters != testLetters {
		t.Errorf("expected letters=%v, got=%v", testLetters, g.Letters)
	}
	if g.Score != 0 {
		t.Errorf("expected initial score=0, got=%d", g.Score)
	}
	if len(g.Found) != 0 {
		t.Errorf("expected no found words, got=%d", len(g.Found))
	}
	if len(g.AllValid) == 0 {
		t.Log("warning: AllValid is empty — dictionary may not have words for these letters")
	}
	if g.MaxScore < 0 {
		t.Errorf("MaxScore should be non-negative, got=%d", g.MaxScore)
	}
}

func TestGuessValid(t *testing.T) {
	g := NewGameWithLetters(testLetters)
	if len(g.AllValid) == 0 {
		t.Skip("no valid words for test letters — dictionary may not be loaded")
	}

	word := g.AllValid[0]
	points, _, err := g.Guess(word)
	if err != nil {
		t.Fatalf("expected no error for valid word %q, got: %v", word, err)
	}
	if points <= 0 {
		t.Errorf("expected positive points, got=%d", points)
	}
	if g.Score != points {
		t.Errorf("expected score=%d, got=%d", points, g.Score)
	}
}

func TestGuessPangram(t *testing.T) {
	// Create letters that spell a known pangram
	letters := [7]rune{'b', 'a', 'c', 'k', 'i', 'n', 'g'}
	g := NewGameWithLetters(letters)

	// "backing" uses all 7 letters
	points, pangram, err := g.Guess("backing")
	if err != nil {
		t.Skipf("'backing' not in dictionary or not valid: %v", err)
	}
	if !pangram {
		t.Error("expected pangram=true for 'backing'")
	}
	// 7-letter word (7 points) + pangram bonus (7) = 14
	expected := 7 + 7
	if points != expected {
		t.Errorf("expected %d points for pangram, got=%d", expected, points)
	}
}

func TestGuessTooShort(t *testing.T) {
	g := NewGameWithLetters(testLetters)
	_, _, err := g.Guess("abc")
	if err == nil {
		t.Error("expected error for 3-letter word")
	}
}

func TestGuessMissingCenter(t *testing.T) {
	g := NewGameWithLetters(testLetters) // center is 'a'
	_, _, err := g.Guess("bcde")
	if err == nil {
		t.Error("expected error for word missing center letter")
	}
}

func TestGuessInvalidLetter(t *testing.T) {
	g := NewGameWithLetters(testLetters) // letters a-g only
	_, _, err := g.Guess("abcz")
	if err == nil {
		t.Error("expected error for word with letter not in puzzle")
	}
}

func TestGuessDuplicate(t *testing.T) {
	g := NewGameWithLetters(testLetters)
	if len(g.AllValid) == 0 {
		t.Skip("no valid words for test letters")
	}

	word := g.AllValid[0]
	_, _, err := g.Guess(word)
	if err != nil {
		t.Skipf("first guess failed: %v", err)
	}

	_, _, err = g.Guess(word)
	if err == nil {
		t.Error("expected error for duplicate guess")
	}
}

func TestIsPangram(t *testing.T) {
	letters := [7]rune{'b', 'a', 'c', 'k', 'i', 'n', 'g'}
	g := NewGameWithLetters(letters)

	if !g.IsPangram("backing") {
		t.Error("expected 'backing' to be a pangram")
	}
	if g.IsPangram("back") {
		t.Error("expected 'back' NOT to be a pangram")
	}
}

func TestScoreWord(t *testing.T) {
	tests := []struct {
		word     string
		pangram  bool
		expected int
	}{
		{"abcd", false, 1},       // 4-letter: 1 point
		{"abcde", false, 5},      // 5-letter: 5 points
		{"abcdef", false, 6},     // 6-letter: 6 points
		{"abcdefg", false, 7},    // 7-letter: 7 points
		{"abcdefg", true, 14},    // 7-letter pangram: 7 + 7
		{"abcde", true, 12},      // 5-letter pangram: 5 + 7
		{"abcd", true, 8},        // 4-letter pangram: 1 + 7
	}
	for _, tc := range tests {
		got := ScoreWord(tc.word, tc.pangram)
		if got != tc.expected {
			t.Errorf("ScoreWord(%q, %v) = %d, want %d", tc.word, tc.pangram, got, tc.expected)
		}
	}
}

func TestRank(t *testing.T) {
	g := &Game{MaxScore: 100, foundSet: make(map[string]bool)}

	tests := []struct {
		score    int
		expected string
	}{
		{0, "Beginner"},
		{5, "Beginner"},
		{10, "Good"},
		{25, "Nice"},
		{50, "Great"},
		{70, "Amazing"},
		{85, "Genius"},
		{100, "Queen Bee"},
	}
	for _, tc := range tests {
		g.Score = tc.score
		got := g.Rank()
		if got != tc.expected {
			t.Errorf("Rank() at score %d = %q, want %q", tc.score, got, tc.expected)
		}
	}
}

func TestProgress(t *testing.T) {
	g := &Game{
		AllValid: []string{"aaaa", "bbbb", "cccc", "dddd"},
		Found:    []string{},
		foundSet: make(map[string]bool),
	}

	if g.Progress() != 0.0 {
		t.Errorf("expected progress=0.0, got=%f", g.Progress())
	}

	g.Found = append(g.Found, "aaaa")
	g.foundSet["aaaa"] = true
	if g.Progress() != 0.25 {
		t.Errorf("expected progress=0.25, got=%f", g.Progress())
	}

	g.Found = append(g.Found, "bbbb")
	g.foundSet["bbbb"] = true
	if g.Progress() != 0.5 {
		t.Errorf("expected progress=0.5, got=%f", g.Progress())
	}
}

func TestNewGamePerformance(t *testing.T) {
	start := time.Now()
	g := NewGame()
	elapsed := time.Since(start)

	if elapsed > 2*time.Second {
		t.Errorf("NewGame() took %v, want < 2s", elapsed)
	}
	if len(g.AllValid) == 0 {
		t.Error("NewGame() produced no valid words")
	}
	t.Logf("NewGame() completed in %v with %d valid words", elapsed, len(g.AllValid))
}

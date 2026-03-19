package hive

import (
	"strings"
	"testing"
)

func TestDictSize(t *testing.T) {
	size := DictSize()
	if size < 100000 {
		t.Errorf("DictSize() = %d, want >= 100000", size)
	}
	t.Logf("Dictionary contains %d words", size)
}

func TestDictAllWordsValid(t *testing.T) {
	for word := range dictionary {
		if len(word) < 4 {
			t.Errorf("word %q has length %d, want >= 4", word, len(word))
		}
		if word != strings.ToUpper(word) {
			t.Errorf("word %q is not uppercase", word)
		}
		for _, r := range word {
			if r < 'A' || r > 'Z' {
				t.Errorf("word %q contains non-letter rune %c", word, r)
			}
		}
	}
}

func TestDictIsWord(t *testing.T) {
	validWords := []string{"ABLE", "HELLO", "WORLD", "GAME", "QUEEN", "ZONE", "ABOUT", "HOUSE"}
	for _, w := range validWords {
		if !IsWord(w) {
			t.Errorf("IsWord(%q) = false, want true", w)
		}
	}
	// Test case-insensitivity
	if !IsWord("able") {
		t.Error("IsWord(\"able\") = false, want true (case-insensitive)")
	}

	invalidWords := []string{"XYZ", "QQQ", "ZZZ", "ABCDEFGHIJKLMNOP", ""}
	for _, w := range invalidWords {
		if IsWord(w) {
			t.Errorf("IsWord(%q) = true, want false", w)
		}
	}
}

func TestDictFindValidWords(t *testing.T) {
	// Use letters that should produce known results
	// Letters: A, B, C, D, E, L, K with center = A
	letters := [7]rune{'A', 'B', 'C', 'D', 'E', 'L', 'K'}
	center := 'A'

	valid := FindValidWords(letters, center)
	if len(valid) == 0 {
		t.Fatal("FindValidWords returned no words; expected at least some")
	}

	// Every returned word must contain the center letter
	for _, w := range valid {
		if !strings.ContainsRune(w, center) {
			t.Errorf("word %q does not contain center letter %c", w, center)
		}
	}

	// Every returned word must only use allowed letters
	allowed := map[rune]bool{'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'L': true, 'K': true}
	for _, w := range valid {
		for _, r := range w {
			if !allowed[r] {
				t.Errorf("word %q contains disallowed letter %c", w, r)
			}
		}
	}

	t.Logf("Found %d valid words for letters %c (center %c)", len(valid), letters, center)
}

func TestDictFindValidWordsRequiresCenter(t *testing.T) {
	// Letters: A, B, C, D, E, F, G with center = G
	letters := [7]rune{'A', 'B', 'C', 'D', 'E', 'F', 'G'}
	center := 'G'

	valid := FindValidWords(letters, center)
	for _, w := range valid {
		if !strings.ContainsRune(w, center) {
			t.Errorf("word %q missing required center letter %c", w, center)
		}
	}
}

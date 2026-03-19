package guess

import (
	"testing"
	"unicode"
)

func TestWordsCount(t *testing.T) {
	if len(ValidWords) < 500 {
		t.Errorf("expected at least 500 words, got %d", len(ValidWords))
	}
	t.Logf("word list contains %d words", len(ValidWords))
}

func TestWordsLength(t *testing.T) {
	for _, w := range ValidWords {
		if len(w) != 5 {
			t.Errorf("word %q has length %d, expected 5", w, len(w))
		}
	}
}

func TestWordsUppercase(t *testing.T) {
	for _, w := range ValidWords {
		for _, r := range w {
			if !unicode.IsUpper(r) {
				t.Errorf("word %q contains non-uppercase character", w)
				break
			}
		}
	}
}

func TestIsValidWord(t *testing.T) {
	if !IsValidWord("ABOUT") {
		t.Error("expected ABOUT to be valid")
	}
	if !IsValidWord("about") {
		t.Error("expected about (lowercase) to be valid via case-insensitive lookup")
	}
	if IsValidWord("ZZZZZ") {
		t.Error("expected ZZZZZ to be invalid")
	}
	if IsValidWord("XYZAB") {
		t.Error("expected XYZAB to be invalid")
	}
}

func TestDailyWordDeterministic(t *testing.T) {
	w1 := DailyWord()
	w2 := DailyWord()
	if w1 != w2 {
		t.Errorf("DailyWord not deterministic: got %q and %q", w1, w2)
	}
}

func TestRandomWordValid(t *testing.T) {
	for i := 0; i < 100; i++ {
		w := RandomWord()
		if !IsValidWord(w) {
			t.Errorf("RandomWord returned invalid word: %q", w)
		}
		if len(w) != 5 {
			t.Errorf("RandomWord returned word of length %d: %q", len(w), w)
		}
	}
}

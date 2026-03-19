package hang

import (
	"testing"
	"unicode"
)

func TestCorrectGuessRevealsLetter(t *testing.T) {
	g := NewGameWithWord("HELLO")
	correct, err := g.GuessLetter('H')
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !correct {
		t.Fatal("expected correct=true for letter in word")
	}
	if g.WrongCount() != 0 {
		t.Fatalf("wrong count should be 0, got %d", g.WrongCount())
	}
	// 'H' revealed, rest blanks
	want := "H _ _ _ _"
	if got := g.Display(); got != want {
		t.Fatalf("Display() = %q, want %q", got, want)
	}
}

func TestWrongGuessIncrementsCount(t *testing.T) {
	g := NewGameWithWord("HELLO")
	correct, err := g.GuessLetter('Z')
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if correct {
		t.Fatal("expected correct=false for letter not in word")
	}
	if g.WrongCount() != 1 {
		t.Fatalf("wrong count should be 1, got %d", g.WrongCount())
	}
}

func TestDuplicateGuessReturnsError(t *testing.T) {
	g := NewGameWithWord("HELLO")
	_, _ = g.GuessLetter('H')
	_, err := g.GuessLetter('H')
	if err == nil {
		t.Fatal("expected error for duplicate guess")
	}
}

func TestInvalidCharacterReturnsError(t *testing.T) {
	g := NewGameWithWord("HELLO")
	_, err := g.GuessLetter('1')
	if err == nil {
		t.Fatal("expected error for non-letter character")
	}
}

func TestWinCondition(t *testing.T) {
	g := NewGameWithWord("HI")
	g.GuessLetter('H')
	g.GuessLetter('I')
	if !g.IsWon() {
		t.Fatal("expected game to be won")
	}
	if !g.IsOver() {
		t.Fatal("expected game to be over")
	}
	if g.IsLost() {
		t.Fatal("expected game to not be lost")
	}
}

func TestLoseCondition(t *testing.T) {
	g := NewGameWithWord("HI")
	wrongs := []rune{'A', 'B', 'C', 'D', 'E', 'F'}
	for _, r := range wrongs {
		_, err := g.GuessLetter(r)
		if err != nil {
			t.Fatalf("unexpected error guessing '%c': %v", r, err)
		}
	}
	if !g.IsLost() {
		t.Fatal("expected game to be lost after 6 wrong guesses")
	}
	if !g.IsOver() {
		t.Fatal("expected game to be over")
	}
	if g.IsWon() {
		t.Fatal("expected game to not be won")
	}
}

func TestGuessAfterGameOverReturnsError(t *testing.T) {
	g := NewGameWithWord("HI")
	g.GuessLetter('H')
	g.GuessLetter('I')
	_, err := g.GuessLetter('A')
	if err == nil {
		t.Fatal("expected error when guessing after game over")
	}
}

func TestLowercaseGuessWorks(t *testing.T) {
	g := NewGameWithWord("HELLO")
	correct, err := g.GuessLetter('h')
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !correct {
		t.Fatal("expected lowercase 'h' to match uppercase target")
	}
}

func TestDisplayShowsAllBlanksInitially(t *testing.T) {
	g := NewGameWithWord("CAT")
	want := "_ _ _"
	if got := g.Display(); got != want {
		t.Fatalf("Display() = %q, want %q", got, want)
	}
}

func TestDisplayFullyRevealed(t *testing.T) {
	g := NewGameWithWord("CAT")
	g.GuessLetter('C')
	g.GuessLetter('A')
	g.GuessLetter('T')
	want := "C A T"
	if got := g.Display(); got != want {
		t.Fatalf("Display() = %q, want %q", got, want)
	}
}

func TestWrongLettersTracked(t *testing.T) {
	g := NewGameWithWord("CAT")
	g.GuessLetter('Z')
	g.GuessLetter('X')
	wrongs := g.WrongLetters()
	if len(wrongs) != 2 {
		t.Fatalf("expected 2 wrong letters, got %d", len(wrongs))
	}
	if wrongs[0] != 'Z' || wrongs[1] != 'X' {
		t.Fatalf("wrong letters = %v, want [Z X]", wrongs)
	}
}

func TestNewGamePicksValidWord(t *testing.T) {
	g := NewGame()
	word := g.Target
	if len(word) < 5 || len(word) > 8 {
		t.Fatalf("target word %q length %d not in [5,8]", word, len(word))
	}
	for _, r := range word {
		if !unicode.IsUpper(r) {
			t.Fatalf("target word %q contains non-uppercase letter", word)
		}
	}
}

func TestIsGuessedAndIsCorrect(t *testing.T) {
	g := NewGameWithWord("HELLO")
	g.GuessLetter('H')
	g.GuessLetter('Z')

	if !g.IsGuessed('H') {
		t.Fatal("expected H to be guessed")
	}
	if !g.IsCorrect('H') {
		t.Fatal("expected H to be correct")
	}
	if !g.IsGuessed('Z') {
		t.Fatal("expected Z to be guessed")
	}
	if !g.IsWrongGuess('Z') {
		t.Fatal("expected Z to be a wrong guess")
	}
	if g.IsGuessed('A') {
		t.Fatal("expected A to not be guessed")
	}
}

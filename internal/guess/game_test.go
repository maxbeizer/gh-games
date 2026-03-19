package guess

import (
	"testing"
)

func TestCheckGuess(t *testing.T) {
	tests := []struct {
		name   string
		target string
		guess  string
		want   [5]LetterState
	}{
		{
			name:   "all correct",
			target: "HELLO",
			guess:  "HELLO",
			want:   [5]LetterState{Correct, Correct, Correct, Correct, Correct},
		},
		{
			name:   "all absent",
			target: "HELLO",
			guess:  "DUCKS",
			want:   [5]LetterState{Absent, Absent, Absent, Absent, Absent},
		},
		{
			name:   "mixed states",
			target: "HELLO",
			guess:  "HEAPS",
			want:   [5]LetterState{Correct, Correct, Absent, Absent, Absent},
		},
		{
			name:   "present letters with one correct",
			target: "HELLO",
			guess:  "OLELH",
			want:   [5]LetterState{Present, Present, Present, Correct, Present},
		},
		{
			name:   "duplicate letters in guess - ABBEY vs BABBY",
			target: "ABBEY",
			guess:  "BABBY",
			// B@0→Present(matched to target[1]), A@1→Present(target[0]),
			// B@2→Correct, B@3→Absent(no more B's), Y@4→Correct
			want: [5]LetterState{Present, Present, Correct, Absent, Correct},
		},
		{
			name:   "duplicate guess letter with one correct and one present",
			target: "ABACK",
			guess:  "AALII",
			// A@0→Correct, A@1→Present(target[2]), L@2→Absent, I@3→Absent, I@4→Absent
			want: [5]LetterState{Correct, Present, Absent, Absent, Absent},
		},
		{
			name:   "extra duplicate in guess should be absent",
			target: "STARE",
			guess:  "SEEDS",
			// S@0→Correct, E@1→Present(target[4]), E@2→Absent(no more E's),
			// D@3→Absent, S@4→Absent(target[0] already matched)
			want: [5]LetterState{Correct, Present, Absent, Absent, Absent},
		},
		{
			name:   "case insensitive",
			target: "HELLO",
			guess:  "hello",
			want:   [5]LetterState{Correct, Correct, Correct, Correct, Correct},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckGuess(tt.target, tt.guess)
			if got != tt.want {
				t.Errorf("CheckGuess(%q, %q) = %v, want %v", tt.target, tt.guess, got, tt.want)
			}
		})
	}
}

func TestNewGame(t *testing.T) {
	g := NewGame("hello")
	if g.Target != "HELLO" {
		t.Errorf("Target = %q, want HELLO", g.Target)
	}
	if g.MaxTurns != 6 {
		t.Errorf("MaxTurns = %d, want 6", g.MaxTurns)
	}
	if len(g.Guesses) != 0 {
		t.Errorf("Guesses should be empty")
	}
	if g.Keyboard == nil {
		t.Error("Keyboard should be initialized")
	}
}

func TestGuessValidation(t *testing.T) {
	t.Run("too short", func(t *testing.T) {
		g := NewGame("HELLO")
		_, err := g.Guess("HI")
		if err == nil {
			t.Error("expected error for short guess")
		}
	})

	t.Run("too long", func(t *testing.T) {
		g := NewGame("HELLO")
		_, err := g.Guess("HELLOO")
		if err == nil {
			t.Error("expected error for long guess")
		}
	})

	t.Run("valid guess", func(t *testing.T) {
		g := NewGame("HELLO")
		result, err := g.Guess("HEAPS")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Word != "HEAPS" {
			t.Errorf("Word = %q, want HEAPS", result.Word)
		}
	})
}

func TestGameFlow(t *testing.T) {
	t.Run("win on first guess", func(t *testing.T) {
		g := NewGame("HELLO")
		g.Guess("HELLO")
		if !g.IsWon() {
			t.Error("expected game to be won")
		}
		if g.IsLost() {
			t.Error("should not be lost")
		}
		if !g.IsOver() {
			t.Error("should be over")
		}
	})

	t.Run("lose after 6 wrong guesses", func(t *testing.T) {
		g := NewGame("HELLO")
		for i := 0; i < 6; i++ {
			g.Guess("DUCKS")
		}
		if g.IsWon() {
			t.Error("should not be won")
		}
		if !g.IsLost() {
			t.Error("expected game to be lost")
		}
		if !g.IsOver() {
			t.Error("should be over")
		}
	})

	t.Run("cannot guess after game over", func(t *testing.T) {
		g := NewGame("HELLO")
		g.Guess("HELLO")
		_, err := g.Guess("WORLD")
		if err == nil {
			t.Error("expected error when guessing after game over")
		}
	})
}

func TestKeyboardTracking(t *testing.T) {
	g := NewGame("HELLO")
	g.Guess("HEAPS")
	if g.Keyboard['H'] != Correct {
		t.Errorf("H should be Correct, got %d", g.Keyboard['H'])
	}
	if g.Keyboard['A'] != Absent {
		t.Errorf("A should be Absent, got %d", g.Keyboard['A'])
	}

	// Keyboard should upgrade from Absent/Present to Correct
	g.Guess("HELLO")
	if g.Keyboard['L'] != Correct {
		t.Errorf("L should be Correct after exact match, got %d", g.Keyboard['L'])
	}
}

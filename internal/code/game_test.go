package code

import "testing"

func TestComputeFeedback_AllExact(t *testing.T) {
	secret := [4]Color{Red, Green, Blue, Yellow}
	guess := [4]Color{Red, Green, Blue, Yellow}
	fb := ComputeFeedback(secret, guess)
	if fb.Exact != 4 || fb.Misplaced != 0 {
		t.Errorf("all exact: got Exact=%d Misplaced=%d, want 4/0", fb.Exact, fb.Misplaced)
	}
}

func TestComputeFeedback_NoMatches(t *testing.T) {
	secret := [4]Color{Red, Red, Red, Red}
	guess := [4]Color{Blue, Blue, Blue, Blue}
	fb := ComputeFeedback(secret, guess)
	if fb.Exact != 0 || fb.Misplaced != 0 {
		t.Errorf("no matches: got Exact=%d Misplaced=%d, want 0/0", fb.Exact, fb.Misplaced)
	}
}

func TestComputeFeedback_AllMisplaced(t *testing.T) {
	secret := [4]Color{Red, Green, Blue, Yellow}
	guess := [4]Color{Yellow, Blue, Green, Red}
	fb := ComputeFeedback(secret, guess)
	if fb.Exact != 0 || fb.Misplaced != 4 {
		t.Errorf("all misplaced: got Exact=%d Misplaced=%d, want 0/4", fb.Exact, fb.Misplaced)
	}
}

func TestComputeFeedback_MixedExactAndMisplaced(t *testing.T) {
	secret := [4]Color{Red, Green, Blue, Yellow}
	guess := [4]Color{Red, Blue, Green, Yellow}
	fb := ComputeFeedback(secret, guess)
	// R@0 exact, Y@3 exact, B and G misplaced
	if fb.Exact != 2 || fb.Misplaced != 2 {
		t.Errorf("mixed: got Exact=%d Misplaced=%d, want 2/2", fb.Exact, fb.Misplaced)
	}
}

func TestComputeFeedback_DuplicateColors_ExtraGuessesDontCount(t *testing.T) {
	// Secret: R R G B, Guess: R G G G
	// Exact: R@0, G@2. The G@1 and G@3 have no unmatched secret slots → 0 misplaced.
	secret := [4]Color{Red, Red, Green, Blue}
	guess := [4]Color{Red, Green, Green, Green}
	fb := ComputeFeedback(secret, guess)
	if fb.Exact != 2 || fb.Misplaced != 0 {
		t.Errorf("duplicate extra guesses: got Exact=%d Misplaced=%d, want 2/0", fb.Exact, fb.Misplaced)
	}
}

func TestComputeFeedback_DuplicateColors_SecretHasDuplicates(t *testing.T) {
	// Secret: R R G B, Guess: G R R R
	// Exact: R@1. Misplaced: R@2 matches secret@0, G@0 matches secret@2 → 2 misplaced.
	// R@3 has no unmatched secret R left → not counted.
	secret := [4]Color{Red, Red, Green, Blue}
	guess := [4]Color{Green, Red, Red, Red}
	fb := ComputeFeedback(secret, guess)
	if fb.Exact != 1 || fb.Misplaced != 2 {
		t.Errorf("duplicate secret: got Exact=%d Misplaced=%d, want 1/2", fb.Exact, fb.Misplaced)
	}
}

func TestComputeFeedback_AllSameColor(t *testing.T) {
	secret := [4]Color{Red, Red, Red, Red}
	guess := [4]Color{Red, Red, Red, Red}
	fb := ComputeFeedback(secret, guess)
	if fb.Exact != 4 || fb.Misplaced != 0 {
		t.Errorf("all same: got Exact=%d Misplaced=%d, want 4/0", fb.Exact, fb.Misplaced)
	}
}

func TestComputeFeedback_OneMisplacedWithDuplicates(t *testing.T) {
	// Secret: R G B Y, Guess: P P P R → R is misplaced (secret@0), rest wrong
	secret := [4]Color{Red, Green, Blue, Yellow}
	guess := [4]Color{Purple, Purple, Purple, Red}
	fb := ComputeFeedback(secret, guess)
	if fb.Exact != 0 || fb.Misplaced != 1 {
		t.Errorf("one misplaced: got Exact=%d Misplaced=%d, want 0/1", fb.Exact, fb.Misplaced)
	}
}

func TestGame_Win(t *testing.T) {
	secret := [4]Color{Red, Green, Blue, Yellow}
	g := NewGameWithSecret(secret)
	g.MakeGuess([4]Color{Purple, Purple, Purple, Purple}) // wrong
	g.MakeGuess(secret)                                    // correct
	if !g.IsWon() {
		t.Error("expected IsWon() to be true")
	}
	if g.IsLost() {
		t.Error("expected IsLost() to be false")
	}
	if !g.IsOver() {
		t.Error("expected IsOver() to be true")
	}
	if len(g.Guesses) != 2 {
		t.Errorf("expected 2 guesses, got %d", len(g.Guesses))
	}
}

func TestGame_Lose(t *testing.T) {
	secret := [4]Color{Red, Green, Blue, Yellow}
	g := NewGameWithSecret(secret)
	wrong := [4]Color{Purple, Purple, Purple, Purple}
	for i := 0; i < 10; i++ {
		g.MakeGuess(wrong)
	}
	if g.IsWon() {
		t.Error("expected IsWon() to be false")
	}
	if !g.IsLost() {
		t.Error("expected IsLost() to be true")
	}
	if !g.IsOver() {
		t.Error("expected IsOver() to be true")
	}
}

func TestGame_NotOverYet(t *testing.T) {
	g := NewGameWithSecret([4]Color{Red, Green, Blue, Yellow})
	g.MakeGuess([4]Color{Red, Red, Red, Red})
	if g.IsOver() {
		t.Error("game should not be over after 1 wrong guess")
	}
}

func TestNewGame_CreatesValidSecret(t *testing.T) {
	g := NewGame()
	for i, c := range g.Secret {
		if c < 0 || c >= NumColors {
			t.Errorf("secret[%d] = %d, out of valid color range", i, c)
		}
	}
	if g.MaxTurns != 10 {
		t.Errorf("expected MaxTurns=10, got %d", g.MaxTurns)
	}
}

func TestColorName(t *testing.T) {
	names := []string{"Red", "Green", "Blue", "Yellow", "Purple", "Orange"}
	for i, want := range names {
		if got := ColorName(Color(i)); got != want {
			t.Errorf("ColorName(%d) = %q, want %q", i, got, want)
		}
	}
}

func TestColorSymbol(t *testing.T) {
	symbols := []string{"🔴", "🟢", "🔵", "🟡", "🟣", "🟠"}
	for i, want := range symbols {
		if got := ColorSymbol(Color(i)); got != want {
			t.Errorf("ColorSymbol(%d) = %q, want %q", i, got, want)
		}
	}
}

package jumble

import (
	"testing"
	"time"
)

func TestScrambleWord(t *testing.T) {
	t.Run("produces different arrangement", func(t *testing.T) {
		word := "JUMBLE"
		scrambled := ScrambleWord(word)
		if scrambled == word {
			t.Errorf("ScrambleWord(%q) returned same word", word)
		}
		if len(scrambled) != len(word) {
			t.Errorf("ScrambleWord length %d != original %d", len(scrambled), len(word))
		}
	})

	t.Run("contains same letters", func(t *testing.T) {
		word := "STORM"
		scrambled := ScrambleWord(word)
		sortedWord := sortString(word)
		sortedScrambled := sortString(scrambled)
		if sortedWord != sortedScrambled {
			t.Errorf("letters differ: %q vs %q", sortedWord, sortedScrambled)
		}
	})

	t.Run("short word still scrambles", func(t *testing.T) {
		word := "AB"
		scrambled := ScrambleWord(word)
		if scrambled == word {
			t.Errorf("ScrambleWord(%q) returned same word", word)
		}
	})
}

func TestNewGame(t *testing.T) {
	g := NewGame()
	if len(g.Rounds) != 5 {
		t.Fatalf("expected 5 rounds, got %d", len(g.Rounds))
	}
	expectedLengths := []int{4, 5, 6, 7, 8}
	for i, r := range g.Rounds {
		if len(r.Target) != expectedLengths[i] {
			t.Errorf("round %d: target length %d, want %d", i, len(r.Target), expectedLengths[i])
		}
		if r.Solved {
			t.Errorf("round %d should not be solved", i)
		}
		if r.Scrambled == r.Target {
			t.Errorf("round %d: scrambled should differ from target", i)
		}
	}
}

func TestGuessCorrect(t *testing.T) {
	g := NewGame()
	target := g.Rounds[0].Target
	correct, points := g.Guess(target)
	if !correct {
		t.Error("expected correct guess")
	}
	if points <= 0 {
		t.Errorf("expected positive points, got %d", points)
	}
	if !g.Rounds[0].Solved {
		t.Error("round should be marked solved")
	}
}

func TestGuessIncorrect(t *testing.T) {
	g := NewGame()
	correct, points := g.Guess("ZZZZ")
	if correct {
		t.Error("expected incorrect guess")
	}
	if points != 0 {
		t.Errorf("expected 0 points, got %d", points)
	}
}

func TestGuessCaseInsensitive(t *testing.T) {
	g := NewGame()
	target := g.Rounds[0].Target
	lower := ""
	for _, r := range target {
		lower += string(r + 32) // convert to lowercase
	}
	correct, _ := g.Guess(lower)
	if !correct {
		t.Error("guess should be case-insensitive")
	}
}

func TestHint(t *testing.T) {
	g := NewGame()
	letter, pos, err := g.Hint()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	target := g.Rounds[0].Target
	if rune(target[pos]) != letter {
		t.Errorf("hint letter %c at pos %d doesn't match target %q", letter, pos, target)
	}
	if g.Rounds[0].HintsUsed != 1 {
		t.Errorf("expected 1 hint used, got %d", g.Rounds[0].HintsUsed)
	}
}

func TestHintFixesPosition(t *testing.T) {
	g := NewGame()
	_, pos, _ := g.Hint()
	r := g.Rounds[0]
	if r.Scrambled[pos] != r.Target[pos] {
		t.Errorf("hinted position %d: scrambled has %c, target has %c",
			pos, r.Scrambled[pos], r.Target[pos])
	}
}

func TestHintAllRevealed(t *testing.T) {
	g := NewGame()
	wordLen := len(g.Rounds[0].Target)
	for i := 0; i < wordLen; i++ {
		_, _, err := g.Hint()
		if err != nil {
			t.Fatalf("hint %d failed: %v", i, err)
		}
	}
	_, _, err := g.Hint()
	if err == nil {
		t.Error("expected error when all letters revealed")
	}
}

func TestHintPenalty(t *testing.T) {
	g := NewGame()
	g.Hint()
	g.Hint()
	// Eliminate speed bonus by backdating start time
	g.Rounds[0].StartTime = time.Now().Add(-30 * time.Second)
	target := g.Rounds[0].Target
	_, points := g.Guess(target)

	// With 2 hints and no speed bonus: base (400) - penalty (100) = 300
	base := len(target) * 100
	if points >= base {
		t.Errorf("points %d should be less than base %d due to hint penalty", points, base)
	}
}

func TestScoring(t *testing.T) {
	g := NewGame()
	// Solve immediately for max speed bonus
	target := g.Rounds[0].Target
	_, points := g.Guess(target)

	// 4-letter word base = 400, speed bonus up to 200
	if points < 400 {
		t.Errorf("expected at least base points (400), got %d", points)
	}
}

func TestScoringMinimum(t *testing.T) {
	g := NewGame()
	// Use all possible hints to drive score down
	wordLen := len(g.Rounds[0].Target)
	for i := 0; i < wordLen-1; i++ {
		g.Hint()
	}
	// Wait a bit to eliminate speed bonus
	g.Rounds[0].StartTime = time.Now().Add(-30 * time.Second)
	target := g.Rounds[0].Target
	_, points := g.Guess(target)

	if points < 10 {
		t.Errorf("points should be at least 10 (minimum), got %d", points)
	}
}

func TestShufflePreservesHints(t *testing.T) {
	g := NewGame()
	_, pos, _ := g.Hint()
	g.Shuffle()
	r := g.Rounds[0]
	if r.Scrambled[pos] != r.Target[pos] {
		t.Errorf("shuffle broke hint at pos %d: got %c, want %c",
			pos, r.Scrambled[pos], r.Target[pos])
	}
}

func TestNextRound(t *testing.T) {
	g := NewGame()
	g.Rounds[0].Solved = true
	ok := g.NextRound()
	if !ok {
		t.Error("expected NextRound to succeed")
	}
	if g.CurrentRound != 1 {
		t.Errorf("expected round 1, got %d", g.CurrentRound)
	}
}

func TestIsComplete(t *testing.T) {
	g := NewGame()
	if g.IsComplete() {
		t.Error("game should not be complete at start")
	}
	// Solve all rounds
	for i := 0; i < 5; i++ {
		g.Guess(g.Rounds[i].Target)
		if i < 4 {
			g.NextRound()
		}
	}
	if !g.IsComplete() {
		t.Error("game should be complete after all rounds solved")
	}
}

func TestWordsByLength(t *testing.T) {
	for length := 4; length <= 8; length++ {
		count := WordCount(length)
		if count < 50 {
			t.Errorf("expected at least 50 words of length %d, got %d", length, count)
		}
	}
}

func TestRandomWordOfLength(t *testing.T) {
	for length := 4; length <= 8; length++ {
		w := RandomWordOfLength(length)
		if len(w) != length {
			t.Errorf("RandomWordOfLength(%d) returned %q (length %d)", length, w, len(w))
		}
	}
}

// helper to sort a string's characters for comparison
func sortString(s string) string {
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		for j := i + 1; j < len(runes); j++ {
			if runes[j] < runes[i] {
				runes[i], runes[j] = runes[j], runes[i]
			}
		}
	}
	return string(runes)
}

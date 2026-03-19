package guess

import (
	"fmt"
	"strings"

	"github.com/maxbeizer/gh-games/internal/common"
)

type LetterState int

const (
	Unknown LetterState = iota
	Absent
	Present
	Correct
)

type GuessResult struct {
	Word   string
	States [5]LetterState
}

type Game struct {
	Target   string
	Guesses  []GuessResult
	MaxTurns int
	Keyboard map[rune]LetterState
}

func NewGame(target string) *Game {
	return &Game{
		Target:   strings.ToUpper(target),
		Guesses:  []GuessResult{},
		MaxTurns: 6,
		Keyboard: make(map[rune]LetterState),
	}
}

func (g *Game) Guess(word string) (GuessResult, error) {
	word = strings.ToUpper(word)
	if len(word) != 5 {
		return GuessResult{}, fmt.Errorf("guess must be exactly 5 letters, got %d", len(word))
	}
	if g.IsOver() {
		return GuessResult{}, fmt.Errorf("game is already over")
	}

	states := CheckGuess(g.Target, word)
	result := GuessResult{Word: word, States: states}

	// Update keyboard with best known state per letter
	for i, r := range word {
		prev := g.Keyboard[r]
		cur := states[i]
		if cur > prev {
			g.Keyboard[r] = cur
		}
	}

	g.Guesses = append(g.Guesses, result)
	return result, nil
}

func (g *Game) IsWon() bool {
	if len(g.Guesses) == 0 {
		return false
	}
	last := g.Guesses[len(g.Guesses)-1]
	for _, s := range last.States {
		if s != Correct {
			return false
		}
	}
	return true
}

func (g *Game) IsLost() bool {
	return len(g.Guesses) >= g.MaxTurns && !g.IsWon()
}

func (g *Game) IsOver() bool {
	return g.IsWon() || g.IsLost()
}

// CheckGuess computes letter states handling duplicate letters correctly.
// First pass: mark exact (Correct) matches.
// Second pass: for remaining letters, mark Present if the letter exists
// in an unmatched target position, otherwise Absent.
func CheckGuess(target, guess string) [5]LetterState {
	var states [5]LetterState
	targetRunes := []rune(strings.ToUpper(target))
	guessRunes := []rune(strings.ToUpper(guess))

	// Track which target positions are still available
	matched := [5]bool{}

	// First pass: mark Correct
	for i := 0; i < 5; i++ {
		if guessRunes[i] == targetRunes[i] {
			states[i] = Correct
			matched[i] = true
		}
	}

	// Second pass: mark Present or Absent
	for i := 0; i < 5; i++ {
		if states[i] == Correct {
			continue
		}
		found := false
		for j := 0; j < 5; j++ {
			if !matched[j] && guessRunes[i] == targetRunes[j] {
				found = true
				matched[j] = true
				break
			}
		}
		if found {
			states[i] = Present
		} else {
			states[i] = Absent
		}
	}

	return states
}

// Summary returns a spoiler-free shareable result.
func (g *Game) Summary() common.ShareResult {
	title := fmt.Sprintf("🟩 Guess %d/6", len(g.Guesses))
	if g.IsLost() {
		title += " ❌"
	}

	lines := make([]string, len(g.Guesses))
	for i, gr := range g.Guesses {
		var b strings.Builder
		for _, s := range gr.States {
			switch s {
			case Correct:
				b.WriteString("🟩")
			case Present:
				b.WriteString("🟨")
			default:
				b.WriteString("⬛")
			}
		}
		lines[i] = b.String()
	}

	return common.ShareResult{
		Game:  "🟩 Guess",
		Title: title,
		Lines: lines,
	}
}

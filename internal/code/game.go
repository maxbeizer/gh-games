package code

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/maxbeizer/gh-games/internal/common"
)

// Color represents one of six peg colors.
type Color int

const (
	Red Color = iota
	Green
	Blue
	Yellow
	Purple
	Orange
	NumColors = 6
	CodeLen   = 4
)

// ColorName returns the display name of a color.
func ColorName(c Color) string {
	switch c {
	case Red:
		return "Red"
	case Green:
		return "Green"
	case Blue:
		return "Blue"
	case Yellow:
		return "Yellow"
	case Purple:
		return "Purple"
	case Orange:
		return "Orange"
	default:
		return "?"
	}
}

// ColorSymbol returns the emoji circle for a color.
func ColorSymbol(c Color) string {
	switch c {
	case Red:
		return "🔴"
	case Green:
		return "🟢"
	case Blue:
		return "🔵"
	case Yellow:
		return "🟡"
	case Purple:
		return "🟣"
	case Orange:
		return "🟠"
	default:
		return "?"
	}
}

// ColorLetter returns the single-letter abbreviation for a color.
func ColorLetter(c Color) string {
	switch c {
	case Red:
		return "R"
	case Green:
		return "G"
	case Blue:
		return "B"
	case Yellow:
		return "Y"
	case Purple:
		return "P"
	case Orange:
		return "O"
	default:
		return "?"
	}
}

// Feedback holds the result of evaluating a guess against the secret.
type Feedback struct {
	Exact     int // right color, right position
	Misplaced int // right color, wrong position
}

// Guess stores a submitted guess and its feedback.
type Guess struct {
	Code     [CodeLen]Color
	Feedback Feedback
}

// Game holds the state for a Code Breaker round.
type Game struct {
	Secret   [CodeLen]Color
	Guesses  []Guess
	MaxTurns int
}

// NewGame creates a new game with a random secret code.
func NewGame() *Game {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var secret [CodeLen]Color
	for i := range secret {
		secret[i] = Color(r.Intn(NumColors))
	}
	return NewGameWithSecret(secret)
}

// NewGameWithSecret creates a new game with a specific secret (for testing).
func NewGameWithSecret(secret [CodeLen]Color) *Game {
	return &Game{
		Secret:   secret,
		Guesses:  make([]Guess, 0, 10),
		MaxTurns: 10,
	}
}

// MakeGuess evaluates a guess, stores it, and returns the feedback.
func (g *Game) MakeGuess(code [CodeLen]Color) Feedback {
	fb := ComputeFeedback(g.Secret, code)
	g.Guesses = append(g.Guesses, Guess{Code: code, Feedback: fb})
	return fb
}

// IsWon returns true if the last guess was an exact match.
func (g *Game) IsWon() bool {
	if len(g.Guesses) == 0 {
		return false
	}
	return g.Guesses[len(g.Guesses)-1].Feedback.Exact == CodeLen
}

// IsLost returns true if all turns are used without winning.
func (g *Game) IsLost() bool {
	return len(g.Guesses) >= g.MaxTurns && !g.IsWon()
}

// IsOver returns true if the game has ended.
func (g *Game) IsOver() bool {
	return g.IsWon() || g.IsLost()
}

// ComputeFeedback is a pure function that compares a guess against a secret.
// It correctly handles duplicate colors using a two-pass algorithm.
func ComputeFeedback(secret, guess [CodeLen]Color) Feedback {
	var fb Feedback
	secretUsed := [CodeLen]bool{}
	guessUsed := [CodeLen]bool{}

	// First pass: exact matches
	for i := 0; i < CodeLen; i++ {
		if guess[i] == secret[i] {
			fb.Exact++
			secretUsed[i] = true
			guessUsed[i] = true
		}
	}

	// Second pass: misplaced (right color, wrong position)
	for i := 0; i < CodeLen; i++ {
		if guessUsed[i] {
			continue
		}
		for j := 0; j < CodeLen; j++ {
			if secretUsed[j] {
				continue
			}
			if guess[i] == secret[j] {
				fb.Misplaced++
				secretUsed[j] = true
				break
			}
		}
	}

	return fb
}

// Summary returns a spoiler-free shareable result.
func (g *Game) Summary() common.ShareResult {
	title := "🔐 Code: "
	if g.IsWon() {
		title += fmt.Sprintf("Cracked in %d", len(g.Guesses))
	} else {
		title += "Failed ❌"
	}

	lines := make([]string, len(g.Guesses))
	for i, guess := range g.Guesses {
		var b strings.Builder
		for j := 0; j < guess.Feedback.Exact; j++ {
			b.WriteString("●")
		}
		for j := 0; j < guess.Feedback.Misplaced; j++ {
			b.WriteString("○")
		}
		nothings := CodeLen - guess.Feedback.Exact - guess.Feedback.Misplaced
		for j := 0; j < nothings; j++ {
			b.WriteString("⚫")
		}
		lines[i] = b.String()
	}

	return common.ShareResult{
		Game:  "🔐 Code",
		Title: title,
		Lines: lines,
	}
}

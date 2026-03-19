package jumble

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/maxbeizer/gh-games/internal/common"
)

// roundLengths defines the word length for each of the 5 rounds (progressive difficulty).
var roundLengths = []int{4, 5, 6, 7, 8}

// Round represents a single round in the Jumble game.
type Round struct {
	Target    string        // the correct word (UPPERCASE)
	Scrambled string        // current scrambled display
	Hints     []int         // positions that have been revealed
	StartTime time.Time
	SolveTime time.Duration // how long it took to solve
	Solved    bool
	HintsUsed int
}

// Game holds the state for a complete Jumble game (5 rounds).
type Game struct {
	Rounds       []Round
	CurrentRound int
	TotalScore   int
}

// NewGame creates a new 5-round Jumble game with progressively longer words.
func NewGame() *Game {
	rounds := make([]Round, 5)
	for i, length := range roundLengths {
		word := RandomWordOfLength(length)
		rounds[i] = Round{
			Target:    word,
			Scrambled: ScrambleWord(word),
			Hints:     []int{},
			StartTime: time.Now(),
		}
	}
	return &Game{
		Rounds:       rounds,
		CurrentRound: 0,
	}
}

// CurrentRoundRef returns a pointer to the current round.
func (g *Game) CurrentRoundRef() *Round {
	if g.CurrentRound >= len(g.Rounds) {
		return nil
	}
	return &g.Rounds[g.CurrentRound]
}

// Guess checks the player's answer. Returns whether correct and points earned.
func (g *Game) Guess(word string) (correct bool, points int) {
	r := g.CurrentRoundRef()
	if r == nil || r.Solved {
		return false, 0
	}

	word = strings.ToUpper(strings.TrimSpace(word))
	if word != r.Target {
		return false, 0
	}

	r.Solved = true
	r.SolveTime = time.Since(r.StartTime)

	// Scoring: base points scale with word length
	base := len(r.Target) * 100

	// Speed bonus: up to 200 points if solved in under 10 seconds
	elapsed := r.SolveTime.Seconds()
	speedBonus := 0
	if elapsed < 10 {
		speedBonus = int((10 - elapsed) * 20)
	}

	// Hint penalty: 50 points per hint used
	hintPenalty := r.HintsUsed * 50

	points = base + speedBonus - hintPenalty
	if points < 10 {
		points = 10 // minimum points for solving
	}

	g.TotalScore += points

	return true, points
}

// Hint reveals one unrevealed letter in its correct position.
// Returns the letter, position (0-indexed), and any error.
func (g *Game) Hint() (rune, int, error) {
	r := g.CurrentRoundRef()
	if r == nil || r.Solved {
		return 0, 0, errors.New("no active round")
	}

	// Find positions not yet hinted
	revealed := make(map[int]bool)
	for _, pos := range r.Hints {
		revealed[pos] = true
	}

	var available []int
	for i := 0; i < len(r.Target); i++ {
		if !revealed[i] {
			available = append(available, i)
		}
	}

	if len(available) == 0 {
		return 0, 0, errors.New("all letters already revealed")
	}

	pos := available[rand.Intn(len(available))]
	letter := rune(r.Target[pos])
	r.Hints = append(r.Hints, pos)
	r.HintsUsed++

	// Update scrambled display: fix the hinted letter in its correct position
	r.Scrambled = applyHints(r.Target, r.Scrambled, r.Hints)

	return letter, pos, nil
}

// applyHints rebuilds the scrambled word with hinted positions fixed.
func applyHints(target, scrambled string, hints []int) string {
	revealed := make(map[int]bool)
	for _, pos := range hints {
		revealed[pos] = true
	}

	// Collect the letters that are NOT in hinted positions from the target
	var freeTargetLetters []byte
	for i := 0; i < len(target); i++ {
		if !revealed[i] {
			freeTargetLetters = append(freeTargetLetters, target[i])
		}
	}

	// Shuffle the free letters
	rand.Shuffle(len(freeTargetLetters), func(i, j int) {
		freeTargetLetters[i], freeTargetLetters[j] = freeTargetLetters[j], freeTargetLetters[i]
	})

	// Build the result
	result := make([]byte, len(target))
	freeIdx := 0
	for i := 0; i < len(target); i++ {
		if revealed[i] {
			result[i] = target[i]
		} else {
			result[i] = freeTargetLetters[freeIdx]
			freeIdx++
		}
	}

	return string(result)
}

// Shuffle re-scrambles the current round's word (preserving hints).
func (g *Game) Shuffle() {
	r := g.CurrentRoundRef()
	if r == nil || r.Solved {
		return
	}

	if len(r.Hints) > 0 {
		r.Scrambled = applyHints(r.Target, r.Scrambled, r.Hints)
	} else {
		r.Scrambled = ScrambleWord(r.Target)
	}
}

// NextRound advances to the next round. Returns false if game is complete.
func (g *Game) NextRound() bool {
	if g.CurrentRound >= len(g.Rounds)-1 {
		return false
	}
	g.CurrentRound++
	g.Rounds[g.CurrentRound].StartTime = time.Now()
	return true
}

// IsComplete returns true if all 5 rounds have been solved or attempted.
func (g *Game) IsComplete() bool {
	return g.CurrentRound >= len(g.Rounds)-1 && g.Rounds[g.CurrentRound].Solved
}

// ScrambleWord shuffles letters of a word, ensuring the result differs from the original.
func ScrambleWord(word string) string {
	runes := []rune(word)
	for attempts := 0; attempts < 100; attempts++ {
		shuffled := make([]rune, len(runes))
		copy(shuffled, runes)
		rand.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})
		if string(shuffled) != word {
			return string(shuffled)
		}
	}
	// Fallback: rotate by 1 position
	rotated := make([]rune, len(runes))
	copy(rotated, runes[1:])
	rotated[len(rotated)-1] = runes[0]
	return string(rotated)
}

// Summary returns a spoiler-free shareable result.
func (g *Game) Summary() common.ShareResult {
	solved := 0
	hints := 0
	for _, r := range g.Rounds {
		if r.Solved {
			solved++
		}
		hints += r.HintsUsed
	}

	return common.ShareResult{
		Game:  "🔀 Jumble",
		Title: "🔀 Jumble",
		Lines: []string{
			fmt.Sprintf("Score: %d", g.TotalScore),
			fmt.Sprintf("Rounds: %d/%d solved", solved, len(g.Rounds)),
			fmt.Sprintf("Hints used: %d", hints),
		},
	}
}

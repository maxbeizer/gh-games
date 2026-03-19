package hang

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/maxbeizer/gh-games/internal/common"
)

const MaxWrong = 6

// Game holds the state of a hangman game.
type Game struct {
	Target    string
	guessed   map[rune]bool
	correct   map[rune]bool
	wrong     map[rune]bool
	wrongList []rune
	wrongCnt  int
}

// NewGame picks a random word and returns a new game.
func NewGame() *Game {
	return NewGameWithWord(RandomWord())
}

// NewGameWithWord creates a game with a specific target word (uppercased).
func NewGameWithWord(word string) *Game {
	return &Game{
		Target:  strings.ToUpper(word),
		guessed: make(map[rune]bool),
		correct: make(map[rune]bool),
		wrong:   make(map[rune]bool),
	}
}

// GuessLetter processes a single letter guess.
func (g *Game) GuessLetter(r rune) (bool, error) {
	r = unicode.ToUpper(r)
	if r < 'A' || r > 'Z' {
		return false, fmt.Errorf("invalid character: %c", r)
	}
	if g.IsOver() {
		return false, fmt.Errorf("game is already over")
	}
	if g.guessed[r] {
		return false, fmt.Errorf("already guessed '%c'", r)
	}

	g.guessed[r] = true

	if strings.ContainsRune(g.Target, r) {
		g.correct[r] = true
		return true, nil
	}

	g.wrong[r] = true
	g.wrongList = append(g.wrongList, r)
	g.wrongCnt++
	return false, nil
}

// Display returns the word with blanks for unguessed letters (e.g. "_ _ O _ E").
func (g *Game) Display() string {
	var parts []string
	for _, r := range g.Target {
		if g.correct[r] {
			parts = append(parts, string(r))
		} else {
			parts = append(parts, "_")
		}
	}
	return strings.Join(parts, " ")
}

// IsWon returns true if all letters have been revealed.
func (g *Game) IsWon() bool {
	for _, r := range g.Target {
		if !g.correct[r] {
			return false
		}
	}
	return true
}

// IsLost returns true if the player has used all wrong guesses.
func (g *Game) IsLost() bool {
	return g.wrongCnt >= MaxWrong
}

// IsOver returns true if the game is won or lost.
func (g *Game) IsOver() bool {
	return g.IsWon() || g.IsLost()
}

// WrongCount returns the number of incorrect guesses.
func (g *Game) WrongCount() int {
	return g.wrongCnt
}

// WrongLetters returns the incorrectly guessed letters in order.
func (g *Game) WrongLetters() []rune {
	return g.wrongList
}

// IsGuessed reports whether a letter has been guessed.
func (g *Game) IsGuessed(r rune) bool {
	return g.guessed[unicode.ToUpper(r)]
}

// IsCorrect reports whether a guessed letter was correct.
func (g *Game) IsCorrect(r rune) bool {
	return g.correct[unicode.ToUpper(r)]
}

// IsWrongGuess reports whether a guessed letter was wrong.
func (g *Game) IsWrongGuess(r rune) bool {
	return g.wrong[unicode.ToUpper(r)]
}

// Summary returns a spoiler-free shareable result.
func (g *Game) Summary() common.ShareResult {
	var result string
	if g.IsWon() {
		result = "Won! ✓"
	} else {
		result = "Lost ❌"
	}

	return common.ShareResult{
		Game:  "☠️ Hang",
		Title: "☠️ Hang",
		Lines: []string{
			result,
			fmt.Sprintf("Wrong guesses: %d/%d", g.WrongCount(), MaxWrong),
		},
	}
}

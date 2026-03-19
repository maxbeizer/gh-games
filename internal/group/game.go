package group

import (
	"errors"
	"sort"
)

var (
	ErrNotEnoughSelected = errors.New("you must select exactly 4 words")
	ErrGameOver          = errors.New("game is over")
	ErrNoMatch           = errors.New("those 4 words do not belong to the same group")
)

// Game holds the full state for a single round of Group.
type Game struct {
	Categories     [4]Category
	RemainingWords []string
	SolvedGroups   []Category
	Mistakes       int
	MaxMistakes    int
	Selected       map[string]bool
}

// NewGame generates a fresh puzzle and returns a ready-to-play Game.
func NewGame() *Game {
	cats, words := GeneratePuzzle()
	return &Game{
		Categories:     cats,
		RemainingWords: words,
		SolvedGroups:   []Category{},
		Mistakes:       0,
		MaxMistakes:    4,
		Selected:       make(map[string]bool),
	}
}

// ToggleSelect toggles a word's selection state. Only words still in
// RemainingWords can be selected.
func (g *Game) ToggleSelect(word string) {
	if !g.isRemaining(word) {
		return
	}
	if g.Selected[word] {
		delete(g.Selected, word)
	} else {
		g.Selected[word] = true
	}
}

// ClearSelection deselects every word.
func (g *Game) ClearSelection() {
	g.Selected = make(map[string]bool)
}

// SelectedCount returns how many words are currently selected.
func (g *Game) SelectedCount() int {
	return len(g.Selected)
}

// Submit checks the current selection against the answer key.
// On a match it reveals the group; on a miss it records a mistake.
func (g *Game) Submit() (matched *Category, err error) {
	if g.IsOver() {
		return nil, ErrGameOver
	}
	if len(g.Selected) != 4 {
		return nil, ErrNotEnoughSelected
	}

	selected := make([]string, 0, 4)
	for w := range g.Selected {
		selected = append(selected, w)
	}
	sort.Strings(selected)

	for i := range g.Categories {
		cat := &g.Categories[i]
		if g.isSolved(cat) {
			continue
		}
		words := make([]string, len(cat.Words))
		copy(words, cat.Words)
		sort.Strings(words)

		if slicesEqual(selected, words) {
			g.SolvedGroups = append(g.SolvedGroups, *cat)
			g.removeWords(selected)
			g.ClearSelection()
			return cat, nil
		}
	}

	g.Mistakes++
	g.ClearSelection()
	return nil, ErrNoMatch
}

// IsWon returns true when all 4 groups have been solved.
func (g *Game) IsWon() bool {
	return len(g.SolvedGroups) == 4
}

// IsLost returns true when the player has used all allowed mistakes.
func (g *Game) IsLost() bool {
	return g.Mistakes >= g.MaxMistakes
}

// IsOver returns true when the game cannot continue.
func (g *Game) IsOver() bool {
	return g.IsWon() || g.IsLost()
}

// RemainingCategories returns the categories that have not been solved yet.
func (g *Game) RemainingCategories() []Category {
	var out []Category
	for i := range g.Categories {
		if !g.isSolved(&g.Categories[i]) {
			out = append(out, g.Categories[i])
		}
	}
	return out
}

// ── helpers ──────────────────────────────────────────────────────────

func (g *Game) isRemaining(word string) bool {
	for _, w := range g.RemainingWords {
		if w == word {
			return true
		}
	}
	return false
}

func (g *Game) isSolved(cat *Category) bool {
	for i := range g.SolvedGroups {
		if g.SolvedGroups[i].Name == cat.Name {
			return true
		}
	}
	return false
}

func (g *Game) removeWords(words []string) {
	remove := make(map[string]bool, len(words))
	for _, w := range words {
		remove[w] = true
	}
	filtered := g.RemainingWords[:0]
	for _, w := range g.RemainingWords {
		if !remove[w] {
			filtered = append(filtered, w)
		}
	}
	g.RemainingWords = filtered
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

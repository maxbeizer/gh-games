package cross

import (
	"fmt"
	"math/rand"
	"time"
	"unicode"

	"github.com/maxbeizer/gh-games/internal/common"
)

// Direction represents the cursor movement direction.
type Direction int

const (
	Across Direction = iota
	Down
)

// Game holds the full state of a crossword game session.
type Game struct {
	Puzzle    Puzzle
	Player    [5][5]rune // what the player has typed
	CurRow    int
	CurCol    int
	Dir       Direction
	StartTime time.Time
	Completed bool
	Checked   [5][5]bool // true = incorrect after check
}

// NewGame picks a random puzzle and returns a new game.
func NewGame() *Game {
	return NewGameWithPuzzle(Puzzles[rand.Intn(len(Puzzles))])
}

// NewGameWithPuzzle creates a game with a specific puzzle.
func NewGameWithPuzzle(p Puzzle) *Game {
	g := &Game{
		Puzzle:    p,
		StartTime: time.Now(),
	}
	// Place cursor on first non-black cell.
	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			if !p.Grid[r][c].Black {
				g.CurRow = r
				g.CurCol = c
				return g
			}
		}
	}
	return g
}

// SetLetter sets the letter at the current cursor position.
func (g *Game) SetLetter(r rune) {
	if g.Puzzle.Grid[g.CurRow][g.CurCol].Black {
		return
	}
	g.Player[g.CurRow][g.CurCol] = unicode.ToUpper(r)
}

// ClearLetter clears the current cell.
func (g *Game) ClearLetter() {
	if g.Puzzle.Grid[g.CurRow][g.CurCol].Black {
		return
	}
	g.Player[g.CurRow][g.CurCol] = 0
}

// Advance moves the cursor forward in the current direction, skipping black cells.
func (g *Game) Advance() {
	if g.Dir == Across {
		g.moveRight()
	} else {
		g.moveDown()
	}
}

// Retreat moves the cursor backward in the current direction, skipping black cells.
func (g *Game) Retreat() {
	if g.Dir == Across {
		g.moveLeft()
	} else {
		g.moveUp()
	}
}

// MoveCursor moves the cursor in the given arrow direction (0=up,1=down,2=left,3=right).
func (g *Game) MoveCursor(dir int) {
	switch dir {
	case 0:
		g.moveUp()
	case 1:
		g.moveDown()
	case 2:
		g.moveLeft()
	case 3:
		g.moveRight()
	}
}

func (g *Game) moveUp() {
	for r := g.CurRow - 1; r >= 0; r-- {
		if !g.Puzzle.Grid[r][g.CurCol].Black {
			g.CurRow = r
			return
		}
	}
}

func (g *Game) moveDown() {
	for r := g.CurRow + 1; r < 5; r++ {
		if !g.Puzzle.Grid[r][g.CurCol].Black {
			g.CurRow = r
			return
		}
	}
}

func (g *Game) moveLeft() {
	for c := g.CurCol - 1; c >= 0; c-- {
		if !g.Puzzle.Grid[g.CurRow][c].Black {
			g.CurCol = c
			return
		}
	}
}

func (g *Game) moveRight() {
	for c := g.CurCol + 1; c < 5; c++ {
		if !g.Puzzle.Grid[g.CurRow][c].Black {
			g.CurCol = c
			return
		}
	}
}

// ToggleDirection switches between across and down.
func (g *Game) ToggleDirection() {
	if g.Dir == Across {
		g.Dir = Down
	} else {
		g.Dir = Across
	}
}

// IsComplete returns true if all non-black cells are filled.
func (g *Game) IsComplete() bool {
	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			if !g.Puzzle.Grid[r][c].Black && g.Player[r][c] == 0 {
				return false
			}
		}
	}
	return true
}

// IsCorrect returns true if all letters match the answer.
func (g *Game) IsCorrect() bool {
	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			cell := g.Puzzle.Grid[r][c]
			if cell.Black {
				continue
			}
			if unicode.ToUpper(g.Player[r][c]) != unicode.ToUpper(cell.Letter) {
				return false
			}
		}
	}
	return true
}

// Check returns a grid where true means the cell is incorrect.
func (g *Game) Check() [5][5]bool {
	var result [5][5]bool
	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			cell := g.Puzzle.Grid[r][c]
			if cell.Black {
				continue
			}
			if g.Player[r][c] != 0 && unicode.ToUpper(g.Player[r][c]) != unicode.ToUpper(cell.Letter) {
				result[r][c] = true
			}
		}
	}
	return result
}

// Elapsed returns the time elapsed since the game started.
func (g *Game) Elapsed() time.Duration {
	return time.Since(g.StartTime)
}

// CellNumber returns the clue number for a cell, or 0 if none.
func (g *Game) CellNumber(r, c int) int {
	if g.Puzzle.Grid[r][c].Black {
		return 0
	}
	startsAcross := (c == 0 || g.Puzzle.Grid[r][c-1].Black) &&
		c+1 < 5 && !g.Puzzle.Grid[r][c+1].Black
	startsDown := (r == 0 || g.Puzzle.Grid[r-1][c].Black) &&
		r+1 < 5 && !g.Puzzle.Grid[r+1][c].Black
	if !startsAcross && !startsDown {
		return 0
	}
	// Count numbered cells up to this position.
	num := 0
	for rr := 0; rr < 5; rr++ {
		for cc := 0; cc < 5; cc++ {
			if g.Puzzle.Grid[rr][cc].Black {
				continue
			}
			sa := (cc == 0 || g.Puzzle.Grid[rr][cc-1].Black) &&
				cc+1 < 5 && !g.Puzzle.Grid[rr][cc+1].Black
			sd := (rr == 0 || g.Puzzle.Grid[rr-1][cc].Black) &&
				rr+1 < 5 && !g.Puzzle.Grid[rr+1][cc].Black
			if sa || sd {
				num++
			}
			if rr == r && cc == c {
				return num
			}
		}
	}
	return 0
}

// CurrentClue returns the clue matching the cursor position and direction.
func (g *Game) CurrentClue() *Clue {
	num := g.currentClueNumber()
	if num == 0 {
		return nil
	}
	dir := "across"
	if g.Dir == Down {
		dir = "down"
	}
	for i := range g.Puzzle.Clues {
		if g.Puzzle.Clues[i].Number == num && g.Puzzle.Clues[i].Direction == dir {
			return &g.Puzzle.Clues[i]
		}
	}
	return nil
}

func (g *Game) currentClueNumber() int {
	r, c := g.CurRow, g.CurCol
	if g.Dir == Across {
		// Walk left to find start of word.
		for c > 0 && !g.Puzzle.Grid[r][c-1].Black {
			c--
		}
	} else {
		// Walk up to find start of word.
		for r > 0 && !g.Puzzle.Grid[r-1][c].Black {
			r--
		}
	}
	return g.CellNumber(r, c)
}

// Summary returns a spoiler-free shareable result.
func (g *Game) Summary() common.ShareResult {
	var line string
	if g.IsCorrect() {
		d := g.Elapsed()
		mins := int(d.Minutes())
		secs := int(d.Seconds()) % 60
		line = fmt.Sprintf("Solved in %d:%02d", mins, secs)
	} else {
		line = "Incomplete"
	}

	return common.ShareResult{
		Game:  "📰 Cross",
		Title: "📰 Cross",
		Lines: []string{line},
	}
}

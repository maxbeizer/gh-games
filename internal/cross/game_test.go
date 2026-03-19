package cross

import (
	"testing"
)

func testPuzzle() Puzzle {
	// S W A M P
	// H O R S E
	// A ■ E ■ A
	// R I S K S
	// P E A R L
	return Puzzles[0]
}

func TestNewGameWithPuzzle(t *testing.T) {
	g := NewGameWithPuzzle(testPuzzle())
	if g.CurRow != 0 || g.CurCol != 0 {
		t.Errorf("expected cursor at (0,0), got (%d,%d)", g.CurRow, g.CurCol)
	}
	if g.Dir != Across {
		t.Error("expected initial direction Across")
	}
}

func TestSetAndClearLetter(t *testing.T) {
	g := NewGameWithPuzzle(testPuzzle())
	g.SetLetter('x')
	if g.Player[0][0] != 'X' {
		t.Errorf("expected 'X', got '%c'", g.Player[0][0])
	}
	g.ClearLetter()
	if g.Player[0][0] != 0 {
		t.Error("expected cleared cell")
	}
}

func TestSetLetterOnBlackCell(t *testing.T) {
	g := NewGameWithPuzzle(testPuzzle())
	// Move to black cell at (2,1)
	g.CurRow = 2
	g.CurCol = 1
	g.SetLetter('X')
	if g.Player[2][1] != 0 {
		t.Error("should not set letter on black cell")
	}
}

func TestMoveCursorSkipsBlack(t *testing.T) {
	g := NewGameWithPuzzle(testPuzzle())
	// Start at (2,0), move right should skip (2,1) which is black
	g.CurRow = 2
	g.CurCol = 0
	g.MoveCursor(3) // right
	if g.CurCol != 2 {
		t.Errorf("expected col 2 (skip black), got %d", g.CurCol)
	}
}

func TestMoveCursorSkipsBlackDown(t *testing.T) {
	g := NewGameWithPuzzle(testPuzzle())
	// Col 1: rows 0,1 are letters, row 2 is black, rows 3,4 are letters
	g.CurRow = 1
	g.CurCol = 1
	g.MoveCursor(1) // down
	if g.CurRow != 3 {
		t.Errorf("expected row 3 (skip black), got %d", g.CurRow)
	}
}

func TestMoveCursorUp(t *testing.T) {
	g := NewGameWithPuzzle(testPuzzle())
	g.CurRow = 3
	g.CurCol = 1
	g.MoveCursor(0) // up — should skip (2,1) black
	if g.CurRow != 1 {
		t.Errorf("expected row 1 (skip black), got %d", g.CurRow)
	}
}

func TestMoveCursorLeft(t *testing.T) {
	g := NewGameWithPuzzle(testPuzzle())
	g.CurRow = 2
	g.CurCol = 2
	g.MoveCursor(2) // left — should skip (2,1) black
	if g.CurCol != 0 {
		t.Errorf("expected col 0 (skip black), got %d", g.CurCol)
	}
}

func TestToggleDirection(t *testing.T) {
	g := NewGameWithPuzzle(testPuzzle())
	if g.Dir != Across {
		t.Error("expected initial Across")
	}
	g.ToggleDirection()
	if g.Dir != Down {
		t.Error("expected Down after toggle")
	}
	g.ToggleDirection()
	if g.Dir != Across {
		t.Error("expected Across after double toggle")
	}
}

func TestIsComplete(t *testing.T) {
	g := NewGameWithPuzzle(testPuzzle())
	if g.IsComplete() {
		t.Error("empty grid should not be complete")
	}
	// Fill all non-black cells with wrong letters
	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			if !g.Puzzle.Grid[r][c].Black {
				g.Player[r][c] = 'Z'
			}
		}
	}
	if !g.IsComplete() {
		t.Error("fully filled grid should be complete")
	}
}

func TestIsCorrect(t *testing.T) {
	g := NewGameWithPuzzle(testPuzzle())
	// Fill with correct answers
	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			if !g.Puzzle.Grid[r][c].Black {
				g.Player[r][c] = g.Puzzle.Grid[r][c].Letter
			}
		}
	}
	if !g.IsCorrect() {
		t.Error("correctly filled grid should be correct")
	}
	// Change one letter
	g.Player[0][0] = 'Z'
	if g.IsCorrect() {
		t.Error("grid with wrong letter should not be correct")
	}
}

func TestCheck(t *testing.T) {
	g := NewGameWithPuzzle(testPuzzle())
	// Fill correctly except one cell
	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			if !g.Puzzle.Grid[r][c].Black {
				g.Player[r][c] = g.Puzzle.Grid[r][c].Letter
			}
		}
	}
	g.Player[0][0] = 'Z' // wrong
	result := g.Check()
	if !result[0][0] {
		t.Error("expected (0,0) to be marked incorrect")
	}
	if result[0][1] {
		t.Error("expected (0,1) to be marked correct")
	}
}

func TestCheckEmptyCellsNotMarkedWrong(t *testing.T) {
	g := NewGameWithPuzzle(testPuzzle())
	result := g.Check()
	// Empty cells should not be marked incorrect
	if result[0][0] {
		t.Error("empty cell should not be marked incorrect")
	}
}

func TestCellNumber(t *testing.T) {
	g := NewGameWithPuzzle(testPuzzle())
	// (0,0) starts both across and down words — should be number 1
	n := g.CellNumber(0, 0)
	if n != 1 {
		t.Errorf("expected cell (0,0) number 1, got %d", n)
	}
	// Black cell should return 0
	n = g.CellNumber(2, 1)
	if n != 0 {
		t.Errorf("expected black cell number 0, got %d", n)
	}
}

func TestAdvanceAndRetreat(t *testing.T) {
	g := NewGameWithPuzzle(testPuzzle())
	g.CurRow = 0
	g.CurCol = 0
	g.Dir = Across
	g.Advance()
	if g.CurCol != 1 {
		t.Errorf("expected col 1 after advance, got %d", g.CurCol)
	}
	g.Retreat()
	if g.CurCol != 0 {
		t.Errorf("expected col 0 after retreat, got %d", g.CurCol)
	}
}

func TestCurrentClue(t *testing.T) {
	g := NewGameWithPuzzle(testPuzzle())
	g.CurRow = 0
	g.CurCol = 2
	g.Dir = Across
	clue := g.CurrentClue()
	if clue == nil {
		t.Fatal("expected a clue")
	}
	if clue.Number != 1 || clue.Direction != "across" {
		t.Errorf("expected 1-across, got %d-%s", clue.Number, clue.Direction)
	}
}

func TestNewGameRandom(t *testing.T) {
	g := NewGame()
	if g == nil {
		t.Fatal("NewGame returned nil")
	}
}

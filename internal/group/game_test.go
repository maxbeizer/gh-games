package group

import "testing"

// helper: build a game with known categories so tests are deterministic.
func testGame() *Game {
	cats := [4]Category{
		{Name: "Fruits", Words: []string{"APPLE", "BANANA", "GRAPE", "MANGO"}, Difficulty: Easy},
		{Name: "Planets", Words: []string{"MARS", "VENUS", "SATURN", "JUPITER"}, Difficulty: Medium},
		{Name: "Colors", Words: []string{"RED", "BLUE", "YELLOW", "GREEN"}, Difficulty: Hard},
		{Name: "Seasons", Words: []string{"SPRING", "SUMMER", "FALL", "WINTER"}, Difficulty: Expert},
	}
	words := make([]string, 0, 16)
	for _, c := range cats {
		words = append(words, c.Words...)
	}
	return &Game{
		Categories:     cats,
		RemainingWords: words,
		SolvedGroups:   []Category{},
		Mistakes:       0,
		MaxMistakes:    4,
		Selected:       make(map[string]bool),
	}
}

func TestNewGame(t *testing.T) {
	g := NewGame()
	if len(g.RemainingWords) != 16 {
		t.Fatalf("expected 16 remaining words, got %d", len(g.RemainingWords))
	}
	if len(g.SolvedGroups) != 0 {
		t.Fatalf("expected 0 solved groups, got %d", len(g.SolvedGroups))
	}
	if g.Mistakes != 0 {
		t.Fatalf("expected 0 mistakes, got %d", g.Mistakes)
	}
	if g.MaxMistakes != 4 {
		t.Fatalf("expected MaxMistakes=4, got %d", g.MaxMistakes)
	}
}

func TestToggleSelect(t *testing.T) {
	g := testGame()

	g.ToggleSelect("APPLE")
	if !g.Selected["APPLE"] {
		t.Fatal("APPLE should be selected")
	}
	if g.SelectedCount() != 1 {
		t.Fatalf("expected 1 selected, got %d", g.SelectedCount())
	}

	// toggle off
	g.ToggleSelect("APPLE")
	if g.Selected["APPLE"] {
		t.Fatal("APPLE should be deselected")
	}
	if g.SelectedCount() != 0 {
		t.Fatalf("expected 0 selected, got %d", g.SelectedCount())
	}
}

func TestToggleSelectIgnoresNonRemaining(t *testing.T) {
	g := testGame()
	g.ToggleSelect("NOTAWORD")
	if g.SelectedCount() != 0 {
		t.Fatal("should not be able to select a word not in RemainingWords")
	}
}

func TestClearSelection(t *testing.T) {
	g := testGame()
	g.ToggleSelect("APPLE")
	g.ToggleSelect("BANANA")
	g.ClearSelection()
	if g.SelectedCount() != 0 {
		t.Fatalf("expected 0 after clear, got %d", g.SelectedCount())
	}
}

func TestSubmitCorrectGroup(t *testing.T) {
	g := testGame()
	for _, w := range []string{"APPLE", "BANANA", "GRAPE", "MANGO"} {
		g.ToggleSelect(w)
	}

	matched, err := g.Submit()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if matched.Name != "Fruits" {
		t.Fatalf("expected Fruits, got %s", matched.Name)
	}
	if len(g.RemainingWords) != 12 {
		t.Fatalf("expected 12 remaining, got %d", len(g.RemainingWords))
	}
	if len(g.SolvedGroups) != 1 {
		t.Fatalf("expected 1 solved, got %d", len(g.SolvedGroups))
	}
	if g.Mistakes != 0 {
		t.Fatal("mistakes should still be 0")
	}
	if g.SelectedCount() != 0 {
		t.Fatal("selection should be cleared after submit")
	}
}

func TestSubmitWrongGroup(t *testing.T) {
	g := testGame()
	for _, w := range []string{"APPLE", "MARS", "RED", "SPRING"} {
		g.ToggleSelect(w)
	}

	matched, err := g.Submit()
	if err != ErrNoMatch {
		t.Fatalf("expected ErrNoMatch, got %v", err)
	}
	if matched != nil {
		t.Fatal("matched should be nil on wrong guess")
	}
	if g.Mistakes != 1 {
		t.Fatalf("expected 1 mistake, got %d", g.Mistakes)
	}
	if g.SelectedCount() != 0 {
		t.Fatal("selection should be cleared after wrong guess")
	}
}

func TestSubmitNotExactlyFour(t *testing.T) {
	g := testGame()
	g.ToggleSelect("APPLE")
	g.ToggleSelect("BANANA")

	_, err := g.Submit()
	if err != ErrNotEnoughSelected {
		t.Fatalf("expected ErrNotEnoughSelected, got %v", err)
	}
}

func TestWinCondition(t *testing.T) {
	g := testGame()

	for _, cat := range g.Categories {
		for _, w := range cat.Words {
			g.ToggleSelect(w)
		}
		_, err := g.Submit()
		if err != nil {
			t.Fatalf("unexpected error solving %s: %v", cat.Name, err)
		}
	}

	if !g.IsWon() {
		t.Fatal("game should be won")
	}
	if !g.IsOver() {
		t.Fatal("game should be over")
	}
	if len(g.RemainingWords) != 0 {
		t.Fatalf("expected 0 remaining, got %d", len(g.RemainingWords))
	}
}

func TestLossCondition(t *testing.T) {
	g := testGame()
	wrong := []string{"APPLE", "MARS", "RED", "SPRING"}

	for i := 0; i < 4; i++ {
		for _, w := range wrong {
			g.ToggleSelect(w)
		}
		g.Submit()
	}

	if !g.IsLost() {
		t.Fatal("game should be lost")
	}
	if !g.IsOver() {
		t.Fatal("game should be over")
	}
	if g.Mistakes != 4 {
		t.Fatalf("expected 4 mistakes, got %d", g.Mistakes)
	}
}

func TestCannotSubmitAfterGameOver(t *testing.T) {
	g := testGame()
	g.Mistakes = 4 // force game over

	g.ToggleSelect("APPLE")
	g.ToggleSelect("BANANA")
	g.ToggleSelect("GRAPE")
	g.ToggleSelect("MANGO")

	_, err := g.Submit()
	if err != ErrGameOver {
		t.Fatalf("expected ErrGameOver, got %v", err)
	}
}

func TestCannotSelectSolvedWords(t *testing.T) {
	g := testGame()
	// solve Fruits
	for _, w := range []string{"APPLE", "BANANA", "GRAPE", "MANGO"} {
		g.ToggleSelect(w)
	}
	g.Submit()

	// try selecting a solved word
	g.ToggleSelect("APPLE")
	if g.SelectedCount() != 0 {
		t.Fatal("should not select a word that was already removed")
	}
}

func TestRemainingCategories(t *testing.T) {
	g := testGame()
	if len(g.RemainingCategories()) != 4 {
		t.Fatalf("expected 4 remaining categories, got %d", len(g.RemainingCategories()))
	}

	// solve one
	for _, w := range []string{"APPLE", "BANANA", "GRAPE", "MANGO"} {
		g.ToggleSelect(w)
	}
	g.Submit()

	rem := g.RemainingCategories()
	if len(rem) != 3 {
		t.Fatalf("expected 3 remaining categories, got %d", len(rem))
	}
	for _, c := range rem {
		if c.Name == "Fruits" {
			t.Fatal("Fruits should not be in remaining categories")
		}
	}
}

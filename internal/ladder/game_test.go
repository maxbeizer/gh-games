package ladder

import (
	"testing"
	"time"
)

func TestDiffByOne(t *testing.T) {
	tests := []struct {
		a, b string
		want bool
	}{
		{"COLD", "CORD", true},
		{"COLD", "BOLD", true},
		{"COLD", "WARM", false},
		{"COLD", "COLD", false},
		{"COLD", "CO", false},
		{"CARE", "CORE", true},
		{"CARE", "DARE", true},
		{"CARE", "CARS", true},
		{"ABCD", "WXYZ", false},
	}

	for _, tt := range tests {
		got := DiffByOne(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("DiffByOne(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestFindShortestPath(t *testing.T) {
	wordSet := BuildWordSet()

	path := FindShortestPath("COLD", "CORD", wordSet)
	if path == nil {
		t.Fatal("expected path from COLD to CORD, got nil")
	}
	if len(path) < 2 {
		t.Errorf("path too short: %v", path)
	}
	if path[0] != "COLD" {
		t.Errorf("path should start with COLD, got %s", path[0])
	}
	if path[len(path)-1] != "CORD" {
		t.Errorf("path should end with CORD, got %s", path[len(path)-1])
	}

	for i := 1; i < len(path); i++ {
		if !DiffByOne(path[i-1], path[i]) {
			t.Errorf("invalid step in path: %s → %s", path[i-1], path[i])
		}
	}
}

func TestFindShortestPathSameWord(t *testing.T) {
	wordSet := BuildWordSet()
	path := FindShortestPath("COLD", "COLD", wordSet)
	if path == nil || len(path) != 1 {
		t.Errorf("same word path should be [COLD], got %v", path)
	}
}

func TestFindShortestPathNoPath(t *testing.T) {
	wordSet := map[string]bool{"ABCD": true, "WXYZ": true}
	path := FindShortestPath("ABCD", "WXYZ", wordSet)
	if path != nil {
		t.Errorf("expected no path, got %v", path)
	}
}

func TestGameStep(t *testing.T) {
	g := NewGameWithWords("COLD", "CORD")

	err := g.Step("CORD")
	if err != nil {
		t.Fatalf("valid step COLD→CORD should work: %v", err)
	}
	if g.Current != "CORD" {
		t.Errorf("current should be CORD, got %s", g.Current)
	}
	if g.StepCount() != 1 {
		t.Errorf("step count should be 1, got %d", g.StepCount())
	}
}

func TestGameStepInvalidWord(t *testing.T) {
	g := NewGameWithWords("COLD", "CORD")

	err := g.Step("XYZW")
	if err == nil {
		t.Error("step with invalid word should return error")
	}
}

func TestGameStepTooManyLettersChanged(t *testing.T) {
	g := NewGameWithWords("COLD", "WARM")

	err := g.Step("WARM")
	if err == nil {
		t.Error("changing multiple letters should return error")
	}
}

func TestGameStepSameWord(t *testing.T) {
	g := NewGameWithWords("COLD", "CORD")

	err := g.Step("COLD")
	if err == nil {
		t.Error("step with same word should return error")
	}
}

func TestGameStepAlreadyUsed(t *testing.T) {
	g := NewGameWithWords("COLD", "CORD")

	_ = g.Step("CORD")
	err := g.Step("COLD")
	if err == nil {
		t.Error("step with already-used word should return error")
	}
}

func TestGameWin(t *testing.T) {
	g := NewGameWithWords("COLD", "CORD")

	if g.IsWon() {
		t.Error("game should not be won at start")
	}

	_ = g.Step("CORD")

	if !g.IsWon() {
		t.Error("game should be won after reaching end word")
	}
}

func TestGameIsOptimal(t *testing.T) {
	g := NewGameWithWords("COLD", "CORD")
	_ = g.Step("CORD")

	if !g.IsOptimal() {
		t.Errorf("should be optimal (steps=%d, optimal=%d)", g.StepCount(), g.Optimal)
	}
}

func TestNewGameDoesNotHang(t *testing.T) {
	done := make(chan bool, 1)
	go func() {
		g := NewGame()
		if g == nil {
			t.Error("NewGame returned nil")
		}
		done <- true
	}()

	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()

	select {
	case <-done:
		// Success
	case <-timer.C:
		t.Fatal("NewGame hung for more than 10 seconds")
	}
}

func TestBuildWordSet(t *testing.T) {
	ws := BuildWordSet()
	if len(ws) < 100 {
		t.Errorf("word set too small: %d words", len(ws))
	}
	for _, w := range []string{"COLD", "WARM", "LOVE", "HATE"} {
		if !ws[w] {
			t.Errorf("word set should contain %s", w)
		}
	}
}

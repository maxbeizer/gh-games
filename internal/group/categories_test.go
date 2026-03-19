package group

import (
	"strings"
	"testing"
)

func TestAtLeast60Categories(t *testing.T) {
	if len(AllCategories) < 60 {
		t.Errorf("expected at least 60 categories, got %d", len(AllCategories))
	}
}

func TestAtLeast15PerDifficulty(t *testing.T) {
	counts := map[Difficulty]int{}
	for _, c := range AllCategories {
		counts[c.Difficulty]++
	}
	for _, d := range []Difficulty{Easy, Medium, Hard, Expert} {
		if counts[d] < 15 {
			t.Errorf("difficulty %d has %d categories, want at least 15", d, counts[d])
		}
	}
}

func TestEachCategoryHas4Words(t *testing.T) {
	for _, c := range AllCategories {
		if len(c.Words) != 4 {
			t.Errorf("category %q has %d words, want 4", c.Name, len(c.Words))
		}
	}
}

func TestAllWordsUppercase(t *testing.T) {
	for _, c := range AllCategories {
		for _, w := range c.Words {
			if w != strings.ToUpper(w) {
				t.Errorf("category %q: word %q is not uppercase", c.Name, w)
			}
		}
	}
}

func TestAllWordsSingleWord(t *testing.T) {
	for _, c := range AllCategories {
		for _, w := range c.Words {
			if strings.Contains(w, " ") {
				t.Errorf("category %q: word %q contains a space", c.Name, w)
			}
		}
	}
}

func TestGeneratePuzzleReturns4CategoriesAnd16Words(t *testing.T) {
	cats, words := GeneratePuzzle()
	if len(cats) != 4 {
		t.Fatalf("expected 4 categories, got %d", len(cats))
	}
	if len(words) != 16 {
		t.Fatalf("expected 16 words, got %d", len(words))
	}
}

func TestGeneratePuzzleOnePerDifficulty(t *testing.T) {
	cats, _ := GeneratePuzzle()
	seen := map[Difficulty]bool{}
	for _, c := range cats {
		if seen[c.Difficulty] {
			t.Errorf("duplicate difficulty %d", c.Difficulty)
		}
		seen[c.Difficulty] = true
	}
	for _, d := range []Difficulty{Easy, Medium, Hard, Expert} {
		if !seen[d] {
			t.Errorf("missing difficulty %d", d)
		}
	}
}

func TestGeneratePuzzleNoDuplicateWords(t *testing.T) {
	for i := 0; i < 50; i++ {
		_, words := GeneratePuzzle()
		seen := map[string]bool{}
		for _, w := range words {
			if seen[w] {
				t.Fatalf("duplicate word %q in puzzle", w)
			}
			seen[w] = true
		}
	}
}

func TestGeneratePuzzleIsRandom(t *testing.T) {
	first, _ := GeneratePuzzle()
	different := false
	for i := 0; i < 20; i++ {
		other, _ := GeneratePuzzle()
		for d := 0; d < 4; d++ {
			if first[d].Name != other[d].Name {
				different = true
				break
			}
		}
		if different {
			break
		}
	}
	if !different {
		t.Error("GeneratePuzzle returned identical puzzles 20 times in a row")
	}
}

func TestGeneratePuzzleWordsMatchCategories(t *testing.T) {
	cats, words := GeneratePuzzle()
	expected := map[string]bool{}
	for _, c := range cats {
		for _, w := range c.Words {
			expected[w] = true
		}
	}
	for _, w := range words {
		if !expected[w] {
			t.Errorf("word %q in shuffled list not found in any category", w)
		}
	}
	if len(expected) != len(words) {
		t.Errorf("word count mismatch: categories have %d unique words, shuffled list has %d", len(expected), len(words))
	}
}

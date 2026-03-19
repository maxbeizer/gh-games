package ladder

import (
	"fmt"
	"math/rand"
	"strings"
)

// Pre-computed solvable start/end pairs with known optimal path lengths.
// These avoid expensive BFS at game startup.
var solvablePairs = []struct {
	Start, End string
	Optimal    int
}{
	{"COLD", "WARM", 5},
	{"LOVE", "HATE", 5},
	{"HEAD", "TAIL", 5},
	{"HIDE", "SEEK", 6},
	{"RAIN", "SNOW", 5},
	{"DARK", "DAWN", 3},
	{"FAST", "SLOW", 5},
	{"LEAD", "GOLD", 4},
	{"LAND", "SAIL", 3},
	{"MALE", "FAME", 3},
	{"MAKE", "MARE", 3},
	{"COLD", "CORD", 1},
	{"BORE", "BORN", 1},
	{"WINE", "DINE", 1},
	{"MINE", "MILE", 1},
	{"RICE", "RACE", 1},
	{"MARE", "CARE", 1},
	{"BONE", "CONE", 1},
	{"CORN", "CORE", 1},
	{"CORD", "WORD", 1},
}

// Game represents the state of a Word Ladder game.
type Game struct {
	Start   string
	End     string
	Current string
	Steps   []string
	Optimal int
	WordSet map[string]bool
	used    map[string]bool
}

// NewGame creates a new game with a random solvable start/end pair
// that has a path of 3-7 steps.
func NewGame() *Game {
	wordSet := BuildWordSet()

	var candidates []validPair

	for _, p := range solvablePairs {
		if len(p.Start) != 4 || len(p.End) != 4 {
			continue
		}
		if !wordSet[p.Start] || !wordSet[p.End] {
			continue
		}
		path := FindShortestPath(p.Start, p.End, wordSet)
		if path == nil {
			continue
		}
		optimal := len(path) - 1
		if optimal >= 3 && optimal <= 7 {
			candidates = append(candidates, validPair{p.Start, p.End, optimal})
		}
	}

	// If no pre-computed pairs work, dynamically find some
	if len(candidates) == 0 {
		candidates = findSolvablePairs(wordSet, 10)
	}

	if len(candidates) == 0 {
		// Absolute fallback
		return NewGameWithWords("COLD", "CORD")
	}

	pick := candidates[rand.Intn(len(candidates))]
	return &Game{
		Start:   pick.start,
		End:     pick.end,
		Current: pick.start,
		Steps:   []string{pick.start},
		Optimal: pick.optimal,
		WordSet: wordSet,
		used:    map[string]bool{pick.start: true},
	}
}

// findSolvablePairs tries random word pairs to find ones with paths of 3-7 steps.
func findSolvablePairs(wordSet map[string]bool, limit int) []validPair {
	words := make([]string, 0, len(wordSet))
	for w := range wordSet {
		words = append(words, w)
	}

	var results []validPair

	attempts := 0
	maxAttempts := 500
	for attempts < maxAttempts && len(results) < limit {
		a := words[rand.Intn(len(words))]
		b := words[rand.Intn(len(words))]
		if a == b {
			attempts++
			continue
		}
		path := FindShortestPath(a, b, wordSet)
		if path != nil {
			optimal := len(path) - 1
			if optimal >= 3 && optimal <= 7 {
				results = append(results, validPair{a, b, optimal})
			}
		}
		attempts++
	}
	return results
}

type validPair struct {
	start, end string
	optimal    int
}

// NewGameWithWords creates a game with specific start and end words.
func NewGameWithWords(start, end string) *Game {
	start = strings.ToUpper(start)
	end = strings.ToUpper(end)
	wordSet := BuildWordSet()

	wordSet[start] = true
	wordSet[end] = true

	optimal := 0
	path := FindShortestPath(start, end, wordSet)
	if path != nil {
		optimal = len(path) - 1
	}

	return &Game{
		Start:   start,
		End:     end,
		Current: start,
		Steps:   []string{start},
		Optimal: optimal,
		WordSet: wordSet,
		used:    map[string]bool{start: true},
	}
}

// Step attempts to add a word to the ladder.
func (g *Game) Step(word string) error {
	word = strings.ToUpper(word)

	if g.IsWon() {
		return fmt.Errorf("game is already won")
	}

	if len(word) != len(g.Start) {
		return fmt.Errorf("word must be %d letters", len(g.Start))
	}

	if !g.WordSet[word] {
		return fmt.Errorf("%s is not a valid word", word)
	}

	if !DiffByOne(g.Current, word) {
		diff := countDiff(g.Current, word)
		if diff == 0 {
			return fmt.Errorf("must change exactly one letter (same word)")
		}
		return fmt.Errorf("must change exactly one letter (%d changed)", diff)
	}

	if g.used[word] {
		return fmt.Errorf("%s was already used", word)
	}

	g.Steps = append(g.Steps, word)
	g.Current = word
	g.used[word] = true
	return nil
}

// IsWon returns true if the current word matches the end word.
func (g *Game) IsWon() bool {
	return g.Current == g.End
}

// StepCount returns the number of steps taken (not counting the start word).
func (g *Game) StepCount() int {
	return len(g.Steps) - 1
}

// IsOptimal returns true if the game was completed in the optimal number of steps.
func (g *Game) IsOptimal() bool {
	return g.IsWon() && g.Optimal > 0 && g.StepCount() == g.Optimal
}

// DiffByOne returns true if two words differ by exactly one letter.
func DiffByOne(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	return countDiff(a, b) == 1
}

func countDiff(a, b string) int {
	diff := 0
	for i := range a {
		if a[i] != b[i] {
			diff++
		}
	}
	return diff
}

// FindShortestPath uses BFS to find the shortest path between two words.
// Returns nil if no path exists.
func FindShortestPath(start, end string, wordSet map[string]bool) []string {
	if start == end {
		return []string{start}
	}
	if !wordSet[start] || !wordSet[end] {
		return nil
	}

	wordLen := len(start)
	words := make([]string, 0)
	for w := range wordSet {
		if len(w) == wordLen {
			words = append(words, w)
		}
	}

	type node struct {
		word string
		path []string
	}

	visited := map[string]bool{start: true}
	queue := []node{{word: start, path: []string{start}}}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for _, w := range words {
			if visited[w] {
				continue
			}
			if DiffByOne(curr.word, w) {
				newPath := make([]string, len(curr.path)+1)
				copy(newPath, curr.path)
				newPath[len(curr.path)] = w
				if w == end {
					return newPath
				}
				visited[w] = true
				queue = append(queue, node{word: w, path: newPath})
			}
		}
	}

	return nil
}

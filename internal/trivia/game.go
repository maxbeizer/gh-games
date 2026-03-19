package trivia

import "math/rand"

// QuestionResult stores the player's answer alongside the question.
type QuestionResult struct {
	Question Question
	Chosen   int
	Correct  bool
}

// Game manages a single trivia round of 10 questions.
type Game struct {
	questions []Question
	results   []QuestionResult
	current   int
	score     int
}

// NewGame randomly selects 10 questions from the bank and shuffles them.
func NewGame() *Game {
	perm := rand.Perm(len(AllQuestions))
	n := 10
	if len(AllQuestions) < n {
		n = len(AllQuestions)
	}
	selected := make([]Question, n)
	for i := 0; i < n; i++ {
		selected[i] = AllQuestions[perm[i]]
	}
	return &Game{questions: selected}
}

// NewGameWithQuestions creates a game with a specific set of questions (for testing).
func NewGameWithQuestions(questions []Question) *Game {
	return &Game{questions: questions}
}

// CurrentQuestion returns the current question, or nil if the game is complete.
func (g *Game) CurrentQuestion() *Question {
	if g.current >= len(g.questions) {
		return nil
	}
	return &g.questions[g.current]
}

// Answer records the player's choice and advances to the next question.
// Returns true if the answer was correct.
func (g *Game) Answer(choice int) bool {
	if g.IsComplete() {
		return false
	}
	q := g.questions[g.current]
	correct := choice == q.Answer
	if correct {
		g.score++
	}
	g.results = append(g.results, QuestionResult{
		Question: q,
		Chosen:   choice,
		Correct:  correct,
	})
	g.current++
	return correct
}

// IsComplete returns true when all questions have been answered.
func (g *Game) IsComplete() bool {
	return g.current >= len(g.questions)
}

// Score returns the number of correct answers.
func (g *Game) Score() int {
	return g.score
}

// Total returns the total number of questions in this round.
func (g *Game) Total() int {
	return len(g.questions)
}

// Current returns the 0-based index of the current question.
func (g *Game) Current() int {
	return g.current
}

// Results returns the detailed results for each answered question.
func (g *Game) Results() []QuestionResult {
	return g.results
}

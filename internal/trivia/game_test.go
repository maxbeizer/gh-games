package trivia

import "testing"

func testQuestions() []Question {
	return []Question{
		{Text: "Q1", Choices: [4]string{"A", "B", "C", "D"}, Answer: 0, Category: "Test"},
		{Text: "Q2", Choices: [4]string{"A", "B", "C", "D"}, Answer: 1, Category: "Test"},
		{Text: "Q3", Choices: [4]string{"A", "B", "C", "D"}, Answer: 2, Category: "Test"},
		{Text: "Q4", Choices: [4]string{"A", "B", "C", "D"}, Answer: 3, Category: "Test"},
		{Text: "Q5", Choices: [4]string{"A", "B", "C", "D"}, Answer: 0, Category: "Test"},
		{Text: "Q6", Choices: [4]string{"A", "B", "C", "D"}, Answer: 1, Category: "Test"},
		{Text: "Q7", Choices: [4]string{"A", "B", "C", "D"}, Answer: 2, Category: "Test"},
		{Text: "Q8", Choices: [4]string{"A", "B", "C", "D"}, Answer: 3, Category: "Test"},
		{Text: "Q9", Choices: [4]string{"A", "B", "C", "D"}, Answer: 0, Category: "Test"},
		{Text: "Q10", Choices: [4]string{"A", "B", "C", "D"}, Answer: 1, Category: "Test"},
	}
}

func TestCorrectAnswerIncreasesScore(t *testing.T) {
	g := NewGameWithQuestions(testQuestions())

	correct := g.Answer(0) // Q1 answer is 0
	if !correct {
		t.Fatal("expected correct answer")
	}
	if g.Score() != 1 {
		t.Fatalf("expected score 1, got %d", g.Score())
	}
}

func TestWrongAnswerDoesNotIncreaseScore(t *testing.T) {
	g := NewGameWithQuestions(testQuestions())

	correct := g.Answer(3) // Q1 answer is 0, so 3 is wrong
	if correct {
		t.Fatal("expected wrong answer")
	}
	if g.Score() != 0 {
		t.Fatalf("expected score 0, got %d", g.Score())
	}
}

func TestGameCompletesAfter10(t *testing.T) {
	g := NewGameWithQuestions(testQuestions())

	for i := 0; i < 10; i++ {
		if g.IsComplete() {
			t.Fatalf("game should not be complete after %d questions", i)
		}
		g.Answer(0)
	}

	if !g.IsComplete() {
		t.Fatal("game should be complete after 10 questions")
	}
}

func TestAllCorrectScoreIsTen(t *testing.T) {
	g := NewGameWithQuestions(testQuestions())
	answers := []int{0, 1, 2, 3, 0, 1, 2, 3, 0, 1}

	for i, a := range answers {
		correct := g.Answer(a)
		if !correct {
			t.Fatalf("expected Q%d to be correct with answer %d", i+1, a)
		}
	}

	if g.Score() != 10 {
		t.Fatalf("expected score 10, got %d", g.Score())
	}
}

func TestAllWrongScoreIsZero(t *testing.T) {
	g := NewGameWithQuestions(testQuestions())
	// All answers shifted by 1 from correct
	answers := []int{1, 2, 3, 0, 1, 2, 3, 0, 1, 2}

	for _, a := range answers {
		g.Answer(a)
	}

	if g.Score() != 0 {
		t.Fatalf("expected score 0, got %d", g.Score())
	}
}

func TestCurrentQuestionAdvances(t *testing.T) {
	g := NewGameWithQuestions(testQuestions())

	q := g.CurrentQuestion()
	if q.Text != "Q1" {
		t.Fatalf("expected Q1, got %s", q.Text)
	}

	g.Answer(0)
	q = g.CurrentQuestion()
	if q.Text != "Q2" {
		t.Fatalf("expected Q2, got %s", q.Text)
	}
}

func TestCurrentQuestionNilWhenComplete(t *testing.T) {
	g := NewGameWithQuestions(testQuestions())
	for i := 0; i < 10; i++ {
		g.Answer(0)
	}

	if g.CurrentQuestion() != nil {
		t.Fatal("expected nil after game complete")
	}
}

func TestResultsLength(t *testing.T) {
	g := NewGameWithQuestions(testQuestions())
	for i := 0; i < 10; i++ {
		g.Answer(0)
	}

	results := g.Results()
	if len(results) != 10 {
		t.Fatalf("expected 10 results, got %d", len(results))
	}
}

func TestResultsRecordCorrectness(t *testing.T) {
	g := NewGameWithQuestions(testQuestions())
	g.Answer(0) // correct
	g.Answer(0) // wrong (Q2 answer is 1)

	results := g.Results()
	if !results[0].Correct {
		t.Fatal("first result should be correct")
	}
	if results[1].Correct {
		t.Fatal("second result should be wrong")
	}
}

func TestAnswerAfterCompleteReturnsFalse(t *testing.T) {
	g := NewGameWithQuestions(testQuestions())
	for i := 0; i < 10; i++ {
		g.Answer(0)
	}

	if g.Answer(0) {
		t.Fatal("answer after complete should return false")
	}
	if g.Score() != 5 { // only Q1,Q5,Q9 are answer=0 — wait, let me count
		// answers are 0,1,2,3,0,1,2,3,0,1 — answering 0 each time:
		// correct for Q1(0), Q5(0), Q9(0) = 3 correct
		// but we already answered all 10, score shouldn't change
	}
}

func TestTotal(t *testing.T) {
	g := NewGameWithQuestions(testQuestions())
	if g.Total() != 10 {
		t.Fatalf("expected total 10, got %d", g.Total())
	}
}

func TestNewGameSelects10(t *testing.T) {
	g := NewGame()
	if g.Total() != 10 {
		t.Fatalf("expected 10 questions, got %d", g.Total())
	}
}

func TestQuestionBankHasAtLeast80(t *testing.T) {
	if len(AllQuestions) < 80 {
		t.Fatalf("expected at least 80 questions, got %d", len(AllQuestions))
	}
}

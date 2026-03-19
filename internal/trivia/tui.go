package trivia

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maxbeizer/gh-games/internal/common"
)

// TUI-specific styles
var (
	categoryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#F5C518")).
			Bold(true).
			Padding(0, 1)

	questionStyle = lipgloss.NewStyle().
			Bold(true).
			Width(60)

	choiceNormal = lipgloss.NewStyle().
			Padding(0, 1)

	choiceSelected = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#F5C518")).
			Bold(true).
			Padding(0, 1)

	choiceCorrectReveal = lipgloss.NewStyle().
				Foreground(common.ColorWhite).
				Background(common.ColorGreen).
				Bold(true).
				Padding(0, 1)

	choiceWrongReveal = lipgloss.NewStyle().
				Foreground(common.ColorWhite).
				Background(common.ColorRed).
				Bold(true).
				Padding(0, 1)

	scoreStyle = lipgloss.NewStyle().
			Faint(true)

	gradeStyle = lipgloss.NewStyle().
			Bold(true).
			MarginTop(1)
)

type phase int

const (
	phaseAnswering phase = iota
	phaseRevealed
	phaseDone
)

// Model is the Bubbletea model for the Trivia game.
type Model struct {
	game        *Game
	cursor      int   // 0-3 selected choice
	phase       phase // current display phase
	lastCorrect bool
	sharePrompt *common.SharePrompt
}

// NewModel creates a new Trivia TUI model.
func NewModel() Model {
	return Model{
		game: NewGame(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles key input.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.phase {
		case phaseAnswering:
			return m.updateAnswering(msg)
		case phaseRevealed:
			return m.updateRevealed(msg)
		case phaseDone:
			return m.updateDone(msg)
		}
	}
	return m, nil
}

func (m Model) updateAnswering(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return m, tea.Quit
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
	case tea.KeyDown:
		if m.cursor < 3 {
			m.cursor++
		}
	case tea.KeyEnter:
		m.lastCorrect = m.game.Answer(m.cursor)
		if m.game.IsComplete() {
			m.phase = phaseDone
			sp := common.NewSharePrompt(m.game.Summary())
			m.sharePrompt = &sp
		} else {
			m.phase = phaseRevealed
		}
	case tea.KeyRunes:
		switch string(msg.Runes) {
		case "a", "A":
			m.cursor = 0
		case "b", "B":
			m.cursor = 1
		case "c", "C":
			m.cursor = 2
		case "d", "D":
			m.cursor = 3
		}
	}
	return m, nil
}

func (m Model) updateRevealed(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return m, tea.Quit
	case tea.KeyEnter, tea.KeySpace:
		m.phase = phaseAnswering
		m.cursor = 0
	}
	return m, nil
}

func (m Model) updateDone(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.sharePrompt != nil {
		prompt, quit := m.sharePrompt.HandleKey(msg.String())
		m.sharePrompt = &prompt
		if quit {
			return m, tea.Quit
		}
		return m, nil
	}
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc, tea.KeyEnter, tea.KeySpace:
		return m, tea.Quit
	}
	return m, nil
}

// View renders the game UI.
func (m Model) View() string {
	var b strings.Builder

	switch m.phase {
	case phaseAnswering:
		m.renderQuestion(&b)
	case phaseRevealed:
		m.renderReveal(&b)
	case phaseDone:
		m.renderDone(&b)
	}

	return b.String()
}

func (m Model) renderQuestion(b *strings.Builder) {
	q := m.game.CurrentQuestion()
	if q == nil {
		return
	}

	// Title + progress
	b.WriteString(common.TitleStyle.Render("🧠 TRIVIA"))
	b.WriteString("  ")
	b.WriteString(fmt.Sprintf("Question %d/%d", m.game.Current()+1, m.game.Total()))
	b.WriteString("  ")
	b.WriteString(categoryStyle.Render(q.Category))
	b.WriteString("\n\n")

	// Question text
	b.WriteString(questionStyle.Render(q.Text))
	b.WriteString("\n\n")

	// Choices
	labels := [4]string{"A", "B", "C", "D"}
	for i, choice := range q.Choices {
		label := fmt.Sprintf(" %s. %s ", labels[i], choice)
		if i == m.cursor {
			b.WriteString(choiceSelected.Render("▸ " + label))
		} else {
			b.WriteString(choiceNormal.Render("  " + label))
		}
		b.WriteString("\n")
	}

	// Score bar
	b.WriteString("\n")
	b.WriteString(scoreStyle.Render(fmt.Sprintf("Score: %d/%d", m.game.Score(), m.game.Total())))
	b.WriteString("\n\n")

	// Help
	b.WriteString(common.HelpStyle.Render("A/B/C/D or ↑↓ to select · Enter to confirm · ESC to quit"))
	b.WriteString("\n")
}

func (m Model) renderReveal(b *strings.Builder) {
	results := m.game.Results()
	last := results[len(results)-1]
	q := last.Question

	b.WriteString(common.TitleStyle.Render("🧠 TRIVIA"))
	b.WriteString("  ")
	b.WriteString(fmt.Sprintf("Question %d/%d", m.game.Current(), m.game.Total()))
	b.WriteString("  ")
	b.WriteString(categoryStyle.Render(q.Category))
	b.WriteString("\n\n")

	b.WriteString(questionStyle.Render(q.Text))
	b.WriteString("\n\n")

	// Show choices with correct/wrong highlighting
	labels := [4]string{"A", "B", "C", "D"}
	for i, choice := range q.Choices {
		label := fmt.Sprintf(" %s. %s ", labels[i], choice)
		if i == q.Answer {
			b.WriteString(choiceCorrectReveal.Render("✓ " + label))
		} else if i == last.Chosen && !last.Correct {
			b.WriteString(choiceWrongReveal.Render("✗ " + label))
		} else {
			b.WriteString(choiceNormal.Render("  " + label))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	if last.Correct {
		b.WriteString(common.SuccessStyle.Render("✓ Correct!"))
	} else {
		b.WriteString(common.ErrorStyle.Render(fmt.Sprintf("✗ Wrong! The answer was %s. %s", labels[q.Answer], q.Choices[q.Answer])))
	}

	b.WriteString("\n\n")
	b.WriteString(scoreStyle.Render(fmt.Sprintf("Score: %d/%d", m.game.Score(), m.game.Total())))
	b.WriteString("\n\n")
	b.WriteString(common.HelpStyle.Render("Press Enter or Space to continue"))
	b.WriteString("\n")
}

func (m Model) renderDone(b *strings.Builder) {
	score := m.game.Score()
	total := m.game.Total()

	b.WriteString(common.TitleStyle.Render("🧠 TRIVIA — Final Score"))
	b.WriteString("\n\n")

	// Grade
	grade := gradeEmoji(score, total)
	b.WriteString(gradeStyle.Render(fmt.Sprintf("%s  %d / %d", grade, score, total)))
	b.WriteString("\n\n")

	// Score bar
	filled := score
	empty := total - score
	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	b.WriteString(fmt.Sprintf("  [%s]  %d%%", bar, score*100/total))
	b.WriteString("\n\n")

	// Wrong answers review
	results := m.game.Results()
	wrongs := 0
	for _, r := range results {
		if !r.Correct {
			wrongs++
		}
	}

	if wrongs == 0 {
		b.WriteString(common.SuccessStyle.Render("🎉 Perfect score! You're a trivia champion!"))
		b.WriteString("\n")
	} else {
		b.WriteString(fmt.Sprintf("Missed %d question(s):\n\n", wrongs))
		labels := [4]string{"A", "B", "C", "D"}
		for i, r := range results {
			if r.Correct {
				continue
			}
			b.WriteString(common.ErrorStyle.Render(fmt.Sprintf(
				"  Q%d: %s", i+1, r.Question.Text,
			)))
			b.WriteString("\n")
			b.WriteString(fmt.Sprintf(
				"      You chose: %s. %s\n",
				labels[r.Chosen], r.Question.Choices[r.Chosen],
			))
			b.WriteString(common.SuccessStyle.Render(fmt.Sprintf(
				"      Correct:   %s. %s",
				labels[r.Question.Answer], r.Question.Choices[r.Question.Answer],
			)))
			b.WriteString("\n\n")
		}
	}

	if m.sharePrompt != nil {
		b.WriteString(m.sharePrompt.View())
	} else {
		b.WriteString(common.HelpStyle.Render("Press Enter or ESC to quit"))
	}
	b.WriteString("\n")
}

func gradeEmoji(score, total int) string {
	pct := score * 100 / total
	switch {
	case pct == 100:
		return "🏆 Perfect!"
	case pct >= 90:
		return "🌟 Amazing!"
	case pct >= 80:
		return "🎉 Great!"
	case pct >= 70:
		return "👍 Good job!"
	case pct >= 50:
		return "😅 Not bad!"
	case pct >= 30:
		return "📚 Keep studying!"
	default:
		return "💀 Oof!"
	}
}

// Run starts the trivia TUI.
func Run() error {
	m := NewModel()
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}

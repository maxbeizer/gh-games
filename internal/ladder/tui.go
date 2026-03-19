package ladder

import (
	"fmt"
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maxbeizer/gh-games/internal/common"
)

var (
	chainStyle = lipgloss.NewStyle().Bold(true).Foreground(common.ColorWhite)
	arrowStyle = lipgloss.NewStyle().Faint(true)
	matchStyle = lipgloss.NewStyle().Foreground(common.ColorGreen).Bold(true)
	diffStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	inputStyle = lipgloss.NewStyle().Bold(true).Foreground(common.ColorBlue)
	placeholderStyle = lipgloss.NewStyle().Faint(true)
	targetStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Bold(true)
	stepInfoStyle = lipgloss.NewStyle().Faint(true)
)

// Model is the Bubbletea model for the Ladder game.
type Model struct {
	Game  *Game
	input string
	err   string
}

// NewModel creates a new Ladder game TUI model.
func NewModel() Model {
	return Model{
		Game: NewGame(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.Game.IsWon() {
			if msg.Type == tea.KeyEsc || msg.Type == tea.KeyCtrlC {
				return m, tea.Quit
			}
			return m, nil
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
				m.err = ""
			}

		case tea.KeyEnter:
			if m.input == "" {
				return m, nil
			}
			word := strings.ToUpper(m.input)
			err := m.Game.Step(word)
			if err != nil {
				m.err = err.Error()
			} else {
				m.err = ""
			}
			m.input = ""

		case tea.KeyRunes:
			for _, r := range msg.Runes {
				if unicode.IsLetter(r) && len(m.input) < 4 {
					m.input += string(unicode.ToUpper(r))
					m.err = ""
				}
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	// Title
	b.WriteString(common.TitleStyle.Render("🪜 LADDER"))
	b.WriteString("\n\n")

	// Goal
	b.WriteString(fmt.Sprintf("  Get from %s to %s\n",
		chainStyle.Render(m.Game.Start),
		targetStyle.Render(m.Game.End)))
	b.WriteString(fmt.Sprintf("  Change one letter at a time — each step must be a real word\n\n"))

	// Word chain
	b.WriteString(m.renderChain())
	b.WriteString("\n")

	// Game state
	if m.Game.IsWon() {
		b.WriteString("\n")
		msg := fmt.Sprintf("🎉 Solved in %d steps!", m.Game.StepCount())
		b.WriteString(common.SuccessStyle.Render(msg))
		if m.Game.IsOptimal() {
			b.WriteString(common.SuccessStyle.Render(" ⭐ Optimal!"))
		}
		b.WriteString("\n")
		b.WriteString(common.HelpStyle.Render("Press ESC to quit."))
	} else {
		// Step counter
		optStr := ""
		if m.Game.Optimal > 0 {
			optStr = fmt.Sprintf(" (optimal: %d)", m.Game.Optimal)
		}
		b.WriteString(stepInfoStyle.Render(
			fmt.Sprintf("  Steps: %d%s\n\n", m.Game.StepCount(), optStr)))

		// Input
		if m.input == "" {
			b.WriteString(placeholderStyle.Render("  Type next word..."))
		} else {
			b.WriteString("  > " + inputStyle.Render(strings.ToUpper(m.input)))
		}
		b.WriteString("\n")

		// Error
		if m.err != "" {
			b.WriteString(common.ErrorStyle.Render("  ✗ " + m.err))
			b.WriteString("\n")
		}

		b.WriteString("\n")
		b.WriteString(common.HelpStyle.Render("  Type next word + Enter · ESC to quit"))
	}

	b.WriteString("\n")
	return b.String()
}

func (m Model) renderChain() string {
	var lines []string

	for i, word := range m.Game.Steps {
		var rendered string
		if i == 0 {
			rendered = "  " + chainStyle.Render(word)
		} else {
			prev := m.Game.Steps[i-1]
			rendered = "  " + renderWithHighlight(word, prev)
		}
		lines = append(lines, rendered)

		if i < len(m.Game.Steps)-1 {
			lines = append(lines, arrowStyle.Render("    ↓"))
		}
	}

	// Show target if not yet won
	if !m.Game.IsWon() {
		if len(m.Game.Steps) > 0 {
			lines = append(lines, arrowStyle.Render("    ↓"))
		}
		// Show remaining distance as question marks
		lines = append(lines, arrowStyle.Render("   ···"))
		lines = append(lines, arrowStyle.Render("    ↓"))
		lines = append(lines, "  "+targetStyle.Render(m.Game.End))
	}

	return strings.Join(lines, "\n")
}

// renderWithHighlight renders a word, highlighting the letter that changed.
func renderWithHighlight(word, prev string) string {
	var parts []string
	for i := range word {
		ch := string(word[i])
		if word[i] != prev[i] {
			parts = append(parts, diffStyle.Render(ch))
		} else {
			parts = append(parts, matchStyle.Render(ch))
		}
	}
	return strings.Join(parts, "")
}

// Run creates and runs the Ladder game TUI.
func Run() error {
	m := NewModel()
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}

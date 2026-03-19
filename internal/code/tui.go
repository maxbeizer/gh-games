package code

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maxbeizer/gh-games/internal/common"
)

// Lipgloss styles for peg colors
var pegStyles = map[Color]lipgloss.Style{
	Red:    lipgloss.NewStyle().Background(lipgloss.Color("#CC3333")).Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Padding(0, 1),
	Green:  lipgloss.NewStyle().Background(lipgloss.Color("#538D4E")).Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Padding(0, 1),
	Blue:   lipgloss.NewStyle().Background(lipgloss.Color("#4A90D9")).Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Padding(0, 1),
	Yellow: lipgloss.NewStyle().Background(lipgloss.Color("#B59F3B")).Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Padding(0, 1),
	Purple: lipgloss.NewStyle().Background(lipgloss.Color("#9B59B6")).Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Padding(0, 1),
	Orange: lipgloss.NewStyle().Background(lipgloss.Color("#E67E22")).Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Padding(0, 1),
}

var emptySlotStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#666666")).
	Padding(0, 1)

// Model is the Bubbletea model for the Code Breaker TUI.
type Model struct {
	Game  *Game
	input []Color // current guess being built (0-4 colors)
	err   string
}

// NewModel creates a new TUI model with a random secret.
func NewModel() Model {
	return Model{
		Game:  NewGame(),
		input: make([]Color, 0, CodeLen),
	}
}

// NewModelWithGame creates a TUI model with a specific game (for testing).
func NewModelWithGame(g *Game) Model {
	return Model{
		Game:  g,
		input: make([]Color, 0, CodeLen),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.Game.IsOver() {
			return m, tea.Quit
		}

		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit

		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
				m.err = ""
			}

		case "enter":
			if len(m.input) != CodeLen {
				m.err = fmt.Sprintf("Need %d colors (have %d)", CodeLen, len(m.input))
				return m, nil
			}
			var code [CodeLen]Color
			copy(code[:], m.input)
			m.Game.MakeGuess(code)
			m.input = m.input[:0]
			m.err = ""

		default:
			c, ok := keyToColor(msg.String())
			if ok && len(m.input) < CodeLen {
				m.input = append(m.input, c)
				m.err = ""
			}
		}
	}

	return m, nil
}

func keyToColor(key string) (Color, bool) {
	switch strings.ToLower(key) {
	case "1", "r":
		return Red, true
	case "2", "g":
		return Green, true
	case "3", "b":
		return Blue, true
	case "4", "y":
		return Yellow, true
	case "5", "p":
		return Purple, true
	case "6", "o":
		return Orange, true
	default:
		return 0, false
	}
}

func (m Model) View() string {
	var b strings.Builder

	// Title
	b.WriteString(common.TitleStyle.Render("🔐 CODE"))
	b.WriteString("\n\n")

	// Secret row
	if m.Game.IsOver() {
		b.WriteString("Secret: ")
		for _, c := range m.Game.Secret {
			b.WriteString(renderPeg(c))
			b.WriteString(" ")
		}
	} else {
		b.WriteString("Secret: ")
		for i := 0; i < CodeLen; i++ {
			b.WriteString(emptySlotStyle.Render("?"))
			b.WriteString(" ")
		}
	}
	b.WriteString("\n\n")

	// Guess history
	for i, g := range m.Game.Guesses {
		b.WriteString(fmt.Sprintf(" %2d │ ", i+1))
		for _, c := range g.Code {
			b.WriteString(renderPeg(c))
			b.WriteString(" ")
		}
		b.WriteString("│ ")
		b.WriteString(renderFeedback(g.Feedback))
		b.WriteString("\n")
	}

	// Game over messages
	if m.Game.IsWon() {
		b.WriteString("\n")
		b.WriteString(common.SuccessStyle.Render(
			fmt.Sprintf("🎉 Cracked it in %d guess%s!",
				len(m.Game.Guesses),
				pluralize(len(m.Game.Guesses)))))
		b.WriteString("\n\n")
		b.WriteString(common.HelpStyle.Render("Press any key to quit"))
		b.WriteString("\n")
		return b.String()
	}

	if m.Game.IsLost() {
		b.WriteString("\n")
		b.WriteString(common.ErrorStyle.Render("😔 The code was: "))
		for _, c := range m.Game.Secret {
			b.WriteString(ColorSymbol(c))
		}
		b.WriteString("\n\n")
		b.WriteString(common.HelpStyle.Render("Press any key to quit"))
		b.WriteString("\n")
		return b.String()
	}

	// Current input
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Guess %d/%d: ", len(m.Game.Guesses)+1, m.Game.MaxTurns))
	for i := 0; i < CodeLen; i++ {
		if i < len(m.input) {
			b.WriteString(renderPeg(m.input[i]))
		} else {
			b.WriteString(emptySlotStyle.Render("_"))
		}
		b.WriteString(" ")
	}
	b.WriteString("\n")

	// Error
	if m.err != "" {
		b.WriteString(common.ErrorStyle.Render(m.err))
		b.WriteString("\n")
	}

	// Legend
	b.WriteString("\n")
	b.WriteString("1:")
	b.WriteString(ColorSymbol(Red))
	b.WriteString(" 2:")
	b.WriteString(ColorSymbol(Green))
	b.WriteString(" 3:")
	b.WriteString(ColorSymbol(Blue))
	b.WriteString(" 4:")
	b.WriteString(ColorSymbol(Yellow))
	b.WriteString(" 5:")
	b.WriteString(ColorSymbol(Purple))
	b.WriteString(" 6:")
	b.WriteString(ColorSymbol(Orange))
	b.WriteString("\n")

	// Help
	b.WriteString(common.HelpStyle.Render("1-6 or R/G/B/Y/P/O to pick colors · Backspace to undo · Enter to guess · ESC to quit"))
	b.WriteString("\n")

	return b.String()
}

func renderPeg(c Color) string {
	style, ok := pegStyles[c]
	if !ok {
		return "?"
	}
	return style.Render(ColorLetter(c))
}

func renderFeedback(fb Feedback) string {
	var parts []string
	for i := 0; i < fb.Exact; i++ {
		parts = append(parts, "●")
	}
	for i := 0; i < fb.Misplaced; i++ {
		parts = append(parts, "○")
	}
	// Pad to CodeLen so alignment is consistent
	for len(parts) < CodeLen {
		parts = append(parts, " ")
	}
	return strings.Join(parts, " ")
}

func pluralize(n int) string {
	if n == 1 {
		return ""
	}
	return "es"
}

// Run starts the Code Breaker TUI.
func Run() error {
	m := NewModel()
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}

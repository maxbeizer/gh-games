package group

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maxbeizer/gh-games/internal/common"
)

const gridCols = 4

// Model is the Bubbletea model for the Group game.
type Model struct {
	game     *Game
	cursor   int
	message  string
	gameOver bool
}

// NewModel creates a Model with a fresh game.
func NewModel() Model {
	return Model{game: NewGame()}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.gameOver {
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.String() {
			case "q", "esc", "ctrl+c":
				return m, tea.Quit
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "up":
			if m.cursor >= gridCols {
				m.cursor -= gridCols
			}
		case "down":
			total := len(m.game.RemainingWords)
			if m.cursor+gridCols < total {
				m.cursor += gridCols
			}
		case "left":
			if m.cursor > 0 {
				m.cursor--
			}
		case "right":
			if m.cursor < len(m.game.RemainingWords)-1 {
				m.cursor++
			}

		case " ":
			if len(m.game.RemainingWords) > 0 {
				word := m.game.RemainingWords[m.cursor]
				m.game.ToggleSelect(word)
				m.message = ""
			}

		case "enter":
			matched, err := m.game.Submit()
			if err != nil {
				switch err {
				case ErrNotEnoughSelected:
					m.message = "Select exactly 4 words before submitting"
				case ErrNoMatch:
					m.message = "Not a valid group — try again"
				case ErrGameOver:
					m.message = "Game is already over"
				}
			} else {
				m.message = fmt.Sprintf("✓ Solved: %s", matched.Name)
			}
			// Clamp cursor after words are removed
			if m.cursor >= len(m.game.RemainingWords) && len(m.game.RemainingWords) > 0 {
				m.cursor = len(m.game.RemainingWords) - 1
			}
			if m.game.IsOver() {
				m.gameOver = true
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	// Title
	b.WriteString(common.TitleStyle.Render("🔗 GROUP"))
	b.WriteString("\n\n")

	// Solved groups
	for _, cat := range m.game.SolvedGroups {
		style := common.TierStyles[cat.Difficulty]
		banner := style.Padding(0, 1).Render(
			fmt.Sprintf(" %s: %s ", cat.Name, strings.Join(cat.Words, ", ")),
		)
		b.WriteString(banner)
		b.WriteString("\n")
	}
	if len(m.game.SolvedGroups) > 0 {
		b.WriteString("\n")
	}

	// Game over screen
	if m.gameOver {
		if m.game.IsWon() {
			if m.game.Mistakes == 0 {
				b.WriteString(common.SuccessStyle.Render("🎉 Perfect!"))
			} else {
				b.WriteString(common.SuccessStyle.Render(
					fmt.Sprintf("🎉 Solved with %d mistake%s!", m.game.Mistakes, plural(m.game.Mistakes)),
				))
			}
		} else {
			// Reveal remaining groups
			for _, cat := range m.game.RemainingCategories() {
				style := common.TierStyles[cat.Difficulty]
				banner := style.Padding(0, 1).Render(
					fmt.Sprintf(" %s: %s ", cat.Name, strings.Join(cat.Words, ", ")),
				)
				b.WriteString(banner)
				b.WriteString("\n")
			}
			b.WriteString("\n")
			b.WriteString(common.ErrorStyle.Render("😔 Better luck next time!"))
		}
		b.WriteString("\n\n")
		b.WriteString(common.HelpStyle.Render("Press Q or ESC to quit"))
		return b.String()
	}

	// Word grid
	words := m.game.RemainingWords
	cellWidth := 12
	for i, w := range words {
		isSelected := m.game.Selected[w]
		isCursor := i == m.cursor

		cell := common.CellStyle(cellWidth, 1)
		label := w

		if isCursor && isSelected {
			cell = cell.
				Background(lipgloss.Color("#555555")).
				Foreground(common.ColorWhite).
				Bold(true).
				Border(lipgloss.ThickBorder()).
				BorderForeground(common.ColorWhite)
			label = "▸" + label
		} else if isCursor {
			cell = cell.
				Background(lipgloss.Color("#555555")).
				Foreground(common.ColorWhite).
				Bold(true)
			label = "▸" + label
		} else if isSelected {
			cell = cell.
				Border(lipgloss.ThickBorder()).
				BorderForeground(common.ColorWhite)
		}

		b.WriteString(cell.Render(label))
		if (i+1)%gridCols == 0 {
			b.WriteString("\n")
		} else {
			b.WriteString(" ")
		}
	}
	b.WriteString("\n")

	// Mistakes indicator
	filled := m.game.Mistakes
	empty := m.game.MaxMistakes - filled
	dots := strings.Repeat("●", filled) + strings.Repeat("○", empty)
	b.WriteString(fmt.Sprintf("Mistakes: %s (%d/%d)\n", dots, m.game.Mistakes, m.game.MaxMistakes))

	// Status/Error message
	if m.message != "" {
		if strings.HasPrefix(m.message, "✓") {
			b.WriteString(common.SuccessStyle.Render(m.message))
		} else {
			b.WriteString(common.ErrorStyle.Render(m.message))
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Help
	b.WriteString(common.HelpStyle.Render("↑↓←→ Navigate • Space Select • Enter Submit • ESC Quit"))
	b.WriteString("\n")

	return b.String()
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// Run is a convenience function that runs the Group game TUI.
func Run() error {
	p := tea.NewProgram(NewModel())
	_, err := p.Run()
	return err
}

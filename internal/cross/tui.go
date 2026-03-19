package cross

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maxbeizer/gh-games/internal/common"
)

// tickMsg is sent every second to update the timer.
type tickMsg time.Time

// Model is the bubbletea model for the crossword TUI.
type Model struct {
	Game        *Game
	status      string
	checked     bool
	solved      bool
	sharePrompt *common.SharePrompt
}

// NewModel creates a new crossword TUI model.
func NewModel() Model {
	return Model{
		Game: NewGame(),
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return tickCmd()
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if !m.solved {
			return m, tickCmd()
		}
		return m, nil
	case tea.KeyMsg:
		if m.solved {
			if m.sharePrompt != nil {
				prompt, quit := m.sharePrompt.HandleKey(msg.String())
				m.sharePrompt = &prompt
				if quit {
					return m, tea.Quit
				}
				return m, nil
			}
			switch msg.String() {
			case "q", "esc", "ctrl+c":
				return m, tea.Quit
			}
			return m, nil
		}
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "ctrl+k":
			checked := m.Game.Check()
			m.Game.Checked = checked
			m.checked = true
			hasErrors := false
			for r := 0; r < 5; r++ {
				for c := 0; c < 5; c++ {
					if checked[r][c] {
						hasErrors = true
					}
				}
			}
			if !hasErrors && m.Game.IsComplete() {
				m.solved = true
				m.Game.Completed = true
				m.status = "🎉 Solved!"
				sp := common.NewSharePrompt(m.Game.Summary())
				m.sharePrompt = &sp
			} else if hasErrors {
				m.status = "❌ Some letters are wrong (shown in red)"
			} else {
				m.status = "✅ All entered letters are correct so far"
			}
			return m, nil
		case "tab":
			m.Game.ToggleDirection()
			m.checked = false
			return m, nil
		case "backspace":
			m.Game.ClearLetter()
			m.Game.Retreat()
			m.checked = false
			return m, nil
		case "up":
			m.Game.MoveCursor(0)
			m.checked = false
			return m, nil
		case "down":
			m.Game.MoveCursor(1)
			m.checked = false
			return m, nil
		case "left":
			m.Game.MoveCursor(2)
			m.checked = false
			return m, nil
		case "right":
			m.Game.MoveCursor(3)
			m.checked = false
			return m, nil
		default:
			s := msg.String()
			if len(s) == 1 && unicode.IsLetter(rune(s[0])) {
				m.Game.SetLetter(rune(s[0]))
				m.checked = false
				// Auto-check on completion
				if m.Game.IsComplete() && m.Game.IsCorrect() {
					m.solved = true
					m.Game.Completed = true
					m.status = "🎉 Solved!"
					sp := common.NewSharePrompt(m.Game.Summary())
					m.sharePrompt = &sp
					return m, nil
				}
				m.Game.Advance()
			}
		}
	}
	return m, nil
}

// View implements tea.Model.
func (m Model) View() string {
	var sb strings.Builder

	// Styles
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	highlight := lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57"))
	blackCell := lipgloss.NewStyle().Background(lipgloss.Color("236")).Foreground(lipgloss.Color("236"))
	normalCell := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	wrongCell := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	clueNormal := lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	clueHighlight := lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Bold(true)
	dirStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("117")).Bold(true)

	// Header
	elapsed := m.Game.Elapsed().Truncate(time.Second)
	sb.WriteString(title.Render("📰 CROSSWORD"))
	sb.WriteString("  ")
	sb.WriteString(dim.Render(fmt.Sprintf("⏱ %s", elapsed)))
	dirLabel := "ACROSS →"
	if m.Game.Dir == Down {
		dirLabel = "DOWN ↓"
	}
	sb.WriteString("  ")
	sb.WriteString(dirStyle.Render(dirLabel))
	sb.WriteString("\n\n")

	// Build grid lines
	gridLines := make([]string, 0, 7)
	for r := 0; r < 5; r++ {
		var row strings.Builder
		for c := 0; c < 5; c++ {
			cell := m.Game.Puzzle.Grid[r][c]
			isCursor := r == m.Game.CurRow && c == m.Game.CurCol

			if cell.Black {
				row.WriteString(blackCell.Render("███"))
			} else {
				letter := m.Game.Player[r][c]
				num := m.Game.CellNumber(r, c)
				display := " "
				if letter != 0 {
					display = string(letter)
				}
				// Build the cell content: number (small) + letter
				content := ""
				if num > 0 {
					numStr := fmt.Sprintf("%d", num)
					if len(numStr) == 1 {
						content = numStr + display + " "
					} else {
						content = numStr + display
					}
				} else {
					content = " " + display + " "
				}

				if isCursor {
					row.WriteString(highlight.Render(content))
				} else if m.checked && m.Game.Checked[r][c] {
					row.WriteString(wrongCell.Render(content))
				} else {
					row.WriteString(normalCell.Render(content))
				}
			}
			if c < 4 {
				row.WriteString(" ")
			}
		}
		gridLines = append(gridLines, row.String())
	}

	// Build clue lines
	clueLines := make([]string, 0)
	currentClue := m.Game.CurrentClue()

	clueLines = append(clueLines, dim.Render("  ACROSS"))
	for _, cl := range m.Game.Puzzle.Clues {
		if cl.Direction != "across" {
			continue
		}
		text := fmt.Sprintf("  %2d. %s", cl.Number, cl.Text)
		if currentClue != nil && cl.Number == currentClue.Number && cl.Direction == currentClue.Direction {
			clueLines = append(clueLines, clueHighlight.Render(text))
		} else {
			clueLines = append(clueLines, clueNormal.Render(text))
		}
	}
	clueLines = append(clueLines, "")
	clueLines = append(clueLines, dim.Render("  DOWN"))
	for _, cl := range m.Game.Puzzle.Clues {
		if cl.Direction != "down" {
			continue
		}
		text := fmt.Sprintf("  %2d. %s", cl.Number, cl.Text)
		if currentClue != nil && cl.Number == currentClue.Number && cl.Direction == currentClue.Direction {
			clueLines = append(clueLines, clueHighlight.Render(text))
		} else {
			clueLines = append(clueLines, clueNormal.Render(text))
		}
	}

	// Join grid and clues side by side
	maxLines := len(gridLines)
	if len(clueLines) > maxLines {
		maxLines = len(clueLines)
	}
	for i := 0; i < maxLines; i++ {
		gridLine := ""
		if i < len(gridLines) {
			gridLine = gridLines[i]
		}
		clueLine := ""
		if i < len(clueLines) {
			clueLine = clueLines[i]
		}
		sb.WriteString(fmt.Sprintf("  %-24s%s\n", gridLine, clueLine))
	}

	// Status
	sb.WriteString("\n")
	if m.status != "" {
		sb.WriteString("  " + m.status + "\n")
	}

	// Help
	if m.sharePrompt != nil {
		sb.WriteString("  " + m.sharePrompt.View() + "\n")
	} else {
		help := "Type to fill · ↑↓←→ move · Tab toggle direction · Ctrl+K check · ESC quit"
		sb.WriteString("  " + dim.Render(help) + "\n")
	}

	return sb.String()
}

// Run starts the crossword TUI.
func Run() error {
	m := NewModel()
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}

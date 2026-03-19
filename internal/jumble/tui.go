package jumble

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maxbeizer/gh-games/internal/common"
)

type tickMsg time.Time

// Model is the Bubbletea model for the Jumble game.
type Model struct {
	Game        *Game
	input       string
	message     string
	messageGood bool
	lastPoints  int
	showSummary bool
}

// NewModel creates a new Jumble game TUI model.
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

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		return m, tickCmd()

	case tea.KeyMsg:
		// Summary screen — any key quits
		if m.showSummary {
			return m, tea.Quit
		}

		r := m.Game.CurrentRoundRef()

		// After solving, Enter advances to next round
		if r != nil && r.Solved {
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				return m, tea.Quit
			case tea.KeyEnter:
				if !m.Game.NextRound() {
					m.showSummary = true
				}
				m.input = ""
				m.message = ""
				return m, nil
			}
			return m, nil
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
				m.message = ""
			}

		case tea.KeyEnter:
			if len(m.input) == 0 {
				return m, nil
			}
			correct, points := m.Game.Guess(m.input)
			if correct {
				m.message = fmt.Sprintf("✅ Correct! +%d points", points)
				m.messageGood = true
				m.lastPoints = points
			} else {
				m.message = "❌ Not quite — try again!"
				m.messageGood = false
			}
			if !correct {
				m.input = ""
			}

		case tea.KeyTab:
			m.Game.Shuffle()

		case tea.KeyRunes:
			for _, r := range msg.Runes {
				if r == '?' {
					letter, pos, err := m.Game.Hint()
					if err == nil {
						m.message = fmt.Sprintf("💡 Hint: letter %d is '%c'", pos+1, letter)
						m.messageGood = true
					}
					continue
				}
				if unicode.IsLetter(r) {
					curRound := m.Game.CurrentRoundRef()
					if curRound != nil && len(m.input) < len(curRound.Target) {
						m.input += string(unicode.ToUpper(r))
						m.message = ""
					}
				}
			}
		}
	}
	return m, nil
}

var (
	letterStyle = lipgloss.NewStyle().
			Bold(true).
			Width(5).
			Height(1).
			Align(lipgloss.Center, lipgloss.Center).
			Background(common.ColorBlue).
			Foreground(common.ColorWhite)

	hintLetterStyle = lipgloss.NewStyle().
			Bold(true).
			Width(5).
			Height(1).
			Align(lipgloss.Center, lipgloss.Center).
			Background(common.ColorGreen).
			Foreground(common.ColorWhite)

	emptyLetterStyle = lipgloss.NewStyle().
				Width(5).
				Height(1).
				Align(lipgloss.Center, lipgloss.Center).
				Border(lipgloss.NormalBorder()).
				BorderForeground(common.ColorGray)

	timerStyle = lipgloss.NewStyle().
			Foreground(common.ColorYellow).
			Bold(true)

	scoreStyle = lipgloss.NewStyle().
			Foreground(common.ColorGreen).
			Bold(true)

	roundStyle = lipgloss.NewStyle().
			Foreground(common.ColorPurple).
			Bold(true)

	summaryHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(common.ColorBlue).
				MarginBottom(1)

	summaryRowStyle = lipgloss.NewStyle().
			PaddingLeft(2)
)

func (m Model) View() string {
	if m.showSummary {
		return m.renderSummary()
	}

	var b strings.Builder

	r := m.Game.CurrentRoundRef()
	if r == nil {
		return "Game error\n"
	}

	// Title + Round indicator
	title := common.TitleStyle.Render("🔀 JUMBLE")
	round := roundStyle.Render(fmt.Sprintf("Round %d/5", m.Game.CurrentRound+1))
	b.WriteString(fmt.Sprintf("%s  %s\n\n", title, round))

	// Scrambled word display
	b.WriteString(m.renderScrambled(r))
	b.WriteString("\n\n")

	// Timer + Score
	elapsed := time.Since(r.StartTime)
	if r.Solved {
		elapsed = r.SolveTime
	}
	timer := timerStyle.Render(fmt.Sprintf("⏱  %ds", int(elapsed.Seconds())))
	score := scoreStyle.Render(fmt.Sprintf("Score: %d", m.Game.TotalScore))
	hints := ""
	if r.HintsUsed > 0 {
		hints = fmt.Sprintf("  💡×%d", r.HintsUsed)
	}
	b.WriteString(fmt.Sprintf("%s    %s%s\n\n", timer, score, hints))

	// Input or solved state
	if r.Solved {
		if m.message != "" {
			b.WriteString(common.SuccessStyle.Render(m.message))
		}
		b.WriteString("\n")
		if m.Game.CurrentRound < len(m.Game.Rounds)-1 {
			b.WriteString(common.HelpStyle.Render("Press Enter for next round · ESC to quit"))
		} else {
			b.WriteString(common.HelpStyle.Render("Press Enter to see results · ESC to quit"))
		}
	} else {
		// Input line
		b.WriteString(fmt.Sprintf("  > %s", m.input))
		if m.message != "" {
			b.WriteString("\n")
			if m.messageGood {
				b.WriteString(common.SuccessStyle.Render(m.message))
			} else {
				b.WriteString(common.ErrorStyle.Render(m.message))
			}
		}
		b.WriteString("\n\n")
		b.WriteString(common.HelpStyle.Render("Type answer + Enter · ? for hint · Tab to shuffle · ESC to quit"))
	}

	b.WriteString("\n")
	return b.String()
}

func (m Model) renderScrambled(r *Round) string {
	hintSet := make(map[int]bool)
	for _, pos := range r.Hints {
		hintSet[pos] = true
	}

	var cells []string
	for i, ch := range r.Scrambled {
		if hintSet[i] {
			cells = append(cells, hintLetterStyle.Render(string(ch)))
		} else {
			cells = append(cells, letterStyle.Render(string(ch)))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, cells...)
}

func (m Model) renderSummary() string {
	var b strings.Builder

	b.WriteString(common.TitleStyle.Render("🔀 JUMBLE — Final Results"))
	b.WriteString("\n\n")

	for i, r := range m.Game.Rounds {
		status := "⬜"
		if r.Solved {
			status = "✅"
		}
		wordLen := len(r.Target)
		line := fmt.Sprintf("%s Round %d (%d letters): %s", status, i+1, wordLen, r.Target)
		if r.Solved {
			line += fmt.Sprintf("  ⏱ %ds", int(r.SolveTime.Seconds()))
			if r.HintsUsed > 0 {
				line += fmt.Sprintf("  💡×%d", r.HintsUsed)
			}
		}
		b.WriteString(summaryRowStyle.Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(summaryHeaderStyle.Render(
		fmt.Sprintf("🏆 Total Score: %d", m.Game.TotalScore)))
	b.WriteString("\n\n")
	b.WriteString(common.HelpStyle.Render("Press any key to quit."))
	b.WriteString("\n")

	return b.String()
}

// Run creates and runs the Jumble game TUI.
func Run() error {
	m := NewModel()
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}

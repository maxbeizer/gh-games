package hive

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maxbeizer/gh-games/internal/common"
)

// Hive-specific styles
var (
	centerLetterStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#F5C518")).
				Padding(0, 1)

	outerLetterStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#3A3A3C")).
				Padding(0, 1)

	inputStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F5C518"))

	placeholderStyle = lipgloss.NewStyle().
				Faint(true).
				Italic(true)

	pangramStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F5C518"))

	progressFilled = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F5C518"))

	progressEmpty = lipgloss.NewStyle().
			Faint(true)

	giveUpStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#CC3333"))
)

// Model is the Bubbletea model for the Hive game.
type Model struct {
	Game           *Game
	input          string
	status         string
	statusIsError  bool
	showGiveUp     bool
	givenUp        bool
	letterSet      map[rune]bool
}

// NewModel creates a new Hive game TUI model.
func NewModel() Model {
	g := NewGame()
	ls := make(map[rune]bool, 7)
	for _, l := range g.Letters {
		ls[l] = true
	}
	return Model{
		Game:      g,
		letterSet: ls,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// If given up, only allow quit
		if m.givenUp {
			if msg.Type == tea.KeyEsc || msg.Type == tea.KeyCtrlC {
				return m, tea.Quit
			}
			return m, nil
		}

		// Give-up confirmation mode
		if m.showGiveUp {
			switch msg.Type {
			case tea.KeyEnter:
				m.showGiveUp = false
				m.givenUp = true
				m.status = ""
				return m, nil
			case tea.KeyEsc, tea.KeyCtrlC:
				m.showGiveUp = false
				m.status = ""
				return m, nil
			default:
				if msg.Type == tea.KeyRunes {
					r := unicode.ToLower(msg.Runes[0])
					if r == 'y' {
						m.showGiveUp = false
						m.givenUp = true
						m.status = ""
						return m, nil
					}
				}
				m.showGiveUp = false
				m.status = ""
				return m, nil
			}
		}

		// Won the game — only allow quit
		if m.Game.Score == m.Game.MaxScore {
			if msg.Type == tea.KeyEsc || msg.Type == tea.KeyCtrlC {
				return m, tea.Quit
			}
			return m, nil
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyCtrlG:
			m.showGiveUp = true
			m.status = ""
			return m, nil

		case tea.KeyTab:
			m.shuffleLetters()
			return m, nil

		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}

		case tea.KeyEnter:
			if m.input == "" {
				return m, nil
			}
			points, pangram, err := m.Game.Guess(m.input)
			if err != nil {
				m.status = "✗ " + capitalize(err.Error())
				m.statusIsError = true
			} else {
				word := strings.ToUpper(m.input)
				msg := fmt.Sprintf("✓ %s +%d", word, points)
				if pangram {
					msg += " 🐝 PANGRAM!"
				}
				m.status = msg
				m.statusIsError = false
			}
			m.input = ""

		case tea.KeyRunes:
			for _, r := range msg.Runes {
				if !unicode.IsLetter(r) {
					continue
				}
				upper := unicode.ToUpper(r)
				if !m.letterSet[upper] {
					m.status = fmt.Sprintf("✗ Letter '%c' is not in the puzzle", upper)
					m.statusIsError = true
					continue
				}
				m.input += string(upper)
				// Clear status on new typing
				if m.statusIsError {
					m.status = ""
				}
			}
		}
	}
	return m, nil
}

func (m *Model) shuffleLetters() {
	outer := m.Game.Letters[1:]
	rand.Shuffle(len(outer), func(i, j int) {
		outer[i], outer[j] = outer[j], outer[i]
	})
}

func (m Model) View() string {
	var b strings.Builder

	// Title
	b.WriteString(common.TitleStyle.Render("🐝 HIVE"))
	b.WriteString("\n\n")

	// Honeycomb
	b.WriteString(m.renderHoneycomb())
	b.WriteString("\n\n")

	// Input
	if m.givenUp {
		b.WriteString(giveUpStyle.Render("Game over — all words revealed"))
	} else if m.Game.Score == m.Game.MaxScore {
		b.WriteString(common.SuccessStyle.Render("🎉 Queen Bee! You found all words!"))
	} else if m.showGiveUp {
		b.WriteString(giveUpStyle.Render("Give up? Press Y or Enter to confirm, any other key to cancel"))
	} else {
		if m.input == "" {
			b.WriteString(placeholderStyle.Render("  Type a word..."))
		} else {
			b.WriteString("  " + inputStyle.Render(strings.ToUpper(m.input)))
		}
	}
	b.WriteString("\n")

	// Status
	if m.status != "" {
		if m.statusIsError {
			b.WriteString(common.ErrorStyle.Render(m.status))
		} else {
			b.WriteString(common.SuccessStyle.Render(m.status))
		}
		b.WriteString("\n")
	} else {
		b.WriteString("\n")
	}

	// Score & Rank
	b.WriteString(fmt.Sprintf("Score: %d/%d · Rank: %s · %d/%d words found",
		m.Game.Score, m.Game.MaxScore, m.Game.Rank(),
		len(m.Game.Found), len(m.Game.AllValid)))
	b.WriteString("\n")

	// Progress bar
	b.WriteString(m.renderProgressBar(20))
	b.WriteString("\n\n")

	// Found words
	b.WriteString(m.renderFoundWords())
	b.WriteString("\n")

	// Give-up: show unfound words
	if m.givenUp {
		b.WriteString(m.renderUnfoundWords())
		b.WriteString("\n")
	}

	// Help
	if m.givenUp || m.Game.Score == m.Game.MaxScore {
		b.WriteString(common.HelpStyle.Render("Press ESC to quit."))
	} else {
		b.WriteString(common.HelpStyle.Render(
			"Type a word + Enter · Backspace to delete · Shuffle: Tab · Give up: Ctrl+G · Quit: ESC"))
	}
	b.WriteString("\n")

	return b.String()
}

func (m Model) renderHoneycomb() string {
	letters := m.Game.Letters
	fmtOuter := func(r rune) string {
		return outerLetterStyle.Render(fmt.Sprintf(" %c ", unicode.ToUpper(r)))
	}
	fmtCenter := func(r rune) string {
		return centerLetterStyle.Render(fmt.Sprintf("*%c*", unicode.ToUpper(r)))
	}

	row1 := fmt.Sprintf("      %s  %s", fmtOuter(letters[1]), fmtOuter(letters[2]))
	row2 := fmt.Sprintf("   %s  %s  %s", fmtOuter(letters[3]), fmtCenter(letters[0]), fmtOuter(letters[4]))
	row3 := fmt.Sprintf("      %s  %s", fmtOuter(letters[5]), fmtOuter(letters[6]))

	return lipgloss.JoinVertical(lipgloss.Left, row1, row2, row3)
}

func (m Model) renderProgressBar(width int) string {
	pct := 0.0
	if m.Game.MaxScore > 0 {
		pct = float64(m.Game.Score) / float64(m.Game.MaxScore)
	}
	filled := int(pct * float64(width))
	if filled > width {
		filled = width
	}

	bar := progressFilled.Render(strings.Repeat("█", filled)) +
		progressEmpty.Render(strings.Repeat("░", width-filled))

	return fmt.Sprintf("[%s] %d%%", bar, int(pct*100))
}

func (m Model) renderFoundWords() string {
	if len(m.Game.Found) == 0 {
		return common.HelpStyle.Render("No words found yet.\n")
	}

	sorted := make([]string, len(m.Game.Found))
	copy(sorted, m.Game.Found)
	sort.Strings(sorted)

	cols := 3
	if len(sorted) < 6 {
		cols = 2
	}
	if len(sorted) < 3 {
		cols = 1
	}

	rows := (len(sorted) + cols - 1) / cols
	colWidth := 16

	var columns []string
	for c := 0; c < cols; c++ {
		var lines []string
		for r := 0; r < rows; r++ {
			idx := c*rows + r
			if idx >= len(sorted) {
				break
			}
			w := strings.ToUpper(sorted[idx])
			if m.Game.IsPangram(sorted[idx]) {
				w = pangramStyle.Render("🐝 " + w)
			}
			// Pad to column width
			for len([]rune(w)) < colWidth {
				w += " "
			}
			lines = append(lines, w)
		}
		columns = append(columns, strings.Join(lines, "\n"))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, columns...)
}

func (m Model) renderUnfoundWords() string {
	foundSet := make(map[string]bool, len(m.Game.Found))
	for _, w := range m.Game.Found {
		foundSet[w] = true
	}

	var unfound []string
	for _, w := range m.Game.AllValid {
		if !foundSet[w] {
			unfound = append(unfound, w)
		}
	}

	if len(unfound) == 0 {
		return ""
	}

	header := giveUpStyle.Render(fmt.Sprintf("Unfound words (%d):", len(unfound)))

	cols := 3
	if len(unfound) < 6 {
		cols = 2
	}
	if len(unfound) < 3 {
		cols = 1
	}

	rows := (len(unfound) + cols - 1) / cols
	colWidth := 16

	var columns []string
	for c := 0; c < cols; c++ {
		var lines []string
		for r := 0; r < rows; r++ {
			idx := c*rows + r
			if idx >= len(unfound) {
				break
			}
			w := strings.ToUpper(unfound[idx])
			if m.Game.IsPangram(unfound[idx]) {
				w = pangramStyle.Render("🐝 " + w)
			}
			for len([]rune(w)) < colWidth {
				w += " "
			}
			lines = append(lines, w)
		}
		columns = append(columns, strings.Join(lines, "\n"))
	}

	return header + "\n" + lipgloss.JoinHorizontal(lipgloss.Top, columns...) + "\n"
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// Run creates and runs the Hive game TUI.
func Run() error {
	m := NewModel()
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}

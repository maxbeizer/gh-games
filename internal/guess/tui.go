package guess

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maxbeizer/gh-games/internal/common"
)

const cellWidth = 5

var keyboardRows = []string{"QWERTYUIOP", "ASDFGHJKL", "ZXCVBNM"}

// Model is the Bubbletea model for the Guess game.
type Model struct {
	Game          *Game
	input         string
	err           string
	validateWords bool
	sharePrompt   *common.SharePrompt
}

// NewModel creates a new Guess game TUI model.
func NewModel(target string, validateWords bool) Model {
	return Model{
		Game:          NewGame(target),
		validateWords: validateWords,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.sharePrompt != nil {
			prompt, quit := m.sharePrompt.HandleKey(msg.String())
			m.sharePrompt = &prompt
			if quit {
				return m, tea.Quit
			}
			return m, nil
		}

		if m.Game.IsOver() {
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
			if len(m.input) != 5 {
				m.err = "Word must be exactly 5 letters"
				return m, nil
			}
			word := strings.ToUpper(m.input)
			if m.validateWords && !IsValidWord(word) {
				m.err = fmt.Sprintf("%s is not a valid word", word)
				return m, nil
			}
			_, err := m.Game.Guess(word)
			if err != nil {
				m.err = err.Error()
				return m, nil
			}
			m.input = ""
			m.err = ""
			if m.Game.IsOver() {
				sp := common.NewSharePrompt(m.Game.Summary())
				m.sharePrompt = &sp
			}

		case tea.KeyRunes:
			for _, r := range msg.Runes {
				if unicode.IsLetter(r) && len(m.input) < 5 {
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
	b.WriteString(common.TitleStyle.Render("🟩 GUESS"))
	b.WriteString("\n\n")

	// Grid
	b.WriteString(m.renderGrid())
	b.WriteString("\n\n")

	// Keyboard
	b.WriteString(m.renderKeyboard())
	b.WriteString("\n\n")

	// Game over or input area
	if m.Game.IsOver() {
		if m.Game.IsWon() {
			b.WriteString(common.SuccessStyle.Render(
				fmt.Sprintf("🎉 You got it in %d/6!", len(m.Game.Guesses))))
		} else {
			b.WriteString(common.ErrorStyle.Render(
				fmt.Sprintf("😔 The word was %s", m.Game.Target)))
		}
		b.WriteString("\n")
		if m.sharePrompt != nil {
			b.WriteString(m.sharePrompt.View())
		} else {
			b.WriteString(common.HelpStyle.Render("Press ESC to quit."))
		}
	} else {
		if m.err != "" {
			b.WriteString(common.ErrorStyle.Render(m.err))
		} else {
			b.WriteString(fmt.Sprintf("  > %s", m.input))
		}
		b.WriteString("\n")
		b.WriteString(common.HelpStyle.Render("Type a 5-letter word and press Enter. ESC to quit."))
	}

	b.WriteString("\n")
	return b.String()
}

func (m Model) renderGrid() string {
	cell := common.CellStyle(cellWidth, 1)
	var rows []string

	for i := 0; i < m.Game.MaxTurns; i++ {
		var cells []string
		if i < len(m.Game.Guesses) {
			// Filled row
			gr := m.Game.Guesses[i]
			for j, r := range gr.Word {
				s := stateStyle(gr.States[j]).
					Width(cellWidth).
					Align(lipgloss.Center)
				cells = append(cells, s.Render(string(r)))
			}
		} else if i == len(m.Game.Guesses) && !m.Game.IsOver() {
			// Current input row
			for j := 0; j < 5; j++ {
				if j < len(m.input) {
					cells = append(cells, cell.Render(string(m.input[j])))
				} else {
					cells = append(cells, common.EmptyStyle.Copy().
						Width(cellWidth).
						Align(lipgloss.Center).
						Render(" "))
				}
			}
		} else {
			// Empty row
			for j := 0; j < 5; j++ {
				cells = append(cells, common.EmptyStyle.Copy().
					Width(cellWidth).
					Align(lipgloss.Center).
					Render(" "))
			}
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Center, cells...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m Model) renderKeyboard() string {
	var rows []string
	for _, row := range keyboardRows {
		var keys []string
		for _, r := range row {
			state := m.Game.Keyboard[r]
			s := stateStyle(state).Width(3).Align(lipgloss.Center)
			keys = append(keys, s.Render(string(r)))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Center, keys...))
	}
	return lipgloss.JoinVertical(lipgloss.Center, rows...)
}

func stateStyle(s LetterState) lipgloss.Style {
	switch s {
	case Correct:
		return common.CorrectStyle
	case Present:
		return common.PresentStyle
	case Absent:
		return common.AbsentStyle
	default:
		return lipgloss.NewStyle()
	}
}

// Run creates and runs the Guess game TUI.
func Run(target string, validateWords bool) error {
	m := NewModel(target, validateWords)
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}

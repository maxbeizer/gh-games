package hang

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/maxbeizer/gh-games/internal/common"

	tea "github.com/charmbracelet/bubbletea"
)

// ASCII art gallows stages (0 = empty, 6 = dead).
var gallows = [7]string{
	// 0 — empty
	`
  ┌───┐
  │
  │
  │
  │
══╧══`,
	// 1 — head
	`
  ┌───┐
  │   O
  │
  │
  │
══╧══`,
	// 2 — body
	`
  ┌───┐
  │   O
  │   |
  │
  │
══╧══`,
	// 3 — left arm
	`
  ┌───┐
  │   O
  │  /|
  │
  │
══╧══`,
	// 4 — right arm
	`
  ┌───┐
  │   O
  │  /|\
  │
  │
══╧══`,
	// 5 — left leg
	`
  ┌───┐
  │   O
  │  /|\
  │  /
  │
══╧══`,
	// 6 — both legs (dead)
	`
  ┌───┐
  │   O
  │  /|\
  │  / \
  │
══╧══`,
}

var (
	gallowsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA"))

	wordDisplayStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF"))

	letterCorrect = lipgloss.NewStyle().
			Bold(true).
			Foreground(common.ColorGreen)

	letterWrong = lipgloss.NewStyle().
			Bold(true).
			Foreground(common.ColorRed)

	letterDim = lipgloss.NewStyle().
			Faint(true)

	statusStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("#B59F3B"))
)

// Model is the Bubbletea model for hangman.
type Model struct {
	Game        *Game
	status      string
	quitted     bool
	sharePrompt *common.SharePrompt
}

// NewModel creates a new hangman TUI model.
func NewModel() Model {
	return Model{Game: NewGame()}
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

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitted = true
			return m, tea.Quit
		case tea.KeyRunes:
			if m.Game.IsOver() {
				return m, nil
			}
			for _, r := range msg.Runes {
				correct, err := m.Game.GuessLetter(r)
				if err != nil {
					m.status = err.Error()
					return m, nil
				}
				if correct {
					if m.Game.IsWon() {
						m.status = "🎉 You got it!"
						sp := common.NewSharePrompt(m.Game.Summary())
						m.sharePrompt = &sp
					} else {
						m.status = fmt.Sprintf("✓ '%c' is in the word!", r)
					}
				} else {
					if m.Game.IsLost() {
						m.status = fmt.Sprintf("💀 The word was: %s", m.Game.Target)
						sp := common.NewSharePrompt(m.Game.Summary())
						m.sharePrompt = &sp
					} else {
						m.status = fmt.Sprintf("✗ No '%c' — %d/%d wrong", r, m.Game.WrongCount(), MaxWrong)
					}
				}
				break // only process first rune
			}
		case tea.KeyEnter:
			if m.Game.IsOver() {
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.quitted {
		return ""
	}

	var b strings.Builder

	// Title
	b.WriteString(common.TitleStyle.Render("☠️  HANGMAN"))
	b.WriteString("\n")

	// Gallows
	stage := m.Game.WrongCount()
	if stage > 6 {
		stage = 6
	}
	b.WriteString(gallowsStyle.Render(gallows[stage]))
	b.WriteString("\n\n")

	// Word display
	b.WriteString("  " + wordDisplayStyle.Render(m.Game.Display()))
	b.WriteString("\n\n")

	// Alphabet keyboard
	b.WriteString(m.renderKeyboard())
	b.WriteString("\n\n")

	// Status
	if m.status != "" {
		b.WriteString("  " + statusStyle.Render(m.status))
		b.WriteString("\n\n")
	}

	// Game over or help
	if m.Game.IsWon() {
		b.WriteString(common.SuccessStyle.Render("  🎉 You win!"))
		b.WriteString("\n")
		if m.sharePrompt != nil {
			b.WriteString("  ")
			b.WriteString(m.sharePrompt.View())
			b.WriteString("\n")
		} else {
			b.WriteString(common.HelpStyle.Render("  Press ENTER or ESC to exit."))
			b.WriteString("\n")
		}
	} else if m.Game.IsLost() {
		b.WriteString(common.ErrorStyle.Render("  💀 Game over!"))
		b.WriteString("\n")
		if m.sharePrompt != nil {
			b.WriteString("  ")
			b.WriteString(m.sharePrompt.View())
			b.WriteString("\n")
		} else {
			b.WriteString(common.HelpStyle.Render("  Press ENTER or ESC to exit."))
			b.WriteString("\n")
		}
	} else {
		b.WriteString(common.HelpStyle.Render("  Press a letter to guess · ESC to quit"))
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) renderKeyboard() string {
	rows := []string{"QWERTYUIOP", "ASDFGHJKL", "ZXCVBNM"}
	var rendered []string
	for _, row := range rows {
		var letters []string
		for _, r := range row {
			var s string
			if m.Game.IsCorrect(r) {
				s = letterCorrect.Render(string(r))
			} else if m.Game.IsWrongGuess(r) {
				s = letterWrong.Render(string(r))
			} else {
				s = letterDim.Render(string(r))
			}
			letters = append(letters, s)
		}
		rendered = append(rendered, "  "+strings.Join(letters, " "))
	}
	return strings.Join(rendered, "\n")
}

// Run starts the hangman TUI.
func Run() error {
	m := NewModel()
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}

package common

// SharePromptState tracks the share prompt lifecycle.
type SharePromptState int

const (
	ShareHidden    SharePromptState = iota // not showing
	SharePrompting                         // showing options
	ShareDone                              // showing confirmation
)

// SharePrompt is a reusable Bubbletea sub-model for post-game sharing.
type SharePrompt struct {
	State      SharePromptState
	Result     ShareResult
	Config     Config
	Message    string // confirmation message like "✓ Copied!"
	slackReady bool   // gh-slack installed AND channel configured
}

// NewSharePrompt creates a share prompt ready to display options.
func NewSharePrompt(result ShareResult) SharePrompt {
	cfg := LoadConfig()
	slackReady := IsGhSlackInstalled() && cfg.Share.SlackChannel != ""
	return SharePrompt{
		State:      SharePrompting,
		Result:     result,
		Config:     cfg,
		slackReady: slackReady,
	}
}

// HandleKey processes a keypress. Returns (updated prompt, should quit).
func (s SharePrompt) HandleKey(key string) (SharePrompt, bool) {
	switch s.State {
	case SharePrompting:
		switch key {
		case "c":
			err := CopyToClipboard(s.Result.String())
			if err != nil {
				s.Message = "✗ Could not copy: " + err.Error()
			} else {
				s.Message = "✓ Copied to clipboard!"
			}
			s.State = ShareDone
		case "s":
			if s.slackReady {
				err := PostViaGhSlack(s.Config.Share.SlackChannel, s.Config.Share.SlackTeam, s.Result.String())
				if err != nil {
					s.Message = "✗ " + err.Error()
				} else {
					s.Message = "✓ Posted to #" + s.Config.Share.SlackChannel + "!"
				}
				s.State = ShareDone
			}
		case "n", "esc":
			return s, true
		}
	case ShareDone:
		return s, true // any key quits after confirmation
	}
	return s, false
}

// View renders the share prompt.
func (s SharePrompt) View() string {
	switch s.State {
	case SharePrompting:
		line := "Share results? [C]opy to clipboard"
		if s.slackReady {
			line += " · [S]lack"
		}
		line += " · [N]o"
		return HelpStyle.Render(line)
	case ShareDone:
		return SuccessStyle.Render(s.Message) + "\n" + HelpStyle.Render("Press any key to exit")
	default:
		return ""
	}
}

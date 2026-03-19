package common

import "fmt"

// SharePromptState tracks the share prompt lifecycle.
type SharePromptState int

const (
	ShareHidden    SharePromptState = iota // not showing
	SharePrompting                         // showing options
	ShareDone                              // showing confirmation
)

// SharePrompt is a reusable Bubbletea sub-model for post-game sharing.
type SharePrompt struct {
	State   SharePromptState
	Result  ShareResult
	Config  Config
	Message string // confirmation message like "✓ Copied!"
}

// NewSharePrompt creates a share prompt ready to display options.
func NewSharePrompt(result ShareResult) SharePrompt {
	cfg := LoadConfig()
	return SharePrompt{
		State:  SharePrompting,
		Result: result,
		Config: cfg,
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
		case "n", "esc":
			return s, true
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			idx := int(key[0] - '1')
			if idx < len(s.Config.Share.Webhooks) {
				wh := s.Config.Share.Webhooks[idx]
				err := PostToWebhook(wh.URL, s.Result.String())
				if err != nil {
					s.Message = "✗ Failed: " + err.Error()
				} else {
					s.Message = "✓ Posted to " + wh.Name + "!"
				}
				s.State = ShareDone
			}
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
		for i, wh := range s.Config.Share.Webhooks {
			line += fmt.Sprintf(" · [%d] %s", i+1, wh.Name)
		}
		line += " · [N]o"
		return HelpStyle.Render(line)
	case ShareDone:
		return SuccessStyle.Render(s.Message) + "\n" + HelpStyle.Render("Press any key to exit")
	default:
		return ""
	}
}

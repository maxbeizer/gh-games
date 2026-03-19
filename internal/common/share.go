package common

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// ShareResult holds the data needed to render a shareable game result.
type ShareResult struct {
	Game  string   // emoji + name, e.g. "🟩 Guess"
	Title string   // headline, e.g. "🟩 Guess 4/6"
	Lines []string // body lines (emoji grid, scores, etc.)
}

// String renders the full shareable text with attribution line.
func (r ShareResult) String() string {
	var b strings.Builder
	b.WriteString(r.Title)
	b.WriteByte('\n')
	for _, l := range r.Lines {
		b.WriteString(l)
		b.WriteByte('\n')
	}
	b.WriteByte('\n')
	b.WriteString("gh games · github.com/maxbeizer/gh-games")
	return b.String()
}

// CopyToClipboard copies text to the system clipboard.
// macOS: pbcopy, Linux: xclip -selection clipboard, Windows: clip.exe.
// Returns nil on success, error if the clipboard tool is not found.
func CopyToClipboard(text string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case "windows":
		cmd = exec.Command("clip.exe")
	default:
		return fmt.Errorf("unsupported OS for clipboard: %s", runtime.GOOS)
	}
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// IsGhSlackInstalled reports whether the gh-slack extension is available.
func IsGhSlackInstalled() bool {
	_, err := exec.LookPath("gh")
	if err != nil {
		return false
	}
	out, err := exec.Command("gh", "extension", "list").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "gh-slack")
}

// PostViaGhSlack sends a message to a Slack channel using the gh-slack extension.
func PostViaGhSlack(channel, team, message string) error {
	channel = strings.TrimPrefix(channel, "#")
	args := []string{"slack", "send", "-c", channel, "-m", message}
	if team != "" {
		args = append(args, "-t", team)
	}
	cmd := exec.Command("gh", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gh-slack send failed: %s: %w", string(output), err)
	}
	return nil
}

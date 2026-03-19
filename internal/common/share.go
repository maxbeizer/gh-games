package common

import (
	"fmt"
	"net/http"
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

// WebhookFormat determines the payload format from a URL.
type WebhookFormat int

const (
	FormatSlack   WebhookFormat = iota // {"text": "..."}
	FormatDiscord                      // {"content": "..."}
	FormatGeneric                      // {"text": "..."}
)

// DetectFormat guesses the webhook payload format from the URL.
func DetectFormat(url string) WebhookFormat {
	if strings.Contains(url, "discord.com") {
		return FormatDiscord
	}
	return FormatSlack
}

// PostToWebhook POSTs the share text to a webhook URL with auto-detected format.
func PostToWebhook(url, text string) error {
	format := DetectFormat(url)

	var key string
	switch format {
	case FormatDiscord:
		key = "content"
	default:
		key = "text"
	}

	// Escape for JSON: backslash, double-quote, newlines.
	escaped := strings.ReplaceAll(text, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `"`, `\"`)
	escaped = strings.ReplaceAll(escaped, "\n", `\n`)
	payload := fmt.Sprintf(`{"%s": "%s"}`, key, escaped)

	resp, err := http.Post(url, "application/json", strings.NewReader(payload)) //nolint:gosec
	if err != nil {
		return fmt.Errorf("webhook POST failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}
	return nil
}

package common

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestShareResultString(t *testing.T) {
	r := ShareResult{
		Game:  "🟩 Guess",
		Title: "🟩 Guess 4/6",
		Lines: []string{"🟩🟨⬛⬛⬛", "🟩🟩🟩⬛⬛", "🟩🟩🟩🟩⬛", "🟩🟩🟩🟩🟩"},
	}
	got := r.String()

	if !strings.HasPrefix(got, "🟩 Guess 4/6\n") {
		t.Errorf("expected title as first line, got:\n%s", got)
	}
	for _, line := range r.Lines {
		if !strings.Contains(got, line) {
			t.Errorf("missing body line %q", line)
		}
	}
	const attribution = "gh games · github.com/maxbeizer/gh-games"
	if !strings.HasSuffix(got, attribution) {
		t.Errorf("missing attribution line, got:\n%s", got)
	}
	if !strings.Contains(got, attribution) {
		t.Errorf("attribution not found anywhere in output")
	}
}

func TestShareResultString_EmptyLines(t *testing.T) {
	r := ShareResult{
		Game:  "🎯 Darts",
		Title: "🎯 Darts 100",
		Lines: nil,
	}
	got := r.String()

	if !strings.HasPrefix(got, "🎯 Darts 100\n") {
		t.Errorf("expected title as first line, got:\n%s", got)
	}
	const attribution = "gh games · github.com/maxbeizer/gh-games"
	if !strings.HasSuffix(got, attribution) {
		t.Errorf("expected attribution at end, got:\n%s", got)
	}
	// Should not contain excessive blank lines — just title + blank + attribution.
	lines := strings.Split(got, "\n")
	if len(lines) < 2 {
		t.Errorf("expected at least 2 lines (title + attribution), got %d", len(lines))
	}
}

func TestCopyToClipboard_Smoke(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("clipboard test only runs on macOS")
	}
	// Smoke test: just verify it doesn't panic.
	err := CopyToClipboard("gh-games test")
	if err != nil {
		t.Errorf("CopyToClipboard returned error: %v", err)
	}
}

func TestIsGhSlackInstalled_NoPanic(t *testing.T) {
	// Just verify it runs without panicking; result depends on environment.
	_ = IsGhSlackInstalled()
}

func TestLoadConfig_MissingFile(t *testing.T) {
	cfg := loadConfigFrom("/nonexistent/path/config.toml")
	if cfg.Share.SlackChannel != "" {
		t.Errorf("expected empty slack_channel, got %q", cfg.Share.SlackChannel)
	}
}

func TestLoadConfig_ParsesSlackChannel(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	content := `# gh-games config

[share]
slack_channel = "fun-games"
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := loadConfigFrom(path)

	if cfg.Share.SlackChannel != "fun-games" {
		t.Errorf("slack_channel = %q, want %q", cfg.Share.SlackChannel, "fun-games")
	}
}

func TestLoadConfig_MalformedFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	// Garbled content — should not panic, just return zero-value or partial config.
	content := `not valid toml at all }{{{
=== garbage ===
slack_channel missing equals
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	// Must not panic.
	cfg := loadConfigFrom(path)

	// We don't require specific output, just graceful handling.
	_ = cfg
}

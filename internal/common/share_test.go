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

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want WebhookFormat
	}{
		{"slack webhook", "https://hooks.slack.com/services/T00/B00/xxxx", FormatSlack},
		{"discord webhook", "https://discord.com/api/webhooks/123/abc", FormatDiscord},
		{"generic URL", "https://example.com/webhook", FormatSlack},
		{"empty string", "", FormatSlack},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectFormat(tt.url); got != tt.want {
				t.Errorf("DetectFormat(%q) = %d, want %d", tt.url, got, tt.want)
			}
		})
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

func TestLoadConfig_MissingFile(t *testing.T) {
	cfg := loadConfigFrom("/nonexistent/path/config.toml")
	if len(cfg.Share.Webhooks) != 0 {
		t.Errorf("expected no webhooks, got %d", len(cfg.Share.Webhooks))
	}
	if cfg.Share.Default != "" {
		t.Errorf("expected empty default, got %q", cfg.Share.Default)
	}
}

func TestLoadConfig_ParsesWebhooks(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	content := `# gh-games config

[[share.webhooks]]
name = "slack"
url = "https://hooks.slack.com/services/T/B/x"

[[share.webhooks]]
name = "discord"
url = "https://discord.com/api/webhooks/123/abc"
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := loadConfigFrom(path)

	if len(cfg.Share.Webhooks) != 2 {
		t.Fatalf("expected 2 webhooks, got %d", len(cfg.Share.Webhooks))
	}
	if cfg.Share.Webhooks[0].Name != "slack" {
		t.Errorf("webhook[0].Name = %q, want slack", cfg.Share.Webhooks[0].Name)
	}
	if cfg.Share.Webhooks[0].URL != "https://hooks.slack.com/services/T/B/x" {
		t.Errorf("webhook[0].URL = %q", cfg.Share.Webhooks[0].URL)
	}
	if cfg.Share.Webhooks[1].Name != "discord" {
		t.Errorf("webhook[1].Name = %q, want discord", cfg.Share.Webhooks[1].Name)
	}
	if cfg.Share.Webhooks[1].URL != "https://discord.com/api/webhooks/123/abc" {
		t.Errorf("webhook[1].URL = %q", cfg.Share.Webhooks[1].URL)
	}
}

func TestLoadConfig_MalformedFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	// Garbled content — should not panic, just return zero-value or partial config.
	content := `not valid toml at all }{{{
=== garbage ===
[[share.webhooks
name "missing equals"
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	// Must not panic.
	cfg := loadConfigFrom(path)

	// We don't require specific output, just graceful handling.
	_ = cfg
}

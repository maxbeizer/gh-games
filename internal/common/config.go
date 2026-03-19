package common

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ShareConfig holds sharing preferences.
type ShareConfig struct {
	SlackChannel string // channel name for gh-slack
	SlackTeam    string // Slack workspace/team name for gh-slack
}

// Config is the top-level application configuration.
type Config struct {
	Share ShareConfig
}

// ConfigPath returns the path to the config file.
func ConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "gh-games", "config.toml")
}

// LoadConfig loads config from ~/.config/gh-games/config.toml.
// Returns a zero-value Config if the file doesn't exist.
func LoadConfig() Config {
	return loadConfigFrom(ConfigPath())
}

// SaveConfig writes config to ~/.config/gh-games/config.toml.
func SaveConfig(cfg Config) error {
	path := ConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	var b strings.Builder
	b.WriteString("[share]\n")
	if cfg.Share.SlackChannel != "" {
		b.WriteString(fmt.Sprintf("slack_channel = \"%s\"\n", cfg.Share.SlackChannel))
	}
	if cfg.Share.SlackTeam != "" {
		b.WriteString(fmt.Sprintf("slack_team = \"%s\"\n", cfg.Share.SlackTeam))
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

// DeleteConfig removes the config file.
func DeleteConfig() error {
	return os.Remove(ConfigPath())
}

// loadConfigFrom parses a minimal TOML config from the given path.
func loadConfigFrom(path string) Config {
	f, err := os.Open(path)
	if err != nil {
		return Config{}
	}
	defer f.Close()

	var cfg Config
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, val, ok := parseTOMLKV(line)
		if !ok {
			continue
		}

		if key == "slack_channel" {
			cfg.Share.SlackChannel = val
		}
		if key == "slack_team" {
			cfg.Share.SlackTeam = val
		}
	}

	return cfg
}

// parseTOMLKV extracts key and unquoted value from a "key = "value"" line.
func parseTOMLKV(line string) (key, val string, ok bool) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	key = strings.TrimSpace(parts[0])
	raw := strings.TrimSpace(parts[1])
	// Remove surrounding quotes if present.
	if len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
	}
	return key, raw, true
}

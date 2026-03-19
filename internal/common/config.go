package common

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// ShareConfig holds sharing preferences.
type ShareConfig struct {
	SlackChannel string // channel name for gh-slack
}

// Config is the top-level application configuration.
type Config struct {
	Share ShareConfig
}

// LoadConfig loads config from ~/.config/gh-games/config.toml.
// Returns a zero-value Config if the file doesn't exist.
func LoadConfig() Config {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}
	}
	return loadConfigFrom(filepath.Join(home, ".config", "gh-games", "config.toml"))
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

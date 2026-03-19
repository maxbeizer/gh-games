package common

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// WebhookConfig holds a named webhook destination.
type WebhookConfig struct {
	Name string
	URL  string
}

// ShareConfig holds sharing preferences.
type ShareConfig struct {
	Default  string          // "clipboard" or a webhook name
	Webhooks []WebhookConfig
}

// Config is the top-level application configuration.
type Config struct {
	Share ShareConfig
}

// LoadConfig loads config from ~/.config/gh-games/config.toml.
// Returns a zero-value Config (no webhooks) if the file doesn't exist.
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
	var cur *WebhookConfig

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if line == "[[share.webhooks]]" {
			// Flush previous entry.
			if cur != nil {
				cfg.Share.Webhooks = append(cfg.Share.Webhooks, *cur)
			}
			cur = &WebhookConfig{}
			continue
		}

		key, val, ok := parseTOMLKV(line)
		if !ok {
			continue
		}

		if cur != nil {
			switch key {
			case "name":
				cur.Name = val
			case "url":
				cur.URL = val
			}
		} else {
			// Top-level share keys.
			if key == "share.default" || key == "default" {
				cfg.Share.Default = val
			}
		}
	}

	// Flush last entry.
	if cur != nil {
		cfg.Share.Webhooks = append(cfg.Share.Webhooks, *cur)
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

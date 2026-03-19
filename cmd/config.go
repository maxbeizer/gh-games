package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/maxbeizer/gh-games/internal/common"
	"github.com/spf13/cobra"
)

// NewConfigCmd returns the config subcommand.
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "⚙️ Configure gh-games settings",
		Long:  "Interactive setup for gh-games configuration.\nCurrently configures Slack sharing via gh-slack.",
		RunE: func(cmd *cobra.Command, args []string) error {
			show, _ := cmd.Flags().GetBool("show")
			reset, _ := cmd.Flags().GetBool("reset")

			if show {
				return showConfig()
			}
			if reset {
				return resetConfig()
			}
			return interactiveConfig()
		},
	}
	cmd.Flags().Bool("show", false, "Show current configuration")
	cmd.Flags().Bool("reset", false, "Delete configuration file")
	return cmd
}

func interactiveConfig() error {
	cfgPath := common.ConfigPath()
	fmt.Println("⚙️  gh-games configuration")
	fmt.Printf("Current config file: %s\n\n", cfgPath)
	fmt.Println("── Slack Sharing ──")

	if common.IsGhSlackInstalled() {
		fmt.Println("✓ gh-slack is installed")
	} else {
		fmt.Println("✗ gh-slack is not installed (install with: gh extension install github/gh-slack)")
	}

	fmt.Print("Slack team/workspace name (e.g. github, leave empty to skip): ")
	scanner := bufio.NewScanner(os.Stdin)
	var team string
	if scanner.Scan() {
		team = scanner.Text()
	}

	fmt.Print("Slack channel name (leave empty to skip): ")
	var channel string
	if scanner.Scan() {
		channel = scanner.Text()
	}

	if channel == "" && team == "" {
		fmt.Println("No Slack config set — skipping.")
		return nil
	}

	cfg := common.LoadConfig()
	cfg.Share.SlackChannel = channel
	cfg.Share.SlackTeam = team
	if err := common.SaveConfig(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	fmt.Printf("✓ Config saved to %s\n", cfgPath)
	return nil
}

func showConfig() error {
	cfgPath := common.ConfigPath()
	if _, err := os.Stat(cfgPath); errors.Is(err, os.ErrNotExist) {
		fmt.Println("No config file found. Run 'gh games config' to create one.")
		return nil
	}

	cfg := common.LoadConfig()
	fmt.Printf("⚙️  gh-games config (%s)\n\n", cfgPath)

	ch := cfg.Share.SlackChannel
	if ch == "" {
		ch = "(not set)"
	}
	tm := cfg.Share.SlackTeam
	if tm == "" {
		tm = "(not set)"
	}
	fmt.Printf("Slack team:    %s\n", tm)
	fmt.Printf("Slack channel: %s\n", ch)

	if common.IsGhSlackInstalled() {
		fmt.Println("gh-slack: ✓ installed")
	} else {
		fmt.Println("gh-slack: ✗ not installed")
	}
	return nil
}

func resetConfig() error {
	cfgPath := common.ConfigPath()
	if _, err := os.Stat(cfgPath); errors.Is(err, os.ErrNotExist) {
		fmt.Println("No config file to delete.")
		return nil
	}
	if err := common.DeleteConfig(); err != nil {
		return fmt.Errorf("deleting config: %w", err)
	}
	fmt.Println("✓ Config deleted")
	return nil
}

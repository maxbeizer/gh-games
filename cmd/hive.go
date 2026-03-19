package cmd

import (
	"github.com/maxbeizer/gh-games/internal/hive"
	"github.com/spf13/cobra"
)

func NewHiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hive",
		Short: "🐝 Find words from 7 letters — center letter required",
		Long: `Find as many words as you can using 7 letters.

Every word must include the center letter and be at least 4 letters long.
Letters can be reused. Find the pangram (uses all 7 letters) for bonus points!

Ranks: Beginner → Good → Nice → Great → Amazing → Genius → Queen Bee`,
		Example: "gh games hive",
		RunE: func(cmd *cobra.Command, args []string) error {
			return hive.Run()
		},
	}
	return cmd
}

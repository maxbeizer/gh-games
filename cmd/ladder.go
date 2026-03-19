package cmd

import (
	"github.com/maxbeizer/gh-games/internal/ladder"
	"github.com/spf13/cobra"
)

// NewLadderCmd creates the ladder subcommand.
func NewLadderCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ladder",
		Short: "🪜 Change one letter at a time to reach the target word",
		Long: `Word Ladder — change one letter at a time to transform the start
word into the target word. Every intermediate step must be a valid
English word. Try to match the optimal path length!

Example: COLD → CORD → CORE → BORE → BORN`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return ladder.Run()
		},
	}
}

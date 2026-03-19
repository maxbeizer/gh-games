package cmd

import (
	"github.com/maxbeizer/gh-games/internal/jumble"
	"github.com/spf13/cobra"
)

func NewJumbleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jumble",
		Short: "🔀 Unscramble the jumbled word",
		Long: `Unscramble jumbled words in 5 rounds of increasing difficulty!

Each round presents a scrambled word that gets progressively longer
(4 → 5 → 6 → 7 → 8 letters). Type the correct word to score points.

Scoring:
  • Base points scale with word length (100 per letter)
  • Speed bonus for solving quickly (under 10 seconds)
  • Hint penalty: -50 points per hint used
  • Minimum 10 points for every correct answer

Controls:
  ? — Reveal a letter in its correct position (costs points)
  Tab — Re-shuffle the scrambled letters
  ESC — Quit the game`,
		Example: "gh games jumble",
		RunE: func(cmd *cobra.Command, args []string) error {
			return jumble.Run()
		},
	}

	return cmd
}

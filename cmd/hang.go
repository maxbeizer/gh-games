package cmd

import (
	"github.com/maxbeizer/gh-games/internal/hang"
	"github.com/spf13/cobra"
)

// NewHangCmd returns the cobra command for the hangman game.
func NewHangCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "hang",
		Short: "☠️ Classic hangman — guess the word letter by letter",
		Long: `Classic hangman game — guess the hidden word one letter at a time.

Each wrong guess adds a body part to the gallows. You get 6 wrong
guesses before the hangman is complete and the game is over.

  🟢 Green  – letter is in the word
  🔴 Red    – letter is not in the word`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return hang.Run()
		},
	}
}

package cmd

import (
	"github.com/maxbeizer/gh-games/internal/guess"
	"github.com/spf13/cobra"
)

func NewGuessCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "guess",
		Short: "🟩 Guess the 5-letter word in 6 tries",
		Long: `Guess the hidden 5-letter word in 6 attempts.

After each guess the letters are colored to show how close you are:
  🟩 Green  – correct letter in the correct position
  🟨 Yellow – correct letter in the wrong position
  ⬜ Gray   – letter is not in the word

By default a date-seeded "daily" word is used so everyone gets the
same puzzle each day. Use --random for a fresh word every play.`,
		Example: "gh games guess\ngh games guess --random\ngh games guess --hard",
		RunE: func(cmd *cobra.Command, args []string) error {
			random, _ := cmd.Flags().GetBool("random")
			hard, _ := cmd.Flags().GetBool("hard")

			var target string
			if random {
				target = guess.RandomWord()
			} else {
				target = guess.DailyWord()
			}

			return guess.Run(target, hard)
		},
	}

	cmd.Flags().Bool("daily", true, "Use the daily word (same word for everyone today)")
	cmd.Flags().Bool("random", false, "Use a fully random word each play")
	cmd.Flags().Bool("hard", false, "Hard mode: guesses must be valid words from the list")

	return cmd
}

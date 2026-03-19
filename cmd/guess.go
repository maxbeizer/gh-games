package cmd

import "github.com/spf13/cobra"

func NewGuessCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "guess",
		Short: "🟩 Guess the 5-letter word in 6 tries",
	}
}

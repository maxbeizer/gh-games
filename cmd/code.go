package cmd

import (
	"github.com/maxbeizer/gh-games/internal/code"
	"github.com/spf13/cobra"
)

// NewCodeCmd returns the cobra command for the Code Breaker game.
func NewCodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "code",
		Short: "🔐 Crack the secret color code in 10 guesses",
		Long: `Code Breaker — a Mastermind-style logic game.

A secret code of 4 colored pegs is generated from 6 possible colors.
You have 10 guesses to crack it. After each guess you get feedback:
  ● = right color, right position
  ○ = right color, wrong position

Use deduction to narrow down the code!`,
		Example: "  gh games code",
		RunE: func(cmd *cobra.Command, args []string) error {
			return code.Run()
		},
	}

	return cmd
}

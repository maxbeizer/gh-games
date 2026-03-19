package cmd

import (
	"github.com/maxbeizer/gh-games/internal/trivia"
	"github.com/spf13/cobra"
)

func NewTriviaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trivia",
		Short: "🧠 Test your knowledge with 10 trivia questions",
		Long: `Answer 10 multiple-choice trivia questions from categories like
programming, tech, science, Git, and fun nerd culture.

Use A/B/C/D or arrow keys to select your answer, then press Enter.
See your final score and review any wrong answers at the end!`,
		Example: "gh games trivia",
		RunE: func(cmd *cobra.Command, args []string) error {
			return trivia.Run()
		},
	}
	return cmd
}

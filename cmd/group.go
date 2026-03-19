package cmd

import (
	"github.com/maxbeizer/gh-games/internal/group"
	"github.com/spf13/cobra"
)

func NewGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
		Short: "🔗 Find four groups of four related words",
		Long: `Find four hidden groups among sixteen words.

Select four words that share a connection and submit your guess.
Each group has a difficulty level shown by color:
  🟨 Yellow – easiest
  🟩 Green  – medium
  🟦 Blue   – hard
  🟪 Purple – expert

You get 4 mistakes before it's game over. Good luck!`,
		Example: "gh games group",
		RunE: func(cmd *cobra.Command, args []string) error {
			return group.Run()
		},
	}

	return cmd
}

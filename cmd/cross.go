package cmd

import (
	"github.com/maxbeizer/gh-games/internal/cross"
	"github.com/spf13/cobra"
)

// NewCrossCmd returns the cobra command for the crossword game.
func NewCrossCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cross",
		Short: "📰 Solve a 5×5 mini crossword puzzle",
		Long: `Solve a 5×5 mini crossword puzzle with across and down clues.

Navigate the grid with arrow keys, type letters to fill cells.
Press Tab to toggle between across and down mode.
Press Ctrl+K to check your answers (wrong letters shown in red).
The puzzle auto-completes when all letters are correct!`,
		Example: "gh games cross",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cross.Run()
		},
	}
	return cmd
}

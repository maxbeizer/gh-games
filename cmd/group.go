package cmd

import "github.com/spf13/cobra"

func NewGroupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "group",
		Short: "🔗 Find four groups of four related words",
	}
}

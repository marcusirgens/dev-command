package cmd

import (
	"github.com/spf13/cobra"
)

// highlightCmd represents the highlight command
var highlightCmd = &cobra.Command{
	Use:   "highlight",
	Short: "A set of syntax highlighting commands",
}

func init() {
	rootCmd.AddCommand(highlightCmd)
}

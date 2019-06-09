package cmd

import (
	"github.com/spf13/cobra"
)

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Syntax highlighting for git",
}

func init() {
	highlightCmd.AddCommand(gitCmd)
}

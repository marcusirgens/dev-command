package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

// zshCompletionCmd represents the zshCompletion command
var zshCompletionCmd = &cobra.Command{
	Use:    "zsh-completion",
	Short:  "Generates zsh completion",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.GenZshCompletion(os.Stdout)
		if err != nil {
			cmd.PrintErrln(err)
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(zshCompletionCmd)

}

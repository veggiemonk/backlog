package cmd

import (
	"fmt"

	"github.com/imjasonh/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print the version information",
	Example: VersionExamples.GenerateExampleText(),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Backlog version:\n%s\n", version.Get().String())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

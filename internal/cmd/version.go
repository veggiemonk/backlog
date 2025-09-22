package cmd

import (
	"fmt"

	"github.com/imjasonh/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Example: `
backlog version # Print the version information

# Example output:
#
# Backlog version:
# Revision: 7c989dabd2c61a063a23788c18eb39eca408f6a7
# Version: v0.0.2-0.20250907193624-7c989dabd2c6
# BuildTime: 2025-09-07T19:36:24Z
# Dirty: false
`,

	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Backlog version:\n%s\n", version.Get().String())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

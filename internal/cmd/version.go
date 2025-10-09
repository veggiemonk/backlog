package cmd

import (
	"context"
	"fmt"

	"github.com/imjasonh/version"
	"github.com/urfave/cli/v3"
)

const versionExamples = `
backlog version # Print the version information

# Example output:
#
# Backlog version:
# Revision: 7c989dabd2c61a063a23788c18eb39eca408f6a7
# Version: v0.0.2-0.20250907193624-7c989dabd2c6
# BuildTime: 2025-09-07T19:36:24Z
# Dirty: false
`

func newVersionCommand() *cli.Command {
	return &cli.Command{
		Name:        "version",
		Usage:       "Print the version information",
		Description: "Print the version information.\n\nExamples:\n" + versionExamples,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Fprintf(cmd.Root().Writer, "Backlog version:\n%s\n", version.Get().String())
			return nil
		},
	}
}

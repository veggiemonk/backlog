package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/veggiemonk/backlog/internal/mcp"
)

const instructionsExamples = `
backlog instructions               # outputs the instructions for agents to use the cli.
backlog instructions --mode cli    # outputs the instructions for agents to use the cli.
backlog instructions --mode mcp    # outputs the instructions for agents to use MCP.
backlog instructions >> AGENTS.md  # add instructions to agent base prompt.
`

func newInstructionsCommand() *cli.Command {
	return &cli.Command{
		Name:  "instructions",
		Usage: "instructions for agents to learn to use backlog",
		Description: "Instructions for agents to learn to use backlog by including them into a prompt.\n\nExamples:\n" +
			instructionsExamples,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "mode", Value: "cli", Usage: "which mode the agent will use backlog: (cli|mcp)"},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			switch cmd.String("mode") {
			case "cli":
				fmt.Fprintf(cmd.Root().Writer, "%s\n", mcp.PromptCLIInstructions)
			case "mcp":
				fmt.Fprintf(cmd.Root().Writer, "%s\n", mcp.PromptMCPInstructions)
			default:
				fmt.Fprintf(cmd.Root().Writer, "%s\n", mcp.PromptCLIInstructions)
			}
			return nil
		},
	}
}

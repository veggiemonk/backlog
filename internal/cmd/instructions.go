package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/veggiemonk/backlog/internal/mcp"
)

var instructionsCmd = &cobra.Command{
	Use:   "instructions",
	Short: "instructions for agents to learn to use backlog",
	Long:  `Instructions for agents to learn to use backlog by including them into a prompt.`,
	Example: `
backlog instructions               # outputs the instructions for agents to use the cli.
backlog instructions --mode cli    # outputs the instructions for agents to use the cli.
backlog instructions --mode mcp    # outputs the instructions for agents to use MCP.
backlog instructions >> AGENTS.md  # add instructions to agent base prompt.
`,
	RunE: runInstructions,
}

func init() {
	rootCmd.AddCommand(instructionsCmd)
	setInstructionsFlags(instructionsCmd)
}

var instructionMode string

func setInstructionsFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&instructionMode, "mode", "cli", "which mode the agent will use backlog: (cli|mcp)")
}

func runInstructions(_ *cobra.Command, _ []string) error {
	switch instructionMode {
	case "cli":
		fmt.Printf("%s\n", mcp.PromptCLIInstructions)
	case "mcp":
		fmt.Printf("%s\n", mcp.PromptMCPInstructions)
	default:
		fmt.Printf("%s\n", mcp.PromptCLIInstructions)
	}
	return nil
}

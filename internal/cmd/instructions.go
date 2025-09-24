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
backlog instructions               # outputs the instructions 
backlog instructions >> AGENTS.md  # add instructions to agent base prompt.
`,
	RunE: runInstructions,
}

func init() {
	rootCmd.AddCommand(instructionsCmd)
}

func runInstructions(_ *cobra.Command, _ []string) error {
	fmt.Printf("%s\n", mcp.PromptInstructions)
	return nil
}

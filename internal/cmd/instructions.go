package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/veggiemonk/backlog/internal/mcp"
)

var instructionsCmd = &cobra.Command{
	Use:     "instructions",
	Short:   "instructions for agents to learn to use backlog",
	Long:    `Instructions for agents to learn to use backlog by including them into a prompt.`,
	Example: InstructionsExamples.GenerateExampleText(),
	Run:     runInstructions,
}

func init() {
	rootCmd.AddCommand(instructionsCmd)
}

func runInstructions(cmd *cobra.Command, args []string) {
	fmt.Printf("%s\n", mcp.PromptInstructions)
}

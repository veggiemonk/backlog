package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/cobra"
	"github.com/veggiemonk/backlog/internal/mcp"
)

var instructionsCmd = &cobra.Command{
	Use:     "instructions",
	Short:   "instructions for agents to learn to use backlog",
	Long:    `Instructions for agents to learn to use backlog by including them into a prompt.`,
	Example: generateExampleText(InstructionsExamples),
	Run:     runInstructions,
}

func init() {
	rootCmd.AddCommand(instructionsCmd)
}

func runInstructions(cmd *cobra.Command, args []string) {
	fmt.Printf("%s\n", mcp.PromptInstructions)
}

func outputNonEmpty(t *testing.T, output []byte, _ error) {
	is := is.New(t)
	o := string(output)
	is.True(strings.Contains(o, "backlog"))
}

// InstructionsExamples contains all examples for the instructions command
var InstructionsExamples = CommandExamples{
	Examples: []CommandExample{
		{
			Description:     "Output Instructions",
			Command:         "backlog instructions",
			Comment:         "outputs the instructions",
			OutputValidator: outputNonEmpty,
		},
		{
			Description:     "Save Instructions to File",
			Command:         "backlog instructions >> AGENTS.md",
			Comment:         "add instructions to agent base prompt",
			OutputValidator: outputNonEmpty,
		},
	},
}

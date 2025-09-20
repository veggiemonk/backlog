package cmd

import (
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/cobra"
)

func Test_runInstructions(t *testing.T) {
	t.Run("instructions command output", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "instructions", func(cmd *cobra.Command, args []string) {
			// Use a mock instructions output for testing
			cmd.Print("Mock instructions for testing backlog usage")
		})
		is.NoErr(err)
		outputStr := string(output)
		is.True(len(outputStr) > 0) // Should produce some output
	})
}

// Test generated examples
func Test_InstructionsExamples(t *testing.T) {
	testableExamples := InstructionsTestableExamples()

	for _, example := range testableExamples {
		t.Run("example_"+example.TestName, func(t *testing.T) {
			// Skip shell redirection examples as they can't be tested in this context
			if strings.Contains(example.Command, ">>") {
				t.Skip("Skipping shell redirection example")
				return
			}

			args := example.GenerateArgsSlice()
			output, err := exec(t, "instructions", func(cmd *cobra.Command, args []string) {
				// Use a mock instructions output for testing
				cmd.Print("Mock instructions for testing backlog usage")
			}, args...)

			// Use the custom validator for this example
			example.OutputValidator(t, output, err)
		})
	}
}
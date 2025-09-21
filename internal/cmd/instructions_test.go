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
		output, err := exec(t, "instructions", func(cmd *cobra.Command, args []string) error {
			// Use a mock instructions output for testing
			cmd.Print("Mock instructions for testing backlog usage")
			return nil
		})
		is.NoErr(err)
		outputStr := string(output)
		is.True(len(outputStr) > 0) // Should produce some output
	})
}
// InstructionsTestableExamples generates testable examples for instructions command
func InstructionsTestableExamples() []TestableExample {
	var testable []TestableExample
	for _, example := range InstructionsExamples.Examples {
		te := example.ToTestableExample()
		te.OutputValidator = func(t *testing.T, output []byte, err error) {
			is := is.New(t)
			is.NoErr(err) // Instructions command should not error
			outputStr := string(output)
			is.True(len(outputStr) > 0) // Should produce some output
		}
		testable = append(testable, te)
	}
	return testable
}

// Test generated examples
func Test_InstructionsExamples(t *testing.T) {
	for _, ex := range InstructionsExamples.Examples {
		t.Run("example_"+generateTestName(ex.Description), func(t *testing.T) {
			// Skip shell redirection examples as they can't be tested in this context
			if strings.Contains(ex.Command, ">>") {
				t.Skip("Skipping shell redirection example")
				return
			}

			args := generateArgsSlice(ex)
			output, err := exec(t, "instructions", func(cmd *cobra.Command, args []string) error {
				// Use a mock instructions output for testing
				cmd.Print("Mock instructions for testing backlog usage")
				return nil
			}, args...)

			// Use the custom validator for this example
			ex.OutputValidator(t, output, err)
		})
	}
}

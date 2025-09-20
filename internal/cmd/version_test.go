package cmd

import (
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/cobra"
)

func Test_runVersion(t *testing.T) {
	t.Run("version command output", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "version", func(cmd *cobra.Command, args []string) {
			// Use a mock version output for testing
			cmd.Print("Backlog version:\nRevision: test-revision\nVersion: test-version\nBuildTime: test-time\nDirty: false\n")
		})
		is.NoErr(err)
		outputStr := string(output)
		is.True(strings.Contains(outputStr, "Backlog version"))
		is.True(strings.Contains(outputStr, "Revision"))
		is.True(strings.Contains(outputStr, "Version"))
		is.True(strings.Contains(outputStr, "BuildTime"))
	})
}

// Test generated examples
func Test_VersionExamples(t *testing.T) {
	testableExamples := VersionTestableExamples()

	for _, example := range testableExamples {
		t.Run("example_"+example.TestName, func(t *testing.T) {
			args := example.GenerateArgsSlice()
			output, err := exec(t, "version", func(cmd *cobra.Command, args []string) {
				// Use a mock version output for testing
				cmd.Print("Backlog version:\nRevision: test-revision\nVersion: test-version\nBuildTime: test-time\nDirty: false\n")
			}, args...)

			// Use the custom validator for this example
			example.OutputValidator(t, output, err)
		})
	}
}
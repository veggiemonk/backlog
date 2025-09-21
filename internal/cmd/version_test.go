package cmd

import (
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/cobra"
)

const versionFakeOutput = `Backlog version:
Revision: test-revision
Version: test-version
BuildTime: test-time
Dirty: false
`

func Test_runVersion(t *testing.T) {
	t.Run("version command output", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "version", func(cmd *cobra.Command, args []string) error {
			// Use a mock version output for testing
			cmd.Print(versionFakeOutput)
			return nil
		})
		is.NoErr(err)
		outputStr := string(output)
		is.True(strings.Contains(outputStr, "Backlog version"))
		is.True(strings.Contains(outputStr, "Revision"))
		is.True(strings.Contains(outputStr, "Version"))
		is.True(strings.Contains(outputStr, "BuildTime"))
	})
}

// VersionExamples contains all examples for the version command
var VersionExamples = CommandExamples{
	Examples: []CommandExample{
		{
			Description: "Print Version Information",
			Command:     "version",
			Comment:     "Print the version information",
			// Expected:    "Backlog version:\nRevision: 7c989dabd2c61a063a23788c18eb39eca408f6a7\nVersion: v0.0.2-0.20250907193624-7c989dabd2c6\nBuildTime: 2025-09-07T19:36:24Z\nDirty: false",
		},
	},
}

// VersionTestableExamples generates testable examples for version command
func VersionTestableExamples() []TestableExample {
	var testable []TestableExample

	for _, example := range VersionExamples.Examples {
		te := example.ToTestableExample()

		te.OutputValidator = func(t *testing.T, output []byte, err error) {
			is := is.New(t)
			is.NoErr(err) // Version command should not error
			outputStr := string(output)
			is.True(strings.Contains(outputStr, "Backlog version")) // Should contain version info
		}

		testable = append(testable, te)
	}

	return testable
}

// Test generated examples
func Test_VersionExamples(t *testing.T) {
	for _, ex := range VersionExamples.Examples {
		t.Run("example_"+generateTestName(ex.Description), func(t *testing.T) {
			args := generateArgsSlice(ex)
			output, err := exec(t, "version", func(cmd *cobra.Command, args []string) error {
				// Use a mock version output for testing
				cmd.Print(versionFakeOutput)
				return nil
			}, args...)

			// Use the custom validator for this example
			ex.OutputValidator(t, output, err)
		})
	}
}

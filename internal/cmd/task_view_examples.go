package cmd

import (
	"testing"

	"github.com/matryer/is"
)

// ViewTestableExamples generates testable examples for view command
func ViewTestableExamples() []TestableExample {
	var testable []TestableExample
	for _, example := range ViewExamples.Examples {
		te := example.ToTestableExample()
		// Customize validation based on flags
		if _, hasJSON := example.Flags["json"]; hasJSON {
			te.OutputValidator = ValidateJSONOutput
		} else {
			// Default markdown output
			te.OutputValidator = func(t *testing.T, output []byte, err error) {
				is := is.New(t)
				is.NoErr(err) // Command should not error
				outputStr := string(output)
				// Should contain task content in markdown format
				is.True(len(outputStr) > 0) // Should produce some output
			}
		}
		testable = append(testable, te)
	}
	return testable
}

// ViewExamples contains all examples for the view command
var ViewExamples = CommandExamples{
	Examples: []CommandExample{
		{
			Description:     "View Task in Markdown",
			Command:         "backlog view",
			Args:            []string{"T01"},
			Comment:         "View task T01 in markdown format",
			OutputValidator: outputNonEmpty,
		},
		{
			Description: "View Task in JSON",
			Command:     "backlog view",
			Args:        []string{"T01"},
			Flags: map[string]string{
				"json": "",
			},
			Comment:         "View task T01 in JSON format",
			OutputValidator: outputNonEmpty,
		},
	},
}

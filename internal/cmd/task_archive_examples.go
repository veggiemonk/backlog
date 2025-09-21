package cmd

import (
	"testing"

	"github.com/matryer/is"
)

// ArchiveTestableExamples generates testable examples for archive command
func ArchiveTestableExamples() []TestableExample {
	var testable []TestableExample

	for _, example := range ArchiveExamples.Examples {
		te := example.ToTestableExample()

		te.OutputValidator = func(t *testing.T, output []byte, err error) {
			is := is.New(t)
			is.NoErr(err) // Archive command should not error
			// Archive command doesn't produce output by default, just logs success
		}

		testable = append(testable, te)
	}

	return testable
}

// ArchiveExamples contains all examples for the archive command
var ArchiveExamples = CommandExamples{
	Examples: []CommandExample{
		{
			Description:     "Archive a Task",
			Command:         "backlog archive",
			Args:            []string{"T01"},
			Comment:         "Archive task T01, moving it to the archived directory",
			OutputValidator: outputNonEmpty,
		},
		{
			Description:     "Archive by Partial ID",
			Command:         "backlog archive",
			Args:            []string{"1"},
			Comment:         "Archive task using partial ID",
			OutputValidator: outputNonEmpty,
		},
	},
}

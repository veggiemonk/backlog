package cmd

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

// SearchTestableExamples generates testable examples for search command
func SearchTestableExamples() []TestableExample {
	var testable []TestableExample
	for _, example := range SearchExamples.Examples {
		te := example.ToTestableExample()
		// Customize validation based on flags
		if _, hasJSON := example.Flags["json"]; hasJSON {
			te.OutputValidator = ValidateJSONOutput
		} else if _, hasMarkdown := example.Flags["markdown"]; hasMarkdown {
			te.OutputValidator = ValidateMarkdownOutput
		} else {
			// Default table output
			te.OutputValidator = func(t *testing.T, output []byte, err error) {
				is := is.New(t)
				is.NoErr(err) // Command should not error
				outputStr := string(output)
				// Should contain search results or "No tasks found matching"
				is.True(strings.Contains(outputStr, "Found") || strings.Contains(outputStr, "No tasks found"))
			}
		}
		testable = append(testable, te)
	}
	return testable
}

// SearchExamples contains all examples for the search command
var SearchExamples = CommandExamples{
	Examples: []CommandExample{
		{
			Description: "Basic Search",
			Command:     "backlog search",
			Args:        []string{"login"},
			Comment:     "Search for tasks containing \"login\" in any field",
		},
		{
			Description: "Search by Bug",
			Command:     "backlog search",
			Args:        []string{"bug"},
			Comment:     "Search for tasks containing \"bug\"",
		},
		{
			Description: "Search Assigned Tasks",
			Command:     "backlog search",
			Args:        []string{"@john"},
			Comment:     "Search for tasks assigned to a specific person",
		},
		{
			Description: "Search with Label",
			Command:     "backlog search",
			Args:        []string{"frontend"},
			Comment:     "Search for tasks with specific labels",
		},
		{
			Description: "Search in Acceptance Criteria",
			Command:     "backlog search",
			Args:        []string{"validation"},
			Comment:     "Search in acceptance criteria",
		},
		{
			Description: "Search with Markdown Output",
			Command:     "backlog search",
			Args:        []string{"api"},
			Flags: map[string]string{
				"markdown": "",
			},
		},
		{
			Description: "Search with JSON Output",
			Command:     "backlog search",
			Args:        []string{"api"},
			Flags: map[string]string{
				"json": "",
			},
		},
		{
			Description: "Search with Status Filter",
			Command:     "backlog search",
			Args:        []string{"user"},
			Flags: map[string]string{
				"status": "todo",
			},
		},
		{
			Description: "Search with Pagination",
			Command:     "backlog search",
			Args:        []string{"api"},
			Flags: map[string]string{
				"limit": "5",
			},
			Comment: "Show first 5 search results",
		},
		{
			Description: "Search with Offset",
			Command:     "backlog search",
			Args:        []string{"bug"},
			Flags: map[string]string{
				"limit":  "3",
				"offset": "5",
			},
			Comment: "Show 3 results starting from 6th match",
		},
	},
}

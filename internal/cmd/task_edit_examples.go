package cmd

import (
	"testing"

	"github.com/matryer/is"
)

// EditTestableExamples generates testable examples for edit command
func EditTestableExamples() []TestableExample {
	var testable []TestableExample

	for _, example := range EditExamples.Examples {
		te := example.ToTestableExample()

		// Edit commands need existing tasks, so we might need to create them first
		te.OutputValidator = func(t *testing.T, output []byte, err error) {
			is := is.New(t)
			is.NoErr(err) // Edit command should not error
			// Edit command doesn't produce output by default, just logs success
		}

		testable = append(testable, te)
	}

	return testable
}

// EditExamples contains all examples for the edit command
var EditExamples = CommandExamples{
	Examples: []CommandExample{
		{
			Description: "Change Title",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"title": "Fix the main login button styling",
			},
			Comment:         "Use the -t or --title flag to give the task a new title.",
			OutputValidator: outputNonEmpty,
		},
		{
			Description: "Update Description",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"description": "The login button on the homepage is misaligned on mobile devices. It should be centered.",
			},
			Comment:         "Use the -d or --description flag to replace the existing description with a new one.",
			OutputValidator: outputNonEmpty,
		},
		{
			Description: "Change Status",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"status": "in-progress",
			},
			Comment:         "Update the task's progress by changing its status with the -s or --status flag.",
			OutputValidator: outputNonEmpty,
		},
		{
			Description: "Re-assign to Single Person",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"assigned": "jordan",
			},
			Comment:         "You can change the assigned names for a task using the -a or --assignee flag.",
			OutputValidator: outputNonEmpty,
		},
		{
			Description:     "Re-assign to Multiple People",
			Command:         "backlog edit",
			Args:            []string{"42"},
			Flags:           map[string]string{"assigned": "jordan,casey"},
			OutputValidator: outputNonEmpty,
		},
		{
			Description:     "Update Labels",
			Command:         "backlog edit",
			Args:            []string{"42"},
			Flags:           map[string]string{"labels": "bug,frontend"},
			Comment:         "Use the -l or --labels flag to replace the existing labels.",
			OutputValidator: outputNonEmpty,
		},
		{
			Description:     "Change Priority",
			Command:         "backlog edit",
			Args:            []string{"42"},
			Flags:           map[string]string{"priority": "high"},
			Comment:         "Adjust the task's priority with the --priority flag.",
			OutputValidator: outputNonEmpty,
		},
		{
			Description:     "Add Acceptance Criteria",
			Command:         "backlog edit",
			Args:            []string{"42"},
			Flags:           map[string]string{"ac": "The button is centered on screens smaller than 576px."},
			Comment:         "Add a new AC",
			OutputValidator: outputNonEmpty,
		},
		{
			Description:     "Check Acceptance Criteria",
			Command:         "backlog edit",
			Args:            []string{"42"},
			Flags:           map[string]string{"check-ac": "1"},
			Comment:         "Check the first AC (assuming it's at index 1)",
			OutputValidator: outputNonEmpty,
		},
		{
			Description:     "Uncheck Acceptance Criteria",
			Command:         "backlog edit",
			Args:            []string{"42"},
			Flags:           map[string]string{"uncheck-ac": "1"},
			Comment:         "Uncheck the first AC",
			OutputValidator: outputNonEmpty,
		},
		{
			Description:     "Remove Acceptance Criteria",
			Command:         "backlog edit",
			Args:            []string{"42"},
			Flags:           map[string]string{"remove-ac": "2"},
			Comment:         "Remove the second AC (at index 2)",
			OutputValidator: outputNonEmpty,
		},
		{
			Description:     "Change Parent Task",
			Command:         "backlog edit",
			Args:            []string{"42"},
			Flags:           map[string]string{"parent": "18"},
			Comment:         "Move a task to be a sub-task of a different parent using the -p or --parent flag.",
			OutputValidator: outputNonEmpty,
		},
		{
			Description:     "Remove Parent",
			Command:         "backlog edit",
			Args:            []string{"42"},
			Flags:           map[string]string{"parent": ""},
			Comment:         "To remove a parent, pass an empty string",
			OutputValidator: outputNonEmpty,
		},
		{
			Description: "Add Implementation Notes",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"notes": "The issue is in the 'main.css' file, specifically in the '.login-container' class. Need to adjust the media query.",
			},
			Comment:         "Use the --notes flag to add or update technical notes for implementation.",
			OutputValidator: outputNonEmpty,
		},
		{
			Description:     "Update Implementation Plan",
			Command:         "backlog edit",
			Args:            []string{"42"},
			Flags:           map[string]string{"plan": "1. Refactor login button\n2. Test on mobile\n3. Review with team"},
			Comment:         "Use the --plan flag to add or update the implementation plan for the task.",
			OutputValidator: outputNonEmpty,
		},
		{
			Description:     "Set Single Dependency",
			Command:         "backlog edit",
			Args:            []string{"42"},
			Flags:           map[string]string{"dep": "T15"},
			Comment:         "If you want to make a task depend on another specific task",
			OutputValidator: outputNonEmpty,
		},
		{
			Description:     "Set Multiple Dependencies",
			Command:         "backlog edit",
			Args:            []string{"42"},
			Flags:           map[string]string{"dep": "T15,T18,T20"},
			Comment:         "You can make a task depend on multiple other tasks",
			OutputValidator: outputNonEmpty,
		},
		{
			Description: "Complex Example with Multiple Changes",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"status":   "in-review",
				"assigned": "alex",
				"priority": "critical",
				"notes":    "The fix is ready for review. Please check on both iOS and Android.",
				"check-ac": "1,2",
			},
			Comment:         "You can combine several flags to make multiple changes at once.",
			OutputValidator: outputNonEmpty,
		},
	},
}

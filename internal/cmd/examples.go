package cmd

import (
	"fmt"
	"strings"
	"testing"
)

// CommandExample represents a single example for a command
type CommandExample struct {
	Description string             // Human readable description
	Command     string             // The base command (e.g., "backlog create")
	Args        []string           // Positional arguments
	Flags       map[string]string  // Flag name to value mapping
	Comment     string             // Optional comment to explain the example
	Expected    string             // Optional expected output or behavior description
	Bla         func(t *testing.T) // Optional expected output or behavior description
}

// CommandExamples holds all examples for a command
type CommandExamples struct {
	Examples []CommandExample
}

// FormatExample converts a CommandExample to the CLI example string format
func (ce CommandExample) FormatExample() string {
	var parts []string

	// Add comment if provided
	if ce.Comment != "" {
		parts = append(parts, fmt.Sprintf("# %s", ce.Comment))
	}

	// Build command line
	cmdLine := ce.Command
	for _, arg := range ce.Args {
		cmdLine += fmt.Sprintf(" %q", arg)
	}

	// Add flags
	for flag, value := range ce.Flags {
		if value == "" {
			cmdLine += fmt.Sprintf(" --%s", flag)
		} else {
			cmdLine += fmt.Sprintf(" --%s %q", flag, value)
		}
	}

	parts = append(parts, cmdLine)

	// Add expected output comment if provided
	if ce.Expected != "" {
		parts = append(parts, fmt.Sprintf("# Expected: %s", ce.Expected))
	}

	return strings.Join(parts, "\n")
}

// GenerateExampleText creates the full Example field content for a cobra command
func (ces CommandExamples) GenerateExampleText() string {
	if len(ces.Examples) == 0 {
		return ""
	}
	var examples []string
	for _, example := range ces.Examples {
		examples = append(examples, example.FormatExample())
	}
	return strings.Join(examples, "\n\n")
}

// These functions should exist in the package
var expectedTestFunctions = []string{
	"CreateTestableExamples",
	"ListTestableExamples",
	"SearchTestableExamples",
	"EditTestableExamples",
	"ViewTestableExamples",
	"ArchiveTestableExamples",
	"VersionTestableExamples",
	"InstructionsTestableExamples",
	"MCPTestableExamples",
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
			Comment: "Use the -t or --title flag to give the task a new title.",
		},
		{
			Description: "Update Description",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"description": "The login button on the homepage is misaligned on mobile devices. It should be centered.",
			},
			Comment: "Use the -d or --description flag to replace the existing description with a new one.",
		},
		{
			Description: "Change Status",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"status": "in-progress",
			},
			Comment: "Update the task's progress by changing its status with the -s or --status flag.",
		},
		{
			Description: "Re-assign to Single Person",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"assigned": "jordan",
			},
			Comment: "You can change the assigned names for a task using the -a or --assignee flag.",
		},
		{
			Description: "Re-assign to Multiple People",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"assigned": "jordan,casey",
			},
		},
		{
			Description: "Update Labels",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"labels": "bug,frontend",
			},
			Comment: "Use the -l or --labels flag to replace the existing labels.",
		},
		{
			Description: "Change Priority",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"priority": "high",
			},
			Comment: "Adjust the task's priority with the --priority flag.",
		},
		{
			Description: "Add Acceptance Criteria",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"ac": "The button is centered on screens smaller than 576px.",
			},
			Comment: "Add a new AC",
		},
		{
			Description: "Check Acceptance Criteria",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"check-ac": "1",
			},
			Comment: "Check the first AC (assuming it's at index 1)",
		},
		{
			Description: "Uncheck Acceptance Criteria",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"uncheck-ac": "1",
			},
			Comment: "Uncheck the first AC",
		},
		{
			Description: "Remove Acceptance Criteria",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"remove-ac": "2",
			},
			Comment: "Remove the second AC (at index 2)",
		},
		{
			Description: "Change Parent Task",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"parent": "18",
			},
			Comment: "Move a task to be a sub-task of a different parent using the -p or --parent flag.",
		},
		{
			Description: "Remove Parent",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"parent": "",
			},
			Comment: "To remove a parent, pass an empty string",
		},
		{
			Description: "Add Implementation Notes",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"notes": "The issue is in the 'main.css' file, specifically in the '.login-container' class. Need to adjust the media query.",
			},
			Comment: "Use the --notes flag to add or update technical notes for implementation.",
		},
		{
			Description: "Update Implementation Plan",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"plan": "1. Refactor login button\n2. Test on mobile\n3. Review with team",
			},
			Comment: "Use the --plan flag to add or update the implementation plan for the task.",
		},
		{
			Description: "Set Single Dependency",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"dep": "T15",
			},
			Comment: "If you want to make a task depend on another specific task",
		},
		{
			Description: "Set Multiple Dependencies",
			Command:     "backlog edit",
			Args:        []string{"42"},
			Flags: map[string]string{
				"dep": "T15,T18,T20",
			},
			Comment: "You can make a task depend on multiple other tasks",
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
			Comment: "You can combine several flags to make multiple changes at once.",
		},
	},
}

// ViewExamples contains all examples for the view command
var ViewExamples = CommandExamples{
	Examples: []CommandExample{
		{
			Description: "View Task in Markdown",
			Command:     "backlog view",
			Args:        []string{"T01"},
			Comment:     "View task T01 in markdown format",
		},
		{
			Description: "View Task in JSON",
			Command:     "backlog view",
			Args:        []string{"T01"},
			Flags: map[string]string{
				"json": "",
			},
			Comment: "View task T01 in JSON format",
		},
	},
}

// ArchiveExamples contains all examples for the archive command
var ArchiveExamples = CommandExamples{
	Examples: []CommandExample{
		{
			Description: "Archive a Task",
			Command:     "backlog archive",
			Args:        []string{"T01"},
			Comment:     "Archive task T01, moving it to the archived directory",
		},
		{
			Description: "Archive by Partial ID",
			Command:     "backlog archive",
			Args:        []string{"1"},
			Comment:     "Archive task using partial ID",
		},
	},
}

// VersionExamples contains all examples for the version command
var VersionExamples = CommandExamples{
	Examples: []CommandExample{
		{
			Description: "Print Version Information",
			Command:     "backlog version",
			Comment:     "Print the version information",
			Expected:    "Backlog version:\nRevision: 7c989dabd2c61a063a23788c18eb39eca408f6a7\nVersion: v0.0.2-0.20250907193624-7c989dabd2c6\nBuildTime: 2025-09-07T19:36:24Z\nDirty: false",
		},
	},
}

// InstructionsExamples contains all examples for the instructions command
var InstructionsExamples = CommandExamples{
	Examples: []CommandExample{
		{
			Description: "Output Instructions",
			Command:     "backlog instructions",
			Comment:     "outputs the instructions",
		},
		{
			Description: "Save Instructions to File",
			Command:     "backlog instructions >> AGENTS.md",
			Comment:     "add instructions to agent base prompt",
		},
	},
}

// MCPExamples contains all examples for the MCP command
var MCPExamples = CommandExamples{
	Examples: []CommandExample{
		{
			Description: "Start MCP Server with HTTP",
			Command:     "backlog mcp",
			Flags: map[string]string{
				"http": "",
			},
			Comment: "Start the MCP server using HTTP transport on default port 8106",
		},
		{
			Description: "Start MCP Server with Custom Port",
			Command:     "backlog mcp",
			Flags: map[string]string{
				"http": "",
				"port": "4321",
			},
			Comment: "Start the MCP server using HTTP transport on port 4321",
		},
		{
			Description: "Start MCP Server with Stdio",
			Command:     "backlog mcp",
			Comment:     "Start the MCP server using stdio transport",
		},
	},
}


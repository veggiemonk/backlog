package cmd

import (
	"fmt"
	"strings"
)

// CommandExample represents a single example for a command
type CommandExample struct {
	Description string            // Human readable description
	Command     string            // The base command (e.g., "backlog create")
	Args        []string          // Positional arguments
	Flags       map[string]string // Flag name to value mapping
	Comment     string            // Optional comment to explain the example
	Expected    string            // Optional expected output or behavior description
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

// CreateExamples contains all examples for the create command
var CreateExamples = CommandExamples{
	Examples: []CommandExample{
		{
			Description: "Basic Task Creation",
			Command:     "backlog create",
			Args:        []string{"Fix the login button styling"},
			Comment:     "This is the simplest way to create a task, providing only a title.",
		},
		{
			Description: "Task with Description",
			Command:     "backlog create",
			Args:        []string{"Implement password reset"},
			Flags: map[string]string{
				"description": "Users should be able to request a password reset link via their email. This involves creating a new API endpoint and a front-end form.",
			},
			Comment: "Use the -d or --description flag to add more detailed information about the task.",
		},
		{
			Description: "Assigning to Single Person",
			Command:     "backlog create",
			Args:        []string{"Design the new dashboard"},
			Flags: map[string]string{
				"assigned": "alex",
			},
			Comment: "You can assign a task to one or more team members using the -a or --assigned flag.",
		},
		{
			Description: "Assigning to Multiple People",
			Command:     "backlog create",
			Args:        []string{"Code review for the payment gateway"},
			Flags: map[string]string{
				"assigned": "jordan,casey",
			},
		},
		{
			Description: "Adding Labels",
			Command:     "backlog create",
			Args:        []string{"Update third-party dependencies"},
			Flags: map[string]string{
				"labels": "bug,backend,security",
			},
			Comment: "Use the -l or --labels flag to categorize the task with comma-separated labels.",
		},
		{
			Description: "Setting High Priority",
			Command:     "backlog create",
			Args:        []string{"Hotfix: Production database is down"},
			Flags: map[string]string{
				"priority": "high",
			},
			Comment: "Specify the task's priority with the --priority flag. The default is \"medium\".",
		},
		{
			Description: "Setting Low Priority",
			Command:     "backlog create",
			Args:        []string{"Refactor the old user model"},
			Flags: map[string]string{
				"priority": "low",
			},
		},
		{
			Description: "Defining Acceptance Criteria",
			Command:     "backlog create",
			Args:        []string{"Develop user profile page"},
			Flags: map[string]string{
				"ac": "Users can view their own profile information.,Users can upload a new profile picture.,The page is responsive on mobile devices.",
			},
			Comment: "Use the --ac flag multiple times to list the conditions that must be met for the task to be considered complete.",
		},
		{
			Description: "Creating a Sub-task",
			Command:     "backlog create",
			Args:        []string{"Add Google OAuth login"},
			Flags: map[string]string{
				"parent": "15",
			},
			Comment: "Link a new task to a parent task using the -p or --parent flag. This is useful for breaking down larger tasks.",
		},
		{
			Description: "Setting Single Dependency",
			Command:     "backlog create",
			Args:        []string{"Deploy user authentication"},
			Flags: map[string]string{
				"deps": "T15",
			},
			Comment: "Use the --deps flag to specify that this task depends on other tasks being completed first.",
		},
		{
			Description: "Setting Multiple Dependencies",
			Command:     "backlog create",
			Args:        []string{"Integration testing"},
			Flags: map[string]string{
				"deps": "T15,T18,T20",
			},
			Comment: "This means the task cannot be started until tasks T15, T18, and T20 are completed.",
		},
		{
			Description: "Task with Implementation Notes",
			Command:     "backlog create",
			Args:        []string{"Optimize database queries"},
			Flags: map[string]string{
				"notes": "Focus on the user lookup queries in the authentication module. Consider adding indexes on email and username fields.",
			},
			Comment: "Use the --notes flag to add implementation notes to help with development.",
		},
		{
			Description: "Task with Implementation Plan",
			Command:     "backlog create",
			Args:        []string{"Implement user registration flow"},
			Flags: map[string]string{
				"plan": "1. Design registration form UI\n2. Create user validation logic\n3. Set up email verification\n4. Add password strength requirements\n5. Write integration tests",
			},
			Comment: "Use the --plan flag to add a structured implementation plan.",
		},
		{
			Description: "Complex Example with Multiple Flags",
			Command:     "backlog create",
			Args:        []string{"Build the new reporting feature"},
			Flags: map[string]string{
				"description": "Create a new section in the app that allows users to generate and export monthly performance reports in PDF format.",
				"assigned":    "drew",
				"labels":      "feature,frontend,backend",
				"priority":    "high",
				"ac":          "Report generation logic is accurate.,Users can select a date range for the report.,The exported PDF has the correct branding and layout.",
				"parent":      "23",
			},
			Comment: "Here is a comprehensive example that uses several flags at once to create a very detailed task.",
		},
	},
}

// ListExamples contains all examples for the list command
var ListExamples = CommandExamples{
	Examples: []CommandExample{
		{
			Description: "List All Tasks",
			Command:     "backlog list",
			Comment:     "List all tasks with all columns",
		},
		{
			Description: "Filter by Status",
			Command:     "backlog list",
			Flags: map[string]string{
				"status": "todo",
			},
			Comment: "List tasks with status \"todo\"",
		},
		{
			Description: "Filter by Multiple Statuses",
			Command:     "backlog list",
			Flags: map[string]string{
				"status": "todo,in-progress",
			},
			Comment: "List tasks with status \"todo\" or \"in-progress\"",
		},
		{
			Description: "Filter by Parent",
			Command:     "backlog list",
			Flags: map[string]string{
				"parent": "12345",
			},
			Comment: "List tasks that are sub-tasks of the task with ID \"12345\"",
		},
		{
			Description: "Filter by Assigned User",
			Command:     "backlog list",
			Flags: map[string]string{
				"assigned": "alice",
			},
			Comment: "List tasks assigned to alice",
		},
		{
			Description: "Filter Unassigned Tasks",
			Command:     "backlog list",
			Flags: map[string]string{
				"unassigned": "",
			},
			Comment: "List tasks that have no one assigned",
		},
		{
			Description: "Filter by Labels",
			Command:     "backlog list",
			Flags: map[string]string{
				"labels": "bug,feature",
			},
			Comment: "List tasks containing either \"bug\" or \"feature\" labels",
		},
		{
			Description: "Filter by Priority",
			Command:     "backlog list",
			Flags: map[string]string{
				"priority": "high",
			},
			Comment: "List all high priority tasks",
		},
		{
			Description: "Filter Tasks with Dependencies",
			Command:     "backlog list",
			Flags: map[string]string{
				"has-dependency": "",
			},
			Comment: "List tasks that have at least one dependency",
		},
		{
			Description: "Filter Blocking Tasks",
			Command:     "backlog list",
			Flags: map[string]string{
				"depended-on": "",
				"status":      "todo",
			},
			Comment: "List all the blocking tasks.",
		},
		{
			Description: "Hide Extra Fields",
			Command:     "backlog list",
			Flags: map[string]string{
				"hide-extra": "",
			},
			Comment: "Hide extra fields (labels, priority, assigned)",
		},
		{
			Description: "Sort by Priority",
			Command:     "backlog list",
			Flags: map[string]string{
				"sort": "priority",
			},
			Comment: "Sort tasks by priority",
		},
		{
			Description: "Multiple Sort Fields",
			Command:     "backlog list",
			Flags: map[string]string{
				"sort": "updated,priority",
			},
			Comment: "Sort tasks by updated date, then priority",
		},
		{
			Description: "Reverse Order",
			Command:     "backlog list",
			Flags: map[string]string{
				"reverse": "",
			},
			Comment: "Reverse the order of tasks",
		},
		{
			Description: "Markdown Output",
			Command:     "backlog list",
			Flags: map[string]string{
				"markdown": "",
			},
			Comment: "List tasks in markdown format",
		},
		{
			Description: "JSON Output",
			Command:     "backlog list",
			Flags: map[string]string{
				"json": "",
			},
			Comment: "List tasks in JSON format",
		},
		{
			Description: "Pagination - Limit",
			Command:     "backlog list",
			Flags: map[string]string{
				"limit": "10",
			},
			Comment: "List first 10 tasks",
		},
		{
			Description: "Pagination - Limit and Offset",
			Command:     "backlog list",
			Flags: map[string]string{
				"limit":  "5",
				"offset": "10",
			},
			Comment: "List 5 tasks starting from 11th task",
		},
	},
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
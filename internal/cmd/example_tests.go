package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/strcase"
)

// CommandExample represents a single example for a command
type CommandExample struct {
	Description     string                                       // Human readable description
	Command         string                                       // The base command (e.g., "backlog create")
	Args            []string                                     // Positional arguments
	Flags           map[string]string                            // Flag name to value mapping
	Comment         string                                       // Optional comment to explain the example
	ExpectedError   bool                                         // Whether this example should produce an error
	OutputValidator func(t *testing.T, output []byte, err error) // Custom validation function
}

// CommandExamples holds all examples for a command
type CommandExamples struct {
	Name     string
	Examples []CommandExample
}

// GenerateExampleText creates the full Example field content for a cobra command
func generateExampleText(ces CommandExamples) string {
	if len(ces.Examples) == 0 {
		return ""
	}
	var examples []string
	for _, example := range ces.Examples {
		examples = append(examples, formatExample(example))
	}
	return strings.Join(examples, "\n\n")
}

// TestableExample represents an example that can be tested
type TestableExample struct {
	CommandExample
	TestName        string                                       // Generated test name
	ExpectedError   bool                                         // Whether this example should produce an error
	OutputValidator func(t *testing.T, output []byte, err error) // Custom validation function
}

// GenerateTestName creates a descriptive test name from the example
func generateTestName(description string) string {
	return strcase.Snake(description)
	// // Clean up the description to make it a valid test name
	// name := strings.ReplaceAll(ce.Description, " ", "_")
	// name = strings.ReplaceAll(name, "-", "_")
	// name = strings.ReplaceAll(name, "/", "_")
	// name = strings.ReplaceAll(name, ":", "")
	// name = strings.ReplaceAll(name, ",", "_")
	// name = strings.ReplaceAll(name, "(", "")
	// name = strings.ReplaceAll(name, ")", "")
	// name = strings.ToLower(name)
	// return name
}

// ToTestableExample converts a CommandExample to a TestableExample with default behavior
func (ce CommandExample) ToTestableExample() TestableExample {
	return TestableExample{
		CommandExample: ce,
		TestName:       generateTestName(ce.Description),
		ExpectedError:  false,
		OutputValidator: func(t *testing.T, output []byte, err error) {
			is := is.New(t)
			is.NoErr(err) // Command should not error
		},
	}
}

// GenerateArgsSlice creates the argument slice for exec() function
func generateArgsSlice(ce CommandExample) []string {
	var args []string
	// Add positional arguments
	args = append(args, ce.Args...)
	// Add flags
	for flag, value := range ce.Flags {
		if value == "" {
			// Boolean flag
			args = append(args, fmt.Sprintf("--%s", flag))
		} else {
			// Flag with value
			args = append(args, fmt.Sprintf("--%s", flag), value)
		}
	}
	return args
}

//	func (te TestableExample) GenerateArgsSlice() []string {
//		var args []string
//		// Add positional arguments
//		args = append(args, te.Args...)
//		// Add flags
//		for flag, value := range te.Flags {
//			if value == "" {
//				// Boolean flag
//				args = append(args, fmt.Sprintf("--%s", flag))
//			} else {
//				// Flag with value
//				args = append(args, fmt.Sprintf("--%s", flag), value)
//			}
//		}
//		return args
// }

func testForExamples(ex CommandExamples) []TestableExample {
	var testable []TestableExample
	for _, example := range ex.Examples {
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
				// Should contain table headers or "No tasks found"
				is.True(strings.Contains(outputStr, "ID") || strings.Contains(outputStr, "No tasks found"))
			}
		}
		testable = append(testable, te)
	}
	return testable
}

// formatExample converts a CommandExample to the CLI example string format
func formatExample(ce CommandExample) string {
	var parts []string

	if ce.Description != "" {
		parts = append(parts, fmt.Sprintf("# %s\n", ce.Description))
	}
	// Add comment if provided
	if ce.Comment != "" {
		parts = append(parts, fmt.Sprintf("# %s", ce.Comment))
	}

	// Build command line
	cmdLine := Name + " " + ce.Command
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
	// if ce.Expected != "" {
	// 	parts = append(parts, fmt.Sprintf("# Expected: %s", ce.Expected))
	// }
	//
	return strings.Join(parts, "\n")
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

// ValidateJSONOutput is a common validator for JSON output
func ValidateJSONOutput(t *testing.T, output []byte, err error) {
	is := is.New(t)
	is.NoErr(err) // Command should not error

	// Try to parse as JSON array (for list/search commands)
	var tasks []*core.Task
	if jsonErr := json.Unmarshal(output, &tasks); jsonErr != nil {
		// Try to parse as single task (for view commands)
		var task core.Task
		jsonErr2 := json.Unmarshal(output, &task)
		is.NoErr(jsonErr2) // Should be valid JSON for either array or single task
	}
}

// ValidateNonEmptyOutput ensures command produces some output
func ValidateNonEmptyOutput(t *testing.T, output []byte, err error) {
	is := is.New(t)
	is.NoErr(err)                                                // Command should not error
	is.True(len(output) > 0)                                     // Should produce some output
	is.True(!strings.Contains(string(output), "No tasks found")) // Should not be empty result
}

// ValidateMarkdownOutput checks for markdown table format
func ValidateMarkdownOutput(t *testing.T, output []byte, err error) {
	is := is.New(t)
	is.NoErr(err) // Command should not error
	outputStr := string(output)
	is.True(strings.Contains(outputStr, "|"))    // Markdown tables contain pipes
	is.True(strings.Contains(outputStr, ":---")) // Markdown table separator
}

// MCPTestableExamples generates testable examples for MCP command
func MCPTestableExamples() []TestableExample {
	var testable []TestableExample

	for _, example := range MCPExamples.Examples {
		te := example.ToTestableExample()

		// MCP command starts a server, so we expect it to potentially run indefinitely
		// We'll need special handling for this in tests
		te.ExpectedError = true // Server commands typically don't return normally in tests
		te.OutputValidator = func(t *testing.T, output []byte, err error) {
			// For MCP server, we might just check that it starts without immediate error
			// This will need special test handling
		}

		testable = append(testable, te)
	}

	return testable
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

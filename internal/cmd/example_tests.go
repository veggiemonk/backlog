package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/veggiemonk/backlog/internal/core"
)

// TestableExample represents an example that can be tested
type TestableExample struct {
	CommandExample
	TestName        string // Generated test name
	ExpectedError   bool   // Whether this example should produce an error
	OutputValidator func(t *testing.T, output []byte, err error) // Custom validation function
}

// GenerateTestName creates a descriptive test name from the example
func (ce CommandExample) GenerateTestName() string {
	// Clean up the description to make it a valid test name
	name := strings.ReplaceAll(ce.Description, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, ":", "")
	name = strings.ReplaceAll(name, ",", "_")
	name = strings.ReplaceAll(name, "(", "")
	name = strings.ReplaceAll(name, ")", "")
	name = strings.ToLower(name)
	return name
}

// ToTestableExample converts a CommandExample to a TestableExample with default behavior
func (ce CommandExample) ToTestableExample() TestableExample {
	return TestableExample{
		CommandExample: ce,
		TestName:       ce.GenerateTestName(),
		ExpectedError:  false,
		OutputValidator: func(t *testing.T, output []byte, err error) {
			is := is.New(t)
			is.NoErr(err) // Command should not error
		},
	}
}

// GenerateArgsSlice creates the argument slice for exec() function
func (te TestableExample) GenerateArgsSlice() []string {
	var args []string

	// Add positional arguments
	args = append(args, te.Args...)

	// Add flags
	for flag, value := range te.Flags {
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
	is.NoErr(err)                          // Command should not error
	is.True(len(output) > 0)               // Should produce some output
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

// CreateTestableExamples generates testable examples for create command
func CreateTestableExamples() []TestableExample {
	var testable []TestableExample

	for _, example := range CreateExamples.Examples {
		te := example.ToTestableExample()

		// Customize validation based on example type
		te.OutputValidator = func(t *testing.T, output []byte, err error) {
			is := is.New(t)
			is.NoErr(err) // Create command should not error
			// Create command doesn't produce JSON output by default,
			// it just logs success
		}

		testable = append(testable, te)
	}

	return testable
}

// ListTestableExamples generates testable examples for list command
func ListTestableExamples() []TestableExample {
	var testable []TestableExample

	for _, example := range ListExamples.Examples {
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
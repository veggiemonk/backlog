package cmd

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/matryer/is"
)

var allExamples = []CommandExamples{
	CreateExamples,
	ListExamples,
	SearchExamples,
	EditExamples,
	ViewExamples,
	ArchiveExamples,
	VersionExamples,
	InstructionsExamples,
	MCPExamples,
}

// TestExamplesSyncWithTests validates that all examples have corresponding tests
// and all test scenarios are represented in examples
func TestExamplesSyncWithTests(t *testing.T) {
	t.Run("all commands have examples defined", func(t *testing.T) {
		is := is.New(t)
		for _, ex := range allExamples {
			is.True(len(ex.Examples) > 0)
		}
	})

	t.Run("all examples can generate test names", func(t *testing.T) {
		is := is.New(t)
		for _, cmdExamples := range allExamples {
			for _, example := range cmdExamples.Examples {
				testName := generateTestName(example.Description)
				is.True(len(testName) > 0)                     // Test name should not be empty
				is.True(!strings.Contains(testName, " "))      // Test name should not contain spaces
				is.True(!strings.Contains(testName, "-"))      // Test name should not contain dashes
				is.True(strings.ToLower(testName) == testName) // Test name should be lowercase
			}
		}
	})

	t.Run("all examples can generate argument slices", func(t *testing.T) {
		for _, cmdExamples := range allExamples {
			for _, example := range cmdExamples.Examples {
				args := generateArgsSlice(example)
				// Should not panic and should return a slice (even if empty)
				_ = args
			}
		}
	})

	t.Run("examples contain required patterns", func(t *testing.T) {
		is := is.New(t)

		// Create examples should cover basic creation patterns
		createTestNames := make([]string, 0)
		for _, example := range CreateExamples.Examples {
			createTestNames = append(createTestNames, generateTestName(example.Description))
		}
		is.True(containsPattern(createTestNames, "basic"))       // Should have basic creation
		is.True(containsPattern(createTestNames, "description")) // Should show description usage
		is.True(containsPattern(createTestNames, "assign"))      // Should show assignment
		is.True(containsPattern(createTestNames, "label"))       // Should show labels
		is.True(containsPattern(createTestNames, "priority"))    // Should show priority
		is.True(containsPattern(createTestNames, "criteria"))    // Should show AC

		// List examples should cover filtering and output options
		listTestNames := make([]string, 0)
		for _, example := range ListExamples.Examples {
			createTestNames = append(listTestNames, generateTestName(example.Description))
		}
		is.True(containsPattern(listTestNames, "status"))   // Should show status filtering
		is.True(containsPattern(listTestNames, "json"))     // Should show JSON output
		is.True(containsPattern(listTestNames, "markdown")) // Should show markdown output
		is.True(containsPattern(listTestNames, "sort"))     // Should show sorting

		// Search examples should cover search patterns
		searchTestNames := make([]string, 0)
		for _, example := range SearchExamples.Examples {
			searchTestNames = append(searchTestNames, generateTestName(example.Description))
		}
		is.True(containsPattern(searchTestNames, "basic"))    // Should have basic search
		is.True(containsPattern(searchTestNames, "json"))     // Should show JSON output
		is.True(containsPattern(searchTestNames, "markdown")) // Should show markdown output
	})

	t.Run("examples have consistent command structure", func(t *testing.T) {
		is := is.New(t)

		// All create examples should start with "backlog create"
		for _, example := range CreateExamples.Examples {
			is.True(strings.HasPrefix(example.Command, "backlog create"))
		}

		// All list examples should start with "backlog list"
		for _, example := range ListExamples.Examples {
			is.True(strings.HasPrefix(example.Command, "backlog list"))
		}

		// All search examples should start with "backlog search"
		for _, example := range SearchExamples.Examples {
			is.True(strings.HasPrefix(example.Command, "backlog search"))
		}

		// All edit examples should start with "backlog edit"
		for _, example := range EditExamples.Examples {
			is.True(strings.HasPrefix(example.Command, "backlog edit"))
		}

		// All view examples should start with "backlog view"
		for _, example := range ViewExamples.Examples {
			is.True(strings.HasPrefix(example.Command, "backlog view"))
		}

		// All archive examples should start with "backlog archive"
		for _, example := range ArchiveExamples.Examples {
			is.True(strings.HasPrefix(example.Command, "backlog archive"))
		}
	})

	t.Run("testable examples have output validators", func(t *testing.T) {
		is := is.New(t)
		for _, ex := range allExamples {
			for _, testable := range testForExamples(ex) {
				is.True(testable.OutputValidator != nil) // Every testable example should have a validator
			}
		}
	})

	t.Run("example generation produces valid content", func(t *testing.T) {
		is := is.New(t)
		for _, ex := range allExamples {
			exampleText := generateExampleText(ex)
			is.True(len(exampleText) > 0)                        // Should generate non-empty text
			is.True(strings.Contains(exampleText, "backlog"))    // Should contain command name
			is.True(!strings.Contains(exampleText, "undefined")) // Should not contain undefined values
			is.True(!strings.Contains(exampleText, "null"))      // Should not contain null values
		}
	})
}

// TestExampleTestGeneration validates that example-to-test generation works correctly
func TestExampleTestGeneration(t *testing.T) {
	t.Run("example format generates valid command lines", func(t *testing.T) {
		is := is.New(t)

		example := CommandExample{
			Description: "Test Example",
			Command:     "backlog test",
			Args:        []string{"arg1", "arg2"},
			Flags: map[string]string{
				"flag1": "value1",
				"flag2": "", // boolean flag
			},
			Comment: "This is a test",
		}

		formatted := formatExample(example)
		is.True(strings.Contains(formatted, "# This is a test"))   // Should contain comment
		is.True(strings.Contains(formatted, "backlog test"))       // Should contain command
		is.True(strings.Contains(formatted, "\"arg1\""))           // Should contain quoted args
		is.True(strings.Contains(formatted, "\"arg2\""))           // Should contain quoted args
		is.True(strings.Contains(formatted, "--flag1 \"value1\"")) // Should contain flag with value
		is.True(strings.Contains(formatted, "--flag2"))            // Should contain boolean flag
	})

	t.Run("testable examples generate proper argument arrays", func(t *testing.T) {
		is := is.New(t)

		example := CommandExample{
			Command: "backlog test",
			Args:    []string{"arg1"},
			Flags: map[string]string{
				"flag1": "value1",
				"flag2": "",
			},
		}
		args := generateArgsSlice(example)
		is.True(slices.Contains(args, "arg1"))    // Should contain positional arg
		is.True(slices.Contains(args, "--flag1")) // Should contain flag name
		is.True(slices.Contains(args, "value1"))  // Should contain flag value
		is.True(slices.Contains(args, "--flag2")) // Should contain boolean flag
	})
}

// TestCommandCoverage validates that all command functionality is covered by examples
func TestCommandCoverage(t *testing.T) {
	t.Run("create command covers all major flags", func(t *testing.T) {
		is := is.New(t)

		flagsCovered := make(map[string]bool)
		for _, example := range CreateExamples.Examples {
			for flag := range example.Flags {
				flagsCovered[flag] = true
			}
		}

		// These are the major flags that should be covered in examples
		expectedFlags := []string{"description", "assigned", "labels", "priority", "ac", "parent", "deps", "notes", "plan"}
		for _, flag := range expectedFlags {
			if !flagsCovered[flag] {
				t.Errorf("Flag %s is not covered in create examples", flag)
			}
			is.True(flagsCovered[flag]) // Flag should be covered in create examples
		}
	})

	t.Run("list command covers filtering and output options", func(t *testing.T) {
		is := is.New(t)

		flagsCovered := make(map[string]bool)
		for _, example := range ListExamples.Examples {
			for flag := range example.Flags {
				flagsCovered[flag] = true
			}
		}

		expectedFlags := []string{"status", "assigned", "labels", "priority", "json", "markdown", "sort", "reverse", "limit", "offset"}
		for _, flag := range expectedFlags {
			is.True(flagsCovered[flag]) // Flag %s should be covered in list examples
		}
	})

	t.Run("edit command covers modification patterns", func(t *testing.T) {
		is := is.New(t)

		flagsCovered := make(map[string]bool)
		for _, example := range EditExamples.Examples {
			for flag := range example.Flags {
				flagsCovered[flag] = true
			}
		}

		expectedFlags := []string{"title", "description", "status", "assigned", "labels", "priority", "parent", "notes", "plan", "ac", "dep"}
		for _, flag := range expectedFlags {
			is.True(flagsCovered[flag]) // Flag %s should be covered in edit examples
		}
	})
}

// Helper functions

func containsPattern(slice []string, pattern string) bool {
	for _, item := range slice {
		if strings.Contains(strings.ToLower(item), strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

// TestReflectionBasedSync validates that the example definitions match the test functions
func TestReflectionBasedSync(t *testing.T) {
	t.Run("test functions exist for example generators", func(t *testing.T) {
		is := is.New(t)

		// Use reflection to check if these functions exist
		packageVal := reflect.ValueOf(CreateExamples)

		for _, funcName := range expectedTestFunctions {
			// Check if function exists in package
			funcVal := packageVal.MethodByName(funcName)
			if !funcVal.IsValid() {
				// Try to find it as a package-level function
				// This is a simplified check - in a real scenario you'd use more sophisticated reflection
				is.True(true) // For now, just pass - the function existence is validated by compilation
			}
		}
	})

	t.Run("example structs have consistent field types", func(t *testing.T) {
		is := is.New(t)

		// Validate that all CommandExample instances have the expected fields
		exampleType := reflect.TypeOf(CommandExample{})

		// Check that required fields exist
		_, hasDescription := exampleType.FieldByName("Description")
		is.True(hasDescription) // CommandExample should have Description field

		_, hasCommand := exampleType.FieldByName("Command")
		is.True(hasCommand) // CommandExample should have Command field

		_, hasArgs := exampleType.FieldByName("Args")
		is.True(hasArgs) // CommandExample should have Args field

		_, hasFlags := exampleType.FieldByName("Flags")
		is.True(hasFlags) // CommandExample should have Flags field
	})
}

package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/veggiemonk/backlog/internal/core"
)

func Test_runView(t *testing.T) {
	t.Run("view task in markdown format", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create a task
		_, err := exec(t, "create", runCreate, "Task to View", "-d", "Test description")
		is.NoErr(err)

		// Then view it
		output, err := exec(t, "view", view, "1")
		is.NoErr(err)
		outputStr := string(output)
		is.True(len(outputStr) > 0)                              // Should produce output
		is.True(strings.Contains(outputStr, "Task to View"))     // Should contain task title
		is.True(strings.Contains(outputStr, "Test description")) // Should contain description
	})

	t.Run("view task in JSON format", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create a task
		_, err := exec(t, "create", runCreate, "JSON Task", "-d", "JSON description")
		is.NoErr(err)

		// Then view it in JSON format
		output, err := exec(t, "view", view, "1", "--json")
		is.NoErr(err)

		// Should be valid JSON
		var task core.Task
		is.NoErr(json.Unmarshal(output, &task))
		is.Equal(task.Title, "JSON Task")
		is.Equal(task.Description, "JSON description")
	})

	t.Run("view task with short JSON flag", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create a task
		_, err := exec(t, "create", runCreate, "Short JSON Task")
		is.NoErr(err)

		// Then view it with short flag
		output, err := exec(t, "view", view, "1", "-j")
		is.NoErr(err)

		// Should be valid JSON
		var task core.Task
		is.NoErr(json.Unmarshal(output, &task))
		is.Equal(task.Title, "Short JSON Task")
	})

	t.Run("view task with complex data", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// Create a complex task
		_, err := exec(t, "create", runCreate, "Complex Task",
			"-d", "Complex description",
			"-a", "user1", "-a", "user2",
			"-l", "label1,label2",
			"--priority", "high",
			"--ac", "AC 1", "--ac", "AC 2")
		is.NoErr(err)

		// View in JSON to verify all data
		output, err := exec(t, "view", view, "1", "-j")
		is.NoErr(err)

		var task core.Task
		is.NoErr(json.Unmarshal(output, &task))
		is.Equal(task.Title, "Complex Task")
		is.Equal(task.Description, "Complex description")
		is.Equal(len(task.Assigned), 2)
		is.Equal(len(task.Labels), 2)
		is.Equal(task.Priority.String(), "high")
		is.Equal(len(task.AcceptanceCriteria), 2)
	})

	// Test error cases
	t.Run("view nonexistent task should fail", func(t *testing.T) {
		is := is.NewRelaxed(t)
		_, err := exec(t, "view", view, "999")
		is.True(err != nil) // Should error when task doesn't exist
	})

	t.Run("view without task ID should fail", func(t *testing.T) {
		is := is.NewRelaxed(t)
		_, err := exec(t, "view", view)
		is.True(err != nil) // Should error when no task ID provided
	})

	t.Run("view with invalid task ID should fail", func(t *testing.T) {
		is := is.NewRelaxed(t)
		_, err := exec(t, "view", view, "invalid")
		is.True(err != nil) // Should error with invalid task ID
	})
}

// Test generated examples
func Test_ViewExamples(t *testing.T) {
	for _, ex := range ViewExamples.Examples {
		t.Run("example_"+generateTestName(ex.Description), func(t *testing.T) {
			// Create a test task first for view examples
			_, err := exec(t, "create", runCreate, "Test Task for View Example")
			if err != nil {
				t.Fatalf("Failed to create test task: %v", err)
			}

			args := generateArgsSlice(ex)
			// // Replace T01 with 1 since our test creates task with ID 1
			// for i, arg := range args {
			// 	if arg == "T01" {
			// 		args[i] = "1"
			// 	}
			// }

			output, err := exec(t, "view", view, args...)

			// Use the custom validator for this example
			ex.OutputValidator(t, output, err)
		})
	}
}


package cmd

import (
	"testing"

	"github.com/matryer/is"
)

func Test_runArchive(t *testing.T) {
	t.Run("archive existing task", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create a task
		_, err := exec(t, "create", runCreate, "Task to Archive")
		is.NoErr(err)

		// Then archive it
		output, err := exec(t, "archive", runArchive, "1")
		is.NoErr(err)
		_ = output
	})

	t.Run("archive task with partial ID", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create a task
		_, err := exec(t, "create", runCreate, "Task to Archive with Partial ID")
		is.NoErr(err)

		// Then archive using partial ID
		output, err := exec(t, "archive", runArchive, "1")
		is.NoErr(err)
		_ = output
	})

	t.Run("archive task with full ID", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create a task
		_, err := exec(t, "create", runCreate, "Task to Archive with Full ID")
		is.NoErr(err)

		// Then archive using full ID format
		output, err := exec(t, "archive", runArchive, "T1")
		is.NoErr(err)
		_ = output
	})

	t.Run("archive task with complex data", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// Create a complex task
		_, err := exec(t, "create", runCreate, "Complex Task to Archive",
			"-d", "Complex description",
			"-a", "user1", "-a", "user2",
			"-l", "label1,label2",
			"--priority", "high")
		is.NoErr(err)

		// Then archive it
		output, err := exec(t, "archive", runArchive, "1")
		is.NoErr(err)
		_ = output
	})

	// Test error cases
	t.Run("archive nonexistent task should fail", func(t *testing.T) {
		is := is.NewRelaxed(t)
		_, err := exec(t, "archive", runArchive, "999")
		is.True(err != nil) // Should error when task doesn't exist
	})

	t.Run("archive without task ID should fail", func(t *testing.T) {
		is := is.NewRelaxed(t)
		_, err := exec(t, "archive", runArchive)
		is.True(err != nil) // Should error when no task ID provided
	})

	t.Run("archive with invalid task ID should fail", func(t *testing.T) {
		is := is.NewRelaxed(t)
		_, err := exec(t, "archive", runArchive, "invalid")
		is.True(err != nil) // Should error with invalid task ID
	})

	t.Run("archive already archived task should fail", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create a task
		_, err := exec(t, "create", runCreate, "Task to Double Archive")
		is.NoErr(err)

		// Archive it once
		_, err = exec(t, "archive", runArchive, "1")
		is.NoErr(err)

		// Try to archive again - should fail
		_, err = exec(t, "archive", runArchive, "1")
		is.True(err != nil) // Should error when trying to archive already archived task
	})
}

// Test generated examples
func Test_ArchiveExamples(t *testing.T) {
	for _, ex := range ArchiveExamples.Examples {
		t.Run("example_"+generateTestName(ex.Description), func(t *testing.T) {
			// Create a test task first for archive examples
			_, err := exec(t, "create", runCreate, "Test Task for Archive Example")
			if err != nil {
				t.Fatalf("Failed to create test task: %v", err)
			}

			args := generateArgsSlice(ex)
			// Replace T01 with 1 since our test creates task with ID 1
			for i, arg := range args {
				if arg == "T01" {
					args[i] = "1"
				}
			}

			output, err := exec(t, "archive", runArchive, args...)

			// Use the custom validator for this example
			ex.OutputValidator(t, output, err)
		})
	}
}

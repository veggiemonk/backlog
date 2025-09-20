package cmd

import (
	"testing"

	"github.com/matryer/is"
)

func Test_runCreate(t *testing.T) {
	t.Run("basic task creation", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "create", runCreate, "Test Task")
		is.NoErr(err)
		// Create command doesn't output JSON by default, just logs success
		_ = output
	})

	t.Run("task with description", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "create", runCreate, "Task with description", "-d", "This is a test description")
		is.NoErr(err)
		_ = output
	})

	t.Run("task with single assigned user", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "create", runCreate, "Assigned Task", "-a", "testuser")
		is.NoErr(err)
		_ = output
	})

	t.Run("task with multiple assigned users", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "create", runCreate, "Multi-assigned Task", "-a", "user1", "-a", "user2")
		is.NoErr(err)
		_ = output
	})

	t.Run("task with labels", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "create", runCreate, "Labeled Task", "-l", "bug,frontend,urgent")
		is.NoErr(err)
		_ = output
	})

	t.Run("task with high priority", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "create", runCreate, "High Priority Task", "--priority", "high")
		is.NoErr(err)
		_ = output
	})

	t.Run("task with low priority", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "create", runCreate, "Low Priority Task", "--priority", "low")
		is.NoErr(err)
		_ = output
	})

	t.Run("task with acceptance criteria", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "create", runCreate, "Task with AC", "--ac", "First AC", "--ac", "Second AC")
		is.NoErr(err)
		_ = output
	})

	t.Run("task with parent", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create a parent task
		_, err := exec(t, "create", runCreate, "Parent Task")
		is.NoErr(err)

		// Then create a child task
		output, err := exec(t, "create", runCreate, "Child Task", "-p", "1")
		is.NoErr(err)
		_ = output
	})

	t.Run("task with dependencies", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "create", runCreate, "Task with deps", "--deps", "T1", "--deps", "T2")
		is.NoErr(err)
		_ = output
	})

	t.Run("task with implementation notes", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "create", runCreate, "Task with notes", "--notes", "Implementation notes here")
		is.NoErr(err)
		_ = output
	})

	t.Run("task with implementation plan", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "create", runCreate, "Task with plan", "--plan", "1. Step one\n2. Step two")
		is.NoErr(err)
		_ = output
	})

	t.Run("complex task with multiple flags", func(t *testing.T) {
		is := is.NewRelaxed(t)
		output, err := exec(t, "create", runCreate, "Complex Task",
			"-d", "Complex task description",
			"-a", "user1", "-a", "user2",
			"-l", "feature,backend",
			"--priority", "high",
			"--ac", "AC 1", "--ac", "AC 2",
			"--notes", "Complex notes",
			"--plan", "Complex plan")
		is.NoErr(err)
		_ = output
	})

	// Test error cases
	t.Run("create without title should fail", func(t *testing.T) {
		is := is.NewRelaxed(t)
		_, err := exec(t, "create", runCreate)
		is.True(err != nil) // Should error when no title provided
	})

	t.Run("create with invalid priority should fail", func(t *testing.T) {
		is := is.NewRelaxed(t)
		_, err := exec(t, "create", runCreate, "Task", "--priority", "invalid")
		is.True(err != nil) // Should error with invalid priority
	})
}

// Test generated examples
func Test_CreateExamples(t *testing.T) {
	testableExamples := CreateTestableExamples()

	for _, example := range testableExamples {
		t.Run("example_"+example.TestName, func(t *testing.T) {
			args := example.GenerateArgsSlice()
			output, err := exec(t, "create", runCreate, args...)

			// Use the custom validator for this example
			example.OutputValidator(t, output, err)
		})
	}
}
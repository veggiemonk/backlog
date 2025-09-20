package cmd

import (
	"testing"

	"github.com/matryer/is"
)

func Test_runEdit(t *testing.T) {
	t.Run("edit task title", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create a task
		_, err := exec(t, "create", runCreate, "Original Title")
		is.NoErr(err)

		// Then edit its title
		output, err := exec(t, "edit", runEdit, "1", "-t", "New Title")
		is.NoErr(err)
		_ = output
	})

	t.Run("edit task description", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create a task
		_, err := exec(t, "create", runCreate, "Task to edit")
		is.NoErr(err)

		// Then edit its description
		output, err := exec(t, "edit", runEdit, "1", "-d", "New description")
		is.NoErr(err)
		_ = output
	})

	t.Run("edit task status", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create a task
		_, err := exec(t, "create", runCreate, "Task to edit status")
		is.NoErr(err)

		// Then change its status
		output, err := exec(t, "edit", runEdit, "1", "-s", "in-progress")
		is.NoErr(err)
		_ = output
	})

	t.Run("edit task priority", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create a task
		_, err := exec(t, "create", runCreate, "Task to edit priority")
		is.NoErr(err)

		// Then change its priority
		output, err := exec(t, "edit", runEdit, "1", "--priority", "high")
		is.NoErr(err)
		_ = output
	})

	t.Run("edit task assigned users", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create a task
		_, err := exec(t, "create", runCreate, "Task to edit assigned")
		is.NoErr(err)

		// Then change assigned users
		output, err := exec(t, "edit", runEdit, "1", "-a", "newuser")
		is.NoErr(err)
		_ = output
	})

	t.Run("edit task labels", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create a task
		_, err := exec(t, "create", runCreate, "Task to edit labels")
		is.NoErr(err)

		// Then change labels
		output, err := exec(t, "edit", runEdit, "1", "-l", "newlabel,anotherlabel")
		is.NoErr(err)
		_ = output
	})

	t.Run("edit task parent", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create a parent task
		_, err := exec(t, "create", runCreate, "Parent Task")
		is.NoErr(err)

		// Create a child task
		_, err = exec(t, "create", runCreate, "Child Task")
		is.NoErr(err)

		// Then edit child to set parent
		output, err := exec(t, "edit", runEdit, "2", "-p", "1")
		is.NoErr(err)
		_ = output
	})

	t.Run("edit task dependencies", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create some tasks
		_, err := exec(t, "create", runCreate, "Task 1")
		is.NoErr(err)
		_, err = exec(t, "create", runCreate, "Task 2")
		is.NoErr(err)

		// Then edit task 2 to depend on task 1
		output, err := exec(t, "edit", runEdit, "2", "--dep", "T1")
		is.NoErr(err)
		_ = output
	})

	t.Run("edit task notes", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create a task
		_, err := exec(t, "create", runCreate, "Task to edit notes")
		is.NoErr(err)

		// Then edit its notes
		output, err := exec(t, "edit", runEdit, "1", "--notes", "Updated implementation notes")
		is.NoErr(err)
		_ = output
	})

	t.Run("edit task plan", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create a task
		_, err := exec(t, "create", runCreate, "Task to edit plan")
		is.NoErr(err)

		// Then edit its plan
		output, err := exec(t, "edit", runEdit, "1", "--plan", "Updated implementation plan")
		is.NoErr(err)
		_ = output
	})

	t.Run("add acceptance criteria", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create a task
		_, err := exec(t, "create", runCreate, "Task to add AC")
		is.NoErr(err)

		// Then add AC
		output, err := exec(t, "edit", runEdit, "1", "--ac", "New acceptance criterion")
		is.NoErr(err)
		_ = output
	})

	t.Run("multiple edits at once", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create a task
		_, err := exec(t, "create", runCreate, "Task for multiple edits")
		is.NoErr(err)

		// Then edit multiple fields
		output, err := exec(t, "edit", runEdit, "1",
			"-t", "Updated Title",
			"-d", "Updated Description",
			"-s", "in-progress",
			"--priority", "high",
			"-a", "updateduser")
		is.NoErr(err)
		_ = output
	})

	// Test error cases
	t.Run("edit nonexistent task should fail", func(t *testing.T) {
		is := is.NewRelaxed(t)
		_, err := exec(t, "edit", runEdit, "999", "-t", "New Title")
		is.True(err != nil) // Should error when task doesn't exist
	})

	t.Run("edit without task ID should fail", func(t *testing.T) {
		is := is.NewRelaxed(t)
		_, err := exec(t, "edit", runEdit, "-t", "New Title")
		is.True(err != nil) // Should error when no task ID provided
	})

	t.Run("edit with invalid status should fail", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// First create a task
		_, err := exec(t, "create", runCreate, "Task to test invalid status")
		is.NoErr(err)

		// Then try to set invalid status
		_, err = exec(t, "edit", runEdit, "1", "-s", "invalid-status")
		is.True(err != nil) // Should error with invalid status
	})
}

// Test generated examples
func Test_EditExamples(t *testing.T) {
	testableExamples := EditTestableExamples()

	for _, example := range testableExamples {
		t.Run("example_"+example.TestName, func(t *testing.T) {
			// Create a test task first for edit examples
			_, err := exec(t, "create", runCreate, "Test Task for Edit Example")
			if err != nil {
				t.Fatalf("Failed to create test task: %v", err)
			}

			args := example.GenerateArgsSlice()
			output, err := exec(t, "edit", runEdit, args...)

			// Use the custom validator for this example
			example.OutputValidator(t, output, err)
		})
	}
}
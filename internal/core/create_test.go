package core

import (
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/afero"
)

func TestCreateTask(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	store := NewFileTaskStore(fs, ".backlog")

	params := CreateTaskParams{
		Title:       "Test Task",
		Description: "This is a test description.",
		AC:          []string{"Criterion 1"},
	}

	createdTask, err := store.Create(params)
	is.NoErr(err)
	is.Equal("Test Task", createdTask.Title)
	is.Equal("T01", createdTask.ID.Name())

	// Check if the file was actually created in the memory fs
	filePath := ".backlog/T01-test_task.md"
	exists, err := afero.Exists(fs, filePath)
	is.NoErr(err)
	is.True(exists)

	// Check the content of the created file
	contentBytes, err := afero.ReadFile(fs, filePath)
	is.NoErr(err)
	content := string(contentBytes)

	is.True(strings.Contains(content, "title: Test Task"))
	is.True(strings.Contains(content, "## Description\n\nThis is a test description."))
	is.True(strings.Contains(content, "- [ ] #1 Criterion 1"))

	// Create a subtask
	subtaskParams := CreateTaskParams{
		Title:       "Subtask 1",
		Description: "This is a subtask.",
		Parent:      "T1",
		AC:          []string{"Subtask Criterion"},
	}

	subtask, err := store.Create(subtaskParams)
	is.NoErr(err)
	is.Equal("Subtask 1", subtask.Title)
	is.Equal("T01.01", subtask.ID.Name())

	subtaskFilePath := ".backlog/T01.01-subtask_1.md"
	exists, err = afero.Exists(fs, subtaskFilePath)
	is.NoErr(err)
	is.True(exists)

	subtaskContentBytes, err := afero.ReadFile(fs, subtaskFilePath)
	is.NoErr(err)
	subtaskContent := string(subtaskContentBytes)

	is.True(strings.Contains(subtaskContent, "title: Subtask 1"))
	is.True(strings.Contains(subtaskContent, "## Description\n\nThis is a subtask."))
	is.True(strings.Contains(subtaskContent, "- [ ] #1 Subtask Criterion"))

	// Test: create a task with an invalid parent
}

func TestCreateTask_InvalidParent(t *testing.T) {
	fs := afero.NewMemMapFs()
	store := NewFileTaskStore(fs, ".backlog")

	t.Run("creation with invalid parent", func(t *testing.T) {
		is := is.New(t)
		invalidParent := "T999"
		params := CreateTaskParams{
			Title:       "Should Fail",
			Description: "This should not be created.",
			Parent:      invalidParent,
			AC:          []string{"Invalid Parent Criterion"},
		}
		_, err := store.Create(params)
		is.True(err != nil) // Should error due to invalid parent
	})
}

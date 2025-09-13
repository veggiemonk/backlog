package core

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/afero"
)

func TestArchiveTask(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	store := NewFileTaskStore(fs, ".backlog")

	// First, create a task to archive
	params := CreateTaskParams{
		Title:       "Test Task",
		Description: "This is a test description.",
	}

	createdTask, err := store.Create(params)
	is.NoErr(err)

	// Now, archive the task
	archivedTaskPath, err := store.Archive(createdTask.ID)
	is.NoErr(err)

	b, err := afero.ReadFile(fs, archivedTaskPath)
	is.NoErr(err)
	archivedTask, err := parseTask(b)
	is.NoErr(err)

	is.Equal(archivedTask.Status, StatusArchived)

	// Check that the file has been moved
	archivedDir := filepath.Join(".backlog", "archived")
	archivedFilePath := filepath.Join(archivedDir, createdTask.FileName())
	exists, err := afero.Exists(fs, archivedFilePath)
	is.NoErr(err)
	is.True(exists)

	// Check that the original file is gone
	originalFilePath := store.Path(createdTask)
	exists, err = afero.Exists(fs, originalFilePath)
	is.NoErr(err)
	is.True(!exists)

	// Check that the history has been updated
	hasArchivedEntry := false
	for _, entry := range archivedTask.History {
		if strings.Contains(entry.Change, "archived") {
			hasArchivedEntry = true
			break
		}
	}
	is.True(hasArchivedEntry)
}

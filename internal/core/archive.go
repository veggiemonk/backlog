package core

import (
	"fmt"
	"path/filepath"
)

// Archive moves a task to the archived directory and updates its status.
func (f *FileTaskStore) Archive(id TaskID) (string, error) {
	task, err := f.Get(id.String())
	if err != nil {
		return "", fmt.Errorf("get task %q: %w", id, err)
	}
	if err = f.Update(&task, EditTaskParams{
		NewStatus: ptr(string(StatusArchived)),
	}); err != nil {
		return "", fmt.Errorf("set status archived task %q: %w", id, err)
	}

	// Move the file to the archived directory.
	archivedDir := filepath.Join(f.tasksDir, "archived")
	if err := f.fs.MkdirAll(archivedDir, 0o750); err != nil {
		return "", fmt.Errorf("create archived directory: %w", err)
	}
	oldPath := f.Path(task)
	newPath := filepath.Join(archivedDir, filepath.Base(oldPath))
	if err := f.fs.Rename(oldPath, newPath); err != nil {
		return "", fmt.Errorf("move task file: %w", err)
	}
	return newPath, nil
}

func ptr[T any](v T) *T { return &v }

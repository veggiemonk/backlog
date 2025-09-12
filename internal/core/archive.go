package core

import (
	"fmt"
	"path/filepath"
	"time"
)

// Archive moves a task to the archived directory and updates its status.
func (f *FileTaskStore) Archive(id TaskID) (*Task, error) {
	task, err := f.Get(id.String())
	if err != nil {
		return nil, fmt.Errorf("could not get task %q: %w", id, err)
	}

	// Move the file to the archived directory.
	archivedDir := filepath.Join(f.tasksDir, "archived")
	if err := f.fs.MkdirAll(archivedDir, 0750); err != nil {
		return nil, fmt.Errorf("could not create archived directory: %w", err)
	}

	oldPath := f.Path(task)
	newPath := filepath.Join(archivedDir, filepath.Base(oldPath))

	if err := f.fs.Rename(oldPath, newPath); err != nil {
		return nil, fmt.Errorf("could not move task file: %w", err)
	}

	// Update the task status and history.
	task.Status = StatusArchived
	task.UpdatedAt = time.Now()
	RecordChange(task, "archived")

	return task, nil
}

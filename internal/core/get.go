package core

import (
	"fmt"

	"github.com/spf13/afero"
)

// Get implements TaskStore.
func (f *FileTaskStore) Get(id string) (task Task, err error) {
	taskID, err := parseTaskID(id)
	if err != nil {
		return task, fmt.Errorf("invalid task ID '%s': %w", id, err)
	}

	filePath, err := f.findTaskFileByID(taskID)
	if err != nil {
		return task, fmt.Errorf("find task file: %w", err)
	}
	b, err := afero.ReadFile(f.fs, filePath)
	if err != nil {
		return task, err
	}

	task, err = parseTask(b)
	if err != nil {
		return task, fmt.Errorf("parse task %s: %v", filePath, err)
	}

	return task, nil
}

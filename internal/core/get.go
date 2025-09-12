package core

import (
	"fmt"

	"github.com/spf13/afero"
)

// Get implements TaskStore.
func (f *FileTaskStore) Get(id string) (*Task, error) {
	taskID, err := parseTaskID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID '%s': %w", id, err)
	}

	filePath, err := f.findTaskFileByID(taskID)
	if err != nil {
		return nil, err
	}
	b, err := afero.ReadFile(f.fs, filePath)
	if err != nil {
		return nil, err
	}

	task, err := parseTask(b, filePath)
	if err != nil {
		return nil, err
	}

	return task, nil
}

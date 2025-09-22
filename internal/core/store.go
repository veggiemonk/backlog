package core

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

const (
	descHeader     = "## Description"
	acHeader       = "## Acceptance Criteria"
	planHeader     = "## Implementation Plan"
	notesHeader    = "## Implementation Notes"
	acStartComment = "<!-- AC:BEGIN -->"
	acEndComment   = "<!-- AC:END -->"
)

type FileTaskStore struct {
	fs       afero.Fs
	tasksDir string
	locker   Locker
}

func NewFileTaskStore(fs afero.Fs, tasksDir string, locker Locker) *FileTaskStore {
	return &FileTaskStore{
		fs:       fs,
		tasksDir: tasksDir,
		locker:   locker,
	}
}

func (f *FileTaskStore) Path(t *Task) string {
	return filepath.Join(f.tasksDir, t.FileName())
}

func (f *FileTaskStore) Fs() afero.Fs {
	return f.fs
}

func (f *FileTaskStore) write(task *Task) error {
	// Create the tasks directory if it doesn't exist
	if err := f.fs.MkdirAll(f.tasksDir, 0755); err != nil {
		return err
	}
	fullContent := task.Bytes()
	filePath := f.Path(task)
	if err := afero.WriteFile(f.fs, filePath, fullContent, 0644); err != nil {
		return err
	}
	return nil
}

// getNextTaskID finds the next available task ID in the tasks directory.
func (f *FileTaskStore) getNextTaskID(treePath ...int) (TaskID, error) {
	if err := f.locker.Lock(); err != nil {
		return TaskID{}, fmt.Errorf("could not acquire lock: %w", err)
	}
	defer func() {
		if err := f.locker.Unlock(); err != nil {
			// Log the error, but don't return it as it would mask the primary error
			fmt.Printf("could not release lock: %v\n", err)
		}
	}()
	files, err := afero.ReadDir(f.fs, f.tasksDir)
	if err != nil {
		return TaskID{}, err
	}

	// Collect all TaskIDs that match the given treePath
	var matchingIDs []TaskID
	for _, file := range files {
		if file.IsDir() || !strings.HasPrefix(file.Name(), TaskIDPrefix) || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}
		id, err := parseTaskIDfromFileName(file.Name())
		if err != nil {
			continue // Skip files with invalid IDs
		}
		// Check if the ID matches the desired tree path
		if len(id.seg) < len(treePath) {
			continue
		}
		matches := true
		for i := range treePath {
			if id.seg[i] != treePath[i] {
				matches = false
				break
			}
		}
		if matches && len(id.seg) == len(treePath)+1 {
			matchingIDs = append(matchingIDs, id)
		}
	}

	var nextSeg []int
	if len(treePath) == 0 {
		// Top-level task
		max := 0
		for _, id := range matchingIDs {
			if id.seg[0] > max {
				max = id.seg[0]
			}
		}
		nextSeg = []int{max + 1}
	} else {
		// Subtask or deeper
		max := 0
		for _, id := range matchingIDs {
			last := id.seg[len(id.seg)-1]
			if last > max {
				max = last
			}
		}
		nextSeg = append([]int{}, treePath...)
		nextSeg = append(nextSeg, max+1)
	}

	return TaskID{seg: nextSeg}, nil
}

// FindTaskFileByID searches the tasks directory for a task file matching the given ID.
// The ID can be in the format "T123" or just "123".
func (f *FileTaskStore) findTaskFileByID(id TaskID) (string, error) {
	var foundPath string
	files, err := afero.ReadDir(f.fs, f.tasksDir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasPrefix(file.Name(), id.Name()) && strings.HasSuffix(file.Name(), ".md") {
			foundPath = filepath.Join(f.tasksDir, file.Name())
			break
		}
	}

	if foundPath != "" {
		return foundPath, nil
	}

	return "", fmt.Errorf("task with ID '%s' not found", id)
}
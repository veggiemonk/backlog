package core

import (
	"fmt"

	"github.com/alexflint/go-filemutex"
)

// FileLocker is a wrapper around filemutex.FileMutex that implements the Locker interface.
type FileLocker struct {
	mutex *filemutex.FileMutex
}

// NewFileLocker creates a new FileLocker.
func NewFileLocker(path string) (*FileLocker, error) {
	mutex, err := filemutex.New(path)
	if err != nil {
		return nil, fmt.Errorf("could not create file mutex: %w", err)
	}
	return &FileLocker{mutex: mutex}, nil
}

// Lock acquires the lock.
func (l *FileLocker) Lock() error {
	return l.mutex.Lock()
}

// Unlock releases the lock.
func (l *FileLocker) Unlock() error {
	return l.mutex.Unlock()
}
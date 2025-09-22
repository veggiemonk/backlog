package core

import "sync"

// MockLocker implements the Locker interface using a sync.Mutex for testing.
type MockLocker struct {
	mu sync.Mutex
}

// NewMockLocker creates a new MockLocker.
func NewMockLocker() *MockLocker {
	return &MockLocker{}
}

// Lock acquires the mutex lock.
func (l *MockLocker) Lock() error {
	l.mu.Lock()
	return nil
}

// Unlock releases the mutex lock.
func (l *MockLocker) Unlock() error {
	l.mu.Unlock()
	return nil
}

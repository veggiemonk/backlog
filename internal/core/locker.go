package core

// Locker defines the interface for a lock that can be acquired and released.
type Locker interface {
	Lock() error
	Unlock() error
}

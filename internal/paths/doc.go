// Package paths contains path resolution utilities used by the CLI.
//
// The main entry point is ResolveTasksDir, which determines the directory to
// store and read tasks from in a way that works locally and inside containers.
// It supports absolute paths, relative paths resolved from the current working
// directory, upward ancestor search, optional anchoring to the git repository
// root, and a sensible fallback to the current working directory when no git
// repository is present.
package paths

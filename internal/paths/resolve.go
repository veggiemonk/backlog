package paths

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/veggiemonk/backlog/internal/commit"
)

const (
	DefaultDir       = ".backlog"
	DefaultDirEnvVar = "BACKLOG_FOLDER"
)

// ResolveTasksDir determines the directory to use for tasks based on the
// provided input path and the current execution context (local or container).
//
// Strategy:
//  1. Absolute path => return as-is.
//  2. Relative path:
//     a) If exists from CWD => return it.
//     b) Walk up parents; if <ancestor>/<dir> exists => return it.
//     c) If in a git repo => use <gitRoot>/<dir>.
//     d) Fallback to <CWD>/<dir> (will be created on demand by the store).
func ResolveTasksDir(fs afero.Fs, dir string) (string, error) {
	// If absolute, just return it (creation happens later in store).
	if filepath.IsAbs(dir) {
		return dir, nil
	}

	// 2a) Check relative to current working directory first.
	if exists, err := afero.DirExists(fs, dir); err == nil && exists {
		return dir, nil
	} else if err != nil {
		return dir, fmt.Errorf("check existence of directory %q: %v", dir, err)
	}

	// 2b) Search upwards for an existing directory named `dir`.
	cwd, err := os.Getwd()
	if err != nil {
		return dir, fmt.Errorf("get working directory: %w", err)
	}
	absCWD, err := filepath.Abs(cwd)
	if err != nil {
		return dir, fmt.Errorf("invalid working dir: %w", err)
	}
	probe := absCWD
	for {
		candidate := filepath.Join(probe, dir)
		exists, err := afero.DirExists(fs, candidate)
		if err != nil {
			return dir, fmt.Errorf("check existence of directory %q: %v", candidate, err)
		}
		if exists {
			return candidate, nil
		}
		parent := filepath.Dir(probe)
		if parent == probe { // reached filesystem root
			break
		}
		probe = parent
	}

	// 2c) Attempt to anchor to the git repository root if available.
	if rootDir, err := commit.FindTopLevelGitDir(); err == nil && rootDir != "" {
		return filepath.Join(rootDir, dir), nil
	}

	// 2d) Fallback to <CWD>/<dir> (will be created by the store when needed).
	return filepath.Join(absCWD, dir), nil
}

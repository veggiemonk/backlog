package commit

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/imjasonh/version"
	"github.com/veggiemonk/backlog/internal/logging"
)

func FindTopLevelGitDir() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir, err := filepath.Abs(workingDir)
	if err != nil {
		return "", fmt.Errorf("invalid working dir: %w", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("no git repository found")
		}
		dir = parent
	}
}

// Add stages and commits files.
func Add(path, oldPath, message string) error {
	repoRoot, err := FindTopLevelGitDir()
	if err != nil {
		return err
	}
	repo, err := git.PlainOpenWithOptions(repoRoot, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return fmt.Errorf("not a git repository: %w", err)
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("could not get worktree: %w", err)
	}
	logging.Info("auto-committing changes", "path", path, "oldPath", oldPath, "message", message)
	if path == "" {
		logging.Info("no changes to commit")
		return nil
	}
	if strings.HasPrefix(path, repoRoot) {
		path = path[len(repoRoot)+1:]
	}
	if _, err = worktree.Add(path); err != nil {
		return fmt.Errorf("error staging file %s: %w", path, err)
	}
	// used in case of a rename (change of title)
	if oldPath != "" {
		if oldPath != "" && strings.HasPrefix(oldPath, repoRoot) {
			oldPath = oldPath[len(repoRoot)+1:]
		}
		_, err := worktree.Add(oldPath)
		if err != nil {
			return fmt.Errorf("error staging file %s: %w", path, err)
		}
	}
	opts := &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Backlog CLI",
			Email: fmt.Sprintf("backlog-cli+%s@localhost", version.Get().Version),
			When:  time.Now(),
		},
	}
	if _, err = worktree.Commit(message, opts); err != nil {
		return fmt.Errorf("error creating commit: %w", err)
	}
	logging.Info("changes committed successfully", "path", path, "oldPath", oldPath)
	return nil
}

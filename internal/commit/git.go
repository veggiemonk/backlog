package commit

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/imjasonh/version"
	"github.com/veggiemonk/backlog/internal/logging"
)

type GitHandle struct {
	repo     *git.Repository
	worktree *git.Worktree
}

// NewHandle opens the git repository in the current directory.
func NewHandle() (*GitHandle, error) {
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, fmt.Errorf("not a git repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("could not get worktree: %w", err)
	}

	return &GitHandle{
		repo:     repo,
		worktree: worktree,
	}, nil
}

// Stage adds the given paths to the git staging area.
func (g *GitHandle) Stage(paths []string, oldPaths []string) error {
	if len(oldPaths) > 0 {
		for _, path := range oldPaths {
			_, err := g.worktree.Remove(path)
			if err != nil {
				return fmt.Errorf("error removing file %s: %w", path, err)
			}
		}
	}

	for _, path := range paths {
		_, err := g.worktree.Add(path)
		if err != nil {
			return fmt.Errorf("error staging file %s: %w", path, err)
		}
	}
	return nil
}

// Commit creates a new commit with the given message.
func (g *GitHandle) Commit(message string) error {
	_, err := g.worktree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Backlog CLI",
			Email: fmt.Sprintf("backlog-cli+%s@localhost", version.Get().Version),
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("error creating commit: %w", err)
	}
	return nil
}

// AutoCommit stages and commits files if the auto-commit feature is enabled.
func (g *GitHandle) AutoCommit(paths []string, message string) error {
	logging.Info("auto-committing changes", "paths", paths, "message", message)
	if err := g.Stage(paths, paths); err != nil {
		return fmt.Errorf("auto-commit failed during staging: %w", err)
	}

	if err := g.Commit(message); err != nil {
		return fmt.Errorf("auto-commit failed during commit: %w", err)
	}

	logging.Info("changes committed successfully", "paths", paths)
	return nil
}

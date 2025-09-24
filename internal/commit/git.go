// Package commit only purpose is to commit tasks.
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
	"github.com/spf13/afero"
	"github.com/veggiemonk/backlog/internal/core"
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

// Add stages and commits files. Will not commit if repo is dirty.
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
	status, err := worktree.Status()
	if err != nil {
		return fmt.Errorf("could not get status: %w", err)
	}
	for fpath, fstat := range status {
		logging.Info("git status", "path", fpath, "staging", string(fstat.Staging))
		// If there is another file added, just skip
		if fpath != path && fstat.Staging == git.Added {
			logging.Warn("the repository status is not clean, skip commit", "status", status.String())
			return nil
		}
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

// CheckForTaskConflicts detects and optionally resolves ID conflicts in the backlog
func CheckForTaskConflicts(tasksDir string, autoResolve bool) error {
	fs := afero.NewOsFs()
	detector := core.NewConflictDetector(fs, tasksDir)

	// Detect conflicts
	conflicts, err := detector.DetectConflicts()
	if err != nil {
		return fmt.Errorf("failed to detect conflicts: %w", err)
	}

	if len(conflicts) == 0 {
		logging.Info("no task ID conflicts detected")
		return nil
	}

	// Log detected conflicts
	summary := core.SummarizeConflicts(conflicts)
	logging.Warn("task ID conflicts detected",
		"total", summary.TotalConflicts,
		"duplicate_ids", summary.DuplicateIDs,
		"orphaned_children", summary.OrphanedChildren,
		"invalid_hierarchy", summary.InvalidHierarchy,
	)

	// If auto-resolve is enabled, attempt to resolve conflicts
	if autoResolve {
		return resolveConflictsAutomatically(detector, conflicts, tasksDir)
	}

	// Otherwise, just log the conflicts for manual resolution
	for _, conflict := range conflicts {
		logging.Error("conflict requires manual resolution",
			"type", conflict.Type.String(),
			"id", conflict.ConflictID.String(),
			"files", conflict.Files,
			"description", conflict.Description,
		)
	}

	return fmt.Errorf("found %d task ID conflicts that require manual resolution", len(conflicts))
}

// resolveConflictsAutomatically attempts to resolve conflicts using the chronological strategy
func resolveConflictsAutomatically(detector *core.ConflictDetector, conflicts []core.IDConflict, tasksDir string) error {
	fs := afero.NewOsFs()
	store := core.NewFileTaskStore(fs, tasksDir)
	resolver := core.NewConflictResolver(detector, store)

	// Create resolution plan using chronological strategy (keeps older tasks)
	plan, err := resolver.CreateResolutionPlan(conflicts, core.ResolutionStrategyChronological)
	if err != nil {
		return fmt.Errorf("failed to create resolution plan: %w", err)
	}

	logging.Info("executing automatic conflict resolution", "actions", len(plan.Actions))

	// Execute the plan with reference updates (not dry run)
	results, err := resolver.ExecuteResolutionPlanWithReferences(plan, false)
	if err != nil {
		return fmt.Errorf("failed to execute resolution plan: %w", err)
	}

	// Log the results
	for _, result := range results {
		logging.Info("conflict resolution action", "result", result)
	}

	logging.Info("automatic conflict resolution completed", "actions_executed", len(results))
	return nil
}

// PostMergeConflictCheck should be called after Git merge operations to detect ID conflicts
func PostMergeConflictCheck(tasksDir string) error {
	return CheckForTaskConflicts(tasksDir, true) // Auto-resolve after merges
}

// PreCommitConflictCheck should be called before commits to ensure no conflicts exist
func PreCommitConflictCheck(tasksDir string) error {
	return CheckForTaskConflicts(tasksDir, false) // Don't auto-resolve before commits, just detect
}

// AddWithConflictDetection adds and commits files with automatic conflict detection
func AddWithConflictDetection(path, oldPath, message, tasksDir string) error {
	// First, check for conflicts before committing
	if err := PreCommitConflictCheck(tasksDir); err != nil {
		logging.Warn("conflicts detected before commit", "error", err)
		// Continue with commit even if conflicts exist (they'll be flagged for manual resolution)
	}

	// Perform the normal commit
	if err := Add(path, oldPath, message); err != nil {
		return err
	}

	logging.Info("commit completed with conflict detection")
	return nil
}

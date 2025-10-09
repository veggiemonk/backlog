package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
	"github.com/veggiemonk/backlog/internal/paths"
)

const doctorDescription = `Diagnose and fix task ID conflicts that can occur when creating tasks
in separate Git branches. Conflicts arise when multiple branches generate the same task IDs.

This command provides conflict detection and resolution capabilities to maintain
task ID uniqueness and data integrity.

Conflict types detected:
- Duplicate IDs (same ID in multiple files)
- Orphaned children (tasks with non-existent parents)
- Invalid hierarchy (parent-child ID mismatch)

Examples:
` +
	"```" +
	`

  backlog doctor                          # Detect conflicts in text format
  backlog doctor --json                   # Detect conflicts in JSON format
  backlog doctor -j                       # Detect conflicts in JSON format (short flag)

  backlog doctor --fix                    # Auto-fix using chronological strategy
  backlog doctor --fix --dry-run          # Preview changes without applying
  backlog doctor --fix --strategy=auto    # Use auto-renumbering strategy
  backlog doctor --fix --strategy=manual  # Create manual resolution plan

` + "```"

func newDoctorCommand(rt *runtime) *cli.Command {
	return &cli.Command{
		Name:  "doctor",
		Usage: "Diagnose and fix task ID conflicts",
		Description: doctorDescription,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "json", Aliases: []string{"j"}, Usage: "Output in JSON format"},
			&cli.BoolFlag{Name: "fix", Usage: "Automatically fix detected conflicts"},
			&cli.StringFlag{Name: "strategy", Value: "chronological", Usage: "Resolution strategy when using --fix (chronological|auto|manual)"},
			&cli.BoolFlag{Name: "dry-run", Usage: "Show what would be changed without making changes (use with --fix)"},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fs := afero.NewOsFs()
			tasksDir := rt.tasksDir
			if tasksDir == "" {
				tasksDir = cmd.String(configFolder)
			}

			if cmd.Bool("fix") {
				return resolveConflicts(fs, tasksDir, cmd.String("strategy"), cmd.Bool("dry-run"))
			}
			return detectConflicts(cmd.Root().Writer, fs, tasksDir, cmd.Bool("json"))
		},
	}
}

func detectConflicts(w io.Writer, fs afero.Fs, tasksDir string, jsonOutput bool) error {
	var err error
	tasksDir, err = paths.ResolveTasksDir(fs, tasksDir)
	if err != nil {
		return fmt.Errorf("failed to resolve tasks directory: %w", err)
	}

	detector := core.NewConflictDetector(fs, tasksDir)

	conflicts, err := detector.DetectConflicts()
	if err != nil {
		return fmt.Errorf("failed to detect conflicts: %w", err)
	}

	if jsonOutput {
		output := map[string]any{
			"conflicts": conflicts,
			"summary":   core.SummarizeConflicts(conflicts),
		}
		if err := json.NewEncoder(w).Encode(output); err != nil {
			return fmt.Errorf("failed to encode JSON: %w", err)
		}
		return nil
	}

	summary := core.SummarizeConflicts(conflicts)
	if summary.TotalConflicts == 0 {
		logging.Info("no task ID conflicts detected")
		return nil
	}

	logging.Warn("conflicts", "found", summary.TotalConflicts)

	if summary.DuplicateIDs > 0 {
		conflicts := map[string][]string{}
		for _, conflict := range summary.ConflictsByType[core.ConflictTypeDuplicateID] {
			conflicts[conflict.ConflictID.String()] = conflict.Files
		}
		logging.Warn("duplicate IDs", slog.Int("found", summary.DuplicateIDs), slog.Any("details", conflicts))
	}

	if summary.OrphanedChildren > 0 {
		conflicts := []string{}
		for _, conflict := range summary.ConflictsByType[core.ConflictTypeOrphanedChild] {
			conflicts = append(conflicts, fmt.Sprintf("- %s references non-existent parent\n", conflict.ConflictID.String()))
		}
		logging.Warn("orphaned children", slog.Int("found", summary.OrphanedChildren), slog.Any("references", conflicts))
	}

	if summary.InvalidHierarchy > 0 {
		conflicts := []string{}
		for _, conflict := range summary.ConflictsByType[core.ConflictTypeInvalidHierarchy] {
			conflicts = append(conflicts, fmt.Sprintf("- %s has incorrect parent structure\n", conflict.ConflictID.String()))
		}
		logging.Warn("invalid hierarchy", slog.Int("found", summary.InvalidHierarchy), slog.Any("references", conflicts))
	}

	logging.Info("Run 'backlog doctor --fix' to fix these conflicts.")
	return nil
}

func resolveConflicts(fs afero.Fs, tasksDir string, strategyName string, dryRun bool) error {
	var err error
	tasksDir, err = paths.ResolveTasksDir(fs, tasksDir)
	if err != nil {
		return fmt.Errorf("failed to resolve tasks directory: %w", err)
	}

	detector := core.NewConflictDetector(fs, tasksDir)

	conflicts, err := detector.DetectConflicts()
	if err != nil {
		return fmt.Errorf("failed to detect conflicts: %w", err)
	}

	if len(conflicts) == 0 {
		logging.Info("no conflicts to resolve")
		return nil
	}

	store := core.NewFileTaskStore(fs, tasksDir)
	resolver := core.NewConflictResolver(detector, store)

	var strategy core.ResolutionStrategy
	switch strategyName {
	case "chronological":
		strategy = core.ResolutionStrategyChronological
	case "auto":
		strategy = core.ResolutionStrategyAutoRenumber
	case "manual":
		strategy = core.ResolutionStrategyManual
	default:
		return fmt.Errorf("invalid strategy: %s", strategyName)
	}

	plan, err := resolver.CreateResolutionPlan(conflicts, strategy)
	if err != nil {
		return fmt.Errorf("failed to create resolution plan: %w", err)
	}

	logging.Info("resolution plan", "strategy", strategyName, "plan", plan.Summary)
	if len(plan.Actions) == 0 {
		logging.Info("no actions required")
		return nil
	}

	for i, action := range plan.Actions {
		logging.Info(fmt.Sprintf("%d. %s\n", i+1, action.Description))
		if action.Type == "manual" {
			logging.Info("type: manual intervention required")
		} else {
			logging.Info(fmt.Sprintf("   Type: %s\n", action.Type))
			if !action.OriginalID.IsZero() {
				logging.Info(fmt.Sprintf("   Original ID: %s\n", action.OriginalID.String()))
			}
			if !action.NewID.IsZero() {
				logging.Info(fmt.Sprintf("   New ID: %s\n", action.NewID.String()))
			}
			if action.FilePath != "" {
				logging.Info(fmt.Sprintf("   File: %s\n", action.FilePath))
			}
		}
	}

	if dryRun {
		logging.Info("DRY RUN - No changes were made")
		return nil
	}

	if strategy == core.ResolutionStrategyManual {
		logging.Warn("manual resolution required - no automatic actions taken")
		logging.Warn("please review the conflicts above and resolve them manually")
		return nil
	}

	logging.Info("executing resolution plan")
	results, err := resolver.ExecuteResolutionPlanWithReferences(plan, false)
	if err != nil {
		return fmt.Errorf("failed to execute resolution plan: %w", err)
	}

	for _, result := range results {
		logging.Info(result)
	}

	logging.Info("success", slog.Int("conflicts resolved", len(results)))
	return nil
}

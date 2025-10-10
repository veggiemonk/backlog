package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
	"github.com/veggiemonk/backlog/internal/paths"
)

var (
	doctorJSON     bool
	doctorStrategy string
	doctorDryRun   bool
	doctorFix      bool
)

var doctorDescription = `
Diagnose and fix task ID conflicts that can occur when creating tasks
in separate Git branches. Conflicts arise when multiple branches generate the same task IDs.

This command provides conflict detection and resolution capabilities to maintain
task ID uniqueness and data integrity.

Conflict types detected:
- Duplicate IDs (same ID in multiple files)
- Orphaned children (tasks with non-existent parents)
- Invalid hierarchy (parent-child ID mismatch)
`

var doctorExamples = `
 backlog doctor                    # Detect conflicts in text format
 backlog doctor --json             # Detect conflicts in JSON format
 backlog doctor --fix              # Detect and automatically fix conflicts
 backlog doctor --fix --dry-run    # Show what would be fixed without making changes
 backlog doctor --fix --strategy=auto    # Use auto-renumbering strategy
`

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:     "doctor",
	Short:   "Diagnose and fix task ID conflicts",
	Long:    doctorDescription,
	Example: doctorExamples,
	RunE:    runDoctor,
}

func setDoctorFlags(_ *cobra.Command) {
	// Doctor command flags
	doctorCmd.Flags().BoolVarP(&doctorJSON, "json", "j", false, "Output in JSON format")
	doctorCmd.Flags().BoolVar(&doctorFix, "fix", false, "Automatically fix detected conflicts")
	doctorCmd.Flags().StringVar(&doctorStrategy, "strategy", "chronological", "Resolution strategy when using --fix (chronological|auto|manual)")
	doctorCmd.Flags().BoolVar(&doctorDryRun, "dry-run", false, "Show what would be changed without making changes (use with --fix)")
}

func init() {
	// Set flags
	setDoctorFlags(doctorCmd)

	// Add to root
	rootCmd.AddCommand(doctorCmd)
}

func runDoctor(cmd *cobra.Command, args []string) error {
	fs := afero.NewOsFs()
	tasksDir := viper.GetString("folder")
	if doctorFix {
		// If --fix flag is provided, run resolve instead of detect
		return resolveConflicts(fs, tasksDir)
	}
	// Otherwise, just detect conflicts
	return detectConflicts(cmd.OutOrStdout(), fs, tasksDir)
}

func detectConflicts(w io.Writer, fs afero.Fs, tasksDir string) error {
	// Get tasks directory using the same approach as other commands
	var err error
	tasksDir, err = paths.ResolveTasksDir(fs, tasksDir)
	if err != nil {
		return fmt.Errorf("failed to resolve tasks directory: %w", err)
	}

	detector := core.NewConflictDetector(fs, tasksDir)

	// Detect conflicts
	conflicts, err := detector.DetectConflicts()
	if err != nil {
		return fmt.Errorf("failed to detect conflicts: %w", err)
	}

	if doctorJSON {
		output := map[string]any{
			"conflicts": conflicts,
			"summary":   core.SummarizeConflicts(conflicts),
		}
		if err := json.NewEncoder(w).Encode(output); err != nil {
			return fmt.Errorf("failed to encode JSON: %w", err)
		}
		return nil
	}

	// Text output
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

func resolveConflicts(fs afero.Fs, tasksDir string) error {
	// Get tasks directory using the same approach as other commands
	var err error
	tasksDir, err = paths.ResolveTasksDir(fs, tasksDir)
	if err != nil {
		return fmt.Errorf("failed to resolve tasks directory: %w", err)
	}

	detector := core.NewConflictDetector(fs, tasksDir)

	// Detect conflicts first
	conflicts, err := detector.DetectConflicts()
	if err != nil {
		return fmt.Errorf("failed to detect conflicts: %w", err)
	}

	if len(conflicts) == 0 {
		logging.Info("no conflicts to resolve")
		return nil
	}

	// Create resolver
	store := core.NewFileTaskStore(fs, tasksDir)
	resolver := core.NewConflictResolver(detector, store)

	// Parse strategy
	var strategy core.ResolutionStrategy
	switch doctorStrategy {
	case "chronological":
		strategy = core.ResolutionStrategyChronological
	case "auto":
		strategy = core.ResolutionStrategyAutoRenumber
	case "manual":
		strategy = core.ResolutionStrategyManual
	default:
		return fmt.Errorf("invalid strategy: %s", doctorStrategy)
	}

	// Create resolution plan
	plan, err := resolver.CreateResolutionPlan(conflicts, strategy)
	if err != nil {
		return fmt.Errorf("failed to create resolution plan: %w", err)
	}

	logging.Info("resolution plan", "strategy", doctorStrategy, "plan", plan.Summary)
	if len(plan.Actions) == 0 {
		logging.Info("no actions required")
		return nil
	}

	// Show planned actions
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

	if doctorDryRun {
		logging.Info("DRY RUN - No changes were made")
		return nil
	}

	if strategy == core.ResolutionStrategyManual {
		logging.Warn("manual resolution required - no automatic actions taken")
		logging.Warn("please review the conflicts above and resolve them manually")
		return nil
	}

	// Execute the plan with reference updates
	logging.Info("executing resolution plan")
	results, err := resolver.ExecuteResolutionPlanWithReferences(plan, false)
	if err != nil {
		return fmt.Errorf("failed to execute resolution plan: %w", err)
	}

	// Show results
	for _, result := range results {
		logging.Info(result)
	}

	logging.Info("success", slog.Int("conflicts resolved", len(results)))
	return nil
}

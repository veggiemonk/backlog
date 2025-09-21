package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
	"github.com/veggiemonk/backlog/internal/paths"
)

var (
	conflictsJSON     bool
	conflictsStrategy string
	conflictsDryRun   bool
)

// conflictsCmd represents the conflicts command
var conflictsCmd = &cobra.Command{
	Use:   "conflicts",
	Short: "Manage task ID conflicts",
	Long: `Detect and resolve task ID conflicts that can occur when creating tasks
in separate Git branches. Conflicts arise when multiple branches generate the same task IDs.

This command provides conflict detection and resolution capabilities to maintain
task ID uniqueness and data integrity.`,
}

// conflictsDetectCmd detects conflicts
var conflictsDetectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect task ID conflicts",
	Long: `Scan all task files and identify ID conflicts including:
- Duplicate IDs (same ID in multiple files)
- Orphaned children (tasks with non-existent parents)
- Invalid hierarchy (parent-child ID mismatch)

Examples:
  backlog conflicts detect                 # Detect conflicts in text format
  backlog conflicts detect --json          # Detect conflicts in JSON format`,
	Run: detectConflicts,
}

// conflictsResolveCmd resolves conflicts
var conflictsResolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Resolve task ID conflicts",
	Long: `Create and execute a resolution plan for detected conflicts.

Resolution strategies:
- chronological: Keep older tasks, renumber newer ones (default)
- auto: Automatically renumber conflicting IDs
- manual: Create plan requiring manual intervention

Examples:
  backlog conflicts resolve                            # Resolve using chronological strategy
  backlog conflicts resolve --strategy=auto           # Resolve using auto strategy
  backlog conflicts resolve --dry-run                 # Show what would be changed
  backlog conflicts resolve --strategy=manual         # Create manual resolution plan`,
	Run: resolveConflicts,
}

func detectConflicts(cmd *cobra.Command, args []string) {
	// Get tasks directory using the same approach as other commands
	fs := afero.NewOsFs()
	tasksDir := viper.GetString("folder")
	var err error
	tasksDir, err = paths.ResolveTasksDir(fs, tasksDir)
	if err != nil {
		logging.Error("failed to resolve tasks directory", "error", err)
		os.Exit(1)
	}

	detector := core.NewConflictDetector(fs, tasksDir)

	// Detect conflicts
	conflicts, err := detector.DetectConflicts()
	if err != nil {
		logging.Error("failed to detect conflicts", "error", err)
		os.Exit(1)
	}

	if conflictsJSON {
		output := map[string]any{
			"conflicts": conflicts,
			"summary":   core.SummarizeConflicts(conflicts),
		}
		if err := json.NewEncoder(cmd.OutOrStdout()).Encode(output); err != nil {
			logging.Error("failed to encode JSON", "error", err)
			os.Exit(1)
		}
		return
	}

	// Text output
	summary := core.SummarizeConflicts(conflicts)
	if summary.TotalConflicts == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "✅ No task ID conflicts detected")
		return
	}

	fmt.Fprintf(cmd.OutOrStdout(), "⚠️  Found %d task ID conflicts:\n\n", summary.TotalConflicts)

	if summary.DuplicateIDs > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "📋 Duplicate IDs: %d\n", summary.DuplicateIDs)
		for _, conflict := range summary.ConflictsByType[core.ConflictTypeDuplicateID] {
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s appears in: %v\n", conflict.ConflictID.String(), conflict.Files)
		}
		fmt.Fprintln(cmd.OutOrStdout())
	}

	if summary.OrphanedChildren > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "👤 Orphaned Children: %d\n", summary.OrphanedChildren)
		for _, conflict := range summary.ConflictsByType[core.ConflictTypeOrphanedChild] {
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s references non-existent parent\n", conflict.ConflictID.String())
		}
		fmt.Fprintln(cmd.OutOrStdout())
	}

	if summary.InvalidHierarchy > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "🔗 Invalid Hierarchy: %d\n", summary.InvalidHierarchy)
		for _, conflict := range summary.ConflictsByType[core.ConflictTypeInvalidHierarchy] {
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s has incorrect parent structure\n", conflict.ConflictID.String())
		}
		fmt.Fprintln(cmd.OutOrStdout())
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Run 'backlog conflicts resolve' to fix these conflicts.")
}

func resolveConflicts(cmd *cobra.Command, args []string) {
	if err := resolveConflictsImpl(cmd, args); err != nil {
		logging.Error("failed to resolve conflicts", "error", err)
		os.Exit(1)
	}
}

func resolveConflictsImpl(cmd *cobra.Command, _ []string) error {
	// Get tasks directory using the same approach as other commands
	fs := afero.NewOsFs()
	tasksDir := viper.GetString("folder")
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
		fmt.Fprintln(cmd.OutOrStdout(), "✅ No conflicts to resolve")
		return nil
	}

	// Create resolver
	store := core.NewFileTaskStore(fs, tasksDir)
	resolver := core.NewConflictResolver(detector, store)

	// Parse strategy
	var strategy core.ResolutionStrategy
	switch conflictsStrategy {
	case "chronological":
		strategy = core.ResolutionStrategyChronological
	case "auto":
		strategy = core.ResolutionStrategyAutoRenumber
	case "manual":
		strategy = core.ResolutionStrategyManual
	default:
		return fmt.Errorf("invalid strategy: %s", conflictsStrategy)
	}

	// Create resolution plan
	plan, err := resolver.CreateResolutionPlan(conflicts, strategy)
	if err != nil {
		return fmt.Errorf("failed to create resolution plan: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "📋 Resolution Plan (%s strategy):\n", conflictsStrategy)
	fmt.Fprintf(cmd.OutOrStdout(), "   %s\n\n", plan.Summary)

	if len(plan.Actions) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No actions required")
		return nil
	}

	// Show planned actions
	for i, action := range plan.Actions {
		fmt.Fprintf(cmd.OutOrStdout(), "%d. %s\n", i+1, action.Description)
		if action.Type == "manual" {
			fmt.Fprintln(cmd.OutOrStdout(), "   Type: Manual intervention required")
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "   Type: %s\n", action.Type)
			if !action.OriginalID.IsZero() {
				fmt.Fprintf(cmd.OutOrStdout(), "   Original ID: %s\n", action.OriginalID.String())
			}
			if !action.NewID.IsZero() {
				fmt.Fprintf(cmd.OutOrStdout(), "   New ID: %s\n", action.NewID.String())
			}
			if action.FilePath != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "   File: %s\n", action.FilePath)
			}
		}
		fmt.Fprintln(cmd.OutOrStdout())
	}

	if conflictsDryRun {
		fmt.Fprintln(cmd.OutOrStdout(), "🔍 DRY RUN - No changes were made")
		return nil
	}

	if strategy == core.ResolutionStrategyManual {
		fmt.Fprintln(cmd.OutOrStdout(), "⚠️  Manual resolution required - no automatic actions taken")
		fmt.Fprintln(cmd.OutOrStdout(), "Please review the conflicts above and resolve them manually")
		return nil
	}

	// Execute the plan with reference updates
	fmt.Fprintf(cmd.OutOrStdout(), "🔧 Executing resolution plan...\n\n")
	results, err := resolver.ExecuteResolutionPlanWithReferences(plan, false)
	if err != nil {
		return fmt.Errorf("failed to execute resolution plan: %w", err)
	}

	// Show results
	for _, result := range results {
		fmt.Fprintf(cmd.OutOrStdout(), "✅ %s\n", result)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "\n🎉 Successfully resolved %d conflicts\n", len(results))
	return nil
}

func setConflictsFlags(_ *cobra.Command) {
	// Detect flags
	conflictsDetectCmd.Flags().BoolVarP(&conflictsJSON, "json", "j", false, "Output in JSON format")

	// Resolve flags
	conflictsResolveCmd.Flags().StringVar(&conflictsStrategy, "strategy", "chronological", "Resolution strategy (chronological|auto|manual)")
	conflictsResolveCmd.Flags().BoolVar(&conflictsDryRun, "dry-run", false, "Show what would be changed without making changes")
}

func init() {
	// Add subcommands
	conflictsCmd.AddCommand(conflictsDetectCmd)
	conflictsCmd.AddCommand(conflictsResolveCmd)

	// Set flags
	setConflictsFlags(conflictsCmd)

	// Add to root
	rootCmd.AddCommand(conflictsCmd)
}


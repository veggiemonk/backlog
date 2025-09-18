package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/veggiemonk/backlog/internal/commit"
	"github.com/veggiemonk/backlog/internal/logging"
)

var archiveCmd = &cobra.Command{
	Use:   "archive <task-id>",
	Short: "Archive a task",
	Long:  `Archives a task, moving it to the archived directory.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runArchive,
}

func init() {
	rootCmd.AddCommand(archiveCmd)
}

func runArchive(cmd *cobra.Command, args []string) error {
	store := cmd.Context().Value(ctxKeyStore).(TaskStore)
	task, err := store.Get(args[0])
	if err != nil {
		return fmt.Errorf("get task %q: %v", args[0], err)
	}

	// Save path for commit
	oldPath := store.Path(task)

	newPath, err := store.Archive(task.ID)
	if err != nil {
		return fmt.Errorf("archive task %q: %v", task.ID.String(), err)
	}

	logging.Info("task archived successfully", "task_id", task.ID)
	// fmt.Printf("Task %s archived successfully.\n", archivedTask.ID)

	if !viper.GetBool(configAutoCommit) {
		return nil // Auto-commit is disabled
	}
	// Auto-commit the change if enabled
	commitMsg := fmt.Sprintf("chore(task): archive %s - \"%s\"", task.ID, task.Title)
	if err := commit.Add(newPath, oldPath, commitMsg); err != nil {
		logging.Warn("auto-commit failed", "task_id", task.ID, "error", err)
	}

	return nil
}

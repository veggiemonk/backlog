package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/veggiemonk/backlog/internal/commit"
	"github.com/veggiemonk/backlog/internal/logging"
)

func newArchiveCommand(rt *runtime) *cli.Command {
	return &cli.Command{
		Name:        "archive",
		Usage:       "Archive a task",
		UsageText:   "backlog archive <task-id>",
		ArgsUsage:   "<task-id>",
		Description: "Archives a task, moving it to the archived directory.",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				return cli.Exit("archive requires exactly one <task-id> argument", 1)
			}

			store := rt.store
			if store == nil {
				return fmt.Errorf("task store not initialized")
			}

			taskID := cmd.Args().First()
			task, err := store.Get(taskID)
			if err != nil {
				return fmt.Errorf("get task %q: %v", taskID, err)
			}

			oldPath := store.Path(task)

			newPath, err := store.Archive(task.ID)
			if err != nil {
				return fmt.Errorf("archive task %q: %v", task.ID.String(), err)
			}

			logging.Info("task archived successfully", "task_id", task.ID)

			if !rt.autoCommit {
				return nil
			}

			commitMsg := fmt.Sprintf("chore(task): archive %s - \"%s\"", task.ID, task.Title)
			if err := commit.Add(newPath, oldPath, commitMsg); err != nil {
				logging.Warn("auto-commit failed", "task_id", task.ID, "error", err)
			}
			return nil
		},
	}
}

package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/veggiemonk/backlog/internal/commit"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
)

const editDescription = `
Edit an existing task by providing its ID and flags for the fields to update.

Examples:
` +
	"```" +
	`

  # Change title
  backlog edit T42 -t "Fix the main login button"

  # Update description
  backlog edit T42 -d "The login button is misaligned on mobile"

  # Change status
  backlog edit T42 -s "in-progress"
  backlog edit T42 -s "done"

  # Assign/re-assign users
  backlog edit T42 -a "jordan"                    # Add jordan
  backlog edit T42 -a "jordan" -a "casey"         # Add multiple users
  backlog edit T42 --remove-assigned "alex"       # Remove alex

  # Update labels
  backlog edit T42 -l "bug,frontend"              # Add labels
  backlog edit T42 --remove-labels "backend"      # Remove label

  # Change priority
  backlog edit T42 --priority "high"

  # Manage acceptance criteria
  backlog edit T42 --ac "Button centered on mobile"     # Add new AC
  backlog edit T42 --check-ac 1                         # Check AC #1
  backlog edit T42 --uncheck-ac 1                       # Uncheck AC #1
  backlog edit T42 --remove-ac 2                        # Remove AC #2

  # Change parent
  backlog edit T42 -p "T18"                       # Make it a sub-task of T18
  backlog edit T42 -p ""                          # Remove parent

  # Add implementation notes
  backlog edit T42 --notes "Fixed in main.css, line 234"

  # Update implementation plan
  backlog edit T42 --plan "1. Refactor login\\n2. Test on mobile\\n3. Review"

  # Set dependencies
  backlog edit T42 --deps "T15"                   # Single dependency
  backlog edit T42 --deps "T15" --deps "T18"      # Multiple dependencies

  # Complex example (multiple changes at once)
  backlog edit T42 \
    -s "in-review" \
    -a "alex" \
    --priority "critical" \
    --notes "Ready for review on iOS and Android" \
    --check-ac 1 \
    --check-ac 2

` + "```"

func newEditCommand(rt *runtime) *cli.Command {
	return &cli.Command{
		Name:                  "edit",
		Usage:                 "Edit an existing task",
		UsageText:             "backlog edit <id>",
		ArgsUsage:             "<id>",
		EnableShellCompletion: true,
		Description:           editDescription,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "title", Aliases: []string{"t"}, Usage: "New title for the task"},
			&cli.StringFlag{Name: "description", Aliases: []string{"d"}, Usage: "New description for the task"},
			&cli.StringFlag{Name: "status", Aliases: []string{"s"}, Usage: "New status for the task"},
			&cli.StringFlag{Name: "priority", Usage: "New priority for the task"},
			&cli.StringFlag{Name: "parent", Aliases: []string{"p"}, Usage: "New parent for the task"},
			&cli.StringSliceFlag{Name: "assigned", Aliases: []string{"a"}, Usage: "Add assigned names for the task (can be used multiple times)"},
			&cli.StringSliceFlag{Name: "remove-assigned", Aliases: []string{"A"}, Usage: "Assigned names to remove from the task (can be used multiple times)"},
			&cli.StringSliceFlag{Name: "labels", Aliases: []string{"l"}, Usage: "Add labels for the task (can be used multiple times)"},
			&cli.StringSliceFlag{Name: "remove-labels", Aliases: []string{"L"}, Usage: "Labels to remove from the task (can be used multiple times)"},
			&cli.StringSliceFlag{Name: "deps", Usage: "Set dependencies (can be used multiple times)"},
			&cli.StringFlag{Name: "notes", Usage: "New implementation notes for the task"},
			&cli.StringFlag{Name: "plan", Usage: "New implementation plan for the task"},
			&cli.StringSliceFlag{Name: "ac", Usage: "Add a new acceptance criterion (can be used multiple times)"},
			&cli.IntSliceFlag{Name: "check-ac", Usage: "Check an acceptance criterion by its index"},
			&cli.IntSliceFlag{Name: "uncheck-ac", Usage: "Uncheck an acceptance criterion by its index"},
			&cli.IntSliceFlag{Name: "remove-ac", Usage: "Remove an acceptance criterion by its index"},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				return cli.Exit("edit requires exactly one <id> argument", 1)
			}

			store := rt.store
			if store == nil {
				return fmt.Errorf("task store not initialized")
			}

			params := core.EditTaskParams{ID: cmd.Args().First()}

			if cmd.IsSet("title") {
				v := cmd.String("title")
				params.NewTitle = &v
			}
			if cmd.IsSet("description") {
				v := cmd.String("description")
				params.NewDescription = &v
			}
			if cmd.IsSet("status") {
				v := cmd.String("status")
				params.NewStatus = &v
			}
			if cmd.IsSet("priority") {
				v := cmd.String("priority")
				params.NewPriority = &v
			}
			if cmd.IsSet("parent") {
				v := cmd.String("parent")
				params.NewParent = &v
			}
			if cmd.IsSet("deps") {
				params.NewDependencies = cmd.StringSlice("deps")
			}
			if cmd.IsSet("assigned") {
				params.AddAssigned = cmd.StringSlice("assigned")
			}
			if cmd.IsSet("remove-assigned") {
				params.RemoveAssigned = cmd.StringSlice("remove-assigned")
			}
			if cmd.IsSet("labels") {
				params.AddLabels = cmd.StringSlice("labels")
			}
			if cmd.IsSet("remove-labels") {
				params.RemoveLabels = cmd.StringSlice("remove-labels")
			}
			if cmd.IsSet("notes") {
				v := cmd.String("notes")
				params.NewNotes = &v
			}
			if cmd.IsSet("plan") {
				v := cmd.String("plan")
				params.NewPlan = &v
			}

			params.AddAC = cmd.StringSlice("ac")
			params.CheckAC = cmd.IntSlice("check-ac")
			params.UncheckAC = cmd.IntSlice("uncheck-ac")
			params.RemoveAC = cmd.IntSlice("remove-ac")

			task, err := store.Get(params.ID)
			if err != nil {
				return fmt.Errorf("failed to retrieve task %q: %w", params.ID, err)
			}

			oldFilePath := store.Path(task)

			if err := store.Update(&task, params); err != nil {
				return fmt.Errorf("failed to update task %q: %w", params.ID, err)
			}

			logging.Info("task updated successfully", "task_id", task.ID)

			if !rt.autoCommit {
				return nil
			}

			currentFilePath := store.Path(task)
			if oldFilePath == currentFilePath {
				oldFilePath = ""
			}

			commitMsg := fmt.Sprintf("feat(task): edit %s - \"%s\"", task.ID, task.Title)
			if err := commit.Add(currentFilePath, oldFilePath, commitMsg); err != nil {
				logging.Warn("auto-commit failed", "task_id", task.ID, "error", err)
			}
			return nil
		},
	}
}

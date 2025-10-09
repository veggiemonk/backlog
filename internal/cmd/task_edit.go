package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/veggiemonk/backlog/internal/commit"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
)

const editExamples = `
# Edit tasks using the "backlog edit" command with its different flags.
# Let's assume you have a task with ID "42" that you want to modify.
# Here are some examples of how to use this command effectively:

# 1. Changing the Title
# Use the -t or --title flag to give the task a new title.
backlog edit 42 -t "Fix the main login button styling"

# 2. Updating the Description
# Use the -d or --description flag to replace the existing description with a new one.
backlog edit 42 -d "The login button on the homepage is misaligned on mobile devices. It should be centered."

# 3. Changing the Status
# Update the task's progress by changing its status with the -s or --status flag.
backlog edit 42 -s "in-progress"

# 4. Re-assigning a Task
# You can change the assigned names for a task using the -a or --assignee flag.
# This will replace the current list of assigned names.
# Assign to a single person:
backlog edit 42 -a "jordan"
# Assign to multiple people:
backlog edit 42 -a "jordan" -a "casey"

# 5. Updating Labels
# Use the -l or --labels flag to replace the existing labels.
backlog edit 42 -l "bug,frontend"

# 6. Changing the Priority
# Adjust the task's priority with the --priority flag.
backlog edit 42 --priority "high"

# 7. Managing Acceptance Criteria
# You can add, check, uncheck, or remove acceptance criteria.
# Add a new AC:
backlog edit 42 --ac "The button is centered on screens smaller than 576px."
# Check the first AC (assuming it's at index 1):
backlog edit 42 --check-ac 1
# Uncheck the first AC:
backlog edit 42 --uncheck-ac 1
# Remove the second AC (at index 2):
backlog edit 42 --remove-ac 2

# 8. Changing the Parent Task
# Move a task to be a sub-task of a different parent using the -p or --parent flag.
backlog edit 42 -p "18"
# To remove a parent, pass an empty string:
backlog edit 42 -p ""

# 9. Adding Implementation Notes
# Use the --notes flag to add or update technical notes for implementation.
backlog edit 42 --notes "The issue is in the 'main.css' file, specifically in the '.login-container' class. Need to adjust the media query."

# 10. Complex Example (Combining Multiple Flags)
# You can combine several flags to make multiple changes at once.
backlog edit 42 \
  -s "in-review" \
  -a "alex" \
  --priority "critical" \
  --notes "The fix is ready for review. Please check on both iOS and Android." \
  --check-ac 1 \
  --check-ac 2

# 11. Updating the Implementation Plan
# Use the --plan flag to add or update the implementation plan for the task.
backlog edit 42 --plan "1. Refactor login button\\n2. Test on mobile\\n3. Review with team"

# 12. Adding Dependencies
# Use the --deps flag to add one or more task dependencies.
# This will replace all existing dependencies with the new ones.
backlog edit 42 --deps "T1" --deps "T2"

# 13. Setting a Single Dependency
# If you want to make a task depend on another specific task:
backlog edit 42 --deps "T15"
# This makes task 42 dependent on task T15, meaning T15 must be completed before T42 can be started.

# 14. Setting Multiple Dependencies
# You can make a task depend on multiple other tasks:
backlog edit 42 --deps "T15" --deps "T18" --deps "T20"
# This makes task 42 dependent on tasks T15, T18, and T20.
# 15. Editing the construction plan
backlog edit 42 --plan "1. Dig hole 2. Pour foundation"
`

func newEditCommand(rt *runtime) *cli.Command {
	return &cli.Command{
		Name:      "edit",
		Usage:     "Edit an existing task",
		UsageText: "backlog edit <id>",
		ArgsUsage: "<id>",
		Description: "Edit an existing task by providing its ID and flags for the fields to update.\n\nExamples:\n" +
			editExamples,
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

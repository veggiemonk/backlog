package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/veggiemonk/backlog/internal/commit"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
)

const createExamples = `
# Create tasks using the "backlog create" command with its different flags.
# Here are some examples of how to use this command effectively:
# 1. Basic Task Creation
# This is the simplest way to create a task, providing only a title.
backlog create "Fix the login button styling"

# 2. Task with a Description. Use the -d or --description flag to add more detailed information about the task.
backlog create "Implement password reset" -d "Users should be able to request a password reset link via their email. This involves creating a new API endpoint and a front-end form."

# 3. Assigning a Task. You can assign a task to one or more team members using the -a or --assigned flag.
# Assign to a single person:
backlog create "Design the new dashboard" -a "alex"
# Assign to multiple people:
backlog create "Code review for the payment gateway" -a "jordan" -a "casey"

# 4. Adding Labels. Use the -l or --labels flag to categorize the task with comma-separated labels.
backlog create "Update third-party dependencies" -l "bug,backend,security"

# 5. Setting a Priority
# Specify the task's priority with the --priority flag. The default is "medium".
backlog create "Hotfix: Production database is down" --priority "high"
backlog create "Refactor the old user model" --priority "low"

# 6. Defining Acceptance Criteria
# Use the --ac flag multiple times to list the conditions that must be met for the task to be considered complete.
backlog create "Develop user profile page" \
  --ac "Users can view their own profile information." \
  --ac "Users can upload a new profile picture." \
  --ac "The page is responsive on mobile devices."

# 7. Creating a Sub-task. Link a new task to a parent task using the -p or --parent flag. This is useful for breaking down larger tasks.
# First, create the parent task
backlog create "Implement User Authentication"
# Now, create a sub-task (assuming the parent task ID is 15)
backlog create "Add Google OAuth login" -p "15"

# 8. Setting Task Dependencies
# Use the --deps flag to specify that this task depends on other tasks being completed first.
# Single dependency:
backlog create "Deploy user authentication" --deps "T15"
# Multiple dependencies:
backlog create "Integration testing" --deps "T15" --deps "T18" --deps "T20"
# This means the task cannot be started until tasks T15, T18, and T20 are completed.

# 9. Complex Example (Combining Multiple Flags). Here is a comprehensive example that uses several flags at once to create a very detailed task.
backlog create "Build the new reporting feature" \
  -d "Create a new section in the app that allows users to generate and export monthly performance reports in PDF format." \
  -a "drew" \
  -l "feature,frontend,backend" \
  --priority "high" \
  --ac "Report generation logic is accurate." \
  --ac "Users can select a date range for the report." \
  --ac "The exported PDF has the correct branding and layout." \
  -p "23"
`

func newCreateCommand(rt *runtime) *cli.Command {
	return &cli.Command{
		Name:      "create",
		Usage:     "Create a new task",
		UsageText: "backlog create <title>",
		ArgsUsage: "<title>",
		Description: "Creates a new task in the backlog.\n\nExamples:\n" +
			createExamples,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "description", Aliases: []string{"d"}, Usage: "Description of the task"},
			&cli.StringFlag{Name: "parent", Aliases: []string{"p"}, Usage: "Parent task ID"},
			&cli.StringFlag{Name: "priority", Usage: "Priority of the task (low, medium, high, critical)", Value: "medium"},
			&cli.StringSliceFlag{Name: "assigned", Aliases: []string{"a"}, Usage: "Assignee for the task (can be specified multiple times)"},
			&cli.StringSliceFlag{Name: "labels", Aliases: []string{"l"}, Usage: "Comma-separated labels for the task"},
			&cli.StringSliceFlag{Name: "deps", Usage: "Add a dependency (can be used multiple times)"},
			&cli.StringSliceFlag{Name: "ac", Usage: "Acceptance criterion (can be specified multiple times)"},
			&cli.StringFlag{Name: "plan", Usage: "Implementation plan for the task"},
			&cli.StringFlag{Name: "notes", Usage: "Additional notes for the task"},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				return cli.Exit("create requires exactly one <title> argument", 1)
			}

			store := rt.store
			if store == nil {
				return fmt.Errorf("task store not initialized")
			}

			params := core.CreateTaskParams{
				Title:        cmd.Args().First(),
				Description:  cmd.String("description"),
				Parent:       cmd.String("parent"),
				Priority:     cmd.String("priority"),
				Assigned:     cmd.StringSlice("assigned"),
				Labels:       cmd.StringSlice("labels"),
				Dependencies: cmd.StringSlice("deps"),
				AC:           cmd.StringSlice("ac"),
				Plan:         cmd.String("plan"),
				Notes:        cmd.String("notes"),
			}

			newTask, err := store.Create(params)
			if err != nil {
				return fmt.Errorf("failed to create task: %w", err)
			}

			logging.Info("task created successfully", "task_id", newTask.ID)

			if !rt.autoCommit {
				return nil
			}

			commitMsg := fmt.Sprintf("feat(task): create %s - \"%s\"", newTask.ID, newTask.Title)
			if err := commit.Add(store.Path(newTask), "", commitMsg); err != nil {
				logging.Warn("auto-commit failed", "task_id", newTask.ID, "error", err)
			}
			return nil
		},
	}
}

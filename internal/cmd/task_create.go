package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/veggiemonk/backlog/internal/commit"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
)

const createDescription = `
Creates a new task in the backlog.

Examples:
` +
	"```" +
	`

  # Basic task
  backlog create "Fix the login button styling"

  # Task with description
  backlog create "Implement password reset" \
    -d "Users should be able to request a password reset link via email"

  # Assign to team members
  backlog create "Design the new dashboard" -a "alex"
  backlog create "Code review" -a "jordan" -a "casey"    # Multiple assignees

  # Add labels
  backlog create "Update dependencies" -l "bug,backend,security"

  # Set priority (low, medium, high, critical)
  backlog create "Hotfix: Database down" --priority "high"
  backlog create "Refactor old code" --priority "low"

  # Define acceptance criteria
  backlog create "Develop user profile page" \
    --ac "Users can view their profile" \
    --ac "Users can upload a profile picture" \
    --ac "Page is responsive on mobile"

  # Create sub-task with parent
  backlog create "Implement User Authentication"           # Creates T01
  backlog create "Add Google OAuth" -p "T01"               # Creates T01.01

  # Set dependencies
  backlog create "Deploy to production" --deps "T15"       # Single dependency
  backlog create "Integration testing" \
    --deps "T15" --deps "T18" --deps "T20"                 # Multiple dependencies

  # Complex example with multiple flags
  backlog create "Build reporting feature" \
    -d "Monthly performance reports in PDF format" \
    -a "drew" \
    -l "feature,frontend,backend" \
    --priority "high" \
    --ac "Report generation logic is accurate" \
    --ac "Users can select date range" \
    --ac "PDF has correct branding" \
    -p "23"

` + "```"

func newCreateCommand(rt *runtime) *cli.Command {
	return &cli.Command{
		Name:                  "create",
		Usage:                 "Create a new task",
		UsageText:             "backlog create <title>",
		ArgsUsage:             "<title>",
		Description:           createDescription,
		EnableShellCompletion: true,
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

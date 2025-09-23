package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/veggiemonk/backlog/internal/commit"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
	mcpserver "github.com/veggiemonk/backlog/internal/mcp"
)

var createCmd = &cobra.Command{
	Use:   "create <title>",
	Short: "Create a new task",
	Long:  `Creates a new task in the backlog.`,
	Args:  cobra.ExactArgs(1),
	Example: `
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
	`,
	RunE: runCreate,
}

var (
	description  string
	parent       string
	priority     string
	assigned     []string
	labels       []string
	dependencies []string
	ac           []string
	plan         string
	notes        string
)

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the task")
	createCmd.Flags().StringVarP(&parent, "parent", "p", "", "Parent task ID")
	createCmd.Flags().StringVar(&priority, "priority", "medium", "Priority of the task (low, medium, high, critical)")
	createCmd.Flags().StringSliceVarP(&assigned, "assigned", "a", []string{}, "Assignee for the task (can be specified multiple times)")
	createCmd.Flags().StringSliceVarP(&labels, "labels", "l", []string{}, "Comma-separated labels for the task")
	createCmd.Flags().StringSliceVar(&dependencies, "deps", []string{}, "Add a dependency (can be used multiple times)")
	createCmd.Flags().StringSliceVar(&ac, "ac", []string{}, "Acceptance criterion (can be specified multiple times)")
	createCmd.Flags().StringVar(&plan, "plan", "", "Implementation plan for the task")
	createCmd.Flags().StringVar(&notes, "notes", "", "Additional notes for the task")
}

func runCreate(cmd *cobra.Command, args []string) error {
	params := core.CreateTaskParams{
		Title:        args[0],
		Description:  description,
		Parent:       parent,
		Priority:     priority,
		Assigned:     assigned,
		Labels:       labels,
		Dependencies: dependencies,
		AC:           ac,
		Plan:         plan,
		Notes:        notes,
	}

	store := cmd.Context().Value(ctxKeyStore).(mcpserver.TaskStore)
	newTask, err := store.Create(params)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	logging.Info("task created successfully", "task_id", newTask.ID)

	if !viper.GetBool(configAutoCommit) {
		return nil // Auto-commit is disabled
	}
	// Auto-commit the change if enabled
	filePath := store.Path(newTask)
	commitMsg := fmt.Sprintf("feat(task): create %s - \"%s\"", newTask.ID, newTask.Title)
	if err := commit.Add(filePath, "", commitMsg); err != nil {
		logging.Warn("auto-commit failed", "task_id", newTask.ID, "error", err)
	}
	return nil
}

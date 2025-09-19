package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/veggiemonk/backlog/internal/commit"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
	"github.com/veggiemonk/backlog/internal/validation"
)

var createCmd = &cobra.Command{
	Use:   "create <title>",
	Short: "Create a new task",
	Long:  `Creates a new task in the backlog.`,
	Args: func(cmd *cobra.Command, args []string) error {
		validator := validation.NewCLIValidator()
		return validator.ValidateArgs(cmd, args, 1)
	},
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
	Run: runCreate,
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

func runCreate(cmd *cobra.Command, args []string) {
	// Sanitize input parameters
	sanitizer := validation.NewSanitizer()

	params := core.CreateTaskParams{
		Title:        sanitizer.SanitizeTitle(args[0]),
		Description:  sanitizer.SanitizeDescription(description),
		Parent:       &parent,
		Priority:     priority,
		Assigned:     sanitizer.SanitizeSlice(assigned, sanitizer.SanitizeAssignee),
		Labels:       sanitizer.SanitizeSlice(labels, sanitizer.SanitizeLabel),
		Dependencies: sanitizer.SanitizeSlice(dependencies, sanitizer.SanitizeTaskID),
		AC:           sanitizer.SanitizeSlice(ac, sanitizer.SanitizeAcceptanceCriterion),
		Plan:         &plan,
		Notes:        &notes,
	}

	// Sanitize parent ID if provided
	if parent != "" {
		sanitizedParent := sanitizer.SanitizeTaskID(parent)
		params.Parent = &sanitizedParent
	}

	// Sanitize plan if provided
	if plan != "" {
		sanitizedPlan := sanitizer.SanitizePlan(plan)
		params.Plan = &sanitizedPlan
	}

	// Sanitize notes if provided
	if notes != "" {
		sanitizedNotes := sanitizer.SanitizeNotes(notes)
		params.Notes = &sanitizedNotes
	}

	// Validate input parameters
	validator := validation.NewCLIValidator()
	if validationErrors := validator.ValidateCreateParams(params); validationErrors.HasErrors() {
		logging.Error("validation failed", "errors", validationErrors.Error())
		for _, verr := range validationErrors {
			logging.Error("validation error", "field", verr.Field, "value", verr.Value, "message", verr.Message, "code", verr.Code)
		}
		os.Exit(1)
	}

	store := cmd.Context().Value(ctxKeyStore).(TaskStore)
	newTask, err := store.Create(params)
	if err != nil {
		logging.Error("failed to create task", "error", err)
		os.Exit(1)
	}

	logging.Info("task created successfully", "task_id", newTask.ID)

	if !viper.GetBool(configAutoCommit) {
		return // Auto-commit is disabled
	}
	// Auto-commit the change if enabled
	filePath := store.Path(newTask)
	commitMsg := fmt.Sprintf("feat(task): create %s - \"%s\"", newTask.ID, newTask.Title)
	if err := commit.Add(filePath, "", commitMsg); err != nil {
		logging.Warn("auto-commit failed", "task_id", newTask.ID, "error", err)
	}
}

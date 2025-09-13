package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/veggiemonk/backlog/internal/commit"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
)

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit an existing task",
	Long:  `Edit an existing task by providing its ID and flags for the fields to update.`,
	Args:  cobra.ExactArgs(1),
	Example: `
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
backlog edit 42 --plan "1. Refactor login button\n2. Test on mobile\n3. Review with team"

# 12. Adding Dependencies
# Use the --dep flag to add one or more task dependencies.
# This will replace all existing dependencies with the new ones.
backlog edit 42 --dep "T1" --dep "T2"

# 13. Setting a Single Dependency
# If you want to make a task depend on another specific task:
backlog edit 42 --dep "T15"
# This makes task 42 dependent on task T15, meaning T15 must be completed before T42 can be started.

# 14. Setting Multiple Dependencies
# You can make a task depend on multiple other tasks:
backlog edit 42 --dep "T15" --dep "T18" --dep "T20"
# This makes task 42 dependent on tasks T15, T18, and T20.
	`,
	Run: runEdit,
}

var (
	newTitle        string
	newDescription  string
	newStatus       string
	newPriority     string
	newParent       string
	addAssigned     []string
	removeAssigned  []string
	addLabels       []string
	removeLabels    []string
	newDependencies []string
	newNotes        string
	newPlan         string
	addAC           []string
	checkAC         []int
	uncheckAC       []int
	removeAC        []int
)

func init() {
	rootCmd.AddCommand(editCmd)
	setEditFlags(editCmd)
}

func setEditFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&newTitle, "title", "t", "", "New title for the task")
	cmd.Flags().StringVarP(&newDescription, "description", "d", "", "New description for the task")
	cmd.Flags().StringVarP(&newStatus, "status", "s", "", "New status for the task")
	cmd.Flags().StringVar(&newPriority, "priority", "", "New priority for the task")
	cmd.Flags().StringVarP(&newParent, "parent", "p", "", "New parent for the task")
	cmd.Flags().StringSliceVarP(&addAssigned, "assigned", "a", nil, "Add assigned names for the task (can be used multiple times)")
	cmd.Flags().StringSliceVarP(&removeAssigned, "remove-assigned", "A", nil, "Assigned names to remove from the task (can be used multiple times)")
	cmd.Flags().StringSliceVarP(&addLabels, "labels", "l", nil, "Add labels for the task (can be used multiple times)")
	cmd.Flags().StringSliceVarP(&removeLabels, "remove-labels", "L", nil, "Labels to remove from the task (can be used multiple times)")
	cmd.Flags().StringSliceVar(&newDependencies, "dep", nil, "Set dependencies (can be used multiple times)")
	cmd.Flags().StringVar(&newNotes, "notes", "", "New implementation notes for the task")
	cmd.Flags().StringVar(&newPlan, "plan", "", "New implementation plan for the task")

	// Acceptance Criteria flags
	cmd.Flags().StringSliceVar(&addAC, "ac", nil, "Add a new acceptance criterion (can be used multiple times)")
	cmd.Flags().IntSliceVar(&checkAC, "check-ac", nil, "Check an acceptance criterion by its index")
	cmd.Flags().IntSliceVar(&uncheckAC, "uncheck-ac", nil, "Uncheck an acceptance criterion by its index")
	cmd.Flags().IntSliceVar(&removeAC, "remove-ac", nil, "Remove an acceptance criterion by its index")
}

func runEdit(cmd *cobra.Command, args []string) {
	params := core.EditTaskParams{ID: args[0]}

	// Set optional pointers for fields that were changed
	if cmd.Flags().Changed("title") {
		params.NewTitle = &newTitle
	}
	if cmd.Flags().Changed("description") {
		params.NewDescription = &newDescription
	}
	if cmd.Flags().Changed("status") {
		params.NewStatus = &newStatus
	}
	if cmd.Flags().Changed("priority") {
		params.NewPriority = &newPriority
	}
	if cmd.Flags().Changed("parent") {
		params.NewParent = &newParent
	}
	if cmd.Flags().Changed("dep") {
		params.NewDependencies = newDependencies
	}
	if cmd.Flags().Changed("assigned") {
		params.AddAssigned = addAssigned
	}
	if cmd.Flags().Changed("remove-assigned") {
		params.RemoveAssigned = removeAssigned
	}
	// New labels
	if cmd.Flags().Changed("labels") {
		params.AddLabels = addLabels
	}
	// Remove labels
	if cmd.Flags().Changed("remove-labels") {
		params.RemoveLabels = removeLabels
	}
	// Other optional fields
	if cmd.Flags().Changed("notes") {
		params.NewNotes = &newNotes
	}
	if cmd.Flags().Changed("plan") {
		params.NewPlan = &newPlan
	}

	// AC params
	params.AddAC = addAC
	params.CheckAC = checkAC
	params.UncheckAC = uncheckAC
	params.RemoveAC = removeAC

	// get store from context
	store := cmd.Context().Value(ctxKeyStore).(TaskStore)

	// current task
	task, err := store.Get(params.ID)
	if err != nil {
		logging.Error("failed to retrieve task", "task_id", params.ID, "error", err)
		os.Exit(1)
	}
	// save the old path in case of a rename
	oldFilePath := store.Path(task)

	updatedTask, err := store.Update(task, params)
	if err != nil {
		logging.Error("failed to update task", "task_id", params.ID, "error", err)
		os.Exit(1)
	}

	defer func() {
		logging.Info("task updated successfully", "task_id", updatedTask.ID)
		// fmt.Printf("Task %s updated successfully.\n", updatedTask.ID)
	}()

	if !autoCommit {
		return // autocommit is disabled
	}

	// paths to commit
	currentFilePath := store.Path(updatedTask)
	if oldFilePath == currentFilePath {
		oldFilePath = ""
	}
	// autocommit the change if enabled
	commitMsg := fmt.Sprintf("feat(task): edit %s - \"%s\"", updatedTask.ID, updatedTask.Title)
	if err := commit.Add(currentFilePath, oldFilePath, commitMsg); err != nil {
		logging.Warn("auto-commit failed", "task_id", updatedTask.ID, "error", err)
	}
}

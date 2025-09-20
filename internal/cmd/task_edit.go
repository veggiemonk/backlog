package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/veggiemonk/backlog/internal/commit"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
)

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit an existing task",
	Long:  `Edit an existing task by providing its ID and flags for the fields to update.`,
	Args:  cobra.ExactArgs(1),
	Example: EditExamples.GenerateExampleText(),
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

	if !viper.GetBool(configAutoCommit) {
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

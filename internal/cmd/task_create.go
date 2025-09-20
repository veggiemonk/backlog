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

var createCmd = &cobra.Command{
	Use:   "create <title>",
	Short: "Create a new task",
	Long:  `Creates a new task in the backlog.`,
	Args:  cobra.ExactArgs(1),
	Example: CreateExamples.GenerateExampleText(),
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
	params := core.CreateTaskParams{
		Title:        args[0],
		Description:  description,
		Parent:       &parent,
		Priority:     priority,
		Assigned:     assigned,
		Labels:       labels,
		Dependencies: dependencies,
		AC:           ac,
		Plan:         &plan,
		Notes:        &notes,
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

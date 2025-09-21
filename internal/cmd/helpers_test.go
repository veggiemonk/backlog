package cmd

import (
	"bytes"
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/veggiemonk/backlog/internal/core"
)

type runEFunc func(cmd *cobra.Command, args []string) error

var testTasksDir = ".backlog"

func exec(t *testing.T, use string, run runEFunc, args ...string) ([]byte, error) {
	t.Helper()
	if use == "" {
		return nil, fmt.Errorf("'use' cannot be empty: %v", args)
	}
	fs := afero.NewMemMapFs()
	store := core.NewFileTaskStore(fs, testTasksDir)
	createTestTasks(t, store)
	// Create fresh command to avoid state pollution
	testRootCmd := &cobra.Command{Use: "backlog"}
	testRootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		cmd.SetContext(context.WithValue(cmd.Context(), ctxKeyStore, store))
	}
	var testCmd *cobra.Command

	// Use proper command configuration based on command type
	switch use {
	case "create":
		testCmd = &cobra.Command{Use: use, Args: cobra.ExactArgs(1), RunE: run}
	case "edit":
		testCmd = &cobra.Command{Use: use, Args: cobra.ExactArgs(1), RunE: run}
	case "view":
		testCmd = &cobra.Command{Use: use, Args: cobra.ExactArgs(1), RunE: run}
	case "archive":
		testCmd = &cobra.Command{Use: use, Args: cobra.ExactArgs(1), RunE: run}
	case "search":
		testCmd = &cobra.Command{Use: use, Args: cobra.ExactArgs(1), RunE: run}
	default:
		testCmd = &cobra.Command{Use: use, RunE: run}
	}

	setRootPersistentFlags(testRootCmd)
	testRootCmd.AddCommand(testCmd)
	switch use {
	case "create":
		setCreateFlags(testCmd)
	case "list":
		setListFlags(testCmd)
	case "search":
		setSearchFlags(testCmd)
	case "edit":
		setEditFlags(testCmd)
	case "view":
		setViewFlags(testCmd)
	case "archive":
		// Archive command has no flags
	case "version":
		// Version command has no flags
	case "instructions":
		// Instructions command has no flags
	case "mcp":
		setMCPFlags(testCmd)
	default:
		t.Fatalf("no command called %s", use)
	}
	args = slices.Insert(args, 0, use)
	args = append([]string{"--log-level", "error"}, args...)
	buf := new(bytes.Buffer)
	testRootCmd.SetOut(buf)
	testRootCmd.SetErr(buf)
	testRootCmd.SetArgs(args)

	err := testRootCmd.Execute()
	return bytes.TrimSpace(buf.Bytes()), err
}

func createTestTasks(t *testing.T, store TaskStore) {
	t.Helper()

	base := []string{"First", "Second", "Third", "Fourth", "Fifth"}
	labels := []string{"first", "second", "third", "fourth", "fifth"}
	assigned := []string{"first-user", "second-user", "third-user", "fourth-user", "fifth-user"}

	for i, s := range base {
		t1, _ := store.Create(core.CreateTaskParams{
			Title:       fmt.Sprintf("%s Task", s),
			Assigned:    []string{assigned[i]},
			Labels:      []string{labels[i]},
			Priority:    "medium",
			Notes:       &[]string{fmt.Sprintf("%s implementation notes.", s)}[0],
			Plan:        &[]string{fmt.Sprintf("%s implementation plan.", s)}[0],
			Description: fmt.Sprintf("%s description.", s),
			AC:          []string{fmt.Sprintf("%s AC.", s)},
		})
		_ = t1
	}

	t6, _ := store.Create(core.CreateTaskParams{
		Title:    "Unassigned Task",
		Assigned: nil,
		Labels:   []string{"sixth"},
	})
	_ = t6

	t7, _ := store.Create(core.CreateTaskParams{
		Title:       "Unlabeled Task",
		Description: "unlabeled description.",
		Assigned:    []string{"seventh-user"},
	})
	_ = t7
	t8, _ := store.Create(core.CreateTaskParams{
		Title:       "High Priority Task",
		Description: "high priority description.",
		Assigned:    []string{"eighth-user"},
		Labels:      []string{"eighth"},
		Priority:    "high",
	})
	_ = t8

	t9, err := store.Create(core.CreateTaskParams{
		Title:       "In Progress Task",
		Description: "in progress description.",
		Assigned:    []string{"nineth-user"},
		Labels:      []string{"nineth"},
	})
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	t9, err = store.Update(t9, core.EditTaskParams{NewStatus: ptr("in-progress")})
	if err != nil {
		t.Fatalf("failed to update task: %v", err)
	}
	_ = t9
}

const countTask = 9

func ptr[T any](v T) *T {
	return &v
}

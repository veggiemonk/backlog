package cmd

import (
	"bytes"
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/spf13/afero"
	"github.com/veggiemonk/backlog/internal/core"
	mcpserver "github.com/veggiemonk/backlog/internal/mcp"
)

func exec(t *testing.T, command string, args ...string) ([]byte, error) {
	t.Helper()
	if command == "" {
		return nil, fmt.Errorf("command cannot be empty: %v", args)
	}

	fs := afero.NewMemMapFs()
	tasksDir := ".backlog"
	store := core.NewFileTaskStore(fs, tasksDir)
	createTestTasks(t, store)

	cmd := NewCommand(
		WithFilesystem(fs),
		WithStore(store),
		WithTasksDir(tasksDir),
		WithSkipLogging(true),
	)
	buf := new(bytes.Buffer)
	cmd.Writer = buf
	cmd.ErrWriter = buf

	cliArgs := slices.Insert(args, 0, command)
	cliArgs = slices.Insert(cliArgs, 0, "backlog")

	err := cmd.Run(context.Background(), cliArgs)
	return bytes.TrimSpace(buf.Bytes()), err
}

func createTestTasks(t *testing.T, store mcpserver.TaskStore) {
	t.Helper()

	base := []string{"First", "Second", "Third", "Fourth", "Fifth"}
	labels := []string{"first", "second", "third", "fourth", "fifth"}
	assigned := []string{"first-user", "second-user", "third-user", "fourth-user", "fifth-user"}

	for i, s := range base {
		t1, _ := store.Create(core.CreateTaskParams{
			Title:       s + " Task",
			Assigned:    []string{assigned[i]},
			Labels:      []string{labels[i]},
			Priority:    "medium",
			Notes:       s + " implementation notes.",
			Plan:        s + " implementation plan.",
			Description: s + " description.",
			AC:          []string{s + " AC."},
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
		Assigned:    []string{"ninth-user"},
		Labels:      []string{"ninth"},
	})
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	if err = store.Update(&t9, core.EditTaskParams{NewStatus: ptr("in-progress")}); err != nil {
		t.Fatalf("failed to update task: %v", err)
	}
}

const countTask = 9

func ptr[T any](v T) *T {
	return &v
}

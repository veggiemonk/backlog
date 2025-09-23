package mcp

import (
	"testing"

	"github.com/matryer/is"
	"github.com/veggiemonk/backlog/internal/core"
)

func setupTestData(t *testing.T, store *core.FileTaskStore) {
	t.Helper()
	is := is.New(t)

	// Create a completed task (for weekly summary)
	completedTask, err := store.Create(core.CreateTaskParams{
		Title:       "Completed Feature",
		Description: "A feature that was completed",
		Labels:      []string{"feature", "backend"},
		Priority:    "high",
		Assigned:    []string{"alice"},
	})
	is.NoErr(err)

	// Mark it as done
	doneStatus := "done"
	is.NoErr(store.Update(&completedTask, core.EditTaskParams{
		ID:        completedTask.ID.String(),
		NewStatus: &doneStatus,
	}))

	// Create a high priority todo task
	_, err = store.Create(core.CreateTaskParams{
		Title:       "High Priority Feature",
		Description: "An urgent feature to implement",
		Labels:      []string{"feature", "urgent"},
		Priority:    "high",
		Assigned:    []string{"bob"},
	})
	is.NoErr(err)

	// Create a blocked task
	_, err = store.Create(core.CreateTaskParams{
		Title:       "Blocked Task",
		Description: "This task is blocked waiting for dependencies",
		Labels:      []string{"blocked", "backend"},
		Priority:    "medium",
		Assigned:    []string{"charlie"},
	})
	is.NoErr(err)

	// Create an unassigned task
	_, err = store.Create(core.CreateTaskParams{
		Title:       "Unassigned Task",
		Description: "This task has no assignee",
		Labels:      []string{"feature"},
		Priority:    "low",
	})
	is.NoErr(err)

	// Create an epic parent task
	epic, err := store.Create(core.CreateTaskParams{
		Title:       "Epic: User Authentication",
		Description: "Complete user authentication system",
		Labels:      []string{"epic", "auth"},
		Priority:    "high",
	})
	is.NoErr(err)

	// Create a subtask under the epic
	parentID := epic.ID.String()[1:] // Remove the "T" prefix
	_, err = store.Create(core.CreateTaskParams{
		Title:       "Login API endpoint",
		Description: "Create login endpoint",
		Labels:      []string{"api", "auth"},
		Priority:    "high",
		Parent:      parentID,
		Assigned:    []string{"alice"},
	})
	is.NoErr(err)

	// Create a bug report
	_, err = store.Create(core.CreateTaskParams{
		Title:       "Bug: Login fails on mobile",
		Description: "Users can't login on mobile devices",
		Labels:      []string{"bug", "mobile"},
		Priority:    "critical",
		Assigned:    []string{"david"},
	})
	is.NoErr(err)

	// Create an in-progress task
	inProgressTask, err := store.Create(core.CreateTaskParams{
		Title:       "In Progress Feature",
		Description: "Currently being worked on",
		Labels:      []string{"feature", "frontend"},
		Priority:    "medium",
		Assigned:    []string{"eve"},
	})
	is.NoErr(err)

	// Mark it as in-progress
	inProgressStatus := "in-progress"
	is.NoErr(store.Update(&inProgressTask, core.EditTaskParams{
		ID:        inProgressTask.ID.String(),
		NewStatus: &inProgressStatus,
	}))
}

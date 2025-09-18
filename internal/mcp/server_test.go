package mcp

import (
	"context"
	"sync"
	"testing"

	"github.com/matryer/is"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/afero"
	"github.com/veggiemonk/backlog/internal/core"
)

func TestServer(t *testing.T) {
	// Setup test store
	fs := afero.NewMemMapFs()
	store := core.NewFileTaskStore(fs, ".backlog")

	// Seed the store with test data
	setupTestData(t, store)

	is := is.New(t)

	// Create server
	server, err := NewServer(store, false)
	is.NoErr(err)

	// Test that server is created
	is.True(server != nil)
	is.True(server.mcpServer != nil)
	is.True(server.handler != nil)
}

func setupTestData(t *testing.T, store *core.FileTaskStore) {
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
	_, err = store.Update(completedTask, core.EditTaskParams{
		ID:        completedTask.ID.String(),
		NewStatus: &doneStatus,
	})
	is.NoErr(err)

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
		Parent:      &parentID,
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
	_, err = store.Update(inProgressTask, core.EditTaskParams{
		ID:        inProgressTask.ID.String(),
		NewStatus: &inProgressStatus,
	})
	is.NoErr(err)
}

func TestMCPHandlers(t *testing.T) {
	fs := afero.NewMemMapFs()
	store := core.NewFileTaskStore(fs, ".backlog")
	setupTestData(t, store)

	handler := &handler{store: store, mu: &sync.Mutex{}}
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	t.Run("handleTaskCreate", func(t *testing.T) {
		t.Run("successful_creation", func(t *testing.T) {
			is := is.New(t)
			params := core.CreateTaskParams{
				Title:       "Test Task",
				Description: "A test task description",
				Priority:    "high",
				Assigned:    []string{"testuser"},
				Labels:      []string{"test", "urgent"},
			}
			result, task, err := handler.create(ctx, req, params)
			is.NoErr(err)
			is.True(result == nil)
			is.Equal(task.Title, "Test Task")
			is.Equal(task.Priority.String(), "high")
		})
	})

	t.Run("handleTaskList", func(t *testing.T) {
		t.Run("list_all_tasks", func(t *testing.T) {
			is := is.New(t)
			params := core.ListTasksParams{}
			result, tasks, err := handler.list(ctx, req, params)
			is.NoErr(err)
			is.True(result == nil)
			is.True(len(tasks.Tasks) > 0)
		})

		t.Run("filter_by_status", func(t *testing.T) {
			is := is.New(t)

			params := core.ListTasksParams{
				Status: []string{"done"},
			}
			result, tasks, err := handler.list(ctx, req, params)
			is.NoErr(err)
			is.True(result == nil)
			for _, task := range tasks.Tasks {
				is.Equal(string(task.Status), "done")
			}
		})
	})

	t.Run("handleTaskView", func(t *testing.T) {
		t.Run("view_existing_task", func(t *testing.T) {
			is := is.New(t)

			// First create a task to view
			createParams := core.CreateTaskParams{
				Title: "View Test Task",
			}
			_, task, err := handler.create(ctx, req, createParams)
			is.NoErr(err)

			// Now view the task
			viewParams := ViewParams{
				ID: task.ID.String(),
			}
			result, task, err := handler.view(ctx, req, viewParams)
			is.NoErr(err)
			is.True(result == nil)
			is.Equal(task.ID, task.ID)
			is.Equal(task.Title, "View Test Task")
		})
	})

	t.Run("handleTaskSearch", func(t *testing.T) {
		t.Run("search_with_results", func(t *testing.T) {
			is := is.New(t)

			params := SearchParams{
				Query: "feature",
			}
			result, tasks, err := handler.search(ctx, req, params)
			is.NoErr(err)
			is.True(result == nil)
			is.True(len(tasks.Tasks) > 0)
		})

		t.Run("search_with_no_results", func(t *testing.T) {
			is := is.New(t)

			params := SearchParams{
				Query: "nonexistent",
			}
			result, tasks, err := handler.search(ctx, req, params)
			is.NoErr(err)
			is.True(result != nil)
			is.True(len(tasks.Tasks) == 0) // No results
		})
	})

	t.Run("handleTaskEdit", func(t *testing.T) {
		t.Run("edit_task_title", func(t *testing.T) {
			is := is.New(t)

			// First create a task to edit
			createParams := core.CreateTaskParams{
				Title: "Original Title",
			}
			_, task, err := handler.create(ctx, req, createParams)
			is.NoErr(err)

			// Now edit the task
			newTitle := "Updated Title"
			editParams := core.EditTaskParams{
				ID:       task.ID.String(),
				NewTitle: &newTitle,
			}
			_, task, err = handler.edit(ctx, req, editParams)
			is.NoErr(err)
			is.Equal(task.Title, "Updated Title")

			// Verify the task was updated
			viewParams := ViewParams{ID: task.ID.String()}
			_, task, err = handler.view(ctx, req, viewParams)
			is.NoErr(err)
			is.Equal(task.Title, "Updated Title")
		})
	})
}

package mcp

import (
	"context"
	"strings"
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
			result, _, err := handler.create(ctx, req, params)
			is.NoErr(err)
			is.True(result != nil) // must have result

			task, ok := result.StructuredContent.(core.Task)
			is.True(ok)
			is.Equal(task.Title, "Test Task")
			is.Equal(task.Priority.String(), "high")
		})
	})

	t.Run("handleTaskList", func(t *testing.T) {
		t.Run("list_all_tasks", func(t *testing.T) {
			is := is.New(t)
			result, _, err := handler.list(ctx, req, core.ListTasksParams{})
			is.NoErr(err)
			is.True(result != nil)
			listResult, ok := result.StructuredContent.(*core.ListResult)
			is.True(ok)
			is.Equal(len(listResult.Tasks), 9)
		})

		t.Run("filter_by_status", func(t *testing.T) {
			is := is.New(t)

			params := core.ListTasksParams{
				Status: []string{"done"},
			}
			result, _, err := handler.list(ctx, req, params)
			is.NoErr(err)
			is.True(result != nil)
			listResult, ok := result.StructuredContent.(*core.ListResult)
			is.True(ok)
			is.Equal(len(listResult.Tasks), 1)
			for _, task := range listResult.Tasks {
				is.Equal(string(task.Status), "done")
			}
		})
	})

	t.Run("handleTaskView", func(t *testing.T) {
		t.Run("view_existing_task", func(t *testing.T) {
			is := is.New(t)

			// First create a task to view
			createParams := core.CreateTaskParams{Title: "View Test Task"}
			createResult, _, err := handler.create(ctx, req, createParams)
			is.NoErr(err)
			is.True(createResult != nil)
			createdTask, ok := createResult.StructuredContent.(core.Task)
			is.True(ok)
			// Now view the task
			viewParams := ViewParams{ID: createdTask.ID.String()}
			viewResult, _, err := handler.view(ctx, req, viewParams)
			is.NoErr(err)
			// Compare created task with viewed task
			viewTask, ok := viewResult.StructuredContent.(core.Task)
			is.True(ok)
			// viewTask.Task is now a value, not a pointer, so no nil check needed
			is.Equal(viewTask.ID.String(), createdTask.ID.String())
			is.Equal(viewTask.Title, "View Test Task")
		})
	})

	t.Run("handleTaskEdit", func(t *testing.T) {
		t.Run("edit_task_title", func(t *testing.T) {
			is := is.New(t)

			// First create a task to edit
			createParams := core.CreateTaskParams{Title: "Original Title"}
			createResult, _, err := handler.create(ctx, req, createParams)
			is.NoErr(err)
			is.True(createResult != nil)

			// Parse the created task to get its ID

			createdTask, ok := createResult.StructuredContent.(core.Task)
			is.True(ok)
			// createdTask.Task is now a value, not a pointer, so no nil check needed

			// Now edit the task
			newTitle := "Updated Title"
			editParams := core.EditTaskParams{
				ID:       createdTask.ID.String(),
				NewTitle: &newTitle,
			}
			result, _, err := handler.edit(ctx, req, editParams)
			is.NoErr(err)
			is.True(result != nil)
			task, ok := result.StructuredContent.(core.Task)
			is.True(ok)
			// task.Task is now a value, not a pointer, so no nil check needed

			is.Equal(task.Title, "Updated Title")

			// Verify the task was updated
			viewParams := ViewParams{ID: task.ID.String()}
			result, _, err = handler.view(ctx, req, viewParams)
			is.NoErr(err)
			is.True(result != nil)
			task, ok = result.StructuredContent.(core.Task)
			is.True(ok)
			// task.Task is now a value, not a pointer, so no nil check needed

			is.Equal(task.Title, "Updated Title")
		})
	})

	t.Run("handleTaskArchive", func(t *testing.T) {
		t.Run("archive_existing_task", func(t *testing.T) {
			is := is.New(t)

			// First create a task to archive
			createParams := core.CreateTaskParams{
				Title:       "Task to Archive",
				Description: "This task will be archived",
				Priority:    "medium",
			}
			createResult, _, err := handler.create(ctx, req, createParams)
			is.NoErr(err)

			// Parse the created task to get its ID
			createdTask, ok := createResult.StructuredContent.(core.Task)
			is.True(ok)
			// createdTask.Task is now a value, not a pointer, so no nil check needed

			// Archive the task
			archiveParams := ArchiveParams{ID: createdTask.ID.String()}
			result, _, err := handler.archive(ctx, req, archiveParams)
			is.NoErr(err)
			is.True(result != nil)
			is.Equal(len(result.Content), 1)
			txtContent, ok := result.Content[0].(*mcp.TextContent)
			is.True(ok)
			is.True(strings.Contains(txtContent.Text, "archived successfully"))
		})

		t.Run("archive_nonexistent_task", func(t *testing.T) {
			is := is.New(t)

			archiveParams := ArchiveParams{
				ID: "nonexistent",
			}
			result, _, err := handler.archive(ctx, req, archiveParams)
			is.True(err != nil) // Should return an error
			is.True(result == nil)
		})
	})

	t.Run("handleTaskBatchCreate", func(t *testing.T) {
		t.Run("batch_create_multiple_tasks", func(t *testing.T) {
			is := is.New(t)

			params := ListCreateParams{
				Tasks: []core.CreateTaskParams{
					{
						Title:       "Batch Task 1",
						Description: "First batch task",
						Priority:    "high",
						Labels:      []string{"batch", "test"},
					},
					{
						Title:       "Batch Task 2",
						Description: "Second batch task",
						Priority:    "medium",
						Assigned:    []string{"user1"},
					},
					{
						Title:       "Batch Task 3",
						Description: "Third batch task",
						Priority:    "low",
					},
				},
			}

			result, _, err := handler.batchCreate(ctx, req, params)
			is.NoErr(err)
			is.True(result != nil)

			// Verify the response content
			st, ok := result.StructuredContent.(struct{ Tasks []core.Task })
			is.True(ok)
			is.Equal(len(st.Tasks), 3)

			// Verify task details
			is.Equal(st.Tasks[0].Title, "Batch Task 1")
			is.Equal(st.Tasks[1].Title, "Batch Task 2")
			is.Equal(st.Tasks[2].Title, "Batch Task 3")

			is.Equal(st.Tasks[0].Priority.String(), "high")
			is.Equal(st.Tasks[1].Priority.String(), "medium")
			is.Equal(st.Tasks[2].Priority.String(), "low")
		})

		t.Run("batch_create_empty_list", func(t *testing.T) {
			is := is.New(t)

			params := ListCreateParams{
				Tasks: []core.CreateTaskParams{},
			}

			result, _, err := handler.batchCreate(ctx, req, params)
			is.NoErr(err)
			is.True(result != nil)

			// Verify empty result
			st, ok := result.StructuredContent.(struct{ Tasks []core.Task })
			is.True(ok)
			is.Equal(len(st.Tasks), 0)
		})
	})

	t.Run("handleTaskCommit", func(t *testing.T) {
		t.Run("commit_disabled", func(t *testing.T) {
			is := is.New(t)

			// Test with the existing handler (autoCommit defaults to false)
			err := handler.commit("T1", "Test Task", "/some/path", "", "create")
			is.NoErr(err) // Should not error when autoCommit is false
		})

		// t.Run("commit_behavior_with_autocommit_enabled", func(t *testing.T) {
		// 	// Create a new handler instance with autoCommit enabled for this specific test
		// 	commitHandler := handler
		// 	commitHandler.autoCommit = true
		//
		// 	// Test commit method - expected to fail in test environment (no real git repo)
		// 	err := commitHandler.commit("T1", "Test Task", "/some/path", "", "create")
		// 	// In test environment with in-memory filesystem, this should fail
		// 	// but we're testing that the method doesn't panic and handles errors gracefully
		// 	// The exact error depends on whether we're in a git repo or not
		// 	_ = err // We expect this might error in test environment, which is fine
		// })
	})
}

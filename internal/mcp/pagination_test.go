package mcp

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/afero"
	"github.com/veggiemonk/backlog/internal/core"
)

func TestMCPPaginationHandlers(t *testing.T) {
	is := is.New(t)

	// Create a memory filesystem and setup MCP server
	fs := afero.NewMemMapFs()
	store := core.NewFileTaskStore(fs, "tasks")
	server, err := NewServer(store, false)
	is.NoErr(err)

	// Create test tasks
	testTasks := []core.CreateTaskParams{
		{Title: "Task A", Description: "First task"},
		{Title: "Task B", Description: "Second task"},
		{Title: "Task C", Description: "Third task"},
		{Title: "Task D", Description: "Fourth task"},
		{Title: "Task E", Description: "Fifth task"},
	}

	// Create the tasks
	for _, params := range testTasks {
		_, err := store.Create(params)
		is.NoErr(err)
	}

	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	t.Run("mcp_list_with_pagination", func(t *testing.T) {
		is := is.New(t)
		params := core.ListTasksParams{Limit: 3}

		result, _, err := server.handler.list(ctx, req, params)
		is.NoErr(err)
		is.True(result != nil)
		is.True(result.StructuredContent != nil)

		listResult, ok := result.StructuredContent.(core.ListResult)
		is.True(ok)
		is.Equal(len(listResult.Tasks), 3)

		// Should have pagination info
		is.True(listResult.Pagination != nil)
		is.Equal(listResult.Pagination.TotalResults, 5)
		is.Equal(listResult.Pagination.DisplayedResults, 3)
		is.Equal(listResult.Pagination.Offset, 0)
		is.Equal(listResult.Pagination.Limit, 3)
		is.True(listResult.Pagination.HasMore)
	})

	t.Run("mcp_list_with_offset_and_limit", func(t *testing.T) {
		is := is.New(t)
		params := core.ListTasksParams{Limit: 2, Offset: 2}

		result, _, err := server.handler.list(ctx, req, params)
		is.NoErr(err)
		is.True(result != nil)
		is.True(result.StructuredContent != nil)

		listResult, ok := result.StructuredContent.(core.ListResult)
		is.True(ok)
		is.Equal(len(listResult.Tasks), 2)

		// Should have pagination info
		is.True(listResult.Pagination != nil)
		is.Equal(listResult.Pagination.TotalResults, 5)
		is.Equal(listResult.Pagination.DisplayedResults, 2)
		is.Equal(listResult.Pagination.Offset, 2)
		is.Equal(listResult.Pagination.Limit, 2)
		is.True(listResult.Pagination.HasMore)
	})

	t.Run("mcp_list_without_pagination", func(t *testing.T) {
		is := is.New(t)

		// Test without pagination
		params := core.ListTasksParams{}

		result, _, err := server.handler.list(ctx, req, params)
		is.NoErr(err)
		is.True(result != nil)
		is.True(result.StructuredContent != nil)

		listResult, ok := result.StructuredContent.(core.ListResult)
		is.True(ok)
		is.Equal(len(listResult.Tasks), 5)

		// Should have pagination info, but with default values
		is.True(listResult.Pagination != nil)
		is.Equal(listResult.Pagination.TotalResults, 5)
		is.Equal(listResult.Pagination.DisplayedResults, 5)
		is.Equal(listResult.Pagination.Offset, 0)
		is.Equal(listResult.Pagination.Limit, 0)
		is.True(!listResult.Pagination.HasMore)
	})

	t.Run("mcp_empty_results_with_pagination", func(t *testing.T) {
		is := is.New(t)

		// Test empty results with pagination
		params := core.ListTasksParams{
			Limit:  5,
			Offset: 10, // Beyond available tasks
		}

		result, _, err := server.handler.list(ctx, req, params)
		is.NoErr(err)
		is.True(result != nil)

		// Should return Content instead of StructuredContent for empty results
		is.True(result.Content != nil)
		is.Equal(len(result.Content), 1)
		textContent, ok := result.Content[0].(*mcp.TextContent)
		is.True(ok)
		is.Equal(textContent.Text, "No tasks found.")
	})
}

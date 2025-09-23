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
		
		// Test with limit
		limit := 3
		params := core.ListTasksParams{
			Limit: &limit,
		}
		
		result, _, err := server.handler.list(ctx, req, params)
		is.NoErr(err)
		is.True(result != nil)
		is.True(result.StructuredContent != nil)
		
		listResult, ok := result.StructuredContent.(*core.ListResult)
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
		
		// Test with offset and limit
		limit := 2
		offset := 2
		params := core.ListTasksParams{
			Limit:  &limit,
			Offset: &offset,
		}
		
		result, _, err := server.handler.list(ctx, req, params)
		is.NoErr(err)
		is.True(result != nil)
		is.True(result.StructuredContent != nil)
		
		listResult, ok := result.StructuredContent.(*core.ListResult)
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
		
		listResult, ok := result.StructuredContent.(*core.ListResult)
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
	
	t.Run("mcp_search_with_pagination", func(t *testing.T) {
		is := is.New(t)
		
		// Test search with pagination
		limit := 1
		searchParams := SearchParams{
			Query: "task",
			Filters: &core.ListTasksParams{
				Limit: &limit,
			},
		}
		
		result, _, err := server.handler.search(ctx, req, searchParams)
		is.NoErr(err)
		is.True(result != nil)
		is.True(result.StructuredContent != nil)
		
		listResult, ok := result.StructuredContent.(*core.ListResult)
		is.True(ok)
		is.Equal(len(listResult.Tasks), 1)
		
		// Should have pagination info
		is.True(listResult.Pagination != nil)
		is.Equal(listResult.Pagination.TotalResults, 5)
		is.Equal(listResult.Pagination.DisplayedResults, 1)
		is.Equal(listResult.Pagination.Offset, 0)
		is.Equal(listResult.Pagination.Limit, 1)
		is.True(listResult.Pagination.HasMore)
	})
	
	t.Run("mcp_search_with_offset", func(t *testing.T) {
		is := is.New(t)
		
		// Test search with offset
		limit := 2
		offset := 3
		searchParams := SearchParams{
			Query: "task",
			Filters: &core.ListTasksParams{
				Limit:  &limit,
				Offset: &offset,
			},
		}
		
		result, _, err := server.handler.search(ctx, req, searchParams)
		is.NoErr(err)
		is.True(result != nil)
		is.True(result.StructuredContent != nil)
		
		listResult, ok := result.StructuredContent.(*core.ListResult)
		is.True(ok)
		is.Equal(len(listResult.Tasks), 2) // Should have 2 remaining (5 - 3 offset = 2)
		
		// Should have pagination info
		is.True(listResult.Pagination != nil)
		is.Equal(listResult.Pagination.TotalResults, 5)
		is.Equal(listResult.Pagination.DisplayedResults, 2)
		is.Equal(listResult.Pagination.Offset, 3)
		is.Equal(listResult.Pagination.Limit, 2)
		is.True(!listResult.Pagination.HasMore) // No more after showing 2 out of remaining 2
	})
	
	t.Run("mcp_search_without_pagination", func(t *testing.T) {
		is := is.New(t)
		
		// Test search without pagination
		searchParams := SearchParams{
			Query: "task",
		}
		
		result, _, err := server.handler.search(ctx, req, searchParams)
		is.NoErr(err)
		is.True(result != nil)
		is.True(result.StructuredContent != nil)
		
		listResult, ok := result.StructuredContent.(*core.ListResult)
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
		limit := 5
		offset := 10 // Beyond available tasks
		params := core.ListTasksParams{
			Limit:  &limit,
			Offset: &offset,
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

func TestMCPPaginationSchemaCompliance(t *testing.T) {
	is := is.New(t)
	
	// Create a memory filesystem and setup MCP server
	fs := afero.NewMemMapFs()
	store := core.NewFileTaskStore(fs, "tasks")
	server, err := NewServer(store, false)
	is.NoErr(err)
	
	// Create a test task
	_, err = store.Create(core.CreateTaskParams{Title: "Test Task"})
	is.NoErr(err)
	
	ctx := context.Background()
	req := &mcp.CallToolRequest{}
	
	t.Run("list_result_schema_compliance", func(t *testing.T) {
		is := is.New(t)
		
		// Test with pagination
		limit := 1
		params := core.ListTasksParams{
			Limit: &limit,
		}
		
		result, _, err := server.handler.list(ctx, req, params)
		is.NoErr(err)
		is.True(result != nil)
		is.True(result.StructuredContent != nil)
		
		// Validate the StructuredContent against the expected schema
		expectedSchema := listResultJSONSchema()
		validateStructuredContent(t, expectedSchema, result.StructuredContent)
	})
	
	t.Run("search_result_schema_compliance", func(t *testing.T) {
		is := is.New(t)
		
		// Test search with pagination
		limit := 1
		searchParams := SearchParams{
			Query: "test",
			Filters: &core.ListTasksParams{
				Limit: &limit,
			},
		}
		
		result, _, err := server.handler.search(ctx, req, searchParams)
		is.NoErr(err)
		is.True(result != nil)
		is.True(result.StructuredContent != nil)
		
		// Validate the StructuredContent against the expected schema
		expectedSchema := listResultJSONSchema()
		validateStructuredContent(t, expectedSchema, result.StructuredContent)
	})
}

func TestPaginationEdgeCases(t *testing.T) {
	is := is.New(t)
	
	// Create a memory filesystem and setup MCP server
	fs := afero.NewMemMapFs()
	store := core.NewFileTaskStore(fs, "tasks")
	server, err := NewServer(store, false)
	is.NoErr(err)
	
	ctx := context.Background()
	req := &mcp.CallToolRequest{}
	
	t.Run("large_offset_returns_empty", func(t *testing.T) {
		is := is.New(t)
		
		// Create one task
		_, err := store.Create(core.CreateTaskParams{Title: "Single Task"})
		is.NoErr(err)
		
		// Request with large offset
		limit := 10
		offset := 100
		params := core.ListTasksParams{
			Limit:  &limit,
			Offset: &offset,
		}
		
		result, _, err := server.handler.list(ctx, req, params)
		is.NoErr(err)
		is.True(result != nil)
		is.True(result.Content != nil) // Should return empty content, not structured content
	})
	
	t.Run("zero_limit_returns_all", func(t *testing.T) {
		is := is.New(t)
		
		// Test zero limit (should return all results)
		limit := 0
		params := core.ListTasksParams{
			Limit: &limit,
		}
		
		result, _, err := server.handler.list(ctx, req, params)
		is.NoErr(err)
		is.True(result != nil)
		is.True(result.StructuredContent != nil)
		
		listResult, ok := result.StructuredContent.(*core.ListResult)
		is.True(ok)
		is.Equal(len(listResult.Tasks), 1) // Should return the one task we created
		
		// Should have pagination info
		is.True(listResult.Pagination != nil)
		is.Equal(listResult.Pagination.Limit, 0)
		is.True(!listResult.Pagination.HasMore)
	})
}

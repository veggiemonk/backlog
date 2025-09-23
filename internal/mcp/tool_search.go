package mcp

import (
	"context"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
)

func (s *Server) registerTaskSearch() error {
	inputSchema, err := jsonschema.For[SearchParams](nil)
	if err != nil {
		return err
	}
	tool := &mcp.Tool{
		Name:         "task_search",
		Title:        "Search by content",
		Description:  "Search tasks by content with optional pagination. Returns a list of matching tasks with optional pagination metadata. Use 'limit' and 'offset' in filters for pagination.",
		InputSchema:  inputSchema,
		OutputSchema: listResultJSONSchema(),
	}
	mcp.AddTool(s.mcpServer, tool, s.handler.search)
	return nil
}

type SearchParams struct {
	Query   string                `json:"query" jsonschema:"Required. The search query."`
	Filters *core.ListTasksParams `json:"filters" jsonschema:"Optional. Additional filters for the search."`
}

func (h *handler) search(ctx context.Context, req *mcp.CallToolRequest, params SearchParams) (*mcp.CallToolResult, any, error) {
	var filters core.ListTasksParams
	if params.Filters != nil {
		filters = *params.Filters
	}

	// Get total search count without pagination for metadata
	totalFilters := filters
	totalFilters.Limit = nil
	totalFilters.Offset = nil
	allTasks, err := h.store.Search(params.Query, totalFilters)
	if err != nil {
		return nil, nil, fmt.Errorf("search (total count): %v", err)
	}
	totalCount := len(allTasks)

	// Get paginated search results
	tasks, err := h.store.Search(params.Query, filters)
	if err != nil {
		return nil, nil, fmt.Errorf("search: %v", err)
	}

	// Create result with pagination info
	result := &core.ListResult{
		Tasks: tasks,
	}

	// Add pagination info if pagination was requested
	if filters.Limit != nil || filters.Offset != nil {
		offsetVal := 0
		if filters.Offset != nil {
			offsetVal = *filters.Offset
		}
		limitVal := 0
		if filters.Limit != nil {
			limitVal = *filters.Limit
		}
		hasMore := (offsetVal + len(tasks)) < totalCount

		result.Pagination = &core.PaginationInfo{
			TotalResults:     totalCount,
			DisplayedResults: len(tasks),
			Offset:           offsetVal,
			Limit:            limitVal,
			HasMore:          hasMore,
		}
	}

	if len(tasks) == 0 {
		content := []mcp.Content{&mcp.TextContent{Text: "No matching tasks found."}}
		return &mcp.CallToolResult{Content: content}, result, nil
	}

	res := &mcp.CallToolResult{StructuredContent: result}
	return res, nil, nil
}

package mcp

import (
	"context"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
)

func (s *Server) registerTaskList() error {
	inputSchema, err := jsonschema.For[core.ListTasksParams](nil)
	if err != nil {
		return err
	}
	description := `List tasks, with optional filtering, sorting, and pagination. 
	Returns a list of tasks with optional pagination metadata.
	Use 'limit' and 'offset' parameters for pagination.
`
	tool := &mcp.Tool{
		Name:        "task_list",
		Title:       "List tasks",
		Description: description,
		InputSchema: inputSchema,
		OutputSchema: listResultJSONSchema(),
	}
	mcp.AddTool(s.mcpServer, tool, s.handler.list)
	return nil
}

func (h *handler) list(ctx context.Context, req *mcp.CallToolRequest, params core.ListTasksParams) (*mcp.CallToolResult, any, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Get total count without pagination for metadata
	totalParams := params
	totalParams.Limit = nil
	totalParams.Offset = nil
	allTasks, err := h.store.List(totalParams)
	if err != nil {
		return nil, nil, fmt.Errorf("list (total count): %v", err)
	}
	totalCount := len(allTasks)

	// Get paginated results
	tasks, err := h.store.List(params)
	if err != nil {
		return nil, nil, fmt.Errorf("list: %v", err)
	}
	
	// Create result with pagination info
	result := &core.ListResult{
		Tasks: tasks,
	}
	
	// Add pagination info if pagination was requested
	if params.Limit != nil || params.Offset != nil {
		offsetVal := 0
		if params.Offset != nil {
			offsetVal = *params.Offset
		}
		limitVal := 0
		if params.Limit != nil {
			limitVal = *params.Limit
		}
		hasMore := (offsetVal + len(tasks)) < totalCount
		
		result.Pagination = &core.PaginationInfo{
			TotalResults:    totalCount,
			DisplayedResults: len(tasks),
			Offset:          offsetVal,
			Limit:           limitVal,
			HasMore:         hasMore,
		}
	}

	if len(tasks) == 0 {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "No tasks found."}}}, result, nil
	}

	res := &mcp.CallToolResult{StructuredContent: result}
	return res, nil, nil
}

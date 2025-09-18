package mcp

import (
	"context"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
)

func (s *Server) registerTaskSearch() error {
	inputSchema, err := jsonschema.For[SearchParams](nil)
	if err != nil {
		return err
	}
	outputSchema, err := jsonschema.For[TaskListResponse](nil)
	if err != nil {
		return err
	}
	tool := &mcp.Tool{
		Name:         "task_search",
		Title:        "Search by content",
		Description:  "Search tasks by content. Returns a list of matching tasks.",
		InputSchema:  inputSchema,
		OutputSchema: outputSchema,
	}
	mcp.AddTool(s.mcpServer, tool, s.handler.search)

	return nil
}

type SearchParams struct {
	Query   string                `json:"query" jsonschema:"Required. The search query."`
	Filters *core.ListTasksParams `json:"filters" jsonschema:"Optional. Additional filters for the search."`
}

func (h *handler) search(ctx context.Context, req *mcp.CallToolRequest, params SearchParams) (*mcp.CallToolResult, TaskListResponse, error) {
	var filters core.ListTasksParams
	if params.Filters != nil {
		filters = *params.Filters
	}
	tasks, err := h.store.Search(params.Query, filters)
	if err != nil {
		return nil, TaskListResponse{}, err
	}
	if len(tasks) == 0 {
		content := []mcp.Content{&mcp.TextContent{Text: "No matching tasks found."}}
		return &mcp.CallToolResult{Content: content}, TaskListResponse{}, nil
	}

	results := TaskListResponse{Tasks: make([]core.Task, 0, len(tasks))}
	for _, t := range tasks {
		if t != nil {
			results.Tasks = append(results.Tasks, *t)
		}
	}
	return nil, results, nil
}

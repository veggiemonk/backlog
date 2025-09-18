package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
)

func (s *Server) registerTaskSearch() error {
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "task_search",
		Description: "Search tasks by content. Returns a list of matching tasks.",
	}, s.handler.search)

	return nil
}

type SearchParams struct {
	Query   string               `json:"query" jsonschema:"Required. The search query."`
	Filters core.ListTasksParams `json:"filters" jsonschema:"Optional. Additional filters for the search."`
}

func (h *handler) search(ctx context.Context, req *mcp.CallToolRequest, params SearchParams) (*mcp.CallToolResult, TaskListResponse, error) {
	tasks, err := h.store.Search(params.Query, params.Filters)
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

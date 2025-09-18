package mcp

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
)

func (s *Server) registerTaskSearch() error {
	tool := &mcp.Tool{
		Name:        "task_search",
		Title:       "Search by content",
		Description: "Search tasks by content. Returns a list of matching tasks.",
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
	tasks, err := h.store.Search(params.Query, filters)
	if err != nil {
		return nil, nil, err
	}
	if len(tasks) == 0 {
		content := []mcp.Content{&mcp.TextContent{Text: "No matching tasks found."}}
		return &mcp.CallToolResult{Content: content}, nil, nil
	}

	wrappedTask := struct{ Tasks []*core.Task }{Tasks: tasks}
	b, err := json.Marshal(wrappedTask)
	if err != nil {
		return nil, nil, err
	}
	res := &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
	}
	return res, nil, nil
}

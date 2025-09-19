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
		Name:        "task_search",
		Title:       "Search by content",
		Description: "Search tasks by content. Returns a list of matching tasks.",
		InputSchema: inputSchema,
		OutputSchema: &jsonschema.Schema{Type: "object", Properties: map[string]*jsonschema.Schema{
			"tasks": {Type: "array", Items: taskJSONSchema()},
		}},
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
		return nil, nil, fmt.Errorf("search: %v", err)
	}
	if len(tasks) == 0 {
		content := []mcp.Content{&mcp.TextContent{Text: "No matching tasks found."}}
		return &mcp.CallToolResult{Content: content}, nil, nil
	}
	// Needs to be object, cannot be array
	wrappedTask := struct{ Tasks []*core.Task }{Tasks: tasks}
	res := &mcp.CallToolResult{StructuredContent: wrappedTask}
	return res, nil, nil
}

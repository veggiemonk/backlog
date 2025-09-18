package mcp

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
)

func (s *Server) registerTaskList() error {
	description := `List tasks, with optional filtering and sorting. 
	Returns a list of tasks.
`
	tool := &mcp.Tool{
		Name:        "task_list",
		Title:       "List tasks",
		Description: description,
	}
	mcp.AddTool(s.mcpServer, tool, s.handler.list)
	return nil
}

func (h *handler) list(ctx context.Context, req *mcp.CallToolRequest, params core.ListTasksParams) (*mcp.CallToolResult, any, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	tasks, err := h.store.List(params)
	if err != nil {
		return nil, nil, err
	}
	if len(tasks) == 0 {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "No tasks found."}}}, nil, nil
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

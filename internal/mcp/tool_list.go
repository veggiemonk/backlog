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
	description := `List tasks, with optional filtering and sorting. 
	Returns a list of tasks.
`
	tool := &mcp.Tool{
		Name:        "task_list",
		Title:       "List tasks",
		Description: description,
		InputSchema: inputSchema,
		OutputSchema: wrappedTasksJSONSchema(),
	}
	mcp.AddTool(s.mcpServer, tool, s.handler.list)
	return nil
}

func (h *handler) list(ctx context.Context, req *mcp.CallToolRequest, params core.ListTasksParams) (*mcp.CallToolResult, any, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	tasks, err := h.store.List(params)
	if err != nil {
		return nil, nil, fmt.Errorf("list: %v", err)
	}
	if len(tasks) == 0 {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "No tasks found."}}}, nil, nil
	}
	// Needs to be object, cannot be array
	wrapped := struct{ Tasks []*core.Task }{Tasks: tasks}
	res := &mcp.CallToolResult{StructuredContent: wrapped}
	return res, nil, nil
}

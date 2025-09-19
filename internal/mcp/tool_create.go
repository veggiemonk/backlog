package mcp

import (
	"context"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
)

func (s *Server) registerTaskCreate() error {
	inputSchema, err := jsonschema.For[core.CreateTaskParams](nil)
	if err != nil {
		return err
	}
	description := `Create a new task. 
The task ID is automatically generated. 
Returns the created task.
`
	tool := &mcp.Tool{
		Name:         "task_create",
		Description:  description,
		InputSchema:  inputSchema,
		OutputSchema: taskJSONSchema(),
	}
	mcp.AddTool(s.mcpServer, tool, s.handler.create)
	return nil
}

func (h *handler) create(
	ctx context.Context,
	req *mcp.CallToolRequest,
	params core.CreateTaskParams,
) (*mcp.CallToolResult, any, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	task, err := h.store.Create(params)
	if err != nil {
		return nil, nil, fmt.Errorf("create: %v", err)
	}
	path := h.store.Path(task)
	if err := h.commit(task.ID.Name(), task.Title, path, "", "create"); err != nil {
		// Log the error but do not fail the creation
		logging.Warn("auto-commit failed for task creation", "task_id", task.ID, "error", err)
	}
	wrapped := struct{ Task *core.Task }{Task: task}
	res := &mcp.CallToolResult{StructuredContent: wrapped}
	return res, nil, nil
}

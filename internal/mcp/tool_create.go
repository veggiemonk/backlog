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
	taskCreateParams, err := jsonschema.For[core.CreateTaskParams](nil)
	if err != nil {
		return fmt.Errorf("jsonschema.For[core.CreateTaskParams]: %v", err)
	}
	// taskCreateParams.Properties["id"].Type = "string"
	// taskCreateParams.Properties["new_parent"].Type = "string"

	taskSchema, err := jsonschema.For[core.Task](nil)
	if err != nil {
		return fmt.Errorf("jsonschema.For[core.Task]: %v", err)
	}
	// taskSchema.Properties["id"].Type = "string"
	// taskSchema.Properties["parent"].Type = "string"

	description := `Create a new task. 
The task ID is automatically generated. 
Returns the created task.",
`
	tool := &mcp.Tool{
		Name:         "task_create",
		Description:  description,
		InputSchema:  taskCreateParams,
		OutputSchema: taskSchema,
	}
	mcp.AddTool(s.mcpServer, tool, s.handler.create)
	return nil
}

func (h *handler) create(
	ctx context.Context,
	req *mcp.CallToolRequest,
	params core.CreateTaskParams,
) (*mcp.CallToolResult, core.Task, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	task, err := h.store.Create(params)
	if err != nil {
		return nil, core.Task{}, err
	}
	path := h.store.Path(task)
	if err := h.commit(task.ID.Name(), task.Title, path, "", "create"); err != nil {
		// Log the error but do not fail the creation
		logging.Warn("auto-commit failed for task creation", "task_id", task.ID, "error", err)
	}
	return nil, *task, nil
}

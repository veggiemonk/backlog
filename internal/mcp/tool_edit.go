package mcp

import (
	"context"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
)

func (s *Server) registerTaskEdit() error {
	inputSchema, err := jsonschema.For[core.EditTaskParams](nil)
	if err != nil {
		return err
	}
	description := `Edit an existing task by its ID.
This is a partial update, only the provided fields will be changed. 
Returns the updated task.`

	editTool := &mcp.Tool{
		Name:         "task_edit",
		Description:  description,
		InputSchema:  inputSchema,
		OutputSchema: taskJSONSchema(),
	}

	mcp.AddTool(s.mcpServer, editTool, s.handler.edit)
	return nil
}

func (h *handler) edit(
	ctx context.Context,
	req *mcp.CallToolRequest,
	params core.EditTaskParams,
) (*mcp.CallToolResult, any, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	task, err := h.store.Get(params.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("edit: %v", err)
	}
	updatedTask, err := h.store.Update(task, params)
	if err != nil {
		return nil, nil, fmt.Errorf("edit: %v", err)
	}
	err = h.commit(
		updatedTask.ID.Name(),
		updatedTask.Title,
		h.store.Path(updatedTask),
		h.store.Path(task),
		"edit",
	)
	if err != nil {
		// Log the error but do not fail the edit
		logging.Warn("auto-commit failed for task edit", "task_id", task.ID, "error", err)
	}

	wrapped := struct{ Task *core.Task }{Task: task}
	res := &mcp.CallToolResult{StructuredContent: wrapped}
	return res, nil, nil
}

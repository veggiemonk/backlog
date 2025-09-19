package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
)

func (s *Server) registerTaskEdit() error {
	// inputSchema, err := jsonschema.For[core.EditTaskParams](nil)
	// if err != nil {
	// 	return fmt.Errorf("jsonschema.For[core.EditTaskParams]: %v", err)
	// }
	// outputSchema, err := jsonschema.For[core.Task](nil)
	// if err != nil {
	// 	return fmt.Errorf("jsonschema.For[core.Task]: %v", err)
	// }
	description := `Edit an existing task by its ID.
This is a partial update, only the provided fields will be changed. 
Returns the updated task.`

	editTool := &mcp.Tool{
		Name:        "task_edit",
		Description: description,
		// InputSchema:  inputSchema,
		// OutputSchema: outputSchema,
	}

	mcp.AddTool(s.mcpServer, editTool, s.handler.edit)
	return nil
}

func (h *handler) edit(
	ctx context.Context,
	req *mcp.CallToolRequest,
	params core.EditTaskParams,
) (*mcp.CallToolResult, any, error) {
	operation := "task_edit"

	// Validate input parameters
	if validationErr := h.validator.ValidateEditTaskParams(params); validationErr != nil {
		return h.responder.WrapValidationError(validationErr, operation)
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	task, err := h.store.Get(params.ID)
	if err != nil {
		// Wrap the error with proper categorization
		mcpErr := WrapError(err, operation)
		return h.responder.WrapError(mcpErr)
	}

	updatedTask, err := h.store.Update(task, params)
	if err != nil {
		// Wrap the error with proper categorization
		mcpErr := WrapError(err, operation)
		return h.responder.WrapError(mcpErr)
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

	return h.responder.WrapSuccess(updatedTask, operation)
}

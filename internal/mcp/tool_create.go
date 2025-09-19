package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
)

func (s *Server) registerTaskCreate() error {
	// inputSchema, err := jsonschema.For[core.CreateTaskParams](nil)
	// if err != nil {
	// 	return fmt.Errorf("jsonschema.For[core.CreateTaskParams]: %v", err)
	// }
	// outputSchema, err := jsonschema.For[core.Task](nil)
	// if err != nil {
	// 	return fmt.Errorf("jsonschema.For[core.Task]: %v", err)
	// }

	description := `Create a new task. 
The task ID is automatically generated. 
Returns the created task.",
`
	tool := &mcp.Tool{
		Name:        "task_create",
		Description: description,
		// InputSchema:  inputSchema,
		// OutputSchema: outputSchema,
	}
	mcp.AddTool(s.mcpServer, tool, s.handler.create)
	return nil
}

func (h *handler) create(
	ctx context.Context,
	req *mcp.CallToolRequest,
	params core.CreateTaskParams,
) (*mcp.CallToolResult, any, error) {
	operation := "task_create"

	// Validate input parameters
	if validationErr := h.validator.ValidateCreateTaskParams(params); validationErr != nil {
		return h.responder.WrapValidationError(validationErr, operation)
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	task, err := h.store.Create(params)
	if err != nil {
		// Wrap the error with proper categorization
		mcpErr := WrapError(err, operation)
		return h.responder.WrapError(mcpErr)
	}

	path := h.store.Path(task)
	if err := h.commit(task.ID.Name(), task.Title, path, "", "create"); err != nil {
		// Log the error but do not fail the creation
		logging.Warn("auto-commit failed for task creation", "task_id", task.ID, "error", err)
	}

	return h.responder.WrapSuccess(task, operation)
}

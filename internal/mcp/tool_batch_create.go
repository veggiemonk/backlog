package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
)

func (s *Server) registerTaskBatchCreate() error {
	description := `Create a list of new tasks.
The schema is a list of "task_create" input parameters.
The task ID of each task is automatically generated. Returns the list of created task.
`
	tool := &mcp.Tool{
		Name:        "task_batch_create",
		Description: description,
	}
	mcp.AddTool(s.mcpServer, tool, s.handler.batchCreate)
	return nil
}

type ListCreateParams struct {
	Tasks []core.CreateTaskParams `json:"new_tasks"`
}

func (h *handler) batchCreate(ctx context.Context, req *mcp.CallToolRequest, listParams ListCreateParams) (*mcp.CallToolResult, any, error) {
	operation := "task_batch_create"

	// Validate batch size
	if len(listParams.Tasks) == 0 {
		validationErr := NewValidationError("new_tasks", "At least one task is required for batch creation", len(listParams.Tasks))
		return h.responder.WrapValidationError(validationErr, operation)
	}

	if len(listParams.Tasks) > 50 {
		validationErr := &MCPError{
			Code:     ErrorCodeInvalidInput,
			Message:  "Too many tasks in batch. Maximum allowed is 50",
			Category: CategoryValidation,
			Details: &ErrorDetails{
				Field: "new_tasks",
				Value: len(listParams.Tasks),
				Expected: "array length <= 50",
				Constraints: map[string]interface{}{
					"max_tasks": 50,
					"provided_tasks": len(listParams.Tasks),
				},
			},
		}
		return h.responder.WrapValidationError(validationErr, operation)
	}

	// Validate each task in the batch
	for i, params := range listParams.Tasks {
		if validationErr := h.validator.ValidateCreateTaskParams(params); validationErr != nil {
			// Add context about which task in the batch failed
			validationErr.Details.Context = map[string]interface{}{
				"batch_index": i,
				"batch_total": len(listParams.Tasks),
			}
			validationErr.Message = fmt.Sprintf("Task #%d in batch: %s", i+1, validationErr.Message)
			return h.responder.WrapValidationError(validationErr, operation)
		}
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	tasks := make([]*core.Task, 0, len(listParams.Tasks))
	var creationErrors []string

	for i, params := range listParams.Tasks {
		task, err := h.store.Create(params)
		if err != nil {
			// Collect error but continue with other tasks
			creationErrors = append(creationErrors, fmt.Sprintf("Task #%d (%s): %v", i+1, params.Title, err))
			continue
		}

		tasks = append(tasks, task)
		path := h.store.Path(task)
		if err := h.commit(task.ID.Name(), task.Title, path, "", "create"); err != nil {
			logging.Warn("auto-commit failed for task creation", "task_id", task.ID, "error", err)
		}
	}

	// If some tasks failed to create, return partial success with error details
	if len(creationErrors) > 0 {
		mcpErr := &MCPError{
			Code:     ErrorCodeInvalidOperation,
			Message:  fmt.Sprintf("Batch creation partially failed. %d/%d tasks created successfully", len(tasks), len(listParams.Tasks)),
			Category: CategoryBusiness,
			Details: &ErrorDetails{
				Context: map[string]interface{}{
					"successful_tasks": len(tasks),
					"failed_tasks": len(creationErrors),
					"total_tasks": len(listParams.Tasks),
					"errors": creationErrors,
				},
			},
			Operation: operation,
		}
		return h.responder.WrapError(mcpErr)
	}

	// All tasks created successfully
	response := struct{ Tasks []*core.Task }{Tasks: tasks}
	return h.responder.WrapSuccess(response, operation)
}

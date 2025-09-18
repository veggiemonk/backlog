package mcp

import (
	"context"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
)

func (s *Server) registerEditTask() error {
	taskSchema, err := jsonschema.For[core.Task](nil)
	if err != nil {
		return fmt.Errorf("jsonschema.For[core.Task]: %v", err)
	}
	// taskSchema.Properties["id"].Type = "string"
	// taskSchema.Properties["parent"].Type = "string"

	taskEditSchema, err := jsonschema.For[core.EditTaskParams](nil)
	if err != nil {
		return fmt.Errorf("jsonschema.For[core.EditTaskParams]: %v", err)
	}
	// taskEditSchema.Properties["id"].Type = "string"
	// taskEditSchema.Properties["new_parent"].Type = "string"

	description := `Edit an existing task by its ID.
This is a partial update, only the provided fields will be changed. 
Returns the updated task.`

	editTool := &mcp.Tool{
		Name:         "task_edit",
		Description:  description,
		InputSchema:  taskEditSchema,
		OutputSchema: taskSchema,
	}

	mcp.AddTool(s.mcpServer, editTool, s.handler.edit)
	return nil
}

func (h *handler) edit(ctx context.Context, req *mcp.CallToolRequest, params core.EditTaskParams) (*mcp.CallToolResult, core.Task, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	task, err := h.store.Get(params.ID)
	if err != nil {
		return nil, core.Task{}, err
	}
	// Get a summary of the changes by capturing the last history entry
	// historyBefore := len(task.History)
	updatedTask, err := h.store.Update(task, params)
	if err != nil {
		return nil, core.Task{}, err
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
	// historyAfter := len(task.History)
	// var changes string
	// if historyAfter > historyBefore {
	// 	// Get the last `historyAfter - historyBefore` changes
	// 	for i := historyBefore; i < historyAfter; i++ {
	// 		changes += fmt.Sprintf("- %s\n", task.History[i].Change)
	// 	}
	// } else {
	// 	changes = "No changes were made."
	// }
	// summary := fmt.Sprintf("Task %s updated successfully:\n%s", task.ID, changes)
	// content := []mcp.Content{&mcp.TextContent{Text: summary}}
	return nil, *updatedTask, nil
}

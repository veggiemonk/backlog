package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/logging"
)

func (s *Server) registerTaskArchive() error {
	// inputSchema, err := jsonschema.For[ArchiveParams](nil)
	// if err != nil {
	// 	return err
	// }
	tool := &mcp.Tool{
		Name:  "task_archive",
		Title: "Archive a task",
		Description: `Archive a task.
The task will be moving to the archived directory and setting status to archived. 
Returns the archived task.`,
		// InputSchema: inputSchema,
	}
	mcp.AddTool(s.mcpServer, tool, s.handler.archive)
	return nil
}

type ArchiveParams struct {
	ID string `json:"id" jsonschema:"Required. The ID of the task to archive."`
}

func (h *handler) archive(ctx context.Context, req *mcp.CallToolRequest, params ArchiveParams) (*mcp.CallToolResult, any, error) {
	operation := "task_archive"

	// Validate task ID
	if validationErr := h.validator.ValidateTaskID(params.ID, "id"); validationErr != nil {
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

	oldPath := h.store.Path(task)
	archivedPath, err := h.store.Archive(task.ID)
	if err != nil {
		// Wrap the error with proper categorization
		mcpErr := WrapError(err, operation)
		return h.responder.WrapError(mcpErr)
	}

	if err := h.commit(task.ID.Name(), task.Title, archivedPath, oldPath, "archive"); err != nil {
		// Log the error but do not fail the archive
		logging.Warn("auto-commit failed for task archive", "task_id", task.ID, "error", err)
	}

	// Create structured response for archive result
	response := struct {
		TaskID       string `json:"task_id"`
		Title        string `json:"title"`
		Status       string `json:"status"`
		ArchivedPath string `json:"archived_path"`
		Message      string `json:"message"`
	}{
		TaskID:       task.ID.Name(),
		Title:        task.Title,
		Status:       "archived",
		ArchivedPath: archivedPath,
		Message:      fmt.Sprintf("Task %s archived successfully", task.ID.Name()),
	}

	return h.responder.WrapSuccess(response, operation)
}

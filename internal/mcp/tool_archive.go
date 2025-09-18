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
	h.mu.Lock()
	defer h.mu.Unlock()
	task, err := h.store.Get(params.ID)
	if err != nil {
		return nil, nil, err
	}
	oldPath := h.store.Path(task)
	archivedPath, err := h.store.Archive(task.ID)
	if err != nil {
		return nil, nil, err
	}
	if err := h.commit(task.ID.Name(), task.Title, archivedPath, oldPath, "archive"); err != nil {
		// Log the error but do not fail the archive
		logging.Warn("auto-commit failed for task archive", "task_id", task.ID, "error", err)
	}
	summary := fmt.Sprintf("Task %s archived successfully:\n\n", task.ID.Name())
	summary += fmt.Sprintf("- Title: %s\n", task.Title)
	summary += fmt.Sprintf("- Status: %s\n", task.Status)
	summary += "- The task has been moved to the archived directory\n"
	content := &mcp.TextContent{Text: summary}
	return &mcp.CallToolResult{Content: []mcp.Content{content}}, nil, nil
}

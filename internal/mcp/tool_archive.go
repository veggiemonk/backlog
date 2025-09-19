package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/logging"
)

func (s *Server) registerTaskArchive() error {
	inputSchema, err := jsonschema.For[ArchiveParams](nil)
	if err != nil {
		return err
	}
	description := `Archive a task.
The task will be moving to the archived directory and setting status to archived. 
Returns the archived task.`

	tool := &mcp.Tool{
		Name:         "task_archive",
		Title:        "Archive a task",
		Description:  description,
		InputSchema:  inputSchema,
		OutputSchema: taskJSONSchema(),
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
	// Fetch current task, archive it, then re-fetch to return archived state
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
	archivedTask, err := h.store.Get(params.ID)
	if err != nil {
		return nil, nil, err
	}
	summary := fmt.Sprintf("Task %s archived successfully:\n\n", task.ID.Name())
	summary += fmt.Sprintf("- Title: %s\n", task.Title)
	summary += fmt.Sprintf("- Status: %s\n", archivedTask.Status)
	summary += "- The task has been moved to the archived directory\n"
	content := &mcp.TextContent{Text: summary}
	// Also return JSON of archived task as first content item for consistency
	jb, _ := json.Marshal(archivedTask)
	jsonContent := &mcp.TextContent{Text: string(jb)}
	return &mcp.CallToolResult{Content: []mcp.Content{jsonContent, content}}, archivedTask, nil
}

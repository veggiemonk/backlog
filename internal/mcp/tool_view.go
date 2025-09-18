package mcp

import (
	"context"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
)

func (s *Server) registerTaskView() error {
	taskSchema, err := jsonschema.For[core.Task](nil)
	if err != nil {
		return fmt.Errorf("jsonschema.For[core.Task]: %v", err)
	}
	taskSchema.Properties["id"].Type = "string"
	taskSchema.Properties["parent"].Type = "string"

	tool := &mcp.Tool{
		Name:         "task_view",
		Description:  "View a single task by its ID. Returns the task.",
		OutputSchema: taskSchema,
	}

	mcp.AddTool(s.mcpServer, tool, s.handler.view)
	return nil
}

type ViewParams struct {
	ID string `json:"id" jsonschema:"Required. The ID of the task."`
}

func (h *handler) view(ctx context.Context, req *mcp.CallToolRequest, params ViewParams) (*mcp.CallToolResult, core.Task, error) {
	task, err := h.store.Get(params.ID)
	if err != nil {
		return nil, core.Task{}, err
	}
	return nil, *task, err
}

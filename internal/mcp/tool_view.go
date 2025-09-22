package mcp

import (
	"context"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *Server) registerTaskView() error {
	inputSchema, err := jsonschema.For[ViewParams](nil)
	if err != nil {
		return err
	}

	tool := &mcp.Tool{
		Name:         "task_view",
		Title:        "View a task",
		Description:  "View a single task by its ID. Returns the task.",
		InputSchema:  inputSchema,
		OutputSchema: taskJSONSchema(),
	}

	mcp.AddTool(s.mcpServer, tool, s.handler.view)
	return nil
}

type ViewParams struct {
	ID string `json:"id" jsonschema:"Required. The ID of the task."`
}

func (h *handler) view(ctx context.Context, req *mcp.CallToolRequest, params ViewParams) (*mcp.CallToolResult, any, error) {
	task, err := h.store.Get(params.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("view: %v", err)
	}
	// Needs to be object wrapped in struct as expected by wrappedTaskJSONSchema
	res := &mcp.CallToolResult{StructuredContent: task}
	return res, nil, nil
}

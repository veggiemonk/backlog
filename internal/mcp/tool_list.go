package mcp

import (
	"context"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
)

func (s *Server) registerTaskList() error {
	inputSchema, err := jsonschema.For[core.ListTasksParams](nil)
	if err != nil {
		return err
	}
	description := `List tasks, with optional filtering, sorting, and pagination. 
	Returns a list of tasks with optional pagination metadata.
	Use 'limit' and 'offset' parameters for pagination.
`
	tool := &mcp.Tool{
		Name:         "task_list",
		Title:        "List tasks",
		Description:  description,
		InputSchema:  inputSchema,
		OutputSchema: listResultJSONSchema(),
	}
	mcp.AddTool(s.mcpServer, tool, s.handler.list)
	return nil
}

func (h *handler) list(ctx context.Context, req *mcp.CallToolRequest, params core.ListTasksParams) (*mcp.CallToolResult, any, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	listResult, err := h.store.List(params)
	if err != nil {
		return nil, nil, fmt.Errorf("list: %v", err)
	}

	if len(listResult.Tasks) == 0 {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "No tasks found."}}}, listResult, nil
	}

	res := &mcp.CallToolResult{StructuredContent: listResult}
	return res, nil, nil
}

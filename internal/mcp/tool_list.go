package mcp

import (
	"context"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
)

// TaskListResponse is used to wrap the list of tasks in a JSON object
// to conform to the MCP specification for structuredContent.
type TaskListResponse struct {
	Tasks []core.Task `json:"tasks"`
}

func (s *Server) registerTaskList() error {
	inputSchema, err := jsonschema.For[core.ListTasksParams](nil)
	if err != nil {
		return err
	}
	outputSchema, err := jsonschema.For[TaskListResponse](nil)
	if err != nil {
		return err
	}

	description := `List tasks, with optional filtering and sorting. 
	Returns a list of tasks.
`
	tool := &mcp.Tool{
		Name:         "task_list",
		Title:        "List tasks",
		Description:  description,
		InputSchema:  inputSchema,
		OutputSchema: outputSchema,
	}

	mcp.AddTool(s.mcpServer, tool, s.handler.list)
	return nil
}

func (h *handler) list(ctx context.Context, req *mcp.CallToolRequest, params core.ListTasksParams) (*mcp.CallToolResult, TaskListResponse, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	tasks, err := h.store.List(params)
	if err != nil {
		return nil, TaskListResponse{}, err
	}
	if len(tasks) == 0 {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "No tasks found."}}}, TaskListResponse{}, nil
	}
	results := TaskListResponse{Tasks: make([]core.Task, 0, len(tasks))}
	for _, t := range tasks {
		if t != nil {
			results.Tasks = append(results.Tasks, *t)
		}
	}
	return nil, results, nil
}

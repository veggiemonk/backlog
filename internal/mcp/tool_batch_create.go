package mcp

import (
	"context"
	"encoding/json"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
)

func (s *Server) registerTaskBatchCreate() error {
	inputSchema, err := jsonschema.For[ListCreateParams](nil)
	if err != nil {
		return err
	}
	description := `Create a list of new tasks.
The schema is a list of "task_create" input parameters.
The task ID of each task is automatically generated. Returns the list of created task.
`
	// Output schema: { tasks: Task[] }
	tool := &mcp.Tool{
		Name:        "task_batch_create",
		Description: description,
		InputSchema: inputSchema,
		OutputSchema: &jsonschema.Schema{
			Type: "object", Properties: map[string]*jsonschema.Schema{
				"tasks": {Type: "array", Items: taskJSONSchema()},
			},
		},
	}
	mcp.AddTool(s.mcpServer, tool, s.handler.batchCreate)
	return nil
}

type ListCreateParams struct {
	Tasks []core.CreateTaskParams `json:"new_tasks"`
}

func (h *handler) batchCreate(ctx context.Context, req *mcp.CallToolRequest, listParams ListCreateParams) (*mcp.CallToolResult, any, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	tasks := make([]*core.Task, 0, len(listParams.Tasks))
	for _, params := range listParams.Tasks {
		task, err := h.store.Create(params)
		if err != nil {
			return nil, nil, err
		}
		tasks = append(tasks, task)
		path := h.store.Path(task)
		if err := h.commit(task.ID.Name(), task.Title, path, "", "create"); err != nil {
			logging.Warn("auto-commit failed for task creation", "task_id", task.ID, "error", err)
		}
	}
	wrappedTask := struct{ Tasks []*core.Task }{Tasks: tasks}
	b, err := json.Marshal(wrappedTask)
	if err != nil {
		return nil, nil, err
	}
	res := &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}
	return res, wrappedTask, nil
}

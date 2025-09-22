package mcp

import (
	"context"
	"fmt"
	"reflect"

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
	outputSchema, err := jsonschema.For[[]core.Task](&jsonschema.ForOptions{
		TypeSchemas: map[reflect.Type]*jsonschema.Schema{
			reflect.TypeFor[core.TaskID](): {OneOf: []*jsonschema.Schema{
				{Type: "string"}, {Type: "null"},
			}},
		},
	})
	if err != nil {
		return err
	}
	tool := &mcp.Tool{
		Name:         "task_batch_create",
		Description:  description,
		InputSchema:  inputSchema,
		OutputSchema: outputSchema,
		// OutputSchema: wrappedTasksJSONSchema(),
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
	tasks := make([]core.Task, 0, len(listParams.Tasks))
	for _, params := range listParams.Tasks {
		task, err := h.store.Create(params)
		if err != nil {
			return nil, nil, fmt.Errorf("batch_create: %v", err)
		}
		tasks = append(tasks, task)
		path := h.store.Path(task)
		if err := h.commit(task.ID.Name(), task.Title, path, "", "create"); err != nil {
			logging.Warn("auto-commit failed for task creation", "task_id", task.ID, "error", err)
		}
	}
	// Needs to be object, cannot be array
	res := &mcp.CallToolResult{StructuredContent: tasks}
	return res, nil, nil
}

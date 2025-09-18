package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
)

// // TaskCreateArgs represents arguments for task_create tool
//
//	type TaskCreateArgs struct {
//		Title       string   `json:"title"`
//		Description string   `json:"description,omitempty"`
//		Labels      []string `json:"labels,omitempty"`
//	}
//
// // TaskListArgs represents arguments for task_list tool
//
//	type TaskListArgs struct {
//		Status     string   `json:"status,omitempty"`
//		Sort       []string `json:"sort,omitempty"`
//		Reverse    bool     `json:"reverse,omitempty"`
//		Unassigned bool     `json:"unassigned,omitempty"`
//		Parent     string   `json:"parent,omitempty"`
//		Labels     string   `json:"labels,omitempty"`
//	}
func (s *Server) registerTaskBatchCreate() error {
	// TODO:
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name: "task_batch_create",
		Description: `Create a list of new tasks.
The schema is a list of "task_create" input parameters.
The task ID of each task is automatically generated. Returns the list of created task.
`,
	}, s.handler.batchCreate)
	return nil
}

type ListCreateParams struct {
	Tasks []core.CreateTaskParams `json:"new_tasks"`
}

func (h *handler) batchCreate(ctx context.Context, req *mcp.CallToolRequest, listParams ListCreateParams) (*mcp.CallToolResult, TaskListResponse, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	tasks := make([]*core.Task, 0, len(listParams.Tasks))
	for _, params := range listParams.Tasks {
		task, err := h.store.Create(params)
		if err != nil {
			return nil, TaskListResponse{}, err
		}
		tasks = append(tasks, task)
		path := h.store.Path(task)
		if err := h.commit(task.ID.Name(), task.Title, path, "", "create"); err != nil {
			logging.Warn("auto-commit failed for task creation", "task_id", task.ID, "error", err)
		}
	}
	results := TaskListResponse{Tasks: make([]core.Task, 0, len(tasks))}
	for _, t := range tasks {
		if t != nil {
			results.Tasks = append(results.Tasks, *t)
		}
	}
	return nil, results, nil
}

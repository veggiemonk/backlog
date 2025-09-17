package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/commit"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
)

// addTools adds all MCP tools to the server
func (s *Server) addTools() {
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "task_create",
		Description: "Create a new task. The task ID is automatically generated. Returns the created task.",
	}, s.handler.create)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name: "task_batch_create",
		Description: `Create a list of new tasks.
The schema is a list of "task_create" input parameters.
The task ID of each task is automatically generated. Returns the list of created task.
`,
	}, s.handler.batchCreate)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "task_list",
		Description: "List tasks, with optional filtering and sorting. Returns a list of tasks.",
	}, s.handler.list)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "task_view",
		Description: "View a single task by its ID. Returns the task.",
	}, s.handler.view)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "task_edit",
		Description: "Edit an existing task by its ID. This is a partial update, only the provided fields will be changed. Returns the updated task.",
	}, s.handler.edit)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "task_search",
		Description: "Search tasks by content. Returns a list of matching tasks.",
	}, s.handler.search)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "task_archive",
		Description: "Archive a task by moving it to the archived directory and setting status to archived. Returns the archived task.",
	}, s.handler.archive)
}

func (h *handler) commit(id, title, path, oldPath, msg string) error {
	if h.autoCommit {
		commitMsg := fmt.Sprintf("feat(task): %s %s - \"%s\"", msg, id, title)
		if err := commit.Add(path, oldPath, commitMsg); err != nil {
			return fmt.Errorf("auto-commit failed: %w", err)
		}
	}
	return nil
}

// Tool handler implementations

func (h *handler) create(ctx context.Context, req *mcp.CallToolRequest, params core.CreateTaskParams) (*mcp.CallToolResult, *core.Task, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	task, err := h.store.Create(params)
	if err != nil {
		return nil, nil, err
	}
	path := h.store.Path(task)
	if err := h.commit(task.ID.Name(), task.Title, path, "", "create"); err != nil {
		// Log the error but do not fail the creation
		logging.Warn("auto-commit failed for task creation", "task_id", task.ID, "error", err)
	}

	summary := fmt.Sprintf("Task %s created successfully:\n\n", task.ID.Name())
	summary += fmt.Sprintf("- Title: %s\n", task.Title)
	summary += fmt.Sprintf("- Status: %s\n", task.Status)
	if len(task.Assigned) > 0 {
		summary += fmt.Sprintf("- Assigned: %s\n", strings.Join(task.Assigned, ", "))
	}
	if len(task.Labels) > 0 {
		summary += fmt.Sprintf("- Labels: %s\n", strings.Join(task.Labels, ", "))
	}

	content := &mcp.TextContent{Text: summary}
	return &mcp.CallToolResult{Content: []mcp.Content{content}}, task, nil
}

type (
	ListCreateParams struct {
		List []core.CreateTaskParams `json:"list"`
	}
)

func (h *handler) batchCreate(ctx context.Context, req *mcp.CallToolRequest, listParams ListCreateParams) (*mcp.CallToolResult, TaskListResponse, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	tasks := make([]*core.Task, 0, len(listParams.List))
	summaries := make([]string, 0, len(listParams.List))
	for _, params := range listParams.List {
		task, err := h.store.Create(params)
		if err != nil {
			return nil, TaskListResponse{}, err
		}
		tasks = append(tasks, task)
		path := h.store.Path(task)
		if err := h.commit(task.ID.Name(), task.Title, path, "", "create"); err != nil {
			// Log the error but do not fail the creation
			logging.Warn("auto-commit failed for task creation", "task_id", task.ID, "error", err)
		}

		summary := fmt.Sprintf("Task %s created successfully:\n\n", task.ID.Name())
		summary += fmt.Sprintf("- Title: %s\n", task.Title)
		summary += fmt.Sprintf("- Status: %s\n", task.Status)
		if len(task.Assigned) > 0 {
			summary += fmt.Sprintf("- Assigned: %s\n", strings.Join(task.Assigned, ", "))
		}
		if len(task.Labels) > 0 {
			summary += fmt.Sprintf("- Labels: %s\n", strings.Join(task.Labels, ", "))
		}
		summaries = append(summaries, summary)
	}

	content := &mcp.TextContent{Text: strings.Join(summaries, "|")}
	results := TaskListResponse{Tasks: tasks}
	return &mcp.CallToolResult{Content: []mcp.Content{content}}, results, nil
}

func (h *handler) edit(ctx context.Context, req *mcp.CallToolRequest, params core.EditTaskParams) (*mcp.CallToolResult, *core.Task, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	task, err := h.store.Get(params.ID)
	if err != nil {
		return nil, nil, err
	}
	// Get a summary of the changes by capturing the last history entry
	historyBefore := len(task.History)
	updatedTask, err := h.store.Update(task, params)
	if err != nil {
		return nil, nil, err
	}
	if err := h.commit(
		updatedTask.ID.Name(),
		updatedTask.Title,
		h.store.Path(updatedTask),
		h.store.Path(task),
		"edit",
	); err != nil {
		// Log the error but do not fail the edit
		logging.Warn("auto-commit failed for task edit", "task_id", task.ID, "error", err)
	}
	historyAfter := len(task.History)

	var changes string
	if historyAfter > historyBefore {
		// Get the last `historyAfter - historyBefore` changes
		for i := historyBefore; i < historyAfter; i++ {
			changes += fmt.Sprintf("- %s\n", task.History[i].Change)
		}
	} else {
		changes = "No changes were made."
	}

	summary := fmt.Sprintf("Task %s updated successfully:\n%s", task.ID, changes)

	content := []mcp.Content{&mcp.TextContent{Text: summary}}
	return &mcp.CallToolResult{Content: content}, nil, nil
}

func (h *handler) list(ctx context.Context, req *mcp.CallToolRequest, params core.ListTasksParams) (*mcp.CallToolResult, any, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	tasks, err := h.store.List(params)
	if err != nil {
		return nil, nil, err
	}
	if len(tasks) == 0 {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "No tasks found."}}}, nil, nil
	}

	var table string
	table += "| ID | Title | Status |\n"
	table += "|---|---|---|\n"
	for _, task := range tasks {
		table += fmt.Sprintf("| %s | %s | %s |\n", task.ID.Name(), task.Title, task.Status)
	}

	content := &mcp.TextContent{Text: table}
	return &mcp.CallToolResult{Content: []mcp.Content{content}}, TaskListResponse{Tasks: tasks}, nil
}

type viewParams struct {
	ID string `json:"id" jsonschema:"Required. The ID of the task."`
}

func (h *handler) view(ctx context.Context, req *mcp.CallToolRequest, params viewParams) (*mcp.CallToolResult, any, error) {
	task, err := h.store.Get(params.ID)
	if err != nil {
		return nil, nil, err
	}

	summary := fmt.Sprintf("# Task %s: %s\n\n%s", task.ID.Name(), task.Title, string(task.Bytes()))

	// Return the task struct in the second return value. The SDK will handle serialization.
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: summary}}}, task, nil
}

type SearchParams struct {
	Query   string               `json:"query" jsonschema:"Required. The search query."`
	Filters core.ListTasksParams `json:"filters" jsonschema:"Optional. Additional filters for the search."`
}

func (h *handler) search(ctx context.Context, req *mcp.CallToolRequest, params SearchParams) (*mcp.CallToolResult, any, error) {
	tasks, err := h.store.Search(params.Query, params.Filters)
	if err != nil {
		return nil, nil, err
	}
	if len(tasks) == 0 {
		content := []mcp.Content{&mcp.TextContent{Text: "No matching tasks found."}}
		return &mcp.CallToolResult{Content: content}, nil, nil
	}

	var table string
	table += "| ID | Title | Status |\n"
	table += "|---|---|---|\n"
	for _, task := range tasks {
		table += fmt.Sprintf("| %s | %s | %s |\n", task.ID.Name(), task.Title, task.Status)
	}

	content := []mcp.Content{&mcp.TextContent{Text: table}}
	return &mcp.CallToolResult{Content: content}, TaskListResponse{Tasks: tasks}, nil
}

type archiveParams struct {
	ID string `json:"id" jsonschema:"Required. The ID of the task to archive."`
}

func (h *handler) archive(ctx context.Context, req *mcp.CallToolRequest, params archiveParams) (*mcp.CallToolResult, any, error) {
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
	return &mcp.CallToolResult{Content: []mcp.Content{content}}, task, nil
}

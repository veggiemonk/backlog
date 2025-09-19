package mcp

import (
	"context"
	"encoding/json"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
)

// taskJSONSchema returns an explicit JSON schema for core.Task matching its JSON encoding.
func taskJSONSchema() *jsonschema.Schema {
	mkStringOrStringArray := func() *jsonschema.Schema {
		return &jsonschema.Schema{OneOf: []*jsonschema.Schema{
			{Type: "string"},
			{Type: "array", Items: &jsonschema.Schema{Type: "string"}},
		}}
	}
	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"id":     {Type: "string"},
			"title":  {Type: "string"},
			"status": {Type: "string", Enum: []any{string(core.StatusTodo), string(core.StatusInProgress), string(core.StatusDone), string(core.StatusCancelled), string(core.StatusArchived), string(core.StatusRejected)}},
			"parent": {Type: "string"},
			"assigned": {OneOf: []*jsonschema.Schema{
				mkStringOrStringArray(), {Type: "null"},
			}},
			"labels": {OneOf: []*jsonschema.Schema{
				mkStringOrStringArray(), {Type: "null"},
			}},
			"dependencies": {OneOf: []*jsonschema.Schema{
				mkStringOrStringArray(), {Type: "null"},
			}},
			"priority":   {Type: "string", Enum: []any{"unknown", "low", "medium", "high", "critical"}},
			"created_at": {Type: "string"},
			"updated_at": {OneOf: []*jsonschema.Schema{{Type: "string"}, {Type: "null"}}},
			"history": {OneOf: []*jsonschema.Schema{
				{Type: "array", Items: &jsonschema.Schema{Type: "object", Properties: map[string]*jsonschema.Schema{"timestamp": {Type: "string"}, "change": {Type: "string"}}}},
				{Type: "null"},
			}},
			"description": {Type: "string"},
			"acceptance_criteria": {OneOf: []*jsonschema.Schema{
				{Type: "array", Items: &jsonschema.Schema{Type: "object", Properties: map[string]*jsonschema.Schema{"text": {Type: "string"}, "checked": {Type: "boolean"}, "index": {Type: "integer"}}}},
				{Type: "null"},
			}},
			"implementation_plan":  {OneOf: []*jsonschema.Schema{{Type: "string"}, {Type: "null"}}},
			"implementation_notes": {OneOf: []*jsonschema.Schema{{Type: "string"}, {Type: "null"}}},
		},
	}
}

func (s *Server) registerTaskCreate() error {
	inputSchema, err := jsonschema.For[core.CreateTaskParams](nil)
	if err != nil {
		return err
	}
	description := `Create a new task. 
The task ID is automatically generated. 
Returns the created task.
`
	tool := &mcp.Tool{
		Name:         "task_create",
		Description:  description,
		InputSchema:  inputSchema,
		OutputSchema: taskJSONSchema(),
	}
	mcp.AddTool(s.mcpServer, tool, s.handler.create)
	return nil
}

func (h *handler) create(
	ctx context.Context,
	req *mcp.CallToolRequest,
	params core.CreateTaskParams,
) (*mcp.CallToolResult, any, error) {
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

	// Provide both structured and text content for compatibility
	b, err := json.Marshal(task)
	if err != nil {
		return nil, nil, err
	}
	res := &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}
	return res, task, nil
}

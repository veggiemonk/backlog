package mcp

import (
	"context"
	"encoding/json"

	// "github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *Server) registerTaskView() error {
	// inputSchema, err := jsonschema.For[ViewParams](nil)
	// if err != nil {
	// 	return fmt.Errorf("jsonschema.For[core.Task]: %v", err)
	// }
	// taskIDSchema, err := jsonschema.ForType(reflect.TypeOf(core.TaskID{}), nil)
	// if err != nil {
	// 	return fmt.Errorf("jsonschema.For[core.Task]: %v", err)
	// }
	// taskIDSchema.Type = "string"
	//
	// outputSchema, err := jsonschema.For[core.Task](&jsonschema.ForOptions{
	// 	TypeSchemas: map[any]*jsonschema.Schema{
	// 		reflect.TypeOf(core.TaskID{}): taskIDSchema,
	// 	},
	// })
	// if err != nil {
	// 	return err
	// }
	// outputSchema.Properties["id"].Type = "string"
	// outputSchema.Properties["parent"].Type = "string"

	tool := &mcp.Tool{
		Name:        "task_view",
		Title:       "View a task",
		Description: "View a single task by its ID. Returns the task.",
		// InputSchema:  inputSchema,
		// OutputSchema: outputSchema,
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
		return nil, nil, err
	}
	b, err := json.Marshal(task)
	if err != nil {
		return nil, nil, err
	}
	res := &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
	}
	return res, nil, err
}

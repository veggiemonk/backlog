package mcp

import (
	"context"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
)

func (s *Server) registerTaskSearch() error {
	inputSchema, err := jsonschema.For[SearchParams](nil)
	if err != nil {
		return err
	}
	tool := &mcp.Tool{
		Name:         "task_search",
		Title:        "Search by content",
		Description:  "Search tasks by content with optional pagination. Returns a list of matching tasks with optional pagination metadata. Use 'limit' and 'offset' in filters for pagination.",
		InputSchema:  inputSchema,
		OutputSchema: listResultJSONSchema(),
	}
	mcp.AddTool(s.mcpServer, tool, s.handler.search)
	return nil
}

type SearchParams struct {
	Query   string                `json:"query" jsonschema:"Required. The search query."`
	Filters *core.ListTasksParams `json:"filters" jsonschema:"Optional. Additional filters for the search."`
}

func (h *handler) search(ctx context.Context, req *mcp.CallToolRequest, params SearchParams) (*mcp.CallToolResult, any, error) {
	var filters core.ListTasksParams
	if params.Filters != nil {
		filters = *params.Filters
	}

	listResult, err := h.store.Search(params.Query, filters)
	if err != nil {
		return nil, nil, fmt.Errorf("search: %v", err)
	}

	if len(listResult.Tasks) == 0 {
		content := []mcp.Content{&mcp.TextContent{Text: "No matching tasks found."}}
		return &mcp.CallToolResult{Content: content}, listResult, nil
	}

	res := &mcp.CallToolResult{StructuredContent: listResult}
	return res, nil, nil
}

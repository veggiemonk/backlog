package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
)

func (s *Server) registerTaskSearch() error {
	tool := &mcp.Tool{
		Name:        "task_search",
		Title:       "Search by content",
		Description: "Search tasks by content. Returns a list of matching tasks.",
	}
	mcp.AddTool(s.mcpServer, tool, s.handler.search)
	return nil
}

type SearchParams struct {
	Query   string                `json:"query" jsonschema:"Required. The search query."`
	Filters *core.ListTasksParams `json:"filters" jsonschema:"Optional. Additional filters for the search."`
}

func (h *handler) search(ctx context.Context, req *mcp.CallToolRequest, params SearchParams) (*mcp.CallToolResult, any, error) {
	operation := "task_search"

	// Validate search query
	if params.Query == "" {
		validationErr := NewMissingRequiredError("query")
		return h.responder.WrapValidationError(validationErr, operation)
	}

	// Validate search query length
	if len(params.Query) > 1000 {
		validationErr := &MCPError{
			Code:     ErrorCodeInvalidInput,
			Message:  "Search query too long. Maximum length is 1000 characters",
			Category: CategoryValidation,
			Details: &ErrorDetails{
				Field: "query",
				Value: len(params.Query),
				Expected: "string length <= 1000",
				Constraints: map[string]interface{}{
					"max_length": 1000,
					"actual_length": len(params.Query),
				},
			},
		}
		return h.responder.WrapValidationError(validationErr, operation)
	}

	// Validate filters if provided
	var filters core.ListTasksParams
	if params.Filters != nil {
		if validationErr := h.validator.ValidateListParams(*params.Filters); validationErr != nil {
			return h.responder.WrapValidationError(validationErr, operation)
		}
		filters = *params.Filters
	}

	tasks, err := h.store.Search(params.Query, filters)
	if err != nil {
		// Wrap the error with proper categorization
		mcpErr := WrapError(err, operation)
		return h.responder.WrapError(mcpErr)
	}

	// Return structured response for search results
	response := struct{ Tasks []*core.Task }{Tasks: tasks}
	return h.responder.WrapSuccess(response, operation)
}

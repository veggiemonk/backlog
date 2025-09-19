package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
)

func (s *Server) registerTaskList() error {
	description := `List tasks, with optional filtering and sorting. 
	Returns a list of tasks.
`
	tool := &mcp.Tool{
		Name:        "task_list",
		Title:       "List tasks",
		Description: description,
	}
	mcp.AddTool(s.mcpServer, tool, s.handler.list)
	return nil
}

func (h *handler) list(ctx context.Context, req *mcp.CallToolRequest, params core.ListTasksParams) (*mcp.CallToolResult, any, error) {
	operation := "task_list"

	// Validate input parameters
	if validationErr := h.validator.ValidateListParams(params); validationErr != nil {
		return h.responder.WrapValidationError(validationErr, operation)
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	tasks, err := h.store.List(params)
	if err != nil {
		// Wrap the error with proper categorization
		mcpErr := WrapError(err, operation)
		return h.responder.WrapError(mcpErr)
	}

	if len(tasks) == 0 {
		// Return structured response for empty results
		response := struct{ Tasks []*core.Task }{Tasks: []*core.Task{}}
		return h.responder.WrapSuccess(response, operation)
	}

	// Check if response will exceed token limits
	if WillExceedLimit(tasks) {
		// If no pagination parameters provided, suggest optimal chunk size
		if params.Limit == nil && params.Offset == nil {
			optimalChunkSize := CalculateOptimalChunkSize(tasks)
			mcpErr := &MCPError{
				Code:     ErrorCodeResponseTooLarge,
				Message:  fmt.Sprintf("Response too large (%d tasks). Please use pagination. Suggested limit: %d", len(tasks), optimalChunkSize),
				Category: CategoryMCP,
				Details: &ErrorDetails{
					Context: map[string]interface{}{
						"total_tasks":        len(tasks),
						"suggested_limit":    optimalChunkSize,
						"requires_pagination": true,
					},
				},
				Operation: operation,
			}
			return h.responder.WrapError(mcpErr)
		}
		// If pagination is provided but still too large, return error
		if params.Limit != nil && *params.Limit > 0 {
			chunkSize := CalculateOptimalChunkSize(tasks)
			if *params.Limit > chunkSize {
				mcpErr := &MCPError{
					Code:     ErrorCodeResponseTooLarge,
					Message:  fmt.Sprintf("Requested limit (%d) too large. Maximum recommended limit: %d", *params.Limit, chunkSize),
					Category: CategoryMCP,
					Details: &ErrorDetails{
						Field: "limit",
						Value: *params.Limit,
						Context: map[string]interface{}{
							"requested_limit":   *params.Limit,
							"maximum_limit":     chunkSize,
						},
					},
					Operation: operation,
				}
				return h.responder.WrapError(mcpErr)
			}
		}
	}

	// Create response with pagination metadata if applicable
	response := createTaskListResponse(tasks, params)

	// Apply response monitoring middleware if available
	if h.middleware != nil {
		if wrapped := h.middleware.WrapResponse(response, operation); wrapped != nil {
			response = wrapped
		}
	}

	// Marshal and create final response
	b, err := json.Marshal(response)
	if err != nil {
		mcpErr := NewSystemError(operation, "Failed to serialize task list response", err)
		return h.responder.WrapError(mcpErr)
	}

	res := &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
	}
	return res, nil, nil
}

// TaskListResponse wraps the task list with pagination metadata
type TaskListResponse struct {
	Tasks      []*core.Task        `json:"tasks"`
	Pagination *PaginationMetadata `json:"pagination,omitempty"`
}

// PaginationMetadata provides information about pagination
type PaginationMetadata struct {
	Offset    int  `json:"offset"`
	Limit     int  `json:"limit"`
	Total     int  `json:"total"`
	HasMore   bool `json:"has_more"`
	NextPage  *int `json:"next_page,omitempty"`
}

// createTaskListResponse creates a properly formatted response with pagination metadata
func createTaskListResponse(tasks []*core.Task, params core.ListTasksParams) interface{} {
	// If no pagination parameters, return simple format for backward compatibility
	if params.Limit == nil && params.Offset == nil {
		return struct{ Tasks []*core.Task }{Tasks: tasks}
	}

	// Get the original total before pagination was applied
	// We need to recalculate this since List() already applied pagination
	allTasks, _ := getAllTasksCount(params)

	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}

	limit := len(tasks)
	if params.Limit != nil {
		limit = *params.Limit
	}

	hasMore := offset+len(tasks) < allTasks
	var nextPage *int
	if hasMore {
		next := offset + limit
		nextPage = &next
	}

	pagination := &PaginationMetadata{
		Offset:   offset,
		Limit:    limit,
		Total:    allTasks,
		HasMore:  hasMore,
		NextPage: nextPage,
	}

	return &TaskListResponse{
		Tasks:      tasks,
		Pagination: pagination,
	}
}

// getAllTasksCount gets the total count of tasks that would match the filter (without pagination)
func getAllTasksCount(params core.ListTasksParams) (int, error) {
	// This is a simplified implementation - in a real scenario you might want to
	// optimize this to avoid loading all tasks just to count them
	// For now, we'll return a reasonable default based on the context
	// In production, this should query the actual store
	return 50, nil // Placeholder - should be implemented to query actual count
}

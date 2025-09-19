package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	startTime := time.Now()

	// Validate input parameters
	if validationErr := h.validator.ValidateListParams(params); validationErr != nil {
		return h.responder.WrapValidationError(validationErr, operation)
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Get pagination configuration from handler
	paginationConfig := h.paginationConfig

	// Get total count before applying pagination
	allTasks, err := h.store.List(core.ListTasksParams{
		Parent:        params.Parent,
		Status:        params.Status,
		Assigned:      params.Assigned,
		Labels:        params.Labels,
		Priority:      params.Priority,
		Unassigned:    params.Unassigned,
		DependedOn:    params.DependedOn,
		HasDependency: params.HasDependency,
		Sort:          params.Sort,
		Reverse:       params.Reverse,
		// No pagination parameters for total count
	})
	if err != nil {
		mcpErr := WrapError(err, operation)
		return h.responder.WrapError(mcpErr)
	}

	totalCount := len(allTasks)

	// Check if auto-pagination should be triggered
	if paginationConfig.ShouldAutoPaginate(totalCount) && params.Limit == nil && params.Offset == nil {
		// Estimate average task size
		estimatedSize := EstimateResponseSize(allTasks)
		avgTaskSize := estimatedSize / max(totalCount, 1)

		optimalPageSize := paginationConfig.CalculateOptimalPageSize(totalCount, avgTaskSize)

		mcpErr := &MCPError{
			Code:     ErrorCodeResponseTooLarge,
			Message:  fmt.Sprintf("Large dataset detected (%d tasks). Auto-pagination recommended. Suggested limit: %d", totalCount, optimalPageSize),
			Category: CategoryMCP,
			Details: &ErrorDetails{
				Context: map[string]interface{}{
					"total_tasks":           totalCount,
					"suggested_limit":       optimalPageSize,
					"auto_pagination_triggered": true,
					"estimated_response_size":   estimatedSize,
				},
			},
			Operation: operation,
		}
		return h.responder.WrapError(mcpErr)
	}

	// Apply pagination to parameters if provided
	if params.Limit != nil {
		validatedLimit := paginationConfig.ValidatePageSize(*params.Limit)
		params.Limit = &validatedLimit
	}

	// Get the actual paginated tasks
	tasks, err := h.store.List(params)
	if err != nil {
		mcpErr := WrapError(err, operation)
		return h.responder.WrapError(mcpErr)
	}

	// Check if response will still exceed token limits even with pagination
	if len(tasks) > 0 && WillExceedLimit(tasks) {
		estimatedSize := EstimateResponseSize(tasks)
		avgTaskSize := estimatedSize / len(tasks)
		optimalPageSize := paginationConfig.CalculateOptimalPageSize(totalCount, avgTaskSize)

		mcpErr := &MCPError{
			Code:     ErrorCodeResponseTooLarge,
			Message:  fmt.Sprintf("Response still too large (%d tasks, ~%d tokens). Reduce page size to %d or less", len(tasks), estimatedSize, optimalPageSize),
			Category: CategoryMCP,
			Details: &ErrorDetails{
				Field: "limit",
				Value: len(tasks),
				Context: map[string]interface{}{
					"current_limit":         len(tasks),
					"recommended_limit":     optimalPageSize,
					"estimated_tokens":      estimatedSize,
					"token_limit":          paginationConfig.ResponseSizeConfig.TokenLimit,
				},
			},
			Operation: operation,
		}
		return h.responder.WrapError(mcpErr)
	}

	generationTime := time.Since(startTime)

	// Create enhanced response with advanced pagination metadata
	response := createEnhancedTaskListResponse(tasks, params, totalCount, paginationConfig, generationTime)

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


// createEnhancedTaskListResponse creates a response with advanced pagination metadata
func createEnhancedTaskListResponse(
	tasks []*core.Task,
	params core.ListTasksParams,
	totalCount int,
	config PaginationConfig,
	generationTime time.Duration,
) interface{} {
	// If no pagination parameters, return simple format for backward compatibility
	if params.Limit == nil && params.Offset == nil {
		return struct{ Tasks []*core.Task }{Tasks: tasks}
	}

	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}

	limit := len(tasks)
	if params.Limit != nil {
		limit = *params.Limit
	}

	// Estimate response size for metadata
	estimatedSize := 0
	if len(tasks) > 0 {
		estimatedSize = EstimateResponseSize(tasks)
	}

	// Create advanced pagination metadata
	pagination := CreateAdvancedPaginationMetadata(
		offset,
		limit,
		totalCount,
		len(tasks),
		config,
		estimatedSize,
		generationTime,
	)

	return &EnhancedTaskListResponse{
		Tasks:      tasks,
		Pagination: pagination,
	}
}

// EnhancedTaskListResponse wraps the task list with advanced pagination metadata
type EnhancedTaskListResponse struct {
	Tasks      []*core.Task                `json:"tasks"`
	Pagination *AdvancedPaginationMetadata `json:"pagination,omitempty"`
}

// Helper function for max since Go doesn't have a built-in max for ints
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

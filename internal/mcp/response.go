package mcp

import (
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ResponseWrapper handles response formatting and error handling for MCP tools
type ResponseWrapper struct {
	middleware *ResponseSizeMiddleware
}

// NewResponseWrapper creates a new response wrapper
func NewResponseWrapper(middleware *ResponseSizeMiddleware) *ResponseWrapper {
	return &ResponseWrapper{
		middleware: middleware,
	}
}

// WrapSuccess creates a successful MCP response with proper monitoring
func (w *ResponseWrapper) WrapSuccess(data interface{}, operation string) (*mcp.CallToolResult, any, error) {
	// Marshal data to JSON
	b, err := json.Marshal(data)
	if err != nil {
		mcpErr := NewSystemError(operation, "Failed to serialize response", err)
		return w.WrapError(mcpErr)
	}

	// Create response
	res := &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
	}

	// Monitor response size
	if w.middleware != nil {
		if wrapped := w.middleware.WrapResponse(res, operation); wrapped != nil {
			if wrappedRes, ok := wrapped.(*mcp.CallToolResult); ok {
				res = wrappedRes
			}
		}
	}

	return res, nil, nil
}

// WrapError creates an error response with structured error information
func (w *ResponseWrapper) WrapError(mcpErr *MCPError) (*mcp.CallToolResult, any, error) {
	// For validation and not found errors, return structured error in response
	// For system errors, return as Go error to trigger MCP error response
	switch mcpErr.Category {
	case CategoryValidation, CategoryNotFound, CategoryBusiness, CategoryMCP:
		// Return structured error in the response content
		errorJSON := mcpErr.ToJSON()
		res := &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: errorJSON}},
			IsError: true,
		}
		return res, nil, nil
	case CategorySystem:
		// Return as Go error for system-level failures
		return nil, nil, mcpErr
	default:
		// Default to Go error
		return nil, nil, mcpErr
	}
}

// WrapValidationError creates a validation error response
func (w *ResponseWrapper) WrapValidationError(err *MCPError, operation string) (*mcp.CallToolResult, any, error) {
	if err != nil {
		err.Operation = operation
	}
	return w.WrapError(err)
}

// WrapSystemError creates a system error response
func (w *ResponseWrapper) WrapSystemError(err error, operation, message string) (*mcp.CallToolResult, any, error) {
	mcpErr := NewSystemError(operation, message, err)
	return w.WrapError(mcpErr)
}

// WrapTaskNotFoundError creates a task not found error response
func (w *ResponseWrapper) WrapTaskNotFoundError(taskID, operation string) (*mcp.CallToolResult, any, error) {
	mcpErr := NewTaskNotFoundError(taskID)
	mcpErr.Operation = operation
	return w.WrapError(mcpErr)
}
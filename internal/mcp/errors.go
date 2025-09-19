package mcp

import (
	"encoding/json"
	"fmt"
	"time"
)

// ErrorCode represents different types of errors that can occur in MCP operations
type ErrorCode string

const (
	// Input validation errors
	ErrorCodeInvalidInput     ErrorCode = "INVALID_INPUT"
	ErrorCodeMissingRequired  ErrorCode = "MISSING_REQUIRED_FIELD"
	ErrorCodeInvalidTaskID    ErrorCode = "INVALID_TASK_ID"
	ErrorCodeInvalidPriority  ErrorCode = "INVALID_PRIORITY"
	ErrorCodeInvalidStatus    ErrorCode = "INVALID_STATUS"

	// Task not found errors
	ErrorCodeTaskNotFound     ErrorCode = "TASK_NOT_FOUND"
	ErrorCodeParentNotFound   ErrorCode = "PARENT_TASK_NOT_FOUND"

	// Business logic errors
	ErrorCodeInvalidOperation ErrorCode = "INVALID_OPERATION"
	ErrorCodeCyclicDependency ErrorCode = "CYCLIC_DEPENDENCY"
	ErrorCodePermissionDenied ErrorCode = "PERMISSION_DENIED"

	// System errors
	ErrorCodeSystemError      ErrorCode = "SYSTEM_ERROR"
	ErrorCodeStorageError     ErrorCode = "STORAGE_ERROR"
	ErrorCodeCommitError      ErrorCode = "COMMIT_ERROR"
	ErrorCodeParseError       ErrorCode = "PARSE_ERROR"

	// MCP specific errors
	ErrorCodeResponseTooLarge ErrorCode = "RESPONSE_TOO_LARGE"
	ErrorCodeInternalError    ErrorCode = "INTERNAL_ERROR"
)

// ErrorCategory represents the category of error for handling and logging
type ErrorCategory string

const (
	CategoryValidation ErrorCategory = "validation"
	CategoryNotFound   ErrorCategory = "not_found"
	CategoryBusiness   ErrorCategory = "business_logic"
	CategorySystem     ErrorCategory = "system"
	CategoryMCP        ErrorCategory = "mcp"
)

// ErrorDetails provides additional context about the error
type ErrorDetails struct {
	Field       string                 `json:"field,omitempty"`
	Value       interface{}            `json:"value,omitempty"`
	Expected    string                 `json:"expected,omitempty"`
	Constraints map[string]interface{} `json:"constraints,omitempty"`
	Context     map[string]interface{} `json:"context,omitempty"`
}

// MCPError represents a structured error response for MCP operations
type MCPError struct {
	Code      ErrorCode     `json:"code"`
	Message   string        `json:"message"`
	Category  ErrorCategory `json:"category"`
	Details   *ErrorDetails `json:"details,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
	Operation string        `json:"operation,omitempty"`
}

// Error implements the error interface
func (e *MCPError) Error() string {
	if e.Details != nil && e.Details.Field != "" {
		return fmt.Sprintf("%s: %s (field: %s)", e.Code, e.Message, e.Details.Field)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// ToJSON converts the error to a JSON string for MCP responses
func (e *MCPError) ToJSON() string {
	b, err := json.Marshal(e)
	if err != nil {
		// Fallback to simple error if marshaling fails
		return fmt.Sprintf(`{"code":"INTERNAL_ERROR","message":"Error serialization failed: %s","category":"system","timestamp":"%s"}`,
			err.Error(), time.Now().Format(time.RFC3339))
	}
	return string(b)
}

// NewMCPError creates a new structured MCP error
func NewMCPError(code ErrorCode, message string, category ErrorCategory) *MCPError {
	return &MCPError{
		Code:      code,
		Message:   message,
		Category:  category,
		Timestamp: time.Now(),
	}
}

// NewValidationError creates a validation error with field details
func NewValidationError(field, message string, value interface{}) *MCPError {
	return &MCPError{
		Code:     ErrorCodeInvalidInput,
		Message:  message,
		Category: CategoryValidation,
		Details: &ErrorDetails{
			Field: field,
			Value: value,
		},
		Timestamp: time.Now(),
	}
}

// NewMissingRequiredError creates an error for missing required fields
func NewMissingRequiredError(field string) *MCPError {
	return &MCPError{
		Code:     ErrorCodeMissingRequired,
		Message:  fmt.Sprintf("Required field '%s' is missing", field),
		Category: CategoryValidation,
		Details: &ErrorDetails{
			Field: field,
		},
		Timestamp: time.Now(),
	}
}

// NewTaskNotFoundError creates a task not found error
func NewTaskNotFoundError(taskID string) *MCPError {
	return &MCPError{
		Code:     ErrorCodeTaskNotFound,
		Message:  fmt.Sprintf("Task with ID '%s' not found", taskID),
		Category: CategoryNotFound,
		Details: &ErrorDetails{
			Field: "id",
			Value: taskID,
		},
		Timestamp: time.Now(),
	}
}

// NewInvalidTaskIDError creates an error for invalid task ID format
func NewInvalidTaskIDError(taskID string, reason string) *MCPError {
	return &MCPError{
		Code:     ErrorCodeInvalidTaskID,
		Message:  fmt.Sprintf("Invalid task ID '%s': %s", taskID, reason),
		Category: CategoryValidation,
		Details: &ErrorDetails{
			Field: "id",
			Value: taskID,
			Expected: "Valid task ID format (e.g., 'T1', '1.2', 'T1.2.3')",
		},
		Timestamp: time.Now(),
	}
}

// NewSystemError creates a system error
func NewSystemError(operation, message string, err error) *MCPError {
	var context map[string]interface{}
	if err != nil {
		context = map[string]interface{}{
			"underlying_error": err.Error(),
		}
	}

	return &MCPError{
		Code:      ErrorCodeSystemError,
		Message:   message,
		Category:  CategorySystem,
		Operation: operation,
		Details: &ErrorDetails{
			Context: context,
		},
		Timestamp: time.Now(),
	}
}

// NewStorageError creates a storage-related error
func NewStorageError(operation, message string, err error) *MCPError {
	var context map[string]interface{}
	if err != nil {
		context = map[string]interface{}{
			"underlying_error": err.Error(),
		}
	}

	return &MCPError{
		Code:      ErrorCodeStorageError,
		Message:   message,
		Category:  CategorySystem,
		Operation: operation,
		Details: &ErrorDetails{
			Context: context,
		},
		Timestamp: time.Now(),
	}
}

// WrapError wraps an existing error and converts it to an MCPError
func WrapError(err error, operation string) *MCPError {
	if mcpErr, ok := err.(*MCPError); ok {
		// Already an MCPError, just update operation if not set
		if mcpErr.Operation == "" {
			mcpErr.Operation = operation
		}
		return mcpErr
	}

	// Determine error type based on error message patterns
	errMsg := err.Error()

	// Check for task ID parsing errors
	if contains(errMsg, "invalid task ID") || contains(errMsg, "invalid segment") {
		return &MCPError{
			Code:      ErrorCodeInvalidTaskID,
			Message:   errMsg,
			Category:  CategoryValidation,
			Operation: operation,
			Details: &ErrorDetails{
				Context: map[string]interface{}{
					"underlying_error": errMsg,
				},
			},
			Timestamp: time.Now(),
		}
	}

	// Check for not found errors
	if contains(errMsg, "not found") || contains(errMsg, "no such file") {
		return &MCPError{
			Code:      ErrorCodeTaskNotFound,
			Message:   errMsg,
			Category:  CategoryNotFound,
			Operation: operation,
			Details: &ErrorDetails{
				Context: map[string]interface{}{
					"underlying_error": errMsg,
				},
			},
			Timestamp: time.Now(),
		}
	}

	// Check for parsing errors
	if contains(errMsg, "parse") || contains(errMsg, "unmarshal") || contains(errMsg, "yaml") {
		return &MCPError{
			Code:      ErrorCodeParseError,
			Message:   errMsg,
			Category:  CategorySystem,
			Operation: operation,
			Details: &ErrorDetails{
				Context: map[string]interface{}{
					"underlying_error": errMsg,
				},
			},
			Timestamp: time.Now(),
		}
	}

	// Default to system error
	return &MCPError{
		Code:      ErrorCodeSystemError,
		Message:   errMsg,
		Category:  CategorySystem,
		Operation: operation,
		Details: &ErrorDetails{
			Context: map[string]interface{}{
				"underlying_error": errMsg,
			},
		},
		Timestamp: time.Now(),
	}
}

// contains is a helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		 containsAt(s, substr, 1))))
}

func containsAt(s, substr string, start int) bool {
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
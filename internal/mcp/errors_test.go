package mcp

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestMCPError_Error(t *testing.T) {
	tests := []struct {
		name     string
		error    *MCPError
		expected string
	}{
		{
			name: "error with field details",
			error: &MCPError{
				Code:    ErrorCodeInvalidInput,
				Message: "Invalid value provided",
				Details: &ErrorDetails{
					Field: "title",
				},
			},
			expected: "INVALID_INPUT: Invalid value provided (field: title)",
		},
		{
			name: "error without field details",
			error: &MCPError{
				Code:    ErrorCodeTaskNotFound,
				Message: "Task not found",
			},
			expected: "TASK_NOT_FOUND: Task not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.error.Error(); got != tt.expected {
				t.Errorf("MCPError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMCPError_ToJSON(t *testing.T) {
	err := &MCPError{
		Code:      ErrorCodeInvalidInput,
		Message:   "Test error",
		Category:  CategoryValidation,
		Timestamp: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		Details: &ErrorDetails{
			Field: "test_field",
			Value: "test_value",
		},
	}

	jsonStr := err.ToJSON()

	// Verify it's valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		t.Fatalf("ToJSON() produced invalid JSON: %v", err)
	}

	// Check required fields
	if parsed["code"] != string(ErrorCodeInvalidInput) {
		t.Errorf("Expected code %s, got %v", ErrorCodeInvalidInput, parsed["code"])
	}
	if parsed["message"] != "Test error" {
		t.Errorf("Expected message 'Test error', got %v", parsed["message"])
	}
	if parsed["category"] != string(CategoryValidation) {
		t.Errorf("Expected category %s, got %v", CategoryValidation, parsed["category"])
	}
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("title", "Title is required", nil)

	if err.Code != ErrorCodeInvalidInput {
		t.Errorf("Expected code %s, got %s", ErrorCodeInvalidInput, err.Code)
	}
	if err.Category != CategoryValidation {
		t.Errorf("Expected category %s, got %s", CategoryValidation, err.Category)
	}
	if err.Details.Field != "title" {
		t.Errorf("Expected field 'title', got %s", err.Details.Field)
	}
	if err.Message != "Title is required" {
		t.Errorf("Expected message 'Title is required', got %s", err.Message)
	}
}

func TestNewMissingRequiredError(t *testing.T) {
	err := NewMissingRequiredError("description")

	if err.Code != ErrorCodeMissingRequired {
		t.Errorf("Expected code %s, got %s", ErrorCodeMissingRequired, err.Code)
	}
	if err.Category != CategoryValidation {
		t.Errorf("Expected category %s, got %s", CategoryValidation, err.Category)
	}
	if err.Details.Field != "description" {
		t.Errorf("Expected field 'description', got %s", err.Details.Field)
	}
	expectedMsg := "Required field 'description' is missing"
	if err.Message != expectedMsg {
		t.Errorf("Expected message '%s', got %s", expectedMsg, err.Message)
	}
}

func TestNewTaskNotFoundError(t *testing.T) {
	err := NewTaskNotFoundError("T1.2.3")

	if err.Code != ErrorCodeTaskNotFound {
		t.Errorf("Expected code %s, got %s", ErrorCodeTaskNotFound, err.Code)
	}
	if err.Category != CategoryNotFound {
		t.Errorf("Expected category %s, got %s", CategoryNotFound, err.Category)
	}
	if err.Details.Field != "id" {
		t.Errorf("Expected field 'id', got %s", err.Details.Field)
	}
	if err.Details.Value != "T1.2.3" {
		t.Errorf("Expected value 'T1.2.3', got %v", err.Details.Value)
	}
}

func TestNewInvalidTaskIDError(t *testing.T) {
	err := NewInvalidTaskIDError("invalid-id", "contains invalid characters")

	if err.Code != ErrorCodeInvalidTaskID {
		t.Errorf("Expected code %s, got %s", ErrorCodeInvalidTaskID, err.Code)
	}
	if err.Category != CategoryValidation {
		t.Errorf("Expected category %s, got %s", CategoryValidation, err.Category)
	}
	if err.Details.Field != "id" {
		t.Errorf("Expected field 'id', got %s", err.Details.Field)
	}
	if err.Details.Value != "invalid-id" {
		t.Errorf("Expected value 'invalid-id', got %v", err.Details.Value)
	}
	expectedMsg := "Invalid task ID 'invalid-id': contains invalid characters"
	if err.Message != expectedMsg {
		t.Errorf("Expected message '%s', got %s", expectedMsg, err.Message)
	}
}

func TestWrapError(t *testing.T) {
	tests := []struct {
		name          string
		inputError    error
		operation     string
		expectedCode  ErrorCode
		expectedCat   ErrorCategory
	}{
		{
			name:          "task ID parsing error",
			inputError:    errors.New("invalid task ID 'abc': invalid segment"),
			operation:     "get_task",
			expectedCode:  ErrorCodeInvalidTaskID,
			expectedCat:   CategoryValidation,
		},
		{
			name:          "not found error",
			inputError:    errors.New("task not found"),
			operation:     "get_task",
			expectedCode:  ErrorCodeTaskNotFound,
			expectedCat:   CategoryNotFound,
		},
		{
			name:          "parse error",
			inputError:    errors.New("failed to parse yaml"),
			operation:     "load_task",
			expectedCode:  ErrorCodeParseError,
			expectedCat:   CategorySystem,
		},
		{
			name:          "generic system error",
			inputError:    errors.New("unexpected failure"),
			operation:     "create_task",
			expectedCode:  ErrorCodeSystemError,
			expectedCat:   CategorySystem,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapped := WrapError(tt.inputError, tt.operation)

			if wrapped.Code != tt.expectedCode {
				t.Errorf("Expected code %s, got %s", tt.expectedCode, wrapped.Code)
			}
			if wrapped.Category != tt.expectedCat {
				t.Errorf("Expected category %s, got %s", tt.expectedCat, wrapped.Category)
			}
			if wrapped.Operation != tt.operation {
				t.Errorf("Expected operation %s, got %s", tt.operation, wrapped.Operation)
			}
			if wrapped.Details == nil || wrapped.Details.Context == nil {
				t.Error("Expected context details to be set")
			}
		})
	}
}

func TestWrapError_AlreadyMCPError(t *testing.T) {
	originalErr := NewTaskNotFoundError("T1")
	originalErr.Operation = "original_op"

	wrapped := WrapError(originalErr, "new_op")

	// Should return the same error, not wrap it again
	if wrapped != originalErr {
		t.Error("Expected the same MCPError instance")
	}

	// Operation should remain unchanged since it was already set
	if wrapped.Operation != "original_op" {
		t.Errorf("Expected operation 'original_op', got %s", wrapped.Operation)
	}
}

func TestWrapError_SetOperationIfEmpty(t *testing.T) {
	originalErr := NewTaskNotFoundError("T1")
	originalErr.Operation = "" // Empty operation

	wrapped := WrapError(originalErr, "new_op")

	// Should update the operation since it was empty
	if wrapped.Operation != "new_op" {
		t.Errorf("Expected operation 'new_op', got %s", wrapped.Operation)
	}
}

func TestContainsHelper(t *testing.T) {
	tests := []struct {
		s        string
		substr   string
		expected bool
	}{
		{"hello world", "world", true},
		{"hello world", "hello", true},
		{"hello world", "lo wo", true},
		{"hello world", "xyz", false},
		{"test", "test", true},
		{"test", "testing", false},
		{"", "test", false},
		{"test", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.s+"_contains_"+tt.substr, func(t *testing.T) {
			if got := contains(tt.s, tt.substr); got != tt.expected {
				t.Errorf("contains(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.expected)
			}
		})
	}
}
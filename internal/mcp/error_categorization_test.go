package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/afero"
	"github.com/veggiemonk/backlog/internal/core"
)

// TestErrorCategorization tests that all error types are properly categorized and handled
func TestErrorCategorization(t *testing.T) {
	// Setup test environment
	fs := afero.NewMemMapFs()
	// Create the .backlog directory
	fs.MkdirAll(".backlog", 0o755)
	store := core.NewFileTaskStore(fs, ".backlog")

	// Initialize handler with all middleware
	responseSizeConfig := DefaultResponseSizeConfig()
	middleware := NewResponseSizeMiddleware(responseSizeConfig)
	validator := NewValidationMiddleware()
	responder := NewResponseWrapper(middleware)

	handler := &handler{
		store:      store,
		validator:  validator,
		responder:  responder,
		middleware: middleware,
		mu:         &sync.Mutex{},
	}

	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	t.Run("ValidationErrors", func(t *testing.T) {
		t.Run("MissingRequiredField", func(t *testing.T) {
			// Test missing title in create
			params := core.CreateTaskParams{
				Description: "Test without title",
			}
			result, _, err := handler.create(ctx, req, params)
			// Should return structured error response
			if err != nil {
				t.Fatalf("Expected no Go error, got: %v", err)
			}
			if result == nil {
				t.Fatal("Expected result, got nil")
			}
			if !result.IsError {
				t.Error("Expected error response")
			}

			// Parse error response
			var mcpErr MCPError
			content := result.Content[0].(*mcp.TextContent)
			if err := json.Unmarshal([]byte(content.Text), &mcpErr); err != nil {
				t.Fatalf("Failed to unmarshal error: %v", err)
			}

			// Verify error categorization
			if mcpErr.Code != ErrorCodeMissingRequired {
				t.Errorf("Expected error code %s, got %s", ErrorCodeMissingRequired, mcpErr.Code)
			}
			if mcpErr.Category != CategoryValidation {
				t.Errorf("Expected category %s, got %s", CategoryValidation, mcpErr.Category)
			}
		})

		t.Run("InvalidTaskID", func(t *testing.T) {
			// Test invalid task ID format
			params := ViewParams{ID: "invalid-id-format"}
			result, _, err := handler.view(ctx, req, params)
			// Should return structured error response
			if err != nil {
				t.Fatalf("Expected no Go error, got: %v", err)
			}
			if !result.IsError {
				t.Error("Expected error response")
			}

			// Parse error response
			var mcpErr MCPError
			content := result.Content[0].(*mcp.TextContent)
			if err := json.Unmarshal([]byte(content.Text), &mcpErr); err != nil {
				t.Fatalf("Failed to unmarshal error: %v", err)
			}

			// Verify error categorization
			if mcpErr.Code != ErrorCodeInvalidTaskID {
				t.Errorf("Expected error code %s, got %s", ErrorCodeInvalidTaskID, mcpErr.Code)
			}
			if mcpErr.Category != CategoryValidation {
				t.Errorf("Expected category %s, got %s", CategoryValidation, mcpErr.Category)
			}
		})

		t.Run("InvalidPriority", func(t *testing.T) {
			// Test invalid priority value
			params := core.CreateTaskParams{
				Title:       "Test Task",
				Description: "Test Description",
				Priority:    "invalid-priority",
			}
			result, _, err := handler.create(ctx, req, params)
			// Should return structured error response
			if err != nil {
				t.Fatalf("Expected no Go error, got: %v", err)
			}
			if !result.IsError {
				t.Error("Expected error response")
			}

			// Parse error response
			var mcpErr MCPError
			content := result.Content[0].(*mcp.TextContent)
			if err := json.Unmarshal([]byte(content.Text), &mcpErr); err != nil {
				t.Fatalf("Failed to unmarshal error: %v", err)
			}

			// Verify error categorization
			if mcpErr.Code != ErrorCodeInvalidPriority {
				t.Errorf("Expected error code %s, got %s", ErrorCodeInvalidPriority, mcpErr.Code)
			}
			if mcpErr.Category != CategoryValidation {
				t.Errorf("Expected category %s, got %s", CategoryValidation, mcpErr.Category)
			}
		})
	})

	t.Run("NotFoundErrors", func(t *testing.T) {
		t.Run("TaskNotFound", func(t *testing.T) {
			// Test viewing non-existent task
			params := ViewParams{ID: "T999"}
			result, _, err := handler.view(ctx, req, params)
			// Should return structured error response
			if err != nil {
				t.Fatalf("Expected no Go error, got: %v", err)
			}
			if !result.IsError {
				t.Error("Expected error response")
			}

			// Parse error response
			var mcpErr MCPError
			content := result.Content[0].(*mcp.TextContent)
			if err := json.Unmarshal([]byte(content.Text), &mcpErr); err != nil {
				t.Fatalf("Failed to unmarshal error: %v", err)
			}

			// Verify error categorization
			if mcpErr.Code != ErrorCodeTaskNotFound {
				t.Errorf("Expected error code %s, got %s", ErrorCodeTaskNotFound, mcpErr.Code)
			}
			if mcpErr.Category != CategoryNotFound {
				t.Errorf("Expected category %s, got %s", CategoryNotFound, mcpErr.Category)
			}
		})
	})

	t.Run("BusinessLogicErrors", func(t *testing.T) {
		t.Run("BatchCreationPartialFailure", func(t *testing.T) {
			// Create a batch with some invalid tasks to test partial failure handling
			tasks := []core.CreateTaskParams{
				{Title: "Valid Task", Description: "Valid description"},
				{Title: "", Description: "Invalid - no title"}, // This will fail validation
			}
			params := ListCreateParams{Tasks: tasks}

			result, _, err := handler.batchCreate(ctx, req, params)
			// Should return structured error response for the validation failure
			if err != nil {
				t.Fatalf("Expected no Go error, got: %v", err)
			}
			if !result.IsError {
				t.Error("Expected error response")
			}

			// Parse error response
			var mcpErr MCPError
			content := result.Content[0].(*mcp.TextContent)
			if err := json.Unmarshal([]byte(content.Text), &mcpErr); err != nil {
				t.Fatalf("Failed to unmarshal error: %v", err)
			}

			// Verify error categorization - should be validation error for the invalid task
			if mcpErr.Category != CategoryValidation {
				t.Errorf("Expected category %s, got %s", CategoryValidation, mcpErr.Category)
			}
		})
	})

	t.Run("MCPErrors", func(t *testing.T) {
		t.Run("ResponseTooLarge", func(t *testing.T) {
			// Create many tasks to trigger response size limit
			for range 100 {
				task, _ := store.Create(core.CreateTaskParams{
					Title:       "Large Dataset Task",
					Description: "This is a task for testing large response handling",
				})
				_ = task // Ignore the task, just populate the store
			}

			// Try to list all tasks without pagination
			params := core.ListTasksParams{}
			result, _, err := handler.list(ctx, req, params)
			// Should return structured error response for response too large
			if err != nil {
				t.Fatalf("Expected no Go error, got: %v", err)
			}
			if !result.IsError {
				t.Error("Expected error response")
			}

			// Parse error response
			var mcpErr MCPError
			content := result.Content[0].(*mcp.TextContent)
			if err := json.Unmarshal([]byte(content.Text), &mcpErr); err != nil {
				t.Fatalf("Failed to unmarshal error: %v", err)
			}

			// Verify error categorization
			if mcpErr.Code != ErrorCodeResponseTooLarge {
				t.Errorf("Expected error code %s, got %s", ErrorCodeResponseTooLarge, mcpErr.Code)
			}
			if mcpErr.Category != CategoryMCP {
				t.Errorf("Expected category %s, got %s", CategoryMCP, mcpErr.Category)
			}
		})
	})

	t.Run("SystemErrors", func(t *testing.T) {
		t.Run("ErrorWrapping", func(t *testing.T) {
			// Test error wrapping functionality
			originalErr := errors.New("test system error")
			mcpErr := WrapError(originalErr, "test_operation")

			// Verify error categorization
			if mcpErr.Code != ErrorCodeSystemError {
				t.Errorf("Expected error code %s, got %s", ErrorCodeSystemError, mcpErr.Code)
			}
			if mcpErr.Category != CategorySystem {
				t.Errorf("Expected category %s, got %s", CategorySystem, mcpErr.Category)
			}
			if mcpErr.Operation != "test_operation" {
				t.Errorf("Expected operation 'test_operation', got %s", mcpErr.Operation)
			}
		})

		t.Run("ParseErrors", func(t *testing.T) {
			// Test parse error detection
			parseErr := errors.New("failed to parse yaml: invalid format")
			mcpErr := WrapError(parseErr, "parse_operation")

			// Should be categorized as parse error
			if mcpErr.Code != ErrorCodeParseError {
				t.Errorf("Expected error code %s, got %s", ErrorCodeParseError, mcpErr.Code)
			}
			if mcpErr.Category != CategorySystem {
				t.Errorf("Expected category %s, got %s", CategorySystem, mcpErr.Category)
			}
		})
	})
}

// TestErrorResponseStructure tests that all error responses have consistent structure
func TestErrorResponseStructure(t *testing.T) {
	tests := []struct {
		name           string
		errorCode      ErrorCode
		category       ErrorCategory
		expectedFields []string
	}{
		{
			name:           "ValidationError",
			errorCode:      ErrorCodeInvalidInput,
			category:       CategoryValidation,
			expectedFields: []string{"code", "message", "category", "timestamp"},
		},
		{
			name:           "NotFoundError",
			errorCode:      ErrorCodeTaskNotFound,
			category:       CategoryNotFound,
			expectedFields: []string{"code", "message", "category", "timestamp"},
		},
		{
			name:           "SystemError",
			errorCode:      ErrorCodeSystemError,
			category:       CategorySystem,
			expectedFields: []string{"code", "message", "category", "timestamp"},
		},
		{
			name:           "MCPError",
			errorCode:      ErrorCodeResponseTooLarge,
			category:       CategoryMCP,
			expectedFields: []string{"code", "message", "category", "timestamp"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create error
			mcpErr := NewMCPError(tt.errorCode, "Test error message", tt.category)
			mcpErr.Operation = "test_operation"

			// Convert to JSON
			jsonStr := mcpErr.ToJSON()

			// Parse back to verify structure
			var parsed map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
				t.Fatalf("Failed to parse error JSON: %v", err)
			}

			// Verify all expected fields are present
			for _, field := range tt.expectedFields {
				if _, exists := parsed[field]; !exists {
					t.Errorf("Expected field '%s' not found in error response", field)
				}
			}

			// Verify specific values
			if parsed["code"] != string(tt.errorCode) {
				t.Errorf("Expected code %s, got %v", tt.errorCode, parsed["code"])
			}
			if parsed["category"] != string(tt.category) {
				t.Errorf("Expected category %s, got %v", tt.category, parsed["category"])
			}
			if parsed["operation"] != "test_operation" {
				t.Errorf("Expected operation 'test_operation', got %v", parsed["operation"])
			}
		})
	}
}

// TestErrorContextAndDetails tests that error details and context are properly included
func TestErrorContextAndDetails(t *testing.T) {
	t.Run("ValidationErrorWithDetails", func(t *testing.T) {
		err := NewValidationError("test_field", "Test validation message", "invalid_value")

		if err.Details == nil {
			t.Fatal("Expected error details, got nil")
		}

		if err.Details.Field != "test_field" {
			t.Errorf("Expected field 'test_field', got %s", err.Details.Field)
		}

		if err.Details.Value != "invalid_value" {
			t.Errorf("Expected value 'invalid_value', got %v", err.Details.Value)
		}
	})

	t.Run("SystemErrorWithContext", func(t *testing.T) {
		originalErr := errors.New("underlying error")
		err := NewSystemError("test_operation", "System failure", originalErr)

		if err.Details == nil {
			t.Fatal("Expected error details, got nil")
		}

		if err.Details.Context == nil {
			t.Fatal("Expected error context, got nil")
		}

		if underlying, ok := err.Details.Context["underlying_error"]; !ok {
			t.Error("Expected underlying_error in context")
		} else if underlying != "underlying error" {
			t.Errorf("Expected 'underlying error', got %v", underlying)
		}
	})
}


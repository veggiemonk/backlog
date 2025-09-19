package mcp

import (
	"strings"
	"testing"

	"github.com/veggiemonk/backlog/internal/core"
)

func TestValidationMiddleware_ValidateTaskID(t *testing.T) {
	v := NewValidationMiddleware()

	tests := []struct {
		name      string
		taskID    string
		fieldName string
		wantError bool
		wantCode  ErrorCode
	}{
		{"valid simple ID", "T1", "id", false, ""},
		{"valid numeric ID", "1", "id", false, ""},
		{"valid hierarchical ID", "T1.2.3", "id", false, ""},
		{"valid numeric hierarchical", "1.2.3", "id", false, ""},
		{"empty ID", "", "id", true, ErrorCodeMissingRequired},
		{"invalid characters", "T1-2", "id", true, ErrorCodeInvalidTaskID},
		{"invalid format", "abc", "id", true, ErrorCodeInvalidTaskID},
		{"leading dots", ".1.2", "id", true, ErrorCodeInvalidTaskID},
		{"trailing dots", "1.2.", "id", true, ErrorCodeInvalidTaskID},
		{"double dots", "1..2", "id", true, ErrorCodeInvalidTaskID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateTaskID(tt.taskID, tt.fieldName)
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.wantError && err != nil && err.Code != tt.wantCode {
				t.Errorf("Expected error code %s, got %s", tt.wantCode, err.Code)
			}
		})
	}
}

func TestValidationMiddleware_ValidateRequired(t *testing.T) {
	v := NewValidationMiddleware()

	tests := []struct {
		name      string
		value     string
		fieldName string
		wantError bool
	}{
		{"valid value", "test", "field", false},
		{"empty string", "", "field", true},
		{"whitespace only", "   ", "field", true},
		{"tabs and spaces", "\t  \n", "field", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateRequired(tt.value, tt.fieldName)
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidationMiddleware_ValidateStringLength(t *testing.T) {
	v := NewValidationMiddleware()

	tests := []struct {
		name      string
		value     string
		maxLength int
		wantError bool
	}{
		{"within limit", "test", 10, false},
		{"at limit", "test", 4, false},
		{"exceeds limit", "test", 3, true},
		{"empty string", "", 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateStringLength(tt.value, "field", tt.maxLength)
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidationMiddleware_ValidatePriority(t *testing.T) {
	v := NewValidationMiddleware()

	tests := []struct {
		name      string
		priority  string
		wantError bool
	}{
		{"valid critical", "critical", false},
		{"valid high", "high", false},
		{"valid medium", "medium", false},
		{"valid low", "low", false},
		{"valid uppercase", "HIGH", false},
		{"valid mixed case", "Critical", false},
		{"empty (optional)", "", false},
		{"invalid value", "urgent", true},
		{"whitespace padding", "  high  ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidatePriority(tt.priority, "priority")
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidationMiddleware_ValidateStatus(t *testing.T) {
	v := NewValidationMiddleware()

	tests := []struct {
		name      string
		status    string
		wantError bool
	}{
		{"valid todo", "todo", false},
		{"valid in-progress", "in-progress", false},
		{"valid done", "done", false},
		{"valid archived", "archived", false},
		{"valid uppercase", "TODO", false},
		{"empty (optional)", "", false},
		{"invalid value", "pending", true},
		{"whitespace padding", "  done  ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateStatus(tt.status, "status")
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidationMiddleware_ValidateArray(t *testing.T) {
	v := NewValidationMiddleware()

	tests := []struct {
		name      string
		arr       []string
		maxItems  int
		wantError bool
	}{
		{"within limit", []string{"a", "b"}, 5, false},
		{"at limit", []string{"a", "b", "c"}, 3, false},
		{"exceeds limit", []string{"a", "b", "c", "d"}, 3, true},
		{"empty array", []string{}, 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateArray(tt.arr, "field", tt.maxItems)
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidationMiddleware_ValidateCreateTaskParams(t *testing.T) {
	v := NewValidationMiddleware()

	validParams := core.CreateTaskParams{
		Title:       "Test Task",
		Description: "Test Description",
		Priority:    "high",
		Assigned:    []string{"user1"},
		Labels:      []string{"test"},
	}

	tests := []struct {
		name      string
		params    core.CreateTaskParams
		wantError bool
		wantCode  ErrorCode
	}{
		{
			name:      "valid params",
			params:    validParams,
			wantError: false,
		},
		{
			name: "missing title",
			params: core.CreateTaskParams{
				Description: "Test Description",
			},
			wantError: true,
			wantCode:  ErrorCodeMissingRequired,
		},
		{
			name: "missing description",
			params: core.CreateTaskParams{
				Title: "Test Task",
			},
			wantError: true,
			wantCode:  ErrorCodeMissingRequired,
		},
		{
			name: "title too long",
			params: core.CreateTaskParams{
				Title:       strings.Repeat("a", 201),
				Description: "Test Description",
			},
			wantError: true,
			wantCode:  ErrorCodeInvalidInput,
		},
		{
			name: "invalid priority",
			params: core.CreateTaskParams{
				Title:       "Test Task",
				Description: "Test Description",
				Priority:    "urgent",
			},
			wantError: true,
			wantCode:  ErrorCodeInvalidPriority,
		},
		{
			name: "invalid parent ID",
			params: func() core.CreateTaskParams {
				p := validParams
				parent := "invalid-id"
				p.Parent = &parent
				return p
			}(),
			wantError: true,
			wantCode:  ErrorCodeInvalidTaskID,
		},
		{
			name: "too many assigned",
			params: core.CreateTaskParams{
				Title:       "Test Task",
				Description: "Test Description",
				Assigned:    make([]string, 11), // Max is 10
			},
			wantError: true,
			wantCode:  ErrorCodeInvalidInput,
		},
		{
			name: "invalid dependency ID",
			params: core.CreateTaskParams{
				Title:        "Test Task",
				Description:  "Test Description",
				Dependencies: []string{"invalid-dep"},
			},
			wantError: true,
			wantCode:  ErrorCodeInvalidTaskID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateCreateTaskParams(tt.params)
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.wantError && err != nil && err.Code != tt.wantCode {
				t.Errorf("Expected error code %s, got %s", tt.wantCode, err.Code)
			}
		})
	}
}

func TestValidationMiddleware_ValidateEditTaskParams(t *testing.T) {
	v := NewValidationMiddleware()

	validParams := core.EditTaskParams{
		ID: "T1",
	}

	tests := []struct {
		name      string
		params    core.EditTaskParams
		wantError bool
		wantCode  ErrorCode
	}{
		{
			name:      "valid params",
			params:    validParams,
			wantError: false,
		},
		{
			name: "invalid task ID",
			params: core.EditTaskParams{
				ID: "invalid-id",
			},
			wantError: true,
			wantCode:  ErrorCodeInvalidTaskID,
		},
		{
			name: "invalid new status",
			params: core.EditTaskParams{
				ID: "T1",
				NewStatus: func() *string {
					s := "invalid-status"
					return &s
				}(),
			},
			wantError: true,
			wantCode:  ErrorCodeInvalidStatus,
		},
		{
			name: "invalid new priority",
			params: core.EditTaskParams{
				ID: "T1",
				NewPriority: func() *string {
					s := "urgent"
					return &s
				}(),
			},
			wantError: true,
			wantCode:  ErrorCodeInvalidPriority,
		},
		{
			name: "invalid AC index",
			params: core.EditTaskParams{
				ID:      "T1",
				CheckAC: []int{0}, // Should be 1-based
			},
			wantError: true,
			wantCode:  ErrorCodeInvalidInput,
		},
		{
			name: "AC index too high",
			params: core.EditTaskParams{
				ID:      "T1",
				CheckAC: []int{1001}, // Max is 1000
			},
			wantError: true,
			wantCode:  ErrorCodeInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateEditTaskParams(tt.params)
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.wantError && err != nil && err.Code != tt.wantCode {
				t.Errorf("Expected error code %s, got %s", tt.wantCode, err.Code)
			}
		})
	}
}

func TestValidationMiddleware_ValidateViewParams(t *testing.T) {
	v := NewValidationMiddleware()

	tests := []struct {
		name      string
		taskID    string
		wantError bool
		wantCode  ErrorCode
	}{
		{"valid ID", "T1", false, ""},
		{"empty ID", "", true, ErrorCodeMissingRequired},
		{"invalid ID", "invalid", true, ErrorCodeInvalidTaskID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateViewParams(tt.taskID)
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.wantError && err != nil && err.Code != tt.wantCode {
				t.Errorf("Expected error code %s, got %s", tt.wantCode, err.Code)
			}
		})
	}
}

func TestValidationMiddleware_ValidateListParams(t *testing.T) {
	v := NewValidationMiddleware()

	tests := []struct {
		name      string
		params    core.ListTasksParams
		wantError bool
		wantCode  ErrorCode
	}{
		{
			name:      "valid empty params",
			params:    core.ListTasksParams{},
			wantError: false,
		},
		{
			name: "invalid parent ID",
			params: core.ListTasksParams{
				Parent: func() *string {
					s := "invalid-id"
					return &s
				}(),
			},
			wantError: true,
			wantCode:  ErrorCodeInvalidTaskID,
		},
		{
			name: "invalid priority",
			params: core.ListTasksParams{
				Priority: func() *string {
					s := "urgent"
					return &s
				}(),
			},
			wantError: true,
			wantCode:  ErrorCodeInvalidPriority,
		},
		{
			name: "invalid status",
			params: core.ListTasksParams{
				Status: []string{"invalid-status"},
			},
			wantError: true,
			wantCode:  ErrorCodeInvalidStatus,
		},
		{
			name: "negative limit",
			params: core.ListTasksParams{
				Limit: func() *int {
					i := -1
					return &i
				}(),
			},
			wantError: true,
			wantCode:  ErrorCodeInvalidInput,
		},
		{
			name: "negative offset",
			params: core.ListTasksParams{
				Offset: func() *int {
					i := -1
					return &i
				}(),
			},
			wantError: true,
			wantCode:  ErrorCodeInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateListParams(tt.params)
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.wantError && err != nil && err.Code != tt.wantCode {
				t.Errorf("Expected error code %s, got %s", tt.wantCode, err.Code)
			}
		})
	}
}
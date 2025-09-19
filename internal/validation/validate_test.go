package validation

import (
	"strings"
	"testing"
)

func TestValidator_ValidateTitle(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name      string
		title     string
		wantCode  string
		wantError bool
	}{
		{
			name:      "valid title",
			title:     "Valid task title",
			wantCode:  "",
			wantError: false,
		},
		{
			name:      "empty title",
			title:     "",
			wantCode:  "TITLE_EMPTY",
			wantError: true,
		},
		{
			name:      "whitespace only title",
			title:     "   ",
			wantCode:  "TITLE_EMPTY",
			wantError: true,
		},
		{
			name:      "title too long",
			title:     strings.Repeat("a", MaxTitleLength+1),
			wantCode:  "TITLE_TOO_LONG",
			wantError: true,
		},
		{
			name:      "title with script tag",
			title:     "Task with <script>alert('xss')</script>",
			wantCode:  "TITLE_INVALID_CHARS",
			wantError: true,
		},
		{
			name:      "title with null byte",
			title:     "Task with\x00null",
			wantCode:  "TITLE_INVALID_CHARS",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateTitle(tt.title)
			if tt.wantError && err.Code == "" {
				t.Errorf("ValidateTitle() expected error but got none")
			}
			if !tt.wantError && err.Code != "" {
				t.Errorf("ValidateTitle() unexpected error: %v", err)
			}
			if tt.wantCode != "" && err.Code != tt.wantCode {
				t.Errorf("ValidateTitle() error code = %v, want %v", err.Code, tt.wantCode)
			}
		})
	}
}

func TestValidator_ValidateTaskID(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name      string
		taskID    string
		wantCode  string
		wantError bool
	}{
		{
			name:      "valid task ID with T prefix",
			taskID:    "T1",
			wantCode:  "",
			wantError: false,
		},
		{
			name:      "valid task ID without T prefix",
			taskID:    "1",
			wantCode:  "",
			wantError: false,
		},
		{
			name:      "valid hierarchical task ID",
			taskID:    "T1.2.3",
			wantCode:  "",
			wantError: false,
		},
		{
			name:      "empty task ID",
			taskID:    "",
			wantCode:  "TASK_ID_EMPTY",
			wantError: true,
		},
		{
			name:      "invalid format",
			taskID:    "invalid",
			wantCode:  "TASK_ID_INVALID_FORMAT",
			wantError: true,
		},
		{
			name:      "task ID too long",
			taskID:    strings.Repeat("1.", MaxDependencyIDLength),
			wantCode:  "TASK_ID_TOO_LONG",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateTaskID(tt.taskID)
			if tt.wantError && err.Code == "" {
				t.Errorf("ValidateTaskID() expected error but got none")
			}
			if !tt.wantError && err.Code != "" {
				t.Errorf("ValidateTaskID() unexpected error: %v", err)
			}
			if tt.wantCode != "" && err.Code != tt.wantCode {
				t.Errorf("ValidateTaskID() error code = %v, want %v", err.Code, tt.wantCode)
			}
		})
	}
}

func TestValidator_ValidateLabel(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name      string
		label     string
		wantCode  string
		wantError bool
	}{
		{
			name:      "valid label",
			label:     "bug",
			wantCode:  "",
			wantError: false,
		},
		{
			name:      "valid label with hyphen",
			label:     "high-priority",
			wantCode:  "",
			wantError: false,
		},
		{
			name:      "valid label with underscore",
			label:     "frontend_task",
			wantCode:  "",
			wantError: false,
		},
		{
			name:      "empty label",
			label:     "",
			wantCode:  "LABEL_EMPTY",
			wantError: true,
		},
		{
			name:      "label with spaces",
			label:     "bug fix",
			wantCode:  "LABEL_INVALID_FORMAT",
			wantError: true,
		},
		{
			name:      "label too long",
			label:     strings.Repeat("a", MaxLabelLength+1),
			wantCode:  "LABEL_TOO_LONG",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateLabel(tt.label)
			if tt.wantError && err.Code == "" {
				t.Errorf("ValidateLabel() expected error but got none")
			}
			if !tt.wantError && err.Code != "" {
				t.Errorf("ValidateLabel() unexpected error: %v", err)
			}
			if tt.wantCode != "" && err.Code != tt.wantCode {
				t.Errorf("ValidateLabel() error code = %v, want %v", err.Code, tt.wantCode)
			}
		})
	}
}

func TestValidator_ValidatePriority(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name      string
		priority  string
		wantCode  string
		wantError bool
	}{
		{
			name:      "valid priority low",
			priority:  "low",
			wantCode:  "",
			wantError: false,
		},
		{
			name:      "valid priority medium",
			priority:  "medium",
			wantCode:  "",
			wantError: false,
		},
		{
			name:      "valid priority high",
			priority:  "high",
			wantCode:  "",
			wantError: false,
		},
		{
			name:      "valid priority critical",
			priority:  "critical",
			wantCode:  "",
			wantError: false,
		},
		{
			name:      "empty priority is valid",
			priority:  "",
			wantCode:  "",
			wantError: false,
		},
		{
			name:      "invalid priority",
			priority:  "urgent",
			wantCode:  "PRIORITY_INVALID",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidatePriority(tt.priority)
			if tt.wantError && err.Code == "" {
				t.Errorf("ValidatePriority() expected error but got none")
			}
			if !tt.wantError && err.Code != "" {
				t.Errorf("ValidatePriority() unexpected error: %v", err)
			}
			if tt.wantCode != "" && err.Code != tt.wantCode {
				t.Errorf("ValidatePriority() error code = %v, want %v", err.Code, tt.wantCode)
			}
		})
	}
}

func TestValidator_ValidateFilePath(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name      string
		path      string
		wantCode  string
		wantError bool
	}{
		{
			name:      "valid absolute path",
			path:      "/home/user/tasks",
			wantCode:  "",
			wantError: false,
		},
		{
			name:      "valid relative path",
			path:      "tasks/backlog",
			wantCode:  "",
			wantError: false,
		},
		{
			name:      "empty path",
			path:      "",
			wantCode:  "PATH_EMPTY",
			wantError: true,
		},
		{
			name:      "path with traversal",
			path:      "../../../etc/passwd",
			wantCode:  "PATH_TRAVERSAL",
			wantError: true,
		},
		{
			name:      "path with null byte",
			path:      "/home/user\x00/tasks",
			wantCode:  "PATH_NULL_BYTE",
			wantError: true,
		},
		{
			name:      "path too long",
			path:      "/" + strings.Repeat("a", MaxPathLength),
			wantCode:  "PATH_TOO_LONG",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateFilePath(tt.path)
			if tt.wantError && err.Code == "" {
				t.Errorf("ValidateFilePath() expected error but got none")
			}
			if !tt.wantError && err.Code != "" {
				t.Errorf("ValidateFilePath() unexpected error: %v", err)
			}
			if tt.wantCode != "" && err.Code != tt.wantCode {
				t.Errorf("ValidateFilePath() error code = %v, want %v", err.Code, tt.wantCode)
			}
		})
	}
}

func TestContainsDangerousChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "safe input",
			input:    "This is a safe string",
			expected: false,
		},
		{
			name:     "script tag",
			input:    "<script>alert('xss')</script>",
			expected: true,
		},
		{
			name:     "javascript protocol",
			input:    "javascript:alert('xss')",
			expected: true,
		},
		{
			name:     "null byte",
			input:    "safe\x00unsafe",
			expected: true,
		},
		{
			name:     "eval function",
			input:    "eval(maliciousCode)",
			expected: true,
		},
		{
			name:     "hex escape",
			input:    "\\x41\\x42",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsDangerousChars(tt.input)
			if result != tt.expected {
				t.Errorf("containsDangerousChars() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValidationErrors_Error(t *testing.T) {
	tests := []struct {
		name     string
		errors   ValidationErrors
		expected string
	}{
		{
			name:     "no errors",
			errors:   ValidationErrors{},
			expected: "no validation errors",
		},
		{
			name: "single error",
			errors: ValidationErrors{
				{Field: "title", Value: "", Message: "title cannot be empty", Code: "TITLE_EMPTY"},
			},
			expected: "validation error in field 'title': title cannot be empty (value: '')",
		},
		{
			name: "multiple errors",
			errors: ValidationErrors{
				{Field: "title", Value: "", Message: "title cannot be empty", Code: "TITLE_EMPTY"},
				{Field: "priority", Value: "invalid", Message: "invalid priority", Code: "PRIORITY_INVALID"},
			},
			expected: "multiple validation errors: validation error in field 'title': title cannot be empty (value: ''); validation error in field 'priority': invalid priority (value: 'invalid')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.errors.Error()
			if result != tt.expected {
				t.Errorf("ValidationErrors.Error() = %v, want %v", result, tt.expected)
			}
		})
	}
}
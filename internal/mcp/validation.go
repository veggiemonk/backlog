package mcp

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/veggiemonk/backlog/internal/core"
)

// ValidationMiddleware provides input validation for MCP operations
type ValidationMiddleware struct {
	taskIDRegex    *regexp.Regexp
	priorityValues map[string]bool
	statusValues   map[string]bool
}

// NewValidationMiddleware creates a new validation middleware
func NewValidationMiddleware() *ValidationMiddleware {
	// TaskID regex: allows T prefix, numbers, and dots (e.g., "T1", "1.2", "T1.2.3")
	taskIDRegex := regexp.MustCompile(`^T?\d+(\.\d+)*$`)

	priorityValues := map[string]bool{
		"critical": true,
		"high":     true,
		"medium":   true,
		"low":      true,
	}

	statusValues := map[string]bool{
		"todo":        true,
		"in-progress": true,
		"done":        true,
		"archived":    true,
	}

	return &ValidationMiddleware{
		taskIDRegex:    taskIDRegex,
		priorityValues: priorityValues,
		statusValues:   statusValues,
	}
}

// ValidateTaskID validates a task ID format and returns an MCPError if invalid
func (v *ValidationMiddleware) ValidateTaskID(taskID string, fieldName string) *MCPError {
	if taskID == "" {
		return NewMissingRequiredError(fieldName)
	}

	if !v.taskIDRegex.MatchString(taskID) {
		return NewInvalidTaskIDError(taskID, "must be in format 'T1', '1.2', or 'T1.2.3'")
	}

	return nil
}

// ValidateRequired validates that a required field is present and not empty
func (v *ValidationMiddleware) ValidateRequired(value string, fieldName string) *MCPError {
	if strings.TrimSpace(value) == "" {
		return NewMissingRequiredError(fieldName)
	}
	return nil
}

// ValidateStringLength validates string length constraints
func (v *ValidationMiddleware) ValidateStringLength(value string, fieldName string, maxLength int) *MCPError {
	if len(value) > maxLength {
		return &MCPError{
			Code:     ErrorCodeInvalidInput,
			Message:  fmt.Sprintf("Field '%s' exceeds maximum length of %d characters", fieldName, maxLength),
			Category: CategoryValidation,
			Details: &ErrorDetails{
				Field: fieldName,
				Value: len(value),
				Expected: fmt.Sprintf("string length <= %d", maxLength),
				Constraints: map[string]interface{}{
					"max_length": maxLength,
					"actual_length": len(value),
				},
			},
		}
	}
	return nil
}

// ValidatePriority validates priority field values
func (v *ValidationMiddleware) ValidatePriority(priority string, fieldName string) *MCPError {
	if priority == "" {
		return nil // Priority is optional
	}

	priority = strings.ToLower(strings.TrimSpace(priority))
	if !v.priorityValues[priority] {
		validValues := make([]string, 0, len(v.priorityValues))
		for k := range v.priorityValues {
			validValues = append(validValues, k)
		}

		return &MCPError{
			Code:     ErrorCodeInvalidPriority,
			Message:  fmt.Sprintf("Invalid priority '%s'. Must be one of: %s", priority, strings.Join(validValues, ", ")),
			Category: CategoryValidation,
			Details: &ErrorDetails{
				Field: fieldName,
				Value: priority,
				Expected: strings.Join(validValues, ", "),
				Constraints: map[string]interface{}{
					"valid_values": validValues,
				},
			},
		}
	}
	return nil
}

// ValidateStatus validates status field values
func (v *ValidationMiddleware) ValidateStatus(status string, fieldName string) *MCPError {
	if status == "" {
		return nil // Status might be optional in some contexts
	}

	status = strings.ToLower(strings.TrimSpace(status))
	if !v.statusValues[status] {
		validValues := make([]string, 0, len(v.statusValues))
		for k := range v.statusValues {
			validValues = append(validValues, k)
		}

		return &MCPError{
			Code:     ErrorCodeInvalidStatus,
			Message:  fmt.Sprintf("Invalid status '%s'. Must be one of: %s", status, strings.Join(validValues, ", ")),
			Category: CategoryValidation,
			Details: &ErrorDetails{
				Field: fieldName,
				Value: status,
				Expected: strings.Join(validValues, ", "),
				Constraints: map[string]interface{}{
					"valid_values": validValues,
				},
			},
		}
	}
	return nil
}

// ValidateArray validates array constraints
func (v *ValidationMiddleware) ValidateArray(arr []string, fieldName string, maxItems int) *MCPError {
	if len(arr) > maxItems {
		return &MCPError{
			Code:     ErrorCodeInvalidInput,
			Message:  fmt.Sprintf("Field '%s' exceeds maximum number of items (%d)", fieldName, maxItems),
			Category: CategoryValidation,
			Details: &ErrorDetails{
				Field: fieldName,
				Value: len(arr),
				Expected: fmt.Sprintf("array length <= %d", maxItems),
				Constraints: map[string]interface{}{
					"max_items": maxItems,
					"actual_items": len(arr),
				},
			},
		}
	}
	return nil
}

// ValidateCreateTaskParams validates parameters for task creation
func (v *ValidationMiddleware) ValidateCreateTaskParams(params core.CreateTaskParams) *MCPError {
	// Validate required fields
	if err := v.ValidateRequired(params.Title, "title"); err != nil {
		return err
	}
	if err := v.ValidateRequired(params.Description, "description"); err != nil {
		return err
	}

	// Validate string lengths
	if err := v.ValidateStringLength(params.Title, "title", 200); err != nil {
		return err
	}
	if err := v.ValidateStringLength(params.Description, "description", 5000); err != nil {
		return err
	}

	// Validate priority if provided
	if err := v.ValidatePriority(params.Priority, "priority"); err != nil {
		return err
	}

	// Validate parent task ID if provided
	if params.Parent != nil && *params.Parent != "" {
		if err := v.ValidateTaskID(*params.Parent, "parent"); err != nil {
			return err
		}
	}

	// Validate arrays
	if err := v.ValidateArray(params.Assigned, "assigned", 10); err != nil {
		return err
	}
	if err := v.ValidateArray(params.Labels, "labels", 20); err != nil {
		return err
	}
	if err := v.ValidateArray(params.Dependencies, "dependencies", 50); err != nil {
		return err
	}
	if err := v.ValidateArray(params.AC, "ac", 100); err != nil {
		return err
	}

	// Validate dependencies task IDs
	for i, depID := range params.Dependencies {
		if err := v.ValidateTaskID(depID, fmt.Sprintf("dependencies[%d]", i)); err != nil {
			return err
		}
	}

	// Validate plan and notes lengths if provided
	if params.Plan != nil {
		if err := v.ValidateStringLength(*params.Plan, "plan", 10000); err != nil {
			return err
		}
	}
	if params.Notes != nil {
		if err := v.ValidateStringLength(*params.Notes, "notes", 10000); err != nil {
			return err
		}
	}

	return nil
}

// ValidateEditTaskParams validates parameters for task editing
func (v *ValidationMiddleware) ValidateEditTaskParams(params core.EditTaskParams) *MCPError {
	// Validate task ID
	if err := v.ValidateTaskID(params.ID, "id"); err != nil {
		return err
	}

	// Validate optional string fields if provided
	if params.NewTitle != nil {
		if err := v.ValidateStringLength(*params.NewTitle, "new_title", 200); err != nil {
			return err
		}
	}
	if params.NewDescription != nil {
		if err := v.ValidateStringLength(*params.NewDescription, "new_description", 5000); err != nil {
			return err
		}
	}

	// Validate status if provided
	if params.NewStatus != nil {
		if err := v.ValidateStatus(*params.NewStatus, "new_status"); err != nil {
			return err
		}
	}

	// Validate priority if provided
	if params.NewPriority != nil {
		if err := v.ValidatePriority(*params.NewPriority, "new_priority"); err != nil {
			return err
		}
	}

	// Validate parent task ID if provided
	if params.NewParent != nil && *params.NewParent != "" {
		if err := v.ValidateTaskID(*params.NewParent, "new_parent"); err != nil {
			return err
		}
	}

	// Validate arrays
	if err := v.ValidateArray(params.AddAssigned, "add_assigned", 10); err != nil {
		return err
	}
	if err := v.ValidateArray(params.RemoveAssigned, "remove_assigned", 10); err != nil {
		return err
	}
	if err := v.ValidateArray(params.AddLabels, "add_labels", 20); err != nil {
		return err
	}
	if err := v.ValidateArray(params.RemoveLabels, "remove_labels", 20); err != nil {
		return err
	}
	if err := v.ValidateArray(params.NewDependencies, "new_dependencies", 50); err != nil {
		return err
	}
	if err := v.ValidateArray(params.AddAC, "add_ac", 100); err != nil {
		return err
	}

	// Validate dependencies task IDs
	for i, depID := range params.NewDependencies {
		if err := v.ValidateTaskID(depID, fmt.Sprintf("new_dependencies[%d]", i)); err != nil {
			return err
		}
	}

	// Validate plan and notes lengths if provided
	if params.NewPlan != nil {
		if err := v.ValidateStringLength(*params.NewPlan, "new_plan", 10000); err != nil {
			return err
		}
	}
	if params.NewNotes != nil {
		if err := v.ValidateStringLength(*params.NewNotes, "new_notes", 10000); err != nil {
			return err
		}
	}

	// Validate AC indices
	if err := v.validateACIndices(params.CheckAC, "check_ac"); err != nil {
		return err
	}
	if err := v.validateACIndices(params.UncheckAC, "uncheck_ac"); err != nil {
		return err
	}
	if err := v.validateACIndices(params.RemoveAC, "remove_ac"); err != nil {
		return err
	}

	return nil
}

// validateACIndices validates acceptance criteria indices
func (v *ValidationMiddleware) validateACIndices(indices []int, fieldName string) *MCPError {
	for i, idx := range indices {
		if idx < 1 || idx > 1000 { // Reasonable upper bound
			return &MCPError{
				Code:     ErrorCodeInvalidInput,
				Message:  fmt.Sprintf("Invalid index in '%s': %d. Indices must be 1-based and <= 1000", fieldName, idx),
				Category: CategoryValidation,
				Details: &ErrorDetails{
					Field: fmt.Sprintf("%s[%d]", fieldName, i),
					Value: idx,
					Expected: "1-based index between 1 and 1000",
					Constraints: map[string]interface{}{
						"min_value": 1,
						"max_value": 1000,
					},
				},
			}
		}
	}
	return nil
}

// ValidateViewParams validates parameters for viewing a task
func (v *ValidationMiddleware) ValidateViewParams(taskID string) *MCPError {
	return v.ValidateTaskID(taskID, "id")
}

// ValidateListParams validates parameters for listing tasks
func (v *ValidationMiddleware) ValidateListParams(params core.ListTasksParams) *MCPError {
	// Validate parent task ID if provided
	if params.Parent != nil && *params.Parent != "" {
		if err := v.ValidateTaskID(*params.Parent, "parent"); err != nil {
			return err
		}
	}

	// Validate priority if provided
	if params.Priority != nil {
		if err := v.ValidatePriority(*params.Priority, "priority"); err != nil {
			return err
		}
	}

	// Validate status values if provided
	for i, status := range params.Status {
		if err := v.ValidateStatus(status, fmt.Sprintf("status[%d]", i)); err != nil {
			return err
		}
	}

	// Validate arrays
	if err := v.ValidateArray(params.Assigned, "assigned", 10); err != nil {
		return err
	}
	if err := v.ValidateArray(params.Labels, "labels", 20); err != nil {
		return err
	}
	if err := v.ValidateArray(params.Sort, "sort", 10); err != nil {
		return err
	}

	// Validate limit and offset
	if params.Limit != nil && *params.Limit < 0 {
		return &MCPError{
			Code:     ErrorCodeInvalidInput,
			Message:  "Limit must be >= 0",
			Category: CategoryValidation,
			Details: &ErrorDetails{
				Field: "limit",
				Value: *params.Limit,
				Expected: "integer >= 0",
			},
		}
	}
	if params.Offset != nil && *params.Offset < 0 {
		return &MCPError{
			Code:     ErrorCodeInvalidInput,
			Message:  "Offset must be >= 0",
			Category: CategoryValidation,
			Details: &ErrorDetails{
				Field: "offset",
				Value: *params.Offset,
				Expected: "integer >= 0",
			},
		}
	}

	return nil
}
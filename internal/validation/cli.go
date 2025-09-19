package validation

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/monitoring"
)

// CLIValidator provides validation functions for CLI commands
type CLIValidator struct {
	validator *Validator
}

// NewCLIValidator creates a new CLI validator
func NewCLIValidator() *CLIValidator {
	return &CLIValidator{
		validator: NewValidator(),
	}
}

// ValidateCreateParams validates parameters for task creation
func (cv *CLIValidator) ValidateCreateParams(params core.CreateTaskParams) ValidationErrors {
	var errors ValidationErrors

	// Validate title
	if err := cv.validator.ValidateTitle(params.Title); err.Code != "" {
		monitoring.LogValidationFailure("title", params.Title, err.Message, map[string]interface{}{
			"code": err.Code,
			"operation": "create_task",
		})
		errors = append(errors, err)
	}

	// Validate description
	if err := cv.validator.ValidateDescription(params.Description); err.Code != "" {
		errors = append(errors, err)
	}

	// Validate priority
	if err := cv.validator.ValidatePriority(params.Priority); err.Code != "" {
		errors = append(errors, err)
	}

	// Validate parent ID
	if params.Parent != nil && *params.Parent != "" {
		if err := cv.validator.ValidateTaskID(*params.Parent); err.Code != "" {
			err.Field = "parent"
			errors = append(errors, err)
		}
	}

	// Validate assignees
	if assigneeErrs := cv.validator.ValidateAssignees(params.Assigned); assigneeErrs.HasErrors() {
		errors = append(errors, assigneeErrs...)
	}

	// Validate labels
	if labelErrs := cv.validator.ValidateLabels(params.Labels); labelErrs.HasErrors() {
		errors = append(errors, labelErrs...)
	}

	// Validate dependencies
	if depErrs := cv.validator.ValidateDependencies(params.Dependencies); depErrs.HasErrors() {
		errors = append(errors, depErrs...)
	}

	// Validate acceptance criteria
	if acErrs := cv.validator.ValidateAcceptanceCriteria(params.AC); acErrs.HasErrors() {
		errors = append(errors, acErrs...)
	}

	// Validate plan
	if params.Plan != nil {
		if err := cv.validator.ValidatePlan(*params.Plan); err.Code != "" {
			errors = append(errors, err)
		}
	}

	// Validate notes
	if params.Notes != nil {
		if err := cv.validator.ValidateNotes(*params.Notes); err.Code != "" {
			errors = append(errors, err)
		}
	}

	return errors
}

// ValidateEditParams validates parameters for task editing
func (cv *CLIValidator) ValidateEditParams(params core.EditTaskParams) ValidationErrors {
	var errors ValidationErrors

	// Validate task ID
	if err := cv.validator.ValidateTaskID(params.ID); err.Code != "" {
		err.Field = "id"
		errors = append(errors, err)
	}

	// Validate new title
	if params.NewTitle != nil {
		if err := cv.validator.ValidateTitle(*params.NewTitle); err.Code != "" {
			err.Field = "new_title"
			errors = append(errors, err)
		}
	}

	// Validate new description
	if params.NewDescription != nil {
		if err := cv.validator.ValidateDescription(*params.NewDescription); err.Code != "" {
			err.Field = "new_description"
			errors = append(errors, err)
		}
	}

	// Validate new priority
	if params.NewPriority != nil {
		if err := cv.validator.ValidatePriority(*params.NewPriority); err.Code != "" {
			err.Field = "new_priority"
			errors = append(errors, err)
		}
	}

	// Validate new parent
	if params.NewParent != nil && *params.NewParent != "" {
		if err := cv.validator.ValidateTaskID(*params.NewParent); err.Code != "" {
			err.Field = "new_parent"
			errors = append(errors, err)
		}
	}

	// Validate assignees to add
	if assigneeErrs := cv.validator.ValidateAssignees(params.AddAssigned); assigneeErrs.HasErrors() {
		for i := range assigneeErrs {
			assigneeErrs[i].Field = strings.Replace(assigneeErrs[i].Field, "assignees", "add_assigned", 1)
		}
		errors = append(errors, assigneeErrs...)
	}

	// Validate assignees to remove
	if assigneeErrs := cv.validator.ValidateAssignees(params.RemoveAssigned); assigneeErrs.HasErrors() {
		for i := range assigneeErrs {
			assigneeErrs[i].Field = strings.Replace(assigneeErrs[i].Field, "assignees", "remove_assigned", 1)
		}
		errors = append(errors, assigneeErrs...)
	}

	// Validate labels to add
	if labelErrs := cv.validator.ValidateLabels(params.AddLabels); labelErrs.HasErrors() {
		for i := range labelErrs {
			labelErrs[i].Field = strings.Replace(labelErrs[i].Field, "labels", "add_labels", 1)
		}
		errors = append(errors, labelErrs...)
	}

	// Validate labels to remove
	if labelErrs := cv.validator.ValidateLabels(params.RemoveLabels); labelErrs.HasErrors() {
		for i := range labelErrs {
			labelErrs[i].Field = strings.Replace(labelErrs[i].Field, "labels", "remove_labels", 1)
		}
		errors = append(errors, labelErrs...)
	}

	// Validate new dependencies
	if depErrs := cv.validator.ValidateDependencies(params.NewDependencies); depErrs.HasErrors() {
		for i := range depErrs {
			depErrs[i].Field = strings.Replace(depErrs[i].Field, "dependencies", "new_dependencies", 1)
		}
		errors = append(errors, depErrs...)
	}

	// Validate new notes
	if params.NewNotes != nil {
		if err := cv.validator.ValidateNotes(*params.NewNotes); err.Code != "" {
			err.Field = "new_notes"
			errors = append(errors, err)
		}
	}

	// Validate new plan
	if params.NewPlan != nil {
		if err := cv.validator.ValidatePlan(*params.NewPlan); err.Code != "" {
			err.Field = "new_plan"
			errors = append(errors, err)
		}
	}

	// Validate acceptance criteria to add
	if acErrs := cv.validator.ValidateAcceptanceCriteria(params.AddAC); acErrs.HasErrors() {
		for i := range acErrs {
			acErrs[i].Field = strings.Replace(acErrs[i].Field, "acceptance_criteria", "add_ac", 1)
		}
		errors = append(errors, acErrs...)
	}

	return errors
}

// ValidateConfigParams validates configuration parameters
func (cv *CLIValidator) ValidateConfigParams(logLevel, logFormat, logFile, folder string) ValidationErrors {
	var errors ValidationErrors

	// Validate log level
	if err := cv.validator.ValidateLogLevel(logLevel); err.Code != "" {
		errors = append(errors, err)
	}

	// Validate log format
	if err := cv.validator.ValidateLogFormat(logFormat); err.Code != "" {
		errors = append(errors, err)
	}

	// Validate log file path
	if logFile != "" {
		if err := cv.validator.ValidateFilePath(logFile); err.Code != "" {
			err.Field = "log_file"
			errors = append(errors, err)
		}
	}

	// Validate folder path
	if folder != "" {
		if err := cv.validator.ValidateFilePath(folder); err.Code != "" {
			err.Field = "folder"
			errors = append(errors, err)
		}
	}

	return errors
}

// ValidateListParams validates parameters for task listing
func (cv *CLIValidator) ValidateListParams(params core.ListTasksParams) ValidationErrors {
	var errors ValidationErrors

	// Validate parent ID
	if params.Parent != nil && *params.Parent != "" {
		if err := cv.validator.ValidateTaskID(*params.Parent); err.Code != "" {
			err.Field = "parent"
			errors = append(errors, err)
		}
	}

	// Validate priority
	if params.Priority != nil && *params.Priority != "" {
		if err := cv.validator.ValidatePriority(*params.Priority); err.Code != "" {
			errors = append(errors, err)
		}
	}

	// Validate assignees
	if assigneeErrs := cv.validator.ValidateAssignees(params.Assigned); assigneeErrs.HasErrors() {
		errors = append(errors, assigneeErrs...)
	}

	// Validate labels
	if labelErrs := cv.validator.ValidateLabels(params.Labels); labelErrs.HasErrors() {
		errors = append(errors, labelErrs...)
	}

	// Validate status values
	validStatuses := []string{"pending", "in-progress", "done", "archived"}
	for i, status := range params.Status {
		if !contains(validStatuses, status) {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("status[%d]", i),
				Value:   status,
				Message: fmt.Sprintf("status must be one of: %s", strings.Join(validStatuses, ", ")),
				Code:    "STATUS_INVALID",
			})
		}
	}

	// Validate sort fields
	validSortFields := []string{"id", "title", "status", "priority", "created_at", "updated_at", "parent"}
	for i, field := range params.Sort {
		if !contains(validSortFields, field) {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("sort[%d]", i),
				Value:   field,
				Message: fmt.Sprintf("sort field must be one of: %s", strings.Join(validSortFields, ", ")),
				Code:    "SORT_FIELD_INVALID",
			})
		}
	}

	return errors
}

// ValidateArgs validates command arguments
func (cv *CLIValidator) ValidateArgs(cmd *cobra.Command, args []string, expectedCount int) error {
	if len(args) != expectedCount {
		if expectedCount == 1 {
			return fmt.Errorf("exactly one argument is required, got %d", len(args))
		}
		return fmt.Errorf("exactly %d arguments are required, got %d", expectedCount, len(args))
	}

	// Validate each argument for dangerous content
	for i, arg := range args {
		if err := cv.validator.ValidateTitle(arg); err.Code != "" {
			return fmt.Errorf("argument %d is invalid: %s", i+1, err.Message)
		}
	}

	return nil
}

// ValidateSearchQuery validates search query
func (cv *CLIValidator) ValidateSearchQuery(query string) ValidationError {
	if strings.TrimSpace(query) == "" {
		return ValidationError{
			Field:   "query",
			Value:   query,
			Message: "search query cannot be empty",
			Code:    "QUERY_EMPTY",
		}
	}

	if utf8.RuneCountInString(query) > 1000 {
		return ValidationError{
			Field:   "query",
			Value:   query,
			Message: "search query cannot exceed 1000 characters",
			Code:    "QUERY_TOO_LONG",
		}
	}

	if containsDangerousChars(query) {
		return ValidationError{
			Field:   "query",
			Value:   query,
			Message: "search query contains potentially dangerous characters",
			Code:    "QUERY_INVALID_CHARS",
		}
	}

	return ValidationError{}
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
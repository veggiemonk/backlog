package validation

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

// ValidationError represents a validation error with context
type ValidationError struct {
	Field   string `json:"field"`
	Value   string `json:"value"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("validation error in field '%s': %s (value: '%s')", v.Field, v.Message, v.Value)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "no validation errors"
	}
	if len(ve) == 1 {
		return ve[0].Error()
	}

	var msgs []string
	for _, err := range ve {
		msgs = append(msgs, err.Error())
	}
	return fmt.Sprintf("multiple validation errors: %s", strings.Join(msgs, "; "))
}

// HasErrors returns true if there are validation errors
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

// Validator provides validation functions for CLI inputs
type Validator struct{}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// Validation constants
const (
	MaxTitleLength        = 200
	MaxDescriptionLength  = 5000
	MaxLabelLength        = 50
	MaxAssigneeLength     = 100
	MaxDependencyIDLength = 20
	MaxACLength           = 500
	MaxPlanLength         = 10000
	MaxNotesLength        = 10000
	MaxPathLength         = 4096
	MaxLogLevelLength     = 10
	MaxLogFormatLength    = 10
)

// Regular expressions for validation
var (
	taskIDRegex    = regexp.MustCompile(`^T?(\d+(?:\.\d+)*)$`)
	labelRegex     = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	assigneeRegex  = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	logLevelRegex  = regexp.MustCompile(`^(debug|info|warn|error)$`)
	logFormatRegex = regexp.MustCompile(`^(json|text)$`)
	priorityRegex  = regexp.MustCompile(`^(low|medium|high|critical)$`)
)

// ValidateTitle validates task title
func (v *Validator) ValidateTitle(title string) ValidationError {
	if strings.TrimSpace(title) == "" {
		return ValidationError{
			Field:   "title",
			Value:   title,
			Message: "title cannot be empty",
			Code:    "TITLE_EMPTY",
		}
	}

	if utf8.RuneCountInString(title) > MaxTitleLength {
		return ValidationError{
			Field:   "title",
			Value:   title,
			Message: fmt.Sprintf("title cannot exceed %d characters", MaxTitleLength),
			Code:    "TITLE_TOO_LONG",
		}
	}

	// Check for potentially dangerous characters
	if containsDangerousChars(title) {
		return ValidationError{
			Field:   "title",
			Value:   title,
			Message: "title contains potentially dangerous characters",
			Code:    "TITLE_INVALID_CHARS",
		}
	}

	return ValidationError{}
}

// ValidateDescription validates task description
func (v *Validator) ValidateDescription(description string) ValidationError {
	if utf8.RuneCountInString(description) > MaxDescriptionLength {
		return ValidationError{
			Field:   "description",
			Value:   description,
			Message: fmt.Sprintf("description cannot exceed %d characters", MaxDescriptionLength),
			Code:    "DESCRIPTION_TOO_LONG",
		}
	}

	if containsDangerousChars(description) {
		return ValidationError{
			Field:   "description",
			Value:   description,
			Message: "description contains potentially dangerous characters",
			Code:    "DESCRIPTION_INVALID_CHARS",
		}
	}

	return ValidationError{}
}

// ValidateTaskID validates task ID format
func (v *Validator) ValidateTaskID(taskID string) ValidationError {
	if taskID == "" {
		return ValidationError{
			Field:   "task_id",
			Value:   taskID,
			Message: "task ID cannot be empty",
			Code:    "TASK_ID_EMPTY",
		}
	}

	if len(taskID) > MaxDependencyIDLength {
		return ValidationError{
			Field:   "task_id",
			Value:   taskID,
			Message: fmt.Sprintf("task ID cannot exceed %d characters", MaxDependencyIDLength),
			Code:    "TASK_ID_TOO_LONG",
		}
	}

	if !taskIDRegex.MatchString(taskID) {
		return ValidationError{
			Field:   "task_id",
			Value:   taskID,
			Message: "task ID must follow format: T1, T1.1, 1, 1.1, etc.",
			Code:    "TASK_ID_INVALID_FORMAT",
		}
	}

	return ValidationError{}
}

// ValidateLabels validates a slice of labels
func (v *Validator) ValidateLabels(labels []string) ValidationErrors {
	var errors ValidationErrors

	for i, label := range labels {
		if err := v.ValidateLabel(label); err.Code != "" {
			err.Field = fmt.Sprintf("labels[%d]", i)
			errors = append(errors, err)
		}
	}

	return errors
}

// ValidateLabel validates a single label
func (v *Validator) ValidateLabel(label string) ValidationError {
	if strings.TrimSpace(label) == "" {
		return ValidationError{
			Field:   "label",
			Value:   label,
			Message: "label cannot be empty",
			Code:    "LABEL_EMPTY",
		}
	}

	if utf8.RuneCountInString(label) > MaxLabelLength {
		return ValidationError{
			Field:   "label",
			Value:   label,
			Message: fmt.Sprintf("label cannot exceed %d characters", MaxLabelLength),
			Code:    "LABEL_TOO_LONG",
		}
	}

	if !labelRegex.MatchString(label) {
		return ValidationError{
			Field:   "label",
			Value:   label,
			Message: "label can only contain alphanumeric characters, hyphens, and underscores",
			Code:    "LABEL_INVALID_FORMAT",
		}
	}

	return ValidationError{}
}

// ValidateAssignees validates a slice of assignees
func (v *Validator) ValidateAssignees(assignees []string) ValidationErrors {
	var errors ValidationErrors

	for i, assignee := range assignees {
		if err := v.ValidateAssignee(assignee); err.Code != "" {
			err.Field = fmt.Sprintf("assignees[%d]", i)
			errors = append(errors, err)
		}
	}

	return errors
}

// ValidateAssignee validates a single assignee
func (v *Validator) ValidateAssignee(assignee string) ValidationError {
	if strings.TrimSpace(assignee) == "" {
		return ValidationError{
			Field:   "assignee",
			Value:   assignee,
			Message: "assignee cannot be empty",
			Code:    "ASSIGNEE_EMPTY",
		}
	}

	if utf8.RuneCountInString(assignee) > MaxAssigneeLength {
		return ValidationError{
			Field:   "assignee",
			Value:   assignee,
			Message: fmt.Sprintf("assignee cannot exceed %d characters", MaxAssigneeLength),
			Code:    "ASSIGNEE_TOO_LONG",
		}
	}

	if !assigneeRegex.MatchString(assignee) {
		return ValidationError{
			Field:   "assignee",
			Value:   assignee,
			Message: "assignee can only contain alphanumeric characters, dots, hyphens, and underscores",
			Code:    "ASSIGNEE_INVALID_FORMAT",
		}
	}

	return ValidationError{}
}

// ValidateDependencies validates a slice of dependency IDs
func (v *Validator) ValidateDependencies(dependencies []string) ValidationErrors {
	var errors ValidationErrors

	for i, dep := range dependencies {
		if err := v.ValidateTaskID(dep); err.Code != "" {
			err.Field = fmt.Sprintf("dependencies[%d]", i)
			errors = append(errors, err)
		}
	}

	return errors
}

// ValidateAcceptanceCriteria validates a slice of acceptance criteria
func (v *Validator) ValidateAcceptanceCriteria(ac []string) ValidationErrors {
	var errors ValidationErrors

	for i, criterion := range ac {
		if err := v.ValidateAcceptanceCriterion(criterion); err.Code != "" {
			err.Field = fmt.Sprintf("acceptance_criteria[%d]", i)
			errors = append(errors, err)
		}
	}

	return errors
}

// ValidateAcceptanceCriterion validates a single acceptance criterion
func (v *Validator) ValidateAcceptanceCriterion(ac string) ValidationError {
	if strings.TrimSpace(ac) == "" {
		return ValidationError{
			Field:   "acceptance_criterion",
			Value:   ac,
			Message: "acceptance criterion cannot be empty",
			Code:    "AC_EMPTY",
		}
	}

	if utf8.RuneCountInString(ac) > MaxACLength {
		return ValidationError{
			Field:   "acceptance_criterion",
			Value:   ac,
			Message: fmt.Sprintf("acceptance criterion cannot exceed %d characters", MaxACLength),
			Code:    "AC_TOO_LONG",
		}
	}

	if containsDangerousChars(ac) {
		return ValidationError{
			Field:   "acceptance_criterion",
			Value:   ac,
			Message: "acceptance criterion contains potentially dangerous characters",
			Code:    "AC_INVALID_CHARS",
		}
	}

	return ValidationError{}
}

// ValidatePriority validates priority value
func (v *Validator) ValidatePriority(priority string) ValidationError {
	if priority == "" {
		return ValidationError{}
	}

	if !priorityRegex.MatchString(priority) {
		return ValidationError{
			Field:   "priority",
			Value:   priority,
			Message: "priority must be one of: low, medium, high, critical",
			Code:    "PRIORITY_INVALID",
		}
	}

	return ValidationError{}
}

// ValidatePlan validates implementation plan
func (v *Validator) ValidatePlan(plan string) ValidationError {
	if utf8.RuneCountInString(plan) > MaxPlanLength {
		return ValidationError{
			Field:   "plan",
			Value:   plan,
			Message: fmt.Sprintf("plan cannot exceed %d characters", MaxPlanLength),
			Code:    "PLAN_TOO_LONG",
		}
	}

	if containsDangerousChars(plan) {
		return ValidationError{
			Field:   "plan",
			Value:   plan,
			Message: "plan contains potentially dangerous characters",
			Code:    "PLAN_INVALID_CHARS",
		}
	}

	return ValidationError{}
}

// ValidateNotes validates implementation notes
func (v *Validator) ValidateNotes(notes string) ValidationError {
	if utf8.RuneCountInString(notes) > MaxNotesLength {
		return ValidationError{
			Field:   "notes",
			Value:   notes,
			Message: fmt.Sprintf("notes cannot exceed %d characters", MaxNotesLength),
			Code:    "NOTES_TOO_LONG",
		}
	}

	if containsDangerousChars(notes) {
		return ValidationError{
			Field:   "notes",
			Value:   notes,
			Message: "notes contains potentially dangerous characters",
			Code:    "NOTES_INVALID_CHARS",
		}
	}

	return ValidationError{}
}

// ValidateFilePath validates file paths for security
func (v *Validator) ValidateFilePath(path string) ValidationError {
	if path == "" {
		return ValidationError{
			Field:   "file_path",
			Value:   path,
			Message: "file path cannot be empty",
			Code:    "PATH_EMPTY",
		}
	}

	if len(path) > MaxPathLength {
		return ValidationError{
			Field:   "file_path",
			Value:   path,
			Message: fmt.Sprintf("file path cannot exceed %d characters", MaxPathLength),
			Code:    "PATH_TOO_LONG",
		}
	}

	// Check for path traversal attempts
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return ValidationError{
			Field:   "file_path",
			Value:   path,
			Message: "file path contains path traversal sequences",
			Code:    "PATH_TRAVERSAL",
		}
	}

	// Check for null bytes
	if strings.Contains(path, "\x00") {
		return ValidationError{
			Field:   "file_path",
			Value:   path,
			Message: "file path contains null bytes",
			Code:    "PATH_NULL_BYTE",
		}
	}

	return ValidationError{}
}

// ValidateLogLevel validates log level
func (v *Validator) ValidateLogLevel(level string) ValidationError {
	if level == "" {
		return ValidationError{}
	}

	if !logLevelRegex.MatchString(level) {
		return ValidationError{
			Field:   "log_level",
			Value:   level,
			Message: "log level must be one of: debug, info, warn, error",
			Code:    "LOG_LEVEL_INVALID",
		}
	}

	return ValidationError{}
}

// ValidateLogFormat validates log format
func (v *Validator) ValidateLogFormat(format string) ValidationError {
	if format == "" {
		return ValidationError{}
	}

	if !logFormatRegex.MatchString(format) {
		return ValidationError{
			Field:   "log_format",
			Value:   format,
			Message: "log format must be one of: json, text",
			Code:    "LOG_FORMAT_INVALID",
		}
	}

	return ValidationError{}
}

// containsDangerousChars checks for potentially dangerous characters
func containsDangerousChars(input string) bool {
	dangerous := []string{
		"\x00",    // null byte
		"<script", // script tags
		"</script>",
		"javascript:",
		"data:",
		"vbscript:",
		"onload=",
		"onerror=",
		"eval(",
		"expression(",
		"\\x", // hex escape sequences
		"\\u", // unicode escape sequences
	}

	lowerInput := strings.ToLower(input)
	for _, danger := range dangerous {
		if strings.Contains(lowerInput, danger) {
			return true
		}
	}

	return false
}


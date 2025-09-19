package validation

import (
	"html"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Sanitizer provides comprehensive input sanitization functions
type Sanitizer struct{}

// NewSanitizer creates a new sanitizer instance
func NewSanitizer() *Sanitizer {
	return &Sanitizer{}
}

// SanitizeText removes potentially dangerous content from text input
func (s *Sanitizer) SanitizeText(input string) string {
	if input == "" {
		return input
	}

	original := input

	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Remove control characters except common whitespace
	var builder strings.Builder
	for _, r := range input {
		if !unicode.IsControl(r) || r == '\n' || r == '\r' || r == '\t' {
			builder.WriteRune(r)
		}
	}
	input = builder.String()

	// HTML escape to prevent XSS
	input = html.EscapeString(input)

	// Normalize whitespace
	input = normalizeWhitespace(input)

	// Log sanitization if content was modified
	if original != input {
		LogSanitizationAlert(original, input, "text sanitization performed")
	}

	return input
}

// SanitizeTitle sanitizes task titles
func (s *Sanitizer) SanitizeTitle(title string) string {
	if title == "" {
		return title
	}

	// Basic text sanitization
	title = s.SanitizeText(title)

	// Remove line breaks from titles
	title = strings.ReplaceAll(title, "\n", " ")
	title = strings.ReplaceAll(title, "\r", " ")
	title = strings.ReplaceAll(title, "\t", " ")

	// Normalize and trim whitespace
	title = normalizeWhitespace(title)
	title = strings.TrimSpace(title)

	return title
}

// SanitizeDescription sanitizes task descriptions
func (s *Sanitizer) SanitizeDescription(description string) string {
	if description == "" {
		return description
	}

	// Basic text sanitization
	description = s.SanitizeText(description)

	// Allow line breaks but normalize them
	description = normalizeLineBreaks(description)

	return description
}

// SanitizeLabel sanitizes labels for consistent formatting
func (s *Sanitizer) SanitizeLabel(label string) string {
	if label == "" {
		return label
	}

	// Remove dangerous characters and normalize
	label = s.SanitizeText(label)

	// Remove spaces and convert to lowercase for consistency
	label = strings.ToLower(label)
	label = strings.ReplaceAll(label, " ", "-")

	// Keep only alphanumeric, hyphens, and underscores
	reg := regexp.MustCompile(`[^a-z0-9_-]`)
	label = reg.ReplaceAllString(label, "")

	// Remove multiple consecutive hyphens/underscores
	reg = regexp.MustCompile(`[-_]+`)
	label = reg.ReplaceAllString(label, "-")

	// Trim leading/trailing hyphens and underscores
	label = strings.Trim(label, "-_")

	return label
}

// SanitizeAssignee sanitizes assignee names
func (s *Sanitizer) SanitizeAssignee(assignee string) string {
	if assignee == "" {
		return assignee
	}

	// Basic text sanitization
	assignee = s.SanitizeText(assignee)

	// Remove line breaks
	assignee = strings.ReplaceAll(assignee, "\n", " ")
	assignee = strings.ReplaceAll(assignee, "\r", " ")
	assignee = strings.ReplaceAll(assignee, "\t", " ")

	// Normalize whitespace and trim
	assignee = normalizeWhitespace(assignee)
	assignee = strings.TrimSpace(assignee)

	// Keep only allowed characters for assignees
	reg := regexp.MustCompile(`[^a-zA-Z0-9\.\-\ ]`)
	assignee = reg.ReplaceAllString(assignee, "")

	return assignee
}

// SanitizeTaskID sanitizes task IDs
func (s *Sanitizer) SanitizeTaskID(taskID string) string {
	if taskID == "" {
		return taskID
	}

	// Remove dangerous characters
	taskID = strings.ReplaceAll(taskID, "\x00", "")

	// Remove control characters
	var builder strings.Builder
	for _, r := range taskID {
		if !unicode.IsControl(r) {
			builder.WriteRune(r)
		}
	}
	taskID = builder.String()

	// Trim whitespace
	taskID = strings.TrimSpace(taskID)

	// Ensure proper format: only digits, dots, and optional T prefix
	reg := regexp.MustCompile(`[^T0-9.]`)
	taskID = reg.ReplaceAllString(taskID, "")

	return taskID
}

// SanitizeFilePath sanitizes file paths to prevent directory traversal
func (s *Sanitizer) SanitizeFilePath(path string) string {
	if path == "" {
		return path
	}

	// Remove null bytes
	path = strings.ReplaceAll(path, "\x00", "")

	// Remove control characters
	var builder strings.Builder
	for _, r := range path {
		if !unicode.IsControl(r) {
			builder.WriteRune(r)
		}
	}
	path = builder.String()

	// Remove dangerous sequences
	path = strings.ReplaceAll(path, "..", "")

	// Remove multiple consecutive slashes
	reg := regexp.MustCompile(`/+`)
	path = reg.ReplaceAllString(path, "/")

	// Trim whitespace
	path = strings.TrimSpace(path)

	return path
}

// SanitizeSearchQuery sanitizes search queries
func (s *Sanitizer) SanitizeSearchQuery(query string) string {
	if query == "" {
		return query
	}

	// Basic text sanitization
	query = s.SanitizeText(query)

	// Normalize whitespace but preserve structure
	query = normalizeWhitespace(query)
	query = strings.TrimSpace(query)

	return query
}

// SanitizeAcceptanceCriterion sanitizes acceptance criteria
func (s *Sanitizer) SanitizeAcceptanceCriterion(ac string) string {
	if ac == "" {
		return ac
	}

	// Basic text sanitization
	ac = s.SanitizeText(ac)

	// Allow line breaks but normalize them
	ac = normalizeLineBreaks(ac)

	return ac
}

// SanitizePlan sanitizes implementation plans
func (s *Sanitizer) SanitizePlan(plan string) string {
	if plan == "" {
		return plan
	}

	// Basic text sanitization
	plan = s.SanitizeText(plan)

	// Allow line breaks but normalize them
	plan = normalizeLineBreaks(plan)

	return plan
}

// SanitizeNotes sanitizes implementation notes
func (s *Sanitizer) SanitizeNotes(notes string) string {
	if notes == "" {
		return notes
	}

	// Basic text sanitization
	notes = s.SanitizeText(notes)

	// Allow line breaks but normalize them
	notes = normalizeLineBreaks(notes)

	return notes
}

// RemoveUnsafeCharacters removes characters that could be used for injection attacks
func (s *Sanitizer) RemoveUnsafeCharacters(input string) string {
	if input == "" {
		return input
	}

	// Define patterns for dangerous sequences
	dangerousPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`), // Script tags
		regexp.MustCompile(`(?i)javascript:[^"\s]*`),        // JavaScript protocol
		regexp.MustCompile(`(?i)vbscript:[^"\s]*`),          // VBScript protocol
		regexp.MustCompile(`(?i)data:[^"\s]*`),              // Data protocol
		regexp.MustCompile(`(?i)on\w+\s*=[^"\s]*`),          // Event handlers
		regexp.MustCompile(`(?i)expression\s*\([^)]*\)`),    // CSS expressions
		regexp.MustCompile(`(?i)eval\s*\([^)]*\)`),          // Eval function
		regexp.MustCompile(`\\x[0-9a-fA-F]{2}`),             // Hex escapes
		regexp.MustCompile(`\\u[0-9a-fA-F]{4}`),             // Unicode escapes
		regexp.MustCompile(`\x00`),                          // Null bytes
	}

	result := input
	for _, pattern := range dangerousPatterns {
		result = pattern.ReplaceAllString(result, "")
	}

	return result
}

// TruncateString safely truncates a string to a maximum length while preserving UTF-8 encoding
func (s *Sanitizer) TruncateString(input string, maxLength int) string {
	if input == "" || maxLength <= 0 {
		return ""
	}

	if utf8.RuneCountInString(input) <= maxLength {
		return input
	}

	runes := []rune(input)
	if len(runes) <= maxLength {
		return input
	}

	return string(runes[:maxLength])
}

// Helper functions

// normalizeWhitespace replaces multiple consecutive whitespace characters with single spaces
func normalizeWhitespace(input string) string {
	reg := regexp.MustCompile(`\s+`)
	return reg.ReplaceAllString(input, " ")
}

// normalizeLineBreaks normalizes different line break styles to Unix style
func normalizeLineBreaks(input string) string {
	// Convert Windows line endings to Unix
	input = strings.ReplaceAll(input, "\r\n", "\n")
	// Convert old Mac line endings to Unix
	input = strings.ReplaceAll(input, "\r", "\n")
	// Remove excessive consecutive line breaks
	reg := regexp.MustCompile(`\n{3,}`)
	input = reg.ReplaceAllString(input, "\n\n")
	return input
}

// IsValidUTF8 checks if a string contains valid UTF-8 sequences
func (s *Sanitizer) IsValidUTF8(input string) bool {
	return utf8.ValidString(input)
}

// SanitizeSlice sanitizes a slice of strings using the provided sanitization function
func (s *Sanitizer) SanitizeSlice(input []string, sanitizeFunc func(string) string) []string {
	if len(input) == 0 {
		return input
	}

	result := make([]string, 0, len(input))
	for _, item := range input {
		sanitized := sanitizeFunc(item)
		if sanitized != "" { // Skip empty strings after sanitization
			result = append(result, sanitized)
		}
	}

	return result
}


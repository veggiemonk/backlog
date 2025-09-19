package validation

import (
	"fmt"
	"maps"
	"runtime"
	"time"

	"github.com/veggiemonk/backlog/internal/logging"
)

// SecurityMonitor provides security monitoring and alerting capabilities
type SecurityMonitor struct {
	enabled bool
}

// NewSecurityMonitor creates a new security monitor instance
func NewSecurityMonitor() *SecurityMonitor {
	return &SecurityMonitor{
		enabled: true,
	}
}

// SecurityEvent represents a security-related event
type SecurityEvent struct {
	Type        string            `json:"type"`
	Severity    string            `json:"severity"`
	Message     string            `json:"message"`
	Details     map[string]any    `json:"details"`
	Timestamp   time.Time         `json:"timestamp"`
	UserContext map[string]string `json:"user_context,omitempty"`
	StackTrace  string            `json:"stack_trace,omitempty"`
}

// Event types
const (
	EventValidationFailure   = "validation_failure"
	EventSanitizationAlert   = "sanitization_alert"
	EventFileAccessViolation = "file_access_violation"
	EventSuspiciousActivity  = "suspicious_activity"
	EventConfigurationError  = "configuration_error"
	EventAuthenticationIssue = "authentication_issue"
	EventCommandExecution    = "command_execution"
	EventDataExfiltration    = "data_exfiltration"
	EventIntegrityViolation  = "integrity_violation"
	EventRateLimitExceeded   = "rate_limit_exceeded"
)

// Severity levels
const (
	SeverityLow      = "low"
	SeverityMedium   = "medium"
	SeverityHigh     = "high"
	SeverityCritical = "critical"
)

// LogValidationFailure logs validation failures with security context
func (sm *SecurityMonitor) LogValidationFailure(field, value, reason string, details map[string]any) {
	if !sm.enabled {
		return
	}

	event := SecurityEvent{
		Type:      EventValidationFailure,
		Severity:  SeverityMedium,
		Message:   fmt.Sprintf("Validation failed for field '%s': %s", field, reason),
		Timestamp: time.Now().UTC(),
		Details: map[string]any{
			"field":  field,
			"value":  value,
			"reason": reason,
		},
	}

	// Add additional details if provided
	maps.Copy(event.Details, details)

	sm.logSecurityEvent(event)
}

// LogSanitizationAlert logs when potentially dangerous content is sanitized
func (sm *SecurityMonitor) LogSanitizationAlert(input, sanitized, reason string) {
	if !sm.enabled {
		return
	}

	severity := SeverityLow
	if len(input) != len(sanitized) {
		severity = SeverityMedium
	}

	event := SecurityEvent{
		Type:      EventSanitizationAlert,
		Severity:  severity,
		Message:   fmt.Sprintf("Input sanitization performed: %s", reason),
		Timestamp: time.Now().UTC(),
		Details: map[string]any{
			"original_input":   input,
			"sanitized_output": sanitized,
			"reason":           reason,
			"chars_removed":    len(input) - len(sanitized),
		},
	}

	sm.logSecurityEvent(event)
}

// LogFileAccessViolation logs unauthorized file access attempts
func (sm *SecurityMonitor) LogFileAccessViolation(operation, path, reason string) {
	if !sm.enabled {
		return
	}

	event := SecurityEvent{
		Type:      EventFileAccessViolation,
		Severity:  SeverityHigh,
		Message:   fmt.Sprintf("File access violation: %s on '%s' - %s", operation, path, reason),
		Timestamp: time.Now().UTC(),
		Details: map[string]any{
			"operation": operation,
			"path":      path,
			"reason":    reason,
		},
		StackTrace: sm.getStackTrace(),
	}

	sm.logSecurityEvent(event)
}

// LogSuspiciousActivity logs suspicious activities that may indicate an attack
func (sm *SecurityMonitor) LogSuspiciousActivity(activity string, details map[string]any) {
	if !sm.enabled {
		return
	}

	event := SecurityEvent{
		Type:       EventSuspiciousActivity,
		Severity:   SeverityHigh,
		Message:    fmt.Sprintf("Suspicious activity detected: %s", activity),
		Timestamp:  time.Now().UTC(),
		Details:    details,
		StackTrace: sm.getStackTrace(),
	}

	sm.logSecurityEvent(event)
}

// LogConfigurationError logs security-related configuration errors
func (sm *SecurityMonitor) LogConfigurationError(config, value, reason string) {
	if !sm.enabled {
		return
	}

	event := SecurityEvent{
		Type:      EventConfigurationError,
		Severity:  SeverityMedium,
		Message:   fmt.Sprintf("Configuration error: %s", reason),
		Timestamp: time.Now().UTC(),
		Details: map[string]any{
			"configuration": config,
			"value":         value,
			"reason":        reason,
		},
	}

	sm.logSecurityEvent(event)
}

// LogCommandExecution logs command executions for audit trail
func (sm *SecurityMonitor) LogCommandExecution(command string, args []string, success bool) {
	if !sm.enabled {
		return
	}

	severity := SeverityLow
	if !success {
		severity = SeverityMedium
	}

	event := SecurityEvent{
		Type:      EventCommandExecution,
		Severity:  severity,
		Message:   fmt.Sprintf("Command executed: %s", command),
		Timestamp: time.Now().UTC(),
		Details: map[string]any{
			"command":   command,
			"arguments": args,
			"success":   success,
		},
	}

	sm.logSecurityEvent(event)
}

// LogIntegrityViolation logs file or data integrity violations
func (sm *SecurityMonitor) LogIntegrityViolation(resource, violation string, details map[string]any) {
	if !sm.enabled {
		return
	}

	event := SecurityEvent{
		Type:       EventIntegrityViolation,
		Severity:   SeverityCritical,
		Message:    fmt.Sprintf("Integrity violation detected on %s: %s", resource, violation),
		Timestamp:  time.Now().UTC(),
		Details:    details,
		StackTrace: sm.getStackTrace(),
	}

	sm.logSecurityEvent(event)
}

// LogRateLimitExceeded logs when rate limits are exceeded
func (sm *SecurityMonitor) LogRateLimitExceeded(operation, identifier string, limit int, attempts int) {
	if !sm.enabled {
		return
	}

	event := SecurityEvent{
		Type:      EventRateLimitExceeded,
		Severity:  SeverityMedium,
		Message:   fmt.Sprintf("Rate limit exceeded for %s: %d attempts (limit: %d)", operation, attempts, limit),
		Timestamp: time.Now().UTC(),
		Details: map[string]any{
			"operation":  operation,
			"identifier": identifier,
			"limit":      limit,
			"attempts":   attempts,
		},
	}

	sm.logSecurityEvent(event)
}

// SetUserContext sets user context for security events
func (sm *SecurityMonitor) SetUserContext(context map[string]string) {
	// In a real implementation, this would store user context
	// for inclusion in security events
}

// Enable enables security monitoring
func (sm *SecurityMonitor) Enable() {
	sm.enabled = true
	logging.Info("Security monitoring enabled")
}

// Disable disables security monitoring
func (sm *SecurityMonitor) Disable() {
	sm.enabled = false
	logging.Info("Security monitoring disabled")
}

// IsEnabled returns whether security monitoring is enabled
func (sm *SecurityMonitor) IsEnabled() bool {
	return sm.enabled
}

// logSecurityEvent logs a security event using the logging system
func (sm *SecurityMonitor) logSecurityEvent(event SecurityEvent) {
	// Log at appropriate level based on severity
	switch event.Severity {
	case SeverityLow:
		logging.Info("security_event",
			"type", event.Type,
			"severity", event.Severity,
			"message", event.Message,
			"details", event.Details,
			"timestamp", event.Timestamp,
		)
	case SeverityMedium:
		logging.Warn("security_event",
			"type", event.Type,
			"severity", event.Severity,
			"message", event.Message,
			"details", event.Details,
			"timestamp", event.Timestamp,
		)
	case SeverityHigh, SeverityCritical:
		logging.Error("security_event",
			"type", event.Type,
			"severity", event.Severity,
			"message", event.Message,
			"details", event.Details,
			"timestamp", event.Timestamp,
			"stack_trace", event.StackTrace,
		)
	}
}

// getStackTrace returns a stack trace for debugging
func (sm *SecurityMonitor) getStackTrace() string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// Global security monitor instance
var globalSecurityMonitor = NewSecurityMonitor()

// GetSecurityMonitor returns the global security monitor instance
func GetSecurityMonitor() *SecurityMonitor {
	return globalSecurityMonitor
}

// Convenience functions for global security monitor

// LogValidationFailure logs validation failures globally
func LogValidationFailure(field, value, reason string, details map[string]any) {
	globalSecurityMonitor.LogValidationFailure(field, value, reason, details)
}

// LogSanitizationAlert logs sanitization alerts globally
func LogSanitizationAlert(input, sanitized, reason string) {
	globalSecurityMonitor.LogSanitizationAlert(input, sanitized, reason)
}

// LogFileAccessViolation logs file access violations globally
func LogFileAccessViolation(operation, path, reason string) {
	globalSecurityMonitor.LogFileAccessViolation(operation, path, reason)
}

// LogSuspiciousActivity logs suspicious activities globally
func LogSuspiciousActivity(activity string, details map[string]any) {
	globalSecurityMonitor.LogSuspiciousActivity(activity, details)
}

// LogConfigurationError logs configuration errors globally
func LogConfigurationError(config, value, reason string) {
	globalSecurityMonitor.LogConfigurationError(config, value, reason)
}

// LogCommandExecution logs command executions globally
func LogCommandExecution(command string, args []string, success bool) {
	globalSecurityMonitor.LogCommandExecution(command, args, success)
}

// LogIntegrityViolation logs integrity violations globally
func LogIntegrityViolation(resource, violation string, details map[string]any) {
	globalSecurityMonitor.LogIntegrityViolation(resource, violation, details)
}

// LogRateLimitExceeded logs rate limit violations globally
func LogRateLimitExceeded(operation, identifier string, limit int, attempts int) {
	globalSecurityMonitor.LogRateLimitExceeded(operation, identifier, limit, attempts)
}

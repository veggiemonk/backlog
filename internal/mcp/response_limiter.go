package mcp

import (
	"encoding/json"
	"time"

	"github.com/veggiemonk/backlog/internal/logging"
)

// ResponseMetrics holds metrics about response sizes
type ResponseMetrics struct {
	TotalResponses    int64
	OversizedCount    int64
	AverageSize       float64
	MaxSize           int
	LastOversizedTime time.Time
}

// ResponseSizeMonitor monitors and tracks response sizes
type ResponseSizeMonitor struct {
	metrics ResponseMetrics
	config  ResponseSizeConfig
}

// NewResponseSizeMonitor creates a new response size monitor
func NewResponseSizeMonitor(config ResponseSizeConfig) *ResponseSizeMonitor {
	return &ResponseSizeMonitor{
		config: config,
	}
}

// MonitorResponse tracks a response and logs if it's oversized
func (m *ResponseSizeMonitor) MonitorResponse(data []byte, operation string) {
	responseSize := len(data)
	estimatedTokens := int(float64(responseSize) * m.config.TokensPerByte)

	// Update metrics
	m.metrics.TotalResponses++

	// Update running average
	oldAvg := m.metrics.AverageSize
	m.metrics.AverageSize = oldAvg + (float64(responseSize)-oldAvg)/float64(m.metrics.TotalResponses)

	if responseSize > m.metrics.MaxSize {
		m.metrics.MaxSize = responseSize
	}

	// Check if oversized
	if estimatedTokens > m.config.TokenLimit {
		m.metrics.OversizedCount++
		m.metrics.LastOversizedTime = time.Now()

		logging.Warn("Oversized MCP response detected",
			"operation", operation,
			"response_size_bytes", responseSize,
			"estimated_tokens", estimatedTokens,
			"token_limit", m.config.TokenLimit,
			"oversized_ratio", float64(estimatedTokens)/float64(m.config.TokenLimit))
	}

	// Log periodic metrics (every 100 responses)
	if m.metrics.TotalResponses%100 == 0 {
		m.logMetrics()
	}
}

// logMetrics logs current response metrics
func (m *ResponseSizeMonitor) logMetrics() {
	oversizedPercentage := float64(m.metrics.OversizedCount) / float64(m.metrics.TotalResponses) * 100

	logging.Info("MCP response size metrics",
		"total_responses", m.metrics.TotalResponses,
		"oversized_count", m.metrics.OversizedCount,
		"oversized_percentage", oversizedPercentage,
		"average_size_bytes", int(m.metrics.AverageSize),
		"max_size_bytes", m.metrics.MaxSize)
}

// GetMetrics returns current metrics
func (m *ResponseSizeMonitor) GetMetrics() ResponseMetrics {
	return m.metrics
}

// ResponseSizeMiddleware wraps MCP tool handlers to monitor response sizes
type ResponseSizeMiddleware struct {
	monitor *ResponseSizeMonitor
}

// NewResponseSizeMiddleware creates new middleware
func NewResponseSizeMiddleware(config ResponseSizeConfig) *ResponseSizeMiddleware {
	return &ResponseSizeMiddleware{
		monitor: NewResponseSizeMonitor(config),
	}
}

// WrapResponse wraps the response and monitors its size
func (m *ResponseSizeMiddleware) WrapResponse(response interface{}, operation string) interface{} {
	// Serialize to measure actual size
	data, err := json.Marshal(response)
	if err != nil {
		logging.Error("Failed to marshal response for monitoring", "error", err, "operation", operation)
		return response
	}

	// Monitor the response
	m.monitor.MonitorResponse(data, operation)

	return response
}

// GetMonitor returns the underlying monitor for metrics access
func (m *ResponseSizeMiddleware) GetMonitor() *ResponseSizeMonitor {
	return m.monitor
}

// ResponseSizeLimiter provides hard limits on response sizes
type ResponseSizeLimiter struct {
	config ResponseSizeConfig
}

// NewResponseSizeLimiter creates a new response size limiter
func NewResponseSizeLimiter(config ResponseSizeConfig) *ResponseSizeLimiter {
	return &ResponseSizeLimiter{config: config}
}

// CheckAndLimitResponse checks if response exceeds limits and truncates if necessary
func (l *ResponseSizeLimiter) CheckAndLimitResponse(response interface{}, operation string) (interface{}, bool) {
	data, err := json.Marshal(response)
	if err != nil {
		logging.Error("Failed to marshal response for size checking", "error", err, "operation", operation)
		return response, false
	}

	responseSize := len(data)
	estimatedTokens := int(float64(responseSize) * l.config.TokensPerByte)

	if estimatedTokens <= l.config.TokenLimit {
		return response, false // Not truncated
	}

	// Response is too large, create a truncated version
	truncatedResponse := map[string]interface{}{
		"error":           "Response too large",
		"estimated_tokens": estimatedTokens,
		"token_limit":     l.config.TokenLimit,
		"operation":       operation,
		"message":         "The response was truncated due to size limits. Please use pagination or filtering to reduce the response size.",
	}

	logging.Warn("Response truncated due to size limits",
		"operation", operation,
		"original_size_bytes", responseSize,
		"estimated_tokens", estimatedTokens,
		"token_limit", l.config.TokenLimit)

	return truncatedResponse, true // Truncated
}

// ResponseSizeStats provides aggregate statistics
type ResponseSizeStats struct {
	TotalRequests      int64     `json:"total_requests"`
	OversizedRequests  int64     `json:"oversized_requests"`
	OversizedPercent   float64   `json:"oversized_percent"`
	AverageSizeBytes   int       `json:"average_size_bytes"`
	MaxSizeBytes       int       `json:"max_size_bytes"`
	LastOversizedTime  time.Time `json:"last_oversized_time,omitempty"`
}

// GetStats returns formatted statistics
func (m *ResponseSizeMonitor) GetStats() ResponseSizeStats {
	metrics := m.GetMetrics()

	var oversizedPercent float64
	if metrics.TotalResponses > 0 {
		oversizedPercent = float64(metrics.OversizedCount) / float64(metrics.TotalResponses) * 100
	}

	return ResponseSizeStats{
		TotalRequests:     metrics.TotalResponses,
		OversizedRequests: metrics.OversizedCount,
		OversizedPercent:  oversizedPercent,
		AverageSizeBytes:  int(metrics.AverageSize),
		MaxSizeBytes:      metrics.MaxSize,
		LastOversizedTime: metrics.LastOversizedTime,
	}
}
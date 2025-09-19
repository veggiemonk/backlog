package mcp

import (
	"encoding/json"

	"github.com/veggiemonk/backlog/internal/core"
)

const (
	// DefaultTokenLimit is the default MCP token limit (25,000 tokens)
	DefaultTokenLimit = 25000
	// SafetyMarginPercent is the safety margin percentage to apply to size estimates
	SafetyMarginPercent = 20
	// EstimatedTokensPerByte is a rough estimate of tokens per byte in JSON
	EstimatedTokensPerByte = 0.75
)

// ResponseSizeConfig holds configuration for response size estimation
type ResponseSizeConfig struct {
	TokenLimit    int
	SafetyMargin  float64
	TokensPerByte float64
}

// DefaultResponseSizeConfig returns the default configuration for response size estimation
func DefaultResponseSizeConfig() ResponseSizeConfig {
	return ResponseSizeConfig{
		TokenLimit:    DefaultTokenLimit,
		SafetyMargin:  SafetyMarginPercent / 100.0,
		TokensPerByte: EstimatedTokensPerByte,
	}
}

// EstimateResponseSize estimates the size of a JSON response for a list of tasks
func EstimateResponseSize(tasks []*core.Task) int {
	return EstimateResponseSizeWithConfig(tasks, DefaultResponseSizeConfig())
}

// EstimateResponseSizeWithConfig estimates response size using custom configuration
func EstimateResponseSizeWithConfig(tasks []*core.Task, config ResponseSizeConfig) int {
	if len(tasks) == 0 {
		// Empty response with wrapper
		return int(float64(50) * (1 + config.SafetyMargin))
	}

	// Sample a few tasks to get average size
	sampleSize := min(len(tasks), 5)
	totalSampleSize := 0

	for i := range sampleSize {
		taskSize := estimateTaskJSONSize(tasks[i])
		totalSampleSize += taskSize
	}

	avgTaskSize := totalSampleSize / sampleSize

	// Estimate total size: average task size * number of tasks + JSON wrapper overhead
	wrapperOverhead := 100 // {"Tasks": [...]} + some padding
	estimatedBytes := avgTaskSize*len(tasks) + wrapperOverhead

	// Apply safety margin
	estimatedBytesWithMargin := int(float64(estimatedBytes) * (1 + config.SafetyMargin))

	// Convert to estimated token count
	estimatedTokens := int(float64(estimatedBytesWithMargin) * config.TokensPerByte)

	return estimatedTokens
}

// estimateTaskJSONSize estimates the JSON size of a single task
func estimateTaskJSONSize(task *core.Task) int {
	// Quick estimation based on field lengths rather than full marshaling
	baseSize := 200 // Basic JSON structure overhead

	// Add estimated sizes for major fields
	baseSize += len(task.ID.String()) * 2     // ID appears in multiple places
	baseSize += len(task.Title) * 2           // Title and filename
	baseSize += len(task.Description)         // Description
	baseSize += len(task.ImplementationPlan)  // Implementation plan
	baseSize += len(task.ImplementationNotes) // Implementation notes

	// Add sizes for arrays
	for _, ac := range task.AcceptanceCriteria {
		baseSize += len(ac.Text) + 50 // AC text + JSON overhead
	}

	for _, label := range task.Labels {
		baseSize += len(label) + 10 // Label + JSON overhead
	}

	for _, assigned := range task.Assigned {
		baseSize += len(assigned) + 10 // Assigned + JSON overhead
	}

	for _, dep := range task.Dependencies {
		baseSize += len(dep) + 10 // Dependency + JSON overhead
	}

	// Add overhead for timestamps and other fields
	baseSize += 200

	return baseSize
}

// WillExceedLimit checks if the estimated response size will exceed the token limit
func WillExceedLimit(tasks []*core.Task) bool {
	return WillExceedLimitWithConfig(tasks, DefaultResponseSizeConfig())
}

// WillExceedLimitWithConfig checks if response will exceed limit using custom config
func WillExceedLimitWithConfig(tasks []*core.Task, config ResponseSizeConfig) bool {
	estimatedTokens := EstimateResponseSizeWithConfig(tasks, config)
	return estimatedTokens > config.TokenLimit
}

// CalculateOptimalChunkSize calculates the optimal number of tasks per chunk
// to stay within token limits
func CalculateOptimalChunkSize(tasks []*core.Task) int {
	return CalculateOptimalChunkSizeWithConfig(tasks, DefaultResponseSizeConfig())
}

// CalculateOptimalChunkSizeWithConfig calculates optimal chunk size with custom config
func CalculateOptimalChunkSizeWithConfig(tasks []*core.Task, config ResponseSizeConfig) int {
	if len(tasks) == 0 {
		return 0
	}

	// If total size is within limits, return all tasks
	if !WillExceedLimitWithConfig(tasks, config) {
		return len(tasks)
	}

	// Binary search for optimal chunk size
	low, high := 1, len(tasks)
	optimalSize := 1

	for low <= high {
		mid := (low + high) / 2

		// Test with a sample chunk of this size
		testChunk := tasks[:min(mid, len(tasks))]

		if !WillExceedLimitWithConfig(testChunk, config) {
			optimalSize = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return optimalSize
}

// AccurateEstimateResponseSize provides a more accurate estimate by actually marshaling a sample
func AccurateEstimateResponseSize(tasks []*core.Task) (int, error) {
	if len(tasks) == 0 {
		return 50, nil
	}

	// Use a sample for accuracy vs performance trade-off
	sampleSize := min(len(tasks), 3)
	sample := tasks[:sampleSize]

	wrappedSample := struct{ Tasks []*core.Task }{Tasks: sample}
	data, err := json.Marshal(wrappedSample)
	if err != nil {
		return 0, err
	}

	// Calculate average task size and extrapolate
	avgTaskSize := len(data) / sampleSize
	wrapperOverhead := 100

	config := DefaultResponseSizeConfig()
	estimatedBytes := avgTaskSize*len(tasks) + wrapperOverhead
	estimatedBytesWithMargin := int(float64(estimatedBytes) * (1 + config.SafetyMargin))
	estimatedTokens := int(float64(estimatedBytesWithMargin) * config.TokensPerByte)

	return estimatedTokens, nil
}

// Helper function for min since Go doesn't have a built-in min for ints
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}


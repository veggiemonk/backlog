package mcp

import (
	"time"
)

// PaginationConfig holds configuration for pagination behavior
type PaginationConfig struct {
	// DefaultPageSize is the default number of items per page when no limit is specified
	DefaultPageSize int `json:"default_page_size"`

	// MaxPageSize is the maximum allowed page size
	MaxPageSize int `json:"max_page_size"`

	// AutoPaginationThreshold is the number of items above which auto-pagination kicks in
	AutoPaginationThreshold int `json:"auto_pagination_threshold"`

	// EnableAutoPagination enables automatic pagination suggestions
	EnableAutoPagination bool `json:"enable_auto_pagination"`

	// ResponseSizeConfig embedded for size-based pagination decisions
	ResponseSizeConfig ResponseSizeConfig `json:"response_size_config"`
}

// DefaultPaginationConfig returns the default pagination configuration
func DefaultPaginationConfig() PaginationConfig {
	return PaginationConfig{
		DefaultPageSize:         50,
		MaxPageSize:             200,
		AutoPaginationThreshold: 100,
		EnableAutoPagination:    true,
		ResponseSizeConfig:      DefaultResponseSizeConfig(),
	}
}

// CustomPaginationConfig creates a custom pagination configuration
func CustomPaginationConfig(defaultPageSize, maxPageSize, autoThreshold int, enableAuto bool) PaginationConfig {
	return PaginationConfig{
		DefaultPageSize:         defaultPageSize,
		MaxPageSize:             maxPageSize,
		AutoPaginationThreshold: autoThreshold,
		EnableAutoPagination:    enableAuto,
		ResponseSizeConfig:      DefaultResponseSizeConfig(),
	}
}

// ValidatePageSize ensures the requested page size is within configured limits
func (c *PaginationConfig) ValidatePageSize(requestedSize int) int {
	if requestedSize <= 0 {
		return c.DefaultPageSize
	}
	if requestedSize > c.MaxPageSize {
		return c.MaxPageSize
	}
	return requestedSize
}

// ShouldAutoPaginate determines if auto-pagination should be triggered
func (c *PaginationConfig) ShouldAutoPaginate(totalItems int) bool {
	return c.EnableAutoPagination && totalItems > c.AutoPaginationThreshold
}

// CalculateOptimalPageSize calculates the optimal page size based on content and limits
func (c *PaginationConfig) CalculateOptimalPageSize(totalItems int, avgItemSize int) int {
	// Start with default page size
	pageSize := c.DefaultPageSize

	// Adjust based on average item size and response size limits
	maxBytesPerResponse := int(float64(c.ResponseSizeConfig.TokenLimit) / c.ResponseSizeConfig.TokensPerByte)
	maxItemsBasedOnSize := maxBytesPerResponse / avgItemSize

	// Use the more restrictive limit
	if maxItemsBasedOnSize < pageSize {
		pageSize = maxItemsBasedOnSize
	}

	// Ensure we don't exceed max page size
	if pageSize > c.MaxPageSize {
		pageSize = c.MaxPageSize
	}

	// Ensure minimum page size of 1
	if pageSize < 1 {
		pageSize = 1
	}

	return pageSize
}

// AdvancedPaginationMetadata provides comprehensive pagination information
type AdvancedPaginationMetadata struct {
	// Basic pagination info
	Offset    int  `json:"offset"`
	Limit     int  `json:"limit"`
	Total     int  `json:"total"`
	HasMore   bool `json:"has_more"`
	NextPage  *int `json:"next_page,omitempty"`
	PrevPage  *int `json:"prev_page,omitempty"`

	// Advanced metadata
	TotalPages    int `json:"total_pages"`
	CurrentPage   int `json:"current_page"`
	ItemsInPage   int `json:"items_in_page"`

	// Performance and sizing info
	EstimatedResponseSize int           `json:"estimated_response_size_tokens,omitempty"`
	ResponseGenerationTime time.Duration `json:"response_generation_time,omitempty"`

	// Optimization suggestions
	SuggestedPageSize   *int   `json:"suggested_page_size,omitempty"`
	OptimizationMessage string `json:"optimization_message,omitempty"`
}

// CreateAdvancedPaginationMetadata creates comprehensive pagination metadata
func CreateAdvancedPaginationMetadata(
	offset, limit, total, itemsInPage int,
	config PaginationConfig,
	estimatedSize int,
	generationTime time.Duration,
) *AdvancedPaginationMetadata {

	// Calculate page numbers (handle zero limit case)
	currentPage := 1
	totalPages := 1
	if limit > 0 {
		currentPage = (offset / limit) + 1
		totalPages = (total + limit - 1) / limit // Ceiling division
	}

	// Calculate next/prev pages
	var nextPage, prevPage *int
	if offset+limit < total {
		next := offset + limit
		nextPage = &next
	}
	if offset > 0 {
		prev := offset - limit
		if prev < 0 {
			prev = 0
		}
		prevPage = &prev
	}

	metadata := &AdvancedPaginationMetadata{
		Offset:          offset,
		Limit:           limit,
		Total:           total,
		HasMore:         offset+limit < total,
		NextPage:        nextPage,
		PrevPage:        prevPage,
		TotalPages:      totalPages,
		CurrentPage:     currentPage,
		ItemsInPage:     itemsInPage,
		EstimatedResponseSize: estimatedSize,
		ResponseGenerationTime: generationTime,
	}

	// Add optimization suggestions
	if estimatedSize > config.ResponseSizeConfig.TokenLimit && itemsInPage > 0 {
		avgItemSize := estimatedSize / itemsInPage
		suggestedSize := config.CalculateOptimalPageSize(total, avgItemSize)
		metadata.SuggestedPageSize = &suggestedSize
		metadata.OptimizationMessage = "Response size exceeds recommended limits. Consider using a smaller page size for better performance."
	}

	return metadata
}

// PaginationStrategy represents different pagination strategies
type PaginationStrategy string

const (
	// OffsetBasedPagination uses offset and limit
	OffsetBasedPagination PaginationStrategy = "offset"

	// CursorBasedPagination uses cursor-based pagination (for future implementation)
	CursorBasedPagination PaginationStrategy = "cursor"

	// TokenBasedPagination optimizes based on response token limits
	TokenBasedPagination PaginationStrategy = "token"
)

// PaginationRequest represents a pagination request
type PaginationRequest struct {
	Strategy    PaginationStrategy `json:"strategy,omitempty"`
	Offset      *int              `json:"offset,omitempty"`
	Limit       *int              `json:"limit,omitempty"`
	Cursor      *string           `json:"cursor,omitempty"`
	MaxTokens   *int              `json:"max_tokens,omitempty"`
}

// NormalizePaginationRequest normalizes and validates a pagination request
func (c *PaginationConfig) NormalizePaginationRequest(req PaginationRequest) PaginationRequest {
	normalized := req

	// Set default strategy
	if normalized.Strategy == "" {
		normalized.Strategy = OffsetBasedPagination
	}

	// Normalize offset
	if normalized.Offset == nil {
		zero := 0
		normalized.Offset = &zero
	} else if *normalized.Offset < 0 {
		zero := 0
		normalized.Offset = &zero
	}

	// Normalize limit
	if normalized.Limit == nil {
		normalized.Limit = &c.DefaultPageSize
	} else {
		validatedLimit := c.ValidatePageSize(*normalized.Limit)
		normalized.Limit = &validatedLimit
	}

	return normalized
}
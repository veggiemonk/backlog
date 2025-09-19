package mcp

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestPaginationConfig(t *testing.T) {
	is := is.New(t)

	t.Run("DefaultConfig", func(t *testing.T) {
		is := is.New(t)

		config := DefaultPaginationConfig()
		is.Equal(config.DefaultPageSize, 50)
		is.Equal(config.MaxPageSize, 200)
		is.Equal(config.AutoPaginationThreshold, 100)
		is.True(config.EnableAutoPagination)
	})

	t.Run("CustomConfig", func(t *testing.T) {
		is := is.New(t)

		config := CustomPaginationConfig(25, 100, 50, false)
		is.Equal(config.DefaultPageSize, 25)
		is.Equal(config.MaxPageSize, 100)
		is.Equal(config.AutoPaginationThreshold, 50)
		is.True(!config.EnableAutoPagination)
	})

	t.Run("ValidatePageSize", func(t *testing.T) {
		is := is.New(t)

		config := DefaultPaginationConfig()

		// Test default when requested size is 0 or negative
		is.Equal(config.ValidatePageSize(0), config.DefaultPageSize)
		is.Equal(config.ValidatePageSize(-10), config.DefaultPageSize)

		// Test valid sizes
		is.Equal(config.ValidatePageSize(25), 25)
		is.Equal(config.ValidatePageSize(100), 100)

		// Test size exceeding maximum
		is.Equal(config.ValidatePageSize(300), config.MaxPageSize)
	})

	t.Run("ShouldAutoPaginate", func(t *testing.T) {
		is := is.New(t)

		config := DefaultPaginationConfig()

		// Should not auto-paginate for small datasets
		is.True(!config.ShouldAutoPaginate(50))
		is.True(!config.ShouldAutoPaginate(100))

		// Should auto-paginate for large datasets
		is.True(config.ShouldAutoPaginate(150))
		is.True(config.ShouldAutoPaginate(1000))

		// Should not auto-paginate when disabled
		config.EnableAutoPagination = false
		is.True(!config.ShouldAutoPaginate(1000))
	})

	t.Run("CalculateOptimalPageSize", func(t *testing.T) {
		is := is.New(t)

		config := DefaultPaginationConfig()

		// Test with small items (should use default page size)
		pageSize := config.CalculateOptimalPageSize(100, 50)
		is.Equal(pageSize, config.DefaultPageSize)

		// Test with large items (should reduce page size)
		largeItemSize := 10000
		pageSize = config.CalculateOptimalPageSize(100, largeItemSize)
		is.True(pageSize < config.DefaultPageSize)
		is.True(pageSize >= 1)

		// Test with very large items (should use minimum page size)
		veryLargeItemSize := 100000
		pageSize = config.CalculateOptimalPageSize(100, veryLargeItemSize)
		is.True(pageSize >= 1)
	})
}

func TestAdvancedPaginationMetadata(t *testing.T) {
	is := is.New(t)

	config := DefaultPaginationConfig()

	t.Run("BasicMetadata", func(t *testing.T) {
		is := is.New(t)

		metadata := CreateAdvancedPaginationMetadata(
			0,    // offset
			10,   // limit
			50,   // total
			10,   // items in page
			config,
			1000, // estimated size
			100*time.Millisecond,
		)

		is.Equal(metadata.Offset, 0)
		is.Equal(metadata.Limit, 10)
		is.Equal(metadata.Total, 50)
		is.Equal(metadata.ItemsInPage, 10)
		is.Equal(metadata.TotalPages, 5)
		is.Equal(metadata.CurrentPage, 1)
		is.True(metadata.HasMore)
		is.True(metadata.NextPage != nil)
		is.Equal(*metadata.NextPage, 10)
		is.True(metadata.PrevPage == nil)
	})

	t.Run("MiddlePage", func(t *testing.T) {
		is := is.New(t)

		metadata := CreateAdvancedPaginationMetadata(
			20,   // offset (page 3)
			10,   // limit
			50,   // total
			10,   // items in page
			config,
			1000, // estimated size
			100*time.Millisecond,
		)

		is.Equal(metadata.CurrentPage, 3)
		is.True(metadata.HasMore)
		is.True(metadata.NextPage != nil)
		is.Equal(*metadata.NextPage, 30)
		is.True(metadata.PrevPage != nil)
		is.Equal(*metadata.PrevPage, 10)
	})

	t.Run("LastPage", func(t *testing.T) {
		is := is.New(t)

		metadata := CreateAdvancedPaginationMetadata(
			40,   // offset (page 5, last page)
			10,   // limit
			50,   // total
			10,   // items in page
			config,
			1000, // estimated size
			100*time.Millisecond,
		)

		is.Equal(metadata.CurrentPage, 5)
		is.True(!metadata.HasMore)
		is.True(metadata.NextPage == nil)
		is.True(metadata.PrevPage != nil)
		is.Equal(*metadata.PrevPage, 30)
	})

	t.Run("OptimizationSuggestions", func(t *testing.T) {
		is := is.New(t)

		// Create metadata with oversized response (much larger than limit)
		oversizedTokens := config.ResponseSizeConfig.TokenLimit * 3 // 3x the limit
		metadata := CreateAdvancedPaginationMetadata(
			0,    // offset
			10,   // limit
			50,   // total
			10,   // items in page
			config,
			oversizedTokens, // estimated size exceeds limit significantly
			100*time.Millisecond,
		)

		is.True(metadata.SuggestedPageSize != nil)
		is.True(*metadata.SuggestedPageSize <= 10) // Should be <= original limit
		is.True(len(metadata.OptimizationMessage) > 0)
	})
}

func TestPaginationRequest(t *testing.T) {
	is := is.New(t)

	config := DefaultPaginationConfig()

	t.Run("NormalizeBasicRequest", func(t *testing.T) {
		is := is.New(t)

		req := PaginationRequest{}
		normalized := config.NormalizePaginationRequest(req)

		is.Equal(normalized.Strategy, OffsetBasedPagination)
		is.True(normalized.Offset != nil)
		is.Equal(*normalized.Offset, 0)
		is.True(normalized.Limit != nil)
		is.Equal(*normalized.Limit, config.DefaultPageSize)
	})

	t.Run("NormalizeNegativeOffset", func(t *testing.T) {
		is := is.New(t)

		offset := -10
		req := PaginationRequest{Offset: &offset}
		normalized := config.NormalizePaginationRequest(req)

		is.Equal(*normalized.Offset, 0)
	})

	t.Run("NormalizeExcessiveLimit", func(t *testing.T) {
		is := is.New(t)

		limit := 500
		req := PaginationRequest{Limit: &limit}
		normalized := config.NormalizePaginationRequest(req)

		is.Equal(*normalized.Limit, config.MaxPageSize)
	})
}

func TestPaginationStrategies(t *testing.T) {
	is := is.New(t)

	t.Run("StrategyConstants", func(t *testing.T) {
		is := is.New(t)

		is.Equal(string(OffsetBasedPagination), "offset")
		is.Equal(string(CursorBasedPagination), "cursor")
		is.Equal(string(TokenBasedPagination), "token")
	})
}
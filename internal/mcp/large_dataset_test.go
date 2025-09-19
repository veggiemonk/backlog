package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/matryer/is"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/afero"
	"github.com/veggiemonk/backlog/internal/core"
)

func TestLargeDatasetScenarios(t *testing.T) {
	tests := []struct {
		name      string
		taskCount int
	}{
		{"100_tasks", 100},
		{"500_tasks", 500},
		{"1000_tasks", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			// Create in-memory filesystem and store
			fs := afero.NewMemMapFs()
			store := core.NewFileTaskStore(fs, ".backlog")

			// Create handler with middleware
			responseSizeConfig := DefaultResponseSizeConfig()
			middleware := NewResponseSizeMiddleware(responseSizeConfig)
			handler := &handler{
				store:      store,
				mu:         &sync.Mutex{},
				middleware: middleware,
			}

			// Create large dataset
			createLargeDataset(t, store, tt.taskCount)

			// Test without pagination (should suggest pagination for large datasets)
			ctx := context.Background()
			req := &mcp.CallToolRequest{}
			params := core.ListTasksParams{}

			result, _, err := handler.list(ctx, req, params)

			// For large datasets, should get structured error response with pagination suggestion
			if tt.taskCount >= 100 {
				// Should return structured error instead of Go error
				is.NoErr(err) // No Go error
				is.True(result != nil) // Should return structured error response
				is.True(result.IsError) // Should be marked as error response

				txtContent, ok := result.Content[0].(*mcp.TextContent)
				is.True(ok)
				// Should contain structured error with pagination suggestion
				is.True(len(txtContent.Text) > 0)
			} else {
				// For smaller datasets, should return successful response
				is.NoErr(err)
				is.True(result != nil)
				is.True(!result.IsError)
			}
		})
	}
}

func TestPaginationFunctionality(t *testing.T) {
	is := is.New(t)

	// Setup
	fs := afero.NewMemMapFs()
	store := core.NewFileTaskStore(fs, ".backlog")
	responseSizeConfig := DefaultResponseSizeConfig()
	middleware := NewResponseSizeMiddleware(responseSizeConfig)
	handler := &handler{
		store:      store,
		mu:         &sync.Mutex{},
		middleware: middleware,
	}

	// Create 50 tasks for pagination testing
	taskCount := 50
	createLargeDataset(t, store, taskCount)

	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	t.Run("pagination_with_limit", func(t *testing.T) {
		is := is.New(t)

		limit := 10
		params := core.ListTasksParams{
			Limit: &limit,
		}

		result, _, err := handler.list(ctx, req, params)
		is.NoErr(err)
		is.True(result != nil)

		// Parse response
		txtContent, ok := result.Content[0].(*mcp.TextContent)
		is.True(ok)

		var response TaskListResponse
		err = json.Unmarshal([]byte(txtContent.Text), &response)
		is.NoErr(err)

		// Should have exactly the requested limit
		is.Equal(len(response.Tasks), limit)

		// Should have pagination metadata
		is.True(response.Pagination != nil)
		is.Equal(response.Pagination.Limit, limit)
		is.Equal(response.Pagination.Offset, 0)
		is.True(response.Pagination.HasMore)
	})

	t.Run("pagination_with_offset_and_limit", func(t *testing.T) {
		is := is.New(t)

		limit := 10
		offset := 20
		params := core.ListTasksParams{
			Limit:  &limit,
			Offset: &offset,
		}

		result, _, err := handler.list(ctx, req, params)
		is.NoErr(err)
		is.True(result != nil)

		// Parse response
		txtContent, ok := result.Content[0].(*mcp.TextContent)
		is.True(ok)

		var response TaskListResponse
		err = json.Unmarshal([]byte(txtContent.Text), &response)
		is.NoErr(err)

		// Should have exactly the requested limit
		is.Equal(len(response.Tasks), limit)

		// Should have correct pagination metadata
		is.True(response.Pagination != nil)
		is.Equal(response.Pagination.Limit, limit)
		is.Equal(response.Pagination.Offset, offset)
		is.True(response.Pagination.HasMore)
	})

	t.Run("pagination_last_page", func(t *testing.T) {
		is := is.New(t)

		limit := 10
		offset := 45 // Should get last 5 tasks
		params := core.ListTasksParams{
			Limit:  &limit,
			Offset: &offset,
		}

		result, _, err := handler.list(ctx, req, params)
		is.NoErr(err)
		is.True(result != nil)

		// Parse response
		txtContent, ok := result.Content[0].(*mcp.TextContent)
		is.True(ok)

		var response TaskListResponse
		err = json.Unmarshal([]byte(txtContent.Text), &response)
		is.NoErr(err)

		// Should have remaining tasks
		is.Equal(len(response.Tasks), 5) // 50 total - 45 offset = 5 remaining

		// Should have correct pagination metadata
		is.True(response.Pagination != nil)
		is.Equal(response.Pagination.Offset, offset)
		is.True(!response.Pagination.HasMore) // No more pages
		is.True(response.Pagination.NextPage == nil)
	})
}

func TestResponseSizeEstimation(t *testing.T) {
	is := is.New(t)

	// Setup
	fs := afero.NewMemMapFs()
	store := core.NewFileTaskStore(fs, ".backlog")

	// Create tasks with varying content sizes
	tasks := createVariedSizeTasks(t, store, 20)

	t.Run("estimate_response_size", func(t *testing.T) {
		is := is.New(t)

		// Test estimation accuracy
		estimatedSize := EstimateResponseSize(tasks)
		is.True(estimatedSize > 0)

		// Test with custom config
		config := ResponseSizeConfig{
			TokenLimit:    10000,
			SafetyMargin:  0.1,
			TokensPerByte: 0.8,
		}
		estimatedSizeCustom := EstimateResponseSizeWithConfig(tasks, config)
		is.True(estimatedSizeCustom > 0)
		is.True(estimatedSizeCustom != estimatedSize) // Should be different with different config
	})

	t.Run("will_exceed_limit", func(t *testing.T) {
		is := is.New(t)

		// Small dataset should not exceed default limits
		smallTasks := tasks[:5]
		is.True(!WillExceedLimit(smallTasks))

		// Test with very low limit
		lowLimitConfig := ResponseSizeConfig{
			TokenLimit:    100, // Very low limit
			SafetyMargin:  0.1,
			TokensPerByte: 0.75,
		}
		is.True(WillExceedLimitWithConfig(tasks, lowLimitConfig))
	})

	t.Run("calculate_optimal_chunk_size", func(t *testing.T) {
		is := is.New(t)

		optimalSize := CalculateOptimalChunkSize(tasks)
		is.True(optimalSize > 0)
		is.True(optimalSize <= len(tasks))

		// Test with low limit config
		lowLimitConfig := ResponseSizeConfig{
			TokenLimit:    1000,
			SafetyMargin:  0.1,
			TokensPerByte: 0.75,
		}
		optimalSizeLowLimit := CalculateOptimalChunkSizeWithConfig(tasks, lowLimitConfig)
		is.True(optimalSizeLowLimit > 0)
		is.True(optimalSizeLowLimit < optimalSize) // Should be smaller with lower limit
	})
}

func TestResponseSizeMonitoring(t *testing.T) {
	is := is.New(t)

	config := DefaultResponseSizeConfig()
	monitor := NewResponseSizeMonitor(config)

	t.Run("monitor_response", func(t *testing.T) {
		is := is.New(t)

		// Create test data
		testData := []byte(`{"tasks": [{"id": "1", "title": "Test Task"}]}`)

		// Monitor response
		monitor.MonitorResponse(testData, "test_operation")

		// Check metrics
		metrics := monitor.GetMetrics()
		is.Equal(metrics.TotalResponses, int64(1))
		is.True(metrics.AverageSize > 0)
		is.Equal(metrics.MaxSize, len(testData))
	})

	t.Run("monitor_multiple_responses", func(t *testing.T) {
		is := is.New(t)

		// Monitor multiple responses
		for i := 0; i < 10; i++ {
			testData := []byte(fmt.Sprintf(`{"tasks": [{"id": "%d", "title": "Test Task %d"}]}`, i, i))
			monitor.MonitorResponse(testData, "test_operation")
		}

		// Check metrics
		metrics := monitor.GetMetrics()
		is.True(metrics.TotalResponses >= 10) // >= because of previous test
		is.True(metrics.AverageSize > 0)
	})

	t.Run("get_stats", func(t *testing.T) {
		is := is.New(t)

		stats := monitor.GetStats()
		is.True(stats.TotalRequests > 0)
		is.True(stats.AverageSizeBytes > 0)
		is.True(stats.OversizedPercent >= 0)
	})
}

func TestResponseSizeMiddleware(t *testing.T) {
	is := is.New(t)

	config := DefaultResponseSizeConfig()
	middleware := NewResponseSizeMiddleware(config)

	t.Run("wrap_response", func(t *testing.T) {
		is := is.New(t)

		testResponse := map[string]interface{}{
			"tasks": []string{"task1", "task2"},
		}

		wrappedResponse := middleware.WrapResponse(testResponse, "test_operation")
		is.True(wrappedResponse != nil)

		// Response should be unchanged (middleware only monitors)
		is.Equal(wrappedResponse, testResponse)
	})

	t.Run("get_monitor", func(t *testing.T) {
		is := is.New(t)

		monitor := middleware.GetMonitor()
		is.True(monitor != nil)

		stats := monitor.GetStats()
		is.True(stats.TotalRequests >= 0)
	})
}

// Helper function to create a large dataset of tasks
func createLargeDataset(t *testing.T, store TaskStore, count int) {
	is := is.New(t)

	for i := 0; i < count; i++ {
		params := core.CreateTaskParams{
			Title:       fmt.Sprintf("Task %d", i+1),
			Description: fmt.Sprintf("Description for task %d with some additional content to make it realistic", i+1),
			Priority:    "medium",
			Labels:      []string{"test", fmt.Sprintf("batch-%d", i/10)},
		}

		if i%3 == 0 {
			params.Assigned = []string{fmt.Sprintf("user%d", i%5)}
		}

		_, err := store.Create(params)
		is.NoErr(err)
	}
}

// Helper function to create tasks with varied content sizes
func createVariedSizeTasks(t *testing.T, store TaskStore, count int) []*core.Task {
	is := is.New(t)
	var tasks []*core.Task

	for i := 0; i < count; i++ {
		// Vary the content size
		descriptionSize := 50 + (i * 20) // Increasing description length
		description := fmt.Sprintf("Description for task %d: %s",
			i+1,
			generateRandomText(descriptionSize))

		params := core.CreateTaskParams{
			Title:       fmt.Sprintf("Varied Task %d", i+1),
			Description: description,
			Priority:    "medium",
			Labels:      []string{"test", "varied-size"},
		}

		task, err := store.Create(params)
		is.NoErr(err)
		tasks = append(tasks, task)
	}

	return tasks
}

// Helper function to generate random text of specified length
func generateRandomText(length int) string {
	text := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. "
	for len(text) < length {
		text += text
	}
	return text[:length]
}
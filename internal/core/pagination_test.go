package core

import (
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/afero"
)

func TestPaginateTasks(t *testing.T) {
	tests := []struct {
		name          string
		tasks         []Task
		limit         *int
		offset        *int
		expectedCount int
		expectedFirst string // ID of first task in result
		expectedLast  string // ID of last task in result
	}{
		{
			name: "no_pagination",
			tasks: []Task{
				{ID: TaskID{seg: []int{1}}, Title: "Task 1"},
				{ID: TaskID{seg: []int{2}}, Title: "Task 2"},
				{ID: TaskID{seg: []int{3}}, Title: "Task 3"},
			},
			limit:         nil,
			offset:        nil,
			expectedCount: 3,
			expectedFirst: "01",
			expectedLast:  "03",
		},
		{
			name: "limit_only",
			tasks: []Task{
				{ID: TaskID{seg: []int{1}}, Title: "Task 1"},
				{ID: TaskID{seg: []int{2}}, Title: "Task 2"},
				{ID: TaskID{seg: []int{3}}, Title: "Task 3"},
			},
			limit:         intPtr(2),
			offset:        nil,
			expectedCount: 2,
			expectedFirst: "01",
			expectedLast:  "02",
		},
		{
			name: "offset_only",
			tasks: []Task{
				{ID: TaskID{seg: []int{1}}, Title: "Task 1"},
				{ID: TaskID{seg: []int{2}}, Title: "Task 2"},
				{ID: TaskID{seg: []int{3}}, Title: "Task 3"},
			},
			limit:         nil,
			offset:        intPtr(1),
			expectedCount: 2,
			expectedFirst: "02",
			expectedLast:  "03",
		},
		{
			name: "limit_and_offset",
			tasks: []Task{
				{ID: TaskID{seg: []int{1}}, Title: "Task 1"},
				{ID: TaskID{seg: []int{2}}, Title: "Task 2"},
				{ID: TaskID{seg: []int{3}}, Title: "Task 3"},
				{ID: TaskID{seg: []int{4}}, Title: "Task 4"},
				{ID: TaskID{seg: []int{5}}, Title: "Task 5"},
			},
			limit:         intPtr(2),
			offset:        intPtr(1),
			expectedCount: 2,
			expectedFirst: "02",
			expectedLast:  "03",
		},
		{
			name: "offset_beyond_end",
			tasks: []Task{
				{ID: TaskID{seg: []int{1}}, Title: "Task 1"},
				{ID: TaskID{seg: []int{2}}, Title: "Task 2"},
			},
			limit:         intPtr(2),
			offset:        intPtr(5),
			expectedCount: 0,
		},
		{
			name: "limit_larger_than_remaining",
			tasks: []Task{
				{ID: TaskID{seg: []int{1}}, Title: "Task 1"},
				{ID: TaskID{seg: []int{2}}, Title: "Task 2"},
				{ID: TaskID{seg: []int{3}}, Title: "Task 3"},
			},
			limit:         intPtr(5),
			offset:        intPtr(1),
			expectedCount: 2,
			expectedFirst: "02",
			expectedLast:  "03",
		},
		{
			name:          "empty_task_list",
			tasks:         []Task{},
			limit:         intPtr(5),
			offset:        intPtr(0),
			expectedCount: 0,
		},
		{
			name: "zero_limit",
			tasks: []Task{
				{ID: TaskID{seg: []int{1}}, Title: "Task 1"},
				{ID: TaskID{seg: []int{2}}, Title: "Task 2"},
			},
			limit:         intPtr(0),
			offset:        nil,
			expectedCount: 2,
			expectedFirst: "01",
			expectedLast:  "02",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			result := PaginateTasks(tt.tasks, tt.limit, tt.offset)
			is.Equal(len(result), tt.expectedCount)

			if tt.expectedCount > 0 {
				is.Equal(result[0].ID.String(), tt.expectedFirst)
				is.Equal(result[len(result)-1].ID.String(), tt.expectedLast)
			}
		})
	}
}

func TestListWithPagination(t *testing.T) {
	is := is.New(t)

	// Create a memory filesystem and task store
	fs := afero.NewMemMapFs()
	store := NewFileTaskStore(fs, "tasks")

	// Create test tasks
	testTasks := []CreateTaskParams{
		{Title: "Task A", Priority: "high"},
		{Title: "Task B", Priority: "medium"},
		{Title: "Task C", Priority: "low"},
		{Title: "Task D", Priority: "high"},
		{Title: "Task E", Priority: "critical"},
	}

	// Create the tasks
	for _, params := range testTasks {
		_, err := store.Create(params)
		is.NoErr(err)
	}

	t.Run("list_with_limit", func(t *testing.T) {
		is := is.New(t)
		limit := 3
		params := ListTasksParams{
			Limit: &limit,
		}

		tasks, err := store.List(params)
		is.NoErr(err)
		is.Equal(len(tasks), 3)
	})

	t.Run("list_with_offset", func(t *testing.T) {
		is := is.New(t)
		offset := 2
		params := ListTasksParams{
			Offset: &offset,
		}

		tasks, err := store.List(params)
		is.NoErr(err)
		is.Equal(len(tasks), 3) // Should have 3 remaining tasks (5 - 2 offset)
	})

	t.Run("list_with_limit_and_offset", func(t *testing.T) {
		is := is.New(t)
		limit := 2
		offset := 1
		params := ListTasksParams{
			Limit:  &limit,
			Offset: &offset,
		}

		tasks, err := store.List(params)
		is.NoErr(err)
		is.Equal(len(tasks), 2)
	})

	t.Run("list_with_filter_and_pagination", func(t *testing.T) {
		is := is.New(t)
		limit := 2
		params := ListTasksParams{
			Status: []string{"todo"},
			Limit:  &limit,
		}

		tasks, err := store.List(params)
		is.NoErr(err)
		is.Equal(len(tasks), 2)

		// All should be "todo" (default status)
		for _, task := range tasks {
			is.Equal(string(task.Status), "todo")
		}
	})
}

func TestSearchWithPagination(t *testing.T) {
	is := is.New(t)

	// Create a memory filesystem and task store
	fs := afero.NewMemMapFs()
	store := NewFileTaskStore(fs, "tasks")

	// Create test tasks with searchable content
	testTasks := []CreateTaskParams{
		{Title: "Feature A", Description: "A new feature implementation"},
		{Title: "Feature B", Description: "Another feature to build"},
		{Title: "Bug Fix A", Description: "Fix for critical bug"},
		{Title: "Feature C", Description: "Yet another feature"},
		{Title: "Bug Fix B", Description: "Another bug fix"},
	}

	// Create the tasks
	for _, params := range testTasks {
		_, err := store.Create(params)
		is.NoErr(err)
	}

	t.Run("search_with_limit", func(t *testing.T) {
		is := is.New(t)
		limit := 2
		params := ListTasksParams{
			Limit: &limit,
		}

		tasks, err := store.Search("feature", params)
		is.NoErr(err)
		is.Equal(len(tasks), 2) // Should limit results to 2

		// All should contain "feature"
		for _, task := range tasks {
			containsFeature := contains(task.Title, "Feature") || contains(task.Description, "feature")
			is.True(containsFeature)
		}
	})

	t.Run("search_with_offset", func(t *testing.T) {
		is := is.New(t)
		offset := 1
		params := ListTasksParams{
			Offset: &offset,
		}

		tasks, err := store.Search("feature", params)
		is.NoErr(err)
		is.Equal(len(tasks), 2) // Should have 2 remaining after offset (3 - 1)
	})

	t.Run("search_with_limit_and_offset", func(t *testing.T) {
		is := is.New(t)
		limit := 1
		offset := 1
		params := ListTasksParams{
			Limit:  &limit,
			Offset: &offset,
		}

		tasks, err := store.Search("feature", params)
		is.NoErr(err)
		is.Equal(len(tasks), 1) // Should have 1 task (skip first, take one)
	})
}

func TestPaginationInfoCalculation(t *testing.T) {
	tests := []struct {
		name              string
		totalTasks        int
		limit             *int
		offset            *int
		expectedTotal     int
		expectedDisplayed int
		expectedOffset    int
		expectedLimit     int
		expectedHasMore   bool
	}{
		{
			name:              "no_pagination",
			totalTasks:        5,
			limit:             nil,
			offset:            nil,
			expectedTotal:     5,
			expectedDisplayed: 5,
			expectedOffset:    0,
			expectedLimit:     0,
			expectedHasMore:   false,
		},
		{
			name:              "with_limit_has_more",
			totalTasks:        10,
			limit:             intPtr(3),
			offset:            nil,
			expectedTotal:     10,
			expectedDisplayed: 3,
			expectedOffset:    0,
			expectedLimit:     3,
			expectedHasMore:   true,
		},
		{
			name:              "with_limit_no_more",
			totalTasks:        3,
			limit:             intPtr(5),
			offset:            nil,
			expectedTotal:     3,
			expectedDisplayed: 3,
			expectedOffset:    0,
			expectedLimit:     5,
			expectedHasMore:   false,
		},
		{
			name:              "with_offset_and_limit",
			totalTasks:        10,
			limit:             intPtr(3),
			offset:            intPtr(2),
			expectedTotal:     10,
			expectedDisplayed: 3,
			expectedOffset:    2,
			expectedLimit:     3,
			expectedHasMore:   true,
		},
		{
			name:              "last_page",
			totalTasks:        10,
			limit:             intPtr(3),
			offset:            intPtr(9),
			expectedTotal:     10,
			expectedDisplayed: 1,
			expectedOffset:    9,
			expectedLimit:     3,
			expectedHasMore:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			// Create mock tasks
			allTasks := make([]Task, tt.totalTasks)
			for i := 0; i < tt.totalTasks; i++ {
				allTasks[i] = Task{
					ID:    TaskID{seg: []int{i + 1}},
					Title: "Task",
				}
			}

			// Apply pagination
			paginatedTasks := PaginateTasks(allTasks, tt.limit, tt.offset)

			// Calculate pagination info
			offsetVal := 0
			if tt.offset != nil {
				offsetVal = *tt.offset
			}
			limitVal := 0
			if tt.limit != nil {
				limitVal = *tt.limit
			}
			hasMore := (offsetVal + len(paginatedTasks)) < tt.totalTasks

			paginationInfo := &PaginationInfo{
				TotalResults:     tt.totalTasks,
				DisplayedResults: len(paginatedTasks),
				Offset:           offsetVal,
				Limit:            limitVal,
				HasMore:          hasMore,
			}

			// Verify calculations
			is.Equal(paginationInfo.TotalResults, tt.expectedTotal)
			is.Equal(paginationInfo.DisplayedResults, tt.expectedDisplayed)
			is.Equal(paginationInfo.Offset, tt.expectedOffset)
			is.Equal(paginationInfo.Limit, tt.expectedLimit)
			is.Equal(paginationInfo.HasMore, tt.expectedHasMore)
		})
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) &&
		(str == substr ||
			(len(str) > len(substr) &&
				(str[:len(substr)] == substr ||
					str[len(str)-len(substr):] == substr ||
					findInString(str, substr))))
}

func findInString(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

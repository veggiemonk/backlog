package core

import (
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/afero"
)

func TestPaginate(t *testing.T) {
	tests := []struct {
		name              string
		tasks             []Task
		limit             int
		offset            int
		expectedCount     int
		expectedFirst     string // ID of first task in result
		expectedLast      string // ID of last task in result
		expectedTotal     int
		expectedDisplayed int
		expectedOffset    int
		expectedLimit     int
		expectedHasMore   bool
	}{
		{
			name: "no_pagination",
			tasks: []Task{
				{ID: TaskID{seg: []int{1}}, Title: "Task 1"},
				{ID: TaskID{seg: []int{2}}, Title: "Task 2"},
				{ID: TaskID{seg: []int{3}}, Title: "Task 3"},
			},
			limit:             0,
			offset:            0,
			expectedCount:     3,
			expectedFirst:     "01",
			expectedLast:      "03",
			expectedTotal:     3,
			expectedDisplayed: 3,
			expectedOffset:    0,
			expectedLimit:     0,
			expectedHasMore:   false,
		},
		{
			name: "limit_only",
			tasks: []Task{
				{ID: TaskID{seg: []int{1}}, Title: "Task 1"},
				{ID: TaskID{seg: []int{2}}, Title: "Task 2"},
				{ID: TaskID{seg: []int{3}}, Title: "Task 3"},
			},
			limit:             2,
			offset:            0,
			expectedCount:     2,
			expectedFirst:     "01",
			expectedLast:      "02",
			expectedTotal:     3,
			expectedDisplayed: 2,
			expectedOffset:    0,
			expectedLimit:     2,
			expectedHasMore:   true,
		},
		{
			name: "offset_only",
			tasks: []Task{
				{ID: TaskID{seg: []int{1}}, Title: "Task 1"},
				{ID: TaskID{seg: []int{2}}, Title: "Task 2"},
				{ID: TaskID{seg: []int{3}}, Title: "Task 3"},
			},
			limit:             0,
			offset:            1,
			expectedCount:     2,
			expectedFirst:     "02",
			expectedLast:      "03",
			expectedTotal:     3,
			expectedDisplayed: 2,
			expectedOffset:    1,
			expectedLimit:     0,
			expectedHasMore:   false,
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
			limit:             2,
			offset:            1,
			expectedCount:     2,
			expectedFirst:     "02",
			expectedLast:      "03",
			expectedTotal:     5,
			expectedDisplayed: 2,
			expectedOffset:    1,
			expectedLimit:     2,
			expectedHasMore:   true,
		},
		{
			name: "offset_beyond_end",
			tasks: []Task{
				{ID: TaskID{seg: []int{1}}, Title: "Task 1"},
				{ID: TaskID{seg: []int{2}}, Title: "Task 2"},
			},
			limit:             2,
			offset:            5,
			expectedCount:     0,
			expectedTotal:     2,
			expectedDisplayed: 0,
			expectedOffset:    5,
			expectedLimit:     2,
			expectedHasMore:   false,
		},
		{
			name: "limit_larger_than_remaining",
			tasks: []Task{
				{ID: TaskID{seg: []int{1}}, Title: "Task 1"},
				{ID: TaskID{seg: []int{2}}, Title: "Task 2"},
				{ID: TaskID{seg: []int{3}}, Title: "Task 3"},
			},
			limit:             5,
			offset:            1,
			expectedCount:     2,
			expectedFirst:     "02",
			expectedLast:      "03",
			expectedTotal:     3,
			expectedDisplayed: 2,
			expectedOffset:    1,
			expectedLimit:     5,
			expectedHasMore:   false,
		},
		{
			name:              "empty_task_list",
			tasks:             []Task{},
			limit:             5,
			offset:            0,
			expectedCount:     0,
			expectedTotal:     0,
			expectedDisplayed: 0,
			expectedOffset:    0,
			expectedLimit:     5,
			expectedHasMore:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			result := Paginate(tt.tasks, tt.limit, tt.offset)
			is.Equal(len(result.Tasks), tt.expectedCount)

			if tt.expectedCount > 0 {
				is.Equal(result.Tasks[0].ID.String(), tt.expectedFirst)
				is.Equal(result.Tasks[len(result.Tasks)-1].ID.String(), tt.expectedLast)
			}

			is.Equal(result.Pagination.TotalResults, tt.expectedTotal)
			is.Equal(result.Pagination.DisplayedResults, tt.expectedDisplayed)
			is.Equal(result.Pagination.Offset, tt.expectedOffset)
			is.Equal(result.Pagination.Limit, tt.expectedLimit)
			is.Equal(result.Pagination.HasMore, tt.expectedHasMore)
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
		params := ListTasksParams{Limit: 3}

		listResult, err := store.List(params)
		is.NoErr(err)
		is.Equal(len(listResult.Tasks), 3)
	})

	t.Run("list_with_offset", func(t *testing.T) {
		is := is.New(t)
		params := ListTasksParams{Offset: 2}

		listResult, err := store.List(params)
		is.NoErr(err)
		is.Equal(len(listResult.Tasks), 3) // Should have 3 remaining tasks (5 - 2 offset)
	})

	t.Run("list_with_limit_and_offset", func(t *testing.T) {
		is := is.New(t)
		params := ListTasksParams{Limit: 2, Offset: 1}

		listResult, err := store.List(params)
		is.NoErr(err)
		is.Equal(len(listResult.Tasks), 2)
	})

	t.Run("list_with_filter_and_pagination", func(t *testing.T) {
		is := is.New(t)
		params := ListTasksParams{Status: []string{"todo"}, Limit: 2}

		listResult, err := store.List(params)
		is.NoErr(err)
		is.Equal(len(listResult.Tasks), 2)

		// All should be "todo" (default status)
		for _, task := range listResult.Tasks {
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
		params := ListTasksParams{Limit: 2}

		listResult, err := store.List(params)
		is.NoErr(err)
		is.Equal(len(listResult.Tasks), 2) // Should limit results to 2

		// All should contain "feature"
		for _, task := range listResult.Tasks {
			is.True(strings.Contains(task.Title, "Feature") || strings.Contains(task.Description, "feature"))
		}
	})

	t.Run("search_with_offset", func(t *testing.T) {
		is := is.New(t)
		params := ListTasksParams{Offset: 1, Query: "feature"}

		listResult, err := store.List(params)
		is.NoErr(err)
		is.Equal(len(listResult.Tasks), 2) // Should have 2 remaining after offset (3 - 1)
	})

	t.Run("search_with_limit_and_offset", func(t *testing.T) {
		is := is.New(t)
		params := ListTasksParams{Limit: 1, Offset: 1, Query: "feature"}

		listResult, err := store.List(params)
		is.NoErr(err)
		is.Equal(len(listResult.Tasks), 1) // Should have 1 task (skip first, take one)
	})
}

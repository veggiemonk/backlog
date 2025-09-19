package core

import (
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/spf13/afero"
)

func TestFilterOptimizer(t *testing.T) {
	is := is.New(t)

	// Create test tasks
	tasks := createTestTasks()

	// Create and test optimizer
	optimizer := NewFilterOptimizer()
	optimizer.BuildIndexes(tasks)

	t.Run("BuildIndexes", func(t *testing.T) {
		is := is.New(t)

		// Check status index
		todoTasks := optimizer.statusIndex[StatusTodo]
		is.True(len(todoTasks) > 0)

		// Check assigned index
		assignedToAlice := optimizer.assignedIndex["alice"]
		is.True(len(assignedToAlice) > 0)

		// Check label index
		bugTasks := optimizer.labelIndex["bug"]
		is.True(len(bugTasks) > 0)
	})

	t.Run("EstimateFilterCost", func(t *testing.T) {
		is := is.New(t)

		params := ListTasksParams{
			Status: []string{"todo"},
		}

		cost := optimizer.EstimateFilterCost(params)
		is.True(cost > 0)
	})

	t.Run("OptimizeFilter", func(t *testing.T) {
		is := is.New(t)

		params := ListTasksParams{
			Status:   []string{"todo"},
			Assigned: []string{"alice"},
		}

		plan := optimizer.OptimizeFilter(params)
		is.True(plan.UseIndexes)
		is.True(plan.EstimatedCost > 0)
		is.True(len(plan.OptimizedOrder) > 0)
	})
}

func TestSmartFilterTasks(t *testing.T) {
	is := is.New(t)

	tasks := createTestTasks()
	optimizer := NewFilterOptimizer()
	optimizer.BuildIndexes(tasks)

	t.Run("NoFilters", func(t *testing.T) {
		is := is.New(t)

		params := ListTasksParams{}
		result, err := SmartFilterTasks(tasks, params, optimizer)

		is.NoErr(err)
		is.Equal(len(result), len(tasks))
	})

	t.Run("StatusFilter", func(t *testing.T) {
		is := is.New(t)

		params := ListTasksParams{
			Status: []string{"todo"},
		}
		result, err := SmartFilterTasks(tasks, params, optimizer)

		is.NoErr(err)
		is.True(len(result) > 0)

		// Verify all results have the correct status
		for _, task := range result {
			is.Equal(task.Status, StatusTodo)
		}
	})

	t.Run("AssignedFilter", func(t *testing.T) {
		is := is.New(t)

		params := ListTasksParams{
			Assigned: []string{"alice"},
		}
		result, err := SmartFilterTasks(tasks, params, optimizer)

		is.NoErr(err)
		is.True(len(result) > 0)

		// Verify all results are assigned to alice
		for _, task := range result {
			is.True(contains(task.Assigned, "alice"))
		}
	})

	t.Run("LabelFilter", func(t *testing.T) {
		is := is.New(t)

		params := ListTasksParams{
			Labels: []string{"bug"},
		}
		result, err := SmartFilterTasks(tasks, params, optimizer)

		is.NoErr(err)
		is.True(len(result) > 0)

		// Verify all results have the bug label
		for _, task := range result {
			is.True(contains(task.Labels, "bug"))
		}
	})

	t.Run("CompoundFilter", func(t *testing.T) {
		is := is.New(t)

		params := ListTasksParams{
			Status:   []string{"todo"},
			Assigned: []string{"alice"},
		}
		result, err := SmartFilterTasks(tasks, params, optimizer)

		is.NoErr(err)

		// Verify all results match both conditions
		for _, task := range result {
			is.Equal(task.Status, StatusTodo)
			is.True(contains(task.Assigned, "alice"))
		}
	})

	t.Run("PriorityFilter", func(t *testing.T) {
		is := is.New(t)

		params := ListTasksParams{
			Priority: stringPtr("high"),
		}
		result, err := SmartFilterTasks(tasks, params, optimizer)

		is.NoErr(err)

		// Verify all results have high priority
		for _, task := range result {
			is.Equal(task.Priority, PriorityHigh)
		}
	})
}

func TestFilterPerformance(t *testing.T) {
	is := is.New(t)

	// Create a large dataset for performance testing
	tasks := createLargeTestDataset(1000)

	t.Run("BenchmarkSmartVsTraditional", func(t *testing.T) {
		is := is.New(t)

		params := ListTasksParams{
			Status: []string{"todo"},
		}

		// Benchmark traditional filtering
		traditionalResult := BenchmarkFilter(tasks, params, false)
		is.True(traditionalResult.Duration > 0)
		is.True(traditionalResult.TasksPerSecond > 0)

		// Benchmark smart filtering
		smartResult := BenchmarkFilter(tasks, params, true)
		is.True(smartResult.Duration > 0)
		is.True(smartResult.TasksPerSecond > 0)

		// Both should return the same number of filtered tasks
		is.Equal(smartResult.FilteredCount, traditionalResult.FilteredCount)
	})
}

func TestFilterBenchmark(t *testing.T) {
	is := is.New(t)

	tasks := createTestTasks()
	params := ListTasksParams{
		Status: []string{"todo"},
	}

	t.Run("TraditionalFilter", func(t *testing.T) {
		is := is.New(t)

		result := BenchmarkFilter(tasks, params, false)
		is.Equal(result.FilterType, "status")
		is.Equal(result.TaskCount, len(tasks))
		is.True(result.Duration > 0)
		is.True(result.TasksPerSecond > 0)
		is.True(!result.OptimizedFilter)
	})

	t.Run("SmartFilter", func(t *testing.T) {
		is := is.New(t)

		result := BenchmarkFilter(tasks, params, true)
		is.Equal(result.FilterType, "status")
		is.Equal(result.TaskCount, len(tasks))
		is.True(result.Duration > 0)
		is.True(result.TasksPerSecond > 0)
		is.True(result.OptimizedFilter)
	})
}

func TestFileTaskStoreWithSmartFiltering(t *testing.T) {
	is := is.New(t)

	// Create in-memory filesystem
	fs := afero.NewMemMapFs()
	store := NewFileTaskStore(fs, ".backlog")

	// Create some test tasks
	task1, err := store.Create(CreateTaskParams{
		Title:       "Task 1",
		Description: "Description 1",
		Priority:    ptr("high"),
		Assigned:    []string{"alice"},
		Labels:      []string{"bug", "urgent"},
	})
	is.NoErr(err)

	_, err = store.Create(CreateTaskParams{
		Title:       "Task 2",
		Description: "Description 2",
		Priority:    ptr("medium"),
		Assigned:    []string{"bob"},
		Labels:      []string{"feature"},
	})
	is.NoErr(err)

	task3, err := store.Create(CreateTaskParams{
		Title:       "Task 3",
		Description: "Description 3",
		Priority:    ptr("high"),
		Assigned:    []string{"alice", "bob"},
		Labels:      []string{"bug"},
	})
	is.NoErr(err)

	t.Run("ListWithStatusFilter", func(t *testing.T) {
		is := is.New(t)

		params := ListTasksParams{
			Status: []string{"todo"},
		}

		tasks, err := store.List(params)
		is.NoErr(err)
		is.Equal(len(tasks), 3) // All tasks should be todo by default
	})

	t.Run("ListWithAssignedFilter", func(t *testing.T) {
		is := is.New(t)

		params := ListTasksParams{
			Assigned: []string{"alice"},
		}

		tasks, err := store.List(params)
		is.NoErr(err)
		is.Equal(len(tasks), 2) // task1 and task3

		// Verify results
		taskIDs := make([]string, len(tasks))
		for i, task := range tasks {
			taskIDs[i] = task.ID.Name()
		}
		is.True(contains(taskIDs, task1.ID.Name()))
		is.True(contains(taskIDs, task3.ID.Name()))
	})

	t.Run("ListWithLabelFilter", func(t *testing.T) {
		is := is.New(t)

		params := ListTasksParams{
			Labels: []string{"bug"},
		}

		tasks, err := store.List(params)
		is.NoErr(err)
		is.Equal(len(tasks), 2) // task1 and task3

		// Verify results
		taskIDs := make([]string, len(tasks))
		for i, task := range tasks {
			taskIDs[i] = task.ID.Name()
		}
		is.True(contains(taskIDs, task1.ID.Name()))
		is.True(contains(taskIDs, task3.ID.Name()))
	})

	t.Run("ListWithCompoundFilter", func(t *testing.T) {
		is := is.New(t)

		params := ListTasksParams{
			Priority: stringPtr("high"),
			Labels:   []string{"bug"},
		}

		tasks, err := store.List(params)
		is.NoErr(err)
		is.Equal(len(tasks), 2) // task1 and task3

		// Verify all results match criteria
		for _, task := range tasks {
			is.Equal(task.Priority, PriorityHigh)
			is.True(contains(task.Labels, "bug"))
		}
	})
}

// Helper functions

func createTestTasks() []*Task {
	now := time.Now()

	return []*Task{
		{
			ID:          TaskID{seg: []int{1}},
			Title:       "Task 1",
			Status:      StatusTodo,
			Priority:    PriorityHigh,
			Assigned:    []string{"alice"},
			Labels:      []string{"bug", "urgent"},
			CreatedAt:   now,
			UpdatedAt:   now,
			Description: "Test task 1",
		},
		{
			ID:          TaskID{seg: []int{2}},
			Title:       "Task 2",
			Status:      StatusInProgress,
			Priority:    PriorityMedium,
			Assigned:    []string{"bob"},
			Labels:      []string{"feature"},
			CreatedAt:   now,
			UpdatedAt:   now,
			Description: "Test task 2",
		},
		{
			ID:          TaskID{seg: []int{3}},
			Title:       "Task 3",
			Status:      StatusTodo,
			Priority:    PriorityHigh,
			Assigned:    []string{"alice", "bob"},
			Labels:      []string{"bug"},
			CreatedAt:   now,
			UpdatedAt:   now,
			Description: "Test task 3",
		},
		{
			ID:          TaskID{seg: []int{4}},
			Title:       "Task 4",
			Status:      StatusDone,
			Priority:    PriorityLow,
			Assigned:    []string{"charlie"},
			Labels:      []string{"documentation"},
			CreatedAt:   now,
			UpdatedAt:   now,
			Description: "Test task 4",
		},
	}
}

func createLargeTestDataset(count int) []*Task {
	tasks := make([]*Task, count)
	now := time.Now()

	statuses := []Status{StatusTodo, StatusInProgress, StatusDone}
	priorities := []Priority{PriorityLow, PriorityMedium, PriorityHigh}
	assignees := []string{"alice", "bob", "charlie", "diana"}
	labels := []string{"bug", "feature", "documentation", "test", "urgent"}

	for i := 0; i < count; i++ {
		tasks[i] = &Task{
			ID:          TaskID{seg: []int{i + 1}},
			Title:       "Generated Task " + string(rune(i+1)),
			Status:      statuses[i%len(statuses)],
			Priority:    priorities[i%len(priorities)],
			Assigned:    []string{assignees[i%len(assignees)]},
			Labels:      []string{labels[i%len(labels)]},
			CreatedAt:   now,
			UpdatedAt:   now,
			Description: "Generated test task",
		}
	}

	return tasks
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func stringPtr(s string) *string {
	return &s
}

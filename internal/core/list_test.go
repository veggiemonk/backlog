package core_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/afero"
	"github.com/veggiemonk/backlog/internal/core"
)

func TestListTasks(t *testing.T) {
	is := is.New(t)
	store := core.NewFileTaskStore(afero.NewMemMapFs(), ".backlog")
	t1Parent := "T1"
	_, _ = store.Create(core.CreateTaskParams{Title: "Task One"})
	taskTwo, _ := store.Create(core.CreateTaskParams{Title: "Task Two"})
	_, _ = store.Create(core.CreateTaskParams{Title: "Task Three", Parent: &t1Parent})
	is.NoErr(store.Update(&taskTwo, core.EditTaskParams{
		ID:        taskTwo.ID.String(),
		NewTitle:  &taskTwo.Title,
		NewStatus: ptr("in-progress"),
	}))

	tests := []struct {
		name          string
		params        core.ListTasksParams
		expectedCount int
		expectedTitle string
	}{
		{"no filter", core.ListTasksParams{}, 3, ""},
		{"filter by status", core.ListTasksParams{Status: []string{"todo"}}, 2, ""},
		{"filter by parent full name", core.ListTasksParams{Parent: ptr("T01")}, 1, "Task Three"},
		{"filter by status in progress", core.ListTasksParams{Status: []string{"in-progress"}}, 1, "Task Two"},
		{"filter by parent no leading 0", core.ListTasksParams{Parent: ptr("T1")}, 1, "Task Three"},
		{"filter by parent just number", core.ListTasksParams{Parent: ptr("1")}, 1, "Task Three"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := store.List(tt.params)
			is.NoErr(err)
			is.Equal(len(tasks), tt.expectedCount)
			if tt.expectedTitle != "" {
				is.Equal(tt.expectedTitle, tasks[0].Title)
			}
		})
	}
}

func TestFilterAndSortTasks(t *testing.T) {
	store := core.NewFileTaskStore(afero.NewMemMapFs(), ".backlog")

	is := is.New(t)
	// Create tasks for testing filtering and sorting
	_, _ = store.Create(core.CreateTaskParams{Title: "Alpha Task", Assigned: []string{"bob"}, Priority: "high"})
	_, _ = store.Create(core.CreateTaskParams{Title: "Bravo Task", Assigned: []string{"alice"}, Priority: "medium"})
	_, _ = store.Create(core.CreateTaskParams{Title: "Charlie Task", Assigned: []string{"alice", "bob"}, Priority: "low"})
	// Update status for one task
	taskTwo, _ := store.Get("T02")
	is.NoErr(store.Update(&taskTwo, core.EditTaskParams{
		ID:        taskTwo.ID.String(),
		NewTitle:  &taskTwo.Title,
		NewStatus: ptr("done"),
	}))

	parentID := "T01"
	_, _ = store.Create(core.CreateTaskParams{Title: "Delta Task", Parent: &parentID, Assigned: []string{"charlie"}})

	t.Run("filter by assignee", func(t *testing.T) {
		is := is.NewRelaxed(t)
		tasks, err := store.List(core.ListTasksParams{Assigned: []string{"alice"}})
		is.NoErr(err)
		is.Equal(len(tasks), 2) // Bravo, Charlie
	})

	t.Run("filter by multiple names assigned", func(t *testing.T) {
		is := is.NewRelaxed(t)
		tasks, err := store.List(core.ListTasksParams{Assigned: []string{"alice", "charlie"}})
		is.NoErr(err)
		is.Equal(len(tasks), 3) // Bravo, Charlie, Delta
	})

	t.Run("filter by parent and status", func(t *testing.T) {
		is := is.NewRelaxed(t)
		// All tasks are 'todo' by default, except T02 which is 'done'
		// T04 is a sub-task of T01
		// So there are no tasks that are sub-tasks of T01 and have status 'done'
		tasks, err := store.List(core.ListTasksParams{Parent: ptr("T01"), Status: []string{"done"}})
		is.NoErr(err)
		is.Equal(len(tasks), 0)

		// T04 is a sub-task of T01 and has status 'todo'
		tasks, err = store.List(core.ListTasksParams{Parent: ptr("T01"), Status: []string{"todo"}})
		is.NoErr(err)
		is.Equal(len(tasks), 1)
		is.Equal(tasks[0].Title, "Delta Task")
	})

	t.Run("sort by title", func(t *testing.T) {
		is := is.NewRelaxed(t)
		tasks, err := store.List(core.ListTasksParams{Sort: []string{"title"}})
		is.NoErr(err)
		is.Equal(len(tasks), 4)
		is.Equal(tasks[0].Title, "Alpha Task")
		is.Equal(tasks[1].Title, "Bravo Task")
		is.Equal(tasks[2].Title, "Charlie Task")
		is.Equal(tasks[3].Title, "Delta Task")
	})

	t.Run("sort by title reversed", func(t *testing.T) {
		is := is.NewRelaxed(t)
		tasks, err := store.List(core.ListTasksParams{Sort: []string{"title"}, Reverse: true})
		is.NoErr(err)
		is.Equal(len(tasks), 4)
		is.Equal(tasks[0].Title, "Delta Task")
		is.Equal(tasks[1].Title, "Charlie Task")
		is.Equal(tasks[2].Title, "Bravo Task")
		is.Equal(tasks[3].Title, "Alpha Task")
	})

	t.Run("sort by priority", func(t *testing.T) {
		is := is.NewRelaxed(t)
		tasks, err := store.List(core.ListTasksParams{Sort: []string{"priority"}})
		is.NoErr(err)
		is.Equal(len(tasks), 4)
		// Sorting is alphabetical on priority string: high, low, medium
		is.Equal(tasks[0].Priority.String(), "high")
		is.Equal(tasks[1].Priority.String(), "medium")
		is.Equal(tasks[2].Priority.String(), "low")
		// T04 has no priority, so it will be last.
	})

	t.Run("sort by id", func(t *testing.T) {
		is := is.New(t)
		tasks, err := store.List(core.ListTasksParams{Sort: []string{"id"}})
		is.NoErr(err)
		is.Equal(len(tasks), 4)
		is.Equal(tasks[0].ID.Name(), "T01")
		is.Equal(tasks[1].ID.Name(), "T01.01")
		is.Equal(tasks[2].ID.Name(), "T02")
		is.Equal(tasks[3].ID.Name(), "T03")
	})

	t.Run("invalid status filter", func(t *testing.T) {
		is := is.New(t)
		_, err := store.List(core.ListTasksParams{Status: []string{"invalid-status"}})
		is.True(err != nil)
	})

	t.Run("invalid parent filter", func(t *testing.T) {
		is := is.New(t)
		_, err := store.List(core.ListTasksParams{Parent: ptr("invalid-parent")})
		is.True(err != nil)
	})
}

func TestFilterUnassignedTasks(t *testing.T) {
	is := is.New(t)
	store := core.NewFileTaskStore(afero.NewMemMapFs(), ".backlog")

	// Create tasks with assigned
	_, _ = store.Create(core.CreateTaskParams{Title: "Task with Assigned", Assigned: []string{"alice"}})
	_, _ = store.Create(core.CreateTaskParams{Title: "Task with Multiple Assigned", Assigned: []string{"alice", "bob"}})

	// Create task without assigned
	_, _ = store.Create(core.CreateTaskParams{Title: "Unassigned Task"})

	t.Run("filter unassigned tasks", func(t *testing.T) {
		tasks, err := store.List(core.ListTasksParams{Unassigned: true})
		is.NoErr(err)
		is.Equal(len(tasks), 1) // Only the task with no one assigned
		is.Equal(tasks[0].Title, "Unassigned Task")
	})

	t.Run("filter unassigned tasks with status", func(t *testing.T) {
		is := is.New(t)
		// All unassigned tasks should have status 'todo' by default
		tasks, err := store.List(core.ListTasksParams{Unassigned: true, Status: []string{"todo"}})
		is.NoErr(err)
		is.Equal(len(tasks), 1) // Only the unassigned task with 'todo' status
		is.Equal(tasks[0].Title, "Unassigned Task")

		// There should be no unassigned tasks with 'done' status
		tasks, err = store.List(core.ListTasksParams{Unassigned: true, Status: []string{"done"}})
		is.NoErr(err)
		is.Equal(len(tasks), 0)
	})
}

func TestDependents(t *testing.T) {
	is := is.New(t)
	store := core.NewFileTaskStore(afero.NewMemMapFs(), ".backlog")

	// Create parent task
	baseTask, err := store.Create(core.CreateTaskParams{Title: "Base level Task"})
	is.NoErr(err)
	mediumTask, err := store.Create(core.CreateTaskParams{Title: "Medium level Task", Dependencies: []string{baseTask.ID.Name()}})
	is.NoErr(err)
	_, err = store.Create(core.CreateTaskParams{Title: "High level Task 1", Dependencies: []string{mediumTask.ID.Name()}})
	is.NoErr(err)
	_, err = store.Create(core.CreateTaskParams{Title: "High level Task 2", Dependencies: []string{mediumTask.ID.Name()}})
	is.NoErr(err)

	t.Run("list dependents tasks", func(t *testing.T) {
		tasks, err := store.List(core.ListTasksParams{DependedOn: true})
		is.NoErr(err)
		is.Equal(len(tasks), 2) // No parent specified, should return no tasks
	})

	t.Run("list tasks with dependencies", func(t *testing.T) {
		tasks, err := store.List(core.ListTasksParams{HasDependency: true})
		is.NoErr(err)
		is.Equal(len(tasks), 3) // Should return all tasks with dependencies
	})

	t.Run("list blocking tasks", func(t *testing.T) {
		tasks, err := store.List(core.ListTasksParams{DependedOn: true, Status: []string{"todo"}})
		is.NoErr(err)
		is.Equal(len(tasks), 2) // Non-existent parent, should return no tasks
	})

	// Mark medium task as done
	is.NoErr(store.Update(&mediumTask, core.EditTaskParams{
		ID:        mediumTask.ID.String(),
		NewStatus: ptr("done"),
	}))
	is.NoErr(err)

	t.Run("list blocking tasks after medium task done", func(t *testing.T) {
		tasks, err := store.List(core.ListTasksParams{DependedOn: true, Status: []string{"todo"}})
		is.NoErr(err)
		is.Equal(len(tasks), 1) // Medium task is done, should still return base level tasks
	})
}

func ptr[T any](v T) *T {
	return &v
}

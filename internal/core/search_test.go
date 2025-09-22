package core_test

import (
	"sort"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/afero"
	"github.com/veggiemonk/backlog/internal/core"
)

func TestSearchTasks(t *testing.T) {
	is := is.New(t)
	store := core.NewFileTaskStore(afero.NewMemMapFs(), ".backlog", core.NewMockLocker())

	// Create tasks with diverse content for searching
	_, _ = store.Create(core.CreateTaskParams{
		Title:       "Implement User Authentication",
		Description: "Add login and registration functionality.",
		Labels:      []string{"security", "feature"},
		Assigned:    []string{"alice"},
	})
	task2, _ := store.Create(core.CreateTaskParams{
		Title:       "Fix API Bug",
		Description: "The endpoint for fetching users is failing.",
		Assigned:    []string{"bob"},
	})
	task3, _ := store.Create(core.CreateTaskParams{
		Title:       "Refactor Database Schema",
		Description: "Update the schema to improve performance.",
		Labels:      []string{"refactoring", "database"},
		AC:          []string{"The new schema should be backward compatible."},
	})

	// Update tasks with implementation plan and notes
	_, err := store.Update(task2, core.EditTaskParams{
		ID:      task2.ID.String(),
		NewPlan: ptr("Investigate the logs and fix the issue."),
	})
	is.NoErr(err)

	_, err = store.Update(task3, core.EditTaskParams{
		ID:       task3.ID.String(),
		NewNotes: ptr("Remember to backup the database before applying changes."),
	})
	is.NoErr(err)

	tests := []struct {
		name           string
		query          string
		expectedCount  int
		expectedTitles []string
	}{
		{"search in title", "Authentication", 1, []string{"Implement User Authentication"}},
		{"search in description", "endpoint", 1, []string{"Fix API Bug"}},
		{"search in implementation plan", "investigate", 1, []string{"Fix API Bug"}},
		{"search in implementation notes", "backup", 1, []string{"Refactor Database Schema"}},
		{"search in acceptance criteria", "backward compatible", 1, []string{"Refactor Database Schema"}},
		{"search in labels", "database", 1, []string{"Refactor Database Schema"}},
		{"search in assignee", "alice", 1, []string{"Implement User Authentication"}},
		{"case-insensitive search", "authentication", 1, []string{"Implement User Authentication"}},
		{"no match", "nonexistent", 0, nil},
		{"query matches multiple fields in same task", "Refactor", 1, []string{"Refactor Database Schema"}},
		{"query matches multiple tasks", "user", 2, []string{"Fix API Bug", "Implement User Authentication"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := store.Search(tt.query, core.ListTasksParams{})
			is.NoErr(err)
			is.Equal(len(tasks), tt.expectedCount)

			if tt.expectedCount > 0 {
				var titles []string
				for _, task := range tasks {
					titles = append(titles, task.Title)
				}
				sort.Strings(titles)
				sort.Strings(tt.expectedTitles)
				is.Equal(titles, tt.expectedTitles)
			}
		})
	}
}

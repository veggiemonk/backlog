package core_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/afero"
	"github.com/veggiemonk/backlog/internal/core"
)

func TestEditTask(t *testing.T) {
	is := is.New(t)
	store := core.NewFileTaskStore(afero.NewMemMapFs(), ".backlog")
	createdTask, _ := store.Create(core.CreateTaskParams{Title: "Original Title"})

	newTitle := "Updated Title"
	params := core.EditTaskParams{
		ID:          createdTask.ID.String(),
		NewTitle:    &newTitle,
		AddAssigned: []string{"test-user"},
	}

	updatedTask, err := store.Update(createdTask, params)
	is.NoErr(err)
	is.Equal("Updated Title", updatedTask.Title)

	// Verify by reading the file again
	rereadTask, err := store.Get(createdTask.ID.String())
	is.NoErr(err)
	is.Equal("Updated Title", rereadTask.Title)
	is.Equal(core.MaybeStringArray{"test-user"}, rereadTask.Assigned)
}

func TestEditTaskHistory(t *testing.T) {
	is := is.New(t)
	store := core.NewFileTaskStore(afero.NewMemMapFs(), ".backlog")
	createdTask, err := store.Create(core.CreateTaskParams{Title: "Original Title"})
	is.NoErr(err)

	newTitle := "Updated Title"
	params := core.EditTaskParams{ID: createdTask.ID.String(), NewTitle: &newTitle}

	updatedTask, err := store.Update(createdTask, params)
	is.NoErr(err)

	is.Equal("Updated Title", updatedTask.Title)

	// Verify by reading the file again
	rereadTask, err := store.Get(createdTask.ID.String())
	is.NoErr(err)
	is.Equal("Updated Title", rereadTask.Title)

	// Check history
	is.Equal(len(rereadTask.History), 1)
	is.True(strings.Contains(rereadTask.History[0].Change, "Title changed from"))

	// Edit again to check multiple history entries
	newTitle2 := "Final Title"
	params2 := core.EditTaskParams{ID: createdTask.ID.String(), NewTitle: &newTitle2}

	_, err = store.Update(createdTask, params2)
	is.NoErr(err)
	rereadTask2, err := store.Get(createdTask.ID.String())
	is.NoErr(err)
	is.Equal("Final Title", rereadTask2.Title)
	is.Equal(len(rereadTask2.History), 2)
	is.True(strings.Contains(rereadTask2.History[1].Change, "Title changed from"))
}

func TestUpdateTaskFields(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	store := core.NewFileTaskStore(fs, ".backlog")

	// Create a parent task for testing parent update
	parentTask, _ := store.Create(core.CreateTaskParams{Title: "Parent Task"})

	t.Run("update various fields", func(t *testing.T) {
		task, _ := store.Create(core.CreateTaskParams{
			Title:       "Initial Task",
			Assigned:    []string{"initial-user"},
			Labels:      []string{"bug"},
			Notes:       &[]string{"Initial implementation notes."}[0],
			Plan:        &[]string{"Initial implementation plan."}[0],
			Priority:    ptr("medium"),
			Description: "This is the initial description.",
		})
		originalPath, err := fs.Stat(store.Path(task))
		is.NoErr(err)

		newDesc := "This is the new description."
		newStatus := "in-progress"
		newLabels := []string{"bug", "urgent"}
		removeLabels := []string{"bug"}
		newPriority := "high"
		newParent := parentTask.ID.String()
		newNotes := "These are the implementation notes."
		newPlan := "This is the implementation plan."
		newTitle := "Updated Task Title"
		newDeps := []string{"T01", "T02"}
		newAssigned := []string{"alice", "bob"}

		params := core.EditTaskParams{
			ID:              task.ID.String(),
			NewTitle:        &newTitle,
			NewDescription:  &newDesc,
			NewStatus:       &newStatus,
			AddLabels:       newLabels,
			RemoveLabels:    removeLabels,
			AddAssigned:     newAssigned,
			RemoveAssigned:  []string{"initial-user"},
			NewPriority:     &newPriority,
			NewParent:       &newParent,
			NewNotes:        &newNotes,
			NewPlan:         &newPlan,
			NewDependencies: newDeps,
		}

		updatedTask, err := store.Update(task, params)
		is.NoErr(err)

		slices.Sort(updatedTask.Assigned)
		// Verify fields are updated
		is.Equal(updatedTask.Title, newTitle)
		is.Equal(updatedTask.Description, newDesc)
		is.Equal(string(updatedTask.Status), newStatus)
		is.Equal(core.MaybeStringArray{"urgent"}, updatedTask.Labels)
		is.Equal(newAssigned, updatedTask.Assigned.ToSlice())
		is.Equal(updatedTask.Priority.String(), newPriority)
		is.True(updatedTask.Parent.Equals(parentTask.ID))
		is.Equal(updatedTask.ImplementationNotes, newNotes)
		is.Equal(strings.TrimSpace(updatedTask.ImplementationPlan), newPlan)
		is.Equal(core.MaybeStringArray(newDeps), updatedTask.Dependencies)

		// Verify old file is removed
		_, err = fs.Stat(originalPath.Name())
		is.True(err != nil) // Should be file not found

		// Verify new file exists
		_, err = fs.Stat(store.Path(updatedTask))
		is.NoErr(err)
	})

	t.Run("invalid status update", func(t *testing.T) {
		task, _ := store.Create(core.CreateTaskParams{Title: "Status Task"})
		invalidStatus := "invalid-status"
		_, err := store.Update(task, core.EditTaskParams{
			ID:        task.ID.String(),
			NewStatus: &invalidStatus,
		})
		is.True(err != nil) // Expecting an error
	})

	t.Run("invalid priority update", func(t *testing.T) {
		task, _ := store.Create(core.CreateTaskParams{Title: "Priority Task"})
		invalidPriority := "invalid-priority"
		_, err := store.Update(task, core.EditTaskParams{
			ID:          task.ID.String(),
			NewPriority: &invalidPriority,
		})
		is.True(err != nil) // Expecting an error
	})

	t.Run("invalid parent update", func(t *testing.T) {
		task, _ := store.Create(core.CreateTaskParams{Title: "Parent Task"})
		invalidParent := "invalid-parent"
		_, err := store.Update(task, core.EditTaskParams{
			ID:        task.ID.String(),
			NewParent: &invalidParent,
		})
		is.True(err != nil) // Expecting an error
	})

	t.Run("invalid priority update", func(t *testing.T) {
		task, _ := store.Create(core.CreateTaskParams{Title: "Dependency Task"})
		invalidDeps := []string{"T@#"}
		_, err := store.Update(task, core.EditTaskParams{
			ID:              task.ID.String(),
			NewDependencies: invalidDeps,
		})
		is.True(err != nil) // Expecting an error
	})
}

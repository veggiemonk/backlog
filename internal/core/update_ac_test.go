package core

import (
	"testing"

	"github.com/matryer/is"
)

func TestACManager(t *testing.T) {
	t.Run("add acceptance criteria", func(t *testing.T) {
		is := is.New(t)
		task := Task{}
		handleACChanges(&task, EditTaskParams{
			AddAC: []string{"AC 1", "AC 2"},
		})

		is.Equal(len(task.AcceptanceCriteria), 2)
		is.Equal(task.AcceptanceCriteria[0].Text, "AC 1")
		is.Equal(task.AcceptanceCriteria[0].Index, 1)
		is.Equal(task.AcceptanceCriteria[1].Text, "AC 2")
		is.Equal(task.AcceptanceCriteria[1].Index, 2)
	})

	t.Run("remove acceptance criteria", func(t *testing.T) {
		is := is.New(t)
		task := Task{
			AcceptanceCriteria: []AcceptanceCriterion{
				{Index: 1, Text: "AC 1"},
				{Index: 2, Text: "AC 2"},
				{Index: 3, Text: "AC 3"},
			},
		}
		handleACChanges(&task, EditTaskParams{
			RemoveAC: []int{2},
		})

		is.Equal(len(task.AcceptanceCriteria), 2)
		is.Equal(task.AcceptanceCriteria[0].Text, "AC 1")
		is.Equal(task.AcceptanceCriteria[0].Index, 1)
		is.Equal(task.AcceptanceCriteria[1].Text, "AC 3")
		is.Equal(task.AcceptanceCriteria[1].Index, 2) // Re-indexed
	})

	t.Run("check and uncheck acceptance criteria", func(t *testing.T) {
		is := is.New(t)
		task := Task{
			AcceptanceCriteria: []AcceptanceCriterion{
				{Index: 1, Text: "AC 1", Checked: false},
				{Index: 2, Text: "AC 2", Checked: true},
			},
		}

		// Check AC 1
		handleACChanges(&task, EditTaskParams{CheckAC: []int{1}})
		is.True(task.AcceptanceCriteria[0].Checked)

		// Uncheck AC 2
		handleACChanges(&task, EditTaskParams{UncheckAC: []int{2}})
		is.True(!task.AcceptanceCriteria[1].Checked)
	})

	t.Run("comprehensive changes", func(t *testing.T) {
		is := is.New(t)
		task := Task{
			AcceptanceCriteria: []AcceptanceCriterion{
				{Index: 1, Text: "Initial AC 1", Checked: false},
				{Index: 2, Text: "Initial AC 2", Checked: true},
				{Index: 3, Text: "Initial AC 3", Checked: false},
			},
		}
		handleACChanges(&task, EditTaskParams{
			AddAC:     []string{"New AC 4"},
			RemoveAC:  []int{2},
			CheckAC:   []int{1},
			UncheckAC: []int{2}, // This will be ignored as AC 2 is removed
		})

		is.Equal(len(task.AcceptanceCriteria), 3)

		// AC 1: was index 1, should be checked, still index 1
		is.Equal(task.AcceptanceCriteria[0].Text, "Initial AC 1")
		is.True(task.AcceptanceCriteria[0].Checked)
		is.Equal(task.AcceptanceCriteria[0].Index, 1)

		// AC 3: was index 3, should be re-indexed to 2
		is.Equal(task.AcceptanceCriteria[1].Text, "Initial AC 3")
		is.True(!task.AcceptanceCriteria[1].Checked)
		is.Equal(task.AcceptanceCriteria[1].Index, 2)

		// New AC 4: should be added with index 3
		is.Equal(task.AcceptanceCriteria[2].Text, "New AC 4")
		is.True(!task.AcceptanceCriteria[2].Checked)
		is.Equal(task.AcceptanceCriteria[2].Index, 3)
	})

	t.Run("remove multiple criteria", func(t *testing.T) {
		is := is.New(t)
		task := Task{
			AcceptanceCriteria: []AcceptanceCriterion{
				{Index: 1, Text: "AC 1"},
				{Index: 2, Text: "AC 2"},
				{Index: 3, Text: "AC 3"},
				{Index: 4, Text: "AC 4"},
			},
		}
		handleACChanges(&task, EditTaskParams{RemoveAC: []int{1, 3}})

		is.Equal(len(task.AcceptanceCriteria), 2)
		is.Equal(task.AcceptanceCriteria[0].Text, "AC 2")
		is.Equal(task.AcceptanceCriteria[0].Index, 1)
		is.Equal(task.AcceptanceCriteria[1].Text, "AC 4")
		is.Equal(task.AcceptanceCriteria[1].Index, 2)
	})

	t.Run("re-indexing after removal", func(t *testing.T) {
		is := is.New(t)
		task := Task{
			AcceptanceCriteria: []AcceptanceCriterion{
				{Index: 1, Text: "AC 1"},
				{Index: 2, Text: "AC 2"},
				{Index: 3, Text: "AC 3"},
			},
		}
		handleACChanges(&task, EditTaskParams{RemoveAC: []int{1}})

		is.Equal(len(task.AcceptanceCriteria), 2)
		is.Equal(task.AcceptanceCriteria[0].Index, 1)
		is.Equal(task.AcceptanceCriteria[1].Index, 2)
	})
}

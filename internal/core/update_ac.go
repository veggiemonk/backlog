package core

import (
	"fmt"
	"sort"
)

// ACManager handles acceptance criteria operations
type ACManager struct{}

// NewACManager creates a new ACManager instance
func NewACManager() *ACManager {
	return &ACManager{}
}

// HandleACChanges processes acceptance criteria changes for a task
func (ac *ACManager) HandleACChanges(task *Task, params EditTaskParams) {
	// 1. Remove ACs
	ac.removeACs(task, params.RemoveAC)

	// 2. Check/Uncheck ACs
	ac.checkACs(task, params.CheckAC)
	ac.uncheckACs(task, params.UncheckAC)

	// 3. Add new ACs
	ac.addACs(task, params.AddAC)

	// 4. Re-index all ACs
	ac.reindexACs(task)
}

// removeACs removes acceptance criteria by index
func (ac *ACManager) removeACs(task *Task, indicesToRemove []int) {
	// Sort in reverse order to avoid index shifting issues
	sort.Sort(sort.Reverse(sort.IntSlice(indicesToRemove)))

	for _, indexToRemove := range indicesToRemove {
		var newACs []AcceptanceCriterion
		for _, criterion := range task.AcceptanceCriteria {
			if criterion.Index != indexToRemove {
				newACs = append(newACs, criterion)
			} else {
				RecordChange(task, fmt.Sprintf("Removed acceptance criterion #%d: %q", criterion.Index, criterion.Text))
			}
		}
		task.AcceptanceCriteria = newACs
	}
}

// checkACs marks acceptance criteria as checked
func (ac *ACManager) checkACs(task *Task, indicesToCheck []int) {
	for _, indexToCheck := range indicesToCheck {
		for i := range task.AcceptanceCriteria {
			if task.AcceptanceCriteria[i].Index == indexToCheck && !task.AcceptanceCriteria[i].Checked {
				task.AcceptanceCriteria[i].Checked = true
				RecordChange(task, fmt.Sprintf("Checked acceptance criterion #%d: %q", task.AcceptanceCriteria[i].Index, task.AcceptanceCriteria[i].Text))
			}
		}
	}
}

// uncheckACs marks acceptance criteria as unchecked
func (ac *ACManager) uncheckACs(task *Task, indicesToUncheck []int) {
	for _, indexToUncheck := range indicesToUncheck {
		for i := range task.AcceptanceCriteria {
			if task.AcceptanceCriteria[i].Index == indexToUncheck && task.AcceptanceCriteria[i].Checked {
				task.AcceptanceCriteria[i].Checked = false
				RecordChange(task, fmt.Sprintf("Unchecked acceptance criterion #%d: %q", task.AcceptanceCriteria[i].Index, task.AcceptanceCriteria[i].Text))
			}
		}
	}
}

// addACs adds new acceptance criteria
func (ac *ACManager) addACs(task *Task, newCriteria []string) {
	for _, newCriterion := range newCriteria {
		highestIndex := 0
		for _, criterion := range task.AcceptanceCriteria {
			if criterion.Index > highestIndex {
				highestIndex = criterion.Index
			}
		}

		criterion := AcceptanceCriterion{
			Text:    newCriterion,
			Checked: false,
			Index:   highestIndex + 1,
		}
		task.AcceptanceCriteria = append(task.AcceptanceCriteria, criterion)
		RecordChange(task, fmt.Sprintf("Added acceptance criterion #%d: %q", criterion.Index, criterion.Text))
	}
}

// reindexACs re-indexes all acceptance criteria to ensure sequential numbering
func (ac *ACManager) reindexACs(task *Task) {
	sort.Slice(task.AcceptanceCriteria, func(i, j int) bool {
		return task.AcceptanceCriteria[i].Index < task.AcceptanceCriteria[j].Index
	})

	for i := range task.AcceptanceCriteria {
		task.AcceptanceCriteria[i].Index = i + 1
	}
}

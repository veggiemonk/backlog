package core

import (
	"fmt"
	"slices"
	"time"
)

// EditTaskParams holds the parameters for editing a task.
type EditTaskParams struct {
	ID              string   `json:"id"                         jsonschema:"Required. The ID of the task to edit."`
	NewTitle        *string  `json:"new_title,omitempty"        jsonschema:"A new title for the task."`
	NewDescription  *string  `json:"new_description,omitempty"  jsonschema:"A new description for the task."`
	NewStatus       *string  `json:"new_status,omitempty"       jsonschema:"A new status (e.g., 'in-progress', 'done')."`
	NewPriority     *string  `json:"new_priority,omitempty"     jsonschema:"A new priority."`
	NewParent       *string  `json:"new_parent,omitempty"       jsonschema:"A new parent task ID."`
	AddAssigned     []string `json:"add_assigned,omitempty"     jsonschema:"A new list of assigned."`
	RemoveAssigned  []string `json:"remove_assigned,omitempty"  jsonschema:"A list of assigned to remove."`
	AddLabels       []string `json:"add_labels,omitempty"       jsonschema:"Add new list of labels."`
	RemoveLabels    []string `json:"remove_labels,omitempty"    jsonschema:"A list of labels to remove."`
	NewDependencies []string `json:"new_dependencies,omitempty" jsonschema:"A new list of dependencies (replaces the old list)."`
	NewNotes        *string  `json:"new_notes,omitempty"        jsonschema:"New implementation notes."`
	NewPlan         *string  `json:"new_plan,omitempty"         jsonschema:"New implementation plan."`
	AddAC           []string `json:"add_ac,omitempty"           jsonschema:"A list of new acceptance criteria to add."`
	CheckAC         []int    `json:"check_ac,omitempty"         jsonschema:"A list of 1-based indices of AC to check."`
	UncheckAC       []int    `json:"uncheck_ac,omitempty"       jsonschema:"A list of 1-based indices of AC to uncheck."`
	RemoveAC        []int    `json:"remove_ac,omitempty"        jsonschema:"A list of 1-based indices of AC to remove."`
}

// Update updates an existing task based on the provided parameters.
func (f *FileTaskStore) Update(task *Task, params EditTaskParams) error {
	var oldFilePath string

	// Update fields based on params
	if params.NewTitle != nil && task.Title != *params.NewTitle {
		oldFilePath = f.Path(*task) // To use when moving the file
		RecordChange(task, fmt.Sprintf("Title changed from %q to %q", task.Title, *params.NewTitle))
		task.Title = *params.NewTitle
	}

	if params.NewDescription != nil && task.Description != *params.NewDescription {
		RecordChange(task, "Description changed")
		task.Description = *params.NewDescription
	}

	if params.NewStatus != nil {
		newStatus, err := ParseStatus(*params.NewStatus)
		if err != nil {
			return fmt.Errorf("invalid status %q: %w", *params.NewStatus, err)
		}
		if task.Status != newStatus {
			RecordChange(task, fmt.Sprintf("Status changed from %q to %q", task.Status, newStatus))
			task.Status = newStatus
		}
	}

	if len(params.AddAssigned) > 0 || len(params.RemoveAssigned) > 0 {
		newAssigned := batchRemoveAdd(task.Assigned, params.RemoveAssigned, params.AddAssigned)
		if newAssigned != nil {
			if !equalStringSlices(task.Assigned, newAssigned) {
				RecordChange(task, fmt.Sprintf("Assigned changed from %q to %q", task.Assigned, newAssigned))
			}
			task.Assigned = newAssigned
		}
	}
	if len(params.RemoveLabels) > 0 || len(params.AddLabels) > 0 {
		newLabels := batchRemoveAdd(task.Labels, params.RemoveLabels, params.AddLabels)
		if newLabels != nil {
			if !equalStringSlices(task.Labels, newLabels) {
				RecordChange(task, fmt.Sprintf("Labels changed from %q to %q", task.Labels, newLabels))
			}
			task.Labels = newLabels
		}
	}

	if params.NewPriority != nil {
		newPriority, err := ParsePriority(*params.NewPriority)
		if err != nil {
			return fmt.Errorf("invalid priority %q: %w", *params.NewPriority, err)
		}
		if task.Priority != newPriority {
			RecordChange(task, fmt.Sprintf("Priority changed from %q to %q", task.Priority, newPriority))
			task.Priority = newPriority
		}
	}

	if params.NewParent != nil {
		newParent, err := parseTaskID(*params.NewParent)
		if err != nil {
			return fmt.Errorf("invalid new parent task ID '%s': %w", *params.NewParent, err)
		}
		if !task.Parent.Equals(newParent) {
			if oldFilePath == "" {
				oldFilePath = f.Path(*task) // Save old file path before ID changes
			}
			RecordChange(task, fmt.Sprintf("Parent changed from %q to %q", task.Parent.String(), newParent.String()))
			task.Parent = newParent
			// Recalculate task ID to be a subtask of the new parent
			nextID, err := f.getNextTaskID(newParent.seg...)
			if err != nil {
				return fmt.Errorf("could not get next task ID for new parent: %w", err)
			}
			task.ID = nextID
		}
	}

	if params.NewNotes != nil && task.ImplementationNotes != *params.NewNotes {
		RecordChange(task, "Implementation notes changed")
		task.ImplementationNotes = *params.NewNotes
	}

	if params.NewPlan != nil && task.ImplementationPlan != *params.NewPlan {
		RecordChange(task, "Implementation plan changed")
		task.ImplementationPlan = fmt.Sprintf("%s\n", *params.NewPlan)
	}

	if params.NewDependencies != nil && !equalStringSlices(task.Dependencies, params.NewDependencies) {
		deps := make([]string, 0, len(params.NewDependencies))
		// check dependencies exists
		for _, depIDStr := range params.NewDependencies {
			depID, err := parseTaskID(depIDStr)
			if err != nil {
				return fmt.Errorf("invalid dependency task ID '%s': %w", depIDStr, err)
			}
			_, err = f.Get(depID.String())
			if err != nil {
				return fmt.Errorf("dependency task ID '%s' does not exist: %w", depIDStr, err)
			}
			deps = append(deps, depID.Name())
		}

		RecordChange(task, fmt.Sprintf("Dependencies changed from %q to %q", task.Dependencies, deps))
		task.Dependencies = deps
	}

	// Handle acceptance criteria changes
	handleACChanges(task, params)

	task.UpdatedAt = time.Now().UTC()

	if err := f.write(*task); err != nil {
		return fmt.Errorf("could not write updated task file: %w", err)
	}
	if oldFilePath != "" {
		if err := f.fs.Remove(oldFilePath); err != nil {
			return fmt.Errorf("could not remove old file: %w", err)
		}
	}
	return nil
}

func batchRemoveAdd(orig []string, toRemove []string, toAdd []string) []string {
	if len(toRemove) > 0 || len(toAdd) > 0 {
		labelSet := make(map[string]struct{})
		for _, l := range orig {
			labelSet[l] = struct{}{}
		}
		for _, l := range toAdd {
			labelSet[l] = struct{}{}
		}
		for _, l := range toRemove {
			delete(labelSet, l)
		}
		newLabels := make([]string, 0, len(labelSet))
		for l := range labelSet {
			newLabels = append(newLabels, l)
		}
		return newLabels
	}
	return nil
}

// equalStringSlices compares two string slices for equality
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	slices.Sort(a)
	slices.Sort(b)
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

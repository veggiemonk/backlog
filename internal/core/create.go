package core

import (
	"fmt"
	"time"

	"github.com/spf13/afero"
)

// CreateTaskParams holds the parameters for creating a new task.
type CreateTaskParams struct {
	Title        string   `json:"title"                  jsonschema:"Required. The title of the task."`
	Description  string   `json:"description"            jsonschema:"A detailed description of the task."`
	Priority     string   `json:"priority,omitempty"     jsonschema:"The priority of the task."`
	Parent       string   `json:"parent,omitempty"       jsonschema:"The ID of the parent task."`
	Assigned     []string `json:"assigned,omitempty"     jsonschema:"A list of names assigned."`
	Labels       []string `json:"labels,omitempty"       jsonschema:"A list of labels."`
	Dependencies []string `json:"dependencies,omitempty" jsonschema:"A list of task IDs that this task depends on."`
	AC           []string `json:"ac,omitempty"           jsonschema:"A list of acceptance criteria."`
	Plan         string   `json:"plan,omitempty"         jsonschema:"The implementation plan."`
	Notes        string   `json:"notes,omitempty"        jsonschema:"Additional notes."`
}

// Create implements TaskStore.
func (f *FileTaskStore) Create(params CreateTaskParams) (newTask Task, err error) {
	exists, err := afero.DirExists(f.fs, f.tasksDir)
	if err != nil {
		return newTask, fmt.Errorf("accessing %s error: %v", f.tasksDir, err)
	}
	if !exists {
		if err := f.fs.MkdirAll(f.tasksDir, 0o750); err != nil {
			return newTask, fmt.Errorf("could not create tasks directory %q: %w", f.tasksDir, err)
		}
	}
	var parentID TaskID
	if params.Parent != "" {
		parentID, err = parseTaskID(params.Parent)
		if err != nil {
			return newTask, fmt.Errorf("invalid parent task ID '%s': %w", params.Parent, err)
		}
		// Check if parent task actually exists
		_, err := f.Get(parentID.String())
		if err != nil {
			return newTask, fmt.Errorf("parent task ID '%s' does not exist: %w", params.Parent, err)
		}
	}
	nextID, err := f.getNextTaskID(parentID.seg...)
	if err != nil {
		return newTask, fmt.Errorf("could not get next task ID: %w", err)
	}

	deps := make([]string, 0, len(params.Dependencies))
	// check dependencies exists
	for _, depIDStr := range params.Dependencies {
		depID, err := parseTaskID(depIDStr)
		if err != nil {
			return newTask, fmt.Errorf("invalid dependency task ID '%s': %w", depIDStr, err)
		}
		_, err = f.Get(depID.String())
		if err != nil {
			return newTask, fmt.Errorf("dependency task ID '%s' does not exist: %w", depIDStr, err)
		}
		deps = append(deps, depID.Name())
	}

	newTask = NewTask()
	newTask.ID = nextID
	newTask.Title = params.Title
	newTask.Description = params.Description
	newTask.Parent = parentID
	newTask.Assigned = params.Assigned
	newTask.Labels = params.Labels
	newTask.Dependencies = deps
	newTask.Priority, err = ParsePriority(params.Priority)
	if err != nil {
		return newTask, fmt.Errorf("invalid priority %q: %w", params.Priority, err)
	}
	if params.Notes != "" {
		newTask.ImplementationNotes = fmt.Sprintf("%s\n", params.Notes)
	}
	if params.Plan != "" {
		newTask.ImplementationPlan = fmt.Sprintf("%s\n", params.Plan)
	}
	newTask.CreatedAt = time.Now().UTC()

	for i, criterion := range params.AC {
		newTask.AcceptanceCriteria = append(newTask.AcceptanceCriteria, AcceptanceCriterion{
			Text:    criterion,
			Checked: false,
			Index:   i + 1,
		})
	}

	if err := f.write(newTask); err != nil {
		return newTask, fmt.Errorf("could not write task file: %w", err)
	}
	return newTask, nil
}

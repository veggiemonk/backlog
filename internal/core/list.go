package core

import (
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/spf13/afero"
)


// ListTasksParams holds the parameters for listing tasks.
type ListTasksParams struct {
	Parent        *string  `json:"parent,omitempty" jsonschema:"Filter tasks by a parent task ID."`
	Status        []string `json:"status,omitempty" jsonschema:"Filter tasks by status."`
	Assigned      []string `json:"assigned,omitempty" jsonschema:"Filter tasks by assignee."`
	Labels        []string `json:"labels,omitempty" jsonschema:"Filter tasks by label."`
	Sort          []string `json:"sort,omitempty" jsonschema:"Fields to sort by."`
	Priority      *string  `json:"priority,omitempty" jsonschema:"Filter tasks by priority."`
	Unassigned    bool     `json:"unassigned,omitempty" jsonschema:"Filter tasks that have no one assigned."`
	DependedOn    bool     `json:"depended_on,omitempty" jsonschema:"Filter tasks that other tasks depend on."`
	HasDependency bool     `json:"has_dependency,omitempty" jsonschema:"Filter tasks that have at least one dependency."`
	Reverse       bool     `json:"reverse,omitempty" jsonschema:"Reverse the sort order."`
	// Pagination
	Limit  *int `json:"limit,omitempty" jsonschema:"Maximum number of tasks to return (0 means no limit)."`
	Offset *int `json:"offset,omitempty" jsonschema:"Number of tasks to skip from the beginning."`
}

// List implements TaskStore.
func (f *FileTaskStore) List(params ListTasksParams) (*ListResult, error) {
	// Load all tasks from filesystem
	tasks, err := f.loadAll()
	if err != nil {
		return nil, err
	}
	filteredTasks, err := FilterTasks(tasks, params)
	if err != nil {
		return nil, err
	}
	SortTasks(filteredTasks, params.Sort, params.Reverse)
	listResult := Paginate(filteredTasks, params.Limit, params.Offset)
	return listResult, nil
}

// LoadAll loads all tasks from the tasks directory.
func (f *FileTaskStore) loadAll() ([]Task, error) {
	exists, err := afero.DirExists(f.fs, f.tasksDir)
	if err != nil {
		return nil, err
	}
	if !exists {
		return []Task{}, nil
	}

	var tasks []Task
	walkErr := afero.Walk(f.fs, f.tasksDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasPrefix(info.Name(), TaskIDPrefix) && strings.HasSuffix(info.Name(), ".md") {
			b, err := afero.ReadFile(f.fs, path)
			if err != nil {
				return err
			}

			task, err := parseTask(b)
			if err != nil {
				return fmt.Errorf("parse task %s: %v", path, err)
			}
			tasks = append(tasks, task)
		}
		return nil
	})

	if walkErr != nil {
		return nil, walkErr
	}

	return tasks, nil
}

// FilterTasks applies filtering logic to a slice of tasks
func FilterTasks(tasks []Task, params ListTasksParams) ([]Task, error) {
	var parentID TaskID
	var statuses []Status
	var assigned []string
	var labels []string
	var priority Priority
	var isParentSet bool
	var isPrioritySet bool
	var err error

	if params.Parent != nil && *params.Parent != "" {
		parentID, err = parseTaskID(*params.Parent)
		if err != nil {
			return nil, fmt.Errorf("parent task ID '%s': %w", *params.Parent, err)
		}
		isParentSet = true
	}
	if params.Priority != nil && *params.Priority != "" {
		priority, err = ParsePriority(*params.Priority)
		if err != nil {
			return nil, fmt.Errorf("priority '%s': %w", *params.Priority, err)
		}
		isPrioritySet = true
	}
	for _, s := range params.Status {
		status, err := ParseStatus(s)
		if err != nil {
			return nil, fmt.Errorf("status '%s': %w", s, err)
		}
		statuses = append(statuses, status)
	}
	for _, a := range params.Assigned {
		assigned = append(assigned, strings.TrimSpace(a)) // clean up to ensure string equality
	}
	for _, l := range params.Labels {
		labels = append(labels, strings.TrimSpace(l)) // clean up to ensure string equality
	}

	if !isParentSet && !isPrioritySet && len(statuses) == 0 && len(assigned) == 0 && len(labels) == 0 && !params.Unassigned && !params.DependedOn && !params.HasDependency {
		return tasks, nil
	}

	if params.DependedOn {
		// Replace the tasks with tasks who are depended on.
		tasks = dependentGraph(tasks)
	}
	filteredTasks := make([]Task, 0, len(tasks))
	for _, t := range tasks {
		if isParentSet && !t.Parent.Equals(parentID) {
			continue
		}
		if isPrioritySet && t.Priority != priority {
			continue
		}
		if params.Unassigned && len(t.Assigned) > 0 {
			continue
		}
		if len(statuses) > 0 && !slices.Contains(statuses, t.Status) {
			continue
		}
		if len(assigned) > 0 && !atLeastOneIntersect(t.Assigned, assigned) {
			continue
		}
		if len(labels) > 0 && !atLeastOneIntersect(t.Labels, labels) {
			continue
		}
		if params.HasDependency && len(t.Dependencies) == 0 {
			continue
		}
		filteredTasks = append(filteredTasks, t)
	}

	return filteredTasks, nil
}

func dependentGraph(tasks []Task) []Task {
	dependents := []Task{}
	// for the dependency graph
	taskIndex := map[string]*Task{}                             // map[taskID]*Task
	taskDependencyIndex := make(map[string]map[string]struct{}) // map[X]map[Y]struct{} where X depends on Y
	for _, t := range tasks {
		// keep index for filling the tree of dependency
		taskIndex[t.ID.Name()] = &t
		// initialize the map
		taskDependencyIndex[t.ID.Name()] = make(map[string]struct{})
		// build the graph ie. X -> Y where X depends on Y
		for _, dep := range t.Dependencies {
			taskDependencyIndex[t.ID.Name()][dep] = struct{}{}
		}
	}
	// Reverse the graph ie. Y -> X where Y is depended on by X
	visited := make(map[string]struct{})
	for _ /*dependee*/, dependent := range taskDependencyIndex {
		for id := range dependent {
			if _, ok := visited[id]; ok {
				continue
			}
			visited[id] = struct{}{}
			if taskIndex[id] != nil {
				// add the task only once
				dependents = append(dependents, *taskIndex[id])
			}
		}
	}
	return dependents
}

func atLeastOneIntersect[S ~[]E, E comparable](got, want S) bool {
	if len(got) == 0 {
		return false
	}
	for _, g := range got {
		if slices.Contains(want, g) {
			return true
		}
	}
	return false
}

// SortTasks sorts the tasks slice based on the provided sort fields.
// Supported sort fields: id, title, status, priority, created, updated
func SortTasks(tasks []Task, sortFields []string, reverse bool) {
	if len(sortFields) == 0 {
		// No sorting requested, but still apply reverse if requested
		if reverse {
			slices.Reverse(tasks)
		}
		return
	}
	sort.Slice(tasks, func(i, j int) bool {
		t1, t2 := tasks[i], tasks[j]

		for _, field := range sortFields {
			field = strings.TrimSpace(strings.ToLower(field))
			var cmp int

			switch field {
			case "id":
				if t1.ID.Less(t2.ID) {
					return true
				} else if t2.ID.Less(t1.ID) {
					return false
				}
				cmp = 0
			case "title":
				cmp = strings.Compare(strings.ToLower(t1.Title), strings.ToLower(t2.Title))
			case "status":
				cmp = strings.Compare(string(t1.Status), string(t2.Status))
			case "priority":
				// Sort by priority value (higher priority comes first)
				if int(t1.Priority) > int(t2.Priority) {
					return true
				} else if int(t1.Priority) < int(t2.Priority) {
					return false
				}
				cmp = 0
			case "created":
				if t1.CreatedAt.Before(t2.CreatedAt) {
					return true
				} else if t1.CreatedAt.After(t2.CreatedAt) {
					return false
				}
				cmp = 0
			case "updated":
				if t1.UpdatedAt.Before(t2.UpdatedAt) {
					return true
				} else if t1.UpdatedAt.After(t2.UpdatedAt) {
					return false
				}
				cmp = 0
			default:
				// Unknown field, skip
				continue
			}

			if cmp < 0 {
				return true
			} else if cmp > 0 {
				return false
			}
			// If equal, continue to next sort field
		}

		// All sort fields are equal
		return false
	})

	if reverse {
		slices.Reverse(tasks)
	}
}


package core

import "strings"

// Search searches for tasks containing the query string in various fields.
func (f *FileTaskStore) Search(query string, listParams ListTasksParams) ([]*Task, error) {
	// Get all tasks and search in memory
	tasks, err := f.List(ListTasksParams{})
	if err != nil {
		return nil, err
	}

	matches := []*Task{}
	queryLower := strings.ToLower(query)

	for _, task := range tasks {
		// Search in task title, description, and other text fields
		if strings.Contains(strings.ToLower(task.Title), queryLower) ||
			strings.Contains(strings.ToLower(task.Description), queryLower) ||
			strings.Contains(strings.ToLower(task.ImplementationPlan), queryLower) ||
			strings.Contains(strings.ToLower(task.ImplementationNotes), queryLower) {
			matches = append(matches, task)
			continue
		}

		// Search in acceptance criteria
		for _, ac := range task.AcceptanceCriteria {
			if strings.Contains(strings.ToLower(ac.Text), queryLower) {
				matches = append(matches, task)
				break
			}
		}

		// Search in labels and assigned names
		for _, label := range task.Labels {
			if strings.Contains(strings.ToLower(label), queryLower) {
				matches = append(matches, task)
				break
			}
		}

		// Search in names assigned
		for _, assignee := range task.Assigned {
			if strings.Contains(strings.ToLower(assignee), queryLower) {
				matches = append(matches, task)
				break
			}
		}

		// Search in priority
		if strings.Contains(strings.ToLower(task.Priority.String()), queryLower) {
			matches = append(matches, task)
		}
	}

	// Apply filtering and sorting to the search results
	filteredMatches, err := FilterTasks(matches, listParams)
	if err != nil {
		return nil, err
	}

	SortTasks(filteredMatches, listParams.Sort, listParams.Reverse)
	return filteredMatches, nil
}

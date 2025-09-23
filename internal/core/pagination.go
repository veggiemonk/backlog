package core

// PaginationInfo contains metadata about pagination results
type PaginationInfo struct {
	TotalResults     int  `json:"total_results"`
	DisplayedResults int  `json:"displayed_results"`
	Offset           int  `json:"offset"`
	Limit            int  `json:"limit"`
	HasMore          bool `json:"has_more"`
}

// ListResult contains the tasks and pagination metadata
type ListResult struct {
	Tasks      []Task          `json:"tasks"`
	Pagination *PaginationInfo `json:"pagination,omitempty"`
}

// Paginate applies pagination to a slice of tasks and returns a ListResult.
func Paginate(tasks []Task, limit *int, offset *int) *ListResult {
	totalResults := len(tasks)

	// Default offset to 0 if nil
	startIndex := 0
	if offset != nil && *offset > 0 {
		startIndex = *offset
	}

	// If offset is beyond the total number of tasks, return empty result
	if startIndex >= totalResults {
		return &ListResult{
			Tasks: []Task{},
			Pagination: &PaginationInfo{
				TotalResults:     totalResults,
				DisplayedResults: 0,
				Offset:           startIndex,
				Limit:            getLimit(limit),
				HasMore:          false,
			},
		}
	}

	// If no limit specified or limit is 0, return from offset to end
	if limit == nil || *limit == 0 {
		paginatedTasks := tasks[startIndex:]
		return &ListResult{
			Tasks: paginatedTasks,
			Pagination: &PaginationInfo{
				TotalResults:     totalResults,
				DisplayedResults: len(paginatedTasks),
				Offset:           startIndex,
				Limit:            getLimit(limit),
				HasMore:          false,
			},
		}
	}

	// Calculate end index
	endIndex := min(startIndex+(*limit), totalResults)
	paginatedTasks := tasks[startIndex:endIndex]

	hasMore := endIndex < totalResults

	return &ListResult{
		Tasks: paginatedTasks,
		Pagination: &PaginationInfo{
			TotalResults:     totalResults,
			DisplayedResults: len(paginatedTasks),
			Offset:           startIndex,
			Limit:            getLimit(limit),
			HasMore:          hasMore,
		},
	}
}

// PaginateTasks applies pagination to a slice of tasks.
// If limit is nil or 0, returns all tasks.
// If offset is nil, defaults to 0.
func PaginateTasks(tasks []Task, limit *int, offset *int) []Task {
	if len(tasks) == 0 {
		return tasks
	}

	// Default offset to 0 if nil
	startIndex := 0
	if offset != nil && *offset > 0 {
		startIndex = *offset
	}

	// If offset is beyond the total number of tasks, return empty slice
	if startIndex >= len(tasks) {
		return []Task{}
	}

	// If no limit specified or limit is 0, return from offset to end
	if limit == nil || *limit == 0 {
		return tasks[startIndex:]
	}

	// Calculate end index
	endIndex := min(startIndex+(*limit), len(tasks))

	return tasks[startIndex:endIndex]
}

func getLimit(limit *int) int {
	if limit == nil {
		return 0
	}
	return *limit
}

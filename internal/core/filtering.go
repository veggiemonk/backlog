package core

import (
	"fmt"
	"time"
)

// FilterOperator represents logical operators for compound filtering
type FilterOperator string

const (
	FilterAND FilterOperator = "AND"
	FilterOR  FilterOperator = "OR"
)

// CompoundFilter represents a compound filter with multiple conditions
type CompoundFilter struct {
	Operator   FilterOperator  `json:"operator"`   // AND or OR
	Conditions []FilterCondition `json:"conditions"` // List of filter conditions
}

// FilterCondition represents a single filter condition
type FilterCondition struct {
	Field    string      `json:"field"`    // Field to filter on (status, assigned, labels, etc.)
	Operator string      `json:"operator"` // Comparison operator (eq, in, contains, etc.)
	Value    any         `json:"value"`    // Value to compare against
}

// FilterOptimizer provides query optimization for common filter patterns
type FilterOptimizer struct {
	// Index of tasks by common fields for faster lookups
	statusIndex     map[Status][]*Task
	assignedIndex   map[string][]*Task
	labelIndex      map[string][]*Task
	priorityIndex   map[Priority][]*Task
	parentIndex     map[string][]*Task
	dependencyGraph map[string][]string // Cached dependency relationships
}

// NewFilterOptimizer creates a new filter optimizer with indexes
func NewFilterOptimizer() *FilterOptimizer {
	return &FilterOptimizer{
		statusIndex:     make(map[Status][]*Task),
		assignedIndex:   make(map[string][]*Task),
		labelIndex:      make(map[string][]*Task),
		priorityIndex:   make(map[Priority][]*Task),
		parentIndex:     make(map[string][]*Task),
		dependencyGraph: make(map[string][]string),
	}
}

// BuildIndexes creates optimized indexes for the given tasks
func (f *FilterOptimizer) BuildIndexes(tasks []*Task) {
	// Clear existing indexes
	f.statusIndex = make(map[Status][]*Task)
	f.assignedIndex = make(map[string][]*Task)
	f.labelIndex = make(map[string][]*Task)
	f.priorityIndex = make(map[Priority][]*Task)
	f.parentIndex = make(map[string][]*Task)
	f.dependencyGraph = make(map[string][]string)

	for _, task := range tasks {
		// Status index
		f.statusIndex[task.Status] = append(f.statusIndex[task.Status], task)

		// Assigned index
		for _, assignee := range task.Assigned {
			f.assignedIndex[assignee] = append(f.assignedIndex[assignee], task)
		}

		// Label index
		for _, label := range task.Labels {
			f.labelIndex[label] = append(f.labelIndex[label], task)
		}

		// Priority index
		f.priorityIndex[task.Priority] = append(f.priorityIndex[task.Priority], task)

		// Parent index
		if !task.Parent.Equals(ZeroTaskID) {
			parentID := task.Parent.Name()
			f.parentIndex[parentID] = append(f.parentIndex[parentID], task)
		}

		// Dependency graph
		taskID := task.ID.Name()
		for _, dep := range task.Dependencies {
			f.dependencyGraph[dep] = append(f.dependencyGraph[dep], taskID)
		}
	}
}

// OptimizeFilter analyzes the filter parameters and returns an optimized execution plan
func (f *FilterOptimizer) OptimizeFilter(params ListTasksParams) FilterExecutionPlan {
	plan := FilterExecutionPlan{
		UseIndexes:     true,
		EstimatedCost:  0,
		OptimizedOrder: []string{},
	}

	// Estimate costs and determine optimal filter order
	if len(params.Status) > 0 {
		cost := f.estimateStatusFilterCost(params.Status)
		plan.EstimatedCost += cost
		plan.OptimizedOrder = append(plan.OptimizedOrder, "status")
	}

	if params.Priority != nil {
		cost := f.estimatePriorityFilterCost(*params.Priority)
		plan.EstimatedCost += cost
		plan.OptimizedOrder = append(plan.OptimizedOrder, "priority")
	}

	if len(params.Assigned) > 0 {
		cost := f.estimateAssignedFilterCost(params.Assigned)
		plan.EstimatedCost += cost
		plan.OptimizedOrder = append(plan.OptimizedOrder, "assigned")
	}

	if len(params.Labels) > 0 {
		cost := f.estimateLabelFilterCost(params.Labels)
		plan.EstimatedCost += cost
		plan.OptimizedOrder = append(plan.OptimizedOrder, "labels")
	}

	if params.Parent != nil {
		cost := f.estimateParentFilterCost(*params.Parent)
		plan.EstimatedCost += cost
		plan.OptimizedOrder = append(plan.OptimizedOrder, "parent")
	}

	// Sort filters by estimated cost (lowest first for better performance)
	// In a real implementation, you would sort based on actual cost calculations
	return plan
}

// FilterExecutionPlan represents an optimized filter execution strategy
type FilterExecutionPlan struct {
	UseIndexes     bool     `json:"use_indexes"`
	EstimatedCost  int      `json:"estimated_cost"`
	OptimizedOrder []string `json:"optimized_order"`
}

// EstimateFilterCost estimates the computational cost of applying filters
func (f *FilterOptimizer) EstimateFilterCost(params ListTasksParams) int {
	cost := 0

	// Base cost for loading tasks
	cost += 100

	// Add cost for each filter type
	if len(params.Status) > 0 {
		cost += f.estimateStatusFilterCost(params.Status)
	}
	if params.Priority != nil {
		cost += f.estimatePriorityFilterCost(*params.Priority)
	}
	if len(params.Assigned) > 0 {
		cost += f.estimateAssignedFilterCost(params.Assigned)
	}
	if len(params.Labels) > 0 {
		cost += f.estimateLabelFilterCost(params.Labels)
	}
	if params.Parent != nil {
		cost += f.estimateParentFilterCost(*params.Parent)
	}

	return cost
}

// Helper methods for cost estimation
func (f *FilterOptimizer) estimateStatusFilterCost(statuses []string) int {
	totalTasks := 0
	for _, statusStr := range statuses {
		if status, err := ParseStatus(statusStr); err == nil {
			totalTasks += len(f.statusIndex[status])
		}
	}
	return totalTasks * 2 // Low cost due to indexing
}

func (f *FilterOptimizer) estimatePriorityFilterCost(priorityStr string) int {
	if priority, err := ParsePriority(priorityStr); err == nil {
		return len(f.priorityIndex[priority]) * 2
	}
	return 50 // Default cost if parsing fails
}

func (f *FilterOptimizer) estimateAssignedFilterCost(assigned []string) int {
	totalTasks := 0
	for _, assignee := range assigned {
		totalTasks += len(f.assignedIndex[assignee])
	}
	return totalTasks * 3 // Slightly higher cost due to potential overlaps
}

func (f *FilterOptimizer) estimateLabelFilterCost(labels []string) int {
	totalTasks := 0
	for _, label := range labels {
		totalTasks += len(f.labelIndex[label])
	}
	return totalTasks * 3 // Similar to assigned
}

func (f *FilterOptimizer) estimateParentFilterCost(parent string) int {
	return len(f.parentIndex[parent]) * 2
}

// SmartFilterTasks applies optimized filtering using indexes and early termination
func SmartFilterTasks(tasks []*Task, params ListTasksParams, optimizer *FilterOptimizer) ([]*Task, error) {
	// Validate filter parameters first (same validation as original FilterTasks)
	err := validateFilterParams(params)
	if err != nil {
		return nil, err
	}

	// If no filters are specified, return all tasks
	if !hasAnyFilters(params) {
		return tasks, nil
	}

	// Build or use existing indexes with original task set
	if optimizer == nil {
		optimizer = NewFilterOptimizer()
		optimizer.BuildIndexes(tasks)
	}

	// Handle special case for DependedOn filter (same logic as original FilterTasks)
	if params.DependedOn {
		tasks = dependentGraph(tasks)
		// Rebuild indexes with the transformed task set if we have other filters
		if hasOtherFilters(params) {
			optimizer.BuildIndexes(tasks)
		}
	}

	// Start with the most selective filter to minimize subsequent operations
	var candidateTasks []*Task

	// Apply the most selective filter first
	candidateTasks = getInitialCandidates(tasks, params, optimizer)

	// Apply remaining filters with early termination
	filteredTasks := make([]*Task, 0, len(candidateTasks))
	for _, task := range candidateTasks {
		if matchesAllFilters(task, params) {
			filteredTasks = append(filteredTasks, task)
		}
	}

	return filteredTasks, nil
}

// validateFilterParams validates the filter parameters and returns an error if any are invalid
func validateFilterParams(params ListTasksParams) error {
	// Validate parent ID
	if params.Parent != nil && *params.Parent != "" {
		_, err := parseTaskID(*params.Parent)
		if err != nil {
			return fmt.Errorf("parent task ID '%s': %w", *params.Parent, err)
		}
	}

	// Validate priority
	if params.Priority != nil && *params.Priority != "" {
		_, err := ParsePriority(*params.Priority)
		if err != nil {
			return fmt.Errorf("priority '%s': %w", *params.Priority, err)
		}
	}

	// Validate status values
	for _, s := range params.Status {
		_, err := ParseStatus(s)
		if err != nil {
			return fmt.Errorf("status '%s': %w", s, err)
		}
	}

	return nil
}

// getInitialCandidates gets the initial set of candidate tasks using the most selective index
func getInitialCandidates(tasks []*Task, params ListTasksParams, optimizer *FilterOptimizer) []*Task {
	// Find the most selective filter and use its index
	if len(params.Status) > 0 {
		return getTasksByStatus(params.Status, optimizer)
	}
	if params.Priority != nil {
		return getTasksByPriority(*params.Priority, optimizer)
	}
	if params.Parent != nil {
		return getTasksByParent(*params.Parent, optimizer)
	}
	if len(params.Assigned) > 0 {
		return getTasksByAssigned(params.Assigned, optimizer)
	}
	if len(params.Labels) > 0 {
		return getTasksByLabels(params.Labels, optimizer)
	}

	// If no indexed filters, return all tasks
	return tasks
}

// Index-based retrieval methods
func getTasksByStatus(statuses []string, optimizer *FilterOptimizer) []*Task {
	var result []*Task
	seen := make(map[string]bool)

	for _, statusStr := range statuses {
		if status, err := ParseStatus(statusStr); err == nil {
			for _, task := range optimizer.statusIndex[status] {
				taskID := task.ID.Name()
				if !seen[taskID] {
					seen[taskID] = true
					result = append(result, task)
				}
			}
		}
	}
	return result
}

func getTasksByPriority(priorityStr string, optimizer *FilterOptimizer) []*Task {
	if priority, err := ParsePriority(priorityStr); err == nil {
		return optimizer.priorityIndex[priority]
	}
	return []*Task{}
}

func getTasksByParent(parent string, optimizer *FilterOptimizer) []*Task {
	// Parse the parent ID to ensure consistent comparison
	parentID, err := parseTaskID(parent)
	if err != nil {
		return []*Task{}
	}
	parentName := parentID.Name()
	return optimizer.parentIndex[parentName]
}

func getTasksByAssigned(assigned []string, optimizer *FilterOptimizer) []*Task {
	var result []*Task
	seen := make(map[string]bool)

	for _, assignee := range assigned {
		for _, task := range optimizer.assignedIndex[assignee] {
			taskID := task.ID.Name()
			if !seen[taskID] {
				seen[taskID] = true
				result = append(result, task)
			}
		}
	}
	return result
}

func getTasksByLabels(labels []string, optimizer *FilterOptimizer) []*Task {
	var result []*Task
	seen := make(map[string]bool)

	for _, label := range labels {
		for _, task := range optimizer.labelIndex[label] {
			taskID := task.ID.Name()
			if !seen[taskID] {
				seen[taskID] = true
				result = append(result, task)
			}
		}
	}
	return result
}

// hasAnyFilters checks if any filtering is requested
func hasAnyFilters(params ListTasksParams) bool {
	return params.Parent != nil ||
		len(params.Status) > 0 ||
		len(params.Assigned) > 0 ||
		len(params.Labels) > 0 ||
		params.Priority != nil ||
		params.Unassigned ||
		params.DependedOn ||
		params.HasDependency
}

// hasOtherFilters checks if there are filters other than DependedOn
func hasOtherFilters(params ListTasksParams) bool {
	return params.Parent != nil ||
		len(params.Status) > 0 ||
		len(params.Assigned) > 0 ||
		len(params.Labels) > 0 ||
		params.Priority != nil ||
		params.Unassigned ||
		params.HasDependency
}

// matchesAllFilters checks if a task matches all the specified filters
func matchesAllFilters(task *Task, params ListTasksParams) bool {
	// Check parent filter
	if params.Parent != nil && *params.Parent != "" {
		parentID, err := parseTaskID(*params.Parent)
		if err != nil || !task.Parent.Equals(parentID) {
			return false
		}
	}

	// Check priority filter
	if params.Priority != nil && *params.Priority != "" {
		priority, err := ParsePriority(*params.Priority)
		if err != nil || task.Priority != priority {
			return false
		}
	}

	// Check unassigned filter
	if params.Unassigned && len(task.Assigned) > 0 {
		return false
	}

	// Check status filter
	if len(params.Status) > 0 {
		found := false
		for _, statusStr := range params.Status {
			if status, err := ParseStatus(statusStr); err == nil && task.Status == status {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check assigned filter
	if len(params.Assigned) > 0 && !atLeastOneIntersect(task.Assigned, params.Assigned) {
		return false
	}

	// Check labels filter
	if len(params.Labels) > 0 && !atLeastOneIntersect(task.Labels, params.Labels) {
		return false
	}

	// Check dependency filter
	if params.HasDependency && len(task.Dependencies) == 0 {
		return false
	}

	// Note: DependedOn is handled by transforming the task list beforehand,
	// so we don't need to check it here

	return true
}

// FilterBenchmarkResult represents the result of a filtering performance benchmark
type FilterBenchmarkResult struct {
	FilterType      string        `json:"filter_type"`
	TaskCount       int           `json:"task_count"`
	FilteredCount   int           `json:"filtered_count"`
	Duration        time.Duration `json:"duration"`
	TasksPerSecond  float64       `json:"tasks_per_second"`
	OptimizedFilter bool          `json:"optimized_filter"`
}

// BenchmarkFilter benchmarks filtering performance for different scenarios
func BenchmarkFilter(tasks []*Task, params ListTasksParams, useOptimized bool) FilterBenchmarkResult {
	start := time.Now()

	var filteredTasks []*Task
	var err error

	if useOptimized {
		optimizer := NewFilterOptimizer()
		optimizer.BuildIndexes(tasks)
		filteredTasks, err = SmartFilterTasks(tasks, params, optimizer)
	} else {
		filteredTasks, err = FilterTasks(tasks, params)
	}

	duration := time.Since(start)

	// Determine filter type for reporting
	filterType := "unknown"
	if len(params.Status) > 0 {
		filterType = "status"
	} else if params.Priority != nil {
		filterType = "priority"
	} else if len(params.Assigned) > 0 {
		filterType = "assigned"
	} else if len(params.Labels) > 0 {
		filterType = "labels"
	}

	tasksPerSecond := 0.0
	if duration.Seconds() > 0 {
		tasksPerSecond = float64(len(tasks)) / duration.Seconds()
	}

	filteredCount := 0
	if err == nil {
		filteredCount = len(filteredTasks)
	}

	return FilterBenchmarkResult{
		FilterType:      filterType,
		TaskCount:       len(tasks),
		FilteredCount:   filteredCount,
		Duration:        duration,
		TasksPerSecond:  tasksPerSecond,
		OptimizedFilter: useOptimized,
	}
}
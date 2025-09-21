package core

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/afero"
)

// ConflictType represents the type of ID conflict detected
type ConflictType int

const (
	ConflictTypeDuplicateID ConflictType = iota
	ConflictTypeInvalidHierarchy
	ConflictTypeOrphanedChild
)

// String returns the string representation of ConflictType
func (ct ConflictType) String() string {
	switch ct {
	case ConflictTypeDuplicateID:
		return "duplicate_id"
	case ConflictTypeInvalidHierarchy:
		return "invalid_hierarchy"
	case ConflictTypeOrphanedChild:
		return "orphaned_child"
	default:
		return "unknown"
	}
}

// IDConflict represents a detected conflict in task IDs
type IDConflict struct {
	Type        ConflictType
	ConflictID  TaskID
	Files       []string
	Tasks       []*Task
	Description string
	DetectedAt  time.Time
}

// ConflictDetector handles detection and resolution of ID conflicts
type ConflictDetector struct {
	fs       afero.Fs
	tasksDir string
}

// NewConflictDetector creates a new conflict detector
func NewConflictDetector(fs afero.Fs, tasksDir string) *ConflictDetector {
	return &ConflictDetector{
		fs:       fs,
		tasksDir: tasksDir,
	}
}

// DetectConflicts scans all task files and identifies ID conflicts
func (cd *ConflictDetector) DetectConflicts() ([]IDConflict, error) {
	var conflicts []IDConflict

	// Get all task files
	files, err := afero.ReadDir(cd.fs, cd.tasksDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks directory: %w", err)
	}

	// Parse all tasks and build ID map
	idToFiles := make(map[string][]string)
	idToTasks := make(map[string][]*Task)
	allTasks := make([]*Task, 0)

	for _, file := range files {
		if file.IsDir() || !strings.HasPrefix(file.Name(), TaskIDPrefix) || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(cd.tasksDir, file.Name())
		task, err := cd.parseTaskFromFile(filePath)
		if err != nil {
			// Skip files that can't be parsed - they might be corrupted
			continue
		}

		idStr := task.ID.String()
		idToFiles[idStr] = append(idToFiles[idStr], filePath)
		idToTasks[idStr] = append(idToTasks[idStr], task)
		allTasks = append(allTasks, task)
	}

	// Detect duplicate IDs
	for idStr, files := range idToFiles {
		if len(files) > 1 {
			id, _ := parseTaskID(idStr)
			conflicts = append(conflicts, IDConflict{
				Type:        ConflictTypeDuplicateID,
				ConflictID:  id,
				Files:       files,
				Tasks:       idToTasks[idStr],
				Description: fmt.Sprintf("Task ID %s appears in multiple files: %v", idStr, files),
				DetectedAt:  time.Now(),
			})
		}
	}

	// Detect hierarchy conflicts (orphaned children)
	for _, task := range allTasks {
		if !task.Parent.IsZero() {
			parentStr := task.Parent.String()
			if _, exists := idToTasks[parentStr]; !exists {
				conflicts = append(conflicts, IDConflict{
					Type:       ConflictTypeOrphanedChild,
					ConflictID: task.ID,
					Files:      []string{cd.getTaskFilePath(task.ID)},
					Tasks:      []*Task{task},
					Description: fmt.Sprintf("Task %s references non-existent parent %s",
						task.ID.String(), parentStr),
					DetectedAt: time.Now(),
				})
			}
		}
	}

	// Detect invalid hierarchy (child ID doesn't match parent structure)
	for _, task := range allTasks {
		if !task.Parent.IsZero() {
			expectedParent := task.ID.Parent()
			if expectedParent != nil && !expectedParent.Equals(task.Parent) {
				conflicts = append(conflicts, IDConflict{
					Type:       ConflictTypeInvalidHierarchy,
					ConflictID: task.ID,
					Files:      []string{cd.getTaskFilePath(task.ID)},
					Tasks:      []*Task{task},
					Description: fmt.Sprintf("Task %s has incorrect parent %s, expected %s based on ID structure",
						task.ID.String(), task.Parent.String(), expectedParent.String()),
					DetectedAt: time.Now(),
				})
			}
		}
	}

	return conflicts, nil
}

// parseTaskFromFile reads and parses a task from a file
func (cd *ConflictDetector) parseTaskFromFile(filePath string) (*Task, error) {
	content, err := afero.ReadFile(cd.fs, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	task, err := parseTask(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse task from %s: %w", filePath, err)
	}

	return task, nil
}

// getTaskFilePath constructs the expected file path for a task ID
func (cd *ConflictDetector) getTaskFilePath(id TaskID) string {
	// This is a simplified version - in reality we'd need to find the actual file
	// since we don't know the title part of the filename
	return filepath.Join(cd.tasksDir, id.Name()+"-*.md")
}

// IsZero checks if TaskID is zero value
func (t TaskID) IsZero() bool {
	return len(t.seg) == 0
}

// ConflictSummary provides a summary of detected conflicts
type ConflictSummary struct {
	TotalConflicts    int
	DuplicateIDs      int
	OrphanedChildren  int
	InvalidHierarchy  int
	ConflictsByType   map[ConflictType][]IDConflict
}

// SummarizeConflicts creates a summary of the provided conflicts
func SummarizeConflicts(conflicts []IDConflict) ConflictSummary {
	summary := ConflictSummary{
		TotalConflicts:  len(conflicts),
		ConflictsByType: make(map[ConflictType][]IDConflict),
	}

	for _, conflict := range conflicts {
		summary.ConflictsByType[conflict.Type] = append(summary.ConflictsByType[conflict.Type], conflict)

		switch conflict.Type {
		case ConflictTypeDuplicateID:
			summary.DuplicateIDs++
		case ConflictTypeOrphanedChild:
			summary.OrphanedChildren++
		case ConflictTypeInvalidHierarchy:
			summary.InvalidHierarchy++
		}
	}

	return summary
}

// ResolutionStrategy defines how conflicts should be resolved
type ResolutionStrategy int

const (
	ResolutionStrategyChronological ResolutionStrategy = iota // Keep older task, renumber newer
	ResolutionStrategyAutoRenumber                            // Automatically renumber conflicting IDs
	ResolutionStrategyManual                                  // Require manual resolution
)

// ResolutionAction represents an action to resolve a conflict
type ResolutionAction struct {
	Type        string            // "renumber", "update_parent", "delete"
	OriginalID  TaskID
	NewID       TaskID
	FilePath    string
	Description string
	Metadata    map[string]any    // Additional metadata for the action
}

// ResolutionPlan contains all actions needed to resolve conflicts
type ResolutionPlan struct {
	Actions     []ResolutionAction
	Summary     string
	Strategy    ResolutionStrategy
	CreatedAt   time.Time
}

// ConflictResolver handles resolution of ID conflicts
type ConflictResolver struct {
	detector *ConflictDetector
	store    *FileTaskStore
}

// NewConflictResolver creates a new conflict resolver
func NewConflictResolver(detector *ConflictDetector, store *FileTaskStore) *ConflictResolver {
	return &ConflictResolver{
		detector: detector,
		store:    store,
	}
}

// CreateResolutionPlan generates a plan to resolve the given conflicts
func (cr *ConflictResolver) CreateResolutionPlan(conflicts []IDConflict, strategy ResolutionStrategy) (*ResolutionPlan, error) {
	plan := &ResolutionPlan{
		Actions:   make([]ResolutionAction, 0),
		Strategy:  strategy,
		CreatedAt: time.Now(),
	}

	switch strategy {
	case ResolutionStrategyChronological:
		return cr.createChronologicalPlan(conflicts, plan)
	case ResolutionStrategyAutoRenumber:
		return cr.createAutoRenumberPlan(conflicts, plan)
	case ResolutionStrategyManual:
		return cr.createManualPlan(conflicts, plan)
	default:
		return nil, fmt.Errorf("unsupported resolution strategy: %v", strategy)
	}
}

// createChronologicalPlan resolves conflicts by keeping older tasks and renumbering newer ones
func (cr *ConflictResolver) createChronologicalPlan(conflicts []IDConflict, plan *ResolutionPlan) (*ResolutionPlan, error) {
	for _, conflict := range conflicts {
		if conflict.Type == ConflictTypeDuplicateID {
			// Sort tasks by creation date
			tasks := make([]*Task, len(conflict.Tasks))
			copy(tasks, conflict.Tasks)

			// Sort by creation date (oldest first)
			for i := 0; i < len(tasks)-1; i++ {
				for j := i + 1; j < len(tasks); j++ {
					if tasks[i].CreatedAt.After(tasks[j].CreatedAt) {
						tasks[i], tasks[j] = tasks[j], tasks[i]
					}
				}
			}

			// Keep the oldest task, renumber the rest
			for i := 1; i < len(tasks); i++ {
				newID, err := cr.findNextAvailableID(tasks[i].ID)
				if err != nil {
					return nil, fmt.Errorf("failed to find available ID for %s: %w", tasks[i].ID.String(), err)
				}

				plan.Actions = append(plan.Actions, ResolutionAction{
					Type:        "renumber",
					OriginalID:  tasks[i].ID,
					NewID:       newID,
					FilePath:    conflict.Files[i],
					Description: fmt.Sprintf("Renumber task %s to %s (chronological resolution)", tasks[i].ID.String(), newID.String()),
					Metadata: map[string]any{
						"reason":      "duplicate_id_chronological",
						"created_at":  tasks[i].CreatedAt,
						"older_task":  tasks[0].ID.String(),
					},
				})
			}
		}
	}

	plan.Summary = fmt.Sprintf("Chronological resolution: %d renumbering actions", len(plan.Actions))
	return plan, nil
}

// createAutoRenumberPlan automatically renumbers conflicting IDs
func (cr *ConflictResolver) createAutoRenumberPlan(conflicts []IDConflict, plan *ResolutionPlan) (*ResolutionPlan, error) {
	for _, conflict := range conflicts {
		switch conflict.Type {
		case ConflictTypeDuplicateID:
			// Renumber all but the first task found
			for i := 1; i < len(conflict.Tasks); i++ {
				newID, err := cr.findNextAvailableID(conflict.Tasks[i].ID)
				if err != nil {
					return nil, fmt.Errorf("failed to find available ID for %s: %w", conflict.Tasks[i].ID.String(), err)
				}

				plan.Actions = append(plan.Actions, ResolutionAction{
					Type:        "renumber",
					OriginalID:  conflict.Tasks[i].ID,
					NewID:       newID,
					FilePath:    conflict.Files[i],
					Description: fmt.Sprintf("Auto-renumber task %s to %s", conflict.Tasks[i].ID.String(), newID.String()),
					Metadata: map[string]any{
						"reason": "duplicate_id_auto",
					},
				})
			}

		case ConflictTypeOrphanedChild:
			// For orphaned children, we can either:
			// 1. Remove the parent reference (make it a top-level task)
			// 2. Find a suitable parent
			// For auto-resolution, we'll remove the parent reference
			plan.Actions = append(plan.Actions, ResolutionAction{
				Type:        "update_parent",
				OriginalID:  conflict.ConflictID,
				NewID:       ZeroTaskID, // Remove parent reference
				FilePath:    conflict.Files[0],
				Description: fmt.Sprintf("Remove invalid parent reference from task %s", conflict.ConflictID.String()),
				Metadata: map[string]any{
					"reason":           "orphaned_child",
					"original_parent":  conflict.Tasks[0].Parent.String(),
				},
			})

		case ConflictTypeInvalidHierarchy:
			// Fix the parent reference to match the ID structure
			expectedParent := conflict.ConflictID.Parent()
			if expectedParent != nil {
				plan.Actions = append(plan.Actions, ResolutionAction{
					Type:        "update_parent",
					OriginalID:  conflict.ConflictID,
					NewID:       *expectedParent,
					FilePath:    conflict.Files[0],
					Description: fmt.Sprintf("Fix parent reference for task %s to %s", conflict.ConflictID.String(), expectedParent.String()),
					Metadata: map[string]any{
						"reason":           "invalid_hierarchy",
						"original_parent":  conflict.Tasks[0].Parent.String(),
						"expected_parent":  expectedParent.String(),
					},
				})
			}
		}
	}

	plan.Summary = fmt.Sprintf("Auto-resolution: %d actions", len(plan.Actions))
	return plan, nil
}

// createManualPlan creates a plan that requires manual intervention
func (cr *ConflictResolver) createManualPlan(conflicts []IDConflict, plan *ResolutionPlan) (*ResolutionPlan, error) {
	for _, conflict := range conflicts {
		plan.Actions = append(plan.Actions, ResolutionAction{
			Type:        "manual",
			OriginalID:  conflict.ConflictID,
			Description: fmt.Sprintf("Manual resolution required for %s: %s", conflict.Type, conflict.Description),
			Metadata: map[string]any{
				"conflict_type": conflict.Type,
				"files":         conflict.Files,
				"description":   conflict.Description,
			},
		})
	}

	plan.Summary = fmt.Sprintf("Manual resolution required for %d conflicts", len(plan.Actions))
	return plan, nil
}

// findNextAvailableID finds the next available ID in the sequence
func (cr *ConflictResolver) findNextAvailableID(conflictID TaskID) (TaskID, error) {
	// Use the store's getNextTaskID method to find the next available ID
	// For hierarchical IDs, we need to consider the parent path
	if conflictID.HasSubTasks() {
		parent := conflictID.Parent()
		if parent != nil {
			return cr.store.getNextTaskID(parent.seg...)
		}
	}

	// For top-level tasks, find the next available top-level ID
	return cr.store.getNextTaskID()
}

// ExecuteResolutionPlan executes the given resolution plan
func (cr *ConflictResolver) ExecuteResolutionPlan(plan *ResolutionPlan, dryRun bool) ([]string, error) {
	var results []string

	if dryRun {
		results = append(results, "DRY RUN MODE - No changes will be made")
	}

	for _, action := range plan.Actions {
		result, err := cr.executeAction(action, dryRun)
		if err != nil {
			return results, fmt.Errorf("failed to execute action %s: %w", action.Description, err)
		}
		results = append(results, result)
	}

	return results, nil
}

// executeAction executes a single resolution action
func (cr *ConflictResolver) executeAction(action ResolutionAction, dryRun bool) (string, error) {
	switch action.Type {
	case "renumber":
		return cr.executeRenumberAction(action, dryRun)
	case "update_parent":
		return cr.executeUpdateParentAction(action, dryRun)
	case "manual":
		return fmt.Sprintf("MANUAL: %s", action.Description), nil
	default:
		return "", fmt.Errorf("unsupported action type: %s", action.Type)
	}
}

// executeRenumberAction executes a renumbering action
func (cr *ConflictResolver) executeRenumberAction(action ResolutionAction, dryRun bool) (string, error) {
	if dryRun {
		return fmt.Sprintf("WOULD RENUMBER: %s -> %s", action.OriginalID.String(), action.NewID.String()), nil
	}

	// Read the task
	task, err := cr.detector.parseTaskFromFile(action.FilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read task: %w", err)
	}

	// Update the task ID
	oldID := task.ID
	task.ID = action.NewID

	// Record the change in history using enhanced tracking
	RecordIDChange(task, oldID, action.NewID, "conflict resolution", action.Metadata)
	task.UpdatedAt = time.Now()

	// Create new file with new ID
	newFilePath := cr.store.Path(task)
	if err := cr.store.write(task); err != nil {
		return "", fmt.Errorf("failed to write updated task: %w", err)
	}

	// Remove old file
	if err := cr.detector.fs.Remove(action.FilePath); err != nil {
		return "", fmt.Errorf("failed to remove old file: %w", err)
	}

	return fmt.Sprintf("RENUMBERED: %s -> %s (file: %s -> %s)", oldID.String(), action.NewID.String(), action.FilePath, newFilePath), nil
}

// executeUpdateParentAction executes a parent update action
func (cr *ConflictResolver) executeUpdateParentAction(action ResolutionAction, dryRun bool) (string, error) {
	if dryRun {
		if action.NewID.IsZero() {
			return fmt.Sprintf("WOULD REMOVE PARENT: %s", action.OriginalID.String()), nil
		}
		return fmt.Sprintf("WOULD UPDATE PARENT: %s -> %s", action.OriginalID.String(), action.NewID.String()), nil
	}

	// Read the task
	task, err := cr.detector.parseTaskFromFile(action.FilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read task: %w", err)
	}

	// Update the parent
	oldParent := task.Parent
	task.Parent = action.NewID

	// Record the change in history using enhanced tracking
	RecordParentChange(task, oldParent, action.NewID, "conflict resolution")
	task.UpdatedAt = time.Now()

	// Write the updated task
	if err := cr.store.write(task); err != nil {
		return "", fmt.Errorf("failed to write updated task: %w", err)
	}

	if action.NewID.IsZero() {
		return fmt.Sprintf("REMOVED PARENT: %s (was %s)", action.OriginalID.String(), oldParent.String()), nil
	}
	return fmt.Sprintf("UPDATED PARENT: %s -> %s", oldParent.String(), action.NewID.String()), nil
}

// ReferenceUpdater handles updating references when task IDs change
type ReferenceUpdater struct {
	detector *ConflictDetector
	store    *FileTaskStore
}

// NewReferenceUpdater creates a new reference updater
func NewReferenceUpdater(detector *ConflictDetector, store *FileTaskStore) *ReferenceUpdater {
	return &ReferenceUpdater{
		detector: detector,
		store:    store,
	}
}

// UpdateReferences updates all references to changed task IDs
func (ru *ReferenceUpdater) UpdateReferences(idChanges map[TaskID]TaskID) error {
	if len(idChanges) == 0 {
		return nil
	}

	// Get all task files
	files, err := afero.ReadDir(ru.detector.fs, ru.detector.tasksDir)
	if err != nil {
		return fmt.Errorf("failed to read tasks directory: %w", err)
	}

	var updatedTasks []*Task
	for _, file := range files {
		if file.IsDir() || !strings.HasPrefix(file.Name(), TaskIDPrefix) || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(ru.detector.tasksDir, file.Name())
		task, err := ru.detector.parseTaskFromFile(filePath)
		if err != nil {
			continue // Skip corrupted files
		}

		updated := false

		// Update parent references
		if newParentID, exists := idChanges[task.Parent]; exists {
			oldParent := task.Parent
			task.Parent = newParentID
			RecordParentChange(task, oldParent, newParentID, "ID change cascade")
			updated = true
		}

		// Update dependencies
		if len(task.Dependencies) > 0 {
			newDeps := make([]string, 0, len(task.Dependencies.ToSlice()))
			depsUpdated := false

			for _, dep := range task.Dependencies.ToSlice() {
				depID, err := parseTaskID(dep)
				if err != nil {
					newDeps = append(newDeps, dep) // Keep as-is if not a valid ID
					continue
				}

				if newDepID, exists := idChanges[depID]; exists {
					newDeps = append(newDeps, newDepID.String())
					depsUpdated = true
				} else {
					newDeps = append(newDeps, dep)
				}
			}

			if depsUpdated {
				// Update dependencies array
				task.Dependencies = MaybeStringArrayFromSlice(newDeps)
				RecordChange(task, fmt.Sprintf("Updated dependencies due to ID changes"))
				updated = true
			}
		}

		if updated {
			task.UpdatedAt = time.Now()
			updatedTasks = append(updatedTasks, task)
		}
	}

	// Write all updated tasks
	for _, task := range updatedTasks {
		if err := ru.store.write(task); err != nil {
			return fmt.Errorf("failed to write updated task %s: %w", task.ID.String(), err)
		}
	}

	return nil
}

// FindTaskReferences finds all tasks that reference the given task ID
func (ru *ReferenceUpdater) FindTaskReferences(targetID TaskID) ([]*Task, error) {
	var referencingTasks []*Task

	// Get all task files
	files, err := afero.ReadDir(ru.detector.fs, ru.detector.tasksDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks directory: %w", err)
	}

	targetIDStr := targetID.String()

	for _, file := range files {
		if file.IsDir() || !strings.HasPrefix(file.Name(), TaskIDPrefix) || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(ru.detector.tasksDir, file.Name())
		task, err := ru.detector.parseTaskFromFile(filePath)
		if err != nil {
			continue // Skip corrupted files
		}

		// Check parent reference
		if task.Parent.String() == targetIDStr {
			referencingTasks = append(referencingTasks, task)
			continue
		}

		// Check dependencies
		for _, dep := range task.Dependencies.ToSlice() {
			if dep == targetIDStr {
				referencingTasks = append(referencingTasks, task)
				break
			}
		}
	}

	return referencingTasks, nil
}

// MaybeStringArrayFromSlice converts a string slice to MaybeStringArray
func MaybeStringArrayFromSlice(slice []string) MaybeStringArray {
	var msa MaybeStringArray
	for _, s := range slice {
		msa = append(msa, s)
	}
	return msa
}

// ExecuteResolutionPlanWithReferences executes a resolution plan and updates all references
func (cr *ConflictResolver) ExecuteResolutionPlanWithReferences(plan *ResolutionPlan, dryRun bool) ([]string, error) {
	var results []string
	idChanges := make(map[TaskID]TaskID)

	if dryRun {
		results = append(results, "DRY RUN MODE - No changes will be made")
	}

	// First pass: execute all actions and collect ID changes
	for _, action := range plan.Actions {
		result, err := cr.executeAction(action, dryRun)
		if err != nil {
			return results, fmt.Errorf("failed to execute action %s: %w", action.Description, err)
		}
		results = append(results, result)

		// Collect ID changes for reference updating
		if action.Type == "renumber" && !dryRun {
			idChanges[action.OriginalID] = action.NewID
		}
	}

	// Second pass: update all references to changed IDs
	if !dryRun && len(idChanges) > 0 {
		updater := NewReferenceUpdater(cr.detector, cr.store)
		if err := updater.UpdateReferences(idChanges); err != nil {
			return results, fmt.Errorf("failed to update references: %w", err)
		}

		results = append(results, fmt.Sprintf("Updated references for %d changed IDs", len(idChanges)))
	}

	return results, nil
}
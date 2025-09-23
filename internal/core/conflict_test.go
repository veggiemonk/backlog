package core

import (
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/spf13/afero"
)

func TestConflictType_String(t *testing.T) {
	is := is.New(t)

	testCases := []struct {
		conflictType ConflictType
		expected     string
	}{
		{ConflictTypeDuplicateID, "duplicate_id"},
		{ConflictTypeInvalidHierarchy, "invalid_hierarchy"},
		{ConflictTypeOrphanedChild, "orphaned_child"},
		{ConflictType(999), "unknown"},
	}

	for _, tc := range testCases {
		is.Equal(tc.conflictType.String(), tc.expected)
	}
}

func TestTaskID_IsZero(t *testing.T) {
	is := is.New(t)

	// Test zero value
	var zeroID TaskID
	is.True(zeroID.IsZero())

	// Test non-zero value
	nonZeroID, _ := parseTaskID("T1")
	is.True(!nonZeroID.IsZero())

	// Test ZeroTaskID constant
	is.True(ZeroTaskID.IsZero())
}

func TestNewConflictDetector(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	tasksDir := ".backlog"

	detector := NewConflictDetector(fs, tasksDir)
	is.True(detector != nil)
	is.Equal(detector.tasksDir, tasksDir)
}

func TestConflictDetector_DetectConflicts_NoDuplicates(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	store := NewFileTaskStore(fs, ".backlog")

	// Create some normal tasks without conflicts
	task1, err := store.Create(CreateTaskParams{
		Title:       "Task 1",
		Description: "First task",
		AC:          []string{"AC 1"},
	})
	is.NoErr(err)

	task2, err := store.Create(CreateTaskParams{
		Title:       "Task 2",
		Description: "Second task",
		AC:          []string{"AC 2"},
	})
	is.NoErr(err)

	detector := NewConflictDetector(fs, ".backlog")
	conflicts, err := detector.DetectConflicts()
	is.NoErr(err)
	is.Equal(len(conflicts), 0)

	// Verify the tasks were created correctly
	is.Equal(task1.ID.String(), "01")
	is.Equal(task2.ID.String(), "02")
}

func TestConflictDetector_DetectConflicts_DuplicateIDs(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	store := NewFileTaskStore(fs, ".backlog")

	// Create a task normally
	_, err := store.Create(CreateTaskParams{
		Title:       "Task 1",
		Description: "First task",
		AC:          []string{"AC 1"},
	})
	is.NoErr(err)

	// Manually create a duplicate file with the same ID
	duplicateContent := `---
id: "01"
title: Duplicate Task
status: todo
created_at: 2024-01-02T10:00:00Z
---

## Description

This is a duplicate task with the same ID.

## Acceptance Criteria

- [ ] #1 Duplicate AC
`
	err = afero.WriteFile(fs, ".backlog/T01-duplicate_task.md", []byte(duplicateContent), 0o644)
	is.NoErr(err)

	detector := NewConflictDetector(fs, ".backlog")
	conflicts, err := detector.DetectConflicts()
	is.NoErr(err)
	is.Equal(len(conflicts), 1)

	conflict := conflicts[0]
	is.Equal(conflict.Type, ConflictTypeDuplicateID)
	is.Equal(conflict.ConflictID.String(), "01")
	is.Equal(len(conflict.Files), 2)
	is.Equal(len(conflict.Tasks), 2)
	is.True(strings.Contains(conflict.Description, "Task ID 01 appears in multiple files"))
}

func TestConflictDetector_DetectConflicts_OrphanedChild(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()

	// Create a task file that references a non-existent parent
	orphanContent := `---
id: "01.01"
title: Orphaned Task
status: todo
parent: "01"
created_at: 2024-01-01T10:00:00Z
---

## Description

This task references a non-existent parent.

## Acceptance Criteria

- [ ] #1 Orphan AC
`
	err := afero.WriteFile(fs, ".backlog/T01.01-orphaned_task.md", []byte(orphanContent), 0o644)
	is.NoErr(err)

	detector := NewConflictDetector(fs, ".backlog")
	conflicts, err := detector.DetectConflicts()
	is.NoErr(err)
	is.Equal(len(conflicts), 1)

	conflict := conflicts[0]
	is.Equal(conflict.Type, ConflictTypeOrphanedChild)
	is.Equal(conflict.ConflictID.String(), "01.01")
	is.True(strings.Contains(conflict.Description, "references non-existent parent"))
}

func TestConflictDetector_DetectConflicts_InvalidHierarchy(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()

	// Create a task with invalid parent hierarchy
	invalidHierarchyContent := `---
id: "01.02"
title: Invalid Hierarchy
status: todo
parent: "02"
created_at: 2024-01-01T10:00:00Z
---

## Description

This task has the wrong parent based on its ID structure.

## Acceptance Criteria

- [ ] #1 Invalid hierarchy AC
`
	err := afero.WriteFile(fs, ".backlog/T01.02-invalid_hierarchy.md", []byte(invalidHierarchyContent), 0o644)
	is.NoErr(err)

	detector := NewConflictDetector(fs, ".backlog")
	conflicts, err := detector.DetectConflicts()
	is.NoErr(err)
	is.Equal(len(conflicts), 2) // Both orphaned child AND invalid hierarchy

	// Should find both conflicts for the same task
	conflictTypes := make(map[ConflictType]bool)
	for _, conflict := range conflicts {
		conflictTypes[conflict.Type] = true
		is.Equal(conflict.ConflictID.String(), "01.02") // Both conflicts are for the same task
	}

	is.True(conflictTypes[ConflictTypeOrphanedChild])
	is.True(conflictTypes[ConflictTypeInvalidHierarchy])
}

func TestConflictDetector_DetectConflicts_MultipleConflictTypes(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	store := NewFileTaskStore(fs, ".backlog")

	// Create a normal task
	_, err := store.Create(CreateTaskParams{
		Title:       "Task 1",
		Description: "First task",
		AC:          []string{"AC 1"},
	})
	is.NoErr(err)

	// Create duplicate ID conflict
	duplicateContent := `---
id: "01"
title: Duplicate Task
status: todo
created_at: 2024-01-02T10:00:00Z
---

## Description

Duplicate task.
`
	err = afero.WriteFile(fs, ".backlog/T01-duplicate.md", []byte(duplicateContent), 0o644)
	is.NoErr(err)

	// Create orphaned child conflict
	orphanContent := `---
id: "02.01"
title: Orphaned Task
status: todo
parent: "02"
created_at: 2024-01-01T10:00:00Z
---

## Description

Orphaned task.
`
	err = afero.WriteFile(fs, ".backlog/T02.01-orphaned.md", []byte(orphanContent), 0o644)
	is.NoErr(err)

	// Create invalid hierarchy conflict
	invalidContent := `---
id: "03.01"
title: Invalid Hierarchy
status: todo
parent: "04"
created_at: 2024-01-01T10:00:00Z
---

## Description

Invalid hierarchy.
`
	err = afero.WriteFile(fs, ".backlog/T03.01-invalid.md", []byte(invalidContent), 0o644)
	is.NoErr(err)

	detector := NewConflictDetector(fs, ".backlog")
	conflicts, err := detector.DetectConflicts()
	is.NoErr(err)
	is.Equal(len(conflicts), 4) // 1 duplicate + 1 orphaned + 2 for invalid hierarchy task (orphaned + invalid)

	// Verify we have all conflict types
	conflictTypes := make(map[ConflictType]bool)
	for _, conflict := range conflicts {
		conflictTypes[conflict.Type] = true
	}
	is.True(conflictTypes[ConflictTypeDuplicateID])
	is.True(conflictTypes[ConflictTypeOrphanedChild])
	is.True(conflictTypes[ConflictTypeInvalidHierarchy])
}

func TestSummarizeConflicts(t *testing.T) {
	is := is.New(t)

	conflicts := []IDConflict{
		{Type: ConflictTypeDuplicateID, ConflictID: mustParseTaskID("01")},
		{Type: ConflictTypeDuplicateID, ConflictID: mustParseTaskID("02")},
		{Type: ConflictTypeOrphanedChild, ConflictID: mustParseTaskID("03.01")},
		{Type: ConflictTypeInvalidHierarchy, ConflictID: mustParseTaskID("04.01")},
	}

	summary := SummarizeConflicts(conflicts)
	is.Equal(summary.TotalConflicts, 4)
	is.Equal(summary.DuplicateIDs, 2)
	is.Equal(summary.OrphanedChildren, 1)
	is.Equal(summary.InvalidHierarchy, 1)
	is.Equal(len(summary.ConflictsByType[ConflictTypeDuplicateID]), 2)
	is.Equal(len(summary.ConflictsByType[ConflictTypeOrphanedChild]), 1)
	is.Equal(len(summary.ConflictsByType[ConflictTypeInvalidHierarchy]), 1)
}

func TestNewConflictResolver(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	detector := NewConflictDetector(fs, ".backlog")
	store := NewFileTaskStore(fs, ".backlog")

	resolver := NewConflictResolver(detector, store)
	is.True(resolver != nil)
	is.Equal(resolver.detector, detector)
	is.Equal(resolver.store, store)
}

func TestConflictResolver_CreateResolutionPlan_ChronologicalStrategy(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	store := NewFileTaskStore(fs, ".backlog")
	detector := NewConflictDetector(fs, ".backlog")
	resolver := NewConflictResolver(detector, store)

	// Create the tasks directory and some tasks so getNextTaskID works
	err := fs.MkdirAll(".backlog", 0o755)
	is.NoErr(err)
	_, err = store.Create(CreateTaskParams{Title: "Existing Task", Description: "test", AC: []string{"test"}})
	is.NoErr(err)

	// Create tasks with timestamps to test chronological ordering
	baseTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	task1 := Task{
		ID:        mustParseTaskID("01"),
		Title:     "First Task",
		CreatedAt: baseTime,
	}
	task2 := Task{
		ID:        mustParseTaskID("01"),
		Title:     "Second Task",
		CreatedAt: baseTime.Add(time.Hour), // Created later
	}

	conflicts := []IDConflict{
		{
			Type:       ConflictTypeDuplicateID,
			ConflictID: mustParseTaskID("01"),
			Tasks:      []Task{task1, task2},
			Files:      []string{".backlog/T01-first.md", ".backlog/T01-second.md"},
		},
	}

	plan, err := resolver.CreateResolutionPlan(conflicts, ResolutionStrategyChronological)
	is.NoErr(err)
	is.True(plan != nil)
	is.Equal(plan.Strategy, ResolutionStrategyChronological)
	is.Equal(len(plan.Actions), 1) // Should renumber the newer task
	is.Equal(plan.Actions[0].Type, "renumber")
	is.Equal(plan.Actions[0].OriginalID.String(), "01")
}

func TestConflictResolver_CreateResolutionPlan_AutoRenumberStrategy(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	store := NewFileTaskStore(fs, ".backlog")
	detector := NewConflictDetector(fs, ".backlog")
	resolver := NewConflictResolver(detector, store)

	// Create the tasks directory and some tasks so getNextTaskID works
	err := fs.MkdirAll(".backlog", 0o755)
	is.NoErr(err)
	_, err = store.Create(CreateTaskParams{Title: "Existing Task", Description: "test", AC: []string{"test"}})
	is.NoErr(err)

	task1 := Task{ID: mustParseTaskID("01"), Title: "First Task"}
	task2 := Task{ID: mustParseTaskID("01"), Title: "Second Task"}

	conflicts := []IDConflict{
		{
			Type:       ConflictTypeDuplicateID,
			ConflictID: mustParseTaskID("01"),
			Tasks:      []Task{task1, task2},
			Files:      []string{".backlog/T01-first.md", ".backlog/T01-second.md"},
		},
		{
			Type:       ConflictTypeOrphanedChild,
			ConflictID: mustParseTaskID("02.01"),
			Tasks:      []Task{{ID: mustParseTaskID("02.01"), Parent: mustParseTaskID("02")}},
			Files:      []string{".backlog/T02.01-orphan.md"},
		},
	}

	plan, err := resolver.CreateResolutionPlan(conflicts, ResolutionStrategyAutoRenumber)
	is.NoErr(err)
	is.True(plan != nil)
	is.Equal(plan.Strategy, ResolutionStrategyAutoRenumber)
	is.Equal(len(plan.Actions), 2) // One renumber, one update_parent
}

func TestConflictResolver_CreateResolutionPlan_ManualStrategy(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	store := NewFileTaskStore(fs, ".backlog")
	detector := NewConflictDetector(fs, ".backlog")
	resolver := NewConflictResolver(detector, store)

	conflicts := []IDConflict{
		{
			Type:       ConflictTypeDuplicateID,
			ConflictID: mustParseTaskID("01"),
		},
	}

	plan, err := resolver.CreateResolutionPlan(conflicts, ResolutionStrategyManual)
	is.NoErr(err)
	is.True(plan != nil)
	is.Equal(plan.Strategy, ResolutionStrategyManual)
	is.Equal(len(plan.Actions), 1)
	is.Equal(plan.Actions[0].Type, "manual")
}

func TestConflictResolver_CreateResolutionPlan_UnsupportedStrategy(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	store := NewFileTaskStore(fs, ".backlog")
	detector := NewConflictDetector(fs, ".backlog")
	resolver := NewConflictResolver(detector, store)

	conflicts := []IDConflict{}
	_, err := resolver.CreateResolutionPlan(conflicts, ResolutionStrategy(999))
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), "unsupported resolution strategy"))
}

func TestConflictResolver_ExecuteResolutionPlan_DryRun(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	store := NewFileTaskStore(fs, ".backlog")
	detector := NewConflictDetector(fs, ".backlog")
	resolver := NewConflictResolver(detector, store)

	plan := &ResolutionPlan{
		Actions: []ResolutionAction{
			{
				Type:        "renumber",
				OriginalID:  mustParseTaskID("01"),
				NewID:       mustParseTaskID("02"),
				Description: "Test renumber",
			},
			{
				Type:        "update_parent",
				OriginalID:  mustParseTaskID("03.01"),
				NewID:       ZeroTaskID,
				Description: "Test update parent",
			},
			{
				Type:        "manual",
				Description: "Test manual",
			},
		},
	}

	results, err := resolver.ExecuteResolutionPlan(plan, true)
	is.NoErr(err)
	is.Equal(len(results), 4) // 3 actions + dry run notice
	is.True(strings.Contains(results[0], "DRY RUN MODE"))
	is.True(strings.Contains(results[1], "WOULD RENUMBER"))
	is.True(strings.Contains(results[2], "WOULD REMOVE PARENT"))
	is.True(strings.Contains(results[3], "MANUAL"))
}

func TestNewReferenceUpdater(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	detector := NewConflictDetector(fs, ".backlog")
	store := NewFileTaskStore(fs, ".backlog")

	updater := NewReferenceUpdater(detector, store)
	is.True(updater != nil)
	is.Equal(updater.detector, detector)
	is.Equal(updater.store, store)
}

func TestReferenceUpdater_UpdateReferences_EmptyChanges(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	detector := NewConflictDetector(fs, ".backlog")
	store := NewFileTaskStore(fs, ".backlog")
	updater := NewReferenceUpdater(detector, store)

	err := updater.UpdateReferences(map[string]TaskID{})
	is.NoErr(err) // Should succeed with no changes
}

func TestReferenceUpdater_FindTaskReferences(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	store := NewFileTaskStore(fs, ".backlog")

	// Create parent task
	parent, err := store.Create(CreateTaskParams{
		Title:       "Parent Task",
		Description: "Parent",
		AC:          []string{"Parent AC"},
	})
	is.NoErr(err)

	// Create child task
	childParams := CreateTaskParams{
		Title:       "Child Task",
		Description: "Child",
		Parent:      parent.ID.String(),
		AC:          []string{"Child AC"},
	}
	child, err := store.Create(childParams)
	is.NoErr(err)

	// Create task with dependency
	depTask, err := store.Create(CreateTaskParams{
		Title:       "Dependent Task",
		Description: "Depends on parent",
		AC:          []string{"Dep AC"},
	})
	is.NoErr(err)

	// Add dependency
	editParams := EditTaskParams{
		ID:              depTask.ID.String(),
		NewDependencies: []string{parent.ID.String()},
	}
	is.NoErr(store.Update(&depTask, editParams))

	detector := NewConflictDetector(fs, ".backlog")
	updater := NewReferenceUpdater(detector, store)

	// Find references to parent task
	references, err := updater.FindTaskReferences(parent.ID)
	is.NoErr(err)
	is.Equal(len(references), 2) // Child and dependent task

	// Verify references
	foundChild := false
	foundDep := false
	for _, ref := range references {
		if ref.ID.Equals(child.ID) {
			foundChild = true
		}
		if ref.ID.Equals(depTask.ID) {
			foundDep = true
		}
	}
	is.True(foundChild)
	is.True(foundDep)
}

func TestMaybeStringArrayFromSlice(t *testing.T) {
	is := is.New(t)

	slice := []string{"a", "b", "c"}
	msa := MaybeStringArrayFromSlice(slice)
	is.Equal(len(msa), 3)
	is.Equal(msa.ToSlice(), slice)

	// Test empty slice
	emptyMsa := MaybeStringArrayFromSlice([]string{})
	is.Equal(len(emptyMsa), 0)
	is.Equal(len(emptyMsa.ToSlice()), 0)
}

func TestConflictDetector_ParseTaskFromFile_InvalidFile(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	detector := NewConflictDetector(fs, ".backlog")

	// Test non-existent file
	_, err := detector.parseTaskFromFile("nonexistent.md")
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), "failed to read file"))

	// Test invalid task content
	invalidContent := "invalid yaml content"
	err = afero.WriteFile(fs, ".backlog/invalid.md", []byte(invalidContent), 0o644)
	is.NoErr(err)

	_, err = detector.parseTaskFromFile(".backlog/invalid.md")
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), "failed to parse task"))
}

func TestConflictDetector_DetectConflicts_EmptyDirectory(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()

	// Create empty tasks directory
	err := fs.MkdirAll(".backlog", 0o755)
	is.NoErr(err)

	detector := NewConflictDetector(fs, ".backlog")
	conflicts, err := detector.DetectConflicts()
	is.NoErr(err)
	is.Equal(len(conflicts), 0)
}

func TestConflictDetector_DetectConflicts_NonExistentDirectory(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()

	detector := NewConflictDetector(fs, ".nonexistent")
	_, err := detector.DetectConflicts()
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), "failed to read tasks directory"))
}

func TestConflictDetector_DetectConflicts_IgnoreNonTaskFiles(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()

	// Create tasks directory
	err := fs.MkdirAll(".backlog", 0o755)
	is.NoErr(err)

	// Create non-task files that should be ignored
	err = afero.WriteFile(fs, ".backlog/README.md", []byte("readme"), 0o644)
	is.NoErr(err)
	err = afero.WriteFile(fs, ".backlog/not-a-task.txt", []byte("text"), 0o644)
	is.NoErr(err)
	err = afero.WriteFile(fs, ".backlog/config.json", []byte("{}"), 0o644)
	is.NoErr(err)

	// Create subdirectory which should be ignored
	err = fs.MkdirAll(".backlog/subdir", 0o755)
	is.NoErr(err)

	detector := NewConflictDetector(fs, ".backlog")
	conflicts, err := detector.DetectConflicts()
	is.NoErr(err)
	is.Equal(len(conflicts), 0)
}

// Helper function for tests
func mustParseTaskID(id string) TaskID {
	taskID, err := parseTaskID(id)
	if err != nil {
		panic(err)
	}
	return taskID
}

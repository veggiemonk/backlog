---
id: "21"
title: Implement backlog move command
status: todo
labels:
    - feature
    - cli
priority: high
created_at: 2025-10-17T13:55:00.979737Z
---
## Description

Add a 'backlog move <task-id> --parent <parent-id>' command to move tasks in the hierarchy. This feature allows changing a task's parent, which triggers an ID change and file rename. Example: 'backlog move 3 --parent 1' moves T03 to become T01.01 (or next available subtask ID).

## Acceptance Criteria
<!-- AC:BEGIN -->

- [ ] #1 CLI command 'backlog move <task-id> --parent <parent-id>' exists and is properly registered
- [ ] #2 Command validates that parent task exists before proceeding
- [ ] #3 Command detects and prevents cycles (cannot move task to be child of its own descendant)
- [ ] #4 Command calculates correct new ID using parent's next available subtask ID
- [ ] #5 Task file is renamed from old ID format to new ID format (e.g.
- [ ] #6  T03-title.md -> T01.01-title.md)
- [ ] #7 All references to moved task are updated (dependencies in other tasks)
- [ ] #8 Child tasks of moved task have their parent reference preserved/updated correctly
- [ ] #9 Task history records the ID change with old and new IDs
- [ ] #10 Command handles case where target ID already exists with clear error message
- [ ] #11 Moved task retains all its metadata (title
- [ ] #12  description
- [ ] #13  AC
- [ ] #14  status
- [ ] #15  etc.)

<!-- AC:END -->

## Implementation Plan



## Implementation Notes



# Backlog MCP: Task Management for Coding Agents

You are a coding agent with access to a task management system called **backlog** via MCP (Model Context Protocol). This system helps you organize, track, and complete development work efficiently.

## Critical Rules - READ THIS FIRST

⚠️ **NEVER directly read, write, modify, or delete task files in the `.backlog/` directory**

✅ **ALWAYS use MCP tools for ALL task operations**: `task_create`, `task_edit`, `task_view`, `task_list`, `task_search`, `task_archive`

❌ **NEVER**: Edit markdown files directly, change checkboxes manually, or bypass the MCP tools

**Why?** Direct file editing breaks metadata synchronization, Git tracking, and task relationships.

## Task Assignment & Work Protocol

### 1. When Starting Work on ANY Task

**Before you write a single line of code**, you MUST:

```json
{
  "name": "task_edit",
  "arguments": {
    "id": "TASK_ID",
    "new_status": "in-progress",
    "new_assigned": ["@coding-agent"]
  }
}
```

Replace `@coding-agent` with your actual agent identifier.

### 2. Create Implementation Plan

After claiming a task, document HOW you will implement it:

```json
{
  "name": "task_edit", 
  "arguments": {
    "id": "TASK_ID",
    "new_plan": "1. Analyze existing code\n2. Implement feature X\n3. Add tests\n4. Update documentation"
  }
}
```

### 3. Track Progress

Mark acceptance criteria as complete as you finish them:

```json
{
  "name": "task_edit",
  "arguments": {
    "id": "TASK_ID", 
    "check_ac": [1, 2]  // Mark criteria #1 and #2 as done
  }
}
```

### 4. Complete the Task

When finished, add implementation notes and mark as done:

```json
// Add implementation summary
{
  "name": "task_edit",
  "arguments": {
    "id": "TASK_ID",
    "new_notes": "Implemented using pattern X. Modified files A, B, C. All tests pass."
  }
}

// Mark as complete
{
  "name": "task_edit", 
  "arguments": {
    "id": "TASK_ID",
    "new_status": "done"
  }
}
```

## Essential MCP Tools

### Finding Work
```json
// Find unassigned tasks
{"name": "task_list", "arguments": {"unassigned": true, "status": "todo"}}

// Find your assigned tasks  
{"name": "task_list", "arguments": {"assigned": "@coding-agent", "status": "in-progress"}}

// Search for specific topics
{"name": "task_search", "arguments": {"query": "authentication API"}}
```

## Advanced Task Operations

### Batch Task Creation

When you need to create multiple related tasks (like breaking down a large feature), use `task_batch_create`:

```json
{
  "name": "task_batch_create",
  "arguments": {
    "new_tasks": [
      {
        "title": "Design API endpoints",
        "description": "Define REST API structure for user management",
        "ac": ["Document endpoints", "Review with team"],
        "labels": ["api", "design"],
        "priority": "high"
      },
      {
        "title": "Implement user authentication", 
        "description": "Add JWT-based authentication system",
        "ac": ["Add login endpoint", "Add token validation", "Add tests"],
        "labels": ["api", "auth"],
        "priority": "high",
        "dependencies": ["T15"]
      },
      {
        "title": "Update documentation",
        "description": "Document the new API endpoints",
        "ac": ["API docs", "Usage examples"],
        "labels": ["docs"],
        "priority": "medium"
      }
    ]
  }
}
```

**Benefits of batch creation:**
- All tasks are created atomically
- Maintains consistent metadata across related tasks
- More efficient than multiple `task_create` calls
- Useful for project planning and feature breakdowns

### Advanced Listing and Search

#### Filtering Options

```json
// Multiple status filters
{"name": "task_list", "arguments": {"status": "todo,in-progress"}}

// Multiple assignee filters  
{"name": "task_list", "arguments": {"assigned": "alice,bob"}}

// Multiple label filters
{"name": "task_list", "arguments": {"labels": "bug,critical"}}

// Unassigned tasks only
{"name": "task_list", "arguments": {"unassigned": true}}

// Tasks with dependencies
{"name": "task_list", "arguments": {"has_dependency": true}}

// Tasks that block others (depended on)
{"name": "task_list", "arguments": {"depended_on": true}}

// Parent/child relationships
{"name": "task_list", "arguments": {"parent": "T15"}}
```

#### Sorting and Display

```json
// Sort by priority (high to low)
{"name": "task_list", "arguments": {"sort": ["priority"], "reverse": true}}

// Sort by multiple fields
{"name": "task_list", "arguments": {"sort": ["status", "priority", "updated"], "reverse": true}}

// Hide extra metadata for cleaner display
{"name": "task_list", "arguments": {"hide_extra": true}}
```

#### Pagination for Large Projects

```json
// Get first 10 tasks
{"name": "task_list", "arguments": {"limit": 10}}

// Get next 10 tasks (second page)
{"name": "task_list", "arguments": {"limit": 10, "offset": 10}}

// High-priority todos, first 5
{"name": "task_list", "arguments": {
  "status": "todo",
  "sort": ["priority"], 
  "reverse": true,
  "limit": 5
}}
```

#### Search with Filters

```json
// Search with status filter
{"name": "task_search", "arguments": {
  "query": "authentication", 
  "filters": {"status": "todo", "limit": 5}
}}

// Search with pagination
{"name": "task_search", "arguments": {
  "query": "API bug",
  "filters": {"labels": "bug", "limit": 3, "offset": 0}
}}

// Search assigned work
{"name": "task_search", "arguments": {
  "query": "refactor",
  "filters": {"assigned": "@coding-agent"}
}}
```

#### Practical Query Patterns

```json
// Find my urgent work
{"name": "task_list", "arguments": {
  "assigned": "@coding-agent",
  "status": "in-progress",
  "sort": ["priority"],
  "reverse": true
}}

// Find blocked tasks I can work on
{"name": "task_list", "arguments": {
  "unassigned": true,
  "status": "todo", 
  "has_dependency": false,
  "sort": ["priority"],
  "reverse": true,
  "limit": 5
}}

// Review what's blocking others
{"name": "task_list", "arguments": {
  "assigned": "@coding-agent",
  "depended_on": true,
  "status": "in-progress"
}}
```

#### Pagination Response Format

When using `limit`, responses include pagination metadata:

```json
{
  "tasks": [...],
  "pagination": {
    "total_results": 45,
    "displayed_results": 10,
    "offset": 0,
    "limit": 10,
    "has_more": true
  }
}
```

Use `has_more` to determine if additional pages exist.

### Reading Tasks
```json
// Get detailed view of a task
{"name": "task_view", "arguments": {"id": "42"}}

// List tasks with filters
{"name": "task_list", "arguments": {"status": "todo", "labels": "bug,critical"}}
```

### Creating Tasks
```json
{
  "name": "task_create",
  "arguments": {
    "title": "Fix authentication bug",
    "description": "Users can't log in with special characters in password",
    "ac": ["Reproduce the bug", "Fix the issue", "Add test cases"],
    "labels": ["bug", "authentication"],
    "priority": "high"
  }
}
```

### Updating Tasks (Most Important)
```json
{
  "name": "task_edit",
  "arguments": {
    "id": "42",
    "new_status": "in-progress",           // Change status
    "new_assigned": ["@coding-agent"],     // Assign to yourself
    "check_ac": [1, 3],                    // Mark criteria as done
    "new_notes": "Implementation details", // Add notes
    "new_plan": "Step-by-step approach"    // Document approach
  }
}
```

## Definition of Done

A task is complete ONLY when:

✅ **All acceptance criteria are checked** via `check_ac`
✅ **Implementation notes added** via `new_notes` 
✅ **Status set to "done"** via `new_status`
✅ **All tests pass** (run your test commands)
✅ **Code is reviewed** (self-review your changes)

## Workflow Example

```json
// 1. Find available work
{"name": "task_list", "arguments": {"status": "todo", "limit": 5}}

// 2. Claim a task
{"name": "task_edit", "arguments": {"id": "15", "new_status": "in-progress", "new_assigned": ["@coding-agent"]}}

// 3. Add implementation plan
{"name": "task_edit", "arguments": {"id": "15", "new_plan": "1. Research API\n2. Write code\n3. Test"}}

// 4. Work on the code (your actual development work)

// 5. Mark progress
{"name": "task_edit", "arguments": {"id": "15", "check_ac": [1, 2]}}

// 6. Finish and document
{"name": "task_edit", "arguments": {"id": "15", "new_notes": "Completed feature X using approach Y", "new_status": "done"}}
```

## Common Mistakes to Avoid

❌ Starting code work without assigning the task to yourself
❌ Forgetting to update task status to "in-progress" 
❌ Not documenting your implementation approach
❌ Marking tasks as done without checking all acceptance criteria
❌ Bypassing MCP tools to edit task files directly

## Quick Reference

| Need to... | Use Tool | Example |
|------------|----------|---------|
| Find work | `task_list` | `{"unassigned": true}` |
| Start work | `task_edit` | `{"id": "X", "new_status": "in-progress", "new_assigned": ["@me"]}` |
| Check progress | `task_edit` | `{"id": "X", "check_ac": [1,2]}` |
| Finish task | `task_edit` | `{"id": "X", "new_status": "done"}` |

Remember: The MCP tools are your interface to the task system. Use them consistently and you'll have a well-organized development workflow with full traceability of your work.

# Instructions for the usage of Backlog MCP Tools


## Backlog.md: Comprehensive Project Management via MCP Tools

### Assistant Objective

Efficiently manage all project tasks, status, and documentation using the Backlog MCP (Model-Controlled Process) tools, ensuring all project metadata remains fully synchronized and up-to-date.

### Core Capabilities

- ✅ **Task Management**: Create, edit, assign, prioritize, and track tasks with full metadata
- ✅ **Acceptance Criteria**: Granular control with add/remove/check/uncheck by index
- ✅ **Git Integration**: Automatic tracking of task states across branches
- ✅ **Dependencies**: Task relationships and subtask hierarchies
- ✅ **AI-Optimized**: Tools return structured data perfect for AI processing

### Why This Matters to You (AI Agent)

1.  **Comprehensive system** - Full project management capabilities through MCP tools.
2.  **The tools are the interface** - All operations go through `task_*` tool calls.
3.  **Unified interaction model** - You can use tools for both reading (`task_view`) and writing (`task_edit`).
4.  **Metadata stays synchronized** - The tools handle all the complex relationships.

### Key Understanding

- **Tasks** live in `.backlog/` as markdown files.
- **You interact via MCP tools only**: `task_create`, `task_edit`, etc.
- **Never bypass the tools** - They handle Git, metadata, file naming, and relationships.

---

# ⚠️ CRITICAL: NEVER EDIT OR DELETE TASK FILES DIRECTLY. Edit Only via MCP Tools

**ALL task operations MUST use the Backlog MCP tool calls.**

- ✅ **DO**: Use `task_edit` and other `task_*` tools.
- ✅ **DO**: Use `task_create` to create new tasks.
- ✅ **DO**: Use `task_edit(id=..., check_ac=[1])` to mark acceptance criteria.
- ❌ **DON'T**: Edit markdown files directly.
- ❌ **DON'T**: Manually change checkboxes in files.
- ❌ **DON'T**: Add or modify text in task files without using the tools.

**Why?** Direct file editing breaks metadata synchronization, Git tracking, and task relationships.

---

## 1. Source of Truth & File Structure

### 📖 **UNDERSTANDING** (What you'll see when reading files)

- Markdown task files live under **`.backlog/`**.
- Files are named using a convention like: `T01.02-my-task-title.md`.
- You DO NOT need to interact with the file system directly for task management.

### 🔧 **ACTING** (How to change things)

- **All task operations MUST use the `task_*` MCP tools.**
- This ensures metadata is correctly updated and the project stays in sync.
- The tools return structured data, so you don't need to parse files.

---

## 2. Common Mistakes to Avoid

### ❌ **WRONG: Direct File Editing**

```python
# DON'T DO THIS:

# 1. Read .backlog/T07-feature.md
# 2. Manually change "- [ ]" to "- [x]" in the content
# 3. Write the modified content back to the file
```

### ✅ **CORRECT: Using MCP Tools**

```json
// DO THIS INSTEAD:

// Mark AC #1 as complete
{
  "name": "task_edit",
  "arguments": {
    "id": "7",
    "check_ac": [1]
  }
}

// Add notes
{
  "name": "task_edit",
  "arguments": {
    "id": "7",
    "new_notes": "Implementation complete"
  }
}

// Multiple changes: change status and assign the task
{
  "name": "task_edit",
  "arguments": {
    "id": "7",
    "new_status": "in-progress",
    "new_assigned": ["@agent-k"]
  }
}
```

---

## 3. Understanding Task Format (Read-Only Reference)

⚠️ **FORMAT REFERENCE ONLY** - The following shows the structure of the underlying data.
**Never edit files directly! Use MCP tools to make changes.**

### Task Structure

```yaml
---
id: "42"
title: "Add GraphQL resolver"
status: "todo"
assigned: ["@sara"]
labels: ["backend", "api"]
---

## Description

Brief explanation of the task purpose.

## Acceptance Criteria

<!-- AC:BEGIN -->
- [ ] #1 First criterion
- [x] #2 Second criterion (completed)
<!-- AC:END -->
```

### How to Modify Each Section

| What You Want to Change | MCP Tool Call Arguments                                   |
| ----------------------- | ---------------------------------------------------------- |
| Title                   | `{"id": "42", "new_title": "New Title"}`                |
| Status                  | `{"id": "42", "new_status": "in-progress"}`             |
| Assigned                | `{"id": "42", "new_assigned": ["@sara"]}`               |
| Labels                  | `{"id": "42", "new_labels": ["backend", "api"]}`        |
| Description             | `{"id": "42", "new_description": "New description"}`    |
| Add AC                  | `{"id": "42", "add_ac": ["New criterion"]}`             |
| Check AC #1             | `{"id": "42", "check_ac": [1]}`                          |
| Uncheck AC #2           | `{"id": "42", "uncheck_ac": [2]}`                        |
| Remove AC #1            | `{"id": "42", "remove_ac": [1]}`                         |
| Add Plan                | `{"id": "42", "new_plan": "1. Step one \n2. Step two"}`  |
| Add Notes               | `{"id": "42", "new_notes": "What I did"}`               |
| Remove Assigned User    | `{"id": "42", "remove_assigned": ["@sara"]}`             |
| Remove Labels           | `{"id": "42", "remove_labels": ["backend", "api"]}`     |

---

## 4. Defining Tasks

### Creating New Tasks

**Always use the `task_create` tool:**

```json
// Example
{
  "name": "task_create",
  "arguments": {
    "title": "Task title",
    "description": "Description of the task.",
    "ac": ["First criterion", "Second criterion"]
  }
}
```

### Acceptance Criteria (The "what")

**Managing Acceptance Criteria via Tools:**

- **Adding criteria (`add_ac`)** takes a list of strings.
- **Checking/unchecking/removing (`check_ac`, `uncheck_ac`, `remove_ac`)** take a list of 1-based indices.
- You can perform multiple operations in a single call.

```json
// Examples

// Add new criteria
{
  "name": "task_edit",
  "arguments": {
    "id": "42",
    "add_ac": ["User can login", "Session persists"]
  }
}

// Check multiple criteria by index
{
  "name": "task_edit",
  "arguments": {
    "id": "42",
    "check_ac": [1, 2, 3]
  }
}

// Uncheck a criterion
{
  "name": "task_edit",
  "arguments": {
    "id": "42",
    "uncheck_ac": [2]
  }
}

// Remove multiple criteria
{
  "name": "task_edit",
  "arguments": {
    "id": "42",
    "remove_ac": [2, 4]
  }
}
// Note: Indices are processed high-to-low

// Mixed operations in a single command
{
  "name": "task_edit",
  "arguments": {
    "id": "42",
    "check_ac": [1],
    "uncheck_ac": [2],
    "remove_ac": [3],
    "add_ac": ["New criterion"]
  }
}
```

---

## 5. Implementing Tasks

### 5.1. First step when implementing a task

The very first things you must do when you take over a task are to set the task to "In Progress" and assign it to yourself.

```json
// Example
{
  "name": "task_edit",
  "arguments": {
    "id": "42",
    "new_status": "in-progress",
    "new_assigned": ["@{myself}"]
  }
}
```

### 5.2. Create an Implementation Plan (The "how")

Once you are familiar with the task, create a plan on **HOW** to tackle it. Write it down in the task so that you can refer to it later.

```json
// Example
{
  "name": "task_edit",
  "arguments": {
    "id": "42",
    "new_plan": "1. Research codebase for references\n2. Research on internet for similar cases\n3. Implement\n4. Test"
  }
}
```

### 5.3. Implementation Notes (PR description)

When you are done implementing a task, write a clean description in the task notes, as if it were a PR description.

```json
// Example
{
  "name": "task_edit",
  "arguments": {
    "id": "42",
    "new_notes": "Implemented using pattern X because of Reason Y. Modified files Z and W."
  }
}
```

---

## 6. Typical Workflow

```json
// 1. Identify work
{"name": "task_list", "arguments": {"status": "todo"}}
{"name": "task_list", "arguments": {"status": "todo,in-progress"}}  // Multiple statuses
{"name": "task_list", "arguments": {"unassigned": true}}  // Find tasks needing assignment
{"name": "task_list", "arguments": {"assigned": "alice"}}  // Tasks assigned to specific person
{"name": "task_list", "arguments": {"assigned": "alice,bob"}}  // Tasks assigned to alice OR bob
{"name": "task_list", "arguments": {"has_dependency": true}}  // Tasks waiting on dependencies
{"name": "task_list", "arguments": {"depended_on": true, "status": "todo"}}  // Blocking tasks
{"name": "task_list", "arguments": {"labels": "bug,critical"}}  // Tasks with specific labels
{"name": "task_list", "arguments": {"status": "todo", "sort": ["priority"], "reverse": true}}  // High priority first

// Pagination examples
{"name": "task_list", "arguments": {"limit": 5}}  // Get first 5 tasks
{"name": "task_list", "arguments": {"status": "todo", "limit": 10}}  // First 10 todo tasks
{"name": "task_search", "arguments": {"query": "feature", "filters": {"limit": 3}}}  // First 3 feature matches

// 2. Read task details
{"name": "task_view", "arguments": {"id": "42"}}

// 3. Start work: assign yourself & change status
{
  "name": "task_edit",
  "arguments": {
    "id": "42",
    "new_status": "in-progress",
    "new_assigned": ["@myself"]
  }
}

// 4. Add implementation plan
{
  "name": "task_edit",
  "arguments": {
    "id": "42",
    "new_plan": "1. Analyze\n2. Refactor\n3. Test"
  }
}

// 5. Work on the task (write code, test, etc.)

// 6. Mark acceptance criteria as complete
{
  "name": "task_edit",
  "arguments": {
    "id": "42",
    "check_ac": [1, 2, 3]
  }
}  // Check all at once

// 7. Add implementation notes (PR Description)
{
  "name": "task_edit",
  "arguments": {
    "id": "42",
    "new_notes": "Refactored using strategy pattern, updated tests."
  }
}

// 8. Mark task as done
{
  "name": "task_edit",
  "arguments": {
    "id": "42",
    "new_status": "done"
  }
}
```

---

## 7. Definition of Done (DoD)

A task is **Done** only when **ALL** of the following are complete:

### ✅ Via MCP Tools:

1.  **All acceptance criteria checked**: Use MCP tool with `{"id": "...", "check_ac": [...]}`.
2.  **Implementation notes added**: Use MCP tool with `{"id": "...", "new_notes": "..."}`.
3.  **Status set to Done**: Use MCP tool with `{"id": "...", "new_status": "done"}`.

### ✅ Via Code/Testing:

4.  **Tests pass**: Run test suite and linting.
5.  **Documentation updated**: Update relevant docs if needed.
6.  **Code reviewed**: Self-review your changes.
7.  **No regressions**: Performance, security checks pass.

⚠️ **NEVER mark a task as Done without completing ALL items above.**

---

## 8. Quick Reference: DO vs DON'T

### Viewing Tasks

| Task       | ✅ DO                                          | ❌ DON'T                         |
| ---------- | ---------------------------------------------- | -------------------------------- |
| View task  | Use `task_view` tool with `{"id": "42"}`       | Open and read .md file directly  |
| List tasks | Use `task_list` tool with `{"status": "todo"}` | Browse the `.backlog` folder     |
| List Tasks | Use `task_list` tool with `{"unassigned": true}` | Browse the `.backlog` folder     |
| List Tasks | Use `task_list` tool with `{"assigned": "alice"}` | Browse the `.backlog` folder     |

### Modifying Tasks

| Task          | ✅ DO                                          | ❌ DON'T                           |
| ------------- | ---------------------------------------------- | ---------------------------------- |
| Check AC      | Use `task_edit` tool with `{"id": "42", "check_ac": [1]}` | Change `- [ ]` to `- [x]` in file  |
| Add notes     | Use `task_edit` tool with `{"id": "42", "new_notes": "..."}` | Type notes into .md file           |
| Change status | Use `task_edit` tool with `{"id": "42", "new_status": "done"}` | Edit status in frontmatter         |
| Add AC        | Use `task_edit` tool with `{"id": "42", "add_ac": ["New criterion"]}` | Add `- [ ] New` to file            |
| Archive task  | Use `task_archive` tool with `{"id": "42"}` | Manually move files to archive folder |

---

## 9. Complete MCP Tool Reference

### `task_create`

Creates a new task. All the fields should be filled.

| Parameter     | Type          | Description                               |
| ------------- | ------------- | ----------------------------------------- |
| `title`       | `string`      | **Required.** The title of the task.      |
| `description` | `string`      | A detailed description of the task.       |
| `parent`      | `string`      | The ID of the parent task.                |
| `ac`          | `list[string]`| A list of acceptance criteria.            |
| `assigned`    | `list[string]`| A list of assigned users.                 |
| `labels`      | `list[string]`| A list of labels.                         |
| `priority`    | `string`      | The priority of the task.                 |
| `dependencies`| `list[string]`| A list of task IDs that this task depends on. |

### `task_edit`

Edits an existing task.

| Parameter         | Type          | Description                                       |
| ----------------- | ------------- | ------------------------------------------------- |
| `id`              | `string`      | **Required.** The ID of the task to edit.         |
| `new_title`       | `string`      | A new title for the task.                         |
| `new_description` | `string`      | A new description for the task.                   |
| `new_status`      | `string`      | A new status (e.g., "in-progress", "done").       |
| `new_dependencies`| `list[string]`| A new list of dependencies (replaces the old list). |
| `new_parent`      | `string`      | A new parent task ID.                             |
| `new_assigned`    | `list[string]`| A new list of assigned users (replaces the old list). |
| `new_labels`      | `list[string]`| A new list of labels (replaces the old list).     |
| `new_priority`    | `string`      | A new priority.                                   |
| `remove_assigned` | `list[string]`| A list of assigned users to remove.               |
| `remove_labels`   | `list[string]`| A list of labels to remove.                       |
| `add_ac`          | `list[string]`| A list of new acceptance criteria to add.         |
| `remove_ac`       | `list[int]`   | A list of 1-based indices of AC to remove.        |
| `check_ac`        | `list[int]`   | A list of 1-based indices of AC to check.         |
| `uncheck_ac`      | `list[int]`   | A list of 1-based indices of AC to uncheck.       |
| `new_plan`        | `string`      | A new implementation plan.                        |
| `new_notes`       | `string`      | New implementation notes.                         |

### `task_list`

Lists tasks with optional filtering, sorting, and pagination. When pagination is used, returns structured results with pagination metadata.

| Parameter        | Type          | Description                                           |
| ---------------- | ------------- | ----------------------------------------------------- |
| `status`         | `string`      | Filter tasks by status (comma-separated for multiple). |
| `parent`         | `string`      | Filter tasks by a parent task ID.                     |
| `assigned`       | `string`      | Filter tasks by assigned user (comma-separated for multiple). |
| `unassigned`     | `bool`        | Filter tasks that have no assigned users.             |
| `labels`         | `string`      | Filter tasks by labels (comma-separated for multiple). |
| `has_dependency` | `bool`        | Filter tasks that have dependencies.                  |
| `depended_on`    | `bool`        | Filter tasks that are depended on by other tasks.     |
| `hide_extra`     | `bool`        | Hide extra fields (labels, priority, assigned).       |
| `sort`           | `list[string]`| Fields to sort by (id, title, status, priority, created, updated). |
| `reverse`        | `bool`        | Reverse the sort order.                               |
| `limit`          | `int`         | Maximum number of tasks to return (0 means no limit). |
| `offset`         | `int`         | Number of tasks to skip from the beginning.          |

### `task_view`

Retrieves and displays the details of a single task.

| Parameter | Type     | Description                       |
| --------- | -------- | --------------------------------- |
| `id`      | `string` | **Required.** The ID of the task. |

### `task_search`

Searches tasks by content with optional filtering and pagination. When pagination is used, returns structured results with pagination metadata.

| Parameter | Type     | Description                       |
| --------- | -------- | --------------------------------- |
| `query`   | `string` | **Required.** The search query.   |
| `filters` | `object` | Optional filters (same as task_list parameters including limit/offset). |

### `task_archive`

Archives a task by moving it to the archived directory and setting status to archived.

| Parameter | Type     | Description                       |
| --------- | -------- | --------------------------------- |
| `id`      | `string` | **Required.** The ID of the task to archive. |

---

## 10. Pagination: Handling Large Task Lists

### When to Use Pagination

Use pagination when:
- Working with projects that have many tasks (>25)
- Performing exploratory queries where you want to see a sample first
- Building interfaces that need to display results in pages
- Avoiding overwhelming output in conversations

### Pagination Parameters

Both `task_list` and `task_search` support pagination:

- **`limit`**: Maximum number of results to return (0 = no limit)
- **`offset`**: Number of results to skip from the beginning

### Pagination Examples

```json
// Get first 10 tasks
{"name": "task_list", "arguments": {"limit": 10}}

// Get next 10 tasks (pagination)
{"name": "task_list", "arguments": {"limit": 10, "offset": 10}}

// Get first 5 high-priority tasks
{"name": "task_list", "arguments": {"status": "todo", "sort": ["priority"], "reverse": true, "limit": 5}}

// Search with pagination
{"name": "task_search", "arguments": {"query": "api", "filters": {"limit": 3}}}

// Search second page
{"name": "task_search", "arguments": {"query": "api", "filters": {"limit": 3, "offset": 3}}}
```

### Pagination Response Format

When pagination is used, responses include metadata:

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

### Best Practices

1. **Start with small limits**: Use `limit: 10` to get an overview
2. **Check `has_more`**: Use pagination metadata to determine if more results exist
3. **Progressive exploration**: Increase offset to see more results
4. **Combine with filtering**: Use pagination with status/priority filters for focused results

### Configuration

Users can configure pagination defaults:
- **Environment variables**: `BACKLOG_PAGE_SIZE`, `BACKLOG_MAX_LIMIT`
- **CLI flags**: `--page-size 25`, `--max-limit 1000`

---

## 11. Advanced Workflows: Batch Task Creation

When a user asks you to perform a multi-step operation like creating a full project plan, your goal is to gather all necessary information from the initial prompt and execute the steps in logical order, using parallel tool calls to be efficient.

### Example High-Level Prompt

A user might provide a comprehensive request like this:

> "Here is our refactoring plan in `plan.md`. Please create all the necessary tasks in the backlog.
>
> -   **Assignee**: `agent-cli`
> -   **Priority**: `high`
> -   **Labels**: Please add relevant labels to each task based on its content (e.g., `refactoring`, `cli`, `documentation`)."

### Your Interpretation and Execution Plan

1.  **Deconstruct the Request**: Identify the separate pieces of information provided:
    *   **Source**: The `plan.md` file.
    *   **Action**: Create tasks.
    *   **Metadata**: Assignee (`agent-cli`), Priority (`high`), and instructions for Labels.

2.  **Formulate a Multi-Step Execution Plan**:
    1.  First, create all the tasks and sub-tasks as defined in the plan. This is necessary to obtain the `id` for each new task.
    2.  Once all tasks are created, execute a second set of **parallel** `task_edit` calls to update the metadata (assignee, priority, labels) for every task you just created.

3.  **Execute Efficiently**:

    ```json
    // Step 1: Create all tasks from the plan to get their IDs
    {"name": "task_create", "arguments": {"title": "Parent Task for Refactoring", ...}}
    {"name": "task_create", "arguments": {"title": "Sub-task 1: Update CLI command", "parent": "T21", ...}}
    {"name": "task_create", "arguments": {"title": "Sub-task 2: Update Documentation", "parent": "T21", ...}}

    // Step 2: In parallel, update metadata for all newly created tasks
    {"name": "task_edit", "arguments": {"id": "T21", "new_assignee": ["agent-cli"], "new_priority": "high", "new_labels": ["refactoring", "cli"]}}
    {"name": "task_edit", "arguments": {"id": "T21.01", "new_assignee": ["agent-cli"], "new_priority": "high", "new_labels": ["refactoring", "cli"]}}
    {"name": "task_edit", "arguments": {"id": "T21.02", "new_assignee": ["agent-cli"], "new_priority": "high", "new_labels": ["documentation"]}}
    ```

### Handling Missing Information

If the user's initial prompt is missing key information (like `assignee` or `priority`), you **must ask** for the missing details before proceeding.

**Example Clarification Question:**

> "I can create the tasks from the plan. Could you please tell me what priority I should set for them and who the assignee should be?"

---

## Remember: The Golden Rule

**🎯 If you want to change ANYTHING in a task, use the `task_edit` tool.**
**📖 Use `task_view` and `task_list` to read tasks. Never write to files directly.**

---

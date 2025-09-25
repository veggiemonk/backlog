# Instructions for the usage of Backlog via MCP Server

## Backlog: Comprehensive Project Management via MCP Tools

### Assistant Objective

Efficiently manage all project tasks, status, and documentation using the Backlog MCP server tools, ensuring all project metadata remains fully synchronized and up-to-date.

### Core Capabilities

- ✅ **Task Management**: Create, edit, assign, prioritize, and track tasks with full metadata
- ✅ **Acceptance Criteria**: Granular control with add/remove/check/uncheck operations
- ✅ **Git Integration**: Automatic tracking of task states across branches
- ✅ **Dependencies**: Task relationships and subtask hierarchies
- ✅ **Structured I/O**: All tools use structured data (JSON), perfect for AI processing

### Why This Matters to You (AI Agent)

1.  **Comprehensive system** - Full project management capabilities through MCP tools.
2.  **The MCP server is the interface** - All operations go through MCP tool calls.
3.  **Unified interaction model** - You can use tools for both reading (`task_view`, `task_list`) and writing (`task_create`, `task_edit`).
4.  **Metadata stays synchronized** - The MCP server handles all the complex relationships.

### Key Understanding

- **Tasks** live in `.backlog/` as markdown files.
- **You interact via MCP tools only**: `task_create`, `task_edit`, etc.
- **Never bypass the MCP server** - It handles Git, metadata, file naming, and relationships.

## ⚠️ CRITICAL: NEVER EDIT TASK FILES DIRECTLY OR USE THE CLI. Edit Only via MCP Tools.

**ALL task operations MUST use the Backlog MCP server tools.**

- ✅ **DO**: Use `task_edit` and other MCP tools.
- ✅ **DO**: Use `task_create` to create new tasks.
- ✅ **DO**: Use `task_edit(id="T1", check_ac=[1])` to mark acceptance criteria.
- ❌ **DON'T**: Edit markdown files directly.
- ❌ **DON'T**: Use the `backlog` CLI commands.
- ❌ **DON'T**: Manually change checkboxes in files.
- ❌ **DON'T**: Add or modify text in task files without using the MCP tools.

**Why?** Direct file editing or using the CLI bypasses the controlled MCP environment, breaking metadata synchronization, Git tracking, and task relationships.


## 1. Source of Truth & File Structure

### 📖 **UNDERSTANDING** (What you'll see when reading files)

- Markdown task files live under **`.backlog/`**.
- Files are named using a convention like: `T01.02-my-task-title.md`.
- Project documentation is in **`docs/`**
- You DO NOT need to interact with the file system directly for task management.

### 🔧 **ACTING** (How to change things)

- **All task operations MUST use the MCP server tools.**
- This ensures metadata is correctly updated and the project stays in sync.
- The tools use and return structured data, so you don't need to parse files.


## 2. Common Mistakes to Avoid

### ❌ **WRONG: Direct File Editing or CLI Usage**

```bash
# DON'T DO THIS (Direct Edit):
# 1. Read .backlog/T07-feature.md
# 2. Manually change "- [ ]" to "- [x]" in the content
# 3. Write the modified content back to the file

# DON'T DO THIS (CLI Usage):
backlog edit 7 --check-ac 1
```

### ✅ **CORRECT: Using MCP Tools**

```python
# DO THIS INSTEAD:

# Mark AC #1 as complete
tools.task_edit(id="T7", check_ac=[1])

# Add notes
tools.task_edit(id="T7", notes="Implementation complete")

# Multiple changes: change status and assign the task
tools.task_edit(id="T7", status="in-progress", assigned=["@agent-k"])
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

## Implementation Plan

1. Research approach
2. Implement solution

## Implementation Notes

Summary of what was done.
```

### How to Modify Each Section

| What You Want to Change | MCP Tool Call                                              |
| ----------------------- | ---------------------------------------------------------- |
| Title                   | `task_edit(id="T42", title="New Title")`                   |
| Status                  | `task_edit(id="T42", status="in-progress")`                |
| Assigned                | `task_edit(id="T42", assigned=["@sara"])`                  |
| Labels                  | `task_edit(id="T42", labels=["backend", "api"])`           |
| Description             | `task_edit(id="T42", description="New description")`       |
| Add AC                  | `task_edit(id="T42", add_ac=["New criterion"])`             |
| Check AC #1             | `task_edit(id="T42", check_ac=[1])`                         |
| Uncheck AC #2           | `task_edit(id="T42", uncheck_ac=[2])`                       |
| Remove AC #1            | `task_edit(id="T42", remove_ac=[1])`                        |
| Add Plan                | `task_edit(id="T42", plan="1. Step one
2. Step two")`      |
| Add Notes               | `task_edit(id="T42", notes="What I did")`                  |
| Remove Assigned User    | `task_edit(id="T42", unassign=["@sara"])`                   |
| Remove Labels           | `task_edit(id="T42", remove_labels=["backend", "api"])`     |

---

## 4. Defining Tasks

### Creating New Tasks

**Always use the `task_create` tool:**

```python
# Example
tools.task_create(
  title="Task title",
  description="Description of the task.",
  assigned=["agent-mcp"],
  labels=["feature", "documentation"],
  priority="medium",
  plan="1. Step1
2. Step 2
3. Step 3
",
  notes="what I did",
  ac=["First criterion", "Second criterion"]
)
```

Replace "agent-mcp" with your name.

### Acceptance Criteria (The "what")

**Managing Acceptance Criteria via MCP tools:**

- **Adding criteria** uses the `add_ac` parameter with a list of strings.
- **Checking/unchecking/removing** use `check_ac`, `uncheck_ac`, `remove_ac` parameters with lists of 1-based indices.
- You can perform multiple operations in a single tool call.

```python
# Examples

# Add new criteria
tools.task_edit(id="T42", add_ac=["User can login", "Session persists"])

# Check multiple criteria by index
tools.task_edit(id="T42", check_ac=[1, 2, 3])

# Uncheck a criterion
tools.task_edit(id="T42", uncheck_ac=[2])

# Remove multiple criteria
tools.task_edit(id="T42", remove_ac=[2, 4])
# Note: Indices are processed high-to-low

# Mixed operations in a single command
tools.task_edit(
  id="T42",
  check_ac=[1],
  uncheck_ac=[2],
  remove_ac=[3],
  add_ac=["New criterion"]
)
```

### Task Breakdown Strategy
1. Identify foundational components first
2. Create tasks in dependency order (foundations before features)
3. Ensure each task delivers value independently
4. Avoid creating tasks that block each other

### Task Requirements
- Tasks must be **atomic** and **testable** or **verifiable**
- Each task should represent a single unit of work for one PR
- **Never** reference future tasks (only tasks with id < current task id)
- Ensure tasks are **independent** and don't depend on future work

---

## 5. Implementing Tasks

### 5.1. First step when implementing a task

The very first things you must do when you take over a task are to set the task to "In Progress" and assign it to yourself.

```python
# Example
tools.task_edit(id="T42", status="in-progress", assigned=["@{myself}"])
```

### 5.2. Create an Implementation Plan (The "how")

Once you are familiar with the task, create a plan on **HOW** to tackle it. Write it down in the task so that you can refer to it later.

```python
# Example
tools.task_edit(id="T42", plan="1. Research codebase for references
2. Research on internet for similar cases
3. Implement
4. Test")
```

### 5.3. Implementation Notes (PR description)

When you are done implementing a task, write a clean description in the task notes, as if it were a PR description. Append notes progressively during implementation using `append_notes`.

```python
# Example
tools.task_edit(id="T42", notes="Implemented using pattern X because of Reason Y. Modified files Z and W.")
```

**IMPORTANT**: Do NOT include an Implementation Plan when creating a task. The plan is added only after you start the implementation.
- Creation phase: provide Title, Description, Acceptance Criteria, and optionally labels/priority/assigned.
- When you begin work, switch to edit, set the task in progress and assign to yourself `tools.task_edit(id="<id>", status="in-progress", assigned=["..."])`.
- Think about how you would solve the task and add the plan: `tools.task_edit(id="<id>", plan="...")`.
- Add Implementation Notes only after completing the work: `tools.task_edit(id="<id>", notes="...")` (replace) or append progressively using `append_notes`.

---

## 6. Typical Workflow

```python
# 1. Identify work
tools.task_list(status=["todo"])
tools.task_list(status=["todo", "in-progress"])  # Multiple statuses
tools.task_list(unassigned=True)  # Find tasks needing assignment
tools.task_list(assigned=["alice"])  # Tasks assigned to specific person
tools.task_list(assigned=["alice", "bob"])  # Tasks assigned to alice OR bob
tools.task_list(has_dependency=True)  # Tasks waiting on dependencies
tools.task_list(depended_on=True, status=["todo"])  # Blocking tasks
tools.task_list(labels=["bug", "critical"])  # Tasks with specific labels
tools.task_list(status=["todo"], sort="priority", reverse=True)  # High priority first

# Pagination examples
tools.task_list(limit=5)  # Get first 5 tasks
tools.task_list(status=["todo"], limit=10)  # First 10 todo tasks
tools.task_list(query="feature", limit=3)  # First 3 feature matches

# 2. Read task details
tools.task_view(id="T42")

# 3. Start work: assign yourself & change status
tools.task_edit(id="T42", status="in-progress", assigned=["@myself"])

# 4. Add implementation plan
tools.task_edit(id="T42", plan="1. Analyze
2. Refactor
3. Test")

# 5. Work on the task (write code, test, etc.)

# 6. Mark acceptance criteria as complete
tools.task_edit(id="T42", check_ac=[1, 2, 3])  # Check all at once

# 7. Add implementation notes (PR Description)
tools.task_edit(id="T42", notes="Refactored using strategy pattern, updated tests.")

# 8. Mark task as done
tools.task_edit(id="T42", status="done")
```

---

## 7. Definition of Done (DoD)

A task is **Done** only when **ALL** of the following are complete:

### ✅ Via MCP Tools:

1.  **All acceptance criteria checked**: Use `task_edit(id="ID", check_ac=[N])` for each criterion.
2.  **Implementation notes added**: Use `task_edit(id="ID", notes="...")`.
3.  **Status set to Done**: Use `task_edit(id="ID", status="done")`.

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
| View task  | Use `task_view(id="T42")`                      | Open and read .md file directly  |
| List tasks | Use `task_list(status=["todo"])`               | Browse the `.backlog` folder     |
| List Tasks | Use `task_list(unassigned=True)`               | Browse the `.backlog` folder     |
| List Tasks | Use `task_list(assigned=["alice"])`            | Browse the `.backlog` folder     |

### Modifying Tasks

| Task          | ✅ DO                                          | ❌ DON'T                           |
| ------------- | ---------------------------------------------- | ---------------------------------- |
| Check AC      | Use `task_edit(id="T42", check_ac=[1])`        | Change `- [ ]` to `- [x]` in file  |
| Add notes     | Use `task_edit(id="T42", notes="...")`         | Type notes into .md file           |
| Change status | Use `task_edit(id="T42", status="done")`       | Edit status in frontmatter         |
| Add AC        | Use `task_edit(id="T42", add_ac=["New"])`       | Add `- [ ] New` to file            |
| Archive task  | Use `task_archive(id="T42")`                   | Manually move files to archive folder |

---

## 9. Complete MCP Tool Reference

### `task_create`

Creates a new task.

| Parameter     | Type           | Description                               |
| ------------- | -------------- | ----------------------------------------- |
| `title`       | `string`       | **Required.** The title of the task.      |
| `description` | `string`       | A detailed description of the task.       |
| `parent`      | `string`       | The ID of the parent task.                |
| `ac`          | `list[string]` | A list of acceptance criteria.            |
| `assigned`    | `list[string]` | A list of assigned users.                 |
| `labels`      | `list[string]` | A list of labels.                         |
| `priority`    | `string`       | The priority of the task.                 |
| `deps`        | `list[string]` | A list of task dependencies.              |

### `task_edit`

Edits an existing task.

| Parameter       | Type           | Description                                       |
| --------------- | -------------- | ------------------------------------------------- |
| `id`            | `string`       | **Required.** The ID of the task to edit.         |
| `title`         | `string`       | A new title for the task.                         |
| `description`   | `string`       | A new description for the task.                   |
| `status`        | `string`       | A new status (e.g., "in-progress", "done").      |
| `deps`          | `list[string]` | Set dependencies (replaces existing).             |
| `parent`        | `string`       | A new parent task ID.                             |
| `assigned`      | `list[string]` | Set assigned users (replaces existing).           |
| `unassign`      | `list[string]` | A list of users to unassign.                      |
| `labels`        | `list[string]` | Set labels (replaces existing).                   |
| `remove_labels` | `list[string]` | A list of labels to remove.                       |
| `priority`      | `string`       | A new priority.                                   |
| `add_ac`        | `list[string]` | A list of new acceptance criteria to add.         |
| `remove_ac`     | `list[int]`    | A list of 1-based indices of AC to remove.        |
| `check_ac`      | `list[int]`    | A list of 1-based indices of AC to check.         |
| `uncheck_ac`    | `list[int]`    | A list of 1-based indices of AC to uncheck.       |
| `plan`          | `string`       | Set implementation plan (replaces existing).      |
| `notes`         | `string`       | Set implementation notes (replaces existing).     |
| `append_notes`  | `string`       | Append to existing implementation notes.          |

### `task_list`

Lists tasks with optional filtering, sorting, and pagination.

| Parameter      | Type           | Description                                                   |
| -------------- | -------------- | ------------------------------------------------------------- |
| `status`       | `list[string]` | Filter by status.                                             |
| `parent`       | `string`       | Filter by parent task ID.                                     |
| `assigned`     | `list[string]` | Filter by assigned user.                                      |
| `unassigned`   | `bool`         | Filter tasks that have no assigned users.                     |
| `labels`       | `list[string]` | Filter by labels.                                             |
| `has_dependency`| `bool`        | Filter tasks that have dependencies.                          |
| `depended_on`  | `bool`         | Filter tasks that are depended on by other tasks.             |
| `sort`         | `string`       | Sort by field (id, title, status, priority, created, updated).|
| `reverse`      | `bool`         | Reverse the sort order.                                       |
| `limit`        | `int`          | Maximum number of tasks to return (0 means no limit).         |
| `offset`       | `int`          | Number of tasks to skip from the beginning.                   |
| `query`        | `string`       | Search query to filter tasks by.                              |

### `task_view`

Retrieves and displays the details of a single task.

| Parameter | Type     | Description                       |
| --------- | -------- | --------------------------------- |
| `id`      | `string` | **Required.** The ID of the task. |

### `task_archive`

Archives a task.

| Parameter | Type     | Description                       |
| --------- | -------- | --------------------------------- |
| `id`      | `string` | **Required.** The ID of the task. |

---

## 10. Pagination: Handling Large Task Lists

### When to Use Pagination

Use pagination when:
- Working with projects that have many tasks (>25)
- Performing exploratory queries where you want to see a sample first
- Building interfaces that need to display results in pages
- Avoiding overwhelming output in conversations

### Pagination Parameters

`task_list` supports pagination:

- **`limit`**: Maximum number of results to return (0 = no limit)
- **`offset`**: Number of results to skip from the beginning

### Pagination Examples

```python
# Get first 10 tasks
tools.task_list(limit=10)

# Get next 10 tasks (pagination)
tools.task_list(limit=10, offset=10)

# Get first 5 high-priority tasks
tools.task_list(status=["todo"], sort="priority", reverse=True, limit=5)

# Search second page
tools.task_list(query="api", limit=3, offset=3)
```

### Pagination Response Format

When pagination is used, the tool output includes metadata showing the total number of tasks and pages.

### Best Practices

1. **Start with small limits**: Use `limit=10` to get an overview.
2. **Check output metadata**: Look at the pagination info to determine if more results exist.
3. **Progressive exploration**: Increase `offset` to see more results.
4. **Combine with filtering**: Use pagination with `status` or `priority` filters for focused results.

---

## Remember: The Golden Rule

**🎯 If you want to change ANYTHING in a task, use the `task_edit` tool.**
**📖 Use `task_view` and `task_list` to read tasks. Never write to files directly or use the CLI.**

---

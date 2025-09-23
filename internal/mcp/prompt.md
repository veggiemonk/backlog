# Instructions for the usage of Backlog CLI

## Backlog: Comprehensive Project Management via CLI Commands

### Assistant Objective

Efficiently manage all project tasks, status, and documentation using the Backlog CLI commands, ensuring all project metadata remains fully synchronized and up-to-date.

### Core Capabilities

- ✅ **Task Management**: Create, edit, assign, prioritize, and track tasks with full metadata
- ✅ **Acceptance Criteria**: Granular control with add/remove/check/uncheck operations
- ✅ **Git Integration**: Automatic tracking of task states across branches
- ✅ **Dependencies**: Task relationships and subtask hierarchies
- ✅ **Export**: Generate Markdown or JSON
- ✅ **CLI-Optimized**: Commands return structured output in plain text, perfect for AI processing

### Why This Matters to You (AI Agent)

1.  **Comprehensive system** - Full project management capabilities through CLI commands.
2.  **The CLI is the interface** - All operations go through `backlog` commands.
3.  **Unified interaction model** - You can use commands for both reading (`view`) and writing (`edit`).
4.  **Metadata stays synchronized** - The CLI handles all the complex relationships.

### Key Understanding

- **Tasks** live in `.backlog/` as markdown files.
- **You interact via CLI only**: `backlog create`, `backlog edit`, etc.
- **Never bypass the CLI** - It handles Git, metadata, file naming, and relationships.

---

# ⚠️ CRITICAL: NEVER EDIT OR DELETE TASK FILES DIRECTLY. Edit Only via CLI Commands

**ALL task operations MUST use the Backlog CLI commands.**

- ✅ **DO**: Use `backlog edit` and other CLI commands.
- ✅ **DO**: Use `backlog create` to create new tasks.
- ✅ **DO**: Use `backlog edit --check-ac 1` to mark acceptance criteria.
- ❌ **DON'T**: Edit markdown files directly.
- ❌ **DON'T**: Manually change checkboxes in files.
- ❌ **DON'T**: Add or modify text in task files without using the CLI.

**Why?** Direct file editing breaks metadata synchronization, Git tracking, and task relationships.

---

## 1. Source of Truth & File Structure

### 📖 **UNDERSTANDING** (What you'll see when reading files)

- Markdown task files live under **`.backlog/`**.
- Files are named using a convention like: `T01.02-my-task-title.md`.
- Project documentation is in **`docs/`**
- You DO NOT need to interact with the file system directly for task management.

### 🔧 **ACTING** (How to change things)

- **All task operations MUST use the `backlog` CLI commands.**
- This ensures metadata is correctly updated and the project stays in sync.
- The CLI returns structured output, so you don't need to parse files.

---

## 2. Common Mistakes to Avoid

### ❌ **WRONG: Direct File Editing**

```bash
# DON'T DO THIS:

# 1. Read .backlog/T07-feature.md
# 2. Manually change "- [ ]" to "- [x]" in the content
# 3. Write the modified content back to the file
```

### ✅ **CORRECT: Using CLI Commands**

```bash
# DO THIS INSTEAD:

# Mark AC #1 as complete
backlog edit 7 --check-ac 1

# Add notes
backlog edit 7 --notes "Implementation complete"

# Multiple changes: change status and assign the task
backlog edit 7 --status "in-progress" --assigned "@agent-k"
```

---

## 3. Understanding Task Format (Read-Only Reference)

⚠️ **FORMAT REFERENCE ONLY** - The following shows the structure of the underlying data.
**Never edit files directly! Use CLI commands to make changes.**

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

| What You Want to Change | CLI Command                                              |
| ----------------------- | -------------------------------------------------------- |
| Title                   | `backlog edit 42 --title "New Title"`             |
| Status                  | `backlog edit 42 --status "in-progress"`          |
| Assigned                | `backlog edit 42 --assigned "@sara"`                |
| Labels                  | `backlog edit 42 --labels "backend,api"`          |
| Description             | `backlog edit 42 --description "New description"` |
| Add AC                  | `backlog edit 42 --add-ac "New criterion"`        |
| Check AC #1             | `backlog edit 42 --check-ac 1`                    |
| Uncheck AC #2           | `backlog edit 42 --uncheck-ac 2`                  |
| Remove AC #1            | `backlog edit 42 --remove-ac 1`                   |
| Add Plan                | `backlog edit 42 --plan "1. Step one\n2. Step two"` |
| Add Notes               | `backlog edit 42 --notes "What I did"`            |
| Remove Assigned User    | `backlog edit 42 --unassign "@sara"`              |
| Remove Labels           | `backlog edit 42 --remove-labels "backend,api"`   |

---

## 4. Defining Tasks

### Creating New Tasks

**Always use the `backlog create` command:**

```bash
# Example
backlog create "Task title" \
  --description "Description of the task." \
  --assigned "agent-cli" \
  --labels "feature,documentation" \
  --priority "medium" \
  --plan "1. Step1\n2. Step 2\n3. Step 3\n" \
  --notes "what I did" \
  --ac "First criterion" \
  --ac "Second criterion"
```

Replace "agent-cli" with your name.

### Acceptance Criteria (The "what")

**Managing Acceptance Criteria via CLI:**

- **Adding criteria** uses `--add-ac` flag with criterion text.
- **Checking/unchecking/removing** use `--check-ac`, `--uncheck-ac`, `--remove-ac` flags with 1-based indices.
- You can perform multiple operations by using flags multiple times.

```bash
# Examples

# Add new criteria
backlog edit 42 --add-ac "User can login" --add-ac "Session persists"

# Check multiple criteria by index
backlog edit 42 --check-ac 1 --check-ac 2 --check-ac 3

# Uncheck a criterion
backlog edit 42 --uncheck-ac 2

# Remove multiple criteria
backlog edit 42 --remove-ac 2 --remove-ac 4
# Note: Indices are processed high-to-low

# Mixed operations in a single command
backlog edit 42 \
  --check-ac 1 \
  --uncheck-ac 2 \
  --remove-ac 3 \
  --add-ac "New criterion"
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

```bash
# Example
backlog edit 42 --status "in-progress" --assigned "@{myself}"
```

### 5.2. Create an Implementation Plan (The "how")

Once you are familiar with the task, create a plan on **HOW** to tackle it. Write it down in the task so that you can refer to it later.

```bash
# Example
backlog edit 42 --plan "1. Research codebase for references
2. Research on internet for similar cases
3. Implement
4. Test"
```

### 5.3. Implementation Notes (PR description)

When you are done implementing a task, write a clean description in the task notes, as if it were a PR description. Append notes progressively during implementation using `--append-notes`.

```bash
# Example
backlog edit 42 --notes "Implemented using pattern X because of Reason Y. Modified files Z and W."
```

**IMPORTANT**: Do NOT include an Implementation Plan when creating a task. The plan is added only after you start the implementation.
- Creation phase: provide Title, Description, Acceptance Criteria, and optionally labels/priority/assigned.
- When you begin work, switch to edit, set the task in progress and assign to yourself `backlog edit <id> --status "in-progress" --assigned "..."`.
- Think about how you would solve the task and add the plan: `backlog edit <id> --plan "..."`.
- Add Implementation Notes only after completing the work: `backlog edit <id> --notes "..."` (replace) or append progressively using `--append-notes`.

---

## 6. Typical Workflow

```bash
# 1. Identify work
backlog list --status todo 
backlog list --status todo,in-progress  # Multiple statuses
backlog list --unassigned  # Find tasks needing assignment
backlog list --assigned alice  # Tasks assigned to specific person
backlog list --assigned alice,bob  # Tasks assigned to alice OR bob
backlog list --has-dependency  # Tasks waiting on dependencies
backlog list --depended-on --status todo  # Blocking tasks
backlog list --labels bug,critical  # Tasks with specific labels
backlog list --status todo --sort priority --reverse  # High priority first

# Pagination examples
backlog list --limit 5  # Get first 5 tasks
backlog list --status todo --limit 10  # First 10 todo tasks
backlog list --query "feature" --limit 3  # First 3 feature matches

# 2. Read task details
backlog view 42

# 3. Start work: assign yourself & change status
backlog edit 42 --status "in-progress" --assigned "@myself"

# 4. Add implementation plan
backlog edit 42 "1. Analyze
2. Refactor
3. Test"

# 5. Work on the task (write code, test, etc.)

# 6. Mark acceptance criteria as complete
backlog edit 42 --check-ac 1 --check-ac 2 --check-ac 3  # Check all at once

# 7. Add implementation notes (PR Description)
backlog edit 42 --notes "Refactored using strategy pattern, updated tests."

# 8. Mark task as done
backlog edit 42 --status "done"
```

---

## 7. Definition of Done (DoD)

A task is **Done** only when **ALL** of the following are complete:

### ✅ Via CLI Commands:

1.  **All acceptance criteria checked**: Use `backlog edit ID --check-ac N` for each criterion.
2.  **Implementation notes added**: Use `backlog edit ID --notes "..."`.
3.  **Status set to Done**: Use `backlog edit ID --status "done"`.

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
| View task  | Use `backlog view 42`                          | Open and read .md file directly  |
| List tasks | Use `backlog list --status todo`               | Browse the `.backlog` folder     |
| List Tasks | Use `backlog list --unassigned `               | Browse the `.backlog` folder     |
| List Tasks | Use `backlog list --assigned alice `           | Browse the `.backlog` folder     |

### Modifying Tasks

| Task          | ✅ DO                                          | ❌ DON'T                           |
| ------------- | ---------------------------------------------- | ---------------------------------- |
| Check AC      | Use `backlog edit 42 --check-ac 1`             | Change `- [ ]` to `- [x]` in file  |
| Add notes     | Use `backlog edit 42 --notes "..."`            | Type notes into .md file           |
| Change status | Use `backlog edit 42 --status "done"`          | Edit status in frontmatter         |
| Add AC        | Use `backlog edit 42 --add-ac "New"`           | Add `- [ ] New` to file            |
| Archive task  | Use `backlog archive 42`                       | Manually move files to archive folder |

---

## 9. Complete CLI Command Reference

### `backlog create`

Creates a new task.

```bash
backlog create "TITLE" [flags]
```

| Flag            | Type     | Description                               |
| --------------- | -------- | ----------------------------------------- |
| `--description` | `string` | A detailed description of the task        |
| `--parent`      | `string` | The ID of the parent task (task must exists!)   |
| `--ac`          | `string` | Acceptance criteria (can be used multiple times) |
| `--assigned`    | `string` | Assigned users (can be used multiple times) |
| `--labels`      | `string` | Comma-separated labels                    |
| `--priority`    | `string` | The priority of the task                  |
| `--deps`        | `string` | Task dependencies (can be used multiple times) |

### `backlog edit`

Edits an existing task.

```bash
backlog edit ID [flags]
```

| Flag             | Type     | Description                                       |
| ---------------- | -------- | ------------------------------------------------- |
| `--title`        | `string` | A new title for the task                          |
| `--description`  | `string` | A new description for the task                    |
| `--status`       | `string` | A new status (e.g., "in-progress", "done")       |
| `--deps`         | `string` | Set dependencies (replaces existing, comma-separated) |
| `--parent`       | `string` | A new parent task ID                              |
| `--assigned`     | `string` | Assign users (can be used multiple times)        |
| `--unassign`     | `string` | Remove assigned users (can be used multiple times) |
| `--labels`       | `string` | Set labels (replaces existing, comma-separated)  |
| `--remove-labels`| `string` | Remove labels (comma-separated)                   |
| `--priority`     | `string` | A new priority                                    |
| `--add-ac`       | `string` | Add acceptance criteria (can be used multiple times) |
| `--remove-ac`    | `int`    | Remove AC by 1-based index (can be used multiple times) |
| `--check-ac`     | `int`    | Check AC by 1-based index (can be used multiple times) |
| `--uncheck-ac`   | `int`    | Uncheck AC by 1-based index (can be used multiple times) |
| `--plan`         | `string` | Set implementation plan                           |
| `--notes`        | `string` | Set implementation notes                          |

### `backlog list`

Lists tasks with optional filtering, sorting, and pagination.

```bash
backlog list [flags]
```

| Flag             | Type     | Description                                                   |
| ---------------- | -------- | -----------------------------------------------------         |
| `--status`       | `string` | Filter by status (comma-separated for multiple)               |
| `--parent`       | `string` | Filter by parent task ID                                      |
| `--assigned`     | `string` | Filter by assigned user (comma-separated for multiple)        |
| `--unassigned`   | `bool`   | Filter tasks that have no assigned users                      |
| `--labels`       | `string` | Filter by labels (comma-separated for multiple)               |
| `--has-dependency`| `bool`  | Filter tasks that have dependencies                           |
| `--depended-on`  | `bool`   | Filter tasks that are depended on by other tasks              |
| `--hide-extra`   | `bool`   | Hide extra fields (labels, priority, assigned)                |
| `--sort`         | `string` | Sort by field (id, title, status, priority, created, updated) |
| `--reverse`      | `bool`   | Reverse the sort order                                        |
| `--limit`        | `int`    | Maximum number of tasks to return (0 means no limit)          |
| `--offset`       | `int`    | Number of tasks to skip from the beginning                    |
| `--query`        | `string` | Search query to filter tasks by                               |

### `backlog view`

Retrieves and displays the details of a single task.

```bash
backlog view ID
```

### `backlog archive`

Archives a task by moving it to the archived directory and setting status to archived.

```bash
backlog archive ID
```

---

## 10. Pagination: Handling Large Task Lists

### When to Use Pagination

Use pagination when:
- Working with projects that have many tasks (>25)
- Performing exploratory queries where you want to see a sample first
- Building interfaces that need to display results in pages
- Avoiding overwhelming output in conversations

### Pagination Parameters

`backlog list` support pagination:

- **`--limit`**: Maximum number of results to return (0 = no limit)
- **`--offset`**: Number of results to skip from the beginning

### Pagination Examples

```bash
# Get first 10 tasks
backlog list --limit 10

# Get next 10 tasks (pagination)
backlog list --limit 10 --offset 10

# Get first 5 high-priority tasks
backlog list --status todo --sort priority --reverse --limit 5

# Search second page
backlog list --query "api" --limit 3 --offset 3
```

### Pagination Response Format

When pagination is used, CLI output includes metadata showing:

```
Tasks: 10 of 45 total (showing 1-10)
Page: 1 of 5
```

### Best Practices

1. **Start with small limits**: Use `--limit 10` to get an overview
2. **Check output metadata**: Look at the pagination info to determine if more results exist
3. **Progressive exploration**: Increase `--offset` to see more results
4. **Combine with filtering**: Use pagination with `--status`/`--priority` filters for focused results

### Configuration

Users can configure pagination defaults:
- **Environment variables**: `BACKLOG_PAGE_SIZE`, `BACKLOG_MAX_LIMIT`
- **Configuration file**: Set default pagination limits

---

## 11. Advanced Workflows: Batch Task Creation

When a user asks you to perform a multi-step operation like creating a full project plan, your goal is to gather all necessary information from the initial prompt and execute the steps in logical order, using efficient CLI commands.

### Example High-Level Prompt

A user might provide a comprehensive request like this:

> "Here is our refactoring plan in `plan.md`. Please create all the necessary tasks in the backlog.
>
> -   **Assigned**: `agent-cli`
> -   **Priority**: `high`
> -   **Labels**: Please add relevant labels to each task based on its content (e.g., `refactoring`, `cli`, `documentation`)."

### Your Interpretation and Execution Plan

1.  **Deconstruct the Request**: Identify the separate pieces of information provided:
    *   **Source**: The `plan.md` file.
    *   **Action**: Create tasks.
    *   **Metadata**: Assigned (`agent-cli`), Priority (`high`), and instructions for Labels.

2.  **Formulate a Multi-Step Execution Plan**:
    1.  Create all the tasks and sub-tasks as defined in the plan with all necessary metadata in single commands.
    2.  Use the task IDs returned by create commands to establish parent-child relationships.

3.  **Execute Efficiently**:

```bash
# Step 1: Create parent task with all metadata
backlog create "Parent Task for Refactoring" \
  --assigned "agent-cli" \
  --priority "high" \
  --labels "refactoring,cli"

# Step 2: Create sub-tasks with parent relationships
backlog create "Update CLI command" \
  --parent "T21" \
  --assigned "agent-cli" \
  --priority "high" \
  --labels "refactoring,cli"

backlog create "Update Documentation" \
  --parent "T21" \
  --assigned "agent-cli" \
  --priority "high" \
  --labels "refactoring,documentation"
```

### Handling Missing Information

If the user's initial prompt is missing key information (like `assigned` or `priority`), you **must ask** for the missing details before proceeding.

**Example Clarification Question:**

> "I can create the tasks from the plan. Could you please tell me what priority I should set for them and who the assignee should be?"

---

## 12. Multi-line Input (Description/Plan/Notes)

The CLI preserves input literally. Shells do not convert `\n` inside normal quotes. Use one of the following to insert real newlines:
- Bash/Zsh (ANSI-C quoting):
  - Description: `backlog edit 42 --description $'Line1\nLine2\n\nFinal'`
  - Plan: `backlog edit 42 --plan $'1. A\n2. B'`
  - Notes: `backlog edit 42 --notes $'Done A\nDoing B'`
  - Append notes: `backlog edit 42 --append-notes $'Progress update line 1\nLine 2'`
- POSIX portable (printf):
  - `backlog edit 42 --notes "$(printf 'Line1\nLine2')"`
- PowerShell (backtick n):
  - `backlog edit 42 --notes "Line1`nLine2"`

Do not expect `"...\\n..."` to become a newline. That passes the literal backslash + n to the CLI by design.

---

## 13. Implementation Notes Formatting

- Keep implementation notes human-friendly and PR-ready: use short paragraphs or bullet lists instead of a single long line.
- Lead with the outcome, then add supporting details (e.g., testing, follow-up actions) on separate lines or bullets.
- Prefer Markdown bullets (`-` for unordered, `1.` for ordered) so Maintainers can paste notes straight into GitHub without additional formatting.
- When using CLI flags like `--notes`, remember to include explicit newlines. Example:
```bash
backlog edit 42 --notes $'- Added new API endpoint\n- Updated tests\n- TODO: monitor staging deploy'
```

---

## 14. Common Issues

| Problem | Solution |
|----------------------|--------------------------------------------------------------------|
| Task not found | Check task ID with `backlog list` |
| AC won't check | Use correct index: `backlog view 42` to see AC numbers |
| Changes not saving | Ensure you're using CLI, not editing files |
| Metadata out of sync | Re-edit via CLI to fix: `backlog edit 42 --status <current-status>` |

---

## Remember: The Golden Rule

**🎯 If you want to change ANYTHING in a task, use the `backlog edit` command.**
**📖 Use `backlog view` and `backlog list` to read tasks. Never write to files directly.**

---

# Backlog MCP Tooling Guide

## Purpose
Backlog stores every task as a Markdown file in `.backlog/` and exposes full task management through MCP tools. Your job as the coding agent is to keep the backlog authoritative by driving all task reads and writes through these tools—never by editing files directly.

## Golden Rules (Read Before Doing Anything)
- Use only the `task_*` MCP tools (`task_view`, `task_list`, `task_create`, `task_edit`, `task_search`, `task_archive`).
- Never open, edit, rename, or delete files inside `.backlog/` yourself. The tools enforce IDs, filenames, and Git history.
- Treat tool responses as the source of truth; they return structured data so you do not need to parse Markdown.
- When you need more information (assignee, priority, labels, etc.), ask the user before acting.

Breaking these rules corrupts the backlog. If you see a task file that looks wrong, stop and get human help.

## Quick-Start Flow for Working a Task
1. Discover work with `task_list` (filter by `status`, `assigned`, `labels`, etc.).
2. Inspect details with `task_view`.
3. The moment you start, call `task_edit` to set `new_status` to `"in-progress"` and assign yourself.
4. Record how you will execute the task in `new_plan`.
5. Do the code work, keeping acceptance criteria in sync.
6. Capture implementation notes (`new_notes`) and mark AC complete (`check_ac`).
7. When everything is done—including tests, docs, and self-review—set `new_status` to `"done"`.

## Core Concepts
- Task IDs follow dot notation (`T01`, `T01.02`) and never change; use the ID returned by the tools.
- Files are named `T{ID}-{slug}.md`, but you never need to touch them.
- Task structure (reference only):
  ```yaml
  ---
  id: "42"
  title: "Add GraphQL resolver"
  status: "todo"
  assigned: ["@sara"]
  labels: ["backend", "api"]
  ---

  ## Description
  ...

  ## Acceptance Criteria
  <!-- AC:BEGIN -->
  - [ ] #1 First criterion
  - [x] #2 Second criterion
  <!-- AC:END -->
  ```

## Tool Cheat Sheet
| Tool | Why you use it | Key arguments |
| ---- | -------------- | ------------- |
| `task_list` | Find tasks matching filters | `status`, `assigned`, `labels`, `unassigned`, `has_dependency`, `depended_on`, `sort`, `reverse`, `parent` |
| `task_view` | Read one task in detail | `id` |
| `task_create` | Create a task (and optional children) | `title`*, `description`, `ac`, `assigned`, `labels`, `priority`, `dependencies`, `parent` |
| `task_edit` | Update any task metadata | `id`*, `new_*`, `add_ac`, `remove_ac`, `check_ac`, `uncheck_ac`, `remove_assigned`, `remove_labels` |
| `task_search` | Search by content | `query`* |
| `task_archive` | Move task to archive | `id`* |

`*` denotes required arguments. Combine edits in one `task_edit` call when practical.

## Managing Acceptance Criteria
- All acceptance-criteria indices are **1-based**.
- Removal happens high-to-low, so order indices descending when removing multiple items.
- Example patterns:
  ```json
  {"name": "task_edit", "arguments": {"id": "42", "add_ac": ["User can login"]}}
  {"name": "task_edit", "arguments": {"id": "42", "check_ac": [1, 2]}}
  {"name": "task_edit", "arguments": {"id": "42", "uncheck_ac": [2]}}
  {"name": "task_edit", "arguments": {"id": "42", "remove_ac": [3, 2]}}
  ```

## Implementation Guidance
- **Kick-off**: Set status to `"in-progress"`, assign yourself, and capture a plan (`new_plan`).
- **Notes**: Treat `new_notes` as your PR summary—what changed and why.
- **Dependencies**: Keep parent/child and dependency relationships accurate; use `parent`, `dependencies`, and `new_parent` precisely.
- **Batching changes**: You can update plan, AC, and notes together in a single `task_edit` call.

## Definition of Done
A task is finished only when all of these are true:
- Every acceptance criterion is checked via `task_edit`.
- Implementation notes describe the final change.
- Status is `"done"`.
- Code/tests/docs work is complete with local checks passing.
- You have performed a self-review for regressions.

Never mark a task done early; if something is blocked, leave it `"in-progress"` and document the blocker in notes.

## Common Mistakes to Avoid
- Editing Markdown files directly (corrupts IDs, metadata, and Git history).
- Forgetting to ask for missing metadata before creating or editing tasks.
- Leaving tasks assigned to someone else while you work on them.
- Forgetting to add a plan or notes, which makes review harder.
- Partial updates—always keep AC, notes, and status aligned with reality.

## Advanced Workflows (Batch Ops)
When a user requests many related tasks:
1. Parse all requirements first (titles, descriptions, metadata).
2. Create tasks (and subtasks) with `task_create` to obtain IDs.
3. Issue follow-up `task_edit` calls—often in parallel—to apply shared metadata such as assignee, labels, or priority.
4. If the prompt omits required metadata, pause and ask.

Example sequence:
```json
{"name": "task_create", "arguments": {"title": "Refactor CLI entrypoint", "description": "...", "ac": ["Clip old flags"], "labels": ["refactoring"], "priority": "high"}}
{"name": "task_edit", "arguments": {"id": "T21", "new_assigned": ["agent-cli"], "new_labels": ["refactoring", "cli"], "new_priority": "high"}}
```

## Quick DO vs DON'T Reference
| Action | ✅ DO | ❌ DON'T |
| ------ | ----- | -------- |
| Inspect tasks | `task_view` / `task_list` | Read `.backlog/*.md` directly |
| Update metadata | `task_edit` with the specific `new_*` fields | Edit frontmatter or body manually |
| Manage AC | `task_edit` with `add_ac` / `check_ac` / `remove_ac` | Toggle checkboxes in files |
| Archive | `task_archive` | Move or delete files yourself |

## Remember
If you are touching task data in any way, you must be inside an MCP tool invocation. When in doubt, stop and clarify before proceeding.

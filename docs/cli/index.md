---
layout: page
title: CLI Reference
nav_order: 3
has_children: true
---

# CLI Reference

Complete command-line interface documentation for the backlog tool.

## Overview

Backlog provides a comprehensive CLI for managing tasks with the following commands:

- [`backlog`](backlog.html) - Main command with global options
- [`backlog create`](backlog_create.html) - Create new tasks
- [`backlog list`](backlog_list.html) - List and filter tasks
- [`backlog view`](backlog_view.html) - View task details
- [`backlog edit`](backlog_edit.html) - Edit existing tasks
- [`backlog search`](backlog_search.html) - Search tasks by content
- [`backlog archive`](backlog_archive.html) - Archive completed tasks
- [`backlog mcp`](backlog_mcp.html) - Start MCP server for AI agents
- [`backlog version`](backlog_version.html) - Show version information

## Quick Examples

### Create Tasks

```bash
# Simple task
backlog create "Fix login bug"

# Task with metadata
backlog create "Implement OAuth" \
  -d "Add Google and GitHub OAuth support" \
  -a "alice,bob" \
  -l "auth,feature" \
  --priority "high" \
  --ac "Google OAuth works" \
  --ac "GitHub OAuth works"

# Subtask
backlog create "OAuth token validation" -p "T01"
```

### List and Filter

```bash
# List all tasks
backlog list

# Filter by status
backlog list --status "in-progress"

# Filter by assignee
backlog list --assignee "alice"

# Show subtasks
backlog list --parent "T01"
```

### View and Edit

```bash
# View task details
backlog view T01

# Update status
backlog edit T01 --status "done"

# Add assignee and labels
backlog edit T01 --add-assignee "bob" --add-label "urgent"
```

### Search and Archive

```bash
# Search tasks
backlog search "authentication"

# Archive completed task
backlog archive T01
```

### AI Integration

```bash
# Start MCP server (STDIO)
backlog mcp

# Start MCP server (HTTP)
backlog mcp --http --port 8106
```

## Global Options

All commands support these global options:

- `--help, -h`: Show help information
- `--version`: Show version information

## Output Formats

Many commands support multiple output formats:

- **Table** (default): Human-readable table format
- **JSON**: Machine-readable JSON format
- **Markdown**: Markdown format for documentation

Use the `--format` flag to specify the output format:

```bash
backlog list --format json
backlog view T01 --format markdown
```

## Task ID Format

Tasks use hierarchical IDs with dot notation:

- `T01` - Root task
- `T01.01` - Subtask of T01
- `T01.01.01` - Sub-subtask of T01.01

IDs are automatically assigned when creating tasks, or you can specify a parent with the `-p` flag.

## File Structure

Tasks are stored as Markdown files in the `.backlog/` directory:

```
.backlog/
├── T01-implement_auth.md
├── T01.01-oauth_setup.md
├── T02-frontend_work.md
└── archived/
    └── T03-completed_task.md
```

Each command in this reference includes detailed usage information, examples, and available options.
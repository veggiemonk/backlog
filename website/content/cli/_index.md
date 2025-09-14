---
title: "CLI Reference"
description: "Complete command-line interface documentation"
weight: 30
geekdocCollapseSection: true
---

# CLI Reference

Complete command-line interface documentation for the backlog tool.

## Overview

Backlog provides a comprehensive CLI for managing tasks with the following commands:

{{< columns >}}

### Core Commands
- [`backlog`](backlog) - Main command with global options
- [`backlog create`](backlog_create) - Create new tasks
- [`backlog list`](backlog_list) - List and filter tasks
- [`backlog view`](backlog_view) - View task details
- [`backlog edit`](backlog_edit) - Edit existing tasks

<--->

### Additional Commands
- [`backlog search`](backlog_search) - Search tasks by content
- [`backlog archive`](backlog_archive) - Archive completed tasks
- [`backlog mcp`](backlog_mcp) - Start MCP server for AI agents
- [`backlog version`](backlog_version) - Show version information

{{< /columns >}}

## Quick Examples

### Create Tasks

{{< tabs "create-examples" >}}
{{< tab "Simple Task" >}}
```bash
backlog create "Fix login bug"
```
{{< /tab >}}
{{< tab "With Metadata" >}}
```bash
backlog create "Implement OAuth" \
  -d "Add Google and GitHub OAuth support" \
  -a "alice,bob" \
  -l "auth,feature" \
  --priority "high" \
  --ac "Google OAuth works" \
  --ac "GitHub OAuth works"
```
{{< /tab >}}
{{< tab "Subtask" >}}
```bash
backlog create "OAuth token validation" -p "T01"
```
{{< /tab >}}
{{< /tabs >}}

### List and Filter

{{< tabs "list-examples" >}}
{{< tab "All Tasks" >}}
```bash
backlog list
```
{{< /tab >}}
{{< tab "By Status" >}}
```bash
backlog list --status "in-progress"
```
{{< /tab >}}
{{< tab "By Assignee" >}}
```bash
backlog list --assignee "alice"
```
{{< /tab >}}
{{< tab "Subtasks" >}}
```bash
backlog list --parent "T01"
```
{{< /tab >}}
{{< /tabs >}}

### View and Edit

{{< columns >}}

#### View Details
```bash
backlog view T01
```

#### Update Status
```bash
backlog edit T01 --status "done"
```

<--->

#### Add Assignee
```bash
backlog edit T01 --add-assignee "bob"
```

#### Add Labels
```bash
backlog edit T01 --add-label "urgent"
```

{{< /columns >}}

### Search and Archive

{{< columns >}}

#### Search Content
```bash
backlog search "authentication"
```

<--->

#### Archive Task
```bash
backlog archive T01
```

{{< /columns >}}

### AI Integration

{{< tabs "mcp-examples" >}}
{{< tab "STDIO Transport" >}}
```bash
backlog mcp
```
{{< /tab >}}
{{< tab "HTTP Transport" >}}
```bash
backlog mcp --http --port 8106
```
{{< /tab >}}
{{< /tabs >}}

## Global Options

{{< hint type=tip >}}
All commands support these global options:
{{< /hint >}}

- `--help, -h`: Show help information
- `--version`: Show version information

## Output Formats

Many commands support multiple output formats:

{{< columns >}}

- **Table** (default): Human-readable table format
- **JSON**: Machine-readable JSON format
- **Markdown**: Markdown format for documentation

<--->

```bash
backlog list --format json
backlog view T01 --format markdown
```

{{< /columns >}}

## Task ID Format

{{< hint type=note >}}
Tasks use hierarchical IDs with dot notation:
{{< /hint >}}

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

{{< hint type=tip >}}
Each command in this reference includes detailed usage information, examples, and available options.
{{< /hint >}}
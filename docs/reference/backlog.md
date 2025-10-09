---
layout: page
title: backlog
---

# NAME

backlog - Backlog is a git-native, markdown-based task manager

# SYNOPSIS

backlog

```
[--auto-commit]
[--folder]=[value]
[--log-file]=[value]
[--log-format]=[value]
[--log-level]=[value]
```

# DESCRIPTION

A Git-native, Markdown-based task manager for developers and AI agents.
Backlog helps you manage tasks within your git repository.

**Usage**:

```
backlog [GLOBAL OPTIONS] [command [COMMAND OPTIONS]] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--auto-commit**: Auto-committing changes to git repository

**--folder**="": Directory for backlog tasks (default: .backlog)

**--log-file**="": Log file path (defaults to stderr)

**--log-format**="": Log format (json, text) (default: text)

**--log-level**="": Log level (debug, info, warn, error) (default: info)


# COMMANDS

## archive

Archive a task

>backlog archive <task-id>

## create

Create a new task


**--ac**="": Acceptance criterion (can be specified multiple times) (default: [])

**--assigned, -a**="": Assignee for the task (can be specified multiple times) (default: [])

**--deps**="": Add a dependency (can be used multiple times) (default: [])

**--description, -d**="": Description of the task

**--labels, -l**="": Comma-separated labels for the task (default: [])

**--notes**="": Additional notes for the task

**--parent, -p**="": Parent task ID

**--plan**="": Implementation plan for the task

**--priority**="": Priority of the task (low, medium, high, critical) (default: medium)

## doctor

Diagnose and fix task ID conflicts

**--dry-run**: Show what would be changed without making changes (use with --fix)

**--fix**: Automatically fix detected conflicts

**--json, -j**: Output in JSON format

**--strategy**="": Resolution strategy when using --fix (chronological|auto|manual) (default: chronological)

## edit

Edit an existing task

>backlog edit <id>

**--ac**="": Add a new acceptance criterion (can be used multiple times) (default: [])

**--assigned, -a**="": Add assigned names for the task (can be used multiple times) (default: [])

**--check-ac**="": Check an acceptance criterion by its index (default: [])

**--deps**="": Set dependencies (can be used multiple times) (default: [])

**--description, -d**="": New description for the task

**--labels, -l**="": Add labels for the task (can be used multiple times) (default: [])

**--notes**="": New implementation notes for the task

**--parent, -p**="": New parent for the task

**--plan**="": New implementation plan for the task

**--priority**="": New priority for the task

**--remove-ac**="": Remove an acceptance criterion by its index (default: [])

**--remove-assigned, -A**="": Assigned names to remove from the task (can be used multiple times) (default: [])

**--remove-labels, -L**="": Labels to remove from the task (can be used multiple times) (default: [])

**--status, -s**="": New status for the task

**--title, -t**="": New title for the task

**--uncheck-ac**="": Uncheck an acceptance criterion by its index (default: [])

## instructions

instructions for agents to learn to use backlog

**--mode**="": which mode the agent will use backlog: (cli|mcp) (default: cli)

## list

List all tasks

**--assigned, -a**="": Filter tasks by assigned names (default: [])

**--depended-on, -d**: Filter tasks that are depended on by other tasks

**--has-dependency, -c**: Filter tasks that have dependencies

**--hide-extra, -e**: Hide extra fields (labels, priority, assigned)

**--json, -j**: Print JSON output

**--labels, -l**="": Filter tasks by labels (default: [])

**--limit**="": Maximum number of tasks to return (0 means no limit) (default: 0)

**--markdown, -m**: Print markdown table

**--offset**="": Number of tasks to skip from the beginning (default: 0)

**--parent, -p**="": Filter tasks by parent ID

**--priority**="": Filter tasks by priority

**--query, -q**="": Search query to filter tasks by

**--reverse, -r**: Reverse the order of tasks

**--sort**="": Sort tasks by comma-separated fields (id, title, status, priority, created, updated)

**--status, -s**="": Filter tasks by status (default: [])

**--unassigned, -u**: Filter tasks that have no one assigned

## mcp

Start the MCP server

**--http**: Use HTTP transport instead of stdio

**--port**="": Port for the MCP server (HTTP transport) (default: 8106)

## version

Print the version information

## view

View a task by providing its ID

>backlog view <id>

**--json, -j**: Print JSON output

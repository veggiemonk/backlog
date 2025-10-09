---
layout: page
title: backlog list
---

# NAME

list - List all tasks

# SYNOPSIS

list

```
[--assigned|-a]=[value]
[--depended-on|-d]
[--has-dependency|-c]
[--hide-extra|-e]
[--json|-j]
[--labels|-l]=[value]
[--limit]=[value]
[--markdown|-m]
[--offset]=[value]
[--parent|-p]=[value]
[--priority]=[value]
[--query|-q]=[value]
[--reverse|-r]
[--sort]=[value]
[--status|-s]=[value]
[--unassigned|-u]
```

# DESCRIPTION


Lists all tasks in the backlog except archived tasks.

Examples:
```

  # List all tasks
  backlog list                                    # All tasks with all columns
  backlog list --status "todo"                    # Tasks with status "todo"
  backlog list --status "todo,in-progress"        # Tasks with status "todo" OR "in-progress"
  backlog list --parent "12345"                   # Sub-tasks of task "12345"
  backlog list --assigned "alice"                 # Tasks assigned to alice
  backlog list --unassigned                       # Tasks with no assignee
  backlog list --labels "bug"                     # Tasks with "bug" label
  backlog list --priority "high"                  # High priority tasks

  # Search
  backlog list --query "refactor"                 # Search for "refactor"

  # Dependency filters
  backlog list --has-dependency                   # Tasks with dependencies
  backlog list --depended-on                      # Tasks depended on by others
  backlog list --depended-on --status "todo"      # Blocking tasks

  # Column visibility
  backlog list --hide-extra                       # Minimal columns
  backlog list -e                                 # Minimal columns (short flag)

  # Sorting
  backlog list --sort "priority"                  # Sort by priority
  backlog list --sort "updated,priority"          # Sort by updated, then priority
  backlog list --reverse                          # Reverse order
  backlog list --sort "priority" --reverse        # Sort by priority, reversed

  # Output format
  backlog list --markdown                         # Markdown table
  backlog list -m                                 # Markdown table (short flag)
  backlog list --json                             # JSON output
  backlog list -j                                 # JSON output (short flag)

  # Pagination
  backlog list --limit 10                         # First 10 tasks
  backlog list --limit 5 --offset 10              # Tasks 11-15
  backlog list --status "todo" --limit 3          # First 3 todo tasks

```

**Usage**:

```
list [GLOBAL OPTIONS] [command [COMMAND OPTIONS]] [ARGUMENTS...]
```

# GLOBAL OPTIONS

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


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

# List all tasks
backlog list                                    # List all tasks with all columns
backlog list --status "todo"                    # List tasks with status "todo"
backlog list --status "todo,in-progress"        # List tasks with status "todo" or "in-progress"
backlog list --status "done"                    # List tasks with status "done"
backlog list --parent "12345"                   # List tasks that are sub-tasks of the task with ID "12345"
backlog list --status "todo" --parent "12345"   # List "todo" sub-tasks of task "12345"
backlog list --assigned "alice"                 # List tasks assigned to alice
backlog list --unassigned                       # List tasks that have no one assigned
backlog list --labels "bug"                     # List tasks containing the label "bug"
backlog list --labels "bug,feature"             # List tasks containing either "bug" or "feature" labels
backlog list --priority "high"                  # List all high priority tasks

# Search
backlog list --query "refactor"                 # Search for tasks with the word "refactor" in them

# dependency filters
backlog list --has-dependency                   # List tasks that have at least one dependency
backlog list --depended-on                      # List tasks that are depended on by other tasks
backlog list --depended-on --status "todo"      # List all the blocking tasks.

# column visibility
backlog list --hide-extra                       # Hide extra fields (labels, priority, assigned)
backlog list -e                                 # Hide extra fields (labels, priority, assigned)
backlog list --status "todo" --hide-extra       # List "todo" tasks with minimal columns

# sorting
backlog list --sort "priority"                  # Sort tasks by priority
backlog list --sort "updated,priority"          # Sort tasks by updated date, then priority
backlog list --sort "status,created"            # Sort tasks by status, then creation date
backlog list --reverse                          # Reverse the order of tasks
backlog list --sort "priority" --reverse        # Sort by priority in reverse order
backlog list --status "todo" \
    --priority "medium"  \
    --sort "priority"    \
    --reverse                               # Combine all options

# output format
backlog list -m                                 # List tasks in markdown format
backlog list -markdown                          # List tasks in markdown format
backlog list --json                             # List tasks in JSON format
backlog list -j                                 # List tasks in JSON format
backlog list --status "todo" --json             # List "todo" tasks in JSON format

# pagination
backlog list --limit 10                         # List first 10 tasks
backlog list --limit 5 --offset 10              # List 5 tasks starting from 11th task
backlog list --status "todo" --limit 3          # List first 3 "todo" tasks
backlog list --sort "priority" --limit 10       # List top 10 tasks by priority


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


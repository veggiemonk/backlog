## backlog list

List all tasks

### Synopsis

Lists all tasks in the backlog except archived tasks.

```
backlog list [flags]
```

### Examples

```
# List All Tasks

# List all tasks with all columns
backlog backlog list

# Filter by Status

# List tasks with status "todo"
backlog backlog list --status "todo"

# Filter by Multiple Statuses

# List tasks with status "todo" or "in-progress"
backlog backlog list --status "todo,in-progress"

# Filter by Parent

# List tasks that are sub-tasks of the task with ID "12345"
backlog backlog list --parent "12345"

# Filter by Assigned User

# List tasks assigned to alice
backlog backlog list --assigned "alice"

# Filter Unassigned Tasks

# List tasks that have no one assigned
backlog backlog list --unassigned

# Filter by Labels

# List tasks containing either "bug" or "feature" labels
backlog backlog list --labels "bug,feature"

# Filter by Priority

# List all high priority tasks
backlog backlog list --priority "high"

# Filter Tasks with Dependencies

# List tasks that have at least one dependency
backlog backlog list --has-dependency

# Filter Blocking Tasks

# List all the blocking tasks.
backlog backlog list --depended-on --status "todo"

# Hide Extra Fields

# Hide extra fields (labels, priority, assigned)
backlog backlog list --hide-extra

# Sort by Priority

# Sort tasks by priority
backlog backlog list --sort "priority"

# Multiple Sort Fields

# Sort tasks by updated date, then priority
backlog backlog list --sort "updated,priority"

# Reverse Order

# Reverse the order of tasks
backlog backlog list --reverse

# Markdown Output

# List tasks in markdown format
backlog backlog list --markdown

# JSON Output

# List tasks in JSON format
backlog backlog list --json

# Pagination - Limit

# List first 10 tasks
backlog backlog list --limit "10"

# Pagination - Limit and Offset

# List 5 tasks starting from 11th task
backlog backlog list --offset "10" --limit "5"
```

### Options

```
  -a, --assigned strings   Filter tasks by assigned names
  -d, --depended-on        Filter tasks that are depended on by other tasks
  -c, --has-dependency     Filter tasks that have dependencies
  -h, --help               help for list
  -e, --hide-extra         Hide extra fields (labels, priority, assigned)
  -j, --json               Print JSON output
  -l, --labels strings     Filter tasks by labels
      --limit int          Maximum number of tasks to return (0 means no limit)
  -m, --markdown           print markdown table
      --offset int         Number of tasks to skip from the beginning
  -p, --parent string      Filter tasks by parent ID
      --priority string    Filter tasks by priority
  -r, --reverse            Reverse the order of tasks
      --sort string        Sort tasks by comma-separated fields (id, title, status, priority, created, updated)
  -s, --status strings     Filter tasks by status
  -u, --unassigned         Filter tasks that have no one assigned
```

### Options inherited from parent commands

```
      --auto-commit         Auto-committing changes to git repository
      --folder string       Directory for backlog tasks (default ".backlog")
      --log-file string     Log file path (defaults to stderr)
      --log-format string   Log format (json, text) (default "text")
      --log-level string    Log level (debug, info, warn, error) (default "info")
      --max-limit int       Maximum limit for pagination (default 1000)
      --page-size int       Default page size for pagination (default 25)
```

### SEE ALSO

* [backlog](backlog.md)	 - backlog is a git-native, markdown-based task manager


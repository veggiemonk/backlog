## backlog search

Search tasks by content

### Synopsis

Search for tasks containing the specified query string.

```
backlog search <query> [flags]
```

### Examples

```
# Basic Search

# Search for tasks containing "login" in any field
backlog backlog search "login"

# Search by Bug

# Search for tasks containing "bug"
backlog backlog search "bug"

# Search Assigned Tasks

# Search for tasks assigned to a specific person
backlog backlog search "@john"

# Search with Label

# Search for tasks with specific labels
backlog backlog search "frontend"

# Search in Acceptance Criteria

# Search in acceptance criteria
backlog backlog search "validation"

# Search with Markdown Output

backlog backlog search "api" --markdown

# Search with JSON Output

backlog backlog search "api" --json

# Search with Status Filter

backlog backlog search "user" --status "todo"

# Search with Pagination

# Show first 5 search results
backlog backlog search "api" --limit "5"

# Search with Offset

# Show 3 results starting from 6th match
backlog backlog search "bug" --limit "3" --offset "5"
```

### Options

```
  -a, --assigned strings   Filter tasks by assigned names
  -d, --depended-on        Filter tasks that are depended on by other tasks
  -c, --has-dependency     Include tasks that have dependencies
  -h, --help               help for search
  -e, --hide-extra         Hide extra fields (labels, priority, assigned)
  -j, --json               Print JSON output
  -l, --labels strings     Filter tasks by labels
      --limit int          Maximum number of tasks to return (0 means no limit)
  -m, --markdown           Print markdown table
      --offset int         Number of tasks to skip from the beginning
  -p, --parent string      Filter tasks by parent ID
  -r, --reverse            Reverse the order of tasks
      --sort string        Sort tasks by comma-separated fields (id, title, status, priority, created, updated)
  -s, --status strings     Filter tasks by status
  -u, --unassigned         List tasks that have no one assigned
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


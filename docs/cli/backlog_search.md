## backlog search

Search tasks by content

### Synopsis

Search for tasks containing the specified query string.

```
backlog search <query> [flags]
```

### Examples

```

# Search for tasks containing "login" in any field
backlog search "login"

# Search for tasks containing "bug" 
backlog search "bug"

# Search for tasks assigned to a specific person
backlog search "@john"

# Search for tasks with specific labels
backlog search "frontend"

# Search in acceptance criteria
backlog search "validation"

# Search with markdown output
backlog search "api" --markdown

# Search with JSON output
backlog search "api" --json

# Search with additional columns displayed
backlog search "user" --labels --priority --assigned

# Search with pagination
backlog search "api" --limit 5                  # Show first 5 search results
backlog search "bug" --limit 3 --offset 5       # Show 3 results starting from 6th match
backlog search "feature" --status "todo" --limit 10  # Show first 10 "todo" feature results
	
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
      --limit int          Maximum number of tasks to return (0 means no limit) (default 25)
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
      --auto-commit         Auto-committing changes to git repository (default true)
      --folder string       Directory for backlog tasks (default ".backlog")
      --log-file string     Log file path (defaults to stderr)
      --log-format string   Log format (json, text) (default "text")
      --log-level string    Log level (debug, info, warn, error) (default "info")
      --max-limit int       Maximum limit for pagination (default 50)
      --page-size int       Default page size for pagination (default 25)
```

### SEE ALSO

* [backlog](backlog.md)	 - Backlog is a git-native, markdown-based task manager


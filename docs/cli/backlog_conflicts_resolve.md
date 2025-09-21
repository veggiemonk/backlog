## backlog conflicts resolve

Resolve task ID conflicts

### Synopsis

Create and execute a resolution plan for detected conflicts.

Resolution strategies:
- chronological: Keep older tasks, renumber newer ones (default)
- auto: Automatically renumber conflicting IDs
- manual: Create plan requiring manual intervention

Examples:
  backlog conflicts resolve                            # Resolve using chronological strategy
  backlog conflicts resolve --strategy=auto           # Resolve using auto strategy
  backlog conflicts resolve --dry-run                 # Show what would be changed
  backlog conflicts resolve --strategy=manual         # Create manual resolution plan

```
backlog conflicts resolve [flags]
```

### Options

```
      --dry-run           Show what would be changed without making changes
  -h, --help              help for resolve
      --strategy string   Resolution strategy (chronological|auto|manual) (default "chronological")
```

### Options inherited from parent commands

```
      --auto-commit         Auto-committing changes to git repository (default true)
      --folder string       Directory for backlog tasks (default ".backlog")
      --log-file string     Log file path (defaults to stderr)
      --log-format string   Log format (json, text) (default "text")
      --log-level string    Log level (debug, info, warn, error) (default "info")
      --max-limit int       Maximum limit for pagination (default 1000)
      --page-size int       Default page size for pagination (default 25)
```

### SEE ALSO

* [backlog conflicts](backlog_conflicts.md)	 - Manage task ID conflicts


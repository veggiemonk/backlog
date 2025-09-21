## backlog conflicts detect

Detect task ID conflicts

### Synopsis

Scan all task files and identify ID conflicts including:
- Duplicate IDs (same ID in multiple files)
- Orphaned children (tasks with non-existent parents)
- Invalid hierarchy (parent-child ID mismatch)

Examples:
  backlog conflicts detect                 # Detect conflicts in text format
  backlog conflicts detect --json          # Detect conflicts in JSON format

```
backlog conflicts detect [flags]
```

### Options

```
  -h, --help   help for detect
  -j, --json   Output in JSON format
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


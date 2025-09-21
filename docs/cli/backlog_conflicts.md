## backlog conflicts

Manage task ID conflicts

### Synopsis

Detect and resolve task ID conflicts that can occur when creating tasks
in separate Git branches. Conflicts arise when multiple branches generate the same task IDs.

This command provides conflict detection and resolution capabilities to maintain
task ID uniqueness and data integrity.

### Options

```
  -h, --help   help for conflicts
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

* [backlog](backlog.md)	 - Backlog is a git-native, markdown-based task manager
* [backlog conflicts detect](backlog_conflicts_detect.md)	 - Detect task ID conflicts
* [backlog conflicts resolve](backlog_conflicts_resolve.md)	 - Resolve task ID conflicts


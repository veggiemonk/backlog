## backlog doctor

Diagnose and fix task ID conflicts

### Synopsis

Diagnose and fix task ID conflicts that can occur when creating tasks
in separate Git branches. Conflicts arise when multiple branches generate the same task IDs.

This command provides conflict detection and resolution capabilities to maintain
task ID uniqueness and data integrity.

Conflict types detected:
- Duplicate IDs (same ID in multiple files)
- Orphaned children (tasks with non-existent parents)
- Invalid hierarchy (parent-child ID mismatch)

Examples:
  backlog doctor                    # Detect conflicts in text format
  backlog doctor --json             # Detect conflicts in JSON format
  backlog doctor --fix              # Detect and automatically fix conflicts
  backlog doctor --fix --dry-run    # Show what would be fixed without making changes
  backlog doctor --fix --strategy=auto    # Use auto-renumbering strategy

```
backlog doctor [flags]
```

### Options

```
      --dry-run           Show what would be changed without making changes (use with --fix)
      --fix               Automatically fix detected conflicts
  -h, --help              help for doctor
  -j, --json              Output in JSON format
      --strategy string   Resolution strategy when using --fix (chronological|auto|manual) (default "chronological")
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


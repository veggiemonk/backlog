## backlog archive

Archive a task

### Synopsis

Archives a task, moving it to the archived directory.

```
backlog archive <task-id> [flags]
```

### Examples

```
# Archive a Task

# Archive task T01, moving it to the archived directory
backlog backlog archive "T01"

# Archive by Partial ID

# Archive task using partial ID
backlog backlog archive "1"
```

### Options

```
  -h, --help   help for archive
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


## backlog version

Print the version information

```
backlog version [flags]
```

### Examples

```
# Print Version Information

# Print the version information
backlog backlog version
# Expected: Backlog version:
Revision: 7c989dabd2c61a063a23788c18eb39eca408f6a7
Version: v0.0.2-0.20250907193624-7c989dabd2c6
BuildTime: 2025-09-07T19:36:24Z
Dirty: false
```

### Options

```
  -h, --help   help for version
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


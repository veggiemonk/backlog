## backlog view

View a task by providing its ID

### Synopsis

View a task by providing its ID. You can output in markdown or JSON format.

```
backlog view <id> [flags]
```

### Examples

```
# View Task in Markdown

# View task T01 in markdown format
backlog backlog view "T01"

# View Task in JSON

# View task T01 in JSON format
backlog backlog view "T01" --json
```

### Options

```
  -h, --help   help for view
  -j, --json   Print JSON output
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


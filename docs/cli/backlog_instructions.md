## backlog instructions

instructions for agents to learn to use backlog

### Synopsis

Instructions for agents to learn to use backlog by including them into a prompt.

```
backlog instructions [flags]
```

### Examples

```
# Output Instructions

# outputs the instructions
backlog backlog instructions

# Save Instructions to File

# add instructions to agent base prompt
backlog backlog instructions >> AGENTS.md
```

### Options

```
  -h, --help   help for instructions
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


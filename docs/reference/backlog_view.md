---
layout: page
title: backlog view
---

## backlog view

View a task by providing its ID

### Synopsis

View a task by providing its ID. You can output in markdown or JSON format.

Examples:
  backlog view T01           # View task T01 in markdown format
  backlog view T01 --json    # View task T01 in JSON format
  backlog view T01 -j        # View task T01 in JSON format (short flag)

```
backlog view <id> [flags]
```


### Options

#### Environment Variables

```
	(name)		(default)
	FOLDER		.backlog
	AUTO-COMMIT	false
	LOG-LEVEL	info
	LOG-FORMAT	text
	LOG-FILE	
```

#### Flags


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
```

### SEE ALSO

* [backlog](backlog.md)	 - Backlog is a git-native, markdown-based task manager


---
layout: page
title: backlog archive
---

## backlog archive

Archive a task

### Synopsis

Archives a task, moving it to the archived directory.

```
backlog archive <task-id> [flags]
```

### Examples

```

backlog archive 10  # archive task 10

```


### Options

#### Environment Variables

```
	(name)				(default)
	BACKLOG_AUTO_COMMIT	false
	BACKLOG_FOLDER		.backlog
	BACKLOG_LOG_FILE	
	BACKLOG_LOG_FORMAT	text
	BACKLOG_LOG_LEVEL	info
```

#### Flags


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
```

### SEE ALSO

* [backlog](backlog.md)	 - Backlog is a git-native, markdown-based task manager


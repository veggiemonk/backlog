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


### Options

#### Environment Variables

```
	(name)		(default)
	AUTO-COMMIT	false
	FOLDER		.backlog
	LOG-FILE	
	LOG-FORMAT	text
	LOG-LEVEL	info
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


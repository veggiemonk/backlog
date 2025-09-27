---
layout: page
title: backlog version
---

## backlog version

Print the version information

```
backlog version [flags]
```

### Examples

```

backlog version # Print the version information

# Example output:
#
# Backlog version:
# Revision: 7c989dabd2c61a063a23788c18eb39eca408f6a7
# Version: v0.0.2-0.20250907193624-7c989dabd2c6
# BuildTime: 2025-09-07T19:36:24Z
# Dirty: false

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
  -h, --help   help for version
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


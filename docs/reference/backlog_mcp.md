---
layout: page
title: backlog mcp
---

## backlog mcp

Start the MCP server

### Synopsis

Starts an MCP server to provide programmatic access to backlog tasks.

```
backlog mcp [flags]
```

### Examples

```

backlog mcp --http             # Start the MCP server using HTTP transport on default port 8106
backlog mcp --http --port 4321 # Start the MCP server using HTTP transport on port 4321
backlog mcp                    # Start the MCP server using stdio transport

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
  -h, --help       help for mcp
      --http       Use HTTP transport instead of stdio
      --port int   Port for the MCP server (HTTP transport) (default 8106)
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


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

```
  -h, --help       help for mcp
      --http       Use HTTP transport instead of stdio
      --port int   Port for the MCP server (HTTP transport) (default 8106)
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


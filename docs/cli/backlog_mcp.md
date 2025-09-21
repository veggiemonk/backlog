## backlog mcp

Start the MCP server

### Synopsis

Starts an MCP server to provide programmatic access to backlog tasks.

```
backlog mcp [flags]
```

### Examples

```
# Start MCP Server with HTTP

# Start the MCP server using HTTP transport on default port 8106
backlog backlog mcp --http

# Start MCP Server with Custom Port

# Start the MCP server using HTTP transport on port 4321
backlog backlog mcp --http --port "4321"

# Start MCP Server with Stdio

# Start the MCP server using stdio transport
backlog backlog mcp
```

### Options

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
      --max-limit int       Maximum limit for pagination (default 1000)
      --page-size int       Default page size for pagination (default 25)
```

### SEE ALSO

* [backlog](backlog.md)	 - backlog is a git-native, markdown-based task manager


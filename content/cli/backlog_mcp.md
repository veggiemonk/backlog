---
title: "backlog mcp"
description: "Starts an MCP server to provide programmatic access to backlog tasks."
weight: 8
---

## backlog mcp

Start the MCP server.

### Synopsis

Starts a Model-Context-Protocol (MCP) server to provide programmatic, tool-based access to your backlog. This is essential for integrating with AI agents.

The server can run in two modes:
- **stdio**: Communicates over standard input/output. This is the default.
- **http**: Communicates over HTTP.

```
backlog mcp [flags]
```

### Examples

```bash
# Start the MCP server using stdio transport (default)
backlog mcp

# Start the MCP server using HTTP transport on the default port (8106)
backlog mcp --http

# Start the MCP server using HTTP on a custom port
backlog mcp --http --port 4321
```

### Options

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | | Help for mcp. |
| `--http` | | Use HTTP transport instead of stdio. |
| `--port` | `8106` | Port for the MCP server (HTTP transport). |

{{< expand "Inherited Options" >}}
| Flag | Default | Description |
| --- | --- | --- |
| `--auto-commit` | `true` | Auto-committing changes to git repository. |
| `--folder` | `.backlog` | Directory for backlog tasks. |
{{< /expand >}}

### SEE ALSO

- [**`backlog`**]({{< relref "backlog" >}}) - The main backlog command.
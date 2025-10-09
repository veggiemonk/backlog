---
layout: page
title: backlog mcp
---

# NAME

mcp - Start the MCP server

# SYNOPSIS

mcp

```
[--http]
[--port]=[value]
```

# DESCRIPTION

Starts an MCP server to provide programmatic access to backlog tasks.

Examples:

backlog mcp --http             # Start the MCP server using HTTP transport on default port 8106
backlog mcp --http --port 4321 # Start the MCP server using HTTP transport on port 4321
backlog mcp                    # Start the MCP server using stdio transport


**Usage**:

```
mcp [GLOBAL OPTIONS] [command [COMMAND OPTIONS]] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--http**: Use HTTP transport instead of stdio

**--port**="": Port for the MCP server (HTTP transport) (default: 8106)


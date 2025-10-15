# AI Agent Integration



Backlog is designed for seamless integration with AI agents through two primary methods: the Model Context Protocol (MCP) server and direct CLI calls.

### 1. MCP Server (Recommended for Tool-based Agents)

Backlog includes an MCP server that exposes task management capabilities as structured tools. This is the recommended approach for AI agents that support tool use, as it provides a clear, typed API for interaction.

#### Available Tools

- `task_create`: Create new tasks with full metadata.
- `task_list`: List and filter tasks.
- `task_view`: Get detailed information for a specific task.
- `task_edit`: Update existing tasks.
- `task_archive`: Archive tasks so they are not displayed in lists but remain in the repository.

#### Usage

To make these tools available to an agent:

1. Instruct your AI tools with `backlog instructions --mcp`
2. Start the MCP server

```bash
# Instruct agent how to interact with backlog.
backlog instructions --mode mcp >> AGENTS.md

# Start MCP server for AI integration (HTTP transport)
backlog mcp --http --port 8106

# Or start the server using STDIO transport
backlog mcp
```

AI agents configured to use this server can then reliably perform task operations like creating tasks, breaking down large projects, and updating progress.

### 2. Direct CLI Usage

For agents that are proficient with shell commands, or for workflows where token usage is a concern, `backlog` can be used directly as a command-line tool. The agent can construct and execute `backlog` commands just as a human developer would.

This method bypasses the MCP server, which can reduce overhead and token consumption associated with JSON-based tool calls.

#### Example CLI Usage by an Agent

An agent could be instructed to use backlog with `backlog instructions --cli`

```bash
# Instruct agent how to interact with backlog.
backlog instructions --mode cli >> AGENTS.md

# Agent creates a new task
backlog create "Refactor the authentication module" -d "The current module is hard to test."

# Agent lists tasks to see what's next
backlog list --status "todo" --priority "high"

# Agent marks a task as complete
backlog edit T05 --status "done"
```

### Choosing the Right Method

- **Use MCP Server when**: Your agent has robust tool-use capabilities and you prefer structured, predictable interactions.
- **Use Direct CLI when**: Your agent is skilled at generating shell commands, you want to minimize token usage, or you need the full flexibility of the CLI flags not exposed via MCP.

##

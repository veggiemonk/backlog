---
layout: home
title: Home
nav_order: 1
---

# Backlog

> Manage the backlog of your project with just markdown in git.
> The goal is to provide a frictionless collaboration between AI agents and developers

Backlog is a zero-configuration task manager written in Go where tasks live inside a Git repository. It leverages plain Markdown files for task storage and a comprehensive command-line interface (CLI) for interaction. This design makes it exceptionally well-suited for AI agents thanks to its MCP (Model Context Protocol) integration.

The system is designed to be offline-first and completely portable, as the entire project state is contained within the Git repository.

## Features

- **Task Management**: Create, edit, list, and view tasks with rich metadata
- **Hierarchical Structure**: Support for parent-child-grandchild task relationships (T01 → T01.01 → T01.01.01)
- **Search & Filter**: Find tasks by content, status, parent relationships, and labels
- **AI-Friendly**: MCP server integration for seamless AI agent collaboration
- **Git Integration**: Tasks are stored as Markdown files with automatic Git commits
- **Offline-First**: Works completely offline with local Git repository storage
- **Portable**: Entire project state contained within the Git repository
- **Zero Configuration**: No setup files or databases required

## Quick Start

### Installation

```bash
# Build from source
git clone https://github.com/veggiemonk/backlog
cd backlog
go build .

# Or install directly
go install github.com/veggiemonk/backlog@latest
```

You can also download the binary from the [release page](https://github.com/veggiemonk/backlog/releases).

### Basic Usage

```bash
# Create a simple task
backlog create "Fix the login button styling"

# Create task with metadata
backlog create "Implement password reset" \
  -d "Users should be able to request a password reset link via email" \
  -a "alex" -l "backend,security" --priority "high"

# List all tasks
backlog list

# View specific task
backlog view T01
```

## AI Agent Integration

Start the MCP server to enable AI agent interaction:

```bash
# Start MCP server for AI agents
backlog mcp --http --port 8106  # HTTP transport
backlog mcp                     # STDIO transport (default)
```

AI agents can then:
- Create tasks from conversation context
- Break down large tasks into subtasks
- Update task status and assignments
- Search and analyze task patterns

## Why Backlog?

**For Humans**: The context window deteriorates rapidly on large-scale projects. A workaround is to ask AI agents to "make a plan for X, write it to a markdown file and keep it updated with the ongoing tasks". This technique has worked incredibly well, making refactoring and other significant code changes more resilient to failures.

**For AI Agents**: Backlog provides a structured MCP server that AI tools can understand and trust to handle task management details, providing a much better experience when using AI tools for complex tasks.

## Learn More

- [Getting Started Guide](getting-started.html) - Detailed setup and usage
- [CLI Reference](cli/) - Complete command documentation
- [AI Integration](ai-integration.html) - MCP server setup for AI agents
- [Development](development/) - Architecture and contribution guide